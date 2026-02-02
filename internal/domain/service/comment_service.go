package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"errors"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

var (
	ErrCommentNotFound     = errors.New("comment not found")
	ErrCommentAccessDenied = errors.New("access denied to comment")
)

type CommentService interface {
	Create(ctx context.Context, comment *entity.Comment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error)
	GetByBlogID(ctx context.Context, blogID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[entity.Comment], error)
	Update(ctx context.Context, comment *entity.Comment, userID uuid.UUID) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type commentService struct {
	commentRepo repository.CommentRepository
}

func NewCommentService(commentRepo repository.CommentRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
	}
}

func (s *commentService) Create(ctx context.Context, comment *entity.Comment) error {
	return s.commentRepo.Create(ctx, comment)
}

func (s *commentService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error) {
	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, ErrCommentNotFound
	}
	return comment, nil
}

func (s *commentService) GetByBlogID(ctx context.Context, blogID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[entity.Comment], error) {
	return s.commentRepo.FindByBlogID(ctx, blogID, repository.Pagination{Page: page, PageSize: pageSize})
}

func (s *commentService) Update(ctx context.Context, comment *entity.Comment, userID uuid.UUID) error {
	// Verify ownership
	existing, err := s.GetByID(ctx, comment.ID)
	if err != nil {
		return err
	}
	if existing.UserID != userID {
		return ErrCommentAccessDenied
	}

	// Update fields
	existing.Content = comment.Content
	return s.commentRepo.Update(ctx, existing)
}

func (s *commentService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if comment == nil {
		return ErrCommentNotFound
	}

	if comment.UserID != userID {
		return ErrCommentAccessDenied
	}

	return s.commentRepo.Delete(ctx, id)
}
