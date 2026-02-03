package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// BlogFilter defines filter options for blog queries
type BlogFilter struct {
	AuthorID        *uuid.UUID
	CategoryID      *uuid.UUID
	Status          *entity.BlogStatus
	Visibility      *entity.BlogVisibility
	TagIDs          []uuid.UUID
	Search          *string // search in title or content
	PublishedBefore *time.Time
}

// BlogRepository defines the interface for blog data operations
type BlogRepository interface {
	Create(ctx context.Context, blog *entity.Blog) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Blog, error)
	FindBySlug(ctx context.Context, authorID uuid.UUID, slug string) (*entity.Blog, error)
	FindAll(ctx context.Context, filter BlogFilter, pagination Pagination) (*PaginatedResult[entity.Blog], error)
	Update(ctx context.Context, blog *entity.Blog) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Tag operations
	AddTags(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error
	RemoveTags(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error
	ReplaceTags(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error

	// Reaction operations
	// UpdateCounts atomically updates the reaction counts for a blog
	UpdateCounts(ctx context.Context, blogID uuid.UUID, upDelta, downDelta int) error
	// React handles the reaction logic (insert/delete/swap) and returns the DELTA (change) in counts
	// rather than the absolute totals, to allow for buffered updates.
	React(ctx context.Context, blogID, userID uuid.UUID, reactionType entity.ReactionType) (upDelta, downDelta int, err error)

	// Recommendation operations
	FindRelated(ctx context.Context, blogID uuid.UUID, limit int) ([]entity.Blog, error)

	// Admin stats
	CountByMonth(ctx context.Context, months int) ([]entity.MonthlyCount, error)

	// ExistsByAuthorAndTag checks if any blog by the author has the specified tag
	ExistsByAuthorAndTag(ctx context.Context, authorID, tagID uuid.UUID) (bool, error)
}
