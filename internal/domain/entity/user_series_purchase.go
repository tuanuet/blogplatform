package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// UserSeriesPurchase represents a user's purchase of a series
type UserSeriesPurchase struct {
	UserID    uuid.UUID       `gorm:"type:uuid;primaryKey" json:"userId"`
	SeriesID  uuid.UUID       `gorm:"type:uuid;primaryKey" json:"seriesId"`
	Amount    decimal.Decimal `gorm:"type:decimal(20,2);not null" json:"amount"`
	CreatedAt time.Time       `gorm:"not null;default:now()" json:"createdAt"`

	// Relationships
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Series *Series `gorm:"foreignKey:SeriesID" json:"series,omitempty"`
}

// TableName returns the table name for UserSeriesPurchase
func (UserSeriesPurchase) TableName() string {
	return "user_series_purchases"
}
