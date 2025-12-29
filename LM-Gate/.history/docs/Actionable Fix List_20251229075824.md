🛠️ قائمة الإصلاحات (Actionable Fix List)
🔴 أولوية عالية جدًا (لازم أولاً)
1️⃣ ضمان ترتيب الـ Chunks

❗ النوع: Data Integrity / Bug
📂 الملفات:

services/file_chunk.go

📍 المكان:

func (s *Manager) isComplete(...)
func (s *Manager) reassemble(...)


🔧 المشكلة:

يعتمد على len(files) == total

لا يتحقق من وجود part_0 → part_(n-1)

✅ الإصلاح:

تحقق من وجود كل ملف بالاسم

part_0, part_1, part_2, ...

2️⃣ منع إعادة كتابة Chunk مرتين (Idempotency)

❗ النوع: Data Consistency
📂 الملفات:

services/file_chunk.go

📍 المكان:

func (s *Manager) OnChunkReceived(...)


🔧 المشكلة:

نفس chunk ممكن ينكتب مرتين

✅ الإصلاح:

قبل WriteFile:

إذا الملف موجود → تجاهل أو تحقق checksum

3️⃣ تجاهل أخطاء JSON خطر

❗ النوع: Error Handling
📂 الملفات:

handlers/dispatcher.go

📍 المكان:

json.Unmarshal(data, &event)


🔧 المشكلة:

لا يتم التحقق من الخطأ

✅ الإصلاح:

if err := json.Unmarshal(...); err != nil {
    return err
}

🟠 أولوية متوسطة
4️⃣ تحسين Dispatcher (OCP)

❗ النوع: Design Improvement
📂 الملفات:

handlers/dispatcher.go

📍 المكان:

switch routingKey { ... }


🔧 المشكلة:

كل Event جديد = تعديل الملف

✅ الإصلاح (لاحقًا):

استبداله بـ:

map[string]Handler

5️⃣ إضافة Context للتحليل

❗ النوع: Stability / Control
📂 الملفات:

services/pcap_uploaded.go

analysis/pcap_analyzer.go

📍 المكان:

func Analyze(...)


🔧 المشكلة:

لا يوجد timeout أو cancel

✅ الإصلاح:

استخدام context.Context

تمريره من handler → service → analyzer

6️⃣ فصل تحليل PCAP عن المسار الأساسي

❗ النوع: Performance
📂 الملفات:

handlers/pcap_uploaded.go

🔧 المشكلة:

التحليل يشتغل فورًا وقد يستهلك موارد

✅ الإصلاح:

تشغيله:

goroutine

أو job queue

🟡 أولوية منخفضة (لكن مهمة)
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