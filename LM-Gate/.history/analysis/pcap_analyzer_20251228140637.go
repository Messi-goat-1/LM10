package analysis

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// 1. دالة الفتح: وظيفتها فقط "الإمساك" بالملف من القرص
func GetFileHandle(filePath string) (*pcap.Handle, error) {
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return nil, fmt.Errorf("فشل الوصول للملف: %v", err)
	}
	return handle, nil
}

// 2. دالة التحليل: وظيفتها استخراج المعلومات (كل شيء)
func RunFullAnalysis(handle *pcap.Handle) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// سنختبر أول 5 حزم فقط لنرى الهيكل
	count := 0
	for packet := range packetSource.Packets() {
		count++
		fmt.Printf("\n--- Packet #%d ---\n", count)

		// هذا السطر سيطبع لك كل تفاصيل الحزمة (IPs, Ports, Payload)
		fmt.Println(packet.String())

		if count >= 5 {
			break
		}
	}
}

// 3. الدالة التي يناديها server.go للربط بينهما
func AnalyzePCAP(fileID string, filePath string) error {
	// تنفيذ الوظيفة الأولى
	handle, err := GetFileHandle(filePath)
	if err != nil {
		return err
	}
	defer handle.Close()

	// تنفيذ الوظيفة الثانية
	RunFullAnalysis(handle)

	return nil
}
