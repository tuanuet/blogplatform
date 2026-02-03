package ranking

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	domainService "github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
)

// RankingUseCase defines the interface for ranking operations
type RankingUseCase interface {
	// GetTrendingUsers gets the trending users by velocity score
	GetTrendingUsers(ctx context.Context, params *dto.RankingFilterParams) (*repository.PaginatedResult[dto.RankedUserResponse], error)

	// GetTopUsers gets top users by total followers
	GetTopUsers(ctx context.Context, params *dto.RankingFilterParams) (*repository.PaginatedResult[dto.RankedUserResponse], error)

	// GetUserRanking gets a specific user's ranking details
	GetUserRanking(ctx context.Context, userID uuid.UUID) (*dto.UserRankingDetailResponse, error)

	// RecalculateAllScores triggers a recalculation of all velocity scores (admin only)
	RecalculateAllScores(ctx context.Context) error
}

type rankingUseCase struct {
	rankingSvc        domainService.RankingService
	velocityScoreRepo repository.UserVelocityScoreRepository
	userRepo          repository.UserRepository
}

// NewRankingUseCase creates a new ranking use case
func NewRankingUseCase(
	rankingSvc domainService.RankingService,
	velocityScoreRepo repository.UserVelocityScoreRepository,
	userRepo repository.UserRepository,
) RankingUseCase {
	return &rankingUseCase{
		rankingSvc:        rankingSvc,
		velocityScoreRepo: velocityScoreRepo,
		userRepo:          userRepo,
	}
}

func (uc *rankingUseCase) GetTrendingUsers(ctx context.Context, params *dto.RankingFilterParams) (*repository.PaginatedResult[dto.RankedUserResponse], error) {
	pagination := repository.Pagination{
		Page:     params.Page,
		PageSize: params.PageSize,
	}

	result, err := uc.velocityScoreRepo.ListRanked(ctx, pagination)
	if err != nil {
		return nil, err
	}

	items := make([]dto.RankedUserResponse, len(result.Data))
	for i, score := range result.Data {
		items[i] = uc.toRankedUserResponse(&score)
	}

	return &repository.PaginatedResult[dto.RankedUserResponse]{
		Data:       items,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (uc *rankingUseCase) GetTopUsers(ctx context.Context, params *dto.RankingFilterParams) (*repository.PaginatedResult[dto.RankedUserResponse], error) {
	// For now, this is similar to GetTrendingUsers but could be extended
	// to use a different sorting mechanism (e.g., by total followers)
	return uc.GetTrendingUsers(ctx, params)
}

func (uc *rankingUseCase) GetUserRanking(ctx context.Context, userID uuid.UUID) (*dto.UserRankingDetailResponse, error) {
	// Ensure user exists
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domainService.ErrUserNotFound
	}

	// Get ranking detail from service
	detail, err := uc.rankingSvc.GetUserRankingDetail(ctx, userID)
	if err != nil {
		return nil, err
	}

	return uc.toUserRankingDetailResponse(detail, user), nil
}

func (uc *rankingUseCase) RecalculateAllScores(ctx context.Context) error {
	// Archive current rankings
	if err := uc.rankingSvc.ArchiveRankings(ctx); err != nil {
		return err
	}

	// Recalculate all scores
	if err := uc.rankingSvc.CalculateAllVelocityScores(ctx); err != nil {
		return err
	}

	// Assign new rank positions
	return uc.rankingSvc.AssignRankPositions(ctx)
}

func (uc *rankingUseCase) toRankedUserResponse(score *entity.UserVelocityScore) dto.RankedUserResponse {
	resp := dto.RankedUserResponse{
		ID:                 score.UserID,
		FollowerCount:      score.FollowerCount,
		FollowerGrowthRate: score.FollowerGrowthRate,
		BlogPostVelocity:   score.BlogPostVelocity,
		CompositeScore:     score.CompositeScore,
		CalculationDate:    score.CalculationDate,
	}

	if score.RankPosition != nil {
		rank := *score.RankPosition
		resp.Rank = &rank
	}

	if score.User != nil {
		resp.Username = score.User.Name
		if score.User.DisplayName != nil {
			resp.DisplayName = *score.User.DisplayName
		}
		if score.User.AvatarURL != nil {
			resp.AvatarURL = *score.User.AvatarURL
		}
	}

	return resp
}

func (uc *rankingUseCase) toUserRankingDetailResponse(detail *domainService.UserRankingDetail, user *entity.User) *dto.UserRankingDetailResponse {
	resp := &dto.UserRankingDetailResponse{
		UserID:     user.ID,
		Username:   user.Name,
		RankChange: detail.RankChange,
		History:    make([]dto.RankingHistoryEntry, len(detail.History)),
	}

	if user.DisplayName != nil {
		resp.DisplayName = *user.DisplayName
	}
	if user.AvatarURL != nil {
		resp.AvatarURL = *user.AvatarURL
	}

	if detail.CurrentScore != nil {
		resp.FollowerCount = detail.CurrentScore.FollowerCount
		resp.FollowerGrowthRate = detail.CurrentScore.FollowerGrowthRate
		resp.BlogPostVelocity = detail.CurrentScore.BlogPostVelocity
		resp.CompositeScore = detail.CurrentScore.CompositeScore
		resp.CalculationDate = detail.CurrentScore.CalculationDate

		if detail.CurrentScore.RankPosition != nil {
			rank := *detail.CurrentScore.RankPosition
			resp.Rank = &rank
		}
	}

	if detail.PreviousRank != nil {
		prevRank := *detail.PreviousRank
		resp.PreviousRank = &prevRank
	}

	for i, h := range detail.History {
		resp.History[i] = dto.RankingHistoryEntry{
			RankPosition:   h.RankPosition,
			CompositeScore: h.CompositeScore,
			FollowerCount:  h.FollowerCount,
			RecordedAt:     h.RecordedAt,
		}
	}

	return resp
}
