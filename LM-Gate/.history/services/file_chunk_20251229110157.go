package services

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manager is responsible for file management and chunk reassembly.
//
// NOTE: This service handles all heavy logic related to files.
// It replaces the old FileService naming for clarity.
// TODO: Make temp and storage paths configurable.
type Manager struct {
	// tempDir stores temporary chunk files.
	tempDir string

	// storageDir stores final assembled files.
	storageDir string
}

// NewManager creates a new Manager instance.
//
// NOTE: This function initializes directory paths.
// TODO: Ensure directories exist during initialization.
func NewManager() *Manager {
	s := &Manager{
		tempDir:    "./temp_chunks",
		storageDir: "./uploads",
	}
	// Directory creation logic can be added here
	return s
}

// OnFileDetected handles a simple file detected event.
//
// NOTE: This is a lightweight notification handler.
// Business logic can be extended here.
func (s *Manager) OnFileDetected(
	fileID string,
	fileName string,
	size int64,
	fileType string,
	checksum string,
) {
	fmt.Printf(
		"📦 [SERVICE] New file detected: %s (ID: %s, Size: %d)\n",
		fileName,
		fileID,
		size,
	)
	// NOTE: checksum and fileType can be used for validation or deduplication.
}

// OnChunkReceived handles an incoming file chunk.
//
// NOTE: This function stores the chunk on disk.
// FIXME: No validation for chunkIndex range or duplicate chunks.
func (s *Manager) OnChunkReceived(fileID string, chunkIndex int, total int, data []byte) error {
	// Create a temporary directory for the file if it does not exist
	fileDir := filepath.Join(s.tempDir, fileID)
	os.MkdirAll(fileDir, 0755)

	// Write the chunk to disk as part_<index>
	chunkPath := filepath.Join(fileDir, fmt.Sprintf("part_%d", chunkIndex))
	err := os.WriteFile(chunkPath, data, 0644)
	if err != nil {
		return err
	}

	// Check if all chunks have been received
	if s.isComplete(fileDir, total) {
		// Reassemble in background
		go s.reassemble(fileID, total)
	}

	return nil
}

// isComplete checks whether all expected chunks are present.
//
// NOTE: This compares file count with total chunks.
// FIXME: Does not validate chunk order or missing indices.
func (s *Manager) isComplete(dir string, total int) bool {
	for i := 0; i < total; i++ {
		partPath := filepath.Join(dir, fmt.Sprintf("part_%d", i))
		if _, err := os.Stat(partPath); err != nil {
			// part_i غير موجود
			return false
		}
	}
	return true
}

// reassemble merges all chunks into a final file.
//
// NOTE: This is an improved version of AssembleFile from server.go.
// TODO: Add checksum verification after reassembly.
// FIXME: No rollback if reassembly fails midway.
func (s *Manager) reassemble(fileID string, totalChunks int) error {
	finalPath := filepath.Join(s.storageDir, fileID)

	// Create the final output file
	dst, err := os.Create(finalPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Read and merge chunks in correct order
	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(
			s.tempDir,
			fileID,
			fmt.Sprintf("part_%d", i),
		)

		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			return err
		}

		// Write chunk data to final file
		dst.Write(chunkData)
	}

	// Cleanup temporary directory after successful merge
	os.RemoveAll(filepath.Join(s.tempDir, fileID))
	return nil
}

/*
FileChunkEvent
   ↓
FileChunkHandler
   ↓
services.Manager.OnChunkReceived
   ↓
isComplete → reassemble → cleanup

*/
