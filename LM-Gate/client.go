package lmgate

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type ChunkMessage struct {
	FileID  string
	ChunkID int
	Total   int
	Data    []byte
	IsEOF   bool
}

var (
	chunkStore        = make(map[string]map[int][]byte)
	ErrInvalidMessage = errors.New("invalid message")
	ErrMissingChunk   = errors.New("missing chunk")
)

// ======= Clint ==========
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

		sender := &MockSender{}
		chunkSize := int64(5 * 1024 * 1024) // 5MB

		if err := UploadFile(path, chunkSize, sender); err != nil {
			return err
		}

		fmt.Println("Upload finished successfully")
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

// GenerateFileID creates a stable identifier based on file name and size.
func GenerateFileID(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s-%d", info.Name(), info.Size())
}

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

type Sender interface {
	Send(msg ChunkMessage) error
}

// UploadFile orchestrates file upload by sending all chunks then EOF.
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

// BuildChunkMessage builds a chunk message for a single file piece.
func BuildChunkMessage(fileID string, chunkID int, total int, data []byte) ChunkMessage {
	return ChunkMessage{
		FileID:  fileID,
		ChunkID: chunkID,
		Total:   total,
		Data:    data,
		IsEOF:   false,
	}
}

// SendEOF notifies the receiver that all chunks have been sent.
func SendEOF(fileID string, sender Sender) error {
	msg := ChunkMessage{
		FileID: fileID,
		IsEOF:  true,
	}
	return sender.Send(msg)
}

func resetChunkStore() {
	chunkStore = make(map[string]map[int][]byte)
}

// تحسين مهم لازم اسوي دالة تاكد انه ترتيب وصل صح يعني من 0 الى رقم
// اسوي دالة تحقق من حجم ملف صح
