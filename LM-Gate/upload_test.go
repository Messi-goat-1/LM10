package lmgate

import "testing"

func TestUploadAndAssemble(t *testing.T) {
	resetChunkStore()

	fakeData := []byte("HELLO WORLD THIS IS TEST DATA")
	fileID := "test-file"

	sender := &FakeSender{}

	// قسم البيانات يدويًا
	chunks := [][]byte{
		fakeData[:5],
		fakeData[5:10],
		fakeData[10:],
	}

	for i, c := range chunks {
		msg := BuildChunkMessage(fileID, i, len(chunks), c)
		if err := sender.Send(msg); err != nil {
			t.Fatal(err)
		}
	}

	// EOF
	if err := sender.Send(ChunkMessage{
		FileID: fileID,
		IsEOF:  true,
	}); err != nil {
		t.Fatal(err)
	}
}
