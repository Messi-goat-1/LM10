package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

/*
========================
CONFIG
========================
*/

func getAPIURL() string {
	return os.Getenv("LM_API_URL")
}

func getAPIKey() string {
	return os.Getenv("LM_API_KEY")
}

/*
========================
CLI HELP
========================
*/

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  lm upload <file.pcap>")
}

/*
========================
UPLOAD FLOW
========================
*/

func uploadFile(filePath string) error {
	file, err := openFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body, contentType, err := buildMultipartBody(
		"pcapfile",
		file,
		filepath.Base(filePath),
	)
	if err != nil {
		return fmt.Errorf("failed to build multipart body: %w", err)
	}

	apiURL := getAPIURL()
	apiKey := getAPIKey()

	if apiURL == "" {
		return fmt.Errorf("API URL is not configured (LM_API_URL)")
	}

	fmt.Println("Using API:", apiURL)

	req, err := createUploadRequest(apiURL, body, contentType, apiKey)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := sendRequest(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	result, err := readResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Println("Server response:")
	fmt.Println(result)
	return nil
}

func saveUploadedFile(file multipart.File, filename string, baseDir string) (string, error) {

	// تأكد من وجود المجلد
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	// حماية من path traversal
	safeName := filepath.Base(filename)
	dstPath := filepath.Join(baseDir, safeName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return dstPath, nil
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		http.Error(w, "invalid form", 400)
		return
	}

	file, header, err := r.FormFile("pcapfile")
	if err != nil {
		http.Error(w, "file missing", 400)
		return
	}
	defer file.Close()

	path, err := saveUploadedFile(
		file,
		header.Filename,
		"/data/uploads",
	)
	if err != nil {
		http.Error(w, "save failed", 500)
		return
	}

	fmt.Fprintf(w, "File saved: %s", path)
}

/*

========================
HELPER FUNCTIONS
========================
*/

// 1️⃣ فتح الملف فقط
func openFile(path string) (*os.File, error) {
	return os.Open(path)
}

// 2️⃣ بناء Multipart Body
func buildMultipartBody(fieldName string, file *os.File, fileName string) (*bytes.Buffer, string, error) {

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, "", err
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, "", err
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return &body, writer.FormDataContentType(), nil
}

// 3️⃣ إنشاء HTTP Request
func createUploadRequest(url string, body *bytes.Buffer, contentType string, apiKey string) (*http.Request, error) {

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-API-Key", apiKey)

	return req, nil
}

// 4️⃣ إرسال الطلب
func sendRequest(req *http.Request) (*http.Response, error) {
	fmt.Println("Uploading file, please wait...")
	client := &http.Client{}
	return client.Do(req)
}

// 5️⃣ قراءة الرد
func readResponse(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	return string(data), err
}

func RunUploadLogic() {
	if len(os.Args) != 3 || os.Args[1] != "upload" {
		printUsage()
		os.Exit(1)
	}

	filePath := os.Args[2]

	if err := uploadFile(filePath); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
