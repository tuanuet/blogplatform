package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PlanManagementService defines the interface for subscription plan business logic
type PlanManagementService interface {
	// UpsertPlans creates or updates multiple subscription plans
	// Returns created/updated plans and any warnings (e.g., price hierarchy violations)
	UpsertPlans(ctx context.Context, authorID uuid.UUID, plans []CreatePlanDTO) ([]entity.SubscriptionPlan, []string, error)

	// GetAuthorPlans retrieves all plans for an author with tag information
	GetAuthorPlans(ctx context.Context, authorID uuid.UUID) ([]PlanWithTags, error)

	// DeactivatePlan deactivates a plan without deleting it
	DeactivatePlan(ctx context.Context, authorID uuid.UUID, tier entity.SubscriptionTier) error

	// ActivatePlan activates a previously deactivated plan
	ActivatePlan(ctx context.Context, authorID uuid.UUID, tier entity.SubscriptionTier) error
}

// CreatePlanDTO represents data for creating/updating a plan
type CreatePlanDTO struct {
	Tier         entity.SubscriptionTier
	Price        decimal.Decimal
	Name         *string
	Description  *string
	DurationDays int
}

// PlanWithTags represents a plan with associated tag information
type PlanWithTags struct {
	Plan     entity.SubscriptionPlan
	Tags     []string
	TagCount int
}
