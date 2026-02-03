package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimit returns a middleware that limits requests per IP using Redis
func RateLimit(redisClient *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s", ip)

		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			// If Redis is down, we might want to allow the request or block it.
			// Usually, we allow it to avoid breaking the service if Redis fails.
			c.Next()
			return
		}

		if count == 1 {
			redisClient.Expire(ctx, key, window)
		}

		if int(count) > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}

		c.Next()
	}
}
