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
