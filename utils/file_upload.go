package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func SaveUploadedFile(c *gin.Context, file *multipart.FileHeader, destFolder string) (string, error) {
	// Create the destination folder if it doesn't exist
	if _, err := os.Stat(destFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(destFolder, 0755); err != nil {
			return "", err
		}
	}

	// Generate a unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(destFolder, filename)

	// Save the file
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}

	return dst, nil
}

// SaveUploadedFileTemp stores the upload in a temp location and returns
// (tempPath, finalPath). Use CommitUploadedFile to move into finalPath
// after DB changes succeed.
func SaveUploadedFileTemp(c *gin.Context, file *multipart.FileHeader, destFolder string) (string, string, error) {
	// Ensure destination folder exists for the final move.
	if _, err := os.Stat(destFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(destFolder, 0755); err != nil {
			return "", "", err
		}
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	finalPath := filepath.Join(destFolder, filename)

	tmpDir := filepath.Join(os.TempDir(), "go-crud-uploads")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return "", "", err
	}
	tempPath := filepath.Join(tmpDir, filename)

	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		return "", "", err
	}

	return tempPath, finalPath, nil
}

// CommitUploadedFile moves the temp file to its final destination.
func CommitUploadedFile(tempPath, finalPath string) error {
	if err := os.Rename(tempPath, finalPath); err == nil {
		return nil
	}

	// Fallback for cross-device rename.
	src, err := os.Open(tempPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(finalPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return os.Remove(tempPath)
}
