package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateCategoryRequest represents the request to create a category
type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Slug        string  `json:"slug" binding:"required,min=1,max=100"`
	Description *string `json:"description,omitempty"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Slug        *string `json:"slug,omitempty" binding:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
}

// CategoryResponse represents a category in API responses
type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
