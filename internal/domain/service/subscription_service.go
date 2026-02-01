package service

import (
	"context"
	"errors"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/google/uuid"
)

var (
	ErrCannotSubscribeToSelf = errors.New("cannot subscribe to yourself")
	ErrAlreadySubscribed     = errors.New("already subscribed")
	ErrSubscriptionNotFound  = errors.New("subscription not found")
)

type SubscriptionService interface {
	Subscribe(ctx context.Context, subscriberID, authorID uuid.UUID) (*entity.Subscription, error)
	Unsubscribe(ctx context.Context, subscriberID, authorID uuid.UUID) error
	IsSubscribed(ctx context.Context, subscriberID, authorID uuid.UUID) (bool, error)
	GetSubscriptions(ctx context.Context, subscriberID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[entity.Subscription], error)
	GetSubscribers(ctx context.Context, authorID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[entity.Subscription], error)
	CountSubscribers(ctx context.Context, authorID uuid.UUID) (int64, error)
	CountSubscriptions(ctx context.Context, subscriberID uuid.UUID) (int64, error)
}

type subscriptionService struct {
	subscriptionRepo repository.SubscriptionRepository
}

func NewSubscriptionService(subscriptionRepo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *subscriptionService) Subscribe(ctx context.Context, subscriberID, authorID uuid.UUID) (*entity.Subscription, error) {
	if subscriberID == authorID {
		return nil, ErrCannotSubscribeToSelf
	}

	exists, _ := s.subscriptionRepo.Exists(ctx, subscriberID, authorID)
	if exists {
		return nil, ErrAlreadySubscribed
	}

	subscription := &entity.Subscription{
		SubscriberID: subscriberID,
		AuthorID:     authorID,
	}

	if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) Unsubscribe(ctx context.Context, subscriberID, authorID uuid.UUID) error {
	exists, _ := s.subscriptionRepo.Exists(ctx, subscriberID, authorID)
	if !exists {
		return ErrSubscriptionNotFound
	}
	return s.subscriptionRepo.Delete(ctx, subscriberID, authorID)
}

func (s *subscriptionService) IsSubscribed(ctx context.Context, subscriberID, authorID uuid.UUID) (bool, error) {
	return s.subscriptionRepo.Exists(ctx, subscriberID, authorID)
}

func (s *subscriptionService) GetSubscriptions(ctx context.Context, subscriberID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[entity.Subscription], error) {
	return s.subscriptionRepo.FindBySubscriber(ctx, subscriberID, repository.Pagination{Page: page, PageSize: pageSize})
}

func (s *subscriptionService) GetSubscribers(ctx context.Context, authorID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[entity.Subscription], error) {
	return s.subscriptionRepo.FindByAuthor(ctx, authorID, repository.Pagination{Page: page, PageSize: pageSize})
}

func (s *subscriptionService) CountSubscribers(ctx context.Context, authorID uuid.UUID) (int64, error) {
	return s.subscriptionRepo.CountSubscribers(ctx, authorID)
}

func (s *subscriptionService) CountSubscriptions(ctx context.Context, subscriberID uuid.UUID) (int64, error) {
	return s.subscriptionRepo.CountBySubscriber(ctx, subscriberID)
}
