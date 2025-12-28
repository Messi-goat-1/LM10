package lmgate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// OnMessage is the server entry point for handling incoming chunk messages.
func OnMessage(msg ChunkMessage) error {
	if err := ValidateMessage(msg); err != nil {
		return err
	}

	if msg.IsEOF {
		if !IsFileComplete(msg.FileID) {
			return ErrMissingChunk
		}

		data, err := AssembleFile(msg.FileID)
		if err != nil {
			return err
		}

		if err := ProcessFile(msg.FileID, data); err != nil {
			return err
		}

		Cleanup(msg.FileID)
		return nil
	}

	return StoreChunk(msg)
}

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
	// 1. تحديد مجلد مؤقت (يمكنك تسميته temp_chunks)
	tempDir := fmt.Sprintf("temp_chunks/%s", msg.FileID)

	// 2. إنشاء المجلد إذا لم يكن موجوداً
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return err
	}

	// 3. تحديد مسار القطعة بناءً على الـ ChunkID
	chunkPath := filepath.Join(tempDir, fmt.Sprintf("part_%d", msg.ChunkID))

	// 4. كتابة البيانات مباشرة للملف (هذا يفرغ الذاكرة فوراً)
	return os.WriteFile(chunkPath, msg.Data, 0644)
}

// IsFileComplete checks whether all chunks for a file are present.
func IsFileComplete(fileID string) bool {
	tempDir := fmt.Sprintf("temp_chunks/%s", fileID)
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return false
	}

	// هنا نفترض أنك تعرف العدد الإجمالي للقطع،
	// أو نتحقق من وجود جميع الملفات المتسلسلة
	return len(files) > 0 // هذا فحص بسيط، يمكن تطويره حسب منطق نظامك
}

// AssembleFile reconstructs the original file from stored chunks.
func AssembleFile(fileID string) ([]byte, error) {
	chunks, ok := chunkStore[fileID]
	if !ok {
		return nil, errors.New("file not found")
	}

	var result []byte
	for i := 0; i < len(chunks); i++ {
		result = append(result, chunks[i]...)
	}
	return result, nil
}

// ProcessFile handles the fully assembled file (analysis, parsing, etc).
func ProcessFile(fileID string, data []byte) error {
	fmt.Printf("Processing file %s (size=%d bytes)\n", fileID, len(data))
	return nil
}

// Cleanup removes all stored data related to a file.
func Cleanup(fileID string) {
	delete(chunkStore, fileID)
}

// ---------------
type FakeSender struct{}

func (f *FakeSender) Send(msg ChunkMessage) error {
	return OnMessage(msg)
}
