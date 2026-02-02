package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/interfaces/http/dto"
	"github.com/google/uuid"
)

// Constants for fraud detection service
const (
	// Risk score range
	minRiskScore = 0
	maxRiskScore = 100

	// Badge statuses
	badgeStatusNone   = "none"
	badgeStatusActive = "active"

	// Admin actions
	actionReviewed = "reviewed"
	actionBanned   = "banned"

	// Job statuses
	jobStatusStarted = "started"

	// Time periods
	period24h = "24h"
	period7d  = "7d"
	period30d = "30d"
	period90d = "90d"

	// Default values
	defaultPeriodDays = 7
)

// fraudDetectionService implements the FraudDetectionService interface
type fraudDetectionService struct {
	repo     FraudDetectionRepository
	notifier NotificationService
	batchJob BatchJobService
}

// NewFraudDetectionService creates a new fraud detection service instance
func NewFraudDetectionService(repo FraudDetectionRepository, notifier NotificationService, batchJob BatchJobService) FraudDetectionService {
	return &fraudDetectionService{
		repo:     repo,
		notifier: notifier,
		batchJob: batchJob,
	}
}

// GetUserRiskScore retrieves the risk score for a specific user
func (s *fraudDetectionService) GetUserRiskScore(ctx context.Context, userID uuid.UUID) (*dto.RiskScoreResponse, error) {
	score, err := s.repo.GetRiskScoreByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if score == nil {
		return nil, nil
	}

	badge, err := s.repo.GetBadgeStatusByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	badgeStatus := badgeStatusNone
	if badge != nil {
		badgeStatus = badge.Status
	}

	return &dto.RiskScoreResponse{
		ID:                        score.ID,
		UserID:                    score.UserID,
		OverallScore:              score.OverallScore,
		FollowerAuthenticityScore: score.FollowerAuthenticityScore,
		EngagementQualityScore:    score.EngagementQualityScore,
		AccountAgeFactor:          score.AccountAgeFactor,
		BadgeStatus:               badgeStatus,
		LastCalculatedAt:          score.LastCalculatedAt,
	}, nil
}

// GetFraudDashboard retrieves paginated list of users for admin dashboard
func (s *fraudDetectionService) GetFraudDashboard(ctx context.Context, req dto.FraudDashboardRequest) (*dto.FraudDashboardResponse, error) {
	minScore := req.MinRiskScore
	maxScore := req.MaxRiskScore
	if maxScore == 0 {
		maxScore = maxRiskScore
	}

	scores, totalCount, err := s.repo.GetUsersByRiskScoreRange(ctx, minScore, maxScore, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	users := make([]dto.FraudDashboardUser, 0, len(scores))
	for _, score := range scores {
		signals, err := s.repo.GetBotSignalsByUser(ctx, score.UserID, false)
		if err != nil {
			continue
		}

		signalSummaries := make([]dto.BotSignalSummary, 0, len(signals))
		for _, signal := range signals {
			signalSummaries = append(signalSummaries, dto.BotSignalSummary{
				SignalType:      signal.SignalType,
				ConfidenceScore: signal.ConfidenceScore,
				DetectedAt:      signal.DetectedAt,
			})
		}

		lastReview, err := s.repo.GetLastReviewByUser(ctx, score.UserID)
		lastAction := ""
		var lastReviewedAt *time.Time
		if err == nil && lastReview != nil {
			lastAction = lastReview.Action
			lastReviewedAt = &lastReview.ReviewedAt
		}

		users = append(users, dto.FraudDashboardUser{
			UserID:                score.UserID,
			OverallScore:          score.OverallScore,
			ActiveSignals:         signalSummaries,
			LastReviewAction:      lastAction,
			LastReviewedAt:        lastReviewedAt,
			RiskScoreCalculatedAt: score.LastCalculatedAt,
		})
	}

	totalPages := (totalCount + req.PageSize - 1) / req.PageSize

	return &dto.FraudDashboardResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// getRiskScoreValue safely retrieves the risk score value for a user
func (s *fraudDetectionService) getRiskScoreValue(ctx context.Context, userID uuid.UUID) (int, error) {
	riskScore, err := s.repo.GetRiskScoreByUser(ctx, userID)
	if err != nil {
		return 0, err
	}

	if riskScore != nil {
		return riskScore.OverallScore, nil
	}
	return 0, nil
}

// createAdminReview creates an admin review record
func (s *fraudDetectionService) createAdminReview(ctx context.Context, adminID, userID uuid.UUID, action string, riskScoreValue int, notes string) (*entity.AdminReview, error) {
	review := &entity.AdminReview{
		ID:                uuid.New(),
		AdminID:           adminID,
		UserID:            userID,
		Action:            action,
		RiskScoreAtReview: riskScoreValue,
		Notes:             notes,
		ReviewedAt:        time.Now(),
	}

	if err := s.repo.CreateAdminReview(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}

// ReviewUser marks a user as reviewed by an admin
func (s *fraudDetectionService) ReviewUser(ctx context.Context, adminID, userID uuid.UUID, req dto.ReviewUserRequest) (*dto.ReviewUserResponse, error) {
	riskScoreValue, err := s.getRiskScoreValue(ctx, userID)
	if err != nil {
		return nil, err
	}

	review, err := s.createAdminReview(ctx, adminID, userID, actionReviewed, riskScoreValue, req.Notes)
	if err != nil {
		return nil, err
	}

	return &dto.ReviewUserResponse{
		ReviewID:   review.ID,
		UserID:     review.UserID,
		AdminID:    review.AdminID,
		Action:     review.Action,
		Notes:      review.Notes,
		ReviewedAt: review.ReviewedAt,
	}, nil
}

// buildBanNotes creates the notes string for a ban action
func buildBanNotes(reason, notes string) string {
	return fmt.Sprintf("Reason: %s. %s", reason, notes)
}

// BanUser bans a user account due to fraudulent activity
func (s *fraudDetectionService) BanUser(ctx context.Context, adminID, userID uuid.UUID, req dto.BanUserRequest) (*dto.BanUserResponse, error) {
	riskScoreValue, err := s.getRiskScoreValue(ctx, userID)
	if err != nil {
		return nil, err
	}

	notes := buildBanNotes(req.Reason, req.Notes)
	review, err := s.createAdminReview(ctx, adminID, userID, actionBanned, riskScoreValue, notes)
	if err != nil {
		return nil, err
	}

	return &dto.BanUserResponse{
		ReviewID: review.ID,
		UserID:   review.UserID,
		AdminID:  review.AdminID,
		Action:   review.Action,
		Reason:   req.Reason,
		Notes:    req.Notes,
		BannedAt: review.ReviewedAt,
	}, nil
}

// calculatePeriodDates calculates the start and end dates for a given period
func calculatePeriodDates(period string) (from, to time.Time) {
	to = time.Now()

	switch period {
	case period24h:
		from = to.AddDate(0, 0, -1)
	case period7d:
		from = to.AddDate(0, 0, -7)
	case period30d:
		from = to.AddDate(0, 0, -30)
	case period90d:
		from = to.AddDate(0, 0, -90)
	default:
		from = to.AddDate(0, 0, -defaultPeriodDays)
	}

	return from, to
}

// GetFraudTrends retrieves analytics data about fraud trends
func (s *fraudDetectionService) GetFraudTrends(ctx context.Context, req dto.FraudTrendsRequest) (*dto.FraudTrendsResponse, error) {
	from, to := calculatePeriodDates(req.Period)

	signalsByType, err := s.repo.GetSignalsCountByType(ctx, from, to)
	if err != nil {
		return nil, err
	}

	totalSignals := 0
	for _, count := range signalsByType {
		totalSignals += count
	}

	newSuspicious, err := s.repo.GetNewSuspiciousAccountsCount(ctx, from, to)
	if err != nil {
		return nil, err
	}

	bannedCount, err := s.repo.GetBannedAccountsCount(ctx, from, to)
	if err != nil {
		return nil, err
	}

	reviewedCount, err := s.repo.GetReviewedAccountsCount(ctx, from, to)
	if err != nil {
		return nil, err
	}

	avgRiskScore, err := s.repo.GetAverageRiskScore(ctx, from, to)
	if err != nil {
		return nil, err
	}

	riskDistribution, err := s.repo.GetRiskScoreDistribution(ctx, from, to)
	if err != nil {
		return nil, err
	}

	dailyStats, err := s.repo.GetDailyFraudStats(ctx, from, to)
	if err != nil {
		return nil, err
	}

	return &dto.FraudTrendsResponse{
		Period:                req.Period,
		FromDate:              from,
		ToDate:                to,
		TotalBotSignals:       totalSignals,
		SignalsByType:         signalsByType,
		NewSuspiciousAccounts: newSuspicious,
		BannedAccounts:        bannedCount,
		ReviewedAccounts:      reviewedCount,
		AverageRiskScore:      avgRiskScore,
		RiskScoreDistribution: riskDistribution,
		DailyStats:            dailyStats,
	}, nil
}

// TriggerBatchAnalysis starts a batch job to analyze followers for bot detection
func (s *fraudDetectionService) TriggerBatchAnalysis(ctx context.Context, req dto.BatchAnalyzeRequest) (*dto.BatchAnalyzeResponse, error) {
	jobID, err := s.batchJob.StartBatchAnalysis(ctx, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}

	return &dto.BatchAnalyzeResponse{
		JobID:     jobID,
		Status:    jobStatusStarted,
		StartedAt: time.Now(),
		Message:   "Batch analysis job started successfully",
	}, nil
}

// GetUserBadgeStatus retrieves the badge status for a user
func (s *fraudDetectionService) GetUserBadgeStatus(ctx context.Context, userID uuid.UUID) (*dto.UserBadgeResponse, error) {
	badge, err := s.repo.GetBadgeStatusByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if badge == nil {
		return nil, nil
	}

	return &dto.UserBadgeResponse{
		UserID:        badge.UserID,
		BadgeType:     badge.BadgeType,
		Status:        badge.Status,
		EligibleSince: badge.EligibleSince,
		ActivatedAt:   badge.ActivatedAt,
	}, nil
}

// GetUserBotNotifications retrieves notifications about flagged bot followers
func (s *fraudDetectionService) GetUserBotNotifications(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]dto.BotFollowerNotificationResponse, error) {
	notifications, err := s.repo.GetBotNotificationsByUser(ctx, userID, unreadOnly)
	if err != nil {
		return nil, err
	}

	response := make([]dto.BotFollowerNotificationResponse, 0, len(notifications))
	for _, notif := range notifications {
		response = append(response, dto.BotFollowerNotificationResponse{
			ID:            notif.ID,
			BotFollowerID: notif.BotFollowerID,
			SentAt:        notif.SentAt,
			ReadAt:        notif.ReadAt,
		})
	}

	return response, nil
}
