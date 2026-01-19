// داخل مجلد events
package events

type FileCollectionPayload struct {
	CollectionID string `json:"collection_id"`
	FileName     string `json:"file_name"`
	FinalPath    string `json:"final_path"`
	Status       string `json:"status"`
}
