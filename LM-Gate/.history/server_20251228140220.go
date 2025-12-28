package lmgate

import (
	"LM-Gate/analysis"
	"fmt"
	"os"
	"path/filepath"
)

// OnMessage is the server entry point for handling incoming chunk messages.
func OnMessage(msg ChunkMessage) error {
	if err := ValidateMessage(msg); err != nil {
		return err
	}

	// إذا كانت الرسالة هي علامة النهاية
	if msg.IsEOF {
		// 1. تجميع الملف والحصول على مساره المادي
		filePath, err := AssembleFile(msg.FileID)
		if err != nil {
			return fmt.Errorf("failed to assemble file: %v", err)
		}

		// 2. البدء في معالجة الملف (التحليل)
		if err := ProcessFile(msg.FileID, filePath); err != nil {
			return err
		}

		// 3. تنظيف القطع المؤقتة بعد التجميع (يتم تنفيذها هنا فقط)
		Cleanup(msg.FileID)
		return nil
	}

	// إذا كانت قطعة عادية، يتم تخزينها على القرص
	return StoreChunk(msg)
}

// تأكد أن هذا هو القوس الوحيد في النهاية
// ValidateMessage performs basic validation on incoming messages.
func ValidateMessage(msg ChunkMessage) error {
	if msg.FileID == "" {
		return ErrInvalidMessage
	}

	if !msg.IsEOF && len(msg.Data) == 0 {
		return ErrInvalidMessage
	}

	return nil
}

// StoreChunk temporarily stores a file chunk in memory.
// StoreChunk - التعديل الجديد لحفظ البيانات على القرص بدلاً من الذاكرة
func StoreChunk(msg ChunkMessage) error {
	tempDir := fmt.Sprintf("temp_chunks/%s", msg.FileID)

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return err
	}

	chunkPath := filepath.Join(tempDir, fmt.Sprintf("part_%d", msg.ChunkID))

	return os.WriteFile(chunkPath, msg.Data, 0644)
}

// IsFileComplete checks whether all chunks for a file are present.
func IsFileComplete(fileID string) bool {
	tempDir := fmt.Sprintf("temp_chunks/%s", fileID)
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return false
	}

	return len(files) > 0
}

// AssembleFile reconstructs the original file from stored chunks.
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

	// دمج القطع بترتيب متسلسل
	for i := 0; ; i++ {
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("part_%d", i))
		data, err := os.ReadFile(chunkPath)
		if err != nil {
			break // توقف عند عدم العثور على القطعة التالية
		}
		out.Write(data)
	}

	return finalPath, nil
}

// ProcessFile handles the fully assembled file (analysis, parsing, etc).
// تعديل لاستقبال مسار الملف (string) بدلاً من []byte
// في ملف server.go
// داخل server.go
func ProcessFile(fileID string, filePath string) error {
	return analysis.AnalyzePCAP(fileID, filePath)
}

// Cleanup removes all stored data related to a file.
func Cleanup(fileID string) {
	os.RemoveAll(filepath.Join("temp_chunks", fileID)) // حذف المجلد من القرص
}

// ---------------
type FakeSender struct{}

func (f *FakeSender) Send(msg ChunkMessage) error {
	return OnMessage(msg)
}
