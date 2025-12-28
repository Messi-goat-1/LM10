package lmgate

import (
	"fmt"
	"os"
	"path/filepath"
	// "your_project_path/analysis" // قم بفك التعليق بعد إنشاء مجلد التحليل
)

// OnMessage هي نقطة الدخول الرئيسية لمعالجة الرسائل القادمة
func OnMessage(msg ChunkMessage) error {
	if err := ValidateMessage(msg); err != nil {
		return err
	}

	if msg.IsEOF {
		// 1. تجميع القطع والحصول على مسار الملف النهائي
		filePath, err := AssembleFile(msg.FileID)
		if err != nil {
			return err
		}

		// 2. استدعاء دالة المعالجة التي ستربطنا بمجلد analysis
		if err := ProcessFile(msg.FileID, filePath); err != nil {
			return err
		}

		// 3. تنظيف المجلد المؤقت بعد التجميع
		Cleanup(msg.FileID)
		return nil
	}

	// تخزين القطعة الحالية على القرص
	return StoreChunk(msg)
}

// ValidateMessage للتحقق من صحة الرسالة
func ValidateMessage(msg ChunkMessage) error {
	if msg.FileID == "" {
		return fmt.Errorf("invalid FileID")
	}
	if !msg.IsEOF && len(msg.Data) == 0 {
		return fmt.Errorf("invalid Data")
	}
	return nil
}

// StoreChunk يحفظ القطعة في مجلد مؤقت على القرص
func StoreChunk(msg ChunkMessage) error {
	tempDir := filepath.Join("temp_chunks", msg.FileID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return err
	}

	chunkPath := filepath.Join(tempDir, fmt.Sprintf("part_%d", msg.ChunkID))
	return os.WriteFile(chunkPath, msg.Data, 0644)
}

// AssembleFile يجمع القطع ويعيد المسار المادي للملف
func AssembleFile(fileID string) (string, error) {
	tempDir := filepath.Join("temp_chunks", fileID)
	finalDir := "uploads"
	os.MkdirAll(finalDir, 0755)

	finalPath := filepath.Join(finalDir, fileID+".pcap")

	out, err := os.Create(finalPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	for i := 0; ; i++ {
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("part_%d", i))
		data, err := os.ReadFile(chunkPath)
		if err != nil {
			break // توقف عند انتهاء القطع
		}
		out.Write(data)
	}

	return finalPath, nil
}

// ProcessFile الجسر الذي يربط السيرفر بمجلد التحليل
func ProcessFile(fileID string, filePath string) error {
	fmt.Printf("🚀 [SERVER] تم التجميع. بدء التحليل للملف: %s\n", filePath)

	// هنا سيتم استدعاء الدالة من المجلد الجديد مستقبلاً
	// err := analysis.AnalyzePCAP(fileID, filePath)
	// return err

	return nil
}

// Cleanup يحذف المجلد المؤقت للملف من القرص
func Cleanup(fileID string) {
	tempDir := filepath.Join("temp_chunks", fileID)
	os.RemoveAll(tempDir)
	fmt.Printf("🧹 [CLEANUP] تم حذف قطع الملف: %s\n", fileID)
}
