package services

import (
	"fmt"
	"time"
)

// FileService مسؤول عن أي شيء له علاقة بالملفات
type FileService struct{}

// NewFileService constructor بسيط
func NewFileService() *FileService {
	return &FileService{}
}

// OnFileDetected هذا أول behavior حقيقي عندك
func (s *FileService) OnFileDetected(fileName string, size int64) {
	fmt.Println("📦 [SERVICE] File detected")
	fmt.Println("   name:", fileName)
	fmt.Println("   size:", size)
	fmt.Println("   time:", time.Now().Format(time.RFC3339))
}
