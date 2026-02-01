package repository

import (
	"context"
	"math"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *gorm.DB) repository.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *entity.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&comment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &comment, err
}

func (r *commentRepository) FindByBlogID(ctx context.Context, blogID uuid.UUID, pagination repository.Pagination) (*repository.PaginatedResult[entity.Comment], error) {
	var comments []entity.Comment
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.Comment{}).
		Where("blog_id = ? AND parent_id IS NULL AND deleted_at IS NULL", blogID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (pagination.Page - 1) * pagination.PageSize
	err := query.
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").Preload("User").Order("created_at ASC")
		}).
		Order("created_at DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&comments).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &repository.PaginatedResult[entity.Comment]{
		Data:       comments,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *commentRepository) FindReplies(ctx context.Context, parentID uuid.UUID) ([]entity.Comment, error) {
	var replies []entity.Comment
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("parent_id = ? AND deleted_at IS NULL", parentID).
		Order("created_at ASC").
		Find(&replies).Error
	return replies, err
}

func (r *commentRepository) Update(ctx context.Context, comment *entity.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *commentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.Comment{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}
