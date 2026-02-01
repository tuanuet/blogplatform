package entity

import (
	"time"

	"github.com/google/uuid"
)

// ReactionType represents the type of user reaction
type ReactionType string

const (
	ReactionTypeUpvote   ReactionType = "upvote"
	ReactionTypeDownvote ReactionType = "downvote"
	ReactionTypeNone     ReactionType = "none" // Used for requests to remove reaction
)

// BlogReaction tracks user reactions (upvotes/downvotes) to blogs
type BlogReaction struct {
	ID        uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BlogID    uuid.UUID    `gorm:"type:uuid;not null;index:idx_blog_reaction_unique,unique" json:"blogId"`
	UserID    uuid.UUID    `gorm:"type:uuid;not null;index:idx_blog_reaction_unique,unique" json:"userId"`
	Type      ReactionType `gorm:"type:varchar(20);not null" json:"type"`
	CreatedAt time.Time    `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time    `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	Blog *Blog `gorm:"foreignKey:BlogID" json:"blog,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for BlogReaction
func (BlogReaction) TableName() string {
	return "blog_reactions"
}
