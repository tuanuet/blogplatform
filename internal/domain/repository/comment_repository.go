package repository

import (
	"context"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/google/uuid"
)

// CommentRepository defines the interface for comment data operations
type CommentRepository interface {
	Create(ctx context.Context, comment *entity.Comment) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error)
	FindByBlogID(ctx context.Context, blogID uuid.UUID, pagination Pagination) (*PaginatedResult[entity.Comment], error)
	FindReplies(ctx context.Context, parentID uuid.UUID) ([]entity.Comment, error)
	Update(ctx context.Context, comment *entity.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error
}
