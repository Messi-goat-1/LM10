🔴 System Test 1: Full File Upload & Reassembly
🧪 اسم الاختبار
TestSystem_FileUpload_Reassemble

🎯 الهدف

التأكد أن الملف:

يُقسم

يُرسل

يُخزّن

يُدمج

ويُحذف المؤقت

📂 الملفات المشاركة
client.go
handlers/dispatcher.go
handlers/file_chunk.go
services/file_chunk.go

🔧 الدوال المستخدمة (للمراقبة)
UploadFile
SplitFile
EventDispatcher.Dispatch
FileChunkHandler.Handle
Manager.OnChunkReceived
Manager.isComplete
Manager.reassemble

✅ التحقق (Assertions)

وجود الملف النهائي في uploads/

حذف temp_chunks/fileID

حجم الملف النهائي == حجم الأصلي

🔴 System Test 2: Upload with Missing Chunk
🧪 اسم الاختبار
TestSystem_FileUpload_MissingChunk

🎯 الهدف

التأكد أن النظام لا يدمج الملف إذا كانت قطعة ناقصة.

📂 الملفات المشاركة
handlers/file_chunk.go
services/file_chunk.go

🔧 الدوال المستخدمة
Manager.OnChunkReceived
Manager.isComplete

✅ التحقق

عدم وجود ملف في uploads/

بقاء temp_chunks/fileID

🔴 System Test 3: Duplicate Chunk Handling
🧪 اسم الاختبار
TestSystem_FileUpload_DuplicateChunk

🎯 الهدف

اختبار إعادة إرسال نفس القطعة.

📂 الملفات المشاركة
services/file_chunk.go

🔧 الدوال المستخدمة
Manager.OnChunkReceived

✅ التحقق

عدم فساد الملف

عدم تشغيل reassemble مرتين

🟠 System Test 4: Event Routing (Dispatcher)
🧪 اسم الاختبار
TestSystem_EventRouting_AllEvents

🎯 الهدف

التأكد أن كل routingKey يذهب للـ handler الصحيح.

📂 الملفات المشاركة
handlers/dispatcher.go
handlers/file_chunk.go
handlers/file_detected.go
handlers/pcap_uploaded.go

🔧 الدوال المستخدمة
EventDispatcher.Dispatch
FileChunkHandler.Handle
FileDetectedHandler.Handle
PCAPAnalyzeHandler.Handle

✅ التحقق

كل handler تم استدعاؤه مرة واحدة

routingKey غير معروف → error

🟠 System Test 5: File Detected → Business Logic
🧪 اسم الاختبار
TestSystem_FileDetected_MetadataFlow

🎯 الهدف

التأكد أن metadata تنتقل من event إلى service.

📂 الملفات المشاركة
events/file_detected.go
handlers/file_detected.go
services/file_detected.go

🔧 الدوال المستخدمة
EventDispatcher.Dispatch
FileDetectedHandler.Handle
FileService.OnFileDetected

✅ التحقق

القيم وصلت كما هي (ID، الاسم، الحجم…)

🟠 System Test 6: PCAP Analysis Trigger
🧪 اسم الاختبار
TestSystem_PCAPAnalyzeFlow

🎯 الهدف

التأكد أن حدث التحليل يشغّل التحليل فعليًا.

📂 الملفات المشاركة
handlers/dispatcher.go
handlers/pcap_uploaded.go
services/pcap_uploaded.go
analysis/pcap_analyzer.go

🔧 الدوال المستخدمة
EventDispatcher.Dispatch
PCAPAnalyzeHandler.Handle
PCAPService.Analyze
GetFileHandle
RunFullAnalysis

✅ التحقق

عدم panic

تحليل بدأ (حتى لو جزئي)

🟡 System Test 7: Invalid Event Data
🧪 اسم الاختبار
TestSystem_InvalidEventData

🎯 الهدف

التأكد أن النظام لا ينهار عند JSON غير صالح.

📂 الملفات المشاركة
handlers/dispatcher.go

🔧 الدوال المستخدمة
EventDispatcher.Dispatch

✅ التحقق

error مُرجع

لا crash

🟡 System Test 8: End-to-End (Happy Path)
🧪 اسم الاختبار
TestSystem_EndToEnd_PCAPFlow

🎯 الهدف

محاكاة السيناريو الحقيقي كامل.

📂 الملفات المشاركة
client.go
handlers/*
services/*
analysis/*

🔧 الدوال المستخدمة
UploadFile
Dispatch
OnChunkReceived
reassemble
Analyze

✅ التحقق

الملف رُفع

دُمج

حُلّل

نُظف المؤقت

🧠 ملخص سريع جدًا
System Test	يغطي أي تدفق
FileUpload_Reassemble	Flow 1 + 2
MissingChunk	Flow 2
DuplicateChunk	Flow 2
EventRouting	Core Routing
FileDetected	Flow 3
PCAPAnalyze	Flow 4
InvalidEvent	Stability
EndToEnd	All Flows