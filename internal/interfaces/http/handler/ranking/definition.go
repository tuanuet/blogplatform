package ranking

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import "github.com/gin-gonic/gin"

type RankingHandler interface {
	GetTrending(c *gin.Context)
	GetTop(c *gin.Context)
	GetUserRanking(c *gin.Context)
	RecalculateScores(c *gin.Context)
}
