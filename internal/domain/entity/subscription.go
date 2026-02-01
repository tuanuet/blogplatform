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
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"createdAt"`

	// Relationships
	Subscriber *User `gorm:"foreignKey:SubscriberID" json:"subscriber,omitempty"`
	Author     *User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

// TableName returns the table name for Subscription
func (Subscription) TableName() string {
	return "subscriptions"
}
