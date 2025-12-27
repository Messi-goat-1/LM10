package handlers

import (
	"LM-Gate/events"
	"LM-Gate/services"
)

// FileDetectedHandler يربط event مع service
type FileDetectedHandler struct {
	fileService *services.FileService
}

// constructor
func NewFileDetectedHandler(fs *services.FileService) *FileDetectedHandler {
	return &FileDetectedHandler{
		fileService: fs,
	}
}

// Handle تُستدعى عند وصول الحدث

func (h *FileDetectedHandler) Handle(event events.FileDetectedEvent) {
	h.fileService.OnFileDetected(
		event.FileName,
		event.Size,
	)
}
