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
	ErrTagNotFound   = errors.New("tag not found")
	ErrTagSlugExists = errors.New("tag slug already exists")
)

type TagService interface {
	Create(ctx context.Context, tag *entity.Tag) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Tag, error)
	List(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[entity.Tag], error)
	Update(ctx context.Context, tag *entity.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type tagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{
		tagRepo: tagRepo,
	}
}

func (s *tagService) Create(ctx context.Context, tag *entity.Tag) error {
	existing, _ := s.tagRepo.FindBySlug(ctx, tag.Slug)
	if existing != nil {
		return ErrTagSlugExists
	}
	return s.tagRepo.Create(ctx, tag)
}

func (s *tagService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error) {
	tag, err := s.tagRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, ErrTagNotFound
	}
	return tag, nil
}

func (s *tagService) GetBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	tag, err := s.tagRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, ErrTagNotFound
	}
	return tag, nil
}

func (s *tagService) List(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[entity.Tag], error) {
	return s.tagRepo.FindAll(ctx, repository.Pagination{Page: page, PageSize: pageSize})
}

func (s *tagService) Update(ctx context.Context, tag *entity.Tag) error {
	existing, _ := s.tagRepo.FindBySlug(ctx, tag.Slug)
	if existing != nil && existing.ID != tag.ID {
		return ErrTagSlugExists
	}
	return s.tagRepo.Update(ctx, tag)
}

func (s *tagService) Delete(ctx context.Context, id uuid.UUID) error {
	existing, err := s.tagRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrTagNotFound
	}
	return s.tagRepo.Delete(ctx, id)
}
