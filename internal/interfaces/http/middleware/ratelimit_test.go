package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRateLimit(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("could not start miniredis: %v", err)
	}
	defer s.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Limit: 5 requests per minute
	r.Use(RateLimit(rdb, 5, time.Minute))

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	t.Run("should allow requests under limit", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "127.0.0.1:1234"
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("should block requests over limit", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:1234"
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("should allow requests after window expires", func(t *testing.T) {
		s.FastForward(time.Minute + time.Second)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:1234"
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
