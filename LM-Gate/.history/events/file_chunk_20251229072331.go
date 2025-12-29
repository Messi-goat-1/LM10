package events

// FileChunkPayload contains raw file chunk data.
//
// NOTE: This structure represents a single piece of a file.
// It is sent from the client to the system (via RabbitMQ or similar).
// TODO: Add checksum field to validate chunk integrity.
type FileChunkPayload struct {
	// FileID is the unique identifier of the original file.
	FileID string `json:"file_id"`

	// ChunkIndex represents the order of this chunk.
	// NOTE: Index starts from 0.
	ChunkIndex int `json:"chunk_index"`

	// TotalChunks is the expected total number of chunks.
	// NOTE: Used to detect when the file is fully received.
	TotalChunks int `json:"total_chunks"`

	// Data contains the raw bytes of the file chunk.
	// FIXME: Large chunks may increase memory usage.
	Data []byte `json:"data"`
}

// FileChunkEvent represents the event received from RabbitMQ.
//
// NOTE: This event wraps the chunk payload.
// It also includes a timestamp to record when the chunk arrived.
// TODO: Use time.Time instead of string for Timestamp.
type FileChunkEvent struct {
	// Payload contains the actual chunk data and metadata.
	Payload FileChunkPayload `json:"payload"`

	// Timestamp records when the event was created or received.
	Timestamp string `json:"timestamp"`
}

/*
Client
  ↓
FileChunkEvent
  ↓
RabbitMQ
  ↓
Handler
  ↓
Service (store → check → reassemble)


*/
