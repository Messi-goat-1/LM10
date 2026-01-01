package handlers

import (
	"LM-Gate/events"
	"LM-Gate/services"
	"encoding/json"
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
	json.Unmarshal(data, &event)

	h.fileService.OnFileDetected(event.Payload)

	return nil
}

/*
file.detected event
   ↓
FileDetectedHandler
   ↓
services.Manager.OnFileDetected



Publisher
 → RabbitMQ
 → FileDetectedHandler
 → FileService.OnFileDetected



*/
