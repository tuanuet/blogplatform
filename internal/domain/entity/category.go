package entity

import (
	"time"

	"github.com/google/uuid"
)

// Category represents a blog category
type Category struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string     `gorm:"size:100;not null;unique" json:"name"`
	Slug        string     `gorm:"size:100;not null;unique;index" json:"slug"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	Blogs []Blog `gorm:"foreignKey:CategoryID" json:"blogs,omitempty"`
}

// TableName returns the table name for Category
func (Category) TableName() string {
	return "categories"
}
