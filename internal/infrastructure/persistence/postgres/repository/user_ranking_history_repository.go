package repository

import (
	"context"
	"time"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRankingHistoryRepository struct {
	db *gorm.DB
}

// NewUserRankingHistoryRepository creates a new user ranking history repository
func NewUserRankingHistoryRepository(db *gorm.DB) repository.UserRankingHistoryRepository {
	return &userRankingHistoryRepository{db: db}
}

func (r *userRankingHistoryRepository) Create(ctx context.Context, history *entity.UserRankingHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *userRankingHistoryRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]entity.UserRankingHistory, error) {
	var history []entity.UserRankingHistory
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		Order("recorded_at DESC").
		Limit(limit).
		Find(&history).Error
	return history, err
}

func (r *userRankingHistoryRepository) FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserRankingHistory, error) {
	var history entity.UserRankingHistory
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		Order("recorded_at DESC").
		First(&history).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &history, err
}

func (r *userRankingHistoryRepository) ListByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]entity.UserRankingHistory, error) {
	var history []entity.UserRankingHistory
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ? AND recorded_at BETWEEN ? AND ?", userID, startDate, endDate).
		Order("recorded_at DESC").
		Find(&history).Error
	return history, err
}
