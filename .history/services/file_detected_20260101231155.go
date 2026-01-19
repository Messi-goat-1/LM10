package services

import (
	"LM-Gate/events"
	"fmt"
	"log"
	"time"
)

// FileService is responsible for executing business logic
// related to detected files.
//
// NOTE: This service focuses on high-level file processing logic.
// It does NOT handle chunk storage or reassembly.
// TODO: Connect this service to a database layer.
type FileService struct{}

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

// نكتفي بنسخة واحدة فقط من الدالة تستقبل الـ Payload
func (s *FileService) OnFileDetected(payload events.FileDetectedPayload) {
	// 1. التحقق من صحة البيانات (Validation) كما هو مقترح في الـ FIXME
	if payload.FileID == "" || payload.SizeBytes <= 0 {
		log.Printf("[WARN] Invalid FileDetectedPayload: %+v", payload)
		return
	}

	fmt.Println("🚀 [SERVICE] Started processing new file")
	fmt.Printf("   ID: %s\n", payload.FileID)
	fmt.Printf("   Name: %s\n", payload.FileName)
	fmt.Printf("   Size: %d bytes\n", payload.SizeBytes)
	fmt.Printf("   Type: %s\n", payload.FileType)
	fmt.Printf("   Checksum: %s\n", payload.Checksum)
	fmt.Printf("   Processing time: %s\n", time.Now().Format(time.RFC3339))

	// هنا تضع خطوات الـ TODO: التحقق من التكرار، التخزين في القاعدة، ونقل الملف
}
