package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
	"github.com/schollz/progressbar/v3"
)

const (
	green = "\033[32m"
	red   = "\033[31m"
	reset = "\033[0m"
)

func main() {
	// Define flags
	maxDays := flag.Int("max-days", 0, "Maximum number of days to keep file")
	maxDownloads := flag.Int("max-downloads", 0, "Maximum number of times that file can be downloaded")
	flag.Parse()

	// Check for filename argument
	if flag.NArg() < 1 {
		fmt.Println("Filename is required")
		flag.Usage()
		return
	}

	filename := flag.Arg(0)

	// Read URL, HTTP auth username, and password from environment variables
	url := os.Getenv("TSH_URL")
	if url == "" {
		fmt.Println("Environment variable TSH_URL must be set")
		return
	}
	user := os.Getenv("TSH_HTTP_AUTH_USER")
	pass := os.Getenv("TSH_HTTP_AUTH_PASS")

	uploadFile(filename, *maxDays, *maxDownloads, url, user, pass)
}

func uploadFile(filename string, maxDays, maxDownloads int, url, user, pass string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Failed to get file info:", err)
		return
	}

	// Set up progress bar
	bar := progressbar.NewOptions64(
		fileInfo.Size(),
		progressbar.OptionSetDescription(filename),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionClearOnFinish(),
	)

	// Create a buffer to store the multipart data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field
	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		fmt.Println("Failed to create form file:", err)
		return
	}

	// Create a reader to read the file and update the progress bar
	reader := io.TeeReader(file, bar)

	// Copy the file data to the form file field
	_, err = io.Copy(part, reader)
	if err != nil {
		fmt.Println("Failed to copy file data:", err)
		return
	}

	// Set headers for max days and downloads if provided
	if maxDays > 0 {
		writer.WriteField("Max-Days", strconv.Itoa(maxDays))
	}
	if maxDownloads > 0 {
		writer.WriteField("Max-Downloads", strconv.Itoa(maxDownloads))
	}

	// Close the multipart writer
	writer.Close()

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Add HTTP basic authentication header if username and password are provided
	if user != "" && pass != "" {
		auth := user + ":" + pass
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	}

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to execute HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		dURL := string(body)
		// Print the download URL and copy to clipboard
		fmt.Printf("👍 Download from here: "+green+"%s"+reset+"\n", dURL)
		clipboard.WriteAll(dURL)
		fmt.Println("😉 It has also been copied to the clipboard!" + reset)
	} else {
		// Print the error message
		fmt.Printf("❌ Failed to upload file: "+red+"%s"+reset+"\n", string(body))
	}
}
