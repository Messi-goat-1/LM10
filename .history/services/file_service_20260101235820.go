package services

import (
	"LM-Gate/events"
	"log"
)

// OnFileCollection يتم استدعاؤها عند اكتمال تجميع أجزاء الملف بنجاح
func (s *FileService) OnFileCollection(payload events.FileCollectionPayload) {
	// التحقق من صحة بيانات التجميع
	if payload.CollectionID == "" || payload.FinalPath == "" {
		log.Printf("[WARN] Received invalid FileCollectionPayload: %+v", payload)
		return
	}

	// طباعة نجاح التجميع باستخدام slog ليكون في سطر واحد
	log.Printf("📦 [SERVICE] File Collection Completed | ID: %s | Name: %s | Path: %s | Status: %s",
		payload.CollectionID, payload.FileName, payload.FinalPath, payload.Status)

	// TODO: هنا يمكنك إضافة منطق إضافي مثل تسجيل العملية في قاعدة البيانات
}
