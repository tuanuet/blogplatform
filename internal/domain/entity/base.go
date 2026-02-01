package entity

import (
	"time"

	"github.com/google/uuid"
)

// BaseEntity contains common fields for all entities
type BaseEntity struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

// BeforeCreate sets the ID if not already set
func (b *BaseEntity) BeforeCreate() {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
}

// BeforeUpdate updates the UpdatedAt timestamp
func (b *BaseEntity) BeforeUpdate() {
	b.UpdatedAt = time.Now()
}
