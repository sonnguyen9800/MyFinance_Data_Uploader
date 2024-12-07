package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var action string

func init() {
	flag.StringVar(&action, "action", "", "action is required, either --upload or --download")
	flag.Parse()
	if action != "upload" && action != "download" {
		fmt.Println("action is required, either --upload or --download")
		os.Exit(1)
	}
}

func main() {
	var filePath string
	if len(os.Args) > 2 {
		filePath = os.Args[2]
	} else {
		fmt.Println("No filepath is provided")
		os.Exit(1)
	}

	if action == "upload" {
		DemoUploadCSV(filePath)
	} else if action == "download" {
		DemoDownloadCSV(filePath)
	}
}
func DemoDownloadCSV(filePath string) {
	apiURL := "http://localhost:8080/api/expenses/download"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFiY0BnbWFpbC5jb20iLCJleHAiOjE3MzM2NTg3OTcsInJvbGUiOiJ1c2VyIiwidXNlcl9pZCI6IjY3NGViZDZjODgzNDU4MTJkOThiNmEwYiJ9.UDpz35nx58XflMzie29glqy6J6m-gQWyULIDXz4OVZI"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		return
	}

	// Save the file to disk
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("CSV file saved to", filePath)
}

func DemoUploadCSV(filePath string) {
	apiURL := "http://localhost:8080/api/expenses/upload"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFiY0BnbWFpbC5jb20iLCJleHAiOjE3MzM2NTg3OTcsInJvbGUiOiJ1c2VyIiwidXNlcl9pZCI6IjY3NGViZDZjODgzNDU4MTJkOThiNmEwYiJ9.UDpz35nx58XflMzie29glqy6J6m-gQWyULIDXz4OVZI"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create a buffer and a multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file to the multipart form data
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Error copying file data: %v\n", err)
		return
	}

	// Close the writer to finalize the multipart form
	err = writer.Close()
	if err != nil {
		fmt.Printf("Error closing writer: %v\n", err)
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Add the Content-Type header with the multipart boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Print the response
	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Response body: %s\n", string(respBody))
}
