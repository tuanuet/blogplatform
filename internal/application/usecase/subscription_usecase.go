package usecase

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	domainService "github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
)

var (
	ErrCannotSubscribeToSelf = domainService.ErrCannotSubscribeToSelf
	ErrAlreadySubscribed     = domainService.ErrAlreadySubscribed
	ErrSubscriptionNotFound  = domainService.ErrSubscriptionNotFound
)

type SubscriptionUseCase interface {
	Subscribe(ctx context.Context, subscriberID, authorID uuid.UUID) (*dto.SubscriptionResponse, error)
	Unsubscribe(ctx context.Context, subscriberID, authorID uuid.UUID) error
	IsSubscribed(ctx context.Context, subscriberID, authorID uuid.UUID) (bool, error)
	GetSubscriptions(ctx context.Context, subscriberID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[dto.SubscriptionResponse], error)
	GetSubscribers(ctx context.Context, authorID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[dto.SubscriptionResponse], error)
	CountSubscribers(ctx context.Context, authorID uuid.UUID) (*dto.SubscriptionCountResponse, error)
	GetSubscriptionCounts(ctx context.Context, userID uuid.UUID) (*dto.SubscriptionCountResponse, error)
}

type subscriptionUseCase struct {
	subscriptionSvc domainService.SubscriptionService
}

func NewSubscriptionUseCase(subscriptionSvc domainService.SubscriptionService) SubscriptionUseCase {
	return &subscriptionUseCase{
		subscriptionSvc: subscriptionSvc,
	}
}

func (uc *subscriptionUseCase) Subscribe(ctx context.Context, subscriberID, authorID uuid.UUID) (*dto.SubscriptionResponse, error) {
	sub, err := uc.subscriptionSvc.Subscribe(ctx, subscriberID, authorID)
	if err != nil {
		return nil, err
	}
	return uc.toSubscriptionResponse(sub), nil
}

func (uc *subscriptionUseCase) Unsubscribe(ctx context.Context, subscriberID, authorID uuid.UUID) error {
	return uc.subscriptionSvc.Unsubscribe(ctx, subscriberID, authorID)
}

func (uc *subscriptionUseCase) IsSubscribed(ctx context.Context, subscriberID, authorID uuid.UUID) (bool, error) {
	return uc.subscriptionSvc.IsSubscribed(ctx, subscriberID, authorID)
}

func (uc *subscriptionUseCase) GetSubscriptions(ctx context.Context, subscriberID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[dto.SubscriptionResponse], error) {
	result, err := uc.subscriptionSvc.GetSubscriptions(ctx, subscriberID, page, pageSize)
	if err != nil {
		return nil, err
	}

	subs := make([]dto.SubscriptionResponse, 0, len(result.Data))
	for _, sub := range result.Data {
		subs = append(subs, *uc.toSubscriptionResponse(&sub))
	}

	return &repository.PaginatedResult[dto.SubscriptionResponse]{
		Data:       subs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (uc *subscriptionUseCase) GetSubscribers(ctx context.Context, authorID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[dto.SubscriptionResponse], error) {
	result, err := uc.subscriptionSvc.GetSubscribers(ctx, authorID, page, pageSize)
	if err != nil {
		return nil, err
	}

	subs := make([]dto.SubscriptionResponse, 0, len(result.Data))
	for _, sub := range result.Data {
		subs = append(subs, *uc.toSubscriptionResponse(&sub))
	}

	return &repository.PaginatedResult[dto.SubscriptionResponse]{
		Data:       subs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (uc *subscriptionUseCase) CountSubscribers(ctx context.Context, authorID uuid.UUID) (*dto.SubscriptionCountResponse, error) {
	count, err := uc.subscriptionSvc.CountSubscribers(ctx, authorID)
	if err != nil {
		return nil, err
	}
	return &dto.SubscriptionCountResponse{
		AuthorID:        authorID,
		SubscriberCount: count,
	}, nil
}

func (uc *subscriptionUseCase) GetSubscriptionCounts(ctx context.Context, userID uuid.UUID) (*dto.SubscriptionCountResponse, error) {
	subscriberCount, err := uc.subscriptionSvc.CountSubscribers(ctx, userID)
	if err != nil {
		return nil, err
	}

	subscriptionCount, err := uc.subscriptionSvc.CountSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.SubscriptionCountResponse{
		AuthorID:          userID,
		SubscriberCount:   subscriberCount,
		SubscriptionCount: subscriptionCount,
	}, nil
}

func (uc *subscriptionUseCase) toSubscriptionResponse(sub *entity.Subscription) *dto.SubscriptionResponse {
	resp := &dto.SubscriptionResponse{
		ID:           sub.ID,
		SubscriberID: sub.SubscriberID,
		AuthorID:     sub.AuthorID,
		CreatedAt:    sub.CreatedAt,
	}

	if sub.Subscriber != nil {
		resp.Subscriber = &dto.UserBriefResponse{
			ID:    sub.Subscriber.ID,
			Name:  sub.Subscriber.Name,
			Email: sub.Subscriber.Email,
		}
	}

	if sub.Author != nil {
		resp.Author = &dto.UserBriefResponse{
			ID:    sub.Author.ID,
			Name:  sub.Author.Name,
			Email: sub.Author.Email,
		}
	}

	return resp
}
