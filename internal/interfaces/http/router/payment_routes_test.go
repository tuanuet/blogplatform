package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aiagent/internal/interfaces/http/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockPaymentHandler struct{}

func (m *mockPaymentHandler) CreatePayment(c *gin.Context) {
	c.Status(http.StatusCreated)
}

type mockWebhookHandler struct{}

func (m *mockWebhookHandler) HandleSePayWebhook(c *gin.Context) {
	c.Status(http.StatusOK)
}

func TestRegisterPaymentRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")

	// Mock middleware that requires Authorization header
	sessionAuth := func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}

	router.RegisterPaymentRoutes(v1, &mockPaymentHandler{}, &mockWebhookHandler{}, sessionAuth)

	t.Run("POST /api/v1/payments should require auth", func(t *testing.T) {
		// No auth header -> 401
		req, _ := http.NewRequest("POST", "/api/v1/payments", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// With auth header -> 201 (handler called)
		reqAuth, _ := http.NewRequest("POST", "/api/v1/payments", nil)
		reqAuth.Header.Set("Authorization", "Bearer token")
		wAuth := httptest.NewRecorder()
		r.ServeHTTP(wAuth, reqAuth)
		assert.Equal(t, http.StatusCreated, wAuth.Code)
	})

	t.Run("POST /api/v1/webhooks/sepay should be public", func(t *testing.T) {
		// No auth header -> 200 (handler called directly)
		req, _ := http.NewRequest("POST", "/api/v1/webhooks/sepay", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
