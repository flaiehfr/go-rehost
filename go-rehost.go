package main

import (
    "bytes"
    "flag"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "math/rand"
    "encoding/json"
    "time"
    "strconv"
)

type UploadResponse struct {
    PicID int64 `json:"picID"`
}

func main() {
    // Define and parse command-line arguments
    filePath := flag.String("file", "", "Path to the file to upload")
    cookieVal := flag.String("cookie", "", "Rehost cookie value")
    flag.Parse()

    if *filePath == "" {
        fmt.Println("File path is required. Use -file to specify the path to the file.")
        return
    }

    var cookie string
    if *cookieVal != "" {
        cookie = ".AspNet.ApplicationCookie=" + *cookieVal
    }
    fmt.Println(cookie)

    err := uploadFile(*filePath, cookie)
    if err != nil {
        fmt.Println("Error uploading file:", err)
    } 
}

func uploadFile(filePath string, cookie string) error {
    // Open the file
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("could not open file: %w", err)
    }
    defer file.Close()

    // Create a buffer to store the file and form data
    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    // Create a form file field
    part, err := writer.CreateFormFile("_"+randomString(9)+".png", filepath.Base(filePath))
    if err != nil {
        return fmt.Errorf("could not create form file: %w", err)
    }

    // Copy the file content to the form file field
    _, err = io.Copy(part, file)
    if err != nil {
        return fmt.Errorf("could not copy file content: %w", err)
    }

    // Close the writer to finalize the form data
    writer.Close()

    // Create the HTTP request
    url := "https://rehost.diberie.com/Host/UploadFiles?PrivateMode=false&SendMail=false&Comment="
    req, err := http.NewRequest("POST", url, &requestBody)
    if err != nil {
        return fmt.Errorf("could not create HTTP request: %w", err)
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())
    req.Header.Set("Cookie", cookie)

    // Perform the HTTP request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("could not perform HTTP request: %w", err)
    }
    defer resp.Body.Close()

    // Check for HTTP response status
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected response status: %s", resp.Status)
    }

    // Read and parse the response body
    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("could not read response body: %w", err)
    }

    // Unmarshal the JSON response
    var uploadResponse UploadResponse
    err = json.Unmarshal(responseBody, &uploadResponse)
    if err != nil {
        return fmt.Errorf("could not parse JSON response: %w", err)
    }

    // Output the constructed URL
    fmt.Println("https://rehost.diberie.com/Picture/Get/f/" + strconv.FormatInt(uploadResponse.PicID, 10))

    return nil
}

// Helper function to generate a random string
func randomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    rand.Seed(time.Now().UnixNano())
    var sb strings.Builder
    for i := 0; i < length; i++ {
        index := rand.Intn(len(charset))
        sb.WriteByte(charset[index])
    }
    return sb.String()
}
