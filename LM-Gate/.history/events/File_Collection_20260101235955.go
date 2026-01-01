package events

import "time"

// FileCollectionPayload يمثل البيانات بعد تجميع الملف بنجاح
type FileCollectionPayload struct {
	CollectionID string `json:"collection_id"` // معرف الملف الفريد
	FileName     string `json:"file_name"`     // اسم الملف النهائي (مثلاً .pcap)
	FinalPath    string `json:"final_path"`    // المسار الذي خزن فيه الملف المجمع
	Status       string `json:"status"`        // حالة التجميع (success)
}

// FileCollectionEvent الغلاف الرئيسي للحدث
type FileCollectionEvent struct {
	Payload   FileCollectionPayload `json:"payload"`
	Timestamp time.Time             `json:"timestamp"`
}
