package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// SubscriptionTier represents the subscription tier levels
type SubscriptionTier string

// SubscriptionTier enum values
const (
	TierFree   SubscriptionTier = "FREE"
	TierBronze SubscriptionTier = "BRONZE"
	TierSilver SubscriptionTier = "SILVER"
	TierGold   SubscriptionTier = "GOLD"
)

// Level returns the numeric level of the tier for comparison
func (t SubscriptionTier) Level() int {
	switch t {
	case TierFree:
		return 0
	case TierBronze:
		return 1
	case TierSilver:
		return 2
	case TierGold:
		return 3
	default:
		return 0
	}
}

// IsValid returns true if the tier is a valid tier value
func (t SubscriptionTier) IsValid() bool {
	switch t {
	case TierFree, TierBronze, TierSilver, TierGold:
		return true
	default:
		return false
	}
}

// String returns the string representation of the tier
func (t SubscriptionTier) String() string {
	return string(t)
}

// SubscriptionPlan represents a subscription pricing plan created by an author
type SubscriptionPlan struct {
	ID           uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AuthorID     uuid.UUID        `gorm:"type:uuid;not null;index:idx_author_tier" json:"authorId"`
	Tier         SubscriptionTier `gorm:"type:varchar(20);not null;index:idx_author_tier" json:"tier"`
	Price        decimal.Decimal  `gorm:"type:decimal(15,2);not null" json:"price"`
	DurationDays int              `gorm:"not null;default:30" json:"durationDays"`
	Name         *string          `gorm:"type:varchar(100)" json:"name,omitempty"`
	Description  *string          `gorm:"type:text" json:"description,omitempty"`
	IsActive     bool             `gorm:"not null;default:true" json:"isActive"`
	CreatedAt    time.Time        `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time        `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`

	// Relationships
	Author *User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

// TableName returns the table name for SubscriptionPlan
func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}
