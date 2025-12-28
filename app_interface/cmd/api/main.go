package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app_interface/internal/config"
	"app_interface/internal/handlers"
	"app_interface/internal/middleware"
	"app_interface/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize services
	coreClient := services.NewCoreClient(cfg.CoreAPIURL)
	storageService := services.NewStorageService(cfg.UploadDir)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	uploadHandler := handlers.NewUploadHandler(storageService)
	postsHandler := handlers.NewPostsHandler(coreClient)

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Routes
	v1 := router.Group("/api/v1")
	{
		// Health check endpoints
		v1.GET("/health", healthHandler.Health)
		v1.GET("/ready", healthHandler.Ready)

		// Upload endpoints
		uploads := v1.Group("/uploads")
		{
			uploads.POST("/json", uploadHandler.UploadJSON)
			uploads.GET("/:id", uploadHandler.GetUpload)
			uploads.GET("", uploadHandler.ListUploads)
		}

		// Posts endpoints (proxy to app_core)
		posts := v1.Group("/posts")
		{
			posts.POST("/sample", postsHandler.SamplePosts)
			posts.POST("/categorize", postsHandler.CategorizePosts)
		}
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Instagram Data Processor - Go Interface",
			"version": "1.0.0",
			"docs":    "/api/v1",
		})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
