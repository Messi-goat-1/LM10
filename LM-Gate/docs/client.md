# client.go – Simple Explanation

This document explains how the client works
and how a file is prepared and sent to the server.

---

## 1. ChunkMessage Data Structure

`ChunkMessage` is the data container used before sending data.

It contains:

- `FileID`  
  A unique identifier for the file.

- `ChunkID`  
  The index of the current chunk (used to keep order).

- `Total`  
  Total number of chunks for the file.

- `Data`  
  The actual chunk data as a byte array.

- `IsEOF`  
  A flag that tells the server:
  “This is the last message, the file is complete.”

NOTE:  
This structure is shared between the client and the server.

TODO:  
Add checksum or hash to detect corrupted chunks.

---

## 2. Entry Point – rootCmd

The client uses the **Cobra** library to create a CLI tool.

Example usage:


What happens:
- The user provides a file path.
- Chunk size is fixed to 5 MB.
- `UploadFile` is called to start the upload process.

NOTE:  
This design makes the client easy to use from the terminal.

TODO:  
Allow chunk size to be passed as a CLI flag.

---

## 3. SplitFile (File Splitting)

`SplitFile` is a technical function that runs in the background
using a goroutine.

What it does:
- Opens the file.
- Reads the file chunk by chunk using `chunkSize`.
- Sends each chunk immediately through a channel.
- Stops when it reaches the end of the file (`io.EOF`).

NOTE:  
Channels allow reading and sending data efficiently.

FIXME:  
Chunks are later collected fully in memory before sending.

---

## 4. GenerateFileID

This function generates a unique identifier for the file.

How it works:
- Combines the original file name and file size.

Why this is important:
- Helps the server know which chunks belong to the same file.

NOTE:  
Simple and fast approach.

TODO:  
Use a hash (e.g. SHA256) for stronger uniqueness.

---

## 5. UploadFile (The Orchestrator)

`UploadFile` controls the full upload process.

Steps:
1. Generate the file ID.
2. Split the file using `SplitFile`.
3. Loop over all chunks.
4. Send each chunk using `BuildChunkMessage`.
5. Send an end-of-file signal using `SendEOF`.

NOTE:  
This function is the “brain” of the client upload logic.

FIXME:  
All chunks are loaded into memory before sending.

TODO:  
Stream chunks directly to the sender without storing them all.

---

## 6. BuildChunkMessage and SendEOF

### BuildChunkMessage
- Wraps chunk data inside a `ChunkMessage`.
- Sets the correct chunk number.
- `IsEOF` is always `false`.

### SendEOF
- Sends a message with `IsEOF = true`.
- Tells the receiver that the file upload is finished.

NOTE:  
No file data is sent with EOF.

---

## How the Whole System Works (Client → Rabbit → Server)

- **Client (`client.go`)**  
  Takes a PCAP file, splits it into 5 MB chunks, and sends them.

- **Broker (`rabbit.go`)**  
  (If used as a sender) receives chunks and places them into
  a RabbitMQ queue.

- **Server (`server.go`)**  
  Receives chunks from the queue, stores them temporarily,
  and when `IsEOF` arrives:
  - Assembles the file.
  - Analyzes the PCAP.
  - Cleans up temporary data.

---

## Summary

- The client prepares and splits the file.
- Chunks are sent in order.
- EOF marks completion.
- The system is modular and easy to extend.

TODO:
- Validate chunk order on the client.
- Verify final file size.
- Add retry logic for failed sends.





File on Disk
   ↓
validatePath
   ↓
SplitFile
   ↓
UploadFile
   ↓
BuildChunkMessage
   ↓
Sender.Send (MockSender)
   ↓
SendEOF
