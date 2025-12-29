package handlers

import (
	"LM-Gate/events"
	"encoding/json"
	"fmt"
)

// EventDispatcher is responsible for routing incoming events
// to the correct handler based on the routing key.
//
// NOTE: This acts as a central event router.
// TODO: Replace switch-case with a dynamic handler registry.
type EventDispatcher struct {
	// fileHandler handles "file.detected" events
	fileHandler *FileDetectedHandler

	// chunkHandler handles "file.chunk" events
	chunkHandler *FileChunkHandler

	// pcapHandler handles "pcap.analyze" events
	// NOTE: This handler is triggered after upload completion.
	pcapHandler *PCAPAnalyzeHandler
}

// NewEventDispatcher creates a new dispatcher and wires all handlers.
//
// NOTE: All handlers are injected to keep the dispatcher decoupled.
// TODO: Add nil checks for handlers.
func NewEventDispatcher(fh *FileDetectedHandler, ch *FileChunkHandler, ph *PCAPAnalyzeHandler) *EventDispatcher {
	return &EventDispatcher{
		fileHandler:  fh,
		chunkHandler: ch,
		pcapHandler:  ph,
	}
}

// Dispatch routes the incoming event to the correct handler
// based on the routing key.
//
// NOTE: routingKey usually comes from RabbitMQ.
// FIXME: json.Unmarshal errors are currently ignored.
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
		// NOTE: This event triggers PCAP analysis.
		json.Unmarshal(data, &event)
		return d.pcapHandler.Handle(event)

	default:
		return fmt.Errorf("unknown routing key: %s", routingKey)
	}

	return nil
}

/*
RabbitMQ
   ↓ (routingKey + data)
EventDispatcher
   ↓
[ FileDetectedHandler | FileChunkHandler | PCAPAnalyzeHandler ]

*/
