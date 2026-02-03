package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aiagent/internal/domain/entity"
	repoMocks "github.com/aiagent/internal/domain/repository/mocks"
	"github.com/aiagent/internal/domain/service"
	cacheMocks "github.com/aiagent/internal/infrastructure/cache/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNotificationAggregator_ShouldAggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockNotificationRepository(ctrl)
	mockCache := cacheMocks.NewMockCache(ctrl)

	t.Run("should return existing notification when similar notification exists within window", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		notifType := entity.NotificationTypeBlogLike
		targetID := uuid.New()

		existingNotif := &entity.Notification{
			ID:           uuid.New(),
			UserID:       userID,
			Type:         notifType,
			GroupedCount: 3,
			Data:         map[string]interface{}{"target_id": targetID.String()},
		}

		mockRepo.EXPECT().
			FindRecentUnread(ctx, userID, notifType, targetID, 5).
			Return(existingNotif, nil)

		// Act
		aggregator := service.NewNotificationAggregator(mockRepo, mockCache)
		result, err := aggregator.ShouldAggregate(ctx, userID, notifType, targetID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, existingNotif.ID, result.ID)
		assert.Equal(t, 3, result.GroupedCount)
	})

	t.Run("should return nil when no similar notification exists within window", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		notifType := entity.NotificationTypeBlogLike
		targetID := uuid.New()

		mockRepo.EXPECT().
			FindRecentUnread(ctx, userID, notifType, targetID, 5).
			Return(nil, nil)

		// Act
		aggregator := service.NewNotificationAggregator(mockRepo, mockCache)
		result, err := aggregator.ShouldAggregate(ctx, userID, notifType, targetID)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		notifType := entity.NotificationTypeBlogLike
		targetID := uuid.New()
		repoErr := errors.New("database error")

		mockRepo.EXPECT().
			FindRecentUnread(ctx, userID, notifType, targetID, 5).
			Return(nil, repoErr)

		// Act
		aggregator := service.NewNotificationAggregator(mockRepo, mockCache)
		result, err := aggregator.ShouldAggregate(ctx, userID, notifType, targetID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to find recent unread")
	})
}

func TestNotificationAggregator_CheckRateLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockNotificationRepository(ctrl)
	mockCache := cacheMocks.NewMockCache(ctrl)

	t.Run("should allow when under limit", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		notifType := entity.NotificationTypeBlogLike
		cacheKey := "rate_limit:" + userID.String() + ":" + string(notifType)

		mockCache.EXPECT().
			Get(ctx, cacheKey, gomock.Any()).
			Return(errors.New("not found"))

		// Act
		aggregator := service.NewNotificationAggregator(mockRepo, mockCache)
		allowed, err := aggregator.CheckRateLimit(ctx, userID, notifType)

		// Assert
		assert.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("should allow when under limit with existing count", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		notifType := entity.NotificationTypeBlogLike
		cacheKey := "rate_limit:" + userID.String() + ":" + string(notifType)

		// Mock Get - return 50 (under 100 limit)
		gomock.InOrder(
			mockCache.EXPECT().
				Get(ctx, cacheKey, gomock.Any()).
				Do(func(ctx context.Context, key string, dest interface{}) {
					if ptr, ok := dest.(*int); ok {
						*ptr = 50
					}
				}),
		)

		// Act
		aggregator := service.NewNotificationAggregator(mockRepo, mockCache)
		allowed, err := aggregator.CheckRateLimit(ctx, userID, notifType)

		// Assert
		assert.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("should deny when limit exceeded", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		notifType := entity.NotificationTypeBlogLike
		cacheKey := "rate_limit:" + userID.String() + ":" + string(notifType)

		// Mock Get - return 100 (at limit)
		gomock.InOrder(
			mockCache.EXPECT().
				Get(ctx, cacheKey, gomock.Any()).
				Do(func(ctx context.Context, key string, dest interface{}) {
					if ptr, ok := dest.(*int); ok {
						*ptr = 100
					}
				}),
		)

		// Act
		aggregator := service.NewNotificationAggregator(mockRepo, mockCache)
		allowed, err := aggregator.CheckRateLimit(ctx, userID, notifType)

		// Assert
		assert.NoError(t, err)
		assert.False(t, allowed)
	})

	t.Run("should return error when cache fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		notifType := entity.NotificationTypeBlogLike
		cacheKey := "rate_limit:" + userID.String() + ":" + string(notifType)
		cacheErr := errors.New("cache error")

		mockCache.EXPECT().
			Get(ctx, cacheKey, gomock.Any()).
			Return(cacheErr)

		// Act
		aggregator := service.NewNotificationAggregator(mockRepo, mockCache)
		allowed, err := aggregator.CheckRateLimit(ctx, userID, notifType)

		// Assert
		assert.Error(t, err)
		assert.False(t, allowed)
		assert.Contains(t, err.Error(), "failed to check rate limit")
	})
}

func TestNotificationAggregator_IncrementRateLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockNotificationRepository(ctrl)
	mockCache := cacheMocks.NewMockCache(ctrl)

	t.Run("should increment rate limit counter", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		notifType := entity.NotificationTypeBlogLike
		cacheKey := "rate_limit:" + userID.String() + ":" + string(notifType)
		ttl := time.Hour

		// Get will return not found (new window)
		mockCache.EXPECT().Get(ctx, cacheKey, gomock.Any()).Return(errors.New("not found"))
		// Then Set will be called
		mockCache.EXPECT().Set(ctx, cacheKey, 1, ttl).Return(nil)

		// Act
		aggregator := service.NewNotificationAggregator(mockRepo, mockCache)
		err := aggregator.IncrementRateLimit(ctx, userID, notifType)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should handle cache set errors", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		notifType := entity.NotificationTypeBlogLike
		cacheKey := "rate_limit:" + userID.String() + ":" + string(notifType)
		ttl := time.Hour
		cacheErr := errors.New("cache error")

		// Get will return not found (new window)
		mockCache.EXPECT().Get(ctx, cacheKey, gomock.Any()).Return(errors.New("not found"))
		// Set will fail
		mockCache.EXPECT().Set(ctx, cacheKey, 1, ttl).Return(cacheErr)

		// Act
		aggregator := service.NewNotificationAggregator(mockRepo, mockCache)
		err := aggregator.IncrementRateLimit(ctx, userID, notifType)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to increment rate limit")
	})
}

// Helper type for mocking the cache interface
// We need a type that can be marshaled/unmarshaled by JSON
type rateLimitCounter struct {
	Count int       `json:"count"`
	Time  time.Time `json:"time"`
}
