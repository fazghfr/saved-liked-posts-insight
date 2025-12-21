package handlers

import (
	"net/http"

	"github.com/fazghfr/saved-liked-posts-insight/app_interface/internal/models"
	"github.com/fazghfr/saved-liked-posts-insight/app_interface/internal/services"
	"github.com/gin-gonic/gin"
)

// PostsHandler handles posts-related endpoints (proxy to app_core)
type PostsHandler struct {
	coreClient *services.CoreClient
}

// NewPostsHandler creates a new posts handler
func NewPostsHandler(coreClient *services.CoreClient) *PostsHandler {
	return &PostsHandler{
		coreClient: coreClient,
	}
}

// SamplePosts proxies the sample posts request to app_core
func (h *PostsHandler) SamplePosts(c *gin.Context) {
	var req models.SamplePostsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call app_core
	response, err := h.coreClient.SamplePosts(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CategorizePosts proxies the categorize request to app_core
func (h *PostsHandler) CategorizePosts(c *gin.Context) {
	var req models.CategorizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call app_core
	response, err := h.coreClient.CategorizePosts(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
