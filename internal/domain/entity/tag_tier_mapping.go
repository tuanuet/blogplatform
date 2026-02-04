package entity

import (
	"time"

	"github.com/google/uuid"
)

// TagTierMapping represents the mapping between a tag and a required subscription tier
type TagTierMapping struct {
	ID           uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AuthorID     uuid.UUID        `gorm:"type:uuid;not null;index:idx_author_tag" json:"authorId"`
	TagID        uuid.UUID        `gorm:"type:uuid;not null;index:idx_author_tag" json:"tagId"`
	RequiredTier SubscriptionTier `gorm:"type:varchar(20);not null" json:"requiredTier"`
	CreatedAt    time.Time        `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time        `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Author *User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Tag    *Tag  `gorm:"foreignKey:TagID" json:"tag,omitempty"`
}

// TableName returns the table name for TagTierMapping
func (TagTierMapping) TableName() string {
	return "tag_tier_mappings"
}
