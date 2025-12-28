package handlers

import (
	"LM-Gate/events"
	"encoding/json"
	"fmt"
)

type EventDispatcher struct {
	fileHandler  *FileDetectedHandler
	chunkHandler *FileChunkHandler
	pcapHandler  *PCAPAnalyzeHandler // 👈 جديد
}

func NewEventDispatcher(fh *FileDetectedHandler, ch *FileChunkHandler) *EventDispatcher {
	return &EventDispatcher{
		fileHandler:  fh,
		chunkHandler: ch,
	}
}

// Dispatch هي الدالة التي تقرر أي Handler ستستدعي
func (d *EventDispatcher) Dispatch(routingKey string, data []byte) error {
	switch routingKey {
	case "file.detected":
		var event events.FileDetectedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return err
		}
		d.fileHandler.Handle(event)

	case "file.chunk":
		var event events.FileChunkEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return err
		}
		d.chunkHandler.Handle(event)

	default:
		return fmt.Errorf("unknown routing key: %s", routingKey)
	}
	return nil
}
