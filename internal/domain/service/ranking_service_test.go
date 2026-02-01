package service

import (
	"testing"
)

func TestCalculateGrowthRate(t *testing.T) {
	svc := &rankingService{
		config: DefaultRankingConfig(),
	}

	tests := []struct {
		name     string
		current  int64
		previous int64
		expected float64
	}{
		{
			name:     "Normal growth rate",
			current:  150,
			previous: 100,
			expected: 0.5, // 50% growth
		},
		{
			name:     "No previous data",
			current:  50,
			previous: 0,
			expected: 0.5, // 50/100 (normalized by min threshold)
		},
		{
			name:     "Small follower count - absolute growth",
			current:  60,
			previous: 50,
			expected: 0.1, // (60-50)/100 = 0.1
		},
		{
			name:     "Zero growth",
			current:  100,
			previous: 100,
			expected: 0.0,
		},
		{
			name:     "Negative growth",
			current:  80,
			previous: 100,
			expected: -0.2, // -20% growth
		},
		{
			name:     "Extreme growth - capped",
			current:  10000,
			previous: 100,
			expected: 10.0, // Capped at 10x (1000%)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.calculateGrowthRate(tt.current, tt.previous)
			if result != tt.expected {
				t.Errorf("calculateGrowthRate(%d, %d) = %f, want %f",
					tt.current, tt.previous, result, tt.expected)
			}
		})
	}
}

func TestDefaultRankingConfig(t *testing.T) {
	config := DefaultRankingConfig()

	if config.FollowerGrowthWeight != 0.6 {
		t.Errorf("Expected FollowerGrowthWeight to be 0.6, got %f", config.FollowerGrowthWeight)
	}

	if config.BlogPostVelocityWeight != 0.4 {
		t.Errorf("Expected BlogPostVelocityWeight to be 0.4, got %f", config.BlogPostVelocityWeight)
	}

	if config.TimeWindowDays != 30 {
		t.Errorf("Expected TimeWindowDays to be 30, got %d", config.TimeWindowDays)
	}

	if config.MinFollowersForRate != 100 {
		t.Errorf("Expected MinFollowersForRate to be 100, got %d", config.MinFollowersForRate)
	}
}
