🧪 Test Matrix (Exact Mapping)
🔴 1) Chunk Management – services
📂 الملف
services/file_chunk.go

🔧 الدالة
func (s *Manager) OnChunkReceived(...)


🧪 ملف الاختبار:

services/file_chunk_test.go


🏷️ أسماء الاختبارات:

TestManager_OnChunkReceived_WriteChunk
TestManager_OnChunkReceived_DuplicateChunk
TestManager_OnChunkReceived_CreateTempDir
TestManager_OnChunkReceived_TriggerReassemble

🔧 الدالة
func (s *Manager) isComplete(dir string, total int) bool


🧪 ملف الاختبار:

services/file_chunk_test.go


🏷️ أسماء الاختبارات:

TestManager_isComplete_AllChunksPresent
TestManager_isComplete_MissingChunk
TestManager_isComplete_ExtraFiles
TestManager_isComplete_WrongOrder

🔧 الدالة
func (s *Manager) reassemble(fileID string, totalChunks int) error


🧪 ملف الاختبار:

services/file_chunk_test.go


🏷️ أسماء الاختبارات:

TestManager_reassemble_Success
TestManager_reassemble_MissingChunk
TestManager_reassemble_CleanupTempDir
TestManager_reassemble_WrongOrder

🔴 2) Dispatcher – handlers
📂 الملف
handlers/dispatcher.go

🔧 الدالة
func (d *EventDispatcher) Dispatch(routingKey string, data []byte) error


🧪 ملف الاختبار:

handlers/dispatcher_test.go


🏷️ أسماء الاختبارات:

TestDispatcher_Dispatch_FileChunk
TestDispatcher_Dispatch_FileDetected
TestDispatcher_Dispatch_PCAPAnalyze
TestDispatcher_Dispatch_UnknownRoutingKey
TestDispatcher_Dispatch_InvalidJSON

🟠 3) Handlers
📂 الملف
handlers/file_chunk.go

🔧 الدالة
func (h *FileChunkHandler) Handle(event events.FileChunkEvent)


🧪 ملف الاختبار:

handlers/file_chunk_test.go


🏷️ أسماء الاختبارات:

TestFileChunkHandler_Handle_ForwardsToManager

📂 الملف
handlers/file_detected.go

🔧 الدالة
func (h *FileDetectedHandler) Handle(event events.FileDetectedEvent)


🧪 ملف الاختبار:

handlers/file_detected_test.go


🏷️ أسماء الاختبارات:

TestFileDetectedHandler_Handle_ForwardsMetadata

📂 الملف
handlers/pcap_uploaded.go

🔧 الدالة
func (h *PCAPAnalyzeHandler) Handle(event events.PCAPAnalyzeEvent) error


🧪 ملف الاختبار:

handlers/pcap_uploaded_test.go


🏷️ أسماء الاختبارات:

TestPCAPAnalyzeHandler_Handle_CallsAnalyze
TestPCAPAnalyzeHandler_Handle_ReturnsError

🟠 4) Events – serialization
📂 الملفات
events/event.go
events/file_chunk.go
events/file_detected.go
events/pcap_uploaded.go


🧪 ملف الاختبار:

events/events_test.go


🏷️ أسماء الاختبارات:

TestEvent_JSON_MarshalUnmarshal
TestFileChunkEvent_JSON_MarshalUnmarshal
TestFileDetectedEvent_JSON_MarshalUnmarshal
TestPCAPAnalyzeEvent_JSON_MarshalUnmarshal

🟡 5) PCAP Analysis
📂 الملف
services/pcap_uploaded.go

🔧 الدالة
func (s *PCAPService) Analyze(fileID string, filePath string) error


🧪 ملف الاختبار:

services/pcap_uploaded_test.go


🏷️ أسماء الاختبارات:

TestPCAPService_Analyze_ValidPCAP
TestPCAPService_Analyze_FileNotFound
TestPCAPService_Analyze_InvalidPCAP

📂 الملف
analysis/pcap_analyzer.go

🔧 الدوال
GetFileHandle
RunFullAnalysis


🧪 ملف الاختبار:

analysis/pcap_analyzer_test.go


🏷️ أسماء الاختبارات:

TestGetFileHandle_ValidFile
TestGetFileHandle_FileNotFound
TestRunFullAnalysis_ReadPackets

🟡 6) Client
📂 الملف
client.go

🔧 الدوال
SplitFile
UploadFile
BuildChunkMessage
SendEOF


🧪 ملف الاختبار:

client_test.go


🏷️ أسماء الاختبارات:

TestSplitFile_SmallFile
TestSplitFile_LargeFile
TestUploadFile_SendsAllChunks
TestUploadFile_SendsEOF
TestBuildChunkMessage_CorrectFields
TestSendEOF_CallsSender
