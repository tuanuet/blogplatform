package middleware

import (
	"net/http"

	"github.com/aiagent/internal/domain/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SessionAuth creates a middleware that checks for a valid session using Redis
func SessionAuth(repo repository.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userIDStr, err := repo.GetUserID(c.Request.Context(), sessionID)
		if err != nil || userIDStr == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", userID)
		c.Set("sessionID", sessionID)
		c.Next()
	}
}
