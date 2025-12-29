package handlers

import (
	"LM-Gate/events"
	"encoding/json"
	"fmt"
)

// Handler هو الواجهة الموحدة لجميع معالجي الأحداث.
// أي Handler جديد يجب أن يحقق هذه الواجهة فقط.
type Handler interface {
	Handle(data []byte) error
}

// EventDispatcher مسؤول عن توجيه الأحداث الواردة بناءً على مفتاح التوجيه.
// تم تحسينه ليكون مغلقاً أمام التعديل ومفتوحاً للإضافة (OCP).
type EventDispatcher struct {
	// handlers تخزن المعالجين المسجلين ديناميكياً
	handlers map[string]Handler
}

// NewEventDispatcher ينشئ نسخة جديدة من الموزع.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string]Handler),
	}
}

// RegisterHandler يسمح بإضافة معالج جديد لأي routing key دون تعديل هذا الملف.
func (d *EventDispatcher) RegisterHandler(routingKey string, handler Handler) {
	d.handlers[routingKey] = handler
}

// Dispatch يقوم بتوجيه البيانات إلى المعالج المناسب.
// إذا لم يوجد معالج مسجل للـ key، فإنه يعيد خطأ.
func (d *EventDispatcher) Dispatch(routingKey string, data []byte) error {
	handler, exists := d.handlers[routingKey]
	if !exists {
		return fmt.Errorf("unknown routing key: %s", routingKey)
	}

	// تنفيذ المعالج المسجل مباشرة
	return handler.Handle(data)
}

func (h *FileChunkHandler) Handle(data []byte) error {
	var event events.FileChunkEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}
	// تنفيذ المنطق الخاص بك هنا
	return nil
}
