package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterReadingHistoryRoutes(v1 *gin.RouterGroup, p Params, sessionAuth gin.HandlerFunc) {
	v1.GET("/me/history", sessionAuth, p.ReadingHistoryHandler.GetHistory)
}
