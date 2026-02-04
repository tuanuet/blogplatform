package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// TagTierMappingRepository defines the interface for tag-tier mapping data operations
type TagTierMappingRepository interface {
	// Create creates a new tag-tier mapping
	Create(ctx context.Context, mapping *entity.TagTierMapping) error

	// FindByID retrieves a mapping by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.TagTierMapping, error)

	// FindByAuthorAndTag retrieves a mapping by author ID and tag ID
	FindByAuthorAndTag(ctx context.Context, authorID, tagID uuid.UUID) (*entity.TagTierMapping, error)

	// FindByAuthor retrieves all mappings for an author
	FindByAuthor(ctx context.Context, authorID uuid.UUID) ([]entity.TagTierMapping, error)

	// FindByTagIDs retrieves mappings for specific tags and author
	FindByTagIDs(ctx context.Context, authorID uuid.UUID, tagIDs []uuid.UUID) ([]entity.TagTierMapping, error)

	// Update updates an existing mapping
	Update(ctx context.Context, mapping *entity.TagTierMapping) error

	// Upsert creates or updates a mapping based on (author_id, tag_id) uniqueness
	Upsert(ctx context.Context, mapping *entity.TagTierMapping) error

	// Delete deletes a tag-tier mapping
	Delete(ctx context.Context, authorID, tagID uuid.UUID) error

	// DeleteByID deletes a mapping by ID
	DeleteByID(ctx context.Context, id uuid.UUID) error

	// CountBlogsByTagAndAuthor counts blogs that have the specified tag from the author
	CountBlogsByTagAndAuthor(ctx context.Context, authorID, tagID uuid.UUID) (int64, error)

	// WithTx returns a new repository with the given transaction
	WithTx(tx interface{}) TagTierMappingRepository
}
