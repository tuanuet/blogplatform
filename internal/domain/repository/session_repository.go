package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"time"
)

// SessionRepository defines the interface for session management
type SessionRepository interface {
	CreateSession(ctx context.Context, sessionID string, userID string, duration time.Duration) error
	GetUserID(ctx context.Context, sessionID string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
}
