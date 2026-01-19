package services

import (
	"LM-Gate/events"
	"log"
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

func NewManager(fs *FileService) *Manager {
	return &Manager{
		fileService: fs,
	}
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
func (s *FileService) OnFileDetected(payload events.FileDetectedPayload) {
	if payload.FileID == "" || payload.SizeBytes <= 0 {
		log.Printf("[WARN] Invalid FileDetectedPayload: %+v", payload)
		return
	}

	// business logic هنا فقط
}

// OnFileDetected handles a newly detected file event.
//
// NOTE: This function is triggered after a file is detected
// and its metadata is available.
// FIXME: No validation is performed on input fields.
