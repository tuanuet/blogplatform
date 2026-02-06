package repository

import (
	"context"
	"time"

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

func (r *seriesRepository) GetHighlighted(ctx context.Context, limit int) ([]repository.HighlightedSeriesResult, error) {
	type highlightedSeriesRaw struct {
		ID              uuid.UUID `gorm:"column:id"`
		Title           string    `gorm:"column:title"`
		Slug            string    `gorm:"column:slug"`
		Description     string    `gorm:"column:description"`
		AuthorID        uuid.UUID `gorm:"column:author_id"`
		CreatedAt       time.Time `gorm:"column:created_at"`
		AuthorName      string    `gorm:"column:author_name"`
		AuthorAvatarURL *string   `gorm:"column:author_avatar_url"`
		SubscriberCount int       `gorm:"column:subscriber_count"`
		BlogCount       int       `gorm:"column:blog_count"`
	}

	var rawResults []highlightedSeriesRaw

	query := `SELECT
          s.id, s.title, s.slug, s.description, s.author_id, s.created_at,
          u.name as author_name, u.avatar_url as author_avatar_url,
          COALESCE(sub.subscriber_count, 0) as subscriber_count,
          COALESCE(bc.blog_count, 0) as blog_count
      FROM series s
      LEFT JOIN users u ON s.author_id = u.id
      LEFT JOIN (
          SELECT series_id, COUNT(*) as subscriber_count
          FROM user_series_purchases
          GROUP BY series_id
      ) sub ON s.id = sub.series_id
      LEFT JOIN (
          SELECT sb.series_id, COUNT(*) as blog_count
          FROM series_blogs sb
          JOIN blogs b ON sb.blog_id = b.id
          WHERE b.deleted_at IS NULL
          GROUP BY sb.series_id
      ) bc ON s.id = bc.series_id
      WHERE s.deleted_at IS NULL
      ORDER BY subscriber_count DESC
      LIMIT ?`

	if err := r.db.WithContext(ctx).Raw(query, limit).Scan(&rawResults).Error; err != nil {
		return nil, err
	}

	results := make([]repository.HighlightedSeriesResult, len(rawResults))
	for i, raw := range rawResults {
		results[i] = repository.HighlightedSeriesResult{
			Series: &entity.Series{
				ID:          raw.ID,
				Title:       raw.Title,
				Slug:        raw.Slug,
				Description: raw.Description,
				AuthorID:    raw.AuthorID,
				CreatedAt:   raw.CreatedAt,
			},
			Author: &entity.User{
				ID:        raw.AuthorID,
				Name:      raw.AuthorName,
				AvatarURL: raw.AuthorAvatarURL,
			},
			SubscriberCount: raw.SubscriberCount,
			BlogCount:       raw.BlogCount,
		}
	}

	return results, nil
}
