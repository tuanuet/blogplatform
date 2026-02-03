package repository

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// UserSeriesPurchaseRepository defines the interface for user series purchase operations
type UserSeriesPurchaseRepository interface {
	Create(ctx context.Context, purchase *entity.UserSeriesPurchase) error
	HasPurchased(ctx context.Context, userID, seriesID uuid.UUID) (bool, error)
	GetUserPurchases(ctx context.Context, userID uuid.UUID) ([]*entity.UserSeriesPurchase, error)
}
