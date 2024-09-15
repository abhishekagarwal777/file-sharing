package files

import (
    "fmt"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
    "time"
)

// SaveFileLocally saves the uploaded file to the local file system
func SaveFileLocally(file multipart.File, header *multipart.FileHeader, localPath string) (string, error) {
    defer file.Close()

    // Create a unique file name
    fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
    filePath := filepath.Join(localPath, fileName)

    // Create directory if not exists
    if err := os.MkdirAll(localPath, os.ModePerm); err != nil {
        return "", fmt.Errorf("failed to create directory: %v", err)
    }

    // Create a file on disk
    outFile, err := os.Create(filePath)
    if err != nil {
        return "", fmt.Errorf("could not create file: %v", err)
    }
    defer outFile.Close()

    // Write the uploaded file to disk
    _, err = io.Copy(outFile, file)
    if err != nil {
        return "", fmt.Errorf("could not save file: %v", err)
    }

    // Return the relative file path
    return fileName, nil
}
