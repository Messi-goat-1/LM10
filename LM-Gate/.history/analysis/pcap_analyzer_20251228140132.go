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

	fmt.Println("📊 [GOPACKET] جاري استخراج كافة تفاصيل الحزم...")

	for packet := range packetSource.Packets() {
		// gopacket هنا تقوم بكل العمل:
		// تحليل الطبقات، العناوين، البروتوكولات، والبيانات (Payload)
		fmt.Println(packet.String()) // هذه تطبع لك "كل المعلومات" كما رأيت في تجربتك

		// ملاحظة: يمكنك وضع شرط توقف إذا كان الملف ضخماً جداً
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
