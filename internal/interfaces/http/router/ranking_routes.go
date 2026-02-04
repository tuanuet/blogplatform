package router

import (
	"github.com/aiagent/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRankingRoutes(v1 *gin.RouterGroup, p Params, auth *middleware.Authorization, sessionAuth gin.HandlerFunc) {
	rankings := v1.Group("/rankings")
	{
		rankings.GET("/trending", p.RankingHandler.GetTrending)
		rankings.GET("/top", p.RankingHandler.GetTop)
		rankings.GET("/users/:userId", p.RankingHandler.GetUserRanking)
		rankings.POST("/recalculate", sessionAuth, auth.RequireAdmin("rankings"), p.RankingHandler.RecalculateScores)
	}
}
