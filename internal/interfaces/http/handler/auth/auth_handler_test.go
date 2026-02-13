package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aiagent/internal/application/dto"
	authUC "github.com/aiagent/internal/application/usecase/auth"
	"github.com/aiagent/internal/application/usecase/auth/mocks"
	"github.com/aiagent/internal/interfaces/http/handler/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupRouter() (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	w := httptest.NewRecorder()
	return r, w
}

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAuthUseCase(ctrl)
	handler := auth.NewAuthHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/register", handler.Register)

		reqBody := dto.RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		expectedResp := &dto.AuthResponse{
			UserID: uuid.New(),
			Email:  reqBody.Email,
			Name:   reqBody.Name,
		}

		mockUseCase.EXPECT().
			Register(gomock.Any(), reqBody).
			Return(expectedResp, nil)

		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, true, response["success"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, reqBody.Email, data["email"])
	})

	t.Run("invalid_request", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/register", handler.Register)

		// Invalid email
		reqBody := dto.RegisterRequest{
			Name:     "Test User",
			Email:    "invalid-email",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase_error", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/register", handler.Register)

		reqBody := dto.RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().
			Register(gomock.Any(), reqBody).
			Return(nil, errors.New("email already exists"))

		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAuthUseCase(ctrl)
	handler := auth.NewAuthHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/login", handler.Login)

		reqBody := dto.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		expectedResp := &dto.AuthResponse{
			SessionID: "session-123",
			UserID:    uuid.New(),
			Email:     reqBody.Email,
			Name:      "Test User",
		}

		mockUseCase.EXPECT().
			Login(gomock.Any(), reqBody).
			Return(expectedResp, nil)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Check for cookie
		cookies := w.Header().Values("Set-Cookie")
		assert.NotEmpty(t, cookies)
		assert.Contains(t, cookies[0], "session_id=session-123")
		assert.Contains(t, cookies[0], "HttpOnly")

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, true, response["success"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "session-123", data["sessionId"])
	})

	t.Run("invalid_credentials", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/login", handler.Login)

		reqBody := dto.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().
			Login(gomock.Any(), reqBody).
			Return(nil, errors.New("invalid credentials"))

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAuthUseCase(ctrl)
	handler := auth.NewAuthHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		r, w := setupRouter()
		// Manual context setting to simulate middleware
		r.POST("/logout", func(c *gin.Context) {
			c.Set("sessionID", "session-123")
			handler.Logout(c)
		})

		mockUseCase.EXPECT().
			Logout(gomock.Any(), "session-123").
			Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/logout", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Check if cookie is cleared
		cookies := w.Header().Values("Set-Cookie")
		assert.NotEmpty(t, cookies)
		assert.Contains(t, cookies[0], "session_id=;")
	})

	t.Run("missing_session", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/logout", handler.Logout)

		req, _ := http.NewRequest(http.MethodPost, "/logout", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthHandler_ForgotPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAuthUseCase(ctrl)
	handler := auth.NewAuthHandler(mockUseCase)

	t.Run("always_200_even_on_bad_json", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/forgot-password", handler.ForgotPassword)

		mockUseCase.EXPECT().ForgotPassword(gomock.Any(), gomock.Any()).Times(0)

		req, _ := http.NewRequest(http.MethodPost, "/forgot-password", bytes.NewBufferString("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("calls_usecase_when_email_present", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/forgot-password", handler.ForgotPassword)

		reqBody := dto.ForgotPasswordRequest{Email: "test@example.com"}
		jsonBody, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().ForgotPassword(gomock.Any(), reqBody.Email).Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/forgot-password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthHandler_ValidateResetPasswordToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAuthUseCase(ctrl)
	handler := auth.NewAuthHandler(mockUseCase)

	t.Run("valid_token_200", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/reset-password/validate", handler.ValidateResetPasswordToken)

		token := "reset:token"
		mockUseCase.EXPECT().ValidateResetPasswordToken(gomock.Any(), token).Return(nil)

		req, _ := http.NewRequest(http.MethodGet, "/reset-password/validate?token="+token, nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid_token_400", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/reset-password/validate", handler.ValidateResetPasswordToken)

		token := "reset:bad"
		mockUseCase.EXPECT().ValidateResetPasswordToken(gomock.Any(), token).Return(authUC.ErrResetTokenInvalid)

		req, _ := http.NewRequest(http.MethodGet, "/reset-password/validate?token="+token, nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_ResetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAuthUseCase(ctrl)
	handler := auth.NewAuthHandler(mockUseCase)

	t.Run("success_200", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/reset-password", handler.ResetPassword)

		reqBody := dto.ResetPasswordRequest{Token: "reset:token", Password: "password123"}
		jsonBody, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().ResetPassword(gomock.Any(), reqBody.Token, reqBody.Password).Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/reset-password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid_token_400", func(t *testing.T) {
		r, w := setupRouter()
		r.POST("/reset-password", handler.ResetPassword)

		reqBody := dto.ResetPasswordRequest{Token: "reset:bad", Password: "password123"}
		jsonBody, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().ResetPassword(gomock.Any(), reqBody.Token, reqBody.Password).Return(authUC.ErrResetTokenInvalid)

		req, _ := http.NewRequest(http.MethodPost, "/reset-password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_SocialLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAuthUseCase(ctrl)
	handler := auth.NewAuthHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/auth/:provider", handler.SocialLogin)

		provider := "google"
		authURL := "https://accounts.google.com/o/oauth2/v2/auth?client_id=...&redirect_uri=..."

		mockUseCase.EXPECT().
			GetSocialAuthURL(gomock.Any(), provider).
			Return(authURL, nil)

		req, _ := http.NewRequest(http.MethodGet, "/auth/"+provider, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Equal(t, authURL, w.Header().Get("Location"))
	})
}

func TestAuthHandler_SocialCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAuthUseCase(ctrl)
	handler := auth.NewAuthHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		r, w := setupRouter()
		r.GET("/auth/:provider/callback", handler.SocialCallback)

		provider := "google"
		code := "auth-code"
		expectedResp := &dto.AuthResponse{
			SessionID: "session-123",
			UserID:    uuid.New(),
			Email:     "user@example.com",
			Name:      "Social User",
		}

		mockUseCase.EXPECT().
			LoginWithSocial(gomock.Any(), dto.LoginWithSocialRequest{
				Provider: provider,
				Code:     code,
			}).
			Return(expectedResp, nil)

		req, _ := http.NewRequest(http.MethodGet, "/auth/"+provider+"/callback?code="+code, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Check for cookie
		cookies := w.Header().Values("Set-Cookie")
		assert.NotEmpty(t, cookies)
		assert.Contains(t, cookies[0], "session_id=session-123")
		assert.Contains(t, cookies[0], "HttpOnly")

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, true, response["success"])
	})
}
