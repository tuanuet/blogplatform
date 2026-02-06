package profile

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	domainService "github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
)

// Use case errors
var (
	ErrUserNotFound      = domainService.ErrUserNotFound
	ErrInvalidFileType   = errors.New("invalid file type, allowed: jpg, jpeg, png, gif, webp")
	ErrFileTooLarge      = errors.New("file too large, max 5MB allowed")
	ErrUploadFailed      = errors.New("file upload failed")
	ErrInvalidAvatarPath = errors.New("invalid avatar path")
)

// Allowed avatar file types
var allowedAvatarTypes = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

const (
	maxAvatarSize = 5 * 1024 * 1024 // 5MB
	avatarDir     = "uploads/avatars"
)

// ProfileUseCase handles profile-related application logic
type ProfileUseCase interface {
	// GetProfile retrieves user's own profile
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.ProfileResponse, error)

	// GetPublicProfile retrieves a user's public profile
	GetPublicProfile(ctx context.Context, userID uuid.UUID) (*dto.PublicProfileResponse, error)

	// UpdateProfile updates user's profile
	UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.ProfileResponse, error)

	// UploadAvatar uploads and updates user's avatar
	UploadAvatar(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (*dto.AvatarUploadResponse, error)
}

type profileUseCase struct {
	userSvc domainService.UserService
}

// NewProfileUseCase creates a new profile use case
func NewProfileUseCase(userSvc domainService.UserService) ProfileUseCase {
	return &profileUseCase{
		userSvc: userSvc,
	}
}

func (uc *profileUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.ProfileResponse, error) {
	user, err := uc.userSvc.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return uc.toProfileResponse(user), nil
}

func (uc *profileUseCase) GetPublicProfile(ctx context.Context, userID uuid.UUID) (*dto.PublicProfileResponse, error) {
	user, err := uc.userSvc.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return uc.toPublicProfileResponse(user), nil
}

func (uc *profileUseCase) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	updates := make(map[string]interface{})

	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.Website != nil {
		updates["website"] = *req.Website
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.TwitterHandle != nil {
		handle := strings.TrimPrefix(*req.TwitterHandle, "@")
		updates["twitter_handle"] = handle
	}
	if req.GithubHandle != nil {
		updates["github_handle"] = *req.GithubHandle
	}
	if req.LinkedinURL != nil {
		updates["linkedin_url"] = *req.LinkedinURL
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.Birthday != nil {
		birthday, err := time.Parse("2006-01-02", *req.Birthday)
		if err != nil {
			return nil, fmt.Errorf("invalid birthday format: %w", err)
		}
		if birthday.After(time.Now()) {
			return nil, errors.New("birthday cannot be in the future")
		}
		updates["birthday"] = birthday
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if err := uc.userSvc.UpdateUser(ctx, userID, updates); err != nil {
		return nil, err
	}

	// Fetch updated user
	user, err := uc.userSvc.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return uc.toProfileResponse(user), nil
}

func (uc *profileUseCase) UploadAvatar(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (*dto.AvatarUploadResponse, error) {
	// Validation
	if file.Size > maxAvatarSize {
		return nil, ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedAvatarTypes[ext] {
		return nil, ErrInvalidFileType
	}

	// Save file logic
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	filename := fmt.Sprintf("%s_%d%s", userID.String(), time.Now().UnixNano(), ext)
	filePath := filepath.Join(avatarDir, filename)

	// Path traversal check
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, ErrInvalidAvatarPath
	}
	absAvatarDir, err := filepath.Abs(avatarDir)
	if err != nil {
		return nil, ErrInvalidAvatarPath
	}
	if !strings.HasPrefix(absPath, absAvatarDir) {
		return nil, ErrInvalidAvatarPath
	}

	// Save file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	// Update DB via Domain Service
	avatarURL := "/" + filePath
	if err := uc.userSvc.UpdateAvatarURL(ctx, userID, avatarURL); err != nil {
		os.Remove(filePath) // Cleanup
		return nil, err
	}

	return &dto.AvatarUploadResponse{
		AvatarURL: avatarURL,
	}, nil
}

// Helpers
func (uc *profileUseCase) toProfileResponse(user *entity.User) *dto.ProfileResponse {
	resp := &dto.ProfileResponse{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		DisplayName: user.GetDisplayName(),
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
	}

	if user.Bio != nil {
		resp.Bio = *user.Bio
	}
	if user.AvatarURL != nil {
		resp.AvatarURL = *user.AvatarURL
	}
	if user.Website != nil {
		resp.Website = *user.Website
	}
	if user.Location != nil {
		resp.Location = *user.Location
	}
	if user.TwitterHandle != nil {
		resp.TwitterHandle = *user.TwitterHandle
	}
	if user.GithubHandle != nil {
		resp.GithubHandle = *user.GithubHandle
	}
	if user.LinkedinURL != nil {
		resp.LinkedinURL = *user.LinkedinURL
	}
	if user.Gender != nil {
		resp.Gender = *user.Gender
	}
	if user.Birthday != nil {
		resp.Birthday = user.Birthday.Format("2006-01-02")
	}
	if user.Description != nil {
		resp.Description = *user.Description
	}

	return resp
}

func (uc *profileUseCase) toPublicProfileResponse(user *entity.User) *dto.PublicProfileResponse {
	resp := &dto.PublicProfileResponse{
		ID:          user.ID,
		DisplayName: user.GetDisplayName(),
	}

	if user.Bio != nil {
		resp.Bio = *user.Bio
	}
	if user.AvatarURL != nil {
		resp.AvatarURL = *user.AvatarURL
	}
	if user.Website != nil {
		resp.Website = *user.Website
	}
	if user.Location != nil {
		resp.Location = *user.Location
	}
	if user.TwitterHandle != nil {
		resp.TwitterHandle = *user.TwitterHandle
	}
	if user.GithubHandle != nil {
		resp.GithubHandle = *user.GithubHandle
	}
	if user.LinkedinURL != nil {
		resp.LinkedinURL = *user.LinkedinURL
	}
	if user.Description != nil {
		resp.Description = *user.Description
	}

	return resp
}
