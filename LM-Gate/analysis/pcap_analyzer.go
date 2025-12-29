package analysis

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// GetFileHandle opens a PCAP file from disk (offline mode).
//
// NOTE: This function is purely technical.
// It only opens the file and prepares it for reading.
// TODO: Add file existence and permission checks before opening.
func GetFileHandle(filePath string) (*pcap.Handle, error) {
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	return handle, nil
}

// RunFullAnalysis performs the actual packet analysis.
//
// NOTE: This function is the analysis engine.
// It reads packets one by one from the PCAP file.
// FIXME: Currently stops after a small number of packets.
// TODO: Analyze all packets and extract structured data.
func RunFullAnalysis(handle *pcap.Handle) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Test only the first few packets to inspect the structure
	count := 0
	for packet := range packetSource.Packets() {
		count++
		fmt.Printf("\n--- Packet #%d ---\n", count)

		// This prints all packet details (IPs, Ports, Payload, etc.)
		fmt.Println(packet.String())

		if count >= 2 {
			break
		}
	}
}

// AnalyzePCAP is the main function called by server.go.
//
// NOTE: This function connects file handling with analysis logic.
// It opens the file first, then runs the full analysis.
// TODO: Pass context for cancellation or timeouts.
func AnalyzePCAP(fileID string, filePath string) error {
	// Step 1: Open the PCAP file
	handle, err := GetFileHandle(filePath)
	if err != nil {
		return err
	}
	defer handle.Close()

	// Step 2: Run packet analysis
	RunFullAnalysis(handle)

	return nil
}
