# PCAP Uploaded – Event Workflow Explanation

This document explains how the **PCAP upload event**
works across different layers of the system.

---

## 1. Events Layer (package events)

This layer defines the **shape of the data**
that moves through the system.

### FileChunkPayload
Contains the basic information of a file chunk:
- `FileID`  
  The original file identifier.
- `ChunkIndex`  
  The order index of the chunk.
- `TotalChunks`  
  Total number of expected chunks (used to ensure no chunk is missing).
- `Data`  
  Raw chunk data.

### FileChunkEvent
Wraps the payload and adds a `Timestamp`
to record when the chunk arrived.

---

## 2. Handlers Layer (package handlers)

This layer acts as the **event receiver**.

Usually, the event arrives from a message broker
such as RabbitMQ (as shown in `rabbit.go`).

### Handle Function
- Receives the event as soon as it arrives.
- Extracts the payload data.
- Does NOT apply business logic.
- Forwards the extracted data to the services layer.

NOTE:  
Handlers are responsible only for routing, not processing.

---

## 3. Services Layer (package services)

This is where the **actual logic and file handling** happens.

### OnChunkReceived Function

When a chunk is received:
- Creates a temporary directory for the file (`FileID`) if it does not exist.
- Saves the chunk data into a separate file:
  - `part_0`, `part_1`, etc.
- Calls a check function to see if all chunks are received.

---

### isComplete Function

- Compares the number of stored chunk files
  with the expected `TotalChunks`.
- Returns true when all chunks are present.

---

### reassemble Function

When the file is complete:
- Opens a final output file.
- Merges all chunk files in correct order.
- Saves the final file in the `uploads` directory.
- Cleans up the temporary directory by deleting all chunk files.

NOTE:  
This step ensures correct order and data integrity.

---

## Event Workflow Summary

1. **Receive**  
   `FileChunkEvent` is received through `rabbit.go`.

2. **Route**  
   The handler forwards the data to the service layer.

3. **Store**  
   The service saves the chunk to  
   `temp_chunks/FileID/part_X`.

4. **Complete Check**  
   If the number of chunks equals `TotalChunks`:
   - All chunks are merged into one file in `uploads/`.
   - Temporary chunk files are removed using cleanup.

---

## Summary

- Events describe incoming chunk data.
- Handlers route the events.
- Services store, verify, and merge chunks.
- Temporary data is cleaned after completion.
- The system ensures correct order and full file reconstruction.
