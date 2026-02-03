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
