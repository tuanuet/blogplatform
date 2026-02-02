package usecase

import (
	"context"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/repository"
)

type AdminUseCase interface {
	GetDashboardStats(ctx context.Context) (*dto.DashboardStatsResponse, error)
}

type adminUseCase struct {
	userRepo    repository.UserRepository
	blogRepo    repository.BlogRepository
	commentRepo repository.CommentRepository
}

func NewAdminUseCase(
	userRepo repository.UserRepository,
	blogRepo repository.BlogRepository,
	commentRepo repository.CommentRepository,
) AdminUseCase {
	return &adminUseCase{
		userRepo:    userRepo,
		blogRepo:    blogRepo,
		commentRepo: commentRepo,
	}
}

func (uc *adminUseCase) GetDashboardStats(ctx context.Context) (*dto.DashboardStatsResponse, error) {
	months := 12

	userCounts, err := uc.userRepo.CountByMonth(ctx, months)
	if err != nil {
		return nil, err
	}
	blogCounts, err := uc.blogRepo.CountByMonth(ctx, months)
	if err != nil {
		return nil, err
	}
	commentCounts, err := uc.commentRepo.CountByMonth(ctx, months)
	if err != nil {
		return nil, err
	}

	// Merge logic
	statsMap := make(map[string]*dto.MonthlyStat)

	// Initialize map with last 12 months keys
	now := time.Now()
	for i := 0; i < months; i++ {
		t := now.AddDate(0, -i, 0)
		monthStr := t.Format("2006-01")
		statsMap[monthStr] = &dto.MonthlyStat{Month: monthStr}
	}

	for _, c := range userCounts {
		if val, ok := statsMap[c.Month]; ok {
			val.NewUsers = c.Count
		}
	}
	for _, c := range blogCounts {
		if val, ok := statsMap[c.Month]; ok {
			val.NewBlogs = c.Count
		}
	}
	for _, c := range commentCounts {
		if val, ok := statsMap[c.Month]; ok {
			val.NewComments = c.Count
		}
	}

	// Convert map to sorted slice (Newest first? Or Oldest first?)
	// User didn't specify, but charts usually go Left->Right (Oldest->Newest).
	// But API response 'stats' usually mimics the order or caller handles it.
	// Gatekeeper Example:
	// "stats": [ { "month": "2023-10" ... } ]
	// Let's return Oldest to Newest (Chronological).

	stats := make([]dto.MonthlyStat, 0, months)
	for i := months - 1; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		monthStr := t.Format("2006-01")
		stats = append(stats, *statsMap[monthStr])
	}

	return &dto.DashboardStatsResponse{Stats: stats}, nil
}
