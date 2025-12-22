package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// StorageService handles file storage operations
type StorageService struct {
	uploadDir string
}

// NewStorageService creates a new storage service
func NewStorageService(uploadDir string) *StorageService {
	return &StorageService{
		uploadDir: uploadDir,
	}
}

// SaveFile saves an uploaded file to the storage directory
func (s *StorageService) SaveFile(c *gin.Context, file *multipart.FileHeader, filename string) error {
	dst := filepath.Join(s.uploadDir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// FileExists checks if a file exists in the storage directory
func (s *StorageService) FileExists(filename string) bool {
	path := filepath.Join(s.uploadDir, filename)
	_, err := os.Stat(path)
	return err == nil
}

// GetFilePath returns the full path to a stored file
func (s *StorageService) GetFilePath(filename string) string {
	return filepath.Join(s.uploadDir, filename)
}

// DeleteFile removes a file from storage
func (s *StorageService) DeleteFile(filename string) error {
	path := filepath.Join(s.uploadDir, filename)
	return os.Remove(path)
}

// FileInfo represents information about a stored file
type FileInfo struct {
	Name    string
	Size    int64
	ModTime string
}

// ListFiles returns a list of all files in the storage directory
func (s *StorageService) ListFiles() ([]FileInfo, error) {
	entries, err := os.ReadDir(s.uploadDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []FileInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue // Skip files we can't read
			}
			files = append(files, FileInfo{
				Name:    entry.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
			})
		}
	}

	return files, nil
}
