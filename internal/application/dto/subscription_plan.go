package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ===== Subscription Plan DTOs =====

// CreatePlanRequest represents a single plan in upsert request
type CreatePlanRequest struct {
	Tier        string          `json:"tier" binding:"required,oneof=BRONZE SILVER GOLD" validate:"required,oneof=BRONZE SILVER GOLD"`
	Price       decimal.Decimal `json:"price" binding:"required" validate:"required,gte=0"`
	Name        *string         `json:"name,omitempty" validate:"omitempty,max=100"`
	Description *string         `json:"description,omitempty"`
}

// UpsertPlansRequest represents the request to create/update subscription plans
type UpsertPlansRequest struct {
	Plans []CreatePlanRequest `json:"plans" validate:"required,min=1,max=3,dive"`
}

// PlanResponse represents a subscription plan in API responses
type PlanResponse struct {
	ID           uuid.UUID       `json:"id"`
	Tier         string          `json:"tier"`
	Price        decimal.Decimal `json:"price"`
	DurationDays int             `json:"durationDays"`
	Name         *string         `json:"name,omitempty"`
	Description  *string         `json:"description,omitempty"`
	IsActive     bool            `json:"isActive"`
	CreatedAt    string          `json:"createdAt"`
	UpdatedAt    string          `json:"updatedAt"`
}

// UpsertPlansResponse represents the response after upserting plans
type UpsertPlansResponse struct {
	Plans    []PlanResponse `json:"plans"`
	Warnings []string       `json:"warnings,omitempty"`
}

// PlanWithTagsResponse represents a plan with associated tags and counts
type PlanWithTagsResponse struct {
	Tier         string          `json:"tier"`
	Price        decimal.Decimal `json:"price"`
	DurationDays int             `json:"durationDays"`
	Name         *string         `json:"name,omitempty"`
	Description  *string         `json:"description,omitempty"`
	TagCount     int             `json:"tagCount"`
	Tags         []string        `json:"tags,omitempty"`
}

// GetAuthorPlansResponse represents the response for getting an author's pricing table
type GetAuthorPlansResponse struct {
	AuthorID uuid.UUID              `json:"authorId"`
	Plans    []PlanWithTagsResponse `json:"plans"`
}

// ===== Tag-Tier Mapping DTOs =====

// AssignTagTierRequest represents the request to assign a tag to a tier
type AssignTagTierRequest struct {
	RequiredTier string `json:"requiredTier" binding:"required,oneof=FREE BRONZE SILVER GOLD" validate:"required,oneof=FREE BRONZE SILVER GOLD"`
}

// AssignTagTierResponse represents the response after assigning a tag to a tier
type AssignTagTierResponse struct {
	TagID              uuid.UUID `json:"tagId"`
	TagName            string    `json:"tagName"`
	RequiredTier       string    `json:"requiredTier"`
	AffectedBlogsCount int64     `json:"affectedBlogsCount"`
}

// UnassignTagTierResponse represents the response after unassigning a tag from a tier
type UnassignTagTierResponse struct {
	Message            string `json:"message"`
	AffectedBlogsCount int64  `json:"affectedBlogsCount"`
}

// TagTierMappingResponse represents a tag-tier mapping in API responses
type TagTierMappingResponse struct {
	TagID        uuid.UUID `json:"tagId"`
	TagName      string    `json:"tagName"`
	RequiredTier string    `json:"requiredTier"`
	BlogCount    int64     `json:"blogCount"`
}

// GetTagTiersResponse represents the response for getting all tag-tier mappings
type GetTagTiersResponse struct {
	Mappings []TagTierMappingResponse `json:"mappings"`
}

// ===== Blog Access DTOs =====

// UpgradeOption represents an available upgrade option
type UpgradeOption struct {
	Tier         string          `json:"tier"`
	Price        decimal.Decimal `json:"price"`
	DurationDays int             `json:"durationDays"`
	PlanID       uuid.UUID       `json:"planId"`
}

// CheckBlogAccessResponse represents the response for checking blog access
type CheckBlogAccessResponse struct {
	Accessible     bool            `json:"accessible"`
	UserTier       string          `json:"userTier"`
	RequiredTier   string          `json:"requiredTier"`
	Reason         string          `json:"reason"`
	UpgradeOptions []UpgradeOption `json:"upgradeOptions,omitempty"`
}
