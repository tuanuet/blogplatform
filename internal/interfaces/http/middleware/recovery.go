package middleware

import (
	"net/http"

	"github.com/aiagent/pkg/logger"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", nil, map[string]interface{}{
					"error": err,
					"path":  c.Request.URL.Path,
				})

				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Success: false,
					Error: &response.ErrorInfo{
						Code:    "INTERNAL_ERROR",
						Message: "An unexpected error occurred",
					},
				})
			}
		}()

		c.Next()
	}
}
