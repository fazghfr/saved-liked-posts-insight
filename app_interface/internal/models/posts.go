package models

import "time"

// Upload represents a file upload record
type Upload struct {
	ID         string    `json:"id"`
	Filename   string    `json:"filename"`
	StoredAs   string    `json:"stored_as"`
	Size       int64     `json:"size"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// SamplePostsRequest matches the Python API request model
type SamplePostsRequest struct {
	Mode       string            `json:"mode" binding:"required"`
	SampleNum  int               `json:"sample_num" binding:"required,min=1,max=100"`
	Seed       int               `json:"seed"`
	CaptionMap map[int]string    `json:"caption_map,omitempty"`
}

// PostData matches the Python API response model
type PostData struct {
	Title     string `json:"title"`
	Link      string `json:"link"`
	Timestamp int    `json:"timestamp"`
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Captions  string `json:"captions"`
}

// SamplePostsResponse matches the Python API response model
type SamplePostsResponse struct {
	Posts []PostData `json:"posts"`
	Count int        `json:"count"`
}

// CategorizeRequest matches the Python API request model
type CategorizeRequest struct {
	Captions    []string `json:"captions" binding:"required"`
	Model       string   `json:"model"`
	SaveResults bool     `json:"save_results"`
}

// CategoryResult matches the Python API response model
type CategoryResult struct {
	RawReasoning     string   `json:"raw_reasoning"`
	SummaryReasoning string   `json:"summary_reasoning"`
	Categories       []string `json:"categories"`
}

// CategorizeResponse matches the Python API response model
type CategorizeResponse struct {
	Results []CategoryResult `json:"results"`
	Count   int              `json:"count"`
}
