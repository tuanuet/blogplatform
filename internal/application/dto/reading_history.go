package dto

import (
	"time"

	"github.com/google/uuid"
)

// RecordViewRequest represents the request to record a blog view
// Note: BlogID is usually taken from the URL path parameter
type RecordViewRequest struct {
	// Empty body is acceptable if only the action of reading is recorded
}

// ReadingHistoryResponse represents a single history item in the response
type ReadingHistoryResponse struct {
	BlogID     uuid.UUID         `json:"blogId"`
	LastReadAt time.Time         `json:"lastReadAt"`
	Blog       *BlogListResponse `json:"blog,omitempty"`
}

// ReadingHistoryListResponse represents the list of reading history
type ReadingHistoryListResponse struct {
	History []ReadingHistoryResponse `json:"history"`
}
