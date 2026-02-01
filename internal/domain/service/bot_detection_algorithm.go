package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/google/uuid"
)

// Constants for bot detection
type BotDetectionConfig struct {
	// Rapid follows detection
	RapidFollowThreshold   int           // Number of follows to trigger rapid follow detection
	RapidFollowTimeWindow  time.Duration // Time window for rapid follows
	RapidFollowMinInterval time.Duration // Minimum time between follows (too fast = bot)

	// IP clustering detection
	IPClusterThreshold      int // Number of accounts from same IP
	IPClusterTimeWindowDays int // Days to look back for IP clustering

	// No profile activity detection
	MinProfileAgeDays    int // Minimum account age to check for activity
	MinPostsForActive    int // Minimum posts to be considered active
	MinCommentsForActive int // Minimum comments to be considered active

	// Suspicious engagement detection
	EngagementVelocityThreshold float64 // Max growth rate per day
	EngagementSpikeMultiplier   float64 // Multiplier over average to flag as spike
}

// DefaultBotDetectionConfig returns the default bot detection configuration
func DefaultBotDetectionConfig() *BotDetectionConfig {
	return &BotDetectionConfig{
		RapidFollowThreshold:        50,
		RapidFollowTimeWindow:       time.Hour,
		RapidFollowMinInterval:      time.Second,
		IPClusterThreshold:          10,
		IPClusterTimeWindowDays:     7,
		MinProfileAgeDays:           30,
		MinPostsForActive:           1,
		MinCommentsForActive:        0,
		EngagementVelocityThreshold: 1000.0,
		EngagementSpikeMultiplier:   10.0,
	}
}

// Constants for confidence scoring
const (
	// Base confidence levels
	confidenceLow      = 0.5
	confidenceMedium   = 0.6
	confidenceHigh     = 0.8
	confidenceVeryHigh = 0.9
	confidenceExtreme  = 0.85

	// Multipliers
	thresholdMultiplier = 2
	riskScoreMultiplier = 2
	accountAgeDivisor   = 200.0

	// Score ranges
	maxAuthenticityScore      = 100
	maxEngagementScore        = 100
	defaultAccountAgeFactor   = 1.0
	defaultCalculationVersion = "v1.0"

	// Signal types
	signalTypeRapidFollows         = "rapid_follows"
	signalTypeIPCluster            = "ip_cluster"
	signalTypeNoProfile            = "no_profile"
	signalTypeSuspiciousEngagement = "suspicious_engagement"
)

// botDetectionAlgorithm implements rule-based bot detection
type botDetectionAlgorithm struct {
	config *BotDetectionConfig
	repo   FraudDetectionRepository
}

// NewBotDetectionAlgorithm creates a new bot detection algorithm instance
func NewBotDetectionAlgorithm(config *BotDetectionConfig, repo FraudDetectionRepository) BotDetectionAlgorithm {
	if config == nil {
		config = DefaultBotDetectionConfig()
	}
	return &botDetectionAlgorithm{
		config: config,
		repo:   repo,
	}
}

// AnalyzeFollower analyzes a follower event and returns bot signals if detected
func (a *botDetectionAlgorithm) AnalyzeFollower(ctx context.Context, event entity.FollowerEvent, recentEvents []entity.FollowerEvent) ([]entity.BotDetectionSignal, error) {
	var signals []entity.BotDetectionSignal

	// Check for rapid follows
	rapidFollowSignal := a.detectRapidFollows(event, recentEvents)
	if rapidFollowSignal != nil {
		signals = append(signals, *rapidFollowSignal)
	}

	// Check for IP clustering
	ipClusterSignal, err := a.detectIPClustering(ctx, event)
	if err != nil {
		return signals, err
	}
	if ipClusterSignal != nil {
		signals = append(signals, *ipClusterSignal)
	}

	return signals, nil
}

// detectRapidFollows detects if a follower is following too many accounts too quickly
func (a *botDetectionAlgorithm) detectRapidFollows(event entity.FollowerEvent, recentEvents []entity.FollowerEvent) *entity.BotDetectionSignal {
	if len(recentEvents) < a.config.RapidFollowThreshold {
		return nil
	}

	// Count follows within the time window
	cutoff := event.Timestamp.Add(-a.config.RapidFollowTimeWindow)
	followCount := 0
	minInterval := time.Duration(0)

	for _, e := range recentEvents {
		if e.Timestamp.After(cutoff) && e.FollowerID == event.FollowerID {
			followCount++
			if minInterval == 0 || e.Timestamp.Sub(event.Timestamp) < minInterval {
				minInterval = e.Timestamp.Sub(event.Timestamp)
			}
		}
	}

	if followCount >= a.config.RapidFollowThreshold {
		confidence := 0.5
		if minInterval < a.config.RapidFollowMinInterval {
			confidence = 0.9 // Very high confidence if follows are too fast
		} else if followCount > a.config.RapidFollowThreshold*2 {
			confidence = 0.8 // High confidence for extreme rapid follows
		}

		evidence := fmt.Sprintf("Followed %d accounts in %v (threshold: %d)",
			followCount, a.config.RapidFollowTimeWindow, a.config.RapidFollowThreshold)

		return &entity.BotDetectionSignal{
			ID:              uuid.New(),
			UserID:          event.FollowerID,
			SignalType:      "rapid_follows",
			ConfidenceScore: confidence,
			DetectedAt:      time.Now(),
			Evidence:        evidence,
			Processed:       false,
		}
	}

	return nil
}

// detectIPClustering detects if multiple accounts are following from the same IP
func (a *botDetectionAlgorithm) detectIPClustering(ctx context.Context, event entity.FollowerEvent) (*entity.BotDetectionSignal, error) {
	if event.IPAddress == "" {
		return nil, nil
	}

	// Get all follower events from this IP in the last N days
	from := time.Now().AddDate(0, 0, -a.config.IPClusterTimeWindowDays)
	to := time.Now()

	events, err := a.repo.GetFollowerEventsByIP(ctx, event.IPAddress, &from, &to)
	if err != nil {
		return nil, err
	}

	// Count unique follower accounts from this IP
	uniqueFollowers := make(map[uuid.UUID]bool)
	for _, e := range events {
		uniqueFollowers[e.FollowerID] = true
	}

	if len(uniqueFollowers) >= a.config.IPClusterThreshold {
		confidence := 0.6
		if len(uniqueFollowers) >= a.config.IPClusterThreshold*2 {
			confidence = 0.85
		}

		// Get related accounts (other suspicious accounts from same IP)
		relatedAccounts := make([]uuid.UUID, 0, len(uniqueFollowers))
		for id := range uniqueFollowers {
			if id != event.FollowerID {
				relatedAccounts = append(relatedAccounts, id)
			}
		}

		evidence := fmt.Sprintf("IP %s has %d unique follower accounts in %d days (threshold: %d)",
			event.IPAddress, len(uniqueFollowers), a.config.IPClusterTimeWindowDays, a.config.IPClusterThreshold)

		return &entity.BotDetectionSignal{
			ID:              uuid.New(),
			UserID:          event.FollowerID,
			SignalType:      "ip_cluster",
			ConfidenceScore: confidence,
			DetectedAt:      time.Now(),
			RelatedAccounts: relatedAccounts,
			Evidence:        evidence,
			Processed:       false,
		}, nil
	}

	return nil, nil
}

// CalculateRiskScore calculates the overall risk score for a user
func (a *botDetectionAlgorithm) CalculateRiskScore(ctx context.Context, userID uuid.UUID, signals []entity.BotDetectionSignal, followerCount int) (*entity.UserRiskScore, error) {
	if len(signals) == 0 {
		// No signals - low risk
		return &entity.UserRiskScore{
			ID:                        uuid.New(),
			UserID:                    userID,
			OverallScore:              0,
			FollowerAuthenticityScore: 100,
			EngagementQualityScore:    100,
			AccountAgeFactor:          1.0,
			CalculationVersion:        "v1.0",
			LastCalculatedAt:          time.Now(),
		}, nil
	}

	// Calculate scores based on signals
	totalConfidence := 0.0
	signalTypeWeights := map[string]float64{
		"rapid_follows":         0.3,
		"ip_cluster":            0.25,
		"no_profile":            0.2,
		"suspicious_engagement": 0.25,
	}

	for _, signal := range signals {
		weight := signalTypeWeights[signal.SignalType]
		if weight == 0 {
			weight = 0.1 // Default weight
		}
		totalConfidence += signal.ConfidenceScore * weight
	}

	// Normalize to 0-100 scale
	riskScore := int(totalConfidence * 100 * 2) // Multiply by 2 to amplify effect
	if riskScore > 100 {
		riskScore = 100
	}

	// Calculate follower authenticity score (inverse of risk)
	followerAuthenticity := 100 - riskScore

	// Calculate engagement quality (simplified)
	engagementQuality := 100 - int(float64(riskScore)*0.5)
	if engagementQuality < 0 {
		engagementQuality = 0
	}

	// Account age factor (older accounts are less likely to be bots)
	// This is a placeholder - in real implementation, you'd fetch user creation date
	accountAgeFactor := 1.0 - (float64(riskScore) / 200.0)
	if accountAgeFactor < 0 {
		accountAgeFactor = 0
	}

	return &entity.UserRiskScore{
		ID:                        uuid.New(),
		UserID:                    userID,
		OverallScore:              riskScore,
		FollowerAuthenticityScore: followerAuthenticity,
		EngagementQualityScore:    engagementQuality,
		AccountAgeFactor:          accountAgeFactor,
		CalculationVersion:        "v1.0",
		LastCalculatedAt:          time.Now(),
	}, nil
}

// DetectCoordinatedBots detects networks of coordinated bot accounts
func (a *botDetectionAlgorithm) DetectCoordinatedBots(ctx context.Context, signals []entity.BotDetectionSignal) ([][]uuid.UUID, error) {
	// Group signals by related accounts
	accountGraph := make(map[uuid.UUID]map[uuid.UUID]bool)

	for _, signal := range signals {
		if len(signal.RelatedAccounts) == 0 {
			continue
		}

		if accountGraph[signal.UserID] == nil {
			accountGraph[signal.UserID] = make(map[uuid.UUID]bool)
		}

		for _, relatedID := range signal.RelatedAccounts {
			accountGraph[signal.UserID][relatedID] = true
			if accountGraph[relatedID] == nil {
				accountGraph[relatedID] = make(map[uuid.UUID]bool)
			}
			accountGraph[relatedID][signal.UserID] = true
		}
	}

	// Find connected components (coordinated networks)
	visited := make(map[uuid.UUID]bool)
	var networks [][]uuid.UUID

	for accountID := range accountGraph {
		if visited[accountID] {
			continue
		}

		// BFS to find connected component
		network := []uuid.UUID{}
		queue := []uuid.UUID{accountID}
		visited[accountID] = true

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			network = append(network, current)

			for neighbor := range accountGraph[current] {
				if !visited[neighbor] {
					visited[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}

		// Only include networks with 3+ accounts
		if len(network) >= 3 {
			networks = append(networks, network)
		}
	}

	return networks, nil
}
