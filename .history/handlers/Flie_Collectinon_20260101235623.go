package handlers

import (
	"LM-Gate/events"
	"LM-Gate/services"
	"encoding/json"
	"log"
)

type FileCollectionHandler struct {
	manager *services.Manager
}

func NewFileCollectionHandler(mgr *services.Manager) *FileCollectionHandler {
	return &FileCollectionHandler{
		manager: mgr,
	}
}

func (h *FileCollectionHandler) Handle(data []byte) error {
	var event events.FileCollectionEvent

	// فك تشفير JSON الخاص بحدث التجميع
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("[ERROR] Failed to unmarshal collection event: %v", err)
		return err
	}

	// تمرير البيانات للمدير
	h.manager.OnFileCollection(event.Payload)
	return nil
}
