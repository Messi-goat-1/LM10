package events

// PCAPAnalyzeEvent represents a request to analyze a PCAP file.
//
// NOTE: This event is triggered after a PCAP file is fully uploaded
// and ready for analysis.
// TODO: Add validation to ensure FilePath exists before processing.
type PCAPAnalyzeEvent struct {
	// FileID is the unique identifier of the PCAP file.
	FileID string `json:"file_id"`

	// FilePath is the absolute or relative path to the PCAP file on disk.
	// FIXME: FilePath may become invalid if the file is moved or deleted.
	FilePath string `json:"file_path"`
}

// Name returns the event name.
//
// NOTE: This is used by the event system to route or identify the event type.
// TODO: Add versioning support for this event name.
func (e PCAPAnalyzeEvent) Name() string {
	return "pcap.analyze"
}
