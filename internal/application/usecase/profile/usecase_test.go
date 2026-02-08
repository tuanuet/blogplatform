package profile_test

import (
	"context"
	"testing"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/profile"
	"github.com/aiagent/internal/domain/entity"
	serviceMocks "github.com/aiagent/internal/domain/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateProfile_WithGender(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSvc := serviceMocks.NewMockUserService(ctrl)
	uc := profile.NewProfileUseCase(mockUserSvc)
	userID := uuid.New()

	gender := "male"
	req := dto.UpdateProfileRequest{
		Gender: &gender,
	}

	mockUserSvc.EXPECT().UpdateUser(gomock.Any(), userID, gomock.Any()).Return(nil)

	user := &entity.User{
		ID:     userID,
		Name:   "Test User",
		Email:  "test@example.com",
		Gender: &gender,
	}
	mockUserSvc.EXPECT().GetUser(gomock.Any(), userID).Return(user, nil)

	// Act
	resp, err := uc.UpdateProfile(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, gender, resp.Gender)
}

func TestUpdateProfile_WithBirthday(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSvc := serviceMocks.NewMockUserService(ctrl)
	uc := profile.NewProfileUseCase(mockUserSvc)
	userID := uuid.New()

	birthdayStr := "1990-01-01"
	req := dto.UpdateProfileRequest{
		Birthday: &birthdayStr,
	}

	mockUserSvc.EXPECT().UpdateUser(gomock.Any(), userID, gomock.Any()).Return(nil)

	parsedBirthday, _ := time.Parse("2006-01-02", birthdayStr)
	user := &entity.User{
		ID:       userID,
		Name:     "Test User",
		Email:    "test@example.com",
		Birthday: &parsedBirthday,
	}
	mockUserSvc.EXPECT().GetUser(gomock.Any(), userID).Return(user, nil)

	// Act
	resp, err := uc.UpdateProfile(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, birthdayStr, resp.Birthday)
}

func TestUpdateProfile_WithDescription(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSvc := serviceMocks.NewMockUserService(ctrl)
	uc := profile.NewProfileUseCase(mockUserSvc)
	userID := uuid.New()

	description := "This is a long description about the user."
	req := dto.UpdateProfileRequest{
		Description: &description,
	}

	mockUserSvc.EXPECT().UpdateUser(gomock.Any(), userID, gomock.Any()).Return(nil)

	user := &entity.User{
		ID:          userID,
		Name:        "Test User",
		Email:       "test@example.com",
		Description: &description,
	}
	mockUserSvc.EXPECT().GetUser(gomock.Any(), userID).Return(user, nil)

	// Act
	resp, err := uc.UpdateProfile(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, description, resp.Description)
}

func TestGetProfile_WithDescription(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSvc := serviceMocks.NewMockUserService(ctrl)
	uc := profile.NewProfileUseCase(mockUserSvc)
	userID := uuid.New()

	description := "User description"
	user := &entity.User{
		ID:          userID,
		Name:        "Test User",
		Email:       "test@example.com",
		Description: &description,
		CreatedAt:   time.Now(),
	}
	mockUserSvc.EXPECT().GetUser(gomock.Any(), userID).Return(user, nil)

	// Act
	resp, err := uc.GetProfile(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, description, resp.Description)
}

func TestGetPublicProfile_WithDescription(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSvc := serviceMocks.NewMockUserService(ctrl)
	uc := profile.NewProfileUseCase(mockUserSvc)
	userID := uuid.New()

	description := "Public user description"
	user := &entity.User{
		ID:          userID,
		Name:        "Test User",
		Description: &description,
	}
	mockUserSvc.EXPECT().GetUser(gomock.Any(), userID).Return(user, nil)

	// Act
	resp, err := uc.GetPublicProfile(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, description, resp.Description)
}

func TestUpdateProfile_WithFacebookURL(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSvc := serviceMocks.NewMockUserService(ctrl)
	uc := profile.NewProfileUseCase(mockUserSvc)
	userID := uuid.New()

	facebookURL := "https://facebook.com/testuser"
	req := dto.UpdateProfileRequest{
		FacebookURL: &facebookURL,
	}

	mockUserSvc.EXPECT().UpdateUser(gomock.Any(), userID, gomock.Any()).Return(nil)

	user := &entity.User{
		ID:          userID,
		Name:        "Test User",
		Email:       "test@example.com",
		FacebookURL: &facebookURL,
	}
	mockUserSvc.EXPECT().GetUser(gomock.Any(), userID).Return(user, nil)

	// Act
	resp, err := uc.UpdateProfile(context.Background(), userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, facebookURL, resp.FacebookURL)
}

func TestGetProfile_WithFacebookURL(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSvc := serviceMocks.NewMockUserService(ctrl)
	uc := profile.NewProfileUseCase(mockUserSvc)
	userID := uuid.New()

	facebookURL := "https://facebook.com/testuser"
	user := &entity.User{
		ID:          userID,
		Name:        "Test User",
		Email:       "test@example.com",
		FacebookURL: &facebookURL,
		CreatedAt:   time.Now(),
	}
	mockUserSvc.EXPECT().GetUser(gomock.Any(), userID).Return(user, nil)

	// Act
	resp, err := uc.GetProfile(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, facebookURL, resp.FacebookURL)
}

func TestGetPublicProfile_WithFacebookURL(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSvc := serviceMocks.NewMockUserService(ctrl)
	uc := profile.NewProfileUseCase(mockUserSvc)
	userID := uuid.New()

	facebookURL := "https://facebook.com/testuser"
	user := &entity.User{
		ID:          userID,
		Name:        "Test User",
		FacebookURL: &facebookURL,
	}
	mockUserSvc.EXPECT().GetUser(gomock.Any(), userID).Return(user, nil)

	// Act
	resp, err := uc.GetPublicProfile(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, facebookURL, resp.FacebookURL)
}
