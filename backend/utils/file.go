package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveUploadedFile saves a multipart file to the uploads directory and returns the file URL or an error
func SaveUploadedFile(file multipart.File, handler *multipart.FileHeader, uploadDir string) (string, error) {
	// Ensure uploads directory exists
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		return "", err
	}

	fileURL := "/uploads/" + filename
	return fileURL, nil
}
