package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateCommentRequest represents the request to create a comment
type CreateCommentRequest struct {
	Content  string  `json:"content" binding:"required,min=1"`
	ParentID *string `json:"parentId,omitempty" binding:"omitempty,uuid"`
}

// UpdateCommentRequest represents the request to update a comment
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

// CommentResponse represents a comment in API responses
type CommentResponse struct {
	ID        uuid.UUID          `json:"id"`
	BlogID    uuid.UUID          `json:"blogId"`
	UserID    uuid.UUID          `json:"userId"`
	User      *UserBriefResponse `json:"user,omitempty"`
	ParentID  *uuid.UUID         `json:"parentId,omitempty"`
	Content   string             `json:"content"`
	Replies   []CommentResponse  `json:"replies,omitempty"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}
