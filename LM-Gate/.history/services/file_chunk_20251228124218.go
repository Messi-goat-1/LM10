package services

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileService مسؤول عن إدارة الملفات وتجميع القطع (Chunks)

// بدلاً من FileService، نستخدم Manager لأنه داخل حزمة الخدمات
type Manager struct {
	tempDir    string
	storageDir string
}

// دالة الإنشاء تصبح NewManager
func NewManager() *Manager {
	s := &Manager{
		tempDir:    "./temp_chunks",
		storageDir: "./uploads",
	}
	// ... منطق إنشاء المجلدات كما هو ...
	return s
}

// OnFileDetected معالجة إشعار وجود ملف (الحدث البسيط)
// تحديث التوقيع لاستقبال 5 باراميترات
func (s *Manager) OnFileDetected(fileID string, fileName string, size int64, fileType string, checksum string) {
	fmt.Printf("📦 [SERVICE] ملف جديد: %s (ID: %s, الحجم: %d)\n", fileName, fileID, size)
	// يمكنك استخدام باقي البيانات هنا (checksum, fileType)
}

// OnChunkReceived معالجة وصول قطعة من ملف (تطبيق منطق server.go)
// في ملف services/file_chunk.go أو file_service.go
func (s *Manager) OnChunkReceived(fileID string, chunkIndex int, total int, data []byte) error {
	fileDir := filepath.Join(s.tempDir, fileID)
	os.MkdirAll(fileDir, 0755)

	chunkPath := filepath.Join(fileDir, fmt.Sprintf("part_%d", chunkIndex))
	err := os.WriteFile(chunkPath, data, 0644)
	if err != nil {
		return err
	}

	// تم حذف الأسطر المكررة هنا لإصلاح خطأ missing return
	if s.isComplete(fileDir, total) {
		go s.reassemble(fileID, total)
	}

	return nil // هذا يغلق الدالة بشكل صحيح ويحل مشكلة الصورة
}

// reassemble يقوم بدمج القطع (تطوير لـ AssembleFile في server.go)
func (s *Manager) reassemble(fileID string, totalChunks int) error {
	finalPath := filepath.Join(s.storageDir, fileID)

	// فتح ملف جديد للكتابة (الملف النهائي)
	dst, err := os.Create(finalPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// قراءة القطع بالترتيب من القرص ودمجها
	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(s.tempDir, fileID, fmt.Sprintf("part_%d", i))

		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			return err
		}

		// كتابة القطعة في الملف النهائي
		dst.Write(chunkData)
	}

	// تنظيف المجلد المؤقت (بديل delete من الـ map)
	os.RemoveAll(filepath.Join(s.tempDir, fileID))
	return nil
}
