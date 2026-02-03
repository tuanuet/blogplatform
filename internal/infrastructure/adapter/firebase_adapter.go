package adapter

import (
	"context"
	"errors"

	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// TargetType represents the type of target entity for notification
type TargetType string

const (
	TargetTypeBlog    TargetType = "blog"
	TargetTypeComment TargetType = "comment"
	TargetTypeUser    TargetType = "user"
)

// NotificationCategory represents the category of notification
type NotificationCategory string

const (
	NotificationCategorySocial  NotificationCategory = "social"
	NotificationCategoryContent NotificationCategory = "content"
	NotificationCategorySystem  NotificationCategory = "system"
)

// validTargetTypes defines the valid target types for notifications
var validTargetTypes = map[TargetType]bool{
	TargetTypeBlog:    true,
	TargetTypeComment: true,
	TargetTypeUser:    true,
}

// validCategories defines the valid notification categories
var validCategories = map[NotificationCategory]bool{
	NotificationCategorySocial:  true,
	NotificationCategoryContent: true,
	NotificationCategorySystem:  true,
}

// FirebaseClient defines the interface for Firebase FCM operations
type FirebaseClient interface {
	Send(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error
}

// firebaseAdapter implements FirebaseAdapter for push notifications
type firebaseAdapter struct {
	deviceRepo repository.DeviceTokenRepository
	firebase   FirebaseClient
}

// NewFirebaseAdapter creates a new Firebase adapter instance
func NewFirebaseAdapter(deviceRepo repository.DeviceTokenRepository, firebase FirebaseClient) *firebaseAdapter {
	return &firebaseAdapter{
		deviceRepo: deviceRepo,
		firebase:   firebase,
	}
}

// SendPush sends a push notification to specific device tokens
func (a *firebaseAdapter) SendPush(ctx context.Context, deviceTokens []string, title, body string, data map[string]interface{}) error {
	// Validate title
	if title == "" {
		return errors.New("title cannot be empty")
	}

	// Return early if no tokens to send to
	if len(deviceTokens) == 0 {
		return nil
	}

	// Validate data payload
	if err := a.validateDataPayload(data); err != nil {
		return err
	}

	// Send push notification
	err := a.firebase.Send(ctx, deviceTokens, title, body, data)
	if err != nil {
		log.Error().Err(err).
			Str("title", title).
			Strs("device_tokens", deviceTokens).
			Msg("failed to send push notification via Firebase")
		return err
	}

	log.Info().
		Str("title", title).
		Int("token_count", len(deviceTokens)).
		Msg("push notification sent successfully")

	return nil
}

// SendToUser sends a push notification to all devices of a user
func (a *firebaseAdapter) SendToUser(ctx context.Context, userID uuid.UUID, title, body string, data map[string]interface{}) error {
	// Validate title
	if title == "" {
		return errors.New("title cannot be empty")
	}

	// Fetch device tokens for the user
	tokens, err := a.deviceRepo.FindByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Str("user_id", userID.String()).
			Msg("failed to get device tokens for user")
		return errors.New("failed to get device tokens")
	}

	// Return early if user has no registered devices
	if len(tokens) == 0 {
		log.Info().
			Str("user_id", userID.String()).
			Msg("no device tokens found for user, skipping push notification")
		return nil
	}

	// Extract device token strings
	deviceTokens := make([]string, 0, len(tokens))
	for _, token := range tokens {
		deviceTokens = append(deviceTokens, token.DeviceToken)
	}

	// Send push notification
	return a.SendPush(ctx, deviceTokens, title, body, data)
}

// validateDataPayload validates the data payload structure
func (a *firebaseAdapter) validateDataPayload(data map[string]interface{}) error {
	// Check if data is nil
	if data == nil {
		return errors.New("data payload is required")
	}

	// Validate target_type
	targetTypeRaw, ok := data["target_type"]
	if !ok {
		return errors.New("target_type is required")
	}

	targetTypeStr, ok := targetTypeRaw.(string)
	if !ok {
		return errors.New("target_type must be a string")
	}

	targetType := TargetType(targetTypeStr)
	if !validTargetTypes[targetType] {
		return errors.New("invalid target_type: must be 'blog', 'comment', or 'user'")
	}

	// Validate target_id
	targetIDRaw, ok := data["target_id"]
	if !ok {
		return errors.New("target_id is required")
	}

	_, ok = targetIDRaw.(string)
	if !ok {
		return errors.New("target_id must be a string")
	}

	// Validate category
	categoryRaw, ok := data["category"]
	if !ok {
		return errors.New("category is required")
	}

	categoryStr, ok := categoryRaw.(string)
	if !ok {
		return errors.New("category must be a string")
	}

	category := NotificationCategory(categoryStr)
	if !validCategories[category] {
		return errors.New("invalid category: must be 'social', 'content', or 'system'")
	}

	// Validate deep_link (optional but recommended)
	if deepLinkRaw, ok := data["deep_link"]; ok {
		_, ok = deepLinkRaw.(string)
		if !ok {
			return errors.New("deep_link must be a string")
		}
	}

	return nil
}

// NewMockFirebaseAdapter creates a mock Firebase adapter for testing
// This function is needed because the entity package defines UserDeviceToken
// and tests need to use entity.UserDeviceToken instead of repository.UserDeviceToken
func NewMockFirebaseAdapter(deviceRepo repository.DeviceTokenRepository, firebase FirebaseClient) *firebaseAdapter {
	return &firebaseAdapter{
		deviceRepo: deviceRepo,
		firebase:   firebase,
	}
}
