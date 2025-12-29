package lmgate

import (
	"LM-Gate/analysis"
	"fmt"
	"os"
	"path/filepath"
)

// OnMessage is the main server entry point.
// It receives chunk messages and decides what to do:
// - store chunks
// - assemble the file
// - process the file
// - clean temporary data
func OnMessage(msg ChunkMessage) error {
	if err := ValidateMessage(msg); err != nil {
		return err
	}

	// If this message is the end-of-file marker
	if msg.IsEOF {
		// 1. Assemble the file and get its path
		filePath, err := AssembleFile(msg.FileID)
		if err != nil {
			return fmt.Errorf("failed to assemble file: %v", err)
		}

		// 2. Start processing the file (analysis step)
		if err := ProcessFile(msg.FileID, filePath); err != nil {
			return err
		}

		// 3. Clean temporary chunks after processing
		Cleanup(msg.FileID)
		return nil
	}

	// If this is a normal chunk, store it on disk
	return StoreChunk(msg)
}

// ValidateMessage performs basic validation on incoming messages.
// It ensures the message has a FileID and valid data.
func ValidateMessage(msg ChunkMessage) error {
	if msg.FileID == "" {
		return ErrInvalidMessage
	}

	if !msg.IsEOF && len(msg.Data) == 0 {
		return ErrInvalidMessage
	}

	return nil
}

// StoreChunk saves a file chunk to disk.
// Each file has its own temporary directory.
// Each chunk is stored as part_<chunk_id>.
func StoreChunk(msg ChunkMessage) error {
	tempDir := fmt.Sprintf("temp_chunks/%s", msg.FileID)

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return err
	}

	chunkPath := filepath.Join(tempDir, fmt.Sprintf("part_%d", msg.ChunkID))

	return os.WriteFile(chunkPath, msg.Data, 0644)
}

// IsFileComplete checks if there are stored chunks for a file.
// This only verifies that chunks exist, not full completion.
func IsFileComplete(fileID string) bool {
	tempDir := fmt.Sprintf("temp_chunks/%s", fileID)
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return false
	}

	return len(files) > 0
}

// AssembleFile rebuilds the original file from stored chunks.
// Chunks are read in order: part_0, part_1, ...
// The final file is saved in the uploads directory.
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

	// Merge chunks sequentially
	for i := 0; ; i++ {
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("part_%d", i))
		data, err := os.ReadFile(chunkPath)
		if err != nil {
			break // Stop when the next chunk is missing
		}
		out.Write(data)
	}

	return finalPath, nil
}

// ProcessFile handles the fully assembled file.
// It passes the file path to the analysis layer.
func ProcessFile(fileID string, filePath string) error {
	return analysis.AnalyzePCAP(fileID, filePath)
}

// Cleanup removes all temporary data related to a file.
func Cleanup(fileID string) {
	os.RemoveAll(filepath.Join("temp_chunks", fileID))
}

// FakeSender is a test helper.
// It simulates sending messages directly to the server.
type FakeSender struct{}

// Send forwards the message directly to OnMessage.
// Used for local testing without networking.
func (f *FakeSender) Send(msg ChunkMessage) error {
	return OnMessage(msg)
}
