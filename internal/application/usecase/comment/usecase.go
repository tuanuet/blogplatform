package comment

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
	ErrCommentNotFound     = domainService.ErrCommentNotFound
	ErrCommentAccessDenied = domainService.ErrCommentAccessDenied
)

type CommentUseCase interface {
	Create(ctx context.Context, userID, blogID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.CommentResponse, error)
	GetByBlogID(ctx context.Context, blogID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[dto.CommentResponse], error)
	Update(ctx context.Context, id, userID uuid.UUID, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type commentUseCase struct {
	commentSvc domainService.CommentService
}

func NewCommentUseCase(commentSvc domainService.CommentService) CommentUseCase {
	return &commentUseCase{
		commentSvc: commentSvc,
	}
}

func (uc *commentUseCase) Create(ctx context.Context, userID, blogID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	comment := &entity.Comment{
		BlogID:  blogID,
		UserID:  userID,
		Content: req.Content,
	}
	if req.ParentID != nil {
		id, err := uuid.Parse(*req.ParentID)
		if err == nil {
			comment.ParentID = &id
		}
	}

	if err := uc.commentSvc.Create(ctx, comment); err != nil {
		return nil, err
	}

	return uc.toCommentResponse(comment), nil
}

func (uc *commentUseCase) GetByID(ctx context.Context, id uuid.UUID) (*dto.CommentResponse, error) {
	comment, err := uc.commentSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toCommentResponse(comment), nil
}

func (uc *commentUseCase) GetByBlogID(ctx context.Context, blogID uuid.UUID, page, pageSize int) (*repository.PaginatedResult[dto.CommentResponse], error) {
	result, err := uc.commentSvc.GetByBlogID(ctx, blogID, page, pageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CommentResponse, len(result.Data))
	for i, c := range result.Data {
		responses[i] = *uc.toCommentResponse(&c)
	}

	return &repository.PaginatedResult[dto.CommentResponse]{
		Data:       responses,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (uc *commentUseCase) Update(ctx context.Context, id, userID uuid.UUID, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error) {
	comment := &entity.Comment{
		ID:      id,
		Content: req.Content,
	}
	if err := uc.commentSvc.Update(ctx, comment, userID); err != nil {
		return nil, err
	}

	// Fetch updated
	updated, err := uc.commentSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toCommentResponse(updated), nil
}

func (uc *commentUseCase) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return uc.commentSvc.Delete(ctx, id, userID)
}

func (uc *commentUseCase) toCommentResponse(comment *entity.Comment) *dto.CommentResponse {
	if comment == nil {
		return nil
	}

	resp := &dto.CommentResponse{
		ID:        comment.ID,
		BlogID:    comment.BlogID,
		UserID:    comment.UserID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		Replies:   make([]dto.CommentResponse, 0),
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	if comment.User != nil {
		resp.User = &dto.UserBriefResponse{
			ID:    comment.User.ID,
			Name:  comment.User.Name,
			Email: comment.User.Email,
		}
	}

	// Add replies if loaded
	for _, reply := range comment.Replies {
		resp.Replies = append(resp.Replies, *uc.toCommentResponse(&reply))
	}

	return resp
}
