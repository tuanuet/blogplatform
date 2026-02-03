package repository

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userSeriesPurchaseRepository struct {
	db *gorm.DB
}

// NewUserSeriesPurchaseRepository creates a new instance of UserSeriesPurchaseRepository
func NewUserSeriesPurchaseRepository(db *gorm.DB) repository.UserSeriesPurchaseRepository {
	return &userSeriesPurchaseRepository{db: db}
}

// Create creates a new user series purchase record
func (r *userSeriesPurchaseRepository) Create(ctx context.Context, purchase *entity.UserSeriesPurchase) error {
	return r.db.WithContext(ctx).Create(purchase).Error
}

// HasPurchased checks if a user has purchased a series
func (r *userSeriesPurchaseRepository) HasPurchased(ctx context.Context, userID, seriesID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.UserSeriesPurchase{}).
		Where("user_id = ? AND series_id = ?", userID, seriesID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserPurchases returns all series purchased by a user, including preloaded Series entity
func (r *userSeriesPurchaseRepository) GetUserPurchases(ctx context.Context, userID uuid.UUID) ([]*entity.UserSeriesPurchase, error) {
	var purchases []*entity.UserSeriesPurchase
	err := r.db.WithContext(ctx).
		Preload("Series").
		Where("user_id = ?", userID).
		Find(&purchases).Error
	if err != nil {
		return nil, err
	}
	return purchases, nil
}

// WithTx returns a new repository with the given transaction
func (r *userSeriesPurchaseRepository) WithTx(tx interface{}) repository.UserSeriesPurchaseRepository {
	if gormDB, ok := tx.(*gorm.DB); ok {
		return &userSeriesPurchaseRepository{db: gormDB}
	}
	return r
}
