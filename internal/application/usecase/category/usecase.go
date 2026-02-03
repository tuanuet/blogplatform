package category

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
	ErrCategoryNotFound   = domainService.ErrCategoryNotFound
	ErrCategorySlugExists = domainService.ErrCategorySlugExists
)

type CategoryUseCase interface {
	Create(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error)
	List(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[dto.CategoryResponse], error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type categoryUseCase struct {
	categorySvc domainService.CategoryService
}

func NewCategoryUseCase(categorySvc domainService.CategoryService) CategoryUseCase {
	return &categoryUseCase{
		categorySvc: categorySvc,
	}
}

func (uc *categoryUseCase) Create(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := &entity.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := uc.categorySvc.Create(ctx, category); err != nil {
		return nil, err
	}

	return uc.toCategoryResponse(category), nil
}

func (uc *categoryUseCase) GetByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error) {
	category, err := uc.categorySvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.toCategoryResponse(category), nil
}

func (uc *categoryUseCase) List(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[dto.CategoryResponse], error) {
	result, err := uc.categorySvc.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.CategoryResponse, len(result.Data))
	for i, cat := range result.Data {
		responses[i] = *uc.toCategoryResponse(&cat)
	}

	return &repository.PaginatedResult[dto.CategoryResponse]{
		Data:       responses,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (uc *categoryUseCase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := uc.categorySvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}
	if req.Description != nil {
		category.Description = req.Description
	}

	if err := uc.categorySvc.Update(ctx, category); err != nil {
		return nil, err
	}

	return uc.toCategoryResponse(category), nil
}

func (uc *categoryUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.categorySvc.Delete(ctx, id)
}

func (uc *categoryUseCase) toCategoryResponse(category *entity.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}
