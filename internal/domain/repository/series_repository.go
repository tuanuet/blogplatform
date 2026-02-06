package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

type HighlightedSeriesResult struct {
	Series          *entity.Series
	Author          *entity.User
	SubscriberCount int
	BlogCount       int
}

type SeriesRepository interface {
	Create(ctx context.Context, series *entity.Series) error
	Update(ctx context.Context, series *entity.Series) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Series, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Series, error)
	List(ctx context.Context, params map[string]interface{}) ([]entity.Series, int64, error)
	AddBlog(ctx context.Context, seriesID, blogID uuid.UUID) error
	RemoveBlog(ctx context.Context, seriesID, blogID uuid.UUID) error
	GetHighlighted(ctx context.Context, limit int) ([]HighlightedSeriesResult, error)
}
