package service

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// EmailService defines the domain service for sending emails
type EmailService interface {
	SendNotification(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, data map[string]interface{}) error
	SendWelcomeEmail(ctx context.Context, userID uuid.UUID, email string, name string) error
	SendVerificationEmail(ctx context.Context, userID uuid.UUID, email string, token string) error
}
