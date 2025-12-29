# Workflow Summary – rabbit.go

This document explains how `rabbit.go` works and how RabbitMQ
is used for sending and receiving messages.

---

## 1. RabbitClient Structure

`RabbitClient` represents the client that communicates with RabbitMQ.

### Fields:
- `conn`  
  Represents the actual TCP connection to the RabbitMQ server.

- `channel`  
  The channel used to send and receive messages.  
  Channels are lighter and cheaper than full connections.

---

## 2. NewRabbitClient Function

This is the constructor function.

### What it does:
- Takes the RabbitMQ server URL.
- Opens a connection to the server.
- Opens a channel on that connection.
- Returns a ready-to-use `RabbitClient` instance.

---

## 3. Close Function

Cleanup function.

### Purpose:
- Safely closes the channel and the connection.
- Prevents resource leaks (memory or open connections).
- Should be called when the server shuts down.

---

## 4. PublishMessage Function

Responsible for **sending messages**.

### Steps:
- Ensures the queue exists using `QueueDeclare`.
- Creates the queue if it does not exist.
- Sends a text message to the specified queue using `PublishWithContext`.

---

## 5. RunHeartbeat Function

Heartbeat (monitoring) function.

### How it works:
- Runs automatically every 1 minute.
- Sends a message to a queue called `server_status`.
- The message indicates that the server is still running and includes the current time.

### Why it is useful:
- Helps monitor server health remotely.
- Confirms that the server connection is alive.

---

## 6. ConsumeMessages Function

Responsible for **receiving messages**.

### What it does:
- Starts listening to a specific queue (`queueName`).
- Uses a goroutine so the program does not block while waiting for messages.
- When a message arrives, its raw data (`d.Body`) is passed to a processing function (`processor`).

The processing logic is defined by the user of this function.

---

## Overall Flow

- **Connection**:  
  Create the client using `NewRabbitClient`.

- **Sending**:  
  Use `PublishMessage` to send data.

- **Monitoring**:  
  `RunHeartbeat` runs in the background to confirm the server is alive.

- **Receiving**:  
  `ConsumeMessages` listens for incoming messages from other systems.

---

## Summary

This design keeps messaging logic simple and clean:
- One client manages the connection and channel.
- Sending and receiving are clearly separated.
- Background tasks (heartbeat, consuming) do not block the main application.
