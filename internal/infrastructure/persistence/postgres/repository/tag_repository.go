package repository

import (
	"context"
	"math"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *gorm.DB) repository.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(ctx context.Context, tag *entity.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *tagRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tag).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &tag, err
}

func (r *tagRepository) FindBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &tag, err
}

func (r *tagRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

func (r *tagRepository) FindAll(ctx context.Context, pagination repository.Pagination) (*repository.PaginatedResult[entity.Tag], error) {
	var tags []entity.Tag
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Tag{})

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.
		Order("name ASC").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &repository.PaginatedResult[entity.Tag]{
		Data:       tags,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *tagRepository) FindOrCreate(ctx context.Context, name, slug string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	if err == nil {
		return &tag, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create new tag
	tag = entity.Tag{Name: name, Slug: slug}
	if err := r.db.WithContext(ctx).Create(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) Update(ctx context.Context, tag *entity.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Tag{}, "id = ?", id).Error
}
