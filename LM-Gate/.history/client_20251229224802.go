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
// NOTE: This command uploads a file using chunk-based upload.
var rootCmd = &cobra.Command{
	Use:   "LM <file>",
	Short: "LM - upload file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		if err := validatePath(path); err != nil {
			return err
		}

		fmt.Println("Uploading:", path)

		// NOTE: MockSender is used for local testing.
		// FIXME: Replace MockSender with real MQ sender (RabbitMQ).
		sender := &MockSender{}

		chunkSize := int64(5 * 1024 * 1024) // 5MB

		sent, err := UploadFile(path, chunkSize, sender)
		if err != nil {
			return err
		}

		fmt.Println("Chunks sent:", sent)
		fmt.Println("Done.")
		return nil
	},
}

// validatePath ensures the given path exists and points to a file.
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
func GenerateFileID(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s-%d", info.Name(), info.Size())
}

// Sender defines the interface for sending chunk messages.
type Sender interface {
	Send(msg ChunkMessage) error
}

// MockSender is a test sender used for local testing.
//
// NOTE: This does NOT send data over network.
type MockSender struct{}

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

// UploadFile orchestrates the full upload process.
//
// Steps:
// 1. Split file into chunks.
// 2. Send each chunk.
// 3. Send EOF message.
func UploadFile(path string, chunkSize int64, sender Sender) (int, error) {
	if sender == nil {
		return 0, errors.New("sender is nil")
	}

	fileID := GenerateFileID(path)
	if fileID == "" {
		return 0, errors.New("failed to generate file id")
	}

	chunks, errs := SplitFile(path, chunkSize)

	var allChunks [][]byte
	for c := range chunks {
		allChunks = append(allChunks, c)
	}

	if err := <-errs; err != nil {
		return 0, err
	}

	total := len(allChunks)
	sent := 0

	for i, data := range allChunks {
		msg := BuildChunkMessage(fileID, i, total, data)
		if err := sender.Send(msg); err != nil {
			return sent, err
		}
		sent++
	}

	if err := SendEOF(fileID, sender); err != nil {
		return sent, err
	}

	return sent, nil
}

// BuildChunkMessage builds a chunk message for one file piece.
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
func SendEOF(fileID string, sender Sender) error {
	msg := ChunkMessage{
		FileID: fileID,
		IsEOF:  true,
	}
	return sender.Send(msg)
}

func Execute() error {
	return rootCmd.Execute()
}
