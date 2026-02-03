package repository

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type notificationPreferenceRepository struct {
	db *gorm.DB
}

func NewNotificationPreferenceRepository(db *gorm.DB) repository.NotificationPreferenceRepository {
	return &notificationPreferenceRepository{db: db}
}

func (r *notificationPreferenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.NotificationPreference, error) {
	var prefs []*entity.NotificationPreference
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&prefs).Error
	return prefs, err
}

func (r *notificationPreferenceRepository) Upsert(ctx context.Context, preferences []*entity.NotificationPreference) error {
	if len(preferences) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "notification_type"}, {Name: "channel"}},
			DoUpdates: clause.AssignmentColumns([]string{"enabled"}),
		}).
		Create(&preferences).Error
}

func (r *notificationPreferenceRepository) IsEnabled(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, channel string) (bool, error) {
	var pref entity.NotificationPreference
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND notification_type = ? AND channel = ?", userID, notifType, channel).
		First(&pref).Error

	if err == gorm.ErrRecordNotFound {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return pref.Enabled, nil
}
