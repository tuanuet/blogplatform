package repository

import (
	"context"
	"fmt"
	"time"

	domainRepo "github.com/aiagent/internal/domain/repository"
	"github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	client *redis.Client
}

// NewSessionRepository creates a new instance of SessionRepository
func NewSessionRepository(client *redis.Client) domainRepo.SessionRepository {
	return &sessionRepository{
		client: client,
	}
}

func (r *sessionRepository) CreateSession(ctx context.Context, sessionID string, userID string, duration time.Duration) error {
	key := r.getKey(sessionID)
	return r.client.Set(ctx, key, userID, duration).Err()
}

func (r *sessionRepository) GetUserID(ctx context.Context, sessionID string) (string, error) {
	key := r.getKey(sessionID)
	userID, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (r *sessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	key := r.getKey(sessionID)
	return r.client.Del(ctx, key).Err()
}

func (r *sessionRepository) getKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}
