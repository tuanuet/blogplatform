package repository

import (
	"context"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/google/uuid"
)

// BookmarkRepository defines the interface for bookmark data operations
type BookmarkRepository interface {
	// Add creates a bookmark for a user on a blog
	Add(ctx context.Context, userID, blogID uuid.UUID) error

	// Remove removes a bookmark for a user on a blog
	Remove(ctx context.Context, userID, blogID uuid.UUID) error

	// IsBookmarked checks if a blog is bookmarked by a user
	IsBookmarked(ctx context.Context, userID, blogID uuid.UUID) (bool, error)

	// CountByBlog returns the number of bookmarks for a blog
	CountByBlog(ctx context.Context, blogID uuid.UUID) (int64, error)

	// FindByUser returns a paginated list of blogs bookmarked by the user
	// Note: We return entity.Blog here, implying a join or separate fetch
	FindByUser(ctx context.Context, userID uuid.UUID, pagination Pagination) (*PaginatedResult[entity.Blog], error)
}
