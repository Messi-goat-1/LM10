package events

import (
	"time"
)

// Event يمثل أي حدث داخل النظام
type Event struct {
	Event      string    `json:"event"`
	Version    int       `json:"version"`
	OccurredAt time.Time `json:"occurred_at"`
	Source     string    `json:"source"`
	Payload    any       `json:"payload"`
}

/*
func DispatchEvent(eventType string, payload []byte) error {
	switch eventType {
	case "file_detected":
		return handlers.HandleFileDetected(payload)

	default:
		return fmt.Errorf("unknown event type: %s", eventType)
	}
}
*/
