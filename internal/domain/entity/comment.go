package entity

import (
	"time"

	"github.com/google/uuid"
)

// Comment represents a blog comment
type Comment struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BlogID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"blogId"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"userId"`
	ParentID  *uuid.UUID `gorm:"type:uuid;index" json:"parentId,omitempty"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	Blog    *Blog     `gorm:"foreignKey:BlogID" json:"blog,omitempty"`
	User    *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Parent  *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// TableName returns the table name for Comment
func (Comment) TableName() string {
	return "comments"
}

// IsReply checks if this comment is a reply to another comment
func (c *Comment) IsReply() bool {
	return c.ParentID != nil
}
