package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BlogStatus represents the status of a blog post
type BlogStatus string

const (
	BlogStatusDraft     BlogStatus = "draft"
	BlogStatusPublished BlogStatus = "published"
)

// BlogVisibility represents the visibility mode of a blog post
type BlogVisibility string

const (
	BlogVisibilityPublic          BlogVisibility = "public"
	BlogVisibilitySubscribersOnly BlogVisibility = "subscribers_only"
)

// Blog represents a blog post entity
type Blog struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AuthorID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"authorId"`
	CategoryID   *uuid.UUID     `gorm:"type:uuid;index" json:"categoryId,omitempty"`
	Title        string         `gorm:"size:255;not null" json:"title"`
	Slug         string         `gorm:"size:255;not null;index" json:"slug"`
	Excerpt      *string        `gorm:"type:text" json:"excerpt,omitempty"`
	Content      string         `gorm:"type:text;not null" json:"content"`
	ThumbnailURL *string        `gorm:"size:500" json:"thumbnailUrl,omitempty"`
	Status       BlogStatus     `gorm:"type:blog_status;not null;default:'draft'" json:"status"`
	Visibility   BlogVisibility `gorm:"type:blog_visibility;not null;default:'public'" json:"visibility"`
	PublishedAt  *time.Time     `json:"publishedAt,omitempty"`
	CreatedAt    time.Time      `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	// Reactions (Denormalized counts)
	UpvoteCount   int `gorm:"not null;default:0" json:"upvoteCount"`
	DownvoteCount int `gorm:"not null;default:0" json:"downvoteCount"`

	// Relationships
	Author   *User     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags     []Tag     `gorm:"many2many:blog_tags" json:"tags,omitempty"`
	Comments []Comment `gorm:"foreignKey:BlogID" json:"comments,omitempty"`
}

// TableName returns the table name for Blog
func (Blog) TableName() string {
	return "blogs"
}

// IsPublished checks if the blog is published
func (b *Blog) IsPublished() bool {
	return b.Status == BlogStatusPublished
}

// IsScheduled checks if the blog is scheduled for future publication
func (b *Blog) IsScheduled() bool {
	return b.Status == BlogStatusPublished && b.PublishedAt != nil && b.PublishedAt.After(time.Now())
}

// IsDraft checks if the blog is a draft
func (b *Blog) IsDraft() bool {
	return b.Status == BlogStatusDraft
}

// IsPublic checks if the blog is public
func (b *Blog) IsPublic() bool {
	return b.Visibility == BlogVisibilityPublic
}

// IsSubscribersOnly checks if the blog is for subscribers only
func (b *Blog) IsSubscribersOnly() bool {
	return b.Visibility == BlogVisibilitySubscribersOnly
}

// Publish publishes the blog
func (b *Blog) Publish(at *time.Time) {
	if at != nil {
		b.PublishedAt = at
	} else {
		now := time.Now()
		b.PublishedAt = &now
	}
	b.Status = BlogStatusPublished
}

// Unpublish reverts the blog to draft
func (b *Blog) Unpublish() {
	b.Status = BlogStatusDraft
	b.PublishedAt = nil
}
