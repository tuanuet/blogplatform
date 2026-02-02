package repository

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type seriesRepository struct {
	db *gorm.DB
}

// NewSeriesRepository creates a new series repository
func NewSeriesRepository(db *gorm.DB) repository.SeriesRepository {
	return &seriesRepository{db: db}
}

func (r *seriesRepository) Create(ctx context.Context, series *entity.Series) error {
	return r.db.WithContext(ctx).Create(series).Error
}

func (r *seriesRepository) Update(ctx context.Context, series *entity.Series) error {
	return r.db.WithContext(ctx).Save(series).Error
}

func (r *seriesRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Series{}, id).Error
}

func (r *seriesRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Series, error) {
	var series entity.Series
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Blogs").
		Preload("Blogs.Author").
		Preload("Blogs.Category").
		Preload("Blogs.Tags").
		First(&series, id).Error
	if err != nil {
		return nil, err
	}
	return &series, nil
}

func (r *seriesRepository) GetBySlug(ctx context.Context, slug string) (*entity.Series, error) {
	var series entity.Series
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Blogs").
		Preload("Blogs.Author").
		Preload("Blogs.Category").
		Preload("Blogs.Tags").
		Where("slug = ?", slug).
		First(&series).Error
	if err != nil {
		return nil, err
	}
	return &series, nil
}

func (r *seriesRepository) List(ctx context.Context, params map[string]interface{}) ([]entity.Series, int64, error) {
	var series []entity.Series
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Series{})

	if authorID, ok := params["author_id"].(string); ok && authorID != "" {
		query = query.Where("author_id = ?", authorID)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit, ok := params["limit"].(int); ok && limit > 0 {
		query = query.Limit(limit)
	}

	if offset, ok := params["offset"].(int); ok && offset >= 0 {
		query = query.Offset(offset)
	}

	err := query.Preload("Author").
		Preload("Blogs").
		Order("created_at DESC").
		Find(&series).Error

	if err != nil {
		return nil, 0, err
	}

	return series, total, nil
}

func (r *seriesRepository) AddBlog(ctx context.Context, seriesID, blogID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.Series{ID: seriesID}).Association("Blogs").Append(&entity.Blog{ID: blogID})
}

func (r *seriesRepository) RemoveBlog(ctx context.Context, seriesID, blogID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.Series{ID: seriesID}).Association("Blogs").Delete(&entity.Blog{ID: blogID})
}
