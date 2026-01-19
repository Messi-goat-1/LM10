package events

import "time"

// FileDetectedPayload contains technical details about a detected file.
//
// NOTE: This payload is used when the system discovers a new file
// and wants to notify other parts of the system.
// TODO: Validate fields (FileID, FileName, SizeBytes) before publishing the event.
type FileDetectedPayload struct {
	// FileID is a unique identifier for the file.
	FileID string `json:"file_id"`

	// FileName is the original name of the file.
	FileName string `json:"file_name"`

	// SizeBytes is the file size in bytes.
	SizeBytes int64 `json:"size_bytes"`

	// FileType describes the type/format (e.g. "pcap", "csv", "json").
	// FIXME: FileType may be unreliable if based only on extension.
	FileType string `json:"file_type"`

	// Checksum is used to verify integrity and detect duplicates.
	// NOTE: Useful for deduplication.
	Checksum string `json:"checksum"`

	// StorageHint suggests where/how the file should be stored (optional).
	// NOTE: Can be used to choose "local", "s3", "fast-disk", etc.
	// TODO: Define allowed values for StorageHint.
	StorageHint string `json:"storage_hint"`
}

// FileDetectedEvent is the main event wrapper.
//
// NOTE: Wraps FileDetectedPayload and includes a timestamp.
// TODO: Use time.Time instead of string for Timestamp.
type FileDetectedEvent struct {
	// Payload holds the file details.
	Payload FileDetectedPayload `json:"payload"`

	// Timestamp records when the file was detected.
	Timestamp time.Time `json:"timestamp"`
}
 -----------
 package handlers

import (
	"LM-Gate/events"
	"LM-Gate/services"
	"encoding/json"
	"log"
)

// FileDetectedHandler connects file detected events
// with the service layer.
//
// NOTE: This handler acts as a simple bridge.
// It does not contain business logic.
type FileDetectedHandler struct {
	// fileService handles file-related business logic.
	fileService *services.Manager
}

// NewFileDetectedHandler creates a new FileDetectedHandler.
//
// NOTE: The service manager is injected to keep layers separated.
// TODO: Add nil validation for the service manager.
func NewFileDetectedHandler(fs *services.Manager) *FileDetectedHandler {
	return &FileDetectedHandler{
		fileService: fs,
	}
}

// Handle is called when a FileDetectedEvent is received.
//
// NOTE: This function extracts data from the event payload
// and forwards it directly to the service layer.
// FIXME: No validation is done on payload values.
// Forward all payload data to the service layer
func (h *FileDetectedHandler) Handle(data []byte) error {
	var event events.FileDetectedEvent

	// فحص الخطأ أثناء تحويل JSON
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[ERROR] Failed to unmarshal event: %v", err)
		return err
	}

	// تمرير الـ Payload إلى الـ Manager
	h.fileService.OnFileDetected(event.Payload)
	return nil
}
----------------
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
