package notification_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/notification/mocks"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/interfaces/http/handler/notification"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTest(t *testing.T) (*gomock.Controller, *mocks.MockNotificationUseCase, *gin.Engine, notification.NotificationHandler) {
	ctrl := gomock.NewController(t)
	mockUC := mocks.NewMockNotificationUseCase(ctrl)
	h := notification.NewNotificationHandler(mockUC)

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

func TestNotificationHandler_List(t *testing.T) {
	ctrl, mockUC, r, h := setupTest(t)
	defer ctrl.Finish()

	r.GET("/notifications", h.List)

	t.Run("success", func(t *testing.T) {
		expectedResult := &repository.PaginatedResult[dto.NotificationResponse]{
			Data: []dto.NotificationResponse{
				{ID: uuid.New(), Title: "Test Notification"},
			},
			Total: 1,
		}

		mockUC.EXPECT().
			List(gomock.Any(), gomock.Any(), 1, 10).
			Return(expectedResult, nil)

		req, _ := http.NewRequest("GET", "/notifications?page=1&page_size=10", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, true, resp["success"])
	})
}

func TestNotificationHandler_GetUnreadCount(t *testing.T) {
	ctrl, mockUC, r, h := setupTest(t)
	defer ctrl.Finish()

	r.GET("/notifications/unread-count", h.GetUnreadCount)

	t.Run("success", func(t *testing.T) {
		mockUC.EXPECT().
			GetUnreadCount(gomock.Any(), gomock.Any()).
			Return(5, nil)

		req, _ := http.NewRequest("GET", "/notifications/unread-count", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, float64(5), data["count"])
	})
}

func TestNotificationHandler_MarkAsRead(t *testing.T) {
	ctrl, mockUC, r, h := setupTest(t)
	defer ctrl.Finish()

	r.PUT("/notifications/:id/read", h.MarkAsRead)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		mockUC.EXPECT().
			MarkAsRead(gomock.Any(), gomock.Any(), id).
			Return(nil)

		req, _ := http.NewRequest("PUT", "/notifications/"+id.String()+"/read", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestNotificationHandler_MarkAllAsRead(t *testing.T) {
	ctrl, mockUC, r, h := setupTest(t)
	defer ctrl.Finish()

	r.PUT("/notifications/read-all", h.MarkAllAsRead)

	t.Run("success", func(t *testing.T) {
		mockUC.EXPECT().
			MarkAllAsRead(gomock.Any(), gomock.Any()).
			Return(nil)

		req, _ := http.NewRequest("PUT", "/notifications/read-all", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestNotificationHandler_GetPreferences(t *testing.T) {
	ctrl, mockUC, r, h := setupTest(t)
	defer ctrl.Finish()

	r.GET("/notifications/preferences", h.GetPreferences)

	t.Run("success", func(t *testing.T) {
		expected := []dto.NotificationPreferenceResponse{
			{NotificationType: "blog_like", Channel: "in_app", Enabled: true},
		}

		mockUC.EXPECT().
			GetPreferences(gomock.Any(), gomock.Any()).
			Return(expected, nil)

		req, _ := http.NewRequest("GET", "/notifications/preferences", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestNotificationHandler_UpdatePreferences(t *testing.T) {
	ctrl, mockUC, r, h := setupTest(t)
	defer ctrl.Finish()

	r.PUT("/notifications/preferences", h.UpdatePreferences)

	t.Run("success", func(t *testing.T) {
		input := dto.UpdatePreferencesRequest{
			Preferences: []dto.NotificationPreferenceItem{
				{NotificationType: "blog_like", Channel: "in_app", Enabled: false},
			},
		}
		body, _ := json.Marshal(input)

		mockUC.EXPECT().
			UpdatePreferences(gomock.Any(), gomock.Any(), gomock.Eq(input)).
			Return(nil)

		req, _ := http.NewRequest("PUT", "/notifications/preferences", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestNotificationHandler_RegisterDeviceToken(t *testing.T) {
	ctrl, mockUC, r, h := setupTest(t)
	defer ctrl.Finish()

	r.POST("/devices/token", h.RegisterDeviceToken)

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
	})
}
