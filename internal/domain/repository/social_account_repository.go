package repository

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// SocialAccountRepository defines the interface for social account persistence
type SocialAccountRepository interface {
	Create(ctx context.Context, socialAccount *entity.SocialAccount) error
	FindByProvider(ctx context.Context, provider string, providerID string) (*entity.SocialAccount, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.SocialAccount, error)
}
