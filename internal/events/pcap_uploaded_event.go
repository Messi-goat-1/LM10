package events

type PcapUploadedEvent struct {
	FileName string
	Path     string
	Size     int64
}
