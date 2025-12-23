package lmgate

import (
	"errors"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/spf13/cobra"
)

//
// ===== CLI =====
//

var rootCmd = &cobra.Command{
	Use:   "LM10 [file]",
	Short: "LM10 PCAP analyzer",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := GetPathFromArgs(args)

		if err := ValidateFilePath(path); err != nil {
			return err
		}

		return RunPCAP(path)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

//
// ===== ARGUMENTS =====
//

func GetPathFromArgs(args []string) string {
	return args[0]
}

//
// ===== VALIDATION =====
//

func ValidateFilePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("file does not exist")
		}
		return err
	}

	if info.IsDir() {
		return errors.New("path is a directory, not a file")
	}

	return nil
}

//
// ===== PIPELINE =====
//

func RunPCAP(path string) error {
	handle, err := OpenPCAP(path)
	if err != nil {
		return err
	}
	defer ClosePCAP(handle)

	report, err := AnalyzePCAP(handle)
	if err != nil {
		return err
	}

	return SaveResults(report)
}

//
// ===== ANALYSIS =====
//

// الورقة الكبيرة
type AnalysisReport struct {
	TotalPackets int            `json:"total_packets"`
	Protocols    map[string]int `json:"protocols"`
}

// دالة تحليل فقط (لا طباعة – لا حفظ)
func AnalyzePCAP(handle *pcap.Handle) (*AnalysisReport, error) {
	report := &AnalysisReport{
		Protocols: make(map[string]int),
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range source.Packets() {
		report.TotalPackets++

		if tl := packet.TransportLayer(); tl != nil {
			report.Protocols[tl.LayerType().String()]++
		}
	}

	return report, nil
}

//
// ===== SAVE =====
//

// تستقبل الورقة فقط – بدون أي تعديل
func SaveResults(report *AnalysisReport) error {
	// المرحلة القادمة:
	// - JSON
	// - DB
	// - AI
	// - Export
	return nil
}

//
// ===== PCAP IO =====
//

func OpenPCAP(path string) (*pcap.Handle, error) {
	return pcap.OpenOffline(path)
}

func ClosePCAP(handle *pcap.Handle) {
	handle.Close()
}
