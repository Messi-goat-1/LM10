# File Chunk Processing – Layered Architecture Explanation

This document explains how file chunks are handled
across the system layers.

---

## 1. Events Layer (package events)

This layer defines the **message format**
that is exchanged across the system (for example via RabbitMQ).

### FileChunkPayload
Contains the raw chunk data and metadata:
- `Data`  
  The raw bytes of the chunk.
- `FileID`  
  The original file identifier.
- `ChunkIndex`  
  The order of this chunk.
- `TotalChunks`  
  Total number of expected chunks.

### FileChunkEvent
This is the outer wrapper of the message.
It contains:
- The payload (`FileChunkPayload`)
- A `Timestamp` indicating when the chunk was sent.

NOTE:  
This layer defines data only, not logic.

---

## 2. Handlers Layer (package handlers)

This layer works as a **receiving bridge only**.

### NewFileChunkHandler
- Prepares the handler.
- Connects it to the service manager.

### Handle
- Triggered when a file chunk arrives.
- Receives the event.
- Extracts all data from the event:
  - File ID
  - Chunk order
  - Total chunks
  - Raw data
- Immediately forwards the data to the services layer.

NOTE:  
Handlers do not perform processing logic.
They only route data.

---

## 3. Services Layer (package services)

This layer is the **core logic** of the system
and handles heavy operations.

### NewManager
- Constructor function.
- Defines:
  - Temporary directory path (`temp_chunks`)
  - Final storage directory path (`uploads`)

---

### OnChunkReceived

When a chunk is received:
- Creates a file-specific temporary directory if it does not exist.
- Writes the chunk data to a separate file on disk:
  - `part_0`, `part_1`, etc.
- Calls `isComplete` to check if all chunks have arrived.

If all chunks are received:
- Starts the reassembly process (`reassemble`) in the background.

---

### isComplete
- Compares the number of chunk files stored in the temp directory
  with the expected `TotalChunks`.
- Returns true when all chunks are present.

---

### reassemble

When the file is complete:
- Opens a new output file in the final storage directory.
- Merges all chunk files in the correct order:
  - `0, 1, 2, ...`
- Deletes the temporary directory (cleanup) to free disk space.

---

## Function Summary by Layer

| Layer     | Function Name      | Main Responsibility |
|----------|-------------------|---------------------|
| Events   | FileChunkEvent     | Define chunk data structure |
| Handlers | Handle              | Receive event and route data |
| Services | OnChunkReceived     | Store chunk and check completion |
| Services | isComplete          | Verify all chunks arrived |
| Services | reassemble          | Merge chunks and cleanup |

---

## Summary

- Events define the message format.
- Handlers receive and forward events.
- Services store chunks, verify completeness, and merge files.
- Temporary data is cleaned after successful assembly.
- The design follows clear separation of responsibilities.
