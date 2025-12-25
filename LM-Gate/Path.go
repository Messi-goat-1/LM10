package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

//========= Clint and server +=======

type ChunkMessage struct {
	FileID  string
	ChunkID int
	Total   int
	Data    []byte
	IsEOF   bool
}

var (
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

// توليد ID للملف
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

// رفع الملف (orchestrator)
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

// بناءرسالة
func BuildChunkMessage(fileID string, chunkID int, total int, data []byte) ChunkMessage {
	return ChunkMessage{
		FileID:  fileID,
		ChunkID: chunkID,
		Total:   total,
		Data:    data,
		IsEOF:   false,
	}
}

// ارسال نهاية الملف
func SendEOF(fileID string, sender Sender) error {
	msg := ChunkMessage{
		FileID: fileID,
		IsEOF:  true,
	}
	return sender.Send(msg)
}

var chunkStore = make(map[string]map[int][]byte)

func resetChunkStore() {
	chunkStore = make(map[string]map[int][]byte)
}

//===== Server ====

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

// التحقق من الرسالة
func ValidateMessage(msg ChunkMessage) error {
	if msg.FileID == "" {
		return ErrInvalidMessage
	}

	if !msg.IsEOF && len(msg.Data) == 0 {
		return ErrInvalidMessage
	}

	return nil
}

// تخزين موقت لل chunks
func StoreChunk(msg ChunkMessage) error {
	if _, ok := chunkStore[msg.FileID]; !ok {
		chunkStore[msg.FileID] = make(map[int][]byte)
	}
	chunkStore[msg.FileID][msg.ChunkID] = msg.Data
	return nil
}

// التحقق من اكتمال الملف
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

// دمج الملف
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

// التحليل - معالجة
func ProcessFile(fileID string, data []byte) error {
	fmt.Printf("Processing file %s (size=%d bytes)\n", fileID, len(data))
	return nil
}

// تنظيف الذاكرة
func Cleanup(fileID string) {
	delete(chunkStore, fileID)
}

// ---------------
type FakeSender struct{}

func (f *FakeSender) Send(msg ChunkMessage) error {
	return OnMessage(msg)
}
func main() {
	sender := &FakeSender{}
	err := UploadFile("test.pcap", 1024*1024, sender)
	if err != nil {
		panic(err)
	}
}
