package handlers

import (
	"LM-Gate/events"
	"LM-Gate/services"
)

type FileChunkHandler struct {
	fileService *services.Manager
}

func NewFileChunkHandler(fs *services.Manager) *FileChunkHandler {
	return &FileChunkHandler{
		fileService: fs,
	}
}

func (h *FileChunkHandler) Handle(event events.FileChunkEvent) {
	// الربط بين بيانات الحدث ووظيفة الخدمة
	h.fileService.OnChunkReceived(
		event.Payload.FileID,
		event.Payload.ChunkIndex,
		event.Payload.TotalChunks,
		event.Payload.Data,
	)
}
