package events

// FileDetectedPayload تفاصيل الملف التقنية
type FileDetectedPayload struct {
	FileID      string `json:"file_id"`
	FileName    string `json:"file_name"`
	SizeBytes   int64  `json:"size_bytes"`
	FileType    string `json:"file_type"`
	Checksum    string `json:"checksum"`
	StorageHint string `json:"storage_hint"`
}

// FileDetectedEvent الحدث الرئيسي
type FileDetectedEvent struct {
	// تم تغيير payload إلى Payload (حرف كبير) لتتمكن مكتبة JSON من قراءته
	// وإضافة وسم json ليتطابق مع البيانات القادمة
	Payload FileDetectedPayload `json:"payload"`

	// يمكنك الاحتفاظ بحقول إضافية خارج الـ Payload إذا أردت
	Timestamp string `json:"timestamp"`
}
