# File Detected – Event Driven Architecture Explanation

This file contains three packages that work together using an
**Event-Driven Architecture** pattern.

The purpose of this code is to detect a new file and process it
through multiple layers.

Below is a simple explanation of how the **event** is created
and how it moves between the three layers.

---

## 1. Events Layer (package events)

This layer defines the **shape of the data** that will be passed
between the system components.

### FileDetectedPayload
This structure is a container that holds technical file details, such as:
- FileID
- File name
- File size (in bytes)
- File type
- Checksum (used to verify file integrity)

### FileDetectedEvent
This is the main event structure.
It wraps the payload and adds a **timestamp** to record
when the file was detected.

NOTE:  
This layer does NOT contain executable logic.  
It only defines data structures (structs).

Structures in this layer:
- `FileDetectedPayload` – holds file technical details.
- `FileDetectedEvent` – combines payload with the detection time.

---

## 2. Handlers Layer (package handlers)

This layer acts as a **middleman** or **dispatcher**.

When the system detects a file, the event is sent to a handler.

The handler does NOT process the data itself.
Its job is to:
- Receive the event.
- Extract the payload data.
- Forward the data to the service layer.

Functions in this layer:
- `NewFileDetectedHandler`  
  Constructor function that creates the handler and connects it to the file service.

- `Handle`  
  The main function that receives a `FileDetectedEvent`,
  extracts its payload, and passes the data to the service layer.

---

## 3. Services Layer (package services)

This is where the **real work (business logic)** happens.

### OnFileDetected
This function receives all file details:
- File ID
- Name
- Size
- Type
- Checksum

Currently, it only prints the data in a structured way.

In a real system, this is where you would:
- Store file information in a database.
- Check for duplicate files using the checksum.
- Move the file to its final storage location.

Functions in this layer:
- `NewFileService`  
  Constructor function that creates a new file service.

- `OnFileDetected`  
  The main business logic function that processes the detected file.

---

## Event Workflow (Step by Step)

1. **Trigger**  
   The system detects a new file and creates a `FileDetectedEvent`.

2. **Dispatch**  
   The event is sent to `FileDetectedHandler`.

3. **Forwarding**  
   The handler extracts the payload and forwards the data to `FileService`.

4. **Execution**  
   The file service executes the required actions
   (printing, storing, validating, etc.).

---

## Summary

- Events define the data.
- Handlers route the event.
- Services execute the business logic.
- Each layer has a single clear responsibility.
- This design keeps the system clean, modular, and easy to extend.
