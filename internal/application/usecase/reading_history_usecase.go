package usecase

import (
	"context"
	"time"

	"github.com/aiagent/boilerplate/internal/application/dto"
	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/google/uuid"
)

type ReadingHistoryUseCase interface {
	MarkAsRead(ctx context.Context, userID, blogID uuid.UUID) error
	GetHistory(ctx context.Context, userID uuid.UUID, limit int) (*dto.ReadingHistoryListResponse, error)
}

type readingHistoryUseCase struct {
	repo repository.ReadingHistoryRepository
}

func NewReadingHistoryUseCase(repo repository.ReadingHistoryRepository) ReadingHistoryUseCase {
	return &readingHistoryUseCase{repo: repo}
}

func (uc *readingHistoryUseCase) MarkAsRead(ctx context.Context, userID, blogID uuid.UUID) error {
	history := &entity.UserReadingHistory{
		UserID:     userID,
		BlogID:     blogID,
		LastReadAt: time.Now(),
	}
	return uc.repo.Upsert(ctx, history)
}

func (uc *readingHistoryUseCase) GetHistory(ctx context.Context, userID uuid.UUID, limit int) (*dto.ReadingHistoryListResponse, error) {
	histories, err := uc.repo.GetRecentByUserID(ctx, userID, limit)
	if err != nil {
		return nil, err
	}

	response := &dto.ReadingHistoryListResponse{
		History: make([]dto.ReadingHistoryResponse, len(histories)),
	}

	for i, h := range histories {
		var blogResp *dto.BlogListResponse
		if h.Blog != nil {
			blogResp = &dto.BlogListResponse{
				ID:            h.Blog.ID,
				AuthorID:      h.Blog.AuthorID,
				Title:         h.Blog.Title,
				Slug:          h.Blog.Slug,
				Excerpt:       h.Blog.Excerpt,
				ThumbnailURL:  h.Blog.ThumbnailURL,
				Status:        h.Blog.Status,
				Visibility:    h.Blog.Visibility,
				PublishedAt:   h.Blog.PublishedAt,
				UpvoteCount:   h.Blog.UpvoteCount,
				DownvoteCount: h.Blog.DownvoteCount,
				CreatedAt:     h.Blog.CreatedAt,
			}

			if h.Blog.Author != nil {
				blogResp.Author = &dto.UserBriefResponse{
					ID:    h.Blog.Author.ID,
					Name:  h.Blog.Author.Name,
					Email: h.Blog.Author.Email,
				}
			}

			if h.Blog.Category != nil {
				blogResp.Category = &dto.CategoryResponse{
					ID:        h.Blog.Category.ID,
					Name:      h.Blog.Category.Name,
					Slug:      h.Blog.Category.Slug,
					CreatedAt: h.Blog.Category.CreatedAt,
					UpdatedAt: h.Blog.Category.UpdatedAt,
				}
				if h.Blog.Category.Description != nil {
					blogResp.Category.Description = h.Blog.Category.Description
				}
			}

			if len(h.Blog.Tags) > 0 {
				blogResp.Tags = make([]dto.TagResponse, len(h.Blog.Tags))
				for j, tag := range h.Blog.Tags {
					blogResp.Tags[j] = dto.TagResponse{
						ID:        tag.ID,
						Name:      tag.Name,
						Slug:      tag.Slug,
						CreatedAt: tag.CreatedAt,
					}
				}
			}
		}

		response.History[i] = dto.ReadingHistoryResponse{
			BlogID:     h.BlogID,
			LastReadAt: h.LastReadAt,
			Blog:       blogResp,
		}
	}

	return response, nil
}
