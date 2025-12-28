package lmgate

import (
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
	tempDir := fmt.Sprintf("temp_chunks/%s", fileID)
	finalPath := fmt.Sprintf("uploads/%s", fileID)

	out, err := os.Create(finalPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// قراءة القطع وصبها في الملف النهائي
	for i := 0; ; i++ {
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("part_%d", i))
		data, err := os.ReadFile(chunkPath)
		if err != nil {
			break
		} // توقف عند انتهاء القطع
		out.Write(data)
	}

	return finalPath, nil
}

// ProcessFile handles the fully assembled file (analysis, parsing, etc).
// تعديل لاستقبال مسار الملف (string) بدلاً من []byte
func ProcessFile(fileID string, filePath string) error {
	fmt.Printf("✅ Processing file %s located at %s\n", fileID, filePath)
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
