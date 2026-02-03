package device_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/notification/mocks"
	"github.com/aiagent/internal/interfaces/http/handler/device"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTest(t *testing.T) (*gomock.Controller, *mocks.MockNotificationUseCase, *gin.Engine, device.DeviceHandler) {
	ctrl := gomock.NewController(t)
	mockUC := mocks.NewMockNotificationUseCase(ctrl)
	h := device.NewDeviceHandler(mockUC)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Mock middleware to set user ID in context
	userID := uuid.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})

	return ctrl, mockUC, r, h
}

func TestDeviceHandler_RegisterDevice(t *testing.T) {
	ctrl, mockUC, r, h := setupTest(t)
	defer ctrl.Finish()

	r.POST("/devices/token", h.RegisterDevice)

	t.Run("success", func(t *testing.T) {
		input := dto.RegisterDeviceTokenRequest{
			DeviceToken: "some-token",
			Platform:    "ios",
		}
		body, _ := json.Marshal(input)

		mockUC.EXPECT().
			RegisterDeviceToken(gomock.Any(), gomock.Any(), gomock.Eq(input)).
			Return(nil)

		req, _ := http.NewRequest("POST", "/devices/token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, true, resp["success"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/devices/token", bytes.NewBuffer([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		input := dto.RegisterDeviceTokenRequest{
			DeviceToken: "",
			Platform:    "",
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/devices/token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeviceHandler_UnregisterDevice(t *testing.T) {
	ctrl, _, r, h := setupTest(t)
	defer ctrl.Finish()

	r.DELETE("/devices/:id", h.UnregisterDevice)

	t.Run("success", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/devices/some-token", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
