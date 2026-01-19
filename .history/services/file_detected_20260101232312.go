package services

import (
	"LM-Gate/events"
	"log/slog" // استيراد المكتبة الجديدة
	"os"
)

type FileService struct {
	logger *slog.Logger
}

func NewFileService() *FileService {
	// إعداد الـ logger ليكون بتنسيق JSON (مناسب جداً للإنتاج)
	// أو TextHandler إذا كنت تفضل القراءة البشرية البسيطة
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &FileService{
		logger: logger,
	}
}

func (s *FileService) OnFileDetected(payload events.FileDetectedPayload) {
	// 1. التحقق من صحة البيانات
	if payload.FileID == "" || payload.SizeBytes <= 0 {
		s.logger.Warn("Invalid FileDetectedPayload",
			slog.String("file_id", payload.FileID),
			slog.Int64("size", payload.SizeBytes),
		)
		return
	}

	// 2. طباعة الحدث بشكل مهيكل (Structured)
	s.logger.Info("🚀 File Detected and processed",
		slog.String("id", payload.FileID),
		slog.String("name", payload.FileName),
		slog.Int64("size_bytes", payload.SizeBytes),
		slog.String("type", payload.FileType),
		slog.String("checksum", payload.Checksum),
	)
}
