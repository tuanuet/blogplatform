package usecase

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	LoginWithSocial(ctx context.Context, req dto.LoginWithSocialRequest) (*dto.AuthResponse, error)
	GetSocialAuthURL(ctx context.Context, provider string) (string, error)
	Logout(ctx context.Context, sessionID string) error
	VerifyEmail(ctx context.Context, token string) error
}

type authUseCase struct {
	userRepo          repository.UserRepository
	sessionRepo       repository.SessionRepository
	socialRepo        repository.SocialAccountRepository
	socialAuthService service.SocialAuthService
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	socialRepo repository.SocialAccountRepository,
	socialAuthService service.SocialAuthService,
) AuthUseCase {
	return &authUseCase{
		userRepo:          userRepo,
		sessionRepo:       sessionRepo,
		socialRepo:        socialRepo,
		socialAuthService: socialAuthService,
	}
}

func (u *authUseCase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entity.User{
		ID:           uuid.New(), // Generate ID here to ensure we have it
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate verification token
	verificationToken := uuid.New().String()
	// Store token in session repo (reuse generic session storage)
	// Using a prefix or separate store would be better, but for now we reuse session repo
	if err := u.sessionRepo.CreateSession(ctx, verificationToken, user.ID.String(), 24*time.Hour); err != nil {
		return nil, fmt.Errorf("failed to create verification token: %w", err)
	}

	// Stub: Send email
	fmt.Printf("[STUB] Sending verification email to %s with token: %s\n", user.Email, verificationToken)

	return &dto.AuthResponse{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
	}, nil
}

func (u *authUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if email is verified
	if user.EmailVerifiedAt == nil {
		return nil, errors.New("email not verified")
	}

	// Generate session
	sessionID := uuid.New().String()
	// Create session (24h validity)
	if err := u.sessionRepo.CreateSession(ctx, sessionID, user.ID.String(), 24*time.Hour); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		SessionID: sessionID,
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
	}, nil
}

func (u *authUseCase) Logout(ctx context.Context, sessionID string) error {
	return u.sessionRepo.DeleteSession(ctx, sessionID)
}

func (u *authUseCase) VerifyEmail(ctx context.Context, token string) error {
	userIDStr, err := u.sessionRepo.GetUserID(ctx, token)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.New("invalid user ID in token")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if user.EmailVerifiedAt != nil {
		return nil
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	user.UpdatedAt = now

	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Invalidate token
	_ = u.sessionRepo.DeleteSession(ctx, token)

	return nil
}

func (u *authUseCase) LoginWithSocial(ctx context.Context, req dto.LoginWithSocialRequest) (*dto.AuthResponse, error) {
	// 1. Get User Info from Provider
	socialInfo, err := u.socialAuthService.GetUserInfo(ctx, req.Provider, req.Code)
	if err != nil {
		return nil, err
	}

	// 2. Check if SocialAccount exists
	socialAccount, err := u.socialRepo.FindByProvider(ctx, req.Provider, socialInfo.ProviderID)
	if err != nil {
		return nil, err
	}

	var user *entity.User

	if socialAccount != nil {
		// Social account exists, find the user
		user, err = u.userRepo.FindByID(ctx, socialAccount.UserID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("user not found for social account")
		}
	} else {
		// Social account does not exist
		// Check if user exists by email
		existingUser, err := u.userRepo.FindByEmail(ctx, socialInfo.Email)
		if err != nil {
			return nil, err
		}

		if existingUser != nil {
			// Link account to existing user
			user = existingUser
		} else {
			// Create new user
			user = &entity.User{
				ID:           uuid.New(),
				Name:         socialInfo.Name,
				Email:        socialInfo.Email,
				PasswordHash: "", // No password for social login users initially
				IsActive:     true,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			if err := u.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}
		}

		// Create SocialAccount
		newSocialAccount := &entity.SocialAccount{
			UserID:     user.ID,
			Provider:   req.Provider,
			ProviderID: socialInfo.ProviderID,
			Email:      socialInfo.Email,
		}
		if err := u.socialRepo.Create(ctx, newSocialAccount); err != nil {
			return nil, err
		}
	}

	// 3. Create Session
	sessionID := uuid.New().String()
	if err := u.sessionRepo.CreateSession(ctx, sessionID, user.ID.String(), 24*time.Hour); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		SessionID: sessionID,
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
	}, nil
}

func (u *authUseCase) GetSocialAuthURL(ctx context.Context, provider string) (string, error) {
	return u.socialAuthService.GetAuthURL(ctx, provider)
}
