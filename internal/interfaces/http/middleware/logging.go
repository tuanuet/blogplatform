package middleware

import (
	"time"

	"github.com/aiagent/boilerplate/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Logging returns a middleware that logs HTTP requests
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request
		logger.Info("HTTP Request", map[string]interface{}{
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"status":     c.Writer.Status(),
			"latency":    latency.String(),
			"latency_ms": latency.Milliseconds(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"errors":     c.Errors.ByType(gin.ErrorTypePrivate).String(),
		})
	}
}
