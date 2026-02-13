package middleware

import (
	"net/url"
	"strings"
	"time"

	"github.com/aiagent/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Logging returns a middleware that logs HTTP requests
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := redactQuery(c.Request.URL.RawQuery)

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

func redactQuery(rawQuery string) string {
	if rawQuery == "" {
		return ""
	}

	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		// As a fallback, don't log raw query when it might contain secrets.
		if strings.Contains(rawQuery, "token=") || strings.Contains(rawQuery, "code=") {
			return "[REDACTED]"
		}
		return rawQuery
	}

	for _, key := range []string{"token", "code"} {
		if values.Has(key) {
			values.Set(key, "[REDACTED]")
		}
	}

	return values.Encode()
}
