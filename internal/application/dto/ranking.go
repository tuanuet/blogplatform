package dto

import (
	"time"

	"github.com/google/uuid"
)

// RankingFilterParams represents filter parameters for ranking queries
type RankingFilterParams struct {
	Page     int    `form:"page" json:"page"`                   // Page number (1-based)
	PageSize int    `form:"pageSize" json:"pageSize"`           // Items per page
	Category string `form:"category" json:"category,omitempty"` // Optional category filter
}

// RankedUserResponse represents a user in the ranking list
type RankedUserResponse struct {
	ID                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	DisplayName        string    `json:"displayName,omitempty"`
	AvatarURL          string    `json:"avatarUrl,omitempty"`
	FollowerCount      int       `json:"followerCount"`
	FollowerGrowthRate float64   `json:"followerGrowthRate"`
	BlogPostVelocity   float64   `json:"blogPostVelocity"`
	CompositeScore     float64   `json:"compositeScore"`
	Rank               *int      `json:"rank,omitempty"`
	CalculationDate    time.Time `json:"calculationDate"`
}

// UserRankingDetailResponse represents detailed ranking information for a user
type UserRankingDetailResponse struct {
	UserID             uuid.UUID             `json:"userId"`
	Username           string                `json:"username"`
	DisplayName        string                `json:"displayName,omitempty"`
	AvatarURL          string                `json:"avatarUrl,omitempty"`
	FollowerCount      int                   `json:"followerCount"`
	FollowerGrowthRate float64               `json:"followerGrowthRate"`
	BlogPostVelocity   float64               `json:"blogPostVelocity"`
	CompositeScore     float64               `json:"compositeScore"`
	Rank               *int                  `json:"rank,omitempty"`
	PreviousRank       *int                  `json:"previousRank,omitempty"`
	RankChange         int                   `json:"rankChange"`
	CalculationDate    time.Time             `json:"calculationDate"`
	History            []RankingHistoryEntry `json:"history"`
}

// RankingHistoryEntry represents a single ranking history record
type RankingHistoryEntry struct {
	RankPosition   int       `json:"rankPosition"`
	CompositeScore float64   `json:"compositeScore"`
	FollowerCount  int       `json:"followerCount"`
	RecordedAt     time.Time `json:"recordedAt"`
}
