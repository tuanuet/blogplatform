package service_test

import (
	"context"
	"testing"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	repoMocks "github.com/aiagent/internal/domain/repository/mocks"
	"github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSubscriptionService_GetTopSubscribedAuthors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockSubscriptionRepository(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	authorID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	expected := &repository.PaginatedResult[entity.TopSubscribedAuthor]{
		Data: []entity.TopSubscribedAuthor{{
			AuthorID:        authorID,
			Username:        "alice",
			DisplayName:     "Alice",
			AvatarURL:       "https://cdn.example.com/alice.png",
			SubscriberCount: 10,
		}},
		Total:      1,
		Page:       1,
		PageSize:   20,
		TotalPages: 1,
	}

	mockRepo.EXPECT().
		FindTopSubscribedAuthors(gomock.Any(), repository.Pagination{Page: 1, PageSize: 20}).
		Return(expected, nil)

	result, err := svc.GetTopSubscribedAuthors(context.Background(), 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSubscriptionService_GetTopSubscribedAuthors_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMocks.NewMockSubscriptionRepository(ctrl)
	svc := service.NewSubscriptionService(mockRepo)

	mockRepo.EXPECT().
		FindTopSubscribedAuthors(gomock.Any(), repository.Pagination{Page: 2, PageSize: 5}).
		Return(nil, assert.AnError)

	result, err := svc.GetTopSubscribedAuthors(context.Background(), 2, 5)

	assert.ErrorIs(t, err, assert.AnError)
	assert.Nil(t, result)
}
