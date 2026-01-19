package services

import (
	"LM-Gate/events"
	"fmt"
	"time"
)

// FileService is responsible for executing business logic
// related to detected files.
//
// NOTE: This service focuses on high-level file processing logic.
// It does NOT handle chunk storage or reassembly.
// TODO: Connect this service to a database layer.
type FileService struct{}
type Manager struct {
	fileService *FileService
}

func (m *Manager) OnFileDetected(payload events.FileDetectedPayload) {
	m.fileService.OnFileDetected(payload)
}

// NewFileService creates a new FileService instance.
//
// NOTE: Currently stateless.
// TODO: Inject dependencies (DB, logger, config) when needed.
func NewFileService() *FileService {
	return &FileService{}
}

// OnFileDetected handles a newly detected file event.
//
// NOTE: This function is triggered after a file is detected
// and its metadata is available.
// FIXME: No validation is performed on input fields.
func (s *FileService) OnFileDetected(fileID string, fileName string, size int64, fileType string, checksum string) {
	fmt.Println("🚀 [SERVICE] Started processing new file")
	fmt.Printf("   ID: %s\n", fileID)
	fmt.Printf("   Name: %s\n", fileName)
	fmt.Printf("   Size: %d bytes\n", size)
	fmt.Printf("   Type: %s\n", fileType)
	fmt.Printf("   Checksum: %s\n", checksum)
	fmt.Printf("   Processing time: %s\n", time.Now().Format(time.RFC3339))

	// TODO:
	// 1. Check for duplicate files using checksum.
	// 2. Store file metadata in the database.
	// 3. Move the file to its final storage location.
}
