package services

import (
	"log"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
	"go-crud/utils"
)

// tempFile is the pending-upload token. The file exists on disk at tempPath
// and will be moved to finalPath only after the DB transaction commits.
type tempFile struct {
	tempPath  string
	finalPath string
}

// FileService wraps utils file-upload helpers and os.Remove so that no
// other layer touches the filesystem directly for product uploads.
type FileService struct{}

func NewFileService() FileService {
	return FileService{}
}

// SaveTemp writes a multipart file to a temp location and returns the pending token.
// *gin.Context is required because utils.SaveUploadedFileTemp calls c.SaveUploadedFile.
func (fs FileService) SaveTemp(c *gin.Context, file *multipart.FileHeader, destFolder string) (tempFile, error) {
	tmp, final, err := utils.SaveUploadedFileTemp(c, file, destFolder)
	if err != nil {
		return tempFile{}, err
	}
	return tempFile{tempPath: tmp, finalPath: final}, nil
}

// Commit moves a single temp file to its final destination.
func (fs FileService) Commit(tf tempFile) error {
	return utils.CommitUploadedFile(tf.tempPath, tf.finalPath)
}

// CommitAll commits every pending file in order. On the first failure it
// cleans up all remaining temp files and returns the error.
// Returns the slice of final paths (same order as input) on success.
func (fs FileService) CommitAll(pending []tempFile) ([]string, error) {
	finalPaths := make([]string, 0, len(pending))
	for i, tf := range pending {
		if err := utils.CommitUploadedFile(tf.tempPath, tf.finalPath); err != nil {
			// clean up: remove this temp file and all subsequent ones
			_ = os.Remove(tf.tempPath)
			for _, remaining := range pending[i+1:] {
				_ = os.Remove(remaining.tempPath)
			}
			return nil, err
		}
		finalPaths = append(finalPaths, tf.finalPath)
	}
	return finalPaths, nil
}

// RemoveFile deletes a file from disk. Errors are logged but not returned
// (matches the existing best-effort behaviour in the original controller).
func (fs FileService) RemoveFile(path string) {
	if path == "" {
		return
	}
	if err := os.Remove(path); err != nil {
		log.Printf("FileService.RemoveFile: %s: %v", path, err)
	}
}
