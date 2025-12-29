package services

import (
	"LM-Gate/analysis" // تأكد أن المسار يطابق مشروعك
	"context"          //
	"fmt"

	"github.com/google/gopacket/pcap"
)

type PCAPService struct{}

func NewPCAPService() *PCAPService {
	return &PCAPService{}
}

// Analyze يقوم بفتح الملف وتمرير السياق لمحرك التحليل
func (s *PCAPService) Analyze(ctx context.Context, fileID string, filePath string) error {
	// 1. فتح ملف الـ PCAP
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return fmt.Errorf("failed to open PCAP (%s): %w", fileID, err)
	}
	defer handle.Close()

	// 2. تمرير الـ ctx والـ handle إلى RunFullAnalysis
	// تأكد من أن RunFullAnalysis في حزمة analysis تستقبل (context.Context, *pcap.Handle)
	return analysis.RunFullAnalysis(ctx, handle)
}
