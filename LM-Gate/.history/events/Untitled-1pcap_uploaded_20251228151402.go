package events

// PCAPAnalyzeEvent يمثل حدث طلب تحليل ملف PCAP
type PCAPAnalyzeEvent struct {
	FileID   string `json:"file_id"`
	FilePath string `json:"file_path"`
}

func (e PCAPAnalyzeEvent) Name() string {
	return "pcap.analyze"
}
