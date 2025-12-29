package handlers

import (
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

// NewEventDispatcher creates a new dispatcher and registers handlers.
func NewEventDispatcher(
	fileHandler *FileDetectedHandler,
	chunkHandler *FileChunkHandler,
	pcapHandler *PCAPAnalyzeHandler,
) *EventDispatcher {

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
