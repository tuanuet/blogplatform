package auth

import (
	"errors"
	"net/http"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/auth"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	VerifyEmail(c *gin.Context)
	ResendVerification(c *gin.Context)
	ForgotPassword(c *gin.Context)
	ValidateResetPasswordToken(c *gin.Context)
	ResetPassword(c *gin.Context)
	SocialLogin(c *gin.Context)
	SocialCallback(c *gin.Context)
}

type authHandler struct {
	authUseCase auth.AuthUseCase
}

func NewAuthHandler(authUseCase auth.AuthUseCase) AuthHandler {
	return &authHandler{
		authUseCase: authUseCase,
	}
}

func (h *authHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authUseCase.Register(c.Request.Context(), req)
	if err != nil {
		// In a real app, we might check error types (e.g. ErrUserAlreadyExists)
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

func (h *authHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authUseCase.Login(c.Request.Context(), req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	// Set session cookie
	secure := gin.Mode() == gin.ReleaseMode
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session_id", resp.SessionID, 3600*24, "/", "", secure, true)

	response.Success(c, http.StatusOK, resp)
}

func (h *authHandler) Logout(c *gin.Context) {
	sessionID, exists := c.Get("sessionID")
	if !exists {
		response.Unauthorized(c, "not logged in")
		return
	}

	if err := h.authUseCase.Logout(c.Request.Context(), sessionID.(string)); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Clear cookie
	secure := gin.Mode() == gin.ReleaseMode
	c.SetCookie("session_id", "", -1, "/", "", secure, true)

	response.Success(c, http.StatusOK, nil)
}

// VerifyEmail godoc
// @Summary Verify email address
// @Description Verify a user's email using the verification token sent via email
// @Tags Auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/verify-email [get]
func (h *authHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.BadRequest(c, "token is required")
		return
	}

	if err := h.authUseCase.VerifyEmail(c.Request.Context(), token); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "email verified"})
}

// ResendVerification godoc
// @Summary Resend verification email
// @Description Resend a verification email (always returns 200 to prevent account enumeration)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResendVerificationRequest true "Email"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/resend-verification [post]
func (h *authHandler) ResendVerification(c *gin.Context) {
	var req dto.ResendVerificationRequest
	_ = c.ShouldBindJSON(&req)

	if req.Email != "" {
		_ = h.authUseCase.ResendVerificationEmail(c.Request.Context(), req.Email)
	}

	response.Success(c, http.StatusOK, gin.H{
		"message": "If an account exists for this email, a verification message will be sent.",
	})
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Always returns 200 to prevent account enumeration. If the user exists, issues a time-limited token and attempts to send email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Email"
// @Success 200 {object} response.Response{data=dto.ForgotPasswordResponse}
// @Router /api/v1/auth/forgot-password [post]
func (h *authHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	_ = c.ShouldBindJSON(&req)

	if req.Email != "" {
		_ = h.authUseCase.ForgotPassword(c.Request.Context(), req.Email)
	}

	response.Success(c, http.StatusOK, dto.ForgotPasswordResponse{Message: "If an account exists for this email, a password reset link will be sent."})
}

// ValidateResetPasswordToken godoc
// @Summary Validate password reset token
// @Description Returns 200 if token is valid, 400 if invalid or expired.
// @Tags Auth
// @Accept json
// @Produce json
// @Param token query string true "Password reset token"
// @Success 200 {object} response.Response{data=dto.ResetPasswordValidateResponse}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/reset-password/validate [get]
func (h *authHandler) ValidateResetPasswordToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.BadRequest(c, "token is required")
		return
	}

	if err := h.authUseCase.ValidateResetPasswordToken(c.Request.Context(), token); err != nil {
		response.BadRequest(c, "invalid or expired reset token")
		return
	}

	response.Success(c, http.StatusOK, dto.ResetPasswordValidateResponse{Valid: true})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Validates token, updates bcrypt password_hash, invalidates token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} response.Response{data=dto.ResetPasswordResponse}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/reset-password [post]
func (h *authHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err := h.authUseCase.ResetPassword(c.Request.Context(), req.Token, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrResetTokenInvalid) || errors.Is(err, auth.ErrPasswordTooShort) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, dto.ResetPasswordResponse{Message: "password reset successful"})
}

func (h *authHandler) SocialLogin(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		response.BadRequest(c, "provider is required")
		return
	}

	authURL, err := h.authUseCase.GetSocialAuthURL(c.Request.Context(), provider)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *authHandler) SocialCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	if code == "" {
		response.BadRequest(c, "code is required")
		return
	}

	req := dto.LoginWithSocialRequest{
		Provider: provider,
		Code:     code,
	}

	resp, err := h.authUseCase.LoginWithSocial(c.Request.Context(), req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	// Set session cookie
	secure := gin.Mode() == gin.ReleaseMode
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session_id", resp.SessionID, 3600*24, "/", "", secure, true)

	response.Success(c, http.StatusOK, resp)
}
