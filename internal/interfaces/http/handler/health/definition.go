package health

import "github.com/gin-gonic/gin"

type HealthHandler interface {
	Check(c *gin.Context)
	Ping(c *gin.Context)
}
