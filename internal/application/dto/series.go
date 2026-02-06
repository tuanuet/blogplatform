package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateSeriesRequest represents the request to create a series
type CreateSeriesRequest struct {
	Title       string `json:"title" binding:"required,max=255"`
	Slug        string `json:"slug" binding:"required,max=255"`
	Description string `json:"description"`
}

// UpdateSeriesRequest represents the request to update a series
type UpdateSeriesRequest struct {
	Title       string `json:"title" binding:"omitempty,max=255"`
	Description string `json:"description"`
}

// SeriesResponse represents a series in API responses
type SeriesResponse struct {
	ID          uuid.UUID          `json:"id"`
	AuthorID    uuid.UUID          `json:"authorId"`
	Title       string             `json:"title"`
	Slug        string             `json:"slug"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	Blogs       []BlogListResponse `json:"blogs,omitempty"`
	Author      *UserBriefResponse `json:"author,omitempty"`
}

// AddBlogToSeriesRequest represents the request to add a blog to a series
type AddBlogToSeriesRequest struct {
	BlogID uuid.UUID `json:"blogId" binding:"required"`
}

// SeriesFilterParams represents query parameters for filtering series
type SeriesFilterParams struct {
	AuthorID *string `form:"authorId"`
	Search   *string `form:"search"`
	Page     int     `form:"page,default=1"`
	PageSize int     `form:"pageSize,default=10"`
}

// HighlightedSeriesResponse represents a highlighted series in API responses
type HighlightedSeriesResponse struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Description     string    `json:"description"`
	AuthorID        uuid.UUID `json:"authorId"`
	AuthorName      string    `json:"authorName"`
	AuthorAvatarURL *string   `json:"authorAvatarUrl,omitempty"`
	SubscriberCount int       `json:"subscriberCount"`
	BlogCount       int       `json:"blogCount"`
	CreatedAt       time.Time `json:"createdAt"`
}
