package events

// FileChunkPayload يحتوي على بيانات قطعة الملف
type FileChunkPayload struct {
	FileID      string `json:"file_id"`      // معرف فريد للملف الأصلي
	ChunkIndex  int    `json:"chunk_index"`  // ترتيب القطعة (يبدأ من 0)
	TotalChunks int    `json:"total_chunks"` // إجمالي عدد القطع
	Data        []byte `json:"data"`         // محتوى القطعة الخام
}

// FileChunkEvent هو الحدث الذي سيتم استقباله من RabbitMQ
type FileChunkEvent struct {
	Payload   FileChunkPayload `json:"payload"`
	Timestamp string           `json:"timestamp"`
}
