package entity

import (
	"time"

	"github.com/google/uuid"
)

// Tag represents a blog tag
type Tag struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"size:50;not null;unique" json:"name"`
	Slug      string    `gorm:"size:50;not null;unique;index" json:"slug"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Blogs []Blog `gorm:"many2many:blog_tags" json:"blogs,omitempty"`
}

// TableName returns the table name for Tag
func (Tag) TableName() string {
	return "tags"
}
