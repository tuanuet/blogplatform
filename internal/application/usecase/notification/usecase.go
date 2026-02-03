package notification

import (
	"context"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

// NotificationUseCase defines the notification use case interface
type NotificationUseCase interface {
	// List returns paginated notifications for a user
	List(ctx context.Context, userID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[dto.NotificationResponse], error)
	// GetUnreadCount returns the count of unread notifications for a user
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, userID uuid.UUID, notificationID uuid.UUID) error
	// MarkAllAsRead marks all notifications as read for a user
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	// GetPreferences returns the notification preferences for a user
	GetPreferences(ctx context.Context, userID uuid.UUID) ([]dto.NotificationPreferenceResponse, error)
	// UpdatePreferences updates the notification preferences for a user
	UpdatePreferences(ctx context.Context, userID uuid.UUID, req dto.UpdatePreferencesRequest) error
	// RegisterDeviceToken registers a device token for push notifications
	RegisterDeviceToken(ctx context.Context, userID uuid.UUID, req dto.RegisterDeviceTokenRequest) error
}
