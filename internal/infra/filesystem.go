package infra

import (
	"os"
	"path/filepath"
)

// FileSystem يعرّف العمليات الأساسية على الملفات
// الهدف: فصل منطق الملفات عن logic / work
type FileSystem interface {
	MkdirAll(path string) error
	Create(path string) (*os.File, error)
	Open(path string) (*os.File, error)
	WriteFile(path string, data []byte) error
}

// LocalFileSystem تنفيذ فعلي باستخدام نظام الملفات المحلي
type LocalFileSystem struct{}

// NewLocalFileSystem ينشئ FileSystem محلي
func NewLocalFileSystem() *LocalFileSystem {
	return &LocalFileSystem{}
}

// MkdirAll ينشئ مجلدات مع المسار الكامل
func (fs *LocalFileSystem) MkdirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// Create ينشئ ملف مع التأكد من وجود المجلد الأب
func (fs *LocalFileSystem) Create(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if err := fs.MkdirAll(dir); err != nil {
		return nil, err
	}
	return os.Create(path)
}

// Open يفتح ملف موجود
func (fs *LocalFileSystem) Open(path string) (*os.File, error) {
	return os.Open(path)
}

// WriteFile يكتب بيانات إلى ملف (وينشئ المسار إن لزم)
func (fs *LocalFileSystem) WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := fs.MkdirAll(dir); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
