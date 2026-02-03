package repository

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionRepository{db: db}
}

// Create creates a new transaction
func (r *transactionRepository) Create(ctx context.Context, tx *entity.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

// FindByRefID finds a transaction by reference code
func (r *transactionRepository) FindByRefID(ctx context.Context, refID string) (*entity.Transaction, error) {
	var tx entity.Transaction
	err := r.db.WithContext(ctx).
		Where("reference_code = ?", refID).
		First(&tx).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &tx, err
}

// FindBySePayID finds a transaction by SePay ID (for idempotency)
func (r *transactionRepository) FindBySePayID(ctx context.Context, sePayID string) (*entity.Transaction, error) {
	var tx entity.Transaction
	err := r.db.WithContext(ctx).
		Where("sepay_id = ?", sePayID).
		First(&tx).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &tx, err
}

// UpdateStatus updates the status of a transaction
func (r *transactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TransactionStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Transaction{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// FindByUserID finds all transactions for a user
func (r *transactionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}
