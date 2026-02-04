package entity_test

import (
	"testing"

	"github.com/aiagent/internal/domain/entity"
)

// TestSubscriptionTier_Level validates hierarchy (FREE=0, BRONZE=1, SILVER=2, GOLD=3)
func TestSubscriptionTier_Level(t *testing.T) {
	tests := []struct {
		name     string
		tier     entity.SubscriptionTier
		expected int
	}{
		{
			name:     "FREE tier has level 0",
			tier:     entity.TierFree,
			expected: 0,
		},
		{
			name:     "BRONZE tier has level 1",
			tier:     entity.TierBronze,
			expected: 1,
		},
		{
			name:     "SILVER tier has level 2",
			tier:     entity.TierSilver,
			expected: 2,
		},
		{
			name:     "GOLD tier has level 3",
			tier:     entity.TierGold,
			expected: 3,
		},
		{
			name:     "invalid tier defaults to level 0",
			tier:     entity.SubscriptionTier("INVALID"),
			expected: 0,
		},
		{
			name:     "empty tier defaults to level 0",
			tier:     entity.SubscriptionTier(""),
			expected: 0,
		},
		{
			name:     "lowercase tier defaults to level 0",
			tier:     entity.SubscriptionTier("free"),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tier.Level()
			if result != tt.expected {
				t.Errorf("SubscriptionTier(%q).Level() = %d, expected %d", tt.tier, result, tt.expected)
			}
		})
	}
}

// TestSubscriptionTier_IsValid validates valid and invalid tier values
func TestSubscriptionTier_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		tier     entity.SubscriptionTier
		expected bool
	}{
		{
			name:     "FREE tier is valid",
			tier:     entity.TierFree,
			expected: true,
		},
		{
			name:     "BRONZE tier is valid",
			tier:     entity.TierBronze,
			expected: true,
		},
		{
			name:     "SILVER tier is valid",
			tier:     entity.TierSilver,
			expected: true,
		},
		{
			name:     "GOLD tier is valid",
			tier:     entity.TierGold,
			expected: true,
		},
		{
			name:     "invalid tier is not valid",
			tier:     entity.SubscriptionTier("INVALID"),
			expected: false,
		},
		{
			name:     "empty tier is not valid",
			tier:     entity.SubscriptionTier(""),
			expected: false,
		},
		{
			name:     "lowercase tier is not valid",
			tier:     entity.SubscriptionTier("free"),
			expected: false,
		},
		{
			name:     "mixed case tier is not valid",
			tier:     entity.SubscriptionTier("Free"),
			expected: false,
		},
		{
			name:     "partial tier name is not valid",
			tier:     entity.SubscriptionTier("BRON"),
			expected: false,
		},
		{
			name:     "tier with extra characters is not valid",
			tier:     entity.SubscriptionTier("BRONZE_X"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tier.IsValid()
			if result != tt.expected {
				t.Errorf("SubscriptionTier(%q).IsValid() = %v, expected %v", tt.tier, result, tt.expected)
			}
		})
	}
}

// TestSubscriptionTier_String returns correct string representation
func TestSubscriptionTier_String(t *testing.T) {
	tests := []struct {
		name     string
		tier     entity.SubscriptionTier
		expected string
	}{
		{
			name:     "FREE tier string is FREE",
			tier:     entity.TierFree,
			expected: "FREE",
		},
		{
			name:     "BRONZE tier string is BRONZE",
			tier:     entity.TierBronze,
			expected: "BRONZE",
		},
		{
			name:     "SILVER tier string is SILVER",
			tier:     entity.TierSilver,
			expected: "SILVER",
		},
		{
			name:     "GOLD tier string is GOLD",
			tier:     entity.TierGold,
			expected: "GOLD",
		},
		{
			name:     "invalid tier returns its string value",
			tier:     entity.SubscriptionTier("INVALID"),
			expected: "INVALID",
		},
		{
			name:     "empty tier returns empty string",
			tier:     entity.SubscriptionTier(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tier.String()
			if result != tt.expected {
				t.Errorf("SubscriptionTier(%q).String() = %q, expected %q", tt.tier, result, tt.expected)
			}
		})
	}
}

// TestSubscriptionTier_TierComparison compares tier levels using Level() method
func TestSubscriptionTier_TierComparison(t *testing.T) {
	tests := []struct {
		name          string
		tier1         entity.SubscriptionTier
		tier2         entity.SubscriptionTier
		shouldBeEqual bool
		shouldBeLess  bool
	}{
		{
			name:          "same tiers are equal",
			tier1:         entity.TierBronze,
			tier2:         entity.TierBronze,
			shouldBeEqual: true,
		},
		{
			name:         "FREE is less than BRONZE",
			tier1:        entity.TierFree,
			tier2:        entity.TierBronze,
			shouldBeLess: true,
		},
		{
			name:         "BRONZE is less than SILVER",
			tier1:        entity.TierBronze,
			tier2:        entity.TierSilver,
			shouldBeLess: true,
		},
		{
			name:         "SILVER is less than GOLD",
			tier1:        entity.TierSilver,
			tier2:        entity.TierGold,
			shouldBeLess: true,
		},
		{
			name:          "GOLD is greater than SILVER",
			tier1:         entity.TierGold,
			tier2:         entity.TierSilver,
			shouldBeEqual: false,
			shouldBeLess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level1 := tt.tier1.Level()
			level2 := tt.tier2.Level()

			if tt.shouldBeEqual {
				if level1 != level2 {
					t.Errorf("Tier %q (level %d) should equal tier %q (level %d)",
						tt.tier1, level1, tt.tier2, level2)
				}
			}

			if tt.shouldBeLess {
				if level1 >= level2 {
					t.Errorf("Tier %q (level %d) should be less than tier %q (level %d)",
						tt.tier1, level1, tt.tier2, level2)
				}
			}
		})
	}
}

// TestSubscriptionTier_Constants verifies tier constants are correctly defined
func TestSubscriptionTier_Constants(t *testing.T) {
	tests := []struct {
		name  string
		tier  entity.SubscriptionTier
		value string
	}{
		{
			name:  "TierFree constant value",
			tier:  entity.TierFree,
			value: "FREE",
		},
		{
			name:  "TierBronze constant value",
			tier:  entity.TierBronze,
			value: "BRONZE",
		},
		{
			name:  "TierSilver constant value",
			tier:  entity.TierSilver,
			value: "SILVER",
		},
		{
			name:  "TierGold constant value",
			tier:  entity.TierGold,
			value: "GOLD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.tier) != tt.value {
				t.Errorf("Tier constant value = %q, expected %q", string(tt.tier), tt.value)
			}
		})
	}
}

// TestSubscriptionTier_LevelProgression verifies tier levels form a proper progression
func TestSubscriptionTier_LevelProgression(t *testing.T) {
	tiers := []entity.SubscriptionTier{
		entity.TierFree,
		entity.TierBronze,
		entity.TierSilver,
		entity.TierGold,
	}

	for i := 1; i < len(tiers); i++ {
		prevLevel := tiers[i-1].Level()
		currLevel := tiers[i].Level()

		if currLevel != prevLevel+1 {
			t.Errorf("Tier %q (level %d) should be one level higher than %q (level %d)",
				tiers[i], currLevel, tiers[i-1], prevLevel)
		}
	}
}

// TestSubscriptionTier_AllTiersAreValid verifies all defined tier constants are valid
func TestSubscriptionTier_AllTiersAreValid(t *testing.T) {
	tiers := []entity.SubscriptionTier{
		entity.TierFree,
		entity.TierBronze,
		entity.TierSilver,
		entity.TierGold,
	}

	for _, tier := range tiers {
		if !tier.IsValid() {
			t.Errorf("Tier constant %q is not valid", tier)
		}
	}
}
