package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateTagRequest represents the request to create a tag
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
	Slug string `json:"slug" binding:"required,min=1,max=50"`
}

// UpdateTagRequest represents the request to update a tag
type UpdateTagRequest struct {
	Name *string `json:"name,omitempty" binding:"omitempty,min=1,max=50"`
	Slug *string `json:"slug,omitempty" binding:"omitempty,min=1,max=50"`
}

// TagResponse represents a tag in API responses
type TagResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
}
