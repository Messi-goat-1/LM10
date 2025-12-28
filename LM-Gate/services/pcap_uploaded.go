package services

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type PCAPService struct{}

func NewPCAPService() *PCAPService {
	return &PCAPService{}
}

// Analyze هو المنطق الحقيقي لتحليل الـ PCAP
func (s *PCAPService) Analyze(fileID string, filePath string) error {
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return fmt.Errorf("فشل فتح PCAP (%s): %w", fileID, err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	count := 0
	for packet := range packetSource.Packets() {
		count++
		fmt.Printf("\n[PCAP:%s] Packet #%d\n", fileID, count)
		fmt.Println(packet.String())

		if count >= 2 {
			break
		}
	}

	return nil
}
