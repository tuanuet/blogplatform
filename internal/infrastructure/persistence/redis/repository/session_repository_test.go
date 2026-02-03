package repository

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestSessionRepository_Integration(t *testing.T) {
	// This requires a running Redis instance.
	// We'll skip if connection fails.

	opts := &redis.Options{
		Addr: "localhost:6379",
	}
	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Skip("Redis is not available")
	}

	repo := NewSessionRepository(client)
	ctx := context.Background()
	sessionID := "test-session-id"
	userID := "test-user-id"
	duration := time.Minute

	// cleanup
	defer client.Del(ctx, "session:"+sessionID)

	// Test CreateSession
	err := repo.CreateSession(ctx, sessionID, userID, duration)
	assert.NoError(t, err)

	// Test GetUserID
	gotUserID, err := repo.GetUserID(ctx, sessionID)
	assert.NoError(t, err)
	assert.Equal(t, userID, gotUserID)

	// Test DeleteSession
	err = repo.DeleteSession(ctx, sessionID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.GetUserID(ctx, sessionID)
	assert.Error(t, err)
}
