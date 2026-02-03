package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/domain/entity"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// FindByID finds a user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// Update updates a user
	Update(ctx context.Context, user *entity.User) error

	// UpdateProfile updates only profile fields
	UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error

	// Interest operations
	GetInterests(ctx context.Context, userID uuid.UUID) ([]entity.Tag, error)
	ReplaceInterests(ctx context.Context, userID uuid.UUID, tagIDs []uuid.UUID) error

	// Admin stats
	CountByMonth(ctx context.Context, months int) ([]entity.MonthlyCount, error)
}
