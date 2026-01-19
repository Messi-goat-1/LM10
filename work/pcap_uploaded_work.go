package work

import (
	"LM-Gate/internal/events"
	"LM-Gate/internal/infra"
	"LM-Gate/internal/logic"

	"encoding/json"
	"log"
)

// RegisterPcapUploaded
// هذه الدالة تربط الحدث مع RabbitMQ
// work هو المسؤول عن "التشغيل"
func RegisterPcapUploaded(rabbit *infra.RabbitClient) {
	rabbit.ConsumeMessages("pcap_processing_queue", func(body []byte) {

		var event events.PcapUploadedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("❌ invalid pcap_uploaded event: %v", err)
			return
		}

		OnPcapUploaded(event)
	})
}

// OnPcapUploaded
func OnPcapUploaded(event events.PcapUploadedEvent) {
	log.Printf("📥 PCAP file received: %s", event.FileName)

	// 1️⃣ إنشاء FileSystem
	fs := infra.NewLocalFileSystem()

	// 2️⃣ فتح ملف PCAP عبر FileSystem
	file, err := fs.Open(event.Path)
	if err != nil {
		log.Printf("❌ failed to open pcap file: %v", err)
		return
	}
	defer file.Close()

	// 3️⃣ تنفيذ المعالجة الفعلية
	if _, err := logic.ProcessPcap(
		fs,
		file,
		event.FileName,
	); err != nil {
		log.Printf("❌ PCAP processing failed: %v", err)
		return
	}

	log.Printf("✅ PCAP processed successfully: %s", event.FileName)
}
