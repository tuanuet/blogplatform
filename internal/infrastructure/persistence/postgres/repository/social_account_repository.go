package repository

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type socialAccountRepository struct {
	db *gorm.DB
}

// NewSocialAccountRepository creates a new social account repository
func NewSocialAccountRepository(db *gorm.DB) repository.SocialAccountRepository {
	return &socialAccountRepository{db: db}
}

func (r *socialAccountRepository) Create(ctx context.Context, socialAccount *entity.SocialAccount) error {
	return r.db.WithContext(ctx).Create(socialAccount).Error
}

func (r *socialAccountRepository) FindByProvider(ctx context.Context, provider string, providerID string) (*entity.SocialAccount, error) {
	var socialAccount entity.SocialAccount
	err := r.db.WithContext(ctx).
		Where("provider = ? AND provider_id = ?", provider, providerID).
		First(&socialAccount).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &socialAccount, nil
}

func (r *socialAccountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.SocialAccount, error) {
	var socialAccounts []*entity.SocialAccount
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&socialAccounts).Error
	if err != nil {
		return nil, err
	}
	return socialAccounts, nil
}
