package handlers

import (
	"LM-Gate/events"
	"yourapp/internal/services"
)

type PCAPUploadedHandler struct {
	PCAPService *services.PCAPService
}

func (h *PCAPUploadedHandler) Handle(e events.Event) error {
	event := e.(events.PCAPUploaded)
	return h.PCAPService.Analyze(event.FileID, event.FilePath)
}
