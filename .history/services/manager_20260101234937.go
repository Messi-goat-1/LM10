package services

import (
	"LM-Gate/events"
)

type Manager struct {
	FileService *FileService
}

func NewManager(fs *FileService) *Manager {
	return &Manager{FileService: fs}
}

// توجيه حدث الاكتشاف
func (m *Manager) OnFileDetected(payload events.FileDetectedPayload) {
	m.FileService.OnFileDetected(payload)
}

// أضف هذه الدالة لتوجيه حدث التجميع
func (m *Manager) OnFileCollection(payload events.FileCollectionPayload) {
	// هنا يمكنك استدعاء خدمة التحليل أو الطباعة
	m.FileService.OnFileCollection(payload)
}
