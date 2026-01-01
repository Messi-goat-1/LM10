package handlers

import (
	"LM-Gate/events"
	"LM-Gate/services"
	"encoding/json"
)

type FileCollectionHandler struct {
	manager *services.Manager
}

func NewFileCollectionHandler(mgr *services.Manager) *FileCollectionHandler {
	return &FileCollectionHandler{manager: mgr}
}

func (h *FileCollectionHandler) Handle(data []byte) error {
	var event events.FileCollectionEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	// إرسال البيانات إلى الـ Manager
	h.manager.OnFileCollection(event.Payload)
	return nil
}
