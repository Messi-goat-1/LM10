package analysis

import (
	"context"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// GetFileHandle opens a PCAP file from disk (offline mode).
func GetFileHandle(filePath string) (*pcap.Handle, error) {
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	return handle, nil
}

// RunFullAnalysis performs the actual packet analysis.
// تمت إضافة Context لدعم الإلغاء أو المهلة الزمانية (Timeout).
func RunFullAnalysis(ctx context.Context, handle *pcap.Handle) error {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	count := 0
	// استخدام for loop مع select لمراقبة الـ Context
	for {
		select {
		case <-ctx.Done():
			// إذا تم إلغاء السياق (بسبب timeout أو إلغاء يدوي)
			return ctx.Err()[cite:1]
		case packet, ok := <-packetSource.Packets():
			if !ok {
				// انتهى الملف
				return nil
			}
			count++
			fmt.Printf("\n--- Packet #%d ---\n", count)
			fmt.Println(packet.String())

			// FIXME: حد مؤقت للتجربة فقط [cite: 13, 14]
			if count >= 2 {
				return nil
			}
		}
	}
}

// AnalyzePCAP is the main function called by server.go.
func AnalyzePCAP(ctx context.Context, fileID string, filePath string) error {
	// الخطوة 1: فتح ملف PCAP
	handle, err := GetFileHandle(filePath)[cite:11]
	if err != nil {
		return err
	}
	defer handle.Close()

	// الخطوة 2: تشغيل التحليل مع تمرير الـ Context
	return RunFullAnalysis(ctx, handle)
}
