package repository

import (
	"context"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type deviceTokenRepository struct {
	db *gorm.DB
}

func NewDeviceTokenRepository(db *gorm.DB) repository.DeviceTokenRepository {
	return &deviceTokenRepository{db: db}
}

func (r *deviceTokenRepository) Upsert(ctx context.Context, token *entity.UserDeviceToken) error {
	token.LastSeenAt = time.Now()

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "device_token"}},
			DoUpdates: clause.AssignmentColumns([]string{"last_seen_at"}),
		}).
		Create(token).Error
}

func (r *deviceTokenRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserDeviceToken, error) {
	var tokens []*entity.UserDeviceToken
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&tokens).Error
	return tokens, err
}

func (r *deviceTokenRepository) Delete(ctx context.Context, userID uuid.UUID, deviceToken string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND device_token = ?", userID, deviceToken).
		Delete(&entity.UserDeviceToken{}).Error
}

func (r *deviceTokenRepository) DeleteStaleTokens(ctx context.Context, days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)

	return r.db.WithContext(ctx).
		Where("last_seen_at < ?", cutoff).
		Delete(&entity.UserDeviceToken{}).Error
}
