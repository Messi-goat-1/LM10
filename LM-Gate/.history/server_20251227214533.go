package lmgate

import (
	"errors"
	"fmt"
)

// OnMessage is the server entry point for handling incoming chunk messages.
func OnMessage(msg ChunkMessage) error {
	if err := ValidateMessage(msg); err != nil {
		return err
	}

	if msg.IsEOF {
		if !IsFileComplete(msg.FileID) {
			return ErrMissingChunk
		}

		data, err := AssembleFile(msg.FileID)
		if err != nil {
			return err
		}

		if err := ProcessFile(msg.FileID, data); err != nil {
			return err
		}

		Cleanup(msg.FileID)
		return nil
	}

	return StoreChunk(msg)
}

// ValidateMessage performs basic validation on incoming messages.
func ValidateMessage(msg ChunkMessage) error {
	if msg.FileID == "" {
		return ErrInvalidMessage
	}

	if !msg.IsEOF && len(msg.Data) == 0 {
		return ErrInvalidMessage
	}

	return nil
}

// StoreChunk temporarily stores a file chunk in memory.
func StoreChunk(msg ChunkMessage) error {
	if _, ok := chunkStore[msg.FileID]; !ok {
		chunkStore[msg.FileID] = make(map[int][]byte)
	}
	chunkStore[msg.FileID][msg.ChunkID] = msg.Data
	return nil
}

// IsFileComplete checks whether all chunks for a file are present.
func IsFileComplete(fileID string) bool {
	chunks, ok := chunkStore[fileID]
	if !ok {
		return false
	}

	for i := 0; i < len(chunks); i++ {
		if _, ok := chunks[i]; !ok {
			return false
		}
	}
	return true
}

// AssembleFile reconstructs the original file from stored chunks.
func AssembleFile(fileID string) ([]byte, error) {
	chunks, ok := chunkStore[fileID]
	if !ok {
		return nil, errors.New("file not found")
	}

	var result []byte
	for i := 0; i < len(chunks); i++ {
		result = append(result, chunks[i]...)
	}
	return result, nil
}

// ProcessFile handles the fully assembled file (analysis, parsing, etc).
func ProcessFile(fileID string, data []byte) error {
	fmt.Printf("Processing file %s (size=%d bytes)\n", fileID, len(data))
	return nil
}

// Cleanup removes all stored data related to a file.
func Cleanup(fileID string) {
	delete(chunkStore, fileID)
}

// ---------------
type FakeSender struct{}

func (f *FakeSender) Send(msg ChunkMessage) error {
	return OnMessage(msg)
}
