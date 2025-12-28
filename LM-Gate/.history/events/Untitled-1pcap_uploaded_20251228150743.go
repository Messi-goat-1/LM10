package events

type PCAPUploaded struct {
	FileID   string
	FilePath string
}

func (e PCAPUploaded) Name() string {
	return "pcap.uploaded"
}
