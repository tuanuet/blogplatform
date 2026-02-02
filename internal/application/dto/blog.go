package dto

import (
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// CreateBlogRequest represents the request to create a blog
type CreateBlogRequest struct {
	Title        string     `json:"title" binding:"required,min=1,max=255"`
	Slug         string     `json:"slug" binding:"required,min=1,max=255"`
	Content      string     `json:"content" binding:"required"`
	Excerpt      *string    `json:"excerpt,omitempty"`
	ThumbnailURL *string    `json:"thumbnailUrl,omitempty" binding:"omitempty,url"`
	CategoryID   *string    `json:"categoryId,omitempty" binding:"omitempty,uuid"`
	TagIDs       []string   `json:"tagIds,omitempty"`
	PublishedAt  *time.Time `json:"publishedAt,omitempty"`
}

// UpdateBlogRequest represents the request to update a blog
type UpdateBlogRequest struct {
	Title        *string    `json:"title,omitempty" binding:"omitempty,min=1,max=255"`
	Slug         *string    `json:"slug,omitempty" binding:"omitempty,min=1,max=255"`
	Content      *string    `json:"content,omitempty"`
	Excerpt      *string    `json:"excerpt,omitempty"`
	ThumbnailURL *string    `json:"thumbnailUrl,omitempty" binding:"omitempty,url"`
	CategoryID   *string    `json:"categoryId,omitempty" binding:"omitempty,uuid"`
	TagIDs       []string   `json:"tagIds,omitempty"`
	PublishedAt  *time.Time `json:"publishedAt,omitempty"`
}

// PublishBlogRequest represents the request to publish a blog
type PublishBlogRequest struct {
	Visibility  string     `json:"visibility" binding:"required,oneof=public subscribers_only"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
}

// ReactionRequest represents the request to react to a blog
type ReactionRequest struct {
	Reaction string `json:"reaction" binding:"required,oneof=upvote downvote none"`
}

// ReactionResponse represents the response after reacting to a blog
type ReactionResponse struct {
	BlogID        uuid.UUID           `json:"blogId"`
	UpvoteCount   int                 `json:"upvoteCount"`
	DownvoteCount int                 `json:"downvoteCount"`
	UserReaction  entity.ReactionType `json:"userReaction"`
}

// BlogResponse represents a blog in API responses
type BlogResponse struct {
	ID            uuid.UUID             `json:"id"`
	AuthorID      uuid.UUID             `json:"authorId"`
	Author        *UserBriefResponse    `json:"author,omitempty"`
	CategoryID    *uuid.UUID            `json:"categoryId,omitempty"`
	Category      *CategoryResponse     `json:"category,omitempty"`
	Title         string                `json:"title"`
	Slug          string                `json:"slug"`
	Excerpt       *string               `json:"excerpt,omitempty"`
	Content       string                `json:"content"`
	ThumbnailURL  *string               `json:"thumbnailUrl,omitempty"`
	Status        entity.BlogStatus     `json:"status"`
	Visibility    entity.BlogVisibility `json:"visibility"`
	PublishedAt   *time.Time            `json:"publishedAt,omitempty"`
	Tags          []TagResponse         `json:"tags"`
	UpvoteCount   int                   `json:"upvoteCount"`
	DownvoteCount int                   `json:"downvoteCount"`
	UserReaction  *entity.ReactionType  `json:"userReaction,omitempty"` // For the current viewer
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
}

// BlogListResponse represents a blog in list view (without full content)
type BlogListResponse struct {
	ID            uuid.UUID             `json:"id"`
	AuthorID      uuid.UUID             `json:"authorId"`
	Author        *UserBriefResponse    `json:"author,omitempty"`
	CategoryID    *uuid.UUID            `json:"categoryId,omitempty"`
	Category      *CategoryResponse     `json:"category,omitempty"`
	Title         string                `json:"title"`
	Slug          string                `json:"slug"`
	Excerpt       *string               `json:"excerpt,omitempty"`
	ThumbnailURL  *string               `json:"thumbnailUrl,omitempty"`
	Status        entity.BlogStatus     `json:"status"`
	Visibility    entity.BlogVisibility `json:"visibility"`
	PublishedAt   *time.Time            `json:"publishedAt,omitempty"`
	Tags          []TagResponse         `json:"tags"`
	UpvoteCount   int                   `json:"upvoteCount"`
	DownvoteCount int                   `json:"downvoteCount"`
	CreatedAt     time.Time             `json:"createdAt"`
}

// BlogFilterParams represents query parameters for filtering blogs
type BlogFilterParams struct {
	AuthorID   *string  `form:"authorId"`
	CategoryID *string  `form:"categoryId"`
	Status     *string  `form:"status" binding:"omitempty,oneof=draft published"`
	Visibility *string  `form:"visibility" binding:"omitempty,oneof=public subscribers_only"`
	TagIDs     []string `form:"tagIds"`
	Search     *string  `form:"search"`
	Page       int      `form:"page,default=1"`
	PageSize   int      `form:"pageSize,default=10"`
}

// UserBriefResponse represents a brief user info for nested responses
type UserBriefResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
