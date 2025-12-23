package types

import (
	"errors"
	"redis/interfaces"
	"sync"
	"time"
)

// التأكد من أن String يحقق واجهة Item للتحليل
var _ interfaces.Item = (*String)(nil)

type String struct {
	mu           sync.RWMutex
	value        string    // القيمة المخزنة (نص أو رقم محول لنص)
	createdAt    time.Time // وقت إنشاء العنصر (مهم جداً للتحليل الزمني)
	lastAccessed time.Time // آخر وقت تم فيه التفاعل مع البيانات
	hits         int64     // عداد التكرار (كم مرة تم استدعاء هذا المفتاح)
}

func NewString(value string) *String {
	now := time.Now()
	return &String{
		value:        value,
		createdAt:    now,
		lastAccessed: now,
		hits:         0,
	}
}

func (s *String) RawValue() interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

func (s *String) CreatedAt() time.Time {
	return s.createdAt
}
func (s *String) AsString() (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value, true
}
func (s *String) LastAccessed() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastAccessed
}

func (s *String) Frequency() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hits
}

func (s *String) Category() string {
	return "analytics_string"
}

// Update هي الدالة الأهم لاستقبال بيانات gRPC ومعالجتها
func (s *String) Update(newValue interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.hits++
	s.lastAccessed = time.Now()

	if val, ok := newValue.(string); ok {
		s.value = val
		return nil
	}

	return errors.New("unsupported type for update")
}

// OnDelete يتم استدعاؤها تلقائياً من ملف db.go عند الحذف
func (s *String) OnDelete(key string, db interfaces.DatabaseOps) {
	// يعنيfmt.Printf("[Final Report] Key: %s, Total Hits: %d, Lifetime: %v\n", key, s.hits, time.Since(s.createdAt))
}

// --- دوال إضافية للتوافق مع الأنظمة القديمة (إذا احتجت) ---

func (s *String) Type() uint64 {
	return 0 // String Type
}

func (s *String) TypeFancy() string {
	return "string"
}
