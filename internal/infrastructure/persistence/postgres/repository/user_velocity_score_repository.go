package repository

import (
	"context"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userVelocityScoreRepository struct {
	db *gorm.DB
}

// NewUserVelocityScoreRepository creates a new user velocity score repository
func NewUserVelocityScoreRepository(db *gorm.DB) repository.UserVelocityScoreRepository {
	return &userVelocityScoreRepository{db: db}
}

func (r *userVelocityScoreRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserVelocityScore, error) {
	var score entity.UserVelocityScore
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		First(&score).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &score, err
}

func (r *userVelocityScoreRepository) FindByRank(ctx context.Context, rank int) (*entity.UserVelocityScore, error) {
	var score entity.UserVelocityScore
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("rank_position = ?", rank).
		First(&score).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &score, err
}

func (r *userVelocityScoreRepository) ListTopRanked(ctx context.Context, limit int) ([]entity.UserVelocityScore, error) {
	var scores []entity.UserVelocityScore
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("rank_position IS NOT NULL").
		Order("rank_position ASC").
		Limit(limit).
		Find(&scores).Error
	return scores, err
}

func (r *userVelocityScoreRepository) ListRanked(ctx context.Context, pagination repository.Pagination) (*repository.PaginatedResult[entity.UserVelocityScore], error) {
	var scores []entity.UserVelocityScore
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&entity.UserVelocityScore{}).
		Where("rank_position IS NOT NULL").
		Count(&total).Error; err != nil {
		return nil, err
	}

	// Fetch paginated results
	offset := (pagination.Page - 1) * pagination.PageSize
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("rank_position IS NOT NULL").
		Order("rank_position ASC").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&scores).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pagination.PageSize
	if int(total)%pagination.PageSize > 0 {
		totalPages++
	}

	return &repository.PaginatedResult[entity.UserVelocityScore]{
		Data:       scores,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *userVelocityScoreRepository) Save(ctx context.Context, score *entity.UserVelocityScore) error {
	return r.db.WithContext(ctx).Save(score).Error
}

func (r *userVelocityScoreRepository) UpdateRankPosition(ctx context.Context, userID uuid.UUID, rank int) error {
	return r.db.WithContext(ctx).
		Model(&entity.UserVelocityScore{}).
		Where("user_id = ?", userID).
		Update("rank_position", rank).Error
}

func (r *userVelocityScoreRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&entity.UserVelocityScore{}).Error
}

func (r *userVelocityScoreRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.UserVelocityScore{}).
		Where("rank_position IS NOT NULL").
		Count(&count).Error
	return count, err
}
