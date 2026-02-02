package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// TagRepository defines the interface for tag data operations
type TagRepository interface {
	Create(ctx context.Context, tag *entity.Tag) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Tag, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Tag, error)
	FindAll(ctx context.Context, pagination Pagination) (*PaginatedResult[entity.Tag], error)
	FindOrCreate(ctx context.Context, name, slug string) (*entity.Tag, error)
	Update(ctx context.Context, tag *entity.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindPopular(ctx context.Context, limit int) ([]entity.Tag, error)
}
