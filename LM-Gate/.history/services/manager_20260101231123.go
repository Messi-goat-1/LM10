package services

import "LM-Gate/events"

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
