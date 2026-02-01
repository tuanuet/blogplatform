package entity

import (
	"time"

	"github.com/google/uuid"
)

// UserVelocityScore represents a user's velocity-based ranking score
type UserVelocityScore struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID             uuid.UUID `gorm:"type:uuid;not null;unique;index" json:"userId"`
	FollowerCount      int       `gorm:"not null;default:0" json:"followerCount"`
	FollowerGrowthRate float64   `gorm:"type:decimal(10,4);not null;default:0" json:"followerGrowthRate"`
	BlogPostVelocity   float64   `gorm:"type:decimal(10,4);not null;default:0" json:"blogPostVelocity"`
	CompositeScore     float64   `gorm:"type:decimal(10,4);not null;default:0" json:"compositeScore"`
	RankPosition       *int      `gorm:"index" json:"rankPosition,omitempty"`
	CalculationDate    time.Time `gorm:"not null;default:now()" json:"calculationDate"`
	CreatedAt          time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"not null;default:now()" json:"updatedAt"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for UserVelocityScore
func (UserVelocityScore) TableName() string {
	return "user_velocity_scores"
}

// GetRankChange calculates the rank change from a previous position
func (s *UserVelocityScore) GetRankChange(previousRank int) int {
	if s.RankPosition == nil {
		return 0
	}
	return previousRank - *s.RankPosition
}

// IsTopRanked checks if the user is in the top rankings
func (s *UserVelocityScore) IsTopRanked(topN int) bool {
	if s.RankPosition == nil {
		return false
	}
	return *s.RankPosition <= topN
}
