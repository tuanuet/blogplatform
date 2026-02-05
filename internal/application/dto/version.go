package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateVersionRequest represents the request to create a new version
type CreateVersionRequest struct {
	ChangeSummary string `json:"changeSummary,omitempty"`
}

// VersionResponse represents a version in list view
type VersionResponse struct {
	ID            uuid.UUID         `json:"id"`
	VersionNumber int               `json:"versionNumber"`
	Title         string            `json:"title"`
	Excerpt       *string           `json:"excerpt,omitempty"`
	Status        string            `json:"status"`
	Visibility    string            `json:"visibility"`
	Editor        UserBriefResponse `json:"editor"`
	ChangeSummary *string           `json:"changeSummary,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
}

// VersionDetailResponse represents a version with full content
type VersionDetailResponse struct {
	ID            uuid.UUID         `json:"id"`
	VersionNumber int               `json:"versionNumber"`
	Title         string            `json:"title"`
	Slug          string            `json:"slug"`
	Excerpt       *string           `json:"excerpt,omitempty"`
	Content       string            `json:"content"`
	ThumbnailURL  *string           `json:"thumbnailUrl,omitempty"`
	Status        string            `json:"status"`
	Visibility    string            `json:"visibility"`
	Category      *CategoryResponse `json:"category,omitempty"`
	Tags          []TagResponse     `json:"tags,omitempty"`
	Editor        UserBriefResponse `json:"editor"`
	ChangeSummary *string           `json:"changeSummary,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
}

// VersionListResponse represents a paginated list of versions
type VersionListResponse struct {
	Data       []VersionResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalPages int               `json:"totalPages"`
}
