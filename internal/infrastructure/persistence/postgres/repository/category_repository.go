package repository

import (
	"context"
	"math"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	var category entity.Category
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&category).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &category, err
}

func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.WithContext(ctx).
		Where("slug = ? AND deleted_at IS NULL", slug).
		First(&category).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &category, err
}

func (r *categoryRepository) FindAll(ctx context.Context, pagination repository.Pagination) (*repository.PaginatedResult[entity.Category], error) {
	var categories []entity.Category
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Category{}).Where("deleted_at IS NULL")

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.
		Order("name ASC").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &repository.PaginatedResult[entity.Category]{
		Data:       categories,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.Category{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}
