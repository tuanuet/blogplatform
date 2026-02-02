package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"errors"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

// Domain errors
var (
	ErrUserNotFound = errors.New("user not found")
)

// UserService handles user-related domain logic
type UserService interface {
	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// UpdateUser updates user fields
	UpdateUser(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error

	// UpdateAvatarURL updates the user's avatar URL
	UpdateAvatarURL(ctx context.Context, id uuid.UUID, avatarURL string) error
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user domain service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	// Verify user exists first
	if _, err := s.GetUser(ctx, id); err != nil {
		return err
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		return s.userRepo.UpdateProfile(ctx, id, updates)
	}
	return nil
}

func (s *userService) UpdateAvatarURL(ctx context.Context, id uuid.UUID, avatarURL string) error {
	updates := map[string]interface{}{
		"avatar_url": avatarURL,
	}
	return s.UpdateUser(ctx, id, updates)
}
