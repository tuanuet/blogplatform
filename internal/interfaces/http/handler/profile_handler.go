package handler

import (
	"errors"
	"net/http"

	"github.com/aiagent/boilerplate/internal/application/dto"
	"github.com/aiagent/boilerplate/internal/application/usecase"
	"github.com/aiagent/boilerplate/pkg/response"
	"github.com/aiagent/boilerplate/pkg/validator"
	"github.com/gin-gonic/gin"
	validatorLib "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ProfileHandler handles profile-related HTTP requests
type ProfileHandler struct {
	profileUseCase usecase.ProfileUseCase
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(profileUseCase usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{
		profileUseCase: profileUseCase,
	}
}

// GetMyProfile godoc
// @Summary Get my profile
// @Description Get the current user's profile
// @Tags Profile
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=dto.ProfileResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/profile [get]
func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Authentication required")
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "Invalid user ID")
		return
	}

	profile, err := h.profileUseCase.GetProfile(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to get profile")
		return
	}

	response.Success(c, http.StatusOK, profile)
}

// UpdateMyProfile godoc
// @Summary Update my profile
// @Description Update the current user's profile
// @Tags Profile
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.UpdateProfileRequest true "Profile update request"
// @Success 200 {object} response.Response{data=dto.ProfileResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/profile [put]
func (h *ProfileHandler) UpdateMyProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Authentication required")
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "Invalid user ID")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Validate using the validator library directly
	validate := validatorLib.New()
	if err := validate.Struct(req); err != nil {
		errs := validator.FormatValidationErrors(err)
		if len(errs) > 0 {
			response.ValidationError(c, errs[0].Message)
			return
		}
		response.ValidationError(c, "Validation failed")
		return
	}

	profile, err := h.profileUseCase.UpdateProfile(c.Request.Context(), uid, req)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to update profile")
		return
	}

	response.Success(c, http.StatusOK, profile)
}

// UploadAvatar godoc
// @Summary Upload avatar
// @Description Upload a new avatar image for the current user
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param avatar formData file true "Avatar image file (max 5MB, jpg/png/gif/webp)"
// @Success 200 {object} response.Response{data=dto.AvatarUploadResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/profile/avatar [post]
func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Authentication required")
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "Invalid user ID")
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		response.BadRequest(c, "Avatar file is required")
		return
	}

	result, err := h.profileUseCase.UploadAvatar(c.Request.Context(), uid, file)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrFileTooLarge):
			response.BadRequest(c, "File too large, max 5MB allowed")
		case errors.Is(err, usecase.ErrInvalidFileType):
			response.BadRequest(c, "Invalid file type, allowed: jpg, jpeg, png, gif, webp")
		case errors.Is(err, usecase.ErrUserNotFound):
			response.NotFound(c, "User not found")
		default:
			response.InternalServerError(c, "Failed to upload avatar")
		}
		return
	}

	response.Success(c, http.StatusOK, result)
}

// GetPublicProfile godoc
// @Summary Get public profile
// @Description Get a user's public profile by ID
// @Tags Profile
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=dto.PublicProfileResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id}/profile [get]
func (h *ProfileHandler) GetPublicProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	profile, err := h.profileUseCase.GetPublicProfile(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to get profile")
		return
	}

	response.Success(c, http.StatusOK, profile)
}
