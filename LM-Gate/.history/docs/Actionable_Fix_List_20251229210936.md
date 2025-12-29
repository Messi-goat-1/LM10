🛠️ قائمة الإصلاحات (Actionable Fix List)
🔴 أولوية عالية جدًا (لازم أولاً)







7️⃣ تفعيل Versioning للأحداث

❗ النوع: Compatibility
📂 الملفات:

events/event.go

handlers/dispatcher.go

📍 المكان:

Event.Version


🔧 المشكلة:

موجود لكنه غير مستخدم

✅ الإصلاح:

تحقق من version داخل dispatcher

8️⃣ توحيد Logging

❗ النوع: Observability
📂 الملفات:

أغلب services/*

handlers/*

🔧 المشكلة:

استخدام fmt.Println فقط

✅ الإصلاح:

Logger موحد

log levels (INFO / ERROR)

9️⃣ اختبارات آلية (Tests)

❗ النوع: Quality
📂 الملفات:

services/file_chunk.go

handlers/dispatcher.go

📍 الأماكن المقترحة للاختبار:

isComplete
reassemble
Dispatch

🧠 خريطة سريعة (Cheat Sheet)
#	المشكلة	الملف
1	ترتيب الـ chunks	services/file_chunk.go
2	إعادة الإرسال	services/file_chunk.go
3	JSON errors	handlers/dispatcher.go
4	Dispatcher OCP	handlers/dispatcher.go
5	Context	services/pcap_uploaded.go
6	Async analysis	handlers/pcap_uploaded.go
7	Event versioning	events/event.go
8	Logging	services / handlers
9	Tests	services / handlers

--------------------------------
🛠️ إصلاح مشاكل الأحداث (Event Hardening)
🔴 المشكلة 1: عدم التحقق من JSON (خطر)
📂 الملف
handlers/dispatcher.go

📍 المكان

داخل الدالة:

func (d *EventDispatcher) Dispatch(routingKey string, data []byte) error

❌ الوضع الحالي (تقريبًا)
json.Unmarshal(data, &event)

✅ الحل

تحقق من الخطأ وأوقف الحدث إذا كان فاسدًا:

if err := json.Unmarshal(data, &event); err != nil {
    return fmt.Errorf("invalid JSON for %s: %w", routingKey, err)
}


📌 النتيجة:

لا event فاسد يمر

لا panic

أخطاء واضحة في اللوق

🔴 المشكلة 2: Version موجود لكن غير مستخدم
📂 الملفات
events/event.go
handlers/dispatcher.go

1️⃣ عرّف version رسميًا
📂 events/event.go
const CurrentEventVersion = 1

2️⃣ تحقق من version داخل dispatcher
📂 handlers/dispatcher.go
if event.Version != events.CurrentEventVersion {
    return fmt.Errorf(
        "unsupported event version %d",
        event.Version,
    )
}


📌 النتيجة:

أي تغيير مستقبلي يكون محمي

Backward compatibility واضح

🔴 المشكلة 3: لا يوجد Contract Tests للأحداث
📂 الملفات الجديدة (مقترحة)
events/contracts_test.go

🧪 اختبارات لازم تضيفها
func TestContract_FileChunkEvent(t *testing.T)
func TestContract_FileDetectedEvent(t *testing.T)
func TestContract_PCAPAnalyzeEvent(t *testing.T)

📍 ماذا تختبر؟

داخل كل اختبار:

Marshal → Unmarshal

تطابق الحقول

عدم فقدان البيانات

مثال بسيط
data, _ := json.Marshal(event)
var decoded events.FileChunkEvent
json.Unmarshal(data, &decoded)

if decoded.Payload.FileID != event.Payload.FileID {
    t.Fatal("contract broken")
}


📌 النتيجة:

أي كسر في العقد ينكشف فورًا

أمان عالي عند refactor

🧠 خريطة الإصلاح السريعة
المشكلة	الملف	الحل
JSON فاسد	handlers/dispatcher.go	تحقق Unmarshal
Version غير مستخدم	events/event.go	const version
Version غير مفحوص	handlers/dispatcher.go	check version
لا Contract Test	events/contracts_test.go	add tests


-------------------------------------------------------
🛠️ إصلاح مشاكل الموديولية (Action Plan)
🔴 المشكلة 1: Dispatcher نقطة مركزية حساسة (Low Modularity)
❗ المشكلة

كل Event جديد يتطلب تعديل switch

هذا يقلل قابلية التوسع

📂 الملف
handlers/dispatcher.go

🔧 الدالة الحالية
func (d *EventDispatcher) Dispatch(routingKey string, data []byte) error

✅ الحل (تحسين الموديولية)
1️⃣ عرّف Interface للـ Handler

📂 ملف جديد:

handlers/handler.go

type Handler interface {
    Handle(data []byte) error
}

2️⃣ غيّر dispatcher من switch إلى registry

📂 handlers/dispatcher.go

type EventDispatcher struct {
    handlers map[string]Handler
}

func NewDispatcher() *EventDispatcher {
    return &EventDispatcher{
        handlers: make(map[string]Handler),
    }
}

func (d *EventDispatcher) Register(routingKey string, h Handler) {
    d.handlers[routingKey] = h
}

func (d *EventDispatcher) Dispatch(routingKey string, data []byte) error {
    h, ok := d.handlers[routingKey]
    if !ok {
        return fmt.Errorf("no handler for %s", routingKey)
    }
    return h.Handle(data)
}


📌 النتيجة:

إضافة Event جديد = Register

لا تعديل dispatcher

Modularity أعلى

🔴 المشكلة 2: Manager يقوم بأكثر من مسؤولية (Fat Service)
❗ المشكلة

Manager:

تخزين chunks

تحقق اكتمال

دمج

تنظيف

هذا مقبول الآن، لكن سيكبر.

📂 الملف
services/file_chunk.go

🔧 الدوال المتأثرة
OnChunkReceived
isComplete
reassemble

✅ الحل (تفكيك تدريجي بدون كسر)
1️⃣ فصل التخزين في Storage

📂 ملف جديد:

services/chunk_storage.go

type ChunkStorage interface {
    SaveChunk(fileID string, index int, data []byte) error
    ListChunks(fileID string) ([]string, error)
    Cleanup(fileID string) error
}

2️⃣ Manager يعتمد على Interface

📂 services/file_chunk.go

type Manager struct {
    storage ChunkStorage
}


📌 الآن:

Manager = Orchestrator

Storage = تنفيذ فعلي

🔴 المشكلة 3: Services تعرف تفاصيل filesystem مباشرة
❗ المشكلة

os.Mkdir

os.WriteFile

os.ReadDir
داخل service

هذا يضعف العزل.

📂 الملفات
services/file_chunk.go
services/pcap_uploaded.go

🔧 الدوال
OnChunkReceived
reassemble
Analyze

✅ الحل
1️⃣ إنشاء FileSystem abstraction

📂 ملف جديد:

services/fs.go

type FileSystem interface {
    WriteFile(path string, data []byte) error
    ReadDir(path string) ([]os.DirEntry, error)
    RemoveAll(path string) error
}

2️⃣ حقنه في services
type Manager struct {
    fs FileSystem
}


📌 النتيجة:

Tests أسهل

موديولية أعلى

قابلية تغيير storage لاحقًا

🟠 المشكلة 4: PCAP Analysis غير معزول كموديول مستقل
❗ المشكلة

Service يستدعي analyzer مباشرة

📂 الملفات
services/pcap_uploaded.go
analysis/pcap_analyzer.go

🔧 الدالة
func (s *PCAPService) Analyze(...)

✅ الحل
1️⃣ تعريف Interface للتحليل

📂 ملف جديد:

analysis/analyzer.go

type Analyzer interface {
    Analyze(path string) error
}

2️⃣ PCAPService يعتمد على Interface

📂 services/pcap_uploaded.go

type PCAPService struct {
    analyzer analysis.Analyzer
}

📌 النتيجة:

تقدر تضيف:

AIAnalyzer

MockAnalyzer

بدون تعديل service

🧠 خريطة الإصلاح النهائية
المشكلة	الملف	الدالة
Dispatcher مركزي	handlers/dispatcher.go	Dispatch
Fat Manager	services/file_chunk.go	OnChunkReceived
FS Coupling	services/file_chunk.go	reassemble
PCAP Coupling	services/pcap_uploaded.go	Analyze