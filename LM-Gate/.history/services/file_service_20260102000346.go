package services

import (
	"LM-Gate/events"
	"log/slog"
)

// OnFileCollection يتم استدعاؤها بعد نجاح دالة AssembleFile
func (s *FileService) OnFileCollection(payload events.FileCollectionPayload) {
	// التحقق من البيانات
	if payload.CollectionID == "" || payload.FinalPath == "" {
		s.logger.Warn("⚠️ Received incomplete collection data", slog.Any("payload", payload))
		return
	}

	// طباعة النجاح باستخدام slog في سطر واحد منظم كما طلبت سابقاً
	s.logger.Info("📦 [SERVICE] File Collection Processed Successfully",
		slog.String("id", payload.CollectionID),
		slog.String("file", payload.FileName),
		slog.String("path", payload.FinalPath),
		slog.String("status", payload.Status),
	)

	// TODO: يمكنك هنا استدعاء ProcessFile(payload.CollectionID, payload.FinalPath) لبدء التحليل
}

func (s *FileService) OnFileCollection(payload events.FileCollectionPayload) {
	s.logger.Info("📦 [SERVICE] تم تجميع الملف بنجاح",
		slog.String("id", payload.CollectionID),
		slog.String("path", payload.FinalPath),
	)
}
