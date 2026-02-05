package repository

import (
	"context"
	"math"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type blogVersionRepository struct {
	db *gorm.DB
}

// NewBlogVersionRepository creates a new blog version repository
func NewBlogVersionRepository(db *gorm.DB) repository.BlogVersionRepository {
	return &blogVersionRepository{db: db}
}

func (r *blogVersionRepository) Create(ctx context.Context, version *entity.BlogVersion) error {
	return r.db.WithContext(ctx).Create(version).Error
}

func (r *blogVersionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.BlogVersion, error) {
	var version entity.BlogVersion
	err := r.db.WithContext(ctx).
		Preload("Editor").
		Preload("Category").
		Preload("Tags").
		Where("id = ?", id).
		First(&version).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *blogVersionRepository) FindByBlogID(ctx context.Context, blogID uuid.UUID, pagination repository.Pagination) (*repository.PaginatedResult[entity.BlogVersion], error) {
	var versions []entity.BlogVersion
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.BlogVersion{}).Where("blog_id = ?", blogID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.
		Preload("Editor").
		Order("version_number DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&versions).Error

	if err != nil {
		return nil, err
	}

	totalPages := 0
	if pagination.PageSize > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pagination.PageSize)))
	}

	return &repository.PaginatedResult[entity.BlogVersion]{
		Data:       versions,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *blogVersionRepository) GetNextVersionNumber(ctx context.Context, blogID uuid.UUID) (int, error) {
	var nextVersion int
	err := r.db.WithContext(ctx).
		Model(&entity.BlogVersion{}).
		Select("COALESCE(MAX(version_number), 0) + 1").
		Where("blog_id = ?", blogID).
		Scan(&nextVersion).Error
	return nextVersion, err
}

func (r *blogVersionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.BlogVersion{}, id).Error
}

func (r *blogVersionRepository) CountByBlogID(ctx context.Context, blogID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.BlogVersion{}).
		Where("blog_id = ?", blogID).
		Count(&count).Error
	return count, err
}

func (r *blogVersionRepository) DeleteOldest(ctx context.Context, blogID uuid.UUID, keep int) error {
	count, err := r.CountByBlogID(ctx, blogID)
	if err != nil {
		return err
	}

	if count <= int64(keep) {
		return nil
	}

	deleteCount := int(count) - keep

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Exec(
			`DELETE FROM "blog_versions" WHERE id IN (SELECT id FROM "blog_versions" WHERE blog_id = ? ORDER BY version_number ASC LIMIT ?)`,
			blogID, deleteCount,
		).Error
	})
}
