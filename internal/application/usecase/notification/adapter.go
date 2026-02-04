package notification

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/google/uuid"
)

// NotificationAdapter defines the interface for notification adapters
type NotificationAdapter interface {
	SendToUser(ctx context.Context, userID uuid.UUID, title, body string, data map[string]interface{}) error
	SendPush(ctx context.Context, deviceTokens []string, title, body string, data map[string]interface{}) error
}
