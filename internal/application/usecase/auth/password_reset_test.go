package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aiagent/internal/application/usecase/auth"
	"github.com/aiagent/internal/domain/entity"
	repoMocks "github.com/aiagent/internal/domain/repository/mocks"
	serviceMocks "github.com/aiagent/internal/domain/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUseCase_ForgotPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repoMocks.NewMockUserRepository(ctrl)
	mockSessionRepo := repoMocks.NewMockSessionRepository(ctrl)
	mockEmailService := serviceMocks.NewMockEmailService(ctrl)
	uc := auth.NewAuthUseCase(mockUserRepo, mockSessionRepo, nil, nil, mockEmailService)

	ctx := context.Background()

	t.Run("AlwaysSucceeds_CreatesTokenAndSendsEmail_WhenUserExists", func(t *testing.T) {
		email := "test@example.com"
		userID := uuid.New()
		user := &entity.User{ID: userID, Email: email, IsActive: true}

		var createdToken string

		gomock.InOrder(
			mockUserRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil),
			mockSessionRepo.EXPECT().CreateSession(ctx, gomock.Any(), userID.String(), 30*time.Minute).
				DoAndReturn(func(ctx context.Context, sessionID, uid string, d time.Duration) error {
					createdToken = sessionID
					assert.Contains(t, sessionID, "reset:")
					return nil
				}),
			mockEmailService.EXPECT().SendPasswordResetEmail(ctx, userID, email, gomock.Any()).
				DoAndReturn(func(ctx context.Context, uid uuid.UUID, email, token string) error {
					assert.Equal(t, createdToken, token)
					return nil
				}),
		)

		err := uc.ForgotPassword(ctx, email)
		assert.NoError(t, err)
	})

	t.Run("AlwaysSucceeds_DoesNothing_WhenUserMissing", func(t *testing.T) {
		email := "missing@example.com"

		mockUserRepo.EXPECT().FindByEmail(ctx, email).Return(nil, nil)
		mockSessionRepo.EXPECT().CreateSession(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockEmailService.EXPECT().SendPasswordResetEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		err := uc.ForgotPassword(ctx, email)
		assert.NoError(t, err)
	})

	t.Run("EmailSendFailure_DeletesToken", func(t *testing.T) {
		email := "test@example.com"
		userID := uuid.New()
		user := &entity.User{ID: userID, Email: email, IsActive: true}

		token := "reset:" + uuid.New().String()

		gomock.InOrder(
			mockUserRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil),
			mockSessionRepo.EXPECT().CreateSession(ctx, gomock.Any(), userID.String(), 30*time.Minute).
				Return(nil),
			mockEmailService.EXPECT().SendPasswordResetEmail(ctx, userID, email, gomock.Any()).
				DoAndReturn(func(ctx context.Context, uid uuid.UUID, email, tkn string) error {
					token = tkn
					return errors.New("smtp down")
				}),
			mockSessionRepo.EXPECT().DeleteSession(ctx, gomock.Any()).
				DoAndReturn(func(ctx context.Context, sessionID string) error {
					assert.Equal(t, token, sessionID)
					return nil
				}),
		)

		err := uc.ForgotPassword(ctx, email)
		assert.NoError(t, err)
	})
}

func TestAuthUseCase_ValidateResetPasswordToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repoMocks.NewMockUserRepository(ctrl)
	mockSessionRepo := repoMocks.NewMockSessionRepository(ctrl)
	mockEmailService := serviceMocks.NewMockEmailService(ctrl)
	uc := auth.NewAuthUseCase(mockUserRepo, mockSessionRepo, nil, nil, mockEmailService)

	ctx := context.Background()

	tests := []struct {
		name  string
		token string
		setup func()
		want  assert.ErrorAssertionFunc
	}{
		{
			name:  "InvalidPrefix",
			token: "verify:abc",
			setup: func() {
				mockSessionRepo.EXPECT().GetUserID(gomock.Any(), gomock.Any()).Times(0)
			},
			want: assert.Error,
		},
		{
			name:  "MissingToken",
			token: "reset:" + uuid.New().String(),
			setup: func() {
				mockSessionRepo.EXPECT().GetUserID(ctx, gomock.Any()).Return("", errors.New("not found"))
			},
			want: assert.Error,
		},
		{
			name:  "ExistingToken",
			token: "reset:" + uuid.New().String(),
			setup: func() {
				mockSessionRepo.EXPECT().GetUserID(ctx, gomock.Any()).Return(uuid.New().String(), nil)
			},
			want: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := uc.ValidateResetPasswordToken(ctx, tt.token)
			tt.want(t, err)
		})
	}
}

func TestAuthUseCase_ResetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repoMocks.NewMockUserRepository(ctrl)
	mockSessionRepo := repoMocks.NewMockSessionRepository(ctrl)
	mockEmailService := serviceMocks.NewMockEmailService(ctrl)
	uc := auth.NewAuthUseCase(mockUserRepo, mockSessionRepo, nil, nil, mockEmailService)

	ctx := context.Background()

	t.Run("RejectsShortPassword", func(t *testing.T) {
		mockSessionRepo.EXPECT().GetUserID(gomock.Any(), gomock.Any()).Times(0)
		mockUserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
		mockSessionRepo.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Times(0)

		err := uc.ResetPassword(ctx, "reset:"+uuid.New().String(), "short")
		assert.Error(t, err)
	})

	t.Run("UpdatesPasswordAndInvalidatesToken", func(t *testing.T) {
		token := "reset:" + uuid.New().String()
		userID := uuid.New()
		oldHash, _ := bcrypt.GenerateFromPassword([]byte("old-password"), bcrypt.DefaultCost)
		user := &entity.User{ID: userID, Email: "test@example.com", PasswordHash: string(oldHash)}
		oldHashStr := user.PasswordHash
		newPassword := "newpassword123"

		gomock.InOrder(
			mockSessionRepo.EXPECT().GetUserID(ctx, token).Return(userID.String(), nil),
			mockUserRepo.EXPECT().FindByID(ctx, userID).Return(user, nil),
			mockUserRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, u *entity.User) error {
				assert.Equal(t, userID, u.ID)
				assert.NotEqual(t, oldHashStr, u.PasswordHash)
				err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(newPassword))
				assert.NoError(t, err)
				return nil
			}),
			mockSessionRepo.EXPECT().DeleteSession(ctx, token).Return(nil),
		)

		err := uc.ResetPassword(ctx, token, newPassword)
		assert.NoError(t, err)
	})
}
