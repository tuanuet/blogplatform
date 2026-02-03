package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// NotificationRepository defines the interface for notification data operations
type NotificationRepository interface {
	// Save creates or updates a notification with upsert logic for grouping
	Save(ctx context.Context, notification *entity.Notification) error

	// FindByUserID retrieves notifications for a user with pagination and total count
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Notification, int64, error)

	// FindUnreadCount returns the count of unread notifications for a user
	FindUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)

	// MarkAsRead marks a specific notification as read
	MarkAsRead(ctx context.Context, userID uuid.UUID, notificationID uuid.UUID) error

	// MarkAllAsRead marks all notifications as read for a user
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error

	// FindRecentUnread finds recent unread notification for aggregation
	FindRecentUnread(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, targetID uuid.UUID, windowMinutes int) (*entity.Notification, error)

	// DeleteExpired removes notifications that have passed their expiration time
	DeleteExpired(ctx context.Context) error
}

// NotificationPreferenceRepository defines the interface for notification preference operations
type NotificationPreferenceRepository interface {
	// GetByUserID retrieves all notification preferences for a user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.NotificationPreference, error)

	// Upsert creates or updates notification preferences
	Upsert(ctx context.Context, preferences []*entity.NotificationPreference) error

	// IsEnabled checks if a specific notification type and channel is enabled for a user
	IsEnabled(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, channel string) (bool, error)
}

// DeviceTokenRepository defines the interface for device token operations
type DeviceTokenRepository interface {
	// Upsert creates or updates a device token
	Upsert(ctx context.Context, token *entity.UserDeviceToken) error

	// FindByUserID retrieves all device tokens for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserDeviceToken, error)

	// Delete removes a specific device token for a user
	Delete(ctx context.Context, userID uuid.UUID, deviceToken string) error

	// DeleteStaleTokens removes device tokens that haven't been used in specified days
	DeleteStaleTokens(ctx context.Context, days int) error
}
