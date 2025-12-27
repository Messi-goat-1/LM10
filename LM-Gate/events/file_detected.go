package events

type FileDetectedPayload struct {
	FileID      string `json:"file_id"`
	FileName    string `json:"file_name"`
	SizeBytes   int64  `json:"size_bytes"`
	FileType    string `json:"file_type"`
	Checksum    string `json:"checksum"`
	StorageHint string `json:"storage_hint"`
}

// FileDetectedEvent = إشعار بوجود ملف
type FileDetectedEvent struct {
	FileName string
	Size     int64
}
