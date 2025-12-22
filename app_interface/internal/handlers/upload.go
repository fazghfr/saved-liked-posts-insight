package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/fazghfr/saved-liked-posts-insight/app_interface/internal/models"
	"github.com/fazghfr/saved-liked-posts-insight/app_interface/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadHandler handles file upload endpoints
type UploadHandler struct {
	storage *services.StorageService
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(storage *services.StorageService) *UploadHandler {
	return &UploadHandler{
		storage: storage,
	}
}

// UploadJSON handles JSON file uploads
func (h *UploadHandler) UploadJSON(c *gin.Context) {
	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JSON files are allowed"})
		return
	}

	// Generate unique ID for the upload
	uploadID := uuid.New().String()
	filename := fmt.Sprintf("%s%s", uploadID, ext)

	// Save file
	if err := h.storage.SaveFile(c, file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file: %v", err)})
		return
	}

	// Create upload record
	upload := models.Upload{
		ID:         uploadID,
		Filename:   file.Filename,
		StoredAs:   filename,
		Size:       file.Size,
		UploadedAt: time.Now(),
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "File uploaded successfully",
		"upload":  upload,
	})
}

// GetUpload retrieves upload information by ID
func (h *UploadHandler) GetUpload(c *gin.Context) {
	uploadID := c.Param("id")

	// TODO: Implement database lookup for upload metadata
	// For now, just check if file exists
	filename := fmt.Sprintf("%s.json", uploadID)
	exists := h.storage.FileExists(filename)

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       uploadID,
		"filename": filename,
		"exists":   exists,
	})
}

// ListUploads lists all uploads
func (h *UploadHandler) ListUploads(c *gin.Context) {
	// Get list of files from storage
	files, err := h.storage.ListFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to list uploads: %v", err),
		})
		return
	}

	// Convert file info to upload models
	uploads := make([]models.Upload, 0, len(files))
	for _, file := range files {
		// Extract upload ID from filename (remove .json extension)
		uploadID := filepath.Base(file.Name)
		if ext := filepath.Ext(uploadID); ext != "" {
			uploadID = uploadID[:len(uploadID)-len(ext)]
		}

		// Parse modification time
		uploadedAt, err := time.Parse("2006-01-02 15:04:05", file.ModTime)
		if err != nil {
			// If parsing fails, use current time as fallback
			uploadedAt = time.Now()
		}

		upload := models.Upload{
			ID:         uploadID,
			Filename:   file.Name, // Original filename is not stored, using stored filename
			StoredAs:   file.Name,
			Size:       file.Size,
			UploadedAt: uploadedAt,
		}
		uploads = append(uploads, upload)
	}

	c.JSON(http.StatusOK, gin.H{
		"uploads": uploads,
		"count":   len(uploads),
	})
}
