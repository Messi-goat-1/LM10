
---

## 1. OnMessage(msg ChunkMessage)

`OnMessage` is the **main entry point** of the server.

It receives all incoming messages and controls the full lifecycle.

### Responsibilities:
- Validate incoming messages.
- Store chunks on disk.
- Detect end-of-file (EOF).
- Assemble the file.
- Send the file for analysis.
- Clean temporary data.

### Behavior:
- If the message is a normal chunk:
  - The chunk is stored using `StoreChunk`.
- If the message is an end-of-file signal (`IsEOF`):
  1. All chunks are assembled using `AssembleFile`.
  2. The completed file is analyzed using `ProcessFile`.
  3. Temporary chunks are deleted using `Cleanup`.

---

## 2. ValidateMessage(msg ChunkMessage)

A simple validation and security function.

### What it checks:
- The `FileID` must not be empty.
- Chunk data must exist unless the message is an EOF signal.

This prevents invalid or corrupted messages from being processed.

---

## 3. StoreChunk(msg ChunkMessage)

Handles **temporary storage** of file chunks.

### What it does:
- Creates a dedicated directory for each file inside `temp_chunks/`.
- Saves each chunk as a separate file named `part_X`.
- Uses disk storage instead of memory to avoid high RAM usage.

This design allows handling very large files safely.

---

## 4. IsFileComplete(fileID string)

A simple check function.

### Purpose:
- Verifies whether chunks already exist for a given file.
- Does NOT verify that all chunks are present.

This function only checks for presence, not completeness.

---

## 5. AssembleFile(fileID string)

Responsible for **rebuilding the original file**.

### Steps:
- Creates the final `.pcap` file inside the `uploads/` directory.
- Reads chunks sequentially (`part_0`, `part_1`, …).
- Merges chunk data into the final file.
- Stops when the next chunk in sequence is missing.

### Output:
- Returns the full path of the assembled file.

---

## 6. ProcessFile(fileID string, filePath string)

Connects the server to the analysis layer.

### What it does:
- Receives the path of the completed file.
- Sends it to the `analysis` package.
- Calls `AnalyzePCAP` to analyze the PCAP file.

The file is processed by path, not loaded into memory.

---

## 7. Cleanup(fileID string)

Handles server housekeeping.

### Purpose:
- Deletes the temporary directory (`temp_chunks/fileID`).
- Frees disk space after processing is complete.
- Prevents disk space leaks.

---

## Summary

- Chunks arrive → `StoreChunk` saves them.
- EOF arrives → `AssembleFile` rebuilds the file.
- Analysis runs → `ProcessFile` analyzes the file.
- Cleanup runs → `Cleanup` removes temporary data.

This design keeps memory usage low, supports large files,
and cleanly separates responsibilities.




ChunkMessage
   ↓
OnMessage
   ↓
ValidateMessage
   ↓
StoreChunk ──→ AssembleFile ──→ ProcessFile ──→ Cleanup
