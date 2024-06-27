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
    filePath := flag.String("file", "", "Path to the file to upload")
    flag.Parse()

    if *filePath == "" {
        fmt.Println("File path is required. Use -file to specify the path to the file.")
        return
    }

    err := uploadFile(*filePath)
    if err != nil {
        fmt.Println("Error uploading file: ", err)
    } 
}

func uploadFile(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("Could not open file: %w", err)
    }
    defer file.Close()

    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    part, err := writer.CreateFormFile("_" + randomString(9) + ".png", filepath.Base(filePath))
    if err != nil {
        return fmt.Errorf("Could not create form file: %w", err)
    }

    _, err = io.Copy(part, file)
    if err != nil {
        return fmt.Errorf("Could not copy file content: %w", err)
    }
    writer.Close()

    url := "https://rehost.diberie.com/Host/UploadFiles?PrivateMode=false&SendMail=false&Comment="
    req, err := http.NewRequest("POST", url, &requestBody)
    if err != nil {
        return fmt.Errorf("could not create HTTP request: %w", err)
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("could not perform HTTP request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected response status: %s", resp.Status)
    }

    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("could not read response body: %w", err)
    }

    var uploadResponse UploadResponse
    err = json.Unmarshal(responseBody, &uploadResponse)
    if err != nil {
        return fmt.Errorf("could not parse JSON response: %w", err)
    }

    fmt.Println("https://rehost.diberie.com/Picture/Get/f/" + strconv.FormatInt(uploadResponse.PicID, 10))

    return nil
}

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
