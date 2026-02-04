package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID           uuid.UUID              `json:"id"`
	UserID       uuid.UUID              `json:"user_id"`
	Type         string                 `json:"type"`
	Category     string                 `json:"category"`
	Title        string                 `json:"title"`
	Body         string                 `json:"body"`
	Data         map[string]interface{} `json:"data"`
	GroupedCount int                    `json:"grouped_count"`
	IsRead       bool                   `json:"is_read"`
	CreatedAt    time.Time              `json:"created_at"`
}

type UnreadCountResponse struct {
	Count int `json:"count"`
}

type NotificationPreferenceResponse struct {
	NotificationType NotificationType `json:"notification_type"`
	Channel          string           `json:"channel"`
	Enabled          bool             `json:"enabled"`
}

type NotificationType string

type UpdatePreferencesRequest struct {
	Preferences []NotificationPreferenceItem `json:"preferences" binding:"required,dive"`
}

type NotificationPreferenceItem struct {
	NotificationType NotificationType `json:"notification_type" binding:"required"`
	Channel          string           `json:"channel" binding:"required"`
	Enabled          bool             `json:"enabled"`
}

type RegisterDeviceTokenRequest struct {
	DeviceToken string `json:"device_token" binding:"required"`
	Platform    string `json:"platform" binding:"required,oneof=ios android web"`
}
