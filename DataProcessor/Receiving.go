package DataProcessor

import (
	pb "LM10/DataProcessor/proto"
	"context"
	"log"
)

type AnalyticsServer struct {
	pb.UnimplementedAnalyticsServiceServer
}

// 1. الدالة الأولى: تستلم وتمرر وتتحرر (تموت) فوراً
func (s *AnalyticsServer) ReceivedBatch(ctx context.Context, req *pb.DataBatch) (*pb.BatchResponse, error) {
	// تمرير الكائن للمسار الداخلي (الكائن يظل حياً هنا)
	go s.InternalPipeline(req)

	// الدالة تنتهي هنا وتتحرر الذاكرة الخاصة بها
	return &pb.BatchResponse{
		Success: true,
		Message: "تم الاستلام، المعالجة بدأت في الخلفية",
	}, nil
}

// 2. الدالة الثانية: المسار الداخلي (تمسك الكائن وتفحصه)
func (s *AnalyticsServer) InternalPipeline(req *pb.DataBatch) {
	// فحص المشاكل (الدالة الثالثة)
	if !s.verifyBatchIntegrity(req) {
		log.Printf("❌ تم رفض الدفعة %s لوجود مشاكل تقنية", req.GetBatchId())
		return // يموت الكائن هنا إذا فشل الفحص
	}

	// الانتقال للتحليل النهائي (ملف analysis.go)
	s.AnalyzeAndSplit(req)

	// بعد انتهاء AnalyzeAndSplit، يموت الكائن نهائياً من الذاكرة
}

// 3. الدالة الثالثة: التأكد من سلامة الإرسال
func (s *AnalyticsServer) verifyBatchIntegrity(req *pb.DataBatch) bool {
	if req == nil || len(req.GetEntries()) == 0 {
		return false
	}
	// تأكد أن البيانات مطابقة للمواصفات
	return true
}
