package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/auth"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/aiagent/internal/domain/service"
	serviceMocks "github.com/aiagent/internal/domain/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUseCase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	authUC := auth.NewAuthUseCase(mockUserRepo, mockSessionRepo, nil, nil)

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		req := dto.RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		// Expect FindByEmail to return nil (user not found)
		mockUserRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)

		// Expect Create to be called
		mockUserRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, u *entity.User) error {
			assert.Equal(t, req.Name, u.Name)
			assert.Equal(t, req.Email, u.Email)
			assert.NotEmpty(t, u.PasswordHash)
			// Verify password hash
			err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))
			assert.NoError(t, err)
			return nil
		})

		// Expect Session Create for Verification Token (Logic to be added)
		mockSessionRepo.EXPECT().CreateSession(ctx, gomock.Any(), gomock.Any(), 24*time.Hour).Return(nil)

		resp, err := authUC.Register(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, req.Email, resp.Email)
		assert.Equal(t, req.Name, resp.Name)
	})

	t.Run("EmailAlreadyExists", func(t *testing.T) {
		req := dto.RegisterRequest{
			Name:     "Test User",
			Email:    "existing@example.com",
			Password: "password123",
		}

		existingUser := &entity.User{
			Email: "existing@example.com",
		}

		mockUserRepo.EXPECT().FindByEmail(ctx, req.Email).Return(existingUser, nil)

		resp, err := authUC.Register(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "email already registered")
	})
}

func TestAuthUseCase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	authUC := auth.NewAuthUseCase(mockUserRepo, mockSessionRepo, nil, nil)

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		now := time.Now()

		user := &entity.User{
			ID:              uuid.New(),
			Email:           "test@example.com",
			Name:            "Test User",
			PasswordHash:    string(hashedPassword),
			EmailVerifiedAt: &now,
		}

		req := dto.LoginRequest{
			Email:    user.Email,
			Password: password,
		}

		mockUserRepo.EXPECT().FindByEmail(ctx, req.Email).Return(user, nil)
		mockSessionRepo.EXPECT().CreateSession(ctx, gomock.Any(), user.ID.String(), 24*time.Hour).Return(nil)

		resp, err := authUC.Login(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.SessionID)
		assert.Equal(t, user.ID, resp.UserID)
	})

	t.Run("EmailNotVerified", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entity.User{
			ID:              uuid.New(),
			Email:           "test@example.com",
			PasswordHash:    string(hashedPassword),
			EmailVerifiedAt: nil,
		}

		req := dto.LoginRequest{
			Email:    user.Email,
			Password: password,
		}

		mockUserRepo.EXPECT().FindByEmail(ctx, req.Email).Return(user, nil)

		resp, err := authUC.Login(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "email not verified")
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		req := dto.LoginRequest{
			Email:    "wrong@example.com",
			Password: "password123",
		}

		mockUserRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)

		resp, err := authUC.Login(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid credentials")
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entity.User{
			ID:           uuid.New(),
			Email:        "test@example.com",
			PasswordHash: string(hashedPassword),
		}

		req := dto.LoginRequest{
			Email:    user.Email,
			Password: "wrongpassword",
		}

		mockUserRepo.EXPECT().FindByEmail(ctx, req.Email).Return(user, nil)

		resp, err := authUC.Login(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid credentials")
	})
}

func TestAuthUseCase_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	authUC := auth.NewAuthUseCase(mockUserRepo, mockSessionRepo, nil, nil)

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		sessionID := "some-session-id"
		mockSessionRepo.EXPECT().DeleteSession(ctx, sessionID).Return(nil)

		err := authUC.Logout(ctx, sessionID)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		sessionID := "some-session-id"
		mockSessionRepo.EXPECT().DeleteSession(ctx, sessionID).Return(errors.New("redis error"))

		err := authUC.Logout(ctx, sessionID)
		assert.Error(t, err)
	})
}

func TestAuthUseCase_LoginWithSocial(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockSocialRepo := mocks.NewMockSocialAccountRepository(ctrl)
	mockSocialAuthService := serviceMocks.NewMockSocialAuthService(ctrl)

	authUC := auth.NewAuthUseCase(mockUserRepo, mockSessionRepo, mockSocialRepo, mockSocialAuthService)
	ctx := context.Background()

	t.Run("SocialAccountExists_Login", func(t *testing.T) {
		req := dto.LoginWithSocialRequest{
			Provider: "google",
			Code:     "auth-code",
		}

		socialInfo := &service.SocialUserInfo{
			ProviderID: "google-123",
			Email:      "test@example.com",
			Name:       "Test User",
		}

		userID := uuid.New()
		socialAccount := &entity.SocialAccount{
			UserID:     userID,
			Provider:   req.Provider,
			ProviderID: socialInfo.ProviderID,
		}

		mockSocialAuthService.EXPECT().GetUserInfo(ctx, req.Provider, req.Code).Return(socialInfo, nil)
		mockSocialRepo.EXPECT().FindByProvider(ctx, req.Provider, socialInfo.ProviderID).Return(socialAccount, nil)
		mockSessionRepo.EXPECT().CreateSession(ctx, gomock.Any(), userID.String(), 24*time.Hour).Return(nil)

		// FindByID is used to populate response
		mockUserRepo.EXPECT().FindByID(ctx, userID).Return(&entity.User{
			ID:    userID,
			Email: socialInfo.Email,
			Name:  socialInfo.Name,
		}, nil)

		resp, err := authUC.LoginWithSocial(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, userID, resp.UserID)
	})

	t.Run("NewUser_RegisterAndLink", func(t *testing.T) {
		req := dto.LoginWithSocialRequest{
			Provider: "google",
			Code:     "auth-code-new",
		}

		socialInfo := &service.SocialUserInfo{
			ProviderID: "google-456",
			Email:      "new@example.com",
			Name:       "New User",
		}

		mockSocialAuthService.EXPECT().GetUserInfo(ctx, req.Provider, req.Code).Return(socialInfo, nil)
		mockSocialRepo.EXPECT().FindByProvider(ctx, req.Provider, socialInfo.ProviderID).Return(nil, nil) // Not found
		mockUserRepo.EXPECT().FindByEmail(ctx, socialInfo.Email).Return(nil, nil)                         // User not found

		// Expect User Creation
		mockUserRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, u *entity.User) error {
			assert.Equal(t, socialInfo.Email, u.Email)
			assert.Equal(t, socialInfo.Name, u.Name)
			return nil
		})

		// Expect Social Account Creation
		mockSocialRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, sa *entity.SocialAccount) error {
			assert.Equal(t, req.Provider, sa.Provider)
			assert.Equal(t, socialInfo.ProviderID, sa.ProviderID)
			return nil
		})

		mockSessionRepo.EXPECT().CreateSession(ctx, gomock.Any(), gomock.Any(), 24*time.Hour).Return(nil)

		resp, err := authUC.LoginWithSocial(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, socialInfo.Email, resp.Email)
	})
}

func TestAuthUseCase_VerifyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	authUC := auth.NewAuthUseCase(mockUserRepo, mockSessionRepo, nil, nil)

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		token := "verification-token"
		userID := uuid.New()
		user := &entity.User{
			ID:              userID,
			Email:           "test@example.com",
			EmailVerifiedAt: nil,
		}

		// 1. Get UserID from token (session)
		mockSessionRepo.EXPECT().GetUserID(ctx, token).Return(userID.String(), nil)

		// 2. Find User
		mockUserRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)

		// 3. Update User
		mockUserRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, u *entity.User) error {
			assert.Equal(t, userID, u.ID)
			assert.NotNil(t, u.EmailVerifiedAt)
			return nil
		})

		// 4. Delete Token (Session)
		mockSessionRepo.EXPECT().DeleteSession(ctx, token).Return(nil)

		err := authUC.VerifyEmail(ctx, token)
		assert.NoError(t, err)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		token := "invalid-token"
		mockSessionRepo.EXPECT().GetUserID(ctx, token).Return("", errors.New("not found"))

		err := authUC.VerifyEmail(ctx, token)
		assert.Error(t, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		token := "valid-token"
		userID := uuid.New()

		mockSessionRepo.EXPECT().GetUserID(ctx, token).Return(userID.String(), nil)
		mockUserRepo.EXPECT().FindByID(ctx, userID).Return(nil, nil)

		err := authUC.VerifyEmail(ctx, token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})
}
