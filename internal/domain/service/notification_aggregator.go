package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"fmt"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/infrastructure/cache"
	"github.com/google/uuid"
)

const (
	// RateLimitMaxNotifications is the maximum number of notifications allowed per hour
	RateLimitMaxNotifications = 100
	// RateLimitWindow is the time window for rate limiting (1 hour)
	RateLimitWindow = time.Hour
	// AggregationWindowMinutes is the time window for notification aggregation (5 minutes)
	AggregationWindowMinutes = 5
)

// NotificationAggregator handles smart notification aggregation and rate limiting
type NotificationAggregator interface {
	// ShouldAggregate checks if a notification should be aggregated with an existing one
	// Returns the existing notification if aggregation should happen, nil otherwise
	ShouldAggregate(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, targetID uuid.UUID) (*entity.Notification, error)

	// CheckRateLimit checks if a user is within their notification rate limit
	// Returns true if allowed, false if limit exceeded
	CheckRateLimit(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType) (bool, error)

	// IncrementRateLimit increments the rate limit counter for a user and notification type
	IncrementRateLimit(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType) error
}

// notificationAggregator implements NotificationAggregator interface
type notificationAggregator struct {
	repo  repository.NotificationRepository
	cache cache.Cache
}

// NewNotificationAggregator creates a new NotificationAggregator instance
func NewNotificationAggregator(repo repository.NotificationRepository, cache cache.Cache) NotificationAggregator {
	return &notificationAggregator{
		repo:  repo,
		cache: cache,
	}
}

// ShouldAggregate checks if a notification should be aggregated with an existing one
func (s *notificationAggregator) ShouldAggregate(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, targetID uuid.UUID) (*entity.Notification, error) {
	recent, err := s.repo.FindRecentUnread(ctx, userID, notifType, targetID, AggregationWindowMinutes)
	if err != nil {
		return nil, fmt.Errorf("failed to find recent unread notification: %w", err)
	}
	return recent, nil
}

// CheckRateLimit checks if a user is within their notification rate limit
func (s *notificationAggregator) CheckRateLimit(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType) (bool, error) {
	cacheKey := s.buildRateLimitKey(userID, notifType)

	var count int
	err := s.cache.Get(ctx, cacheKey, &count)
	if err != nil {
		// Check if it's a "not found" error (new window) vs actual cache error
		if isNotFoundError(err) {
			// New window - allow
			return true, nil
		}
		// Actual cache error - return error
		return false, fmt.Errorf("failed to check rate limit: %w", err)
	}

	return count < RateLimitMaxNotifications, nil
}

// isNotFoundError checks if the error is a cache "not found" error
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	// redis.Nil indicates key not found
	return err.Error() == "redis: nil" || err.Error() == "not found"
}

// IncrementRateLimit increments the rate limit counter for a user and notification type
func (s *notificationAggregator) IncrementRateLimit(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType) error {
	cacheKey := s.buildRateLimitKey(userID, notifType)

	var count int
	_ = s.cache.Get(ctx, cacheKey, &count)
	count++

	if err := s.cache.Set(ctx, cacheKey, count, RateLimitWindow); err != nil {
		return fmt.Errorf("failed to increment rate limit: %w", err)
	}

	return nil
}

// buildRateLimitKey builds a Redis key for rate limiting
func (s *notificationAggregator) buildRateLimitKey(userID uuid.UUID, notifType entity.NotificationType) string {
	return fmt.Sprintf("rate_limit:%s:%s", userID.String(), string(notifType))
}

// FirebaseAdapter defines the interface for Firebase Cloud Messaging (FCM) operations
type FirebaseAdapter interface {
	// SendPush sends a push notification to the specified device tokens
	SendPush(ctx context.Context, deviceTokens []string, title, body string, data map[string]interface{}) error

	// SendPushToUser sends a push notification to all device tokens for a user
	SendPushToUser(ctx context.Context, userID uuid.UUID, title, body string, data map[string]interface{}) error
}
