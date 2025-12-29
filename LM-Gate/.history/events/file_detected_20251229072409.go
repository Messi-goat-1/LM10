package events

// FileDetectedPayload contains technical details about a detected file.
//
// NOTE: This payload is used when the system discovers a new file
// and wants to notify other parts of the system.
// TODO: Validate fields (FileID, FileName, SizeBytes) before publishing the event.
type FileDetectedPayload struct {
	// FileID is a unique identifier for the file.
	FileID string `json:"file_id"`

	// FileName is the original name of the file.
	FileName string `json:"file_name"`

	// SizeBytes is the file size in bytes.
	SizeBytes int64 `json:"size_bytes"`

	// FileType describes the type/format (e.g. "pcap", "csv", "json").
	// FIXME: FileType may be unreliable if based only on extension.
	FileType string `json:"file_type"`

	// Checksum is used to verify integrity and detect duplicates.
	// NOTE: Useful for deduplication.
	Checksum string `json:"checksum"`

	// StorageHint suggests where/how the file should be stored (optional).
	// NOTE: Can be used to choose "local", "s3", "fast-disk", etc.
	// TODO: Define allowed values for StorageHint.
	StorageHint string `json:"storage_hint"`
}

// FileDetectedEvent is the main event wrapper.
//
// NOTE: Wraps FileDetectedPayload and includes a timestamp.
// TODO: Use time.Time instead of string for Timestamp.
type FileDetectedEvent struct {
	// Payload holds the file details.
	Payload FileDetectedPayload `json:"payload"`

	// Timestamp records when the file was detected.
	Timestamp string `json:"timestamp"`
}
