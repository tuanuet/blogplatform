package entity_test

import (
	"testing"

	"github.com/aiagent/boilerplate/internal/domain/entity"
)

func TestUser_GetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		user     entity.User
		expected string
	}{
		{
			name: "returns display name when set",
			user: entity.User{
				Name:        "John Doe",
				DisplayName: strPtr("Johnny"),
			},
			expected: "Johnny",
		},
		{
			name: "falls back to name when display name is nil",
			user: entity.User{
				Name:        "John Doe",
				DisplayName: nil,
			},
			expected: "John Doe",
		},
		{
			name: "falls back to name when display name is empty",
			user: entity.User{
				Name:        "John Doe",
				DisplayName: strPtr(""),
			},
			expected: "John Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.GetDisplayName()
			if result != tt.expected {
				t.Errorf("GetDisplayName() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestUser_HasAvatar(t *testing.T) {
	tests := []struct {
		name     string
		user     entity.User
		expected bool
	}{
		{
			name: "returns true when avatar URL is set",
			user: entity.User{
				AvatarURL: strPtr("/uploads/avatars/test.jpg"),
			},
			expected: true,
		},
		{
			name: "returns false when avatar URL is nil",
			user: entity.User{
				AvatarURL: nil,
			},
			expected: false,
		},
		{
			name: "returns false when avatar URL is empty",
			user: entity.User{
				AvatarURL: strPtr(""),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.HasAvatar()
			if result != tt.expected {
				t.Errorf("HasAvatar() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Helper function for string pointers
func strPtr(s string) *string {
	return &s
}
