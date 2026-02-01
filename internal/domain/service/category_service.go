package service

import (
	"context"
	"errors"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/google/uuid"
)

var (
	ErrCategoryNotFound   = errors.New("category not found")
	ErrCategorySlugExists = errors.New("category slug already exists")
)

type CategoryService interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
	List(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[entity.Category], error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) Create(ctx context.Context, category *entity.Category) error {
	existing, _ := s.categoryRepo.FindBySlug(ctx, category.Slug)
	if existing != nil {
		return ErrCategorySlugExists
	}
	return s.categoryRepo.Create(ctx, category)
}

func (s *categoryService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *categoryService) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	category, err := s.categoryRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *categoryService) List(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[entity.Category], error) {
	return s.categoryRepo.FindAll(ctx, repository.Pagination{Page: page, PageSize: pageSize})
}

func (s *categoryService) Update(ctx context.Context, category *entity.Category) error {
	existing, _ := s.categoryRepo.FindBySlug(ctx, category.Slug)
	if existing != nil && existing.ID != category.ID {
		return ErrCategorySlugExists
	}
	return s.categoryRepo.Update(ctx, category)
}

func (s *categoryService) Delete(ctx context.Context, id uuid.UUID) error {
	// Verify existence
	existing, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrCategoryNotFound
	}
	return s.categoryRepo.Delete(ctx, id)
}
