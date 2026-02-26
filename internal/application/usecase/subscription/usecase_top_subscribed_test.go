package subscription_test

import (
	"context"
	"testing"

	"github.com/aiagent/internal/application/usecase/subscription"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	serviceMocks "github.com/aiagent/internal/domain/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSubscriptionUseCase_GetTopSubscribedAuthors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := serviceMocks.NewMockSubscriptionService(ctrl)
	uc := subscription.NewSubscriptionUseCase(mockSvc)

	firstAuthorID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	secondAuthorID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	mockSvc.EXPECT().
		GetTopSubscribedAuthors(gomock.Any(), 1, 20).
		Return(&repository.PaginatedResult[entity.TopSubscribedAuthor]{
			Data: []entity.TopSubscribedAuthor{
				{
					AuthorID:        firstAuthorID,
					Username:        "alice",
					DisplayName:     "Alice",
					AvatarURL:       "https://cdn.example.com/alice.png",
					SubscriberCount: 12,
				},
				{
					AuthorID:        secondAuthorID,
					Username:        "bob",
					DisplayName:     "Bob",
					AvatarURL:       "",
					SubscriberCount: 7,
				},
			},
			Total:      2,
			Page:       1,
			PageSize:   20,
			TotalPages: 1,
		}, nil)

	result, err := uc.GetTopSubscribedAuthors(context.Background(), 1, 20)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, firstAuthorID, result.Data[0].AuthorID)
	assert.Equal(t, "alice", result.Data[0].Username)
	assert.Equal(t, int64(12), result.Data[0].SubscriberCount)
	assert.Equal(t, secondAuthorID, result.Data[1].AuthorID)
	assert.Equal(t, int64(2), result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 20, result.PageSize)
	assert.Equal(t, 1, result.TotalPages)
}

func TestSubscriptionUseCase_GetTopSubscribedAuthors_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := serviceMocks.NewMockSubscriptionService(ctrl)
	uc := subscription.NewSubscriptionUseCase(mockSvc)

	mockSvc.EXPECT().
		GetTopSubscribedAuthors(gomock.Any(), 2, 5).
		Return(nil, assert.AnError)

	result, err := uc.GetTopSubscribedAuthors(context.Background(), 2, 5)

	assert.ErrorIs(t, err, assert.AnError)
	assert.Nil(t, result)
}
