package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSessionAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSessionRepository(ctrl)
	validUUID := uuid.New()

	tests := []struct {
		name           string
		setupMock      func()
		setupRequest   func(req *http.Request)
		expectedStatus int
		expectedUserID string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().GetUserID(gomock.Any(), "valid_session").Return(validUUID.String(), nil)
			},
			setupRequest: func(req *http.Request) {
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "valid_session"})
			},
			expectedStatus: http.StatusOK,
			expectedUserID: validUUID.String(),
		},
		{
			name: "NoCookie",
			setupMock: func() {
				// No call expected
			},
			setupRequest: func(req *http.Request) {
				// No cookie
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "InvalidSession",
			setupMock: func() {
				mockRepo.EXPECT().GetUserID(gomock.Any(), "invalid_session").Return("", errors.New("not found"))
			},
			setupRequest: func(req *http.Request) {
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "invalid_session"})
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "MalformedUUID",
			setupMock: func() {
				mockRepo.EXPECT().GetUserID(gomock.Any(), "malformed_session").Return("not-a-uuid", nil)
			},
			setupRequest: func(req *http.Request) {
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "malformed_session"})
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest("GET", "/", nil)
			tt.setupRequest(req)
			c.Request = req

			// Initialize middleware
			middleware := SessionAuth(mockRepo)
			middleware(c)

			if tt.expectedStatus == http.StatusOK {
				// Verify userID and sessionID were set in context
				userID, exists := c.Get("userID")
				assert.True(t, exists, "userID should be set in context")
				assert.Equal(t, tt.expectedUserID, userID.(uuid.UUID).String())

				sessionID, exists := c.Get("sessionID")
				assert.True(t, exists, "sessionID should be set in context")
				assert.NotEmpty(t, sessionID)
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
				assert.True(t, c.IsAborted())
			}
		})
	}
}
