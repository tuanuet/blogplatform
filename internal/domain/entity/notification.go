package entity

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeNewFollower          NotificationType = "new_follower"
	NotificationTypeBlogLike             NotificationType = "blog_like"
	NotificationTypeBlogComment          NotificationType = "blog_comment"
	NotificationTypeCommentReply         NotificationType = "comment_reply"
	NotificationTypeMention              NotificationType = "mention"
	NotificationTypeNewBlogFromFollowing NotificationType = "new_blog_from_following"
	NotificationTypeSeriesUpdate         NotificationType = "series_update"
	NotificationTypeBotFollowerDetected  NotificationType = "bot_follower_detected"
	NotificationTypeBadgeStatusChange    NotificationType = "badge_status_change"
)

type NotificationCategory string

const (
	NotificationCategorySocial  NotificationCategory = "social"
	NotificationCategoryContent NotificationCategory = "content"
	NotificationCategorySystem  NotificationCategory = "system"
)

// Notification represents an in-app notification
type Notification struct {
	ID           uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID              `gorm:"type:uuid;not null;index" json:"user_id"`
	Type         NotificationType       `gorm:"size:50;not null" json:"type"`
	Category     NotificationCategory   `gorm:"size:20;not null" json:"category"`
	Title        string                 `gorm:"size:255;not null" json:"title"`
	Body         string                 `gorm:"type:text;not null" json:"body"`
	Data         map[string]interface{} `gorm:"type:jsonb;default:'{}'" json:"data"`
	GroupedCount int                    `gorm:"not null;default:1" json:"grouped_count"`
	IsRead       bool                   `gorm:"not null;default:false" json:"is_read"`
	CreatedAt    time.Time              `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time              `gorm:"not null;default:now()" json:"updated_at"`
	ExpiresAt    time.Time              `gorm:"not null;index" json:"expires_at"`
}

// TableName returns the table name for Notification
func (Notification) TableName() string {
	return "notifications"
}

// NotificationPreference represents a user's notification settings
type NotificationPreference struct {
	UserID           uuid.UUID        `gorm:"type:uuid;not null;primaryKey" json:"user_id"`
	NotificationType NotificationType `gorm:"size:50;not null;primaryKey" json:"notification_type"`
	Channel          string           `gorm:"size:20;not null;primaryKey;default:'in_app'" json:"channel"`
	Enabled          bool             `gorm:"not null;default:true" json:"enabled"`
}

// TableName returns the table name for NotificationPreference
func (NotificationPreference) TableName() string {
	return "notification_preferences"
}

// UserDeviceToken represents a registered device for push notifications
type UserDeviceToken struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	DeviceToken string    `gorm:"size:255;not null" json:"device_token"`
	Platform    string    `gorm:"size:20;not null" json:"platform"`
	LastSeenAt  time.Time `gorm:"not null;default:now()" json:"last_seen_at"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// TableName returns the table name for UserDeviceToken
func (UserDeviceToken) TableName() string {
	return "user_device_tokens"
}

// NotifyRequest represents a request to send a notification
type NotifyRequest struct {
	TargetUserID uuid.UUID
	Type         NotificationType
	Data         map[string]interface{}
}
