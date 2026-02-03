package valueobject

import (
	"time"

	"github.com/google/uuid"
)

// RiskScoreResult represents a user's risk assessment result
type RiskScoreResult struct {
	ID                        uuid.UUID
	UserID                    uuid.UUID
	OverallScore              int     // 0-100
	FollowerAuthenticityScore int     // 0-100
	EngagementQualityScore    int     // 0-100
	AccountAgeFactor          float64 // 0-1
	BadgeStatus               string  // "eligible", "active", "revoked", "none"
	LastCalculatedAt          time.Time
}

// FraudDashboardFilter represents filters for the admin dashboard
type FraudDashboardFilter struct {
	MinRiskScore *int
	MaxRiskScore *int
	SignalTypes  []string
	ReviewStatus string // pending, reviewed, banned, cleared
	FromDate     *time.Time
	ToDate       *time.Time
	Page         int
	PageSize     int
}

// FraudDashboardUser represents a user entry in the admin dashboard
type FraudDashboardUser struct {
	UserID                uuid.UUID
	Username              string
	Email                 string
	OverallScore          int
	FollowerCount         int
	BotFollowerCount      int
	ActiveSignals         []BotSignalSummary
	LastReviewAction      string // "reviewed", "banned", "warned", "cleared", ""
	LastReviewedAt        *time.Time
	RiskScoreCalculatedAt time.Time
}

// BotSignalSummary represents a summary of a bot detection signal
type BotSignalSummary struct {
	SignalType      string
	ConfidenceScore float64
	DetectedAt      time.Time
}

// FraudDashboardResult represents the admin dashboard data
type FraudDashboardResult struct {
	Users      []FraudDashboardUser
	TotalCount int
	Page       int
	PageSize   int
	TotalPages int
}

// ReviewUserCommand represents details to mark a user as reviewed
type ReviewUserCommand struct {
	Notes string
}

// ReviewUserResult represents the result of reviewing a user
type ReviewUserResult struct {
	ReviewID   uuid.UUID
	UserID     uuid.UUID
	AdminID    uuid.UUID
	Action     string
	Notes      string
	ReviewedAt time.Time
}

// BanUserCommand represents details to ban a user
type BanUserCommand struct {
	Reason string
	Notes  string
}

// BanUserResult represents the result of banning a user
type BanUserResult struct {
	ReviewID uuid.UUID
	UserID   uuid.UUID
	AdminID  uuid.UUID
	Action   string // "banned"
	Reason   string
	Notes    string
	BannedAt time.Time
}

// FraudTrendsFilter represents filters for fraud analytics
type FraudTrendsFilter struct {
	Period   string // 24h, 7d, 30d, 90d
	FromDate *time.Time
	ToDate   *time.Time
}

// FraudTrendsResult represents fraud analytics data
type FraudTrendsResult struct {
	Period                string
	FromDate              time.Time
	ToDate                time.Time
	TotalBotSignals       int
	SignalsByType         map[string]int
	NewSuspiciousAccounts int
	BannedAccounts        int
	ReviewedAccounts      int
	AverageRiskScore      float64
	RiskScoreDistribution map[string]int
	DailyStats            []DailyFraudStat
}

// DailyFraudStat represents fraud statistics for a single day
type DailyFraudStat struct {
	Date                  string
	NewSignals            int
	NewSuspiciousAccounts int
	BannedAccounts        int
}

// BatchAnalyzeCommand represents options to trigger batch analysis
type BatchAnalyzeCommand struct {
	DateFrom *time.Time
	DateTo   *time.Time
}

// BatchAnalyzeResult represents the result of batch analysis
type BatchAnalyzeResult struct {
	JobID              uuid.UUID
	Status             string // "started", "completed", "failed"
	StartedAt          time.Time
	CompletedAt        *time.Time
	ProcessedFollowers int
	NewSignalsDetected int
	UsersScored        int
	Message            string
}

// UserBadgeResult represents a user's badge status
type UserBadgeResult struct {
	UserID        uuid.UUID
	BadgeType     string
	Status        string
	EligibleSince *time.Time
	ActivatedAt   *time.Time
}

// BotFollowerNotificationResult represents a notification about flagged followers
type BotFollowerNotificationResult struct {
	ID              uuid.UUID
	BotFollowerID   uuid.UUID
	BotFollowerName string
	SignalType      string
	ConfidenceScore float64
	SentAt          time.Time
	ReadAt          *time.Time
}
