package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	subscriptionmocks "github.com/aiagent/internal/interfaces/http/handler/subscription/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegisterSubscriptionRoutes_TopSubscribedAuthors_IsPublicAndRateLimited(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	t.Run("public_endpoint_does_not_require_session_auth", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHandler := subscriptionmocks.NewMockSubscriptionHandler(ctrl)
		rateLimitCalls := 0
		sessionAuthCalls := 0

		rateLimit := func(c *gin.Context) {
			rateLimitCalls++
			c.Next()
		}
		sessionAuth := func(c *gin.Context) {
			sessionAuthCalls++
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		mockHandler.EXPECT().
			GetTopSubscribedAuthors(gomock.Any()).
			Do(func(c *gin.Context) {
				c.Status(http.StatusOK)
			}).
			Times(1)

		r := gin.New()
		v1 := r.Group("/api/v1")
		RegisterSubscriptionRoutes(v1, Params{SubscriptionHandler: mockHandler}, sessionAuth, rateLimit)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/authors/top-subscribed", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 1, rateLimitCalls)
		assert.Equal(t, 0, sessionAuthCalls)
	})

	t.Run("rate_limit_middleware_can_block_request", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHandler := subscriptionmocks.NewMockSubscriptionHandler(ctrl)
		rateLimitCalls := 0

		rateLimit := func(c *gin.Context) {
			rateLimitCalls++
			c.AbortWithStatus(http.StatusTooManyRequests)
		}
		sessionAuth := func(c *gin.Context) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		mockHandler.EXPECT().GetTopSubscribedAuthors(gomock.Any()).Times(0)

		r := gin.New()
		v1 := r.Group("/api/v1")
		RegisterSubscriptionRoutes(v1, Params{SubscriptionHandler: mockHandler}, sessionAuth, rateLimit)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/authors/top-subscribed", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Equal(t, 1, rateLimitCalls)
	})
}
