package handlers

import (
	"LM-Gate/events"
	"LM-Gate/services"
)

// FileChunkHandler handles incoming file chunk events.
//
// NOTE: This handler acts only as a bridge between
// the event layer and the service layer.
// It does not contain business logic.
type FileChunkHandler struct {
	// fileService manages file storage and reassembly.
	fileService *services.Manager
}

// NewFileChunkHandler creates a new FileChunkHandler.
//
// NOTE: The service manager is injected to keep layers decoupled.
// TODO: Add nil check for fileService.
func NewFileChunkHandler(fs *services.Manager) *FileChunkHandler {
	return &FileChunkHandler{
		fileService: fs,
	}
}

// Handle receives a FileChunkEvent and forwards its data
// to the service layer.
//
// NOTE: This function does not modify data.
// It simply extracts values from the event payload.
// FIXME: No validation is performed on chunk order or size.
func (h *FileChunkHandler) Handle(event events.FileChunkEvent) {
	// Bridge event data to the service logic
	h.fileService.OnChunkReceived(
		event.Payload.FileID,
		event.Payload.ChunkIndex,
		event.Payload.TotalChunks,
		event.Payload.Data,
	)
}
