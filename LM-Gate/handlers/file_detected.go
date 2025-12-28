package handlers

import (
	"LM-Gate/events"
	"LM-Gate/services"
)

// FileDetectedHandler يربط event مع service
type FileDetectedHandler struct {
	fileService *services.Manager
}

// constructor
func NewFileDetectedHandler(fs *services.Manager) *FileDetectedHandler {
	return &FileDetectedHandler{
		fileService: fs,
	}
}

// Handle تُستدعى عند وصول الحدث

// file_detected.go

func (h *FileDetectedHandler) Handle(event events.FileDetectedEvent) {
	// تمرير كل البيانات الجديدة من الـ Payload إلى الـ Service
	h.fileService.OnFileDetected(
		event.Payload.FileID,
		event.Payload.FileName,
		event.Payload.SizeBytes,
		event.Payload.FileType,
		event.Payload.Checksum,
	)
}
