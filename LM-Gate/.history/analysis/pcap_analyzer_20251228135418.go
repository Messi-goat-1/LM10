package analysis

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// الوظيفة الأولى: فتح ملف الـ PCAP وتهيئته
func OpenCapture(filePath string) (*pcap.Handle, error) {
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return nil, fmt.Errorf("خطأ في فتح الملف المادي: %v", err)
	}
	return handle, nil
}

// الوظيفة الثانية: تحليل الحزم الموجودة داخل المقبض
func ProcessPackets(handle *pcap.Handle) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	count := 0
	for packet := range packetSource.Packets() {
		count++

		// منطق التحليل الخاص بك يوضع هنا
		// مثال: طباعة نوع البروتوكول في كل حزمة
		if transportLayer := packet.TransportLayer(); transportLayer != nil {
			// fmt.Println("Protocol:", transportLayer.LayerType())
		}

		if count >= 100 { // للتجربة فقط
			break
		}
	}
	fmt.Printf("✅ تم تحليل %d حزمة بنجاح.\n", count)
}

// الدالة الرئيسية التي تجمعهم (التي يستدعيها server.go)
func AnalyzePCAP(fileID string, filePath string) error {
	fmt.Printf("🔍 بدء المعالجة المنفصلة للملف: %s\n", fileID)

	// 1. استدعاء وظيفة الفتح
	handle, err := OpenCapture(filePath)
	if err != nil {
		return err
	}
	defer handle.Close()

	// 2. استدعاء وظيفة التحليل
	ProcessPackets(handle)

	return nil
}
