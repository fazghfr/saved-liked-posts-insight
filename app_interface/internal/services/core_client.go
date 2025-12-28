package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"app_interface/internal/models"
)

// CoreClient is an HTTP client for communicating with app_core
type CoreClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewCoreClient creates a new core API client
func NewCoreClient(baseURL string) *CoreClient {
	return &CoreClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SamplePosts calls the app_core sample posts endpoint
func (c *CoreClient) SamplePosts(ctx context.Context, req *models.SamplePostsRequest) (*models.SamplePostsResponse, error) {
	url := fmt.Sprintf("%s/posts/sample", c.baseURL)

	var response models.SamplePostsResponse
	if err := c.doRequest(ctx, "POST", url, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CategorizePosts calls the app_core categorize endpoint
func (c *CoreClient) CategorizePosts(ctx context.Context, req *models.CategorizeRequest) (*models.CategorizeResponse, error) {
	url := fmt.Sprintf("%s/posts/categorize", c.baseURL)

	var response models.CategorizeResponse
	if err := c.doRequest(ctx, "POST", url, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// doRequest performs an HTTP request and decodes the response
func (c *CoreClient) doRequest(ctx context.Context, method, url string, body, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("app_core returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
