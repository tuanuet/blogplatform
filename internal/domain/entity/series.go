package entity

import (
	"time"

	"github.com/google/uuid"
)

// Series represents a collection of blog posts
type Series struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AuthorID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"authorId"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Slug        string     `gorm:"size:255;not null;index" json:"slug"`
	Description string     `gorm:"type:text" json:"description"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	Author *User  `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Blogs  []Blog `gorm:"many2many:series_blogs;" json:"blogs,omitempty"`
}

// TableName returns the table name for Series
func (Series) TableName() string {
	return "series"
}
