package notification

import (
	"context"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/notification/mocks"
	"github.com/aiagent/internal/domain/entity"
	repoMocks "github.com/aiagent/internal/domain/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNotificationUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockNotificationRepository(ctrl)
	mockDeviceRepo := repoMocks.NewMockDeviceTokenRepository(ctrl)
	mockPrefRepo := repoMocks.NewMockNotificationPreferenceRepository(ctrl)
	mockAdapter := mocks.NewMockNotificationAdapter(ctrl)

	uc := NewNotificationUseCase(mockRepo, mockDeviceRepo, mockPrefRepo, mockAdapter)

	ctx := context.Background()
	userID := uuid.New()
	page := 1
	pageSize := 10

	t.Run("success", func(t *testing.T) {
		notifications := []*entity.Notification{
			{
				ID:        uuid.New(),
				UserID:    userID,
				Title:     "Test",
				Body:      "Body",
				CreatedAt: time.Now(),
			},
		}
		total := int64(1)

		mockRepo.EXPECT().
			FindByUserID(ctx, userID, pageSize, 0).
			Return(notifications, total, nil)

		result, err := uc.List(ctx, userID, page, pageSize)

		assert.NoError(t, err)
		assert.Equal(t, total, result.Total)
		assert.Len(t, result.Data, 1)
		assert.Equal(t, notifications[0].ID, result.Data[0].ID)
	})
}

func TestNotificationUseCase_GetUnreadCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockNotificationRepository(ctrl)
	mockDeviceRepo := repoMocks.NewMockDeviceTokenRepository(ctrl)
	mockPrefRepo := repoMocks.NewMockNotificationPreferenceRepository(ctrl)
	mockAdapter := mocks.NewMockNotificationAdapter(ctrl)

	uc := NewNotificationUseCase(mockRepo, mockDeviceRepo, mockPrefRepo, mockAdapter)

	ctx := context.Background()
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		count := 5
		mockRepo.EXPECT().
			FindUnreadCount(ctx, userID).
			Return(count, nil)

		result, err := uc.GetUnreadCount(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, count, result)
	})
}

func TestNotificationUseCase_MarkAsRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockNotificationRepository(ctrl)
	mockDeviceRepo := repoMocks.NewMockDeviceTokenRepository(ctrl)
	mockPrefRepo := repoMocks.NewMockNotificationPreferenceRepository(ctrl)
	mockAdapter := mocks.NewMockNotificationAdapter(ctrl)

	uc := NewNotificationUseCase(mockRepo, mockDeviceRepo, mockPrefRepo, mockAdapter)

	ctx := context.Background()
	userID := uuid.New()
	notifID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			MarkAsRead(ctx, userID, notifID).
			Return(nil)

		err := uc.MarkAsRead(ctx, userID, notifID)

		assert.NoError(t, err)
	})
}

func TestNotificationUseCase_MarkAllAsRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockNotificationRepository(ctrl)
	mockDeviceRepo := repoMocks.NewMockDeviceTokenRepository(ctrl)
	mockPrefRepo := repoMocks.NewMockNotificationPreferenceRepository(ctrl)
	mockAdapter := mocks.NewMockNotificationAdapter(ctrl)

	uc := NewNotificationUseCase(mockRepo, mockDeviceRepo, mockPrefRepo, mockAdapter)

	ctx := context.Background()
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			MarkAllAsRead(ctx, userID).
			Return(nil)

		err := uc.MarkAllAsRead(ctx, userID)

		assert.NoError(t, err)
	})
}

func TestNotificationUseCase_RegisterDeviceToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockNotificationRepository(ctrl)
	mockDeviceRepo := repoMocks.NewMockDeviceTokenRepository(ctrl)
	mockPrefRepo := repoMocks.NewMockNotificationPreferenceRepository(ctrl)
	mockAdapter := mocks.NewMockNotificationAdapter(ctrl)

	uc := NewNotificationUseCase(mockRepo, mockDeviceRepo, mockPrefRepo, mockAdapter)

	ctx := context.Background()
	userID := uuid.New()
	req := dto.RegisterDeviceTokenRequest{
		DeviceToken: "token123",
		Platform:    "ios",
	}

	t.Run("success", func(t *testing.T) {
		mockDeviceRepo.EXPECT().
			Upsert(ctx, gomock.Any()).
			Return(nil)

		err := uc.RegisterDeviceToken(ctx, userID, req)

		assert.NoError(t, err)
	})
}

func TestNotificationUseCase_GetPreferences(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockNotificationRepository(ctrl)
	mockDeviceRepo := repoMocks.NewMockDeviceTokenRepository(ctrl)
	mockPrefRepo := repoMocks.NewMockNotificationPreferenceRepository(ctrl)
	mockAdapter := mocks.NewMockNotificationAdapter(ctrl)

	uc := NewNotificationUseCase(mockRepo, mockDeviceRepo, mockPrefRepo, mockAdapter)

	ctx := context.Background()
	userID := uuid.New()

	t.Run("success - returns preferences", func(t *testing.T) {
		prefs := []*entity.NotificationPreference{
			{UserID: userID, NotificationType: "social", Channel: "push", Enabled: true},
		}
		mockPrefRepo.EXPECT().
			GetByUserID(ctx, userID).
			Return(prefs, nil)

		result, err := uc.GetPreferences(ctx, userID)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "social", string(result[0].NotificationType))
	})
}

func TestNotificationUseCase_UpdatePreferences(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockNotificationRepository(ctrl)
	mockDeviceRepo := repoMocks.NewMockDeviceTokenRepository(ctrl)
	mockPrefRepo := repoMocks.NewMockNotificationPreferenceRepository(ctrl)
	mockAdapter := mocks.NewMockNotificationAdapter(ctrl)

	uc := NewNotificationUseCase(mockRepo, mockDeviceRepo, mockPrefRepo, mockAdapter)

	ctx := context.Background()
	userID := uuid.New()
	req := dto.UpdatePreferencesRequest{
		Preferences: []dto.NotificationPreferenceItem{
			{NotificationType: "social", Channel: "push", Enabled: false},
		},
	}

	t.Run("success", func(t *testing.T) {
		mockPrefRepo.EXPECT().
			Upsert(ctx, gomock.Any()).
			Return(nil)

		err := uc.UpdatePreferences(ctx, userID, req)

		assert.NoError(t, err)
	})
}

// NewNotificationUseCase definition is needed for the test to compile,
// even if we haven't implemented it yet. We'll define a dummy one if it doesn't exist
// or just rely on the implementation step.
// Since we are doing TDD, I expect the test to fail compilation first because NewNotificationUseCase doesn't exist.
