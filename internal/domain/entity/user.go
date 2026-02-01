package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity
type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string     `gorm:"size:255;not null;unique;index" json:"email"`
	Name         string     `gorm:"size:255;not null" json:"name"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	IsActive     bool       `gorm:"not null;default:true" json:"isActive"`
	CreatedAt    time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Profile fields
	DisplayName   *string    `gorm:"size:50" json:"displayName,omitempty"`
	Bio           *string    `gorm:"type:text" json:"bio,omitempty"`
	AvatarURL     *string    `gorm:"size:500" json:"avatarUrl,omitempty"`
	Website       *string    `gorm:"size:255" json:"website,omitempty"`
	Location      *string    `gorm:"size:100" json:"location,omitempty"`
	TwitterHandle *string    `gorm:"size:50" json:"twitterHandle,omitempty"`
	GithubHandle  *string    `gorm:"size:50" json:"githubHandle,omitempty"`
	LinkedinURL   *string    `gorm:"size:255" json:"linkedinUrl,omitempty"`
	Gender        *string    `gorm:"size:10" json:"gender,omitempty"`
	Birthday      *time.Time `gorm:"type:date" json:"birthday,omitempty"`

	// Relationships
	Blogs           []Blog         `gorm:"foreignKey:AuthorID" json:"blogs,omitempty"`
	BookmarkedBlogs []Blog         `gorm:"many2many:user_bookmarks;joinForeignKey:user_id;joinReferences:blog_id" json:"bookmarkedBlogs,omitempty"`
	Comments        []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	Subscriptions   []Subscription `gorm:"foreignKey:SubscriberID" json:"subscriptions,omitempty"`
	Subscribers     []Subscription `gorm:"foreignKey:AuthorID" json:"subscribers,omitempty"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// GetDisplayName returns display name or falls back to name
func (u *User) GetDisplayName() string {
	if u.DisplayName != nil && *u.DisplayName != "" {
		return *u.DisplayName
	}
	return u.Name
}

// HasAvatar returns true if user has an avatar
func (u *User) HasAvatar() bool {
	return u.AvatarURL != nil && *u.AvatarURL != ""
}
