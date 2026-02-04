package repository

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type subscriptionPlanRepository struct {
	db *gorm.DB
}

// NewSubscriptionPlanRepository creates a new subscription plan repository
func NewSubscriptionPlanRepository(db *gorm.DB) repository.SubscriptionPlanRepository {
	return &subscriptionPlanRepository{db: db}
}

// Create creates a new subscription plan
func (r *subscriptionPlanRepository) Create(ctx context.Context, plan *entity.SubscriptionPlan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

// FindByID retrieves a subscription plan by ID
func (r *subscriptionPlanRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.SubscriptionPlan, error) {
	var plan entity.SubscriptionPlan
	err := r.db.WithContext(ctx).First(&plan, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}

// FindByAuthorAndTier retrieves a plan by author ID and tier
func (r *subscriptionPlanRepository) FindByAuthorAndTier(ctx context.Context, authorID uuid.UUID, tier entity.SubscriptionTier) (*entity.SubscriptionPlan, error) {
	var plan entity.SubscriptionPlan
	err := r.db.WithContext(ctx).
		Where("author_id = ? AND tier = ?", authorID, tier).
		First(&plan).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}

// FindByAuthor retrieves all plans for an author
func (r *subscriptionPlanRepository) FindByAuthor(ctx context.Context, authorID uuid.UUID) ([]entity.SubscriptionPlan, error) {
	var plans []entity.SubscriptionPlan
	err := r.db.WithContext(ctx).
		Where("author_id = ?", authorID).
		Find(&plans).Error
	return plans, err
}

// FindActiveByAuthor retrieves all active plans for an author
func (r *subscriptionPlanRepository) FindActiveByAuthor(ctx context.Context, authorID uuid.UUID) ([]entity.SubscriptionPlan, error) {
	var plans []entity.SubscriptionPlan
	err := r.db.WithContext(ctx).
		Where("author_id = ? AND is_active = ?", authorID, true).
		Find(&plans).Error
	return plans, err
}

// Update updates an existing subscription plan
func (r *subscriptionPlanRepository) Update(ctx context.Context, plan *entity.SubscriptionPlan) error {
	return r.db.WithContext(ctx).Save(plan).Error
}

// Upsert creates or updates a subscription plan based on (author_id, tier) uniqueness
func (r *subscriptionPlanRepository) Upsert(ctx context.Context, plan *entity.SubscriptionPlan) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "author_id"}, {Name: "tier"}},
			DoUpdates: clause.AssignmentColumns([]string{"price", "duration_days", "name", "description", "is_active", "updated_at", "deleted_at"}),
		}).
		Create(plan).Error
}

// Delete soft deletes a subscription plan
func (r *subscriptionPlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entity.SubscriptionPlan{}, id)
	return result.Error
}

// WithTx returns a new repository with the given transaction
func (r *subscriptionPlanRepository) WithTx(tx interface{}) repository.SubscriptionPlanRepository {
	if gormDB, ok := tx.(*gorm.DB); ok {
		return &subscriptionPlanRepository{db: gormDB}
	}
	return r
}
