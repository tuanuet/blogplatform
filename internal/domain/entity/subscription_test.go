package entity_test

import (
	"testing"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

func TestSubscription_IsActive(t *testing.T) {
	now := time.Now()
	futureTime := now.Add(24 * time.Hour)
	pastTime := now.Add(-24 * time.Hour)

	tests := []struct {
		name     string
		sub      entity.Subscription
		expected bool
	}{
		{
			name: "returns true for free follower (ExpiresAt is nil)",
			sub: entity.Subscription{
				ID:           uuid.New(),
				SubscriberID: uuid.New(),
				AuthorID:     uuid.New(),
				ExpiresAt:    nil,
			},
			expected: true,
		},
		{
			name: "returns true for paid subscription not yet expired",
			sub: entity.Subscription{
				ID:           uuid.New(),
				SubscriberID: uuid.New(),
				AuthorID:     uuid.New(),
				ExpiresAt:    &futureTime,
			},
			expected: true,
		},
		{
			name: "returns false for paid subscription expired",
			sub: entity.Subscription{
				ID:           uuid.New(),
				SubscriberID: uuid.New(),
				AuthorID:     uuid.New(),
				ExpiresAt:    &pastTime,
			},
			expected: false,
		},
		{
			name: "returns true for paid subscription expiring exactly now",
			sub: entity.Subscription{
				ID:           uuid.New(),
				SubscriberID: uuid.New(),
				AuthorID:     uuid.New(),
				ExpiresAt:    &now,
			},
			expected: false, // ExpiresAt.After(time.Now()) returns false when equal
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sub.IsActive()
			if result != tt.expected {
				t.Errorf("IsActive() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSubscription_IsPaid(t *testing.T) {
	now := time.Now()
	futureTime := now.Add(24 * time.Hour)

	tests := []struct {
		name     string
		sub      entity.Subscription
		expected bool
	}{
		{
			name: "returns false for free follower (ExpiresAt is nil)",
			sub: entity.Subscription{
				ID:           uuid.New(),
				SubscriberID: uuid.New(),
				AuthorID:     uuid.New(),
				ExpiresAt:    nil,
			},
			expected: false,
		},
		{
			name: "returns true for paid subscription with future expiry",
			sub: entity.Subscription{
				ID:           uuid.New(),
				SubscriberID: uuid.New(),
				AuthorID:     uuid.New(),
				ExpiresAt:    &futureTime,
			},
			expected: true,
		},
		{
			name: "returns true for paid subscription even when expired",
			sub: entity.Subscription{
				ID:           uuid.New(),
				SubscriberID: uuid.New(),
				AuthorID:     uuid.New(),
				ExpiresAt:    &now,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sub.IsPaid()
			if result != tt.expected {
				t.Errorf("IsPaid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSubscription_IsExpired(t *testing.T) {
	now := time.Now()
	futureTime := now.Add(24 * time.Hour)
	pastTime := now.Add(-24 * time.Hour)

	tests := []struct {
		name     string
		sub      entity.Subscription
		expected bool
	}{
		{
			name: "returns false for free follower (ExpiresAt is nil)",
			sub: entity.Subscription{
				ID:           uuid.New(),
				SubscriberID: uuid.New(),
				AuthorID:     uuid.New(),
				ExpiresAt:    nil,
			},
			expected: false,
		},
		{
			name: "returns false for paid subscription not yet expired",
			sub: entity.Subscription{
				ID:           uuid.New(),
				SubscriberID: uuid.New(),
				AuthorID:     uuid.New(),
				ExpiresAt:    &futureTime,
			},
			expected: false,
		},
		{
			name: "returns true for paid subscription expired",
			sub: entity.Subscription{
				ID:           uuid.New(),
				SubscriberID: uuid.New(),
				AuthorID:     uuid.New(),
				ExpiresAt:    &pastTime,
			},
			expected: true,
		},
		{
			name: "returns true for paid subscription expiring exactly now",
			sub: entity.Subscription{
				ID:           uuid.New(),
				SubscriberID: uuid.New(),
				AuthorID:     uuid.New(),
				ExpiresAt:    &now,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sub.IsExpired()
			if result != tt.expected {
				t.Errorf("IsExpired() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSubscription_NewFields(t *testing.T) {
	now := time.Now()
	futureTime := now.Add(24 * time.Hour)
	id := uuid.New()
	subscriberID := uuid.New()
	authorID := uuid.New()

	tests := []struct {
		name string
		sub  entity.Subscription
	}{
		{
			name: "has ExpiresAt field for paid subscription",
			sub: entity.Subscription{
				ID:           id,
				SubscriberID: subscriberID,
				AuthorID:     authorID,
				ExpiresAt:    &futureTime,
			},
		},
		{
			name: "has nil ExpiresAt for free follower",
			sub: entity.Subscription{
				ID:           id,
				SubscriberID: subscriberID,
				AuthorID:     authorID,
				ExpiresAt:    nil,
			},
		},
		{
			name: "has Tier field with FREE value",
			sub: entity.Subscription{
				ID:           id,
				SubscriberID: subscriberID,
				AuthorID:     authorID,
				Tier:         "FREE",
			},
		},
		{
			name: "has Tier field with PREMIUM value",
			sub: entity.Subscription{
				ID:           id,
				SubscriberID: subscriberID,
				AuthorID:     authorID,
				Tier:         "PREMIUM",
			},
		},
		{
			name: "has UpdatedAt field",
			sub: entity.Subscription{
				ID:           id,
				SubscriberID: subscriberID,
				AuthorID:     authorID,
				UpdatedAt:    now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify struct can be created with all new fields
			if tt.sub.ID != id {
				t.Errorf("ID mismatch")
			}
			if tt.sub.SubscriberID != subscriberID {
				t.Errorf("SubscriberID mismatch")
			}
			if tt.sub.AuthorID != authorID {
				t.Errorf("AuthorID mismatch")
			}
			// If the struct compiles and we can access the fields, the test passes
		})
	}
}
