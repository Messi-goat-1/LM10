# LM-Gate – System Summary (One-Page)

This document provides a **compact summary** of the entire system,
including the **4 main flows** and the **key functions** that drive each flow.

The system follows an **Event-Driven Architecture (EDA)**.

---

## Flow 1: Client Upload Flow
**Purpose:** Split a file and send it in chunks from the client.

**Main File:**
- client.go

**Key Functions:**
```go
rootCmd.RunE        // CLI entry point
validatePath        // Validate file path
SplitFile           // Split file into chunks
UploadFile          // Control upload process
BuildChunkMessage   // Build chunk message
Sender.Send         // Send chunk
SendEOF             // Signal end of file

Flow:
CLI → SplitFile → Send chunks → SendEOF

------------------------
Flow 2: File Chunk Ingestion Flow
Purpose: Receive chunks, store them, verify completion, and reassemble the file.

Main Files:

events/file_chunk.go

handlers/file_chunk.go

services/file_chunk.go

Key Functions:
EventDispatcher.Dispatch
FileChunkHandler.Handle
Manager.OnChunkReceived
Manager.isComplete
Manager.reassemble
Flow:
file.chunk event
→ Dispatcher
→ FileChunkHandler
→ OnChunkReceived
→ isComplete
→ reassemble → cleanup
----------------------------------------------
Flow 3: File Detected Flow

Purpose: Handle file metadata and business-level logic.

Main Files:

events/file_detected.go

handlers/file_detected.go

services/file_detected.go

Key Functions:

EventDispatcher.Dispatch
FileDetectedHandler.Handle
FileService.OnFileDetected

:

file.detected event
→ Dispatcher
→ FileDetectedHandler
→ FileService.OnFileDetected
-----------------------------
Flow 4: PCAP Analysis Flow

Purpose: Analyze a fully uploaded PCAP file.

Main Files:

events/pcap_uploaded.go

handlers/pcap_uploaded.go

services/pcap_uploaded.go

analysis/pcap_analyzer.go

Key Functions:

EventDispatcher.Dispatch
PCAPAnalyzeHandler.Handle
PCAPService.Analyze
pcap.OpenOffline
gopacket.NewPacketSource


Flow:

pcap.analyze event
→ Dispatcher
→ PCAPAnalyzeHandler
→ PCAPService.Analyze
→ PCAP packet analysis

Core Architecture Rules

Events = data only

Handlers = routing only

Services = real logic

Dispatcher = central router

Each flow is independent and extensible


----
Flow Overview Table
#	Flow Name	Responsibility
1	Client Upload	Send file chunks
2	Chunk Ingestion	Store + reassemble
3	File Detected	Business logic
4	PCAP Analysis	Network analysis