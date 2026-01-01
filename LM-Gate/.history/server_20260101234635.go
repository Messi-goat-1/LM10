package lmgate

import (
	"LM-Gate/analysis"
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// OnMessage is the main server entry point.
// It receives chunk messages and decides what to do.
//
// NOTE: This function acts as the orchestrator for the whole flow.
// TODO: Add structured logging for each step.

// ValidateMessage performs basic validation on incoming messages.
//
// NOTE: This prevents invalid or corrupted messages from being processed.
// TODO: Add validation for ChunkID range.
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
//
// NOTE: Chunks are stored on disk to avoid high memory usage.
// TODO: Add checksum validation per chunk.
// FIXME: No limit on disk usage is enforced.
func StoreChunk(msg ChunkMessage) error {
	tempDir := fmt.Sprintf("temp_chunks/%s", msg.FileID)

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return err
	}

	chunkPath := filepath.Join(tempDir, fmt.Sprintf("part_%d", msg.ChunkID))

	return os.WriteFile(chunkPath, msg.Data, 0644)
}

// IsFileComplete checks if there are stored chunks for a file.
//
// NOTE: This only checks if chunks exist, not if all chunks arrived.
// FIXME: This does not verify the expected number of chunks.
func IsFileComplete(fileID string) bool {
	tempDir := fmt.Sprintf("temp_chunks/%s", fileID)
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return false
	}

	return len(files) > 0
}

// AssembleFile rebuilds the original file from stored chunks.
//
// NOTE: Chunks are read sequentially: part_0, part_1, ...
// FIXME: If a chunk is missing, the loop stops silently.
// TODO: Detect missing chunks and return an error.
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

	// Merge chunks in order
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
//
// NOTE: The file is passed by path, not loaded into memory.
// TODO: Add timeout or context support.
func ProcessFile(fileID string, filePath string) error {
	ctx := context.Background() // أو context.WithTimeout
	return analysis.AnalyzePCAP(ctx, fileID, filePath)
}

// Cleanup removes all temporary data related to a file.
//
// NOTE: This helps prevent disk space leaks.
// TODO: Add retry or safety checks before deletion.
func Cleanup(fileID string) {
	os.RemoveAll(filepath.Join("temp_chunks", fileID))
}

// FakeSender is a helper for testing.
//
// NOTE: Used to simulate message sending without network or MQ.
type FakeSender struct{}

// Send forwards the message directly to OnMessage.
//
// TODO: Add test assertions around Send behavior.
func (f *FakeSender) Send(msg ChunkMessage) error {
	return OnMessage(msg)
}
