package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// BlogVersionRepository defines the interface for blog version data operations
type BlogVersionRepository interface {
	// Create saves a new version
	Create(ctx context.Context, version *entity.BlogVersion) error

	// FindByID finds a version by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.BlogVersion, error)

	// FindByBlogID finds all versions for a blog with pagination
	FindByBlogID(ctx context.Context, blogID uuid.UUID, pagination Pagination) (*PaginatedResult[entity.BlogVersion], error)

	// GetNextVersionNumber gets the next version number for a blog
	GetNextVersionNumber(ctx context.Context, blogID uuid.UUID) (int, error)

	// Delete deletes a version by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// CountByBlogID counts versions for a blog
	CountByBlogID(ctx context.Context, blogID uuid.UUID) (int64, error)

	// DeleteOldest deletes oldest versions keeping only 'keep' most recent
	DeleteOldest(ctx context.Context, blogID uuid.UUID, keep int) error
}
