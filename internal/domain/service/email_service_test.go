package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/aiagent/internal/domain/entity"
	repoMocks "github.com/aiagent/internal/domain/repository/mocks"
	adapterMocks "github.com/aiagent/internal/infrastructure/adapter/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type syncTaskRunner struct{}

func (r *syncTaskRunner) Submit(task func(ctx context.Context)) {
	task(context.Background())
}

func TestEmailService_SendWelcomeEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repoMocks.NewMockUserRepository(ctrl)
	mockProvider := adapterMocks.NewMockEmailProvider(ctrl)
	taskRunner := &syncTaskRunner{}

	// We need to find the templates directory.
	// Since tests run in the package directory, we need to go up.
	wd, _ := os.Getwd()
	templateDir := filepath.Join(wd, "../../infrastructure/email/templates")

	service := NewEmailServiceImpl(mockUserRepo, mockProvider, taskRunner, templateDir)

	userID := uuid.New()
	email := "test@example.com"
	name := "Test User"

	mockProvider.EXPECT().
		Send(gomock.Any(), []string{email}, "Welcome to AI Agent!", gomock.Any(), gomock.Any()).
		Return(nil)

	err := service.SendWelcomeEmail(context.Background(), userID, email, name)
	assert.NoError(t, err)
}

func TestEmailService_SendNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repoMocks.NewMockUserRepository(ctrl)
	mockProvider := adapterMocks.NewMockEmailProvider(ctrl)
	taskRunner := &syncTaskRunner{}

	wd, _ := os.Getwd()
	templateDir := filepath.Join(wd, "../../infrastructure/email/templates")

	service := NewEmailServiceImpl(mockUserRepo, mockProvider, taskRunner, templateDir)

	userID := uuid.New()
	user := &entity.User{
		ID:    userID,
		Email: "test@example.com",
	}

	mockUserRepo.EXPECT().
		FindByID(gomock.Any(), userID).
		Return(user, nil)

	mockProvider.EXPECT().
		Send(gomock.Any(), []string{user.Email}, gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	data := map[string]interface{}{
		"message":    "You have a new follower",
		"action_url": "https://aiagent.com/followers",
	}

	err := service.SendNotification(context.Background(), userID, entity.NotificationTypeNewFollower, data)
	assert.NoError(t, err)
}

func TestEmailService_SendVerificationEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repoMocks.NewMockUserRepository(ctrl)
	mockProvider := adapterMocks.NewMockEmailProvider(ctrl)
	taskRunner := &syncTaskRunner{}

	wd, _ := os.Getwd()
	templateDir := filepath.Join(wd, "../../infrastructure/email/templates")

	service := NewEmailServiceImpl(mockUserRepo, mockProvider, taskRunner, templateDir)

	userID := uuid.New()
	email := "test@example.com"
	token := "abc-123"

	mockProvider.EXPECT().
		Send(gomock.Any(), []string{email}, "Verify your email address", gomock.Any(), gomock.Any()).
		Return(nil)

	err := service.SendVerificationEmail(context.Background(), userID, email, token)
	assert.NoError(t, err)
}
