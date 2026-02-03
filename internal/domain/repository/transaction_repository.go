package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// TransactionRepository defines the interface for transaction data operations
type TransactionRepository interface {
	// Create creates a new transaction
	Create(ctx context.Context, tx *entity.Transaction) error

	// FindByRefID finds a transaction by reference code
	FindByRefID(ctx context.Context, refID string) (*entity.Transaction, error)

	// FindBySePayID finds a transaction by SePay ID (for idempotency)
	FindBySePayID(ctx context.Context, sePayID string) (*entity.Transaction, error)

	// UpdateStatus updates the status of a transaction
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TransactionStatus) error

	// Update updates a transaction
	Update(ctx context.Context, tx *entity.Transaction) error

	// FindByUserID finds all transactions for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Transaction, error)

	// WithTx returns a new repository with the given transaction
	WithTx(tx interface{}) TransactionRepository
}
