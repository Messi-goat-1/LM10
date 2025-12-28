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

func NewEventDispatcher(fh *FileDetectedHandler, ch *FileChunkHandler, ph *PCAPAnalyzeHandler) *EventDispatcher {
	return &EventDispatcher{
		fileHandler:  fh,
		chunkHandler: ch,
		pcapHandler:  ph,
	}
}
func (d *EventDispatcher) Dispatch(routingKey string, data []byte) error {
	switch routingKey {

	case "file.detected":
		var event events.FileDetectedEvent
		json.Unmarshal(data, &event)
		d.fileHandler.Handle(event)

	case "file.chunk":
		var event events.FileChunkEvent
		json.Unmarshal(data, &event)
		d.chunkHandler.Handle(event)

	case "pcap.analyze":
		var event events.PCAPAnalyzeEvent
		json.Unmarshal(data, &event)
		return d.pcapHandler.Handle(event)

	default:
		return fmt.Errorf("unknown routing key: %s", routingKey)
	}
	return nil
}
