package services

import (
	"LM-Gate/events"
	"log"
	"log/slog" // استيراد المكتبة الجديدة
)

type FileService struct {
	logger *slog.Logger
}

func NewFileService(logger *slog.Logger) *FileService {
	return &FileService{
		logger: logger,
	}
}

func (s *FileService) OnFileDetected(payload events.FileDetectedPayload) {
	// 1. التحقق من صحة البيانات
	if payload.FileID == "" || payload.SizeBytes <= 0 {
		s.logger.Warn("⚠️ Invalid FileDetectedPayload received",
			slog.String("file_id", payload.FileID),
			slog.Int64("size", payload.SizeBytes),
		)
		return
	}

	// 2. طباعة الحدث في سطر واحد منظم باستخدام slog
	s.logger.Info("🚀 [SERVICE] File Detected and Processed",
		slog.String("id", payload.FileID),
		slog.String("name", payload.FileName),
		slog.Int64("size_bytes", payload.SizeBytes),
		slog.String("type", payload.FileType),
		slog.String("checksum", payload.Checksum),
	)
}

// داخل services/file_service.go
func (s *FileService) OnFileCollection(payload events.FileCollectionPayload) {
	// استخدم slog الذي أضفناه سابقاً لطباعة النجاح
	log.Printf("📦 [SERVICE] Collection Completed: %s", payload.CollectionID)
}
