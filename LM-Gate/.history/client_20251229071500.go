package lmgate

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// ChunkMessage represents a single file chunk message.
// It is sent from the client to the server.
//
// NOTE: The same structure is used on both client and server.
// TODO: Add checksum or hash field for data integrity.
type ChunkMessage struct {
	FileID  string
	ChunkID int
	Total   int
	Data    []byte
	IsEOF   bool
}

var (
	// chunkStore is used for temporary chunk storage (testing or future use).
	// NOTE: Currently not used in the upload flow.
	chunkStore        = make(map[string]map[int][]byte)
	ErrInvalidMessage = errors.New("invalid message")
	ErrMissingChunk   = errors.New("missing chunk")
)

// ================= Client =================

// rootCmd defines the CLI command: LM <file>
//
// NOTE: This command uploads a PCAP file using chunk-based upload.
// TODO: Add flags for chunk size and server address.
var rootCmd = &cobra.Command{
	Use:   "LM <file>",
	Short: "LM - upload PCAP file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		if err := validatePath(path); err != nil {
			return err
		}

		fmt.Println("Uploading:", path)

		// NOTE: MockSender is used for local testing.
		// FIXME: Replace MockSender with real network sender.
		sender := &MockSender{}

		chunkSize := int64(5 * 1024 * 1024) // 5MB
		// TODO: Make chunk size configurable via CLI flag.

		if err := UploadFile(path, chunkSize, sender); err != nil {
			return err
		}

		fmt.Println("Upload finished successfully")
		return nil
	},
}

// validatePath ensures the given path exists and points to a file.
//
// NOTE: Prevents uploading directories or invalid paths.
func validatePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return errors.New("file does not exist")
	}

	if info.IsDir() {
		return errors.New("path is a directory, not a file")
	}

	return nil
}

// SplitFile reads a file and streams it as fixed-size byte chunks.
//
// NOTE: Uses a goroutine to read the file asynchronously.
// FIXME: Entire file is later stored in memory in UploadFile.
func SplitFile(path string, chunkSize int64) (<-chan []byte, <-chan error) {
	chunks := make(chan []byte)
	errs := make(chan error, 1)

	go func() {
		defer close(chunks)
		defer close(errs)

		file, err := os.Open(path)
		if err != nil {
			errs <- err
			return
		}
		defer file.Close()

		for {
			buf := make([]byte, chunkSize)
			n, err := file.Read(buf)

			if n > 0 {
				chunks <- buf[:n]
			}

			if err == io.EOF {
				return
			}

			if err != nil {
				errs <- err
				return
			}
		}
	}()

	return chunks, errs
}

// GenerateFileID creates a stable file identifier.
//
// NOTE: Uses file name and file size.
// TODO: Use hash (SHA256) for stronger uniqueness.
func GenerateFileID(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s-%d", info.Name(), info.Size())
}

// MockSender is a test sender used for local testing.
//
// NOTE: This does NOT send data over network.
type MockSender struct{}

// Send simulates sending a chunk message.
//
// NOTE: Prints message details to stdout.
// FIXME: No real transmission or error handling.
func (m *MockSender) Send(msg ChunkMessage) error {
	if msg.IsEOF {
		fmt.Println("[SEND] EOF for file:", msg.FileID)
		return nil
	}

	fmt.Printf(
		"[SEND] file=%s chunk=%d/%d size=%d\n",
		msg.FileID,
		msg.ChunkID+1,
		msg.Total,
		len(msg.Data),
	)
	return nil
}

// Sender defines the interface for sending chunk messages.
//
// NOTE: Allows swapping MockSender with real sender (RabbitMQ, HTTP, etc.).
type Sender interface {
	Send(msg ChunkMessage) error
}

// UploadFile orchestrates the full upload process.
//
// Steps:
// 1. Split file into chunks.
// 2. Send each chunk.
// 3. Send EOF message.
//
// FIXME: All chunks are loaded into memory before sending.
func UploadFile(path string, chunkSize int64, sender Sender) error {
	if sender == nil {
		return errors.New("sender is nil")
	}

	fileID := GenerateFileID(path)
	if fileID == "" {
		return errors.New("failed to generate file id")
	}

	chunks, errs := SplitFile(path, chunkSize)

	var allChunks [][]byte
	for c := range chunks {
		allChunks = append(allChunks, c)
	}

	if err := <-errs; err != nil {
		return err
	}

	total := len(allChunks)

	for i, data := range allChunks {
		msg := BuildChunkMessage(fileID, i, total, data)
		if err := sender.Send(msg); err != nil {
			return err
		}
	}

	return SendEOF(fileID, sender)
}

// BuildChunkMessage builds a chunk message for one file piece.
//
// NOTE: IsEOF is always false here.
func BuildChunkMessage(fileID string, chunkID int, total int, data []byte) ChunkMessage {
	return ChunkMessage{
		FileID:  fileID,
		ChunkID: chunkID,
		Total:   total,
		Data:    data,
		IsEOF:   false,
	}
}

// SendEOF sends the end-of-file signal to the receiver.
//
// NOTE: Indicates that all chunks were sent.
func SendEOF(fileID string, sender Sender) error {
	msg := ChunkMessage{
		FileID: fileID,
		IsEOF:  true,
	}
	return sender.Send(msg)
}

// resetChunkStore clears the temporary chunk storage.
//
// NOTE: Currently unused.
func resetChunkStore() {
	chunkStore = make(map[string]map[int][]byte)
}

// TODO: Add function to validate chunk order (0 → total-1).
// TODO: Add function to verify final file size after upload.
