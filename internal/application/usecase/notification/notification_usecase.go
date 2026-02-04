package notification

import (
	"context"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

type notificationUseCase struct {
	notificationRepo repository.NotificationRepository
	deviceRepo       repository.DeviceTokenRepository
	prefRepo         repository.NotificationPreferenceRepository
	adapter          NotificationAdapter
}

func NewNotificationUseCase(
	notificationRepo repository.NotificationRepository,
	deviceRepo repository.DeviceTokenRepository,
	prefRepo repository.NotificationPreferenceRepository,
	adapter NotificationAdapter,
) NotificationUseCase {
	return &notificationUseCase{
		notificationRepo: notificationRepo,
		deviceRepo:       deviceRepo,
		prefRepo:         prefRepo,
		adapter:          adapter,
	}
}

func (u *notificationUseCase) List(ctx context.Context, userID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[dto.NotificationResponse], error) {
	offset := (page - 1) * pageSize
	notifications, total, err := u.notificationRepo.FindByUserID(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	data := make([]dto.NotificationResponse, len(notifications))
	for i, n := range notifications {
		data[i] = dto.NotificationResponse{
			ID:        n.ID,
			Title:     n.Title,
			Body:      n.Body,
			Type:      string(n.Type),
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt,
			Data:      n.Data,
		}
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &repository.PaginatedResult[dto.NotificationResponse]{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (u *notificationUseCase) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	return u.notificationRepo.FindUnreadCount(ctx, userID)
}

func (u *notificationUseCase) MarkAsRead(ctx context.Context, userID uuid.UUID, notificationID uuid.UUID) error {
	return u.notificationRepo.MarkAsRead(ctx, userID, notificationID)
}

func (u *notificationUseCase) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return u.notificationRepo.MarkAllAsRead(ctx, userID)
}

func (u *notificationUseCase) GetPreferences(ctx context.Context, userID uuid.UUID) ([]dto.NotificationPreferenceResponse, error) {
	prefs, err := u.prefRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.NotificationPreferenceResponse, 0, len(prefs))
	for _, p := range prefs {
		result = append(result, dto.NotificationPreferenceResponse{
			NotificationType: dto.NotificationType(p.NotificationType),
			Channel:          p.Channel,
			Enabled:          p.Enabled,
		})
	}
	return result, nil
}

func (u *notificationUseCase) UpdatePreferences(ctx context.Context, userID uuid.UUID, req dto.UpdatePreferencesRequest) error {
	prefs := make([]*entity.NotificationPreference, len(req.Preferences))
	for i, p := range req.Preferences {
		prefs[i] = &entity.NotificationPreference{
			UserID:           userID,
			NotificationType: entity.NotificationType(p.NotificationType),
			Channel:          p.Channel,
			Enabled:          p.Enabled,
		}
	}
	return u.prefRepo.Upsert(ctx, prefs)
}

func (u *notificationUseCase) RegisterDeviceToken(ctx context.Context, userID uuid.UUID, req dto.RegisterDeviceTokenRequest) error {
	token := &entity.UserDeviceToken{
		UserID:      userID,
		DeviceToken: req.DeviceToken,
		Platform:    req.Platform,
		LastSeenAt:  time.Now(),
	}
	return u.deviceRepo.Upsert(ctx, token)
}
