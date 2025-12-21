package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health returns the health status of the service
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "app_interface",
	})
}

// Ready returns the readiness status of the service
func (h *HealthHandler) Ready(c *gin.Context) {
	// TODO: Add checks for dependencies (database, app_core, etc.)
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"service": "app_interface",
	})
}
