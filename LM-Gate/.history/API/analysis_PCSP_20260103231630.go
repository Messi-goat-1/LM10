package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

// الإعدادات العامة
const (
	MaxPacketsPerChunk = 1000             // عدد الحزم لكل جزء
	OutputDir          = "chunks"         // مجلد تخزين الأجزاء
	CleanupInterval    = 10 * time.Minute // فحص المجلد كل 10 دقائق
	MaxFileAge         = 30 * time.Minute // حذف الملفات التي عمرها أكثر من 30 دقيقة
)

func main() {
	// 1. إنشاء مجلد التخزين إذا لم يكن موجوداً
	os.MkdirAll(OutputDir, os.ModePerm)

	// 2. تشغيل عامل التنظيف الدوري في الخلفية (Goroutine)
	startCleanupWorker(CleanupInterval, MaxFileAge)

	// 3. إعداد سيرفر الـ API
	r := gin.Default()

	r.POST("/split-pcap", handlePcapSplit)

	fmt.Println("🚀 السيرفر يعمل على المنفذ :8080")
	r.Run(":8080")
}

// --- [ الدالات الخاصة بـ API ] ---

func handlePcapSplit(c *gin.Context) {
	fileHeader, err := c.FormFile("pcapfile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "لم يتم العثور على ملف مرفق"})
		return
	}

	// فتح الملف المرفوع
	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "فشل فتح الملف"})
		return
	}
	defer src.Close()

	// معالجة التقسيم
	createdFiles, err := ProcessPcap(src, fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// الرد على المستخدم بنجاح العملية
	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"total_chunks": len(createdFiles),
		"chunks":       createdFiles,
		"note":         "ستحذف هذه الملفات تلقائياً بعد 30 دقيقة",
	})
}

// --- [ منطق معالجة الـ PCAP ] ---

func ProcessPcap(inputFile io.Reader, originalName string) ([]string, error) {
	reader, err := pcapgo.NewReader(inputFile)
	if err != nil {
		return nil, fmt.Errorf("تنسيق ملف PCAP غير صالح")
	}

	var createdFiles []string
	var currentWriter *pcapgo.Writer
	var currentFile *os.File
	packetCount := 0
	chunkID := 0

	for {
		data, ci, err := reader.ReadPacketData()
		if err == io.EOF {
			break
		}

		// إنشاء ملف جديد عند الوصول للحد الأقصى لكل جزء
		if packetCount%MaxPacketsPerChunk == 0 {
			if currentFile != nil {
				currentFile.Close()
			}

			chunkName := fmt.Sprintf("chunk_%d_%s", chunkID, originalName)
			currentFile, currentWriter, err = createNewChunk(chunkName, reader.LinkType())
			if err != nil {
				return nil, err
			}

			createdFiles = append(createdFiles, chunkName)
			chunkID++
		}

		currentWriter.WritePacket(ci, data)
		packetCount++
	}

	if currentFile != nil {
		currentFile.Close()
	}

	return createdFiles, nil
}

func createNewChunk(name string, linkType layers.LinkType) (*os.File, *pcapgo.Writer, error) {
	path := filepath.Join(OutputDir, name)
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}

	writer := pcapgo.NewWriter(f)
	writer.WriteFileHeader(65536, linkType)

	return f, writer, nil
}

// --- [ دالات التنظيف والحذف ] ---

func startCleanupWorker(interval time.Duration, maxAge time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			files, err := os.ReadDir(OutputDir)
			if err != nil {
				continue
			}
			for _, file := range files {
				path := filepath.Join(OutputDir, file.Name())
				info, err := file.Info()
				if err != nil {
					continue
				}
				if time.Since(info.ModTime()) > maxAge {
					os.Remove(path)
					log.Printf("🗑️ تم حذف ملف قديم: %s", file.Name())
				}
			}
		}
	}()
}
