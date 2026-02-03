package adapter

import (
	"context"
	"testing"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// Mock firebase client interface for testing
type mockFirebaseClient struct {
	sendFn func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error
}

func (m *mockFirebaseClient) Send(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
	if m.sendFn != nil {
		return m.sendFn(ctx, tokens, title, body, data)
	}
	return nil
}

func TestFirebaseAdapter_SendPush_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockClient := &mockFirebaseClient{
		sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
			return nil
		},
	}

	adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
	ctx := context.Background()
	deviceTokens := []string{"token1", "token2", "token3"}
	title := "Test Notification"
	body := "This is a test notification"
	data := map[string]interface{}{
		"target_type": "blog",
		"target_id":   uuid.New().String(),
		"category":    "content",
		"deep_link":   "/blogs/123",
	}

	// Act
	err := adapter.SendPush(ctx, deviceTokens, title, body, data)

	// Assert
	assert.NoError(t, err)
}

func TestFirebaseAdapter_SendToUser_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	deviceTokens := []*entity.UserDeviceToken{
		{DeviceToken: "token1"},
		{DeviceToken: "token2"},
	}

	mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockClient := &mockFirebaseClient{
		sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
			return nil
		},
	}

	mockDeviceRepo.EXPECT().
		FindByUserID(gomock.Any(), userID).
		Return(deviceTokens, nil).
		Times(1)

	adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
	ctx := context.Background()
	title := "New Blog Post"
	body := "Check out the new blog post"
	data := map[string]interface{}{
		"target_type": "blog",
		"target_id":   uuid.New().String(),
		"category":    "content",
		"deep_link":   "/blogs/123",
	}

	// Act
	err := adapter.SendToUser(ctx, userID, title, body, data)

	// Assert
	assert.NoError(t, err)
}

func TestFirebaseAdapter_SendToUser_NoTokens(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()

	mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockClient := &mockFirebaseClient{
		sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
			return nil
		},
	}

	mockDeviceRepo.EXPECT().
		FindByUserID(gomock.Any(), userID).
		Return([]*entity.UserDeviceToken{}, nil).
		Times(1)

	adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
	ctx := context.Background()
	title := "New Blog Post"
	body := "Check out the new blog post"
	data := map[string]interface{}{
		"target_type": "blog",
		"target_id":   uuid.New().String(),
		"category":    "content",
		"deep_link":   "/blogs/123",
	}

	// Act
	err := adapter.SendToUser(ctx, userID, title, body, data)

	// Assert
	assert.NoError(t, err)
}

func TestFirebaseAdapter_SendToUser_TokensFromMultipleDevices(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	deviceTokens := []*entity.UserDeviceToken{
		{DeviceToken: "ios_token_1"},
		{DeviceToken: "ios_token_2"},
		{DeviceToken: "android_token_1"},
		{DeviceToken: "android_token_2"},
	}

	mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockClient := &mockFirebaseClient{
		sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
			return nil
		},
	}

	mockDeviceRepo.EXPECT().
		FindByUserID(gomock.Any(), userID).
		Return(deviceTokens, nil).
		Times(1)

	adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
	ctx := context.Background()
	title := "New Blog Post"
	body := "Check out the new blog post"
	data := map[string]interface{}{
		"target_type": "blog",
		"target_id":   uuid.New().String(),
		"category":    "content",
		"deep_link":   "/blogs/123",
	}

	// Act
	err := adapter.SendToUser(ctx, userID, title, body, data)

	// Assert
	assert.NoError(t, err)
}

func TestFirebaseAdapter_SendToUser_RepositoryError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()

	mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockClient := &mockFirebaseClient{
		sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
			return nil
		},
	}

	mockDeviceRepo.EXPECT().
		FindByUserID(gomock.Any(), userID).
		Return(nil, assert.AnError).
		Times(1)

	adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
	ctx := context.Background()
	title := "New Blog Post"
	body := "Check out the new blog post"
	data := map[string]interface{}{
		"target_type": "blog",
		"target_id":   uuid.New().String(),
		"category":    "content",
		"deep_link":   "/blogs/123",
	}

	// Act
	err := adapter.SendToUser(ctx, userID, title, body, data)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get device tokens")
}

func TestFirebaseAdapter_SendPush_InvalidTargetType(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockClient := &mockFirebaseClient{
		sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
			return nil
		},
	}

	adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
	ctx := context.Background()
	deviceTokens := []string{"token1"}
	title := "Test Notification"
	body := "This is a test notification"
	data := map[string]interface{}{
		"target_type": "invalid_type",
		"target_id":   uuid.New().String(),
		"category":    "content",
		"deep_link":   "/blogs/123",
	}

	// Act
	err := adapter.SendPush(ctx, deviceTokens, title, body, data)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid target_type")
}

func TestFirebaseAdapter_SendPush_InvalidCategory(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockClient := &mockFirebaseClient{
		sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
			return nil
		},
	}

	adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
	ctx := context.Background()
	deviceTokens := []string{"token1"}
	title := "Test Notification"
	body := "This is a test notification"
	data := map[string]interface{}{
		"target_type": "blog",
		"target_id":   uuid.New().String(),
		"category":    "invalid_category",
		"deep_link":   "/blogs/123",
	}

	// Act
	err := adapter.SendPush(ctx, deviceTokens, title, body, data)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid category")
}

func TestFirebaseAdapter_SendPush_EmptyTitle(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockClient := &mockFirebaseClient{
		sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
			return nil
		},
	}

	adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
	ctx := context.Background()
	deviceTokens := []string{"token1"}
	title := ""
	body := "This is a test notification"
	data := map[string]interface{}{
		"target_type": "blog",
		"target_id":   uuid.New().String(),
		"category":    "content",
		"deep_link":   "/blogs/123",
	}

	// Act
	err := adapter.SendPush(ctx, deviceTokens, title, body, data)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title cannot be empty")
}

func TestFirebaseAdapter_SendPush_EmptyDeviceTokens(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockClient := &mockFirebaseClient{
		sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
			return nil
		},
	}

	adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
	ctx := context.Background()
	deviceTokens := []string{}
	title := "Test Notification"
	body := "This is a test notification"
	data := map[string]interface{}{
		"target_type": "blog",
		"target_id":   uuid.New().String(),
		"category":    "content",
		"deep_link":   "/blogs/123",
	}

	// Act
	err := adapter.SendPush(ctx, deviceTokens, title, body, data)

	// Assert
	// Should return early without error if no tokens
	assert.NoError(t, err)
}

func TestFirebaseAdapter_SendPush_ValidTargetTypes(t *testing.T) {
	// Arrange
	validTargetTypes := []string{"blog", "comment", "user"}

	for _, targetType := range validTargetTypes {
		t.Run(targetType, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
			mockClient := &mockFirebaseClient{
				sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
					return nil
				},
			}

			adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
			ctx := context.Background()
			deviceTokens := []string{"token1"}
			title := "Test Notification"
			body := "This is a test notification"
			data := map[string]interface{}{
				"target_type": targetType,
				"target_id":   uuid.New().String(),
				"category":    "content",
				"deep_link":   "/blogs/123",
			}

			// Act
			err := adapter.SendPush(ctx, deviceTokens, title, body, data)

			// Assert
			assert.NoError(t, err)
		})
	}
}

func TestFirebaseAdapter_SendPush_ValidCategories(t *testing.T) {
	// Arrange
	validCategories := []string{"social", "content", "system"}

	for _, category := range validCategories {
		t.Run(category, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
			mockClient := &mockFirebaseClient{
				sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
					return nil
				},
			}

			adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
			ctx := context.Background()
			deviceTokens := []string{"token1"}
			title := "Test Notification"
			body := "This is a test notification"
			data := map[string]interface{}{
				"target_type": "blog",
				"target_id":   uuid.New().String(),
				"category":    category,
				"deep_link":   "/blogs/123",
			}

			// Act
			err := adapter.SendPush(ctx, deviceTokens, title, body, data)

			// Assert
			assert.NoError(t, err)
		})
	}
}

func TestFirebaseAdapter_SendPush_DataPayloadValidation(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockClient := &mockFirebaseClient{
		sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
			return nil
		},
	}

	adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
	ctx := context.Background()
	deviceTokens := []string{"token1"}
	title := "Test Notification"
	body := "This is a test notification"

	// Test with complete data payload
	data := map[string]interface{}{
		"target_type": "blog",
		"target_id":   uuid.New().String(),
		"category":    "content",
		"deep_link":   "/blogs/123",
	}

	// Act
	err := adapter.SendPush(ctx, deviceTokens, title, body, data)

	// Assert
	assert.NoError(t, err)
}

func TestFirebaseAdapter_SendPush_MissingRequiredDataFields(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeviceRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockClient := &mockFirebaseClient{
		sendFn: func(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error {
			return nil
		},
	}

	adapter := NewFirebaseAdapter(mockDeviceRepo, mockClient)
	ctx := context.Background()
	deviceTokens := []string{"token1"}
	title := "Test Notification"
	body := "This is a test notification"

	// Test with missing target_type
	data := map[string]interface{}{
		"target_id": uuid.New().String(),
		"category":  "content",
		"deep_link": "/blogs/123",
	}

	// Act
	err := adapter.SendPush(ctx, deviceTokens, title, body, data)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "target_type is required")
}
