package repository

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type readingHistoryRepository struct {
	db *gorm.DB
}

func NewReadingHistoryRepository(db *gorm.DB) repository.ReadingHistoryRepository {
	return &readingHistoryRepository{db: db}
}

func (r *readingHistoryRepository) Upsert(ctx context.Context, history *entity.UserReadingHistory) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "blog_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_read_at"}),
	}).Create(history).Error
}

func (r *readingHistoryRepository) GetRecentByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.UserReadingHistory, error) {
	var histories []*entity.UserReadingHistory
	err := r.db.WithContext(ctx).
		Preload("Blog").
		Preload("Blog.Author").
		Preload("Blog.Category").
		Preload("Blog.Tags").
		Where("user_id = ?", userID).
		Order("last_read_at DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}
