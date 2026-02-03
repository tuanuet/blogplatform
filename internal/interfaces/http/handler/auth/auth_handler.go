package auth

import (
	"net/http"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	SocialLogin(c *gin.Context)
	SocialCallback(c *gin.Context)
}

type authHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase) AuthHandler {
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
