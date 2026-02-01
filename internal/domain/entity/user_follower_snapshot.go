package entity

import (
	"time"

	"github.com/google/uuid"
)

// UserFollowerSnapshot represents a daily snapshot of a user's follower count
type UserFollowerSnapshot struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	FollowerCount int       `gorm:"not null" json:"followerCount"`
	SnapshotDate  time.Time `gorm:"type:date;not null" json:"snapshotDate"`
	CreatedAt     time.Time `gorm:"not null;default:now()" json:"createdAt"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for UserFollowerSnapshot
func (UserFollowerSnapshot) TableName() string {
	return "user_follower_snapshots"
}

// GetSnapshotDateString returns the snapshot date as a string (YYYY-MM-DD)
func (s *UserFollowerSnapshot) GetSnapshotDateString() string {
	return s.SnapshotDate.Format("2006-01-02")
}
