package entity

import (
	"time"

	"github.com/google/uuid"
)

// Subscription represents a user subscription to an author
type Subscription struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SubscriberID uuid.UUID `gorm:"type:uuid;not null;index" json:"subscriberId"`
	AuthorID     uuid.UUID `gorm:"type:uuid;not null;index" json:"authorId"`

	// Fields for Paid Subscriptions
	ExpiresAt *time.Time `gorm:"index" json:"expiresAt,omitempty"`   // Null = Free Follower
	Tier      string     `gorm:"size:20;default:'FREE'" json:"tier"` // FREE, PREMIUM, VIP

	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Subscriber *User `gorm:"foreignKey:SubscriberID" json:"subscriber,omitempty"`
	Author     *User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

// TableName returns the table name for Subscription
func (Subscription) TableName() string {
	return "subscriptions"
}

// IsActive returns true if subscription is active (free follower or paid and not expired)
func (s *Subscription) IsActive() bool {
	if s.ExpiresAt == nil {
		return true // Free follower
	}
	return s.ExpiresAt.After(time.Now())
}

// IsPaid returns true if this is a paid subscription
func (s *Subscription) IsPaid() bool {
	return s.ExpiresAt != nil
}

// IsExpired returns true if paid subscription has expired
func (s *Subscription) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false // Free followers don't expire
	}
	return s.ExpiresAt.Before(time.Now())
}
