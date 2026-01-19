package services

import (
	"LM-Gate/events"
	"fmt"
	"log/slog"
	"time"
)

// Manager هو "مركز التحكم" الذي يربط كل الخدمات ببعضها
type Manager struct {
	FileService *FileService
	// مستقبلاً يمكنك إضافة خدمات أخرى هنا بسهولة:
	// AuthService  *AuthService
	// LogService   *LogService
}

func NewManager(fs *FileService) *Manager {
	return &Manager{
		FileService: fs,
	}
}

// توجيه الطلبات للخدمة المناسبة
func (m *Manager) OnFileDetected(payload events.FileDetectedPayload) {
	m.FileService.OnFileDetected(payload)
}

func OnMessage(msg ChunkMessage) error {
	if err := ValidateMessage(msg); err != nil {
		return err
	}

	if msg.IsEOF {
		// 1. تجميع الملف
		filePath, err := AssembleFile(msg.FileID)
		if err != nil {
			return fmt.Errorf("failed to assemble file: %v", err)
		}

		// 2. إنشاء حدث التجميع (هنا الإضافة)
		collectionEvent := events.FileCollectionEvent{
			Payload: events.FileCollectionPayload{
				CollectionID: msg.FileID,
				FileName:     msg.FileID + ".pcap",
				FinalPath:    filePath,
				Status:       "assembled_successfully",
			},
			Timestamp: time.Now(),
		}

		// 3. إرسال الحدث إلى الـ Manager أو عبر RabbitMQ
		// (يمكنك استخدام slog هنا لتوثيق نجاح التجميع)
		slog.Info("📦 File successfully assembled and collection event created",
			slog.String("collection_id", msg.FileID))

		// 4. البدء في المعالجة
		return ProcessFile(msg.FileID, filePath)
	}

	return StoreChunk(msg)
}
