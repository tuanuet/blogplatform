package entity

import (
	"time"

	"github.com/google/uuid"
)

// UserReadingHistory represents a record of a user reading a blog post
type UserReadingHistory struct {
	UserID     uuid.UUID `gorm:"type:uuid;primary_key" json:"userId"`
	BlogID     uuid.UUID `gorm:"type:uuid;primary_key" json:"blogId"`
	LastReadAt time.Time `gorm:"not null" json:"lastReadAt"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Blog *Blog `gorm:"foreignKey:BlogID" json:"blog,omitempty"`
}

// TableName returns the table name for UserReadingHistory
func (UserReadingHistory) TableName() string {
	return "user_reading_history"
}
