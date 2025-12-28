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
func (s *Manager) OnChunkReceived(fileID string, index int, total int, data []byte) {
	// 1. التحقق من صحة البيانات (مستوحى من ValidateMessage في server.go)
	if fileID == "" || len(data) == 0 {
		fmt.Println("⚠️ [SERVICE] بيانات القطعة غير صالحة")
		return
	}

	// 2. إنشاء مجلد فرعي لكل FileID لتنظيم القطع
	fileDir := filepath.Join(s.tempDir, fileID)
	os.MkdirAll(fileDir, 0755)

	// 3. تخزين القطعة في القرص (تطوير لـ StoreChunk في server.go)
	chunkPath := filepath.Join(fileDir, fmt.Sprintf("part_%d", index))
	err := os.WriteFile(chunkPath, data, 0644)
	if err != nil {
		fmt.Printf("❌ خطأ أثناء حفظ القطعة %d: %v\n", index, err)
		return
	}

	fmt.Printf("📥 استلام القطعة %d من %d للملف: %s\n", index+1, total, fileID)

	// 4. التحقق من اكتمال جميع القطع (مستوحى من IsFileComplete في server.go)
	if s.isComplete(fileDir, total) {
		fmt.Println("🚀 اكتملت جميع القطع، جاري تجميع الملف...")
		s.reassemble(fileID, fileDir, total)
	}
}

// isComplete يتحقق من عدد الملفات في المجلد المؤقت
func (s *Manager) isComplete(dir string, total int) bool {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	return len(files) == total
}

// reassemble يقوم بدمج القطع (تطوير لـ AssembleFile في server.go)
func (s *Manager) reassemble(fileID string, dir string, total int) {
	finalPath := filepath.Join(s.storageDir, fileID)
	finalFile, err := os.Create(finalPath)
	if err != nil {
		fmt.Printf("❌ فشل إنشاء الملف النهائي: %v\n", err)
		return
	}
	defer finalFile.Close()

	// دمج القطع بالترتيب الصحيح
	for i := 0; i < total; i++ {
		chunkPath := filepath.Join(dir, fmt.Sprintf("part_%d", i))
		content, err := os.ReadFile(chunkPath)
		if err != nil {
			fmt.Printf("❌ فشل قراءة القطعة %d: %v\n", i, err)
			return
		}
		finalFile.Write(content)
	}

	// 5. التنظيف (مستوحى من Cleanup في server.go)
	os.RemoveAll(dir)
	fmt.Printf("✅ تم تجميع الملف بنجاح: %s\n", finalPath)
}
