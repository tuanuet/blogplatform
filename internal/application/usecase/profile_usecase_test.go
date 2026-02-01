package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/aiagent/boilerplate/internal/application/dto"
	"github.com/aiagent/boilerplate/internal/application/usecase"
	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of domainService.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockUserService) UpdateAvatarURL(ctx context.Context, id uuid.UUID, avatarURL string) error {
	args := m.Called(ctx, id, avatarURL)
	return args.Error(0)
}

func TestUpdateProfile_WithGender(t *testing.T) {
	// Arrange
	mockUserSvc := new(MockUserService)
	uc := usecase.NewProfileUseCase(mockUserSvc)
	userID := uuid.New()

	gender := "male"
	req := dto.UpdateProfileRequest{
		Gender: &gender, // This will fail compilation initially
	}

	mockUserSvc.On("UpdateUser", context.Background(), userID, mock.MatchedBy(func(updates map[string]interface{}) bool {
		return updates["gender"] == "male"
	})).Return(nil)

	user := &entity.User{
		ID:     userID,
		Name:   "Test User",
		Email:  "test@example.com",
		Gender: &gender, // This will also fail compilation
	}
	mockUserSvc.On("GetUser", context.Background(), userID).Return(user, nil)

	// Act
	resp, err := uc.UpdateProfile(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, gender, resp.Gender)
	mockUserSvc.AssertExpectations(t)
}

func TestUpdateProfile_WithBirthday(t *testing.T) {
	// Arrange
	mockUserSvc := new(MockUserService)
	uc := usecase.NewProfileUseCase(mockUserSvc)
	userID := uuid.New()

	birthdayStr := "1990-01-01"
	req := dto.UpdateProfileRequest{
		Birthday: &birthdayStr,
	}

	mockUserSvc.On("UpdateUser", context.Background(), userID, mock.MatchedBy(func(updates map[string]interface{}) bool {
		val, ok := updates["birthday"]
		if !ok {
			return false
		}
		tVal, ok := val.(time.Time)
		if !ok {
			return false
		}
		return tVal.Format("2006-01-02") == birthdayStr
	})).Return(nil)

	parsedBirthday, _ := time.Parse("2006-01-02", birthdayStr)
	user := &entity.User{
		ID:       userID,
		Name:     "Test User",
		Email:    "test@example.com",
		Birthday: &parsedBirthday,
	}
	mockUserSvc.On("GetUser", context.Background(), userID).Return(user, nil)

	// Act
	resp, err := uc.UpdateProfile(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, birthdayStr, resp.Birthday)
	mockUserSvc.AssertExpectations(t)
}
