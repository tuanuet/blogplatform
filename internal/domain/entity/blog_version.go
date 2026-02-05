package entity

import (
	"time"

	"github.com/google/uuid"
)

// BlogVersion represents a version of a blog post
type BlogVersion struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BlogID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"blogId"`
	VersionNumber int            `gorm:"not null" json:"versionNumber"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	Slug          string         `gorm:"size:255;not null" json:"slug"`
	Excerpt       *string        `gorm:"type:text" json:"excerpt,omitempty"`
	Content       string         `gorm:"type:text;not null" json:"content"`
	ThumbnailURL  *string        `gorm:"size:500" json:"thumbnailUrl,omitempty"`
	Status        BlogStatus     `gorm:"type:blog_status;not null" json:"status"`
	Visibility    BlogVisibility `gorm:"type:blog_visibility;not null" json:"visibility"`
	CategoryID    *uuid.UUID     `gorm:"type:uuid" json:"categoryId,omitempty"`
	EditorID      uuid.UUID      `gorm:"type:uuid;not null" json:"editorId"`
	ChangeSummary *string        `gorm:"type:text" json:"changeSummary,omitempty"`
	CreatedAt     time.Time      `gorm:"not null;default:now()" json:"createdAt"`

	// Relationships
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Editor   *User     `gorm:"foreignKey:EditorID" json:"editor,omitempty"`
	Tags     []Tag     `gorm:"many2many:blog_version_tags" json:"tags,omitempty"`
}

// TableName returns the table name for BlogVersion
func (BlogVersion) TableName() string {
	return "blog_versions"
}
