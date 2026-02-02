package health

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import "github.com/gin-gonic/gin"

type HealthHandler interface {
	Check(c *gin.Context)
	Ping(c *gin.Context)
}
