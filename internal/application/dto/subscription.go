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

// SubscriptionCountResponse represents subscriber and subscription counts
type SubscriptionCountResponse struct {
	AuthorID          uuid.UUID `json:"authorId,omitempty"`
	SubscriberCount   int64     `json:"subscriberCount"`
	SubscriptionCount int64     `json:"subscriptionCount,omitempty"`
}
