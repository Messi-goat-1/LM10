package DataProcessor

import (
	pb "LM10/DataProcessor/proto"
	"log"
)

// AnalyzeAndSplit وظيفتها فقط: تقسيم الدفعة الكبيرة إلى قطع صغيرة
func (s *AnalyticsServer) AnalyzeAndSplit(req *pb.DataBatch) {
	batchID := req.GetBatchId()
	entries := req.GetEntries()

	// 1. تحديد حجم المجموعة (مثلاً كل 100 سطر في مجموعة)
	chunkSize := 100

	// 2. عملية التقسيم (The Splitting Logic)
	for i := 0; i < len(entries); i += chunkSize {
		end := i + chunkSize
		if end > len(entries) {
			end = len(entries)
		}

		// استخراج المجموعة
		currentChunk := entries[i:end]

		// 3. التمرير الفوري لدالة التحليل
		// دالة AnalyzeAndSplit تنتهي مهمتها هنا وتنتقل للقطعة التالية
		s.PerformAnalysis(batchID, currentChunk)
	}

	// هنا يموت الكائن req تماماً من الذاكرة بعد معالجة آخر قطعة
	log.Printf("✅ تم تقسيم الدفعة [%s] بالكامل وتحويلها للتحليل.", batchID)
}

// PerformAnalysis هي الدالة التي ستنفذ فيها فكرتك المستقبلية (مثل التجميع حسب الـ IP)
func (s *AnalyticsServer) PerformAnalysis(batchID string, data []string) {
	// حالياً: مجرد طباعة للتأكد من وصول المجموعة
	log.Printf("🔬 جاري تحليل مجموعة مصغرة من الدفعة [%s] تحتوي على %d عنصر", batchID, len(data))

	// مستقبلاً: هنا تضع كود (Group by IP)
}
