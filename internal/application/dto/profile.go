package dto

import "github.com/google/uuid"

// UpdateProfileRequest represents the request to update user profile
type UpdateProfileRequest struct {
	DisplayName   *string `json:"displayName" validate:"omitempty,min=3,max=50"`
	Bio           *string `json:"bio" validate:"omitempty,max=500"`
	Website       *string `json:"website" validate:"omitempty,max=255,url"`
	Location      *string `json:"location" validate:"omitempty,max=100"`
	TwitterHandle *string `json:"twitterHandle" validate:"omitempty,max=50"`
	GithubHandle  *string `json:"githubHandle" validate:"omitempty,max=50"`
	LinkedinURL   *string `json:"linkedinUrl" validate:"omitempty,max=255,url"`
}

// ProfileResponse represents the user profile response
type ProfileResponse struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	DisplayName   string    `json:"displayName"`
	Bio           string    `json:"bio,omitempty"`
	AvatarURL     string    `json:"avatarUrl,omitempty"`
	Website       string    `json:"website,omitempty"`
	Location      string    `json:"location,omitempty"`
	TwitterHandle string    `json:"twitterHandle,omitempty"`
	GithubHandle  string    `json:"githubHandle,omitempty"`
	LinkedinURL   string    `json:"linkedinUrl,omitempty"`
	CreatedAt     string    `json:"createdAt"`
}

// PublicProfileResponse represents the public user profile (limited fields)
type PublicProfileResponse struct {
	ID            uuid.UUID `json:"id"`
	DisplayName   string    `json:"displayName"`
	Bio           string    `json:"bio,omitempty"`
	AvatarURL     string    `json:"avatarUrl,omitempty"`
	Website       string    `json:"website,omitempty"`
	Location      string    `json:"location,omitempty"`
	TwitterHandle string    `json:"twitterHandle,omitempty"`
	GithubHandle  string    `json:"githubHandle,omitempty"`
	LinkedinURL   string    `json:"linkedinUrl,omitempty"`
}

// AvatarUploadResponse represents the response after avatar upload
type AvatarUploadResponse struct {
	AvatarURL string `json:"avatarUrl"`
}
