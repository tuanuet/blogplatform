package entity

import (
	"time"

	"github.com/google/uuid"
)

// SocialAccount represents a linked social media account
type SocialAccount struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	Provider   string    `gorm:"size:50;not null;index:idx_social_provider_id" json:"provider"`
	ProviderID string    `gorm:"size:255;not null;index:idx_social_provider_id" json:"providerId"`
	Email      string    `gorm:"size:255" json:"email"`
	CreatedAt  time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// TableName returns the table name for SocialAccount
func (SocialAccount) TableName() string {
	return "social_accounts"
}
