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
	ErrTagNotFound   = domainService.ErrTagNotFound
	ErrTagSlugExists = domainService.ErrTagSlugExists
)

type TagUseCase interface {
	Create(ctx context.Context, req *dto.CreateTagRequest) (*dto.TagResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.TagResponse, error)
	List(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[dto.TagResponse], error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateTagRequest) (*dto.TagResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type tagUseCase struct {
	tagSvc domainService.TagService
}

func NewTagUseCase(tagSvc domainService.TagService) TagUseCase {
	return &tagUseCase{
		tagSvc: tagSvc,
	}
}

func (uc *tagUseCase) Create(ctx context.Context, req *dto.CreateTagRequest) (*dto.TagResponse, error) {
	tag := &entity.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := uc.tagSvc.Create(ctx, tag); err != nil {
		return nil, err
	}

	return uc.toTagResponse(tag), nil
}

func (uc *tagUseCase) GetByID(ctx context.Context, id uuid.UUID) (*dto.TagResponse, error) {
	tag, err := uc.tagSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toTagResponse(tag), nil
}

func (uc *tagUseCase) List(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[dto.TagResponse], error) {
	result, err := uc.tagSvc.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TagResponse, len(result.Data))
	for i, tag := range result.Data {
		responses[i] = *uc.toTagResponse(&tag)
	}

	return &repository.PaginatedResult[dto.TagResponse]{
		Data:       responses,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (uc *tagUseCase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateTagRequest) (*dto.TagResponse, error) {
	tag, err := uc.tagSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		tag.Name = *req.Name
	}
	if req.Slug != nil {
		tag.Slug = *req.Slug
	}

	if err := uc.tagSvc.Update(ctx, tag); err != nil {
		return nil, err
	}

	return uc.toTagResponse(tag), nil
}

func (uc *tagUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.tagSvc.Delete(ctx, id)
}

func (uc *tagUseCase) toTagResponse(tag *entity.Tag) *dto.TagResponse {
	return &dto.TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
	}
}
