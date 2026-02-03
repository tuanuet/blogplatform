package repository

import (
	"context"
	"errors"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) repository.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Save(ctx context.Context, notification *entity.Notification) error {
	recent, err := r.FindRecentUnread(ctx, notification.UserID, notification.Type, extractTargetID(notification.Data), 5)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if recent != nil {
		return r.updateGroupedNotification(ctx, recent, notification)
	}

	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *notificationRepository) updateGroupedNotification(ctx context.Context, existing *entity.Notification, new *entity.Notification) error {
	existing.GroupedCount++
	existing.UpdatedAt = time.Now()

	if existing.Data == nil {
		existing.Data = make(map[string]interface{})
	}

	actorIDs, _ := existing.Data["actor_ids"].([]uuid.UUID)
	actorIDs = append(actorIDs, extractActorID(new.Data))
	existing.Data["actor_ids"] = actorIDs

	actorName := extractActorName(new.Data)
	if actorName != "" {
		existing.Body = generateGroupedBody(actorName, existing.GroupedCount, existing.Type)
	}

	return r.db.WithContext(ctx).Model(&entity.Notification{}).
		Where("id = ?", existing.ID).
		Updates(map[string]interface{}{
			"grouped_count": existing.GroupedCount,
			"body":          existing.Body,
			"data":          existing.Data,
			"updated_at":    existing.UpdatedAt,
		}).Error
}

func (r *notificationRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Notification, int64, error) {
	var notifications []*entity.Notification
	var total int64

	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *notificationRepository) FindUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = ? AND expires_at > ?", userID, false, time.Now()).
		Count(&count).Error
	return int(count), err
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, userID uuid.UUID, notificationID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

func (r *notificationRepository) FindRecentUnread(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, targetID uuid.UUID, windowMinutes int) (*entity.Notification, error) {
	var notification entity.Notification
	window := time.Now().Add(-time.Duration(windowMinutes) * time.Minute)

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND is_read = ? AND created_at > ?", userID, notifType, false, window).
		Where("(data->>'target_id')::text = ?", targetID.String()).
		Order("created_at DESC").
		First(&notification).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entity.Notification{}).Error
}

func extractTargetID(data map[string]interface{}) uuid.UUID {
	if data == nil {
		return uuid.Nil
	}
	if id, ok := data["target_id"].(string); ok {
		if parsed, err := uuid.Parse(id); err == nil {
			return parsed
		}
	}
	return uuid.Nil
}

func extractActorID(data map[string]interface{}) uuid.UUID {
	if data == nil {
		return uuid.Nil
	}
	if id, ok := data["actor_id"].(string); ok {
		if parsed, err := uuid.Parse(id); err == nil {
			return parsed
		}
	}
	return uuid.Nil
}

func extractActorName(data map[string]interface{}) string {
	if data == nil {
		return ""
	}
	if name, ok := data["actor_name"].(string); ok {
		return name
	}
	return ""
}

func generateGroupedBody(actorName string, count int, notifType entity.NotificationType) string {
	if count <= 1 {
		return ""
	}

	otherCount := count - 1
	if otherCount == 1 {
		switch notifType {
		case entity.NotificationTypeBlogLike:
			return actorName + " and 1 other liked your blog"
		default:
			return actorName + " and 1 other"
		}
	}
	return actorName + " and " + string(rune(otherCount)) + " others"
}
