package service

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/interfaces/http/dto"
	"github.com/google/uuid"
)

// FraudDetectionService defines the interface for fraud detection operations
type FraudDetectionService interface {
	// GetUserRiskScore retrieves the risk score for a specific user
	GetUserRiskScore(ctx context.Context, userID uuid.UUID) (*dto.RiskScoreResponse, error)

	// GetFraudDashboard retrieves paginated list of users for admin dashboard
	GetFraudDashboard(ctx context.Context, req dto.FraudDashboardRequest) (*dto.FraudDashboardResponse, error)

	// ReviewUser marks a user as reviewed by an admin
	ReviewUser(ctx context.Context, adminID, userID uuid.UUID, req dto.ReviewUserRequest) (*dto.ReviewUserResponse, error)

	// BanUser bans a user account due to fraudulent activity
	BanUser(ctx context.Context, adminID, userID uuid.UUID, req dto.BanUserRequest) (*dto.BanUserResponse, error)

	// GetFraudTrends retrieves analytics data about fraud trends
	GetFraudTrends(ctx context.Context, req dto.FraudTrendsRequest) (*dto.FraudTrendsResponse, error)

	// TriggerBatchAnalysis starts a batch job to analyze followers for bot detection
	TriggerBatchAnalysis(ctx context.Context, req dto.BatchAnalyzeRequest) (*dto.BatchAnalyzeResponse, error)

	// GetUserBadgeStatus retrieves the badge status for a user
	GetUserBadgeStatus(ctx context.Context, userID uuid.UUID) (*dto.UserBadgeResponse, error)

	// GetUserBotNotifications retrieves notifications about flagged bot followers
	GetUserBotNotifications(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]dto.BotFollowerNotificationResponse, error)
}

// FraudDetectionRepository defines the interface for fraud detection data access
type FraudDetectionRepository interface {
	// Follower Events
	CreateFollowerEvent(ctx context.Context, event *entity.FollowerEvent) error
	GetFollowerEventsByUser(ctx context.Context, userID uuid.UUID, from, to *time.Time) ([]entity.FollowerEvent, error)
	GetFollowerEventsByIP(ctx context.Context, ipAddress string, from, to *time.Time) ([]entity.FollowerEvent, error)

	// Bot Detection Signals
	CreateBotSignal(ctx context.Context, signal *entity.BotDetectionSignal) error
	GetBotSignalsByUser(ctx context.Context, userID uuid.UUID, processed bool) ([]entity.BotDetectionSignal, error)
	GetUnprocessedBotSignals(ctx context.Context, limit int) ([]entity.BotDetectionSignal, error)
	MarkBotSignalAsProcessed(ctx context.Context, signalID uuid.UUID) error

	// User Risk Scores
	CreateOrUpdateRiskScore(ctx context.Context, score *entity.UserRiskScore) error
	GetRiskScoreByUser(ctx context.Context, userID uuid.UUID) (*entity.UserRiskScore, error)
	GetUsersByRiskScoreRange(ctx context.Context, minScore, maxScore int, page, pageSize int) ([]entity.UserRiskScore, int, error)

	// Badge Status
	CreateOrUpdateBadgeStatus(ctx context.Context, status *entity.UserBadgeStatus) error
	GetBadgeStatusByUser(ctx context.Context, userID uuid.UUID) (*entity.UserBadgeStatus, error)

	// Admin Reviews
	CreateAdminReview(ctx context.Context, review *entity.AdminReview) error
	GetAdminReviewsByUser(ctx context.Context, userID uuid.UUID, limit int) ([]entity.AdminReview, error)
	GetLastReviewByUser(ctx context.Context, userID uuid.UUID) (*entity.AdminReview, error)

	// Notifications
	CreateBotNotification(ctx context.Context, notification *entity.BotFollowerNotification) error
	GetBotNotificationsByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]entity.BotFollowerNotification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID uuid.UUID) error

	// Analytics
	GetSignalsCountByType(ctx context.Context, from, to time.Time) (map[string]int, error)
	GetNewSuspiciousAccountsCount(ctx context.Context, from, to time.Time) (int, error)
	GetBannedAccountsCount(ctx context.Context, from, to time.Time) (int, error)
	GetReviewedAccountsCount(ctx context.Context, from, to time.Time) (int, error)
	GetAverageRiskScore(ctx context.Context, from, to time.Time) (float64, error)
	GetRiskScoreDistribution(ctx context.Context, from, to time.Time) (map[string]int, error)
	GetDailyFraudStats(ctx context.Context, from, to time.Time) ([]dto.DailyFraudStat, error)
}

// BotDetectionAlgorithm defines the interface for bot detection algorithms
type BotDetectionAlgorithm interface {
	// AnalyzeFollower analyzes a follower event and returns bot signals if detected
	AnalyzeFollower(ctx context.Context, event entity.FollowerEvent, recentEvents []entity.FollowerEvent) ([]entity.BotDetectionSignal, error)

	// CalculateRiskScore calculates the overall risk score for a user
	CalculateRiskScore(ctx context.Context, userID uuid.UUID, signals []entity.BotDetectionSignal, followerCount int) (*entity.UserRiskScore, error)

	// DetectCoordinatedBots detects networks of coordinated bot accounts
	DetectCoordinatedBots(ctx context.Context, signals []entity.BotDetectionSignal) ([][]uuid.UUID, error)
}

// NotificationService defines the interface for sending notifications
type NotificationService interface {
	// SendBotFollowerNotification sends notification to user about flagged bot followers
	SendBotFollowerNotification(ctx context.Context, userID uuid.UUID, notifications []dto.BotFollowerNotificationResponse) error

	// SendBadgeStatusUpdate notifies user about badge status changes
	SendBadgeStatusUpdate(ctx context.Context, userID uuid.UUID, status string, reason string) error
}

// BatchJobService defines the interface for batch processing jobs
type BatchJobService interface {
	// StartBatchAnalysis starts the batch analysis job
	StartBatchAnalysis(ctx context.Context, dateFrom, dateTo *time.Time) (uuid.UUID, error)

	// GetBatchJobStatus retrieves the status of a batch job
	GetBatchJobStatus(ctx context.Context, jobID uuid.UUID) (*dto.BatchAnalyzeResponse, error)
}
