package events

import (
	"time"
)

// Event represents a generic event inside the system.
//
// NOTE: This is a base event structure used across the platform.
// It follows an Event-Driven Architecture style.
// TODO: Enforce strong typing for Payload instead of using `any`.
// FIXME: No validation exists for Event name or Version.
type Event struct {
	// Event is the name/type of the event (e.g. "file.detected", "chunk.received")
	Event string `json:"event"`

	// Version represents the schema version of the event.
	// NOTE: Useful for backward compatibility.
	Version int `json:"version"`

	// OccurredAt records when the event happened.
	// NOTE: Should always be set when the event is created.
	OccurredAt time.Time `json:"occurred_at"`

	// Source indicates where the event was produced from
	// (e.g. "client", "server", "scanner").
	Source string `json:"source"`

	// Payload contains event-specific data.
	// NOTE: Payload structure depends on the event type.
	// FIXME: Using `any` removes compile-time safety.
	Payload any `json:"payload"`
}
