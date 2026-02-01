package dto

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionResponse represents a subscription in API responses
type SubscriptionResponse struct {
	ID           uuid.UUID          `json:"id"`
	SubscriberID uuid.UUID          `json:"subscriberId"`
	Subscriber   *UserBriefResponse `json:"subscriber,omitempty"`
	AuthorID     uuid.UUID          `json:"authorId"`
	Author       *UserBriefResponse `json:"author,omitempty"`
	CreatedAt    time.Time          `json:"createdAt"`
}

// SubscriptionCountResponse represents subscriber count
type SubscriptionCountResponse struct {
	AuthorID        uuid.UUID `json:"authorId"`
	SubscriberCount int64     `json:"subscriberCount"`
}
