package interfaces

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ---------------------------
// Item Interface
// ---------------------------
type Item interface {
	// الوصول للقيمة الخام (مفيدة للـ Serialization)
	RawValue() interface{}

	// معلومات التحليل (Critical for Analytics)
	CreatedAt() time.Time    // متى دخلت المعلومة؟ (لتحليل السلاسل الزمنية)
	LastAccessed() time.Time // متى تم تحليلها آخر مرة؟
	Frequency() int64        // كم مرة تكرر هذا الحدث أو استُدعي هذا المفتاح؟
	AsString() (string, bool)
	// التصنيف
	Category() string // تصنيف البيانات (مثل: "error", "traffic", "user_action")

	// الوظيفة الحسابية (Aggregation)
	// تتيح لك دمج البيانات الجديدة مع القديمة دون مسحها
	Update(newValue interface{}) error

	// Hooks
	OnDelete(key string, db DatabaseOps)
}

type AnalyticsMetric struct {
	Data       float64
	Created    time.Time
	Hits       int64
	Label      string
	sync.Mutex // لحماية البيانات أثناء التحديث
}

func (m *AnalyticsMetric) RawValue() interface{} {
	return m.Data
}

func (m *AnalyticsMetric) CreatedAt() time.Time {
	return m.Created
}

func (m *AnalyticsMetric) Frequency() int64 {
	return m.Hits
}

func (m *AnalyticsMetric) Category() string {
	return m.Label
}

func (m *AnalyticsMetric) AsString() (string, bool) {
	return fmt.Sprintf("%f", m.Data), true
}

func (m *AnalyticsMetric) Update(newValue interface{}) error {
	m.Lock()
	defer m.Unlock()
	val := newValue.(float64)
	m.Data += val // مثال: جمع القيم للتحليل
	m.Hits++
	return nil
}

func (m *AnalyticsMetric) OnDelete(key string, db DatabaseOps) {
	// يمكنك هنا كتابة الكود لإرسال تقرير نهائي قبل الحذف
}

type CounterItem struct {
	sync.Mutex
	count int64
}

func (c *CounterItem) Value() interface{} { return c.count }

func (c *CounterItem) Update(newValue interface{}) error {
	c.Lock()
	defer c.Unlock()

	// التأكد أن القيمة القادمة هي رقم
	val, ok := newValue.(int64)
	if !ok {
		return errors.New("invalid increment value")
	}
	c.count += val
	return nil
}

func (c *CounterItem) OnDelete(key string, db DatabaseOps) {}
