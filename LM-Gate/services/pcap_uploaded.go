package services

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// PCAPService is responsible for analyzing PCAP files.
//
// NOTE: This service focuses only on PCAP analysis logic.
// It does not handle file upload, storage, or events.
// TODO: Add support for structured output instead of printing packets.
type PCAPService struct{}

// NewPCAPService creates a new PCAPService instance.
//
// NOTE: Currently stateless.
// TODO: Inject configuration or logger if needed.
func NewPCAPService() *PCAPService {
	return &PCAPService{}
}

// Analyze performs the actual PCAP analysis.
//
// NOTE: This is the core analysis logic.
// It opens the PCAP file and reads packets one by one.
// FIXME: Currently stops after analyzing only a few packets.
func (s *PCAPService) Analyze(fileID string, filePath string) error {
	// Open the PCAP file from disk (offline mode)
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return fmt.Errorf("failed to open PCAP (%s): %w", fileID, err)
	}
	defer handle.Close()

	// Create a packet source to read packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	count := 0
	for packet := range packetSource.Packets() {
		count++
		fmt.Printf("\n[PCAP:%s] Packet #%d\n", fileID, count)

		// Print full packet details (IPs, ports, payload, etc.)
		fmt.Println(packet.String())

		// Temporary limit for testing
		if count >= 2 {
			break
		}
	}

	return nil
}
