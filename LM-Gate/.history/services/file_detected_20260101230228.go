package services

import (
	"LM-Gate/events"
	"fmt"
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

// داخل ملف services
func (s *FileService) OnFileDetected(payload events.FileDetectedPayload) {
	// 1. التحقق من صحة البيانات (Validation)
	if payload.FileID == "" || payload.SizeBytes <= 0 {
		log.Printf("[WARN] Invalid FileDetectedPayload: %+v", payload)
		return
	}

	// 2. منطق العمل (Business Logic)
	fmt.Println("🚀 [SERVICE] Started processing new file")
	fmt.Printf("   ID: %s | Name: %s | Size: %d\n", payload.FileID, payload.FileName, payload.SizeBytes)

	// TODO: إضافة العمليات الأخرى مثل حفظ البيانات في القاعدة
}
