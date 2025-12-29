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

// RunFullAnalysis يقوم بالتحليل الفعلي للحزم مع مراقبة السياق
func RunFullAnalysis(ctx context.Context, handle *pcap.Handle) error {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	count := 0
	for {
		select {
		case <-ctx.Done():
			// إرجاع سبب توقف السياق (Timeout أو Cancellation)
			return ctx.Err()
		case packet, ok := <-packetSource.Packets():
			if !ok {
				return nil // انتهى الملف بنجاح
			}
			count++
			fmt.Printf("\n--- Packet #%d ---\n", count)
			fmt.Println(packet.String())

			// FIXME: حد مؤقت للتجربة [cite: 13, 14]
			if count >= 2 {
				return nil
			}
		}
	}
}

// AnalyzePCAP هي الدالة الأساسية التي يتم استدعاؤها
func AnalyzePCAP(ctx context.Context, fileID string, filePath string) error {
	// الخطوة 1: فتح ملف الـ PCAP واستقبال القيمتين (handle و err) [cite: 11]
	handle, err := GetFileHandle(filePath)
	if err != nil {
		return err
	}
	defer handle.Close()

	// الخطوة 2: تشغيل التحليل مع تمرير الـ Context
	return RunFullAnalysis(ctx, handle)
}
