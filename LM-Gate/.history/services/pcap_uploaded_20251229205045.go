package services

import (
	"LM-Gate/analysis" // تأكد من مسار الـ package الصحيح لديك
	"context"
	"fmt"

	"github.com/google/gopacket/pcap"
)

type PCAPService struct{}

func NewPCAPService() *PCAPService {
	return &PCAPService{}
}

// Analyze يقوم بتنفيذ منطق التحليل الأساسي.
// تم تحديث التوقيع لاستقبال Context لضمان استقرار النظام.
func (s *PCAPService) Analyze(ctx context.Context, fileID string, filePath string) error {
	// فتح ملف الـ PCAP من القرص [cite: 8]
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return fmt.Errorf("failed to open PCAP (%s): %w", fileID, err)
	}
	defer handle.Close()

	// استدعاء محرك التحليل وتمرير الـ Context
	return analysis.RunFullAnalysis(ctx, handle)
}
