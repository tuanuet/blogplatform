package entity

import (
	"time"

	"github.com/google/uuid"
)

// UserRankingHistory represents a historical snapshot of a user's ranking
type UserRankingHistory struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	RankPosition   int       `gorm:"not null" json:"rankPosition"`
	CompositeScore float64   `gorm:"type:decimal(10,4);not null" json:"compositeScore"`
	FollowerCount  int       `gorm:"not null;default:0" json:"followerCount"`
	RecordedAt     time.Time `gorm:"not null;default:now()" json:"recordedAt"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for UserRankingHistory
func (UserRankingHistory) TableName() string {
	return "user_ranking_history"
}
