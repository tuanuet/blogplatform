package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"math"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository"
	"github.com/google/uuid"
)

// Reuse ErrUserNotFound from user_service.go

// RankingService defines the interface for ranking operations
type RankingService interface {
	// CalculateVelocityScore calculates the velocity score for a user
	CalculateVelocityScore(ctx context.Context, userID uuid.UUID) (*entity.UserVelocityScore, error)

	// CalculateAllVelocityScores calculates velocity scores for all users
	CalculateAllVelocityScores(ctx context.Context) error

	// AssignRankPositions assigns rank positions based on composite scores
	AssignRankPositions(ctx context.Context) error

	// GetUserRankingDetail gets detailed ranking information for a user
	GetUserRankingDetail(ctx context.Context, userID uuid.UUID) (*UserRankingDetail, error)

	// ArchiveRankings archives current rankings to history
	ArchiveRankings(ctx context.Context) error
}

// UserRankingDetail contains detailed ranking information
type UserRankingDetail struct {
	CurrentScore *entity.UserVelocityScore
	History      []entity.UserRankingHistory
	RankChange   int
	PreviousRank *int
}

// RankingConfig contains configuration for ranking calculations
type RankingConfig struct {
	FollowerGrowthWeight   float64
	BlogPostVelocityWeight float64
	TimeWindowDays         int
	MinFollowersForRate    int
}

// DefaultRankingConfig returns the default ranking configuration
func DefaultRankingConfig() RankingConfig {
	return RankingConfig{
		FollowerGrowthWeight:   0.6,
		BlogPostVelocityWeight: 0.4,
		TimeWindowDays:         30,
		MinFollowersForRate:    100,
	}
}

type rankingService struct {
	velocityScoreRepo    repository.UserVelocityScoreRepository
	rankingHistoryRepo   repository.UserRankingHistoryRepository
	followerSnapshotRepo repository.UserFollowerSnapshotRepository
	config               RankingConfig
}

// NewRankingService creates a new ranking service
func NewRankingService(
	velocityScoreRepo repository.UserVelocityScoreRepository,
	rankingHistoryRepo repository.UserRankingHistoryRepository,
	followerSnapshotRepo repository.UserFollowerSnapshotRepository,
) RankingService {
	return &rankingService{
		velocityScoreRepo:    velocityScoreRepo,
		rankingHistoryRepo:   rankingHistoryRepo,
		followerSnapshotRepo: followerSnapshotRepo,
		config:               DefaultRankingConfig(),
	}
}

func (s *rankingService) CalculateVelocityScore(ctx context.Context, userID uuid.UUID) (*entity.UserVelocityScore, error) {
	// Get current follower count
	currentFollowerCount, err := s.followerSnapshotRepo.CountFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get follower count from 30 days ago
	timeWindowStart := time.Now().AddDate(0, 0, -s.config.TimeWindowDays)
	snapshot30DaysAgo, err := s.followerSnapshotRepo.FindByUserIDAndDate(ctx, userID, timeWindowStart)
	if err != nil {
		return nil, err
	}

	var followerCount30DaysAgo int64 = 0
	if snapshot30DaysAgo != nil {
		followerCount30DaysAgo = int64(snapshot30DaysAgo.FollowerCount)
	}

	// Calculate follower growth rate
	followerGrowthRate := s.calculateGrowthRate(currentFollowerCount, followerCount30DaysAgo)

	// Count blogs published in last 30 days
	now := time.Now()
	blogCount, err := s.followerSnapshotRepo.CountBlogs(ctx, userID, timeWindowStart, now)
	if err != nil {
		return nil, err
	}

	// Calculate blog post velocity (posts per day)
	blogPostVelocity := float64(blogCount) / float64(s.config.TimeWindowDays)

	// Calculate composite score
	compositeScore := (followerGrowthRate * s.config.FollowerGrowthWeight) +
		(blogPostVelocity * s.config.BlogPostVelocityWeight)

	// Create or update velocity score
	score := &entity.UserVelocityScore{
		UserID:             userID,
		FollowerCount:      int(currentFollowerCount),
		FollowerGrowthRate: followerGrowthRate,
		BlogPostVelocity:   blogPostVelocity,
		CompositeScore:     compositeScore,
		CalculationDate:    now,
	}

	// Check if score already exists
	existingScore, err := s.velocityScoreRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if existingScore != nil {
		score.ID = existingScore.ID
		score.RankPosition = existingScore.RankPosition
		score.CreatedAt = existingScore.CreatedAt
	}

	if err := s.velocityScoreRepo.Save(ctx, score); err != nil {
		return nil, err
	}

	return score, nil
}

func (s *rankingService) CalculateAllVelocityScores(ctx context.Context) error {
	// This would typically iterate through all users
	// For now, we'll return nil as the actual implementation would require
	// a user repository to get all user IDs
	return nil
}

func (s *rankingService) AssignRankPositions(ctx context.Context) error {
	// Get all velocity scores ordered by composite score
	scores, err := s.velocityScoreRepo.ListTopRanked(ctx, 10000) // Large limit to get all
	if err != nil {
		return err
	}

	// Assign rank positions
	for i, score := range scores {
		rank := i + 1
		if err := s.velocityScoreRepo.UpdateRankPosition(ctx, score.UserID, rank); err != nil {
			return err
		}
	}

	return nil
}

func (s *rankingService) GetUserRankingDetail(ctx context.Context, userID uuid.UUID) (*UserRankingDetail, error) {
	// Get current score
	currentScore, err := s.velocityScoreRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if currentScore == nil {
		return nil, ErrUserNotFound
	}

	// Get ranking history
	history, err := s.rankingHistoryRepo.FindByUserID(ctx, userID, 30)
	if err != nil {
		return nil, err
	}

	// Calculate rank change
	var rankChange int
	var previousRank *int
	if len(history) > 0 {
		previousRank = &history[0].RankPosition
		if currentScore.RankPosition != nil {
			rankChange = *previousRank - *currentScore.RankPosition
		}
	}

	return &UserRankingDetail{
		CurrentScore: currentScore,
		History:      history,
		RankChange:   rankChange,
		PreviousRank: previousRank,
	}, nil
}

func (s *rankingService) ArchiveRankings(ctx context.Context) error {
	// Get all current rankings
	scores, err := s.velocityScoreRepo.ListTopRanked(ctx, 10000)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, score := range scores {
		if score.RankPosition != nil {
			history := &entity.UserRankingHistory{
				UserID:         score.UserID,
				RankPosition:   *score.RankPosition,
				CompositeScore: score.CompositeScore,
				FollowerCount:  score.FollowerCount,
				RecordedAt:     now,
			}
			if err := s.rankingHistoryRepo.Create(ctx, history); err != nil {
				return err
			}
		}
	}

	return nil
}

// calculateGrowthRate calculates the growth rate with special handling for small follower counts
func (s *rankingService) calculateGrowthRate(current, previous int64) float64 {
	if previous == 0 {
		// If no previous data, use absolute growth normalized by minimum threshold
		return float64(current) / float64(s.config.MinFollowersForRate)
	}

	if previous < int64(s.config.MinFollowersForRate) {
		// For small follower counts, use absolute growth normalized
		growth := float64(current - previous)
		return growth / float64(s.config.MinFollowersForRate)
	}

	// Normal percentage growth rate
	growthRate := float64(current-previous) / float64(previous)

	// Cap extreme growth rates to prevent gaming
	maxGrowthRate := 10.0 // 1000% growth cap
	return math.Min(growthRate, maxGrowthRate)
}
