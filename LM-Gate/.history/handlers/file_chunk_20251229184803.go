package handlers

import (
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

/*
RabbitMQ
   ↓
FileChunkEvent
   ↓
FileChunkHandler
   ↓
services.Manager.OnChunkReceived


*/
