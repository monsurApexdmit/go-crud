package utils

import (
	"fmt"
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
