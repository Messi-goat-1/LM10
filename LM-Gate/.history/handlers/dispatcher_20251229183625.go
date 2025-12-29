package handlers

import (
	"LM-Gate/events"
	"encoding/json"
	"fmt"
)

// EventHandler defines a generic handler for events.
//
// NOTE: Each handler is responsible for unmarshalling its own event data.
type EventHandler interface {
	Handle(data []byte) error
}

// EventDispatcher routes incoming events to the correct handler
// based on the routing key.
//
// NOTE: Dispatcher does not know event structures.
// It only performs routing + delegation.
type EventDispatcher struct {
	handlers map[string]EventHandler
}

func (h *FileChunkHandler) Handle(data []byte) error {
	var event events.FileChunkEvent

	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	return h.handleEvent(event)
}

// NewEventDispatcher creates a new dispatcher and registers handlers.
func NewEventDispatcher(fileHandler *FileDetectedHandler, chunkHandler *FileChunkHandler, pcapHandler *PCAPAnalyzeHandler) *EventDispatcher {

	return &EventDispatcher{
		handlers: map[string]EventHandler{
			"file.detected": fileHandler,
			"file.chunk":    chunkHandler,
			"pcap.analyze":  pcapHandler,
		},
	}
}

// Dispatch routes the incoming event to the correct handler.
//
// NOTE: routingKey usually comes from RabbitMQ.
func (d *EventDispatcher) Dispatch(routingKey string, data []byte) error {
	handler, ok := d.handlers[routingKey]
	if !ok {
		return fmt.Errorf("unknown routing key: %s", routingKey)
	}

	return handler.Handle(data)
}

/*
RabbitMQ
   ↓ (routingKey + data)
EventDispatcher
   ↓
[ FileDetectedHandler | FileChunkHandler | PCAPAnalyzeHandler ]

*/
