package recommendation

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import "github.com/gin-gonic/gin"

type RecommendationHandler interface {
	GetPopularTags(c *gin.Context)
	UpdateInterests(c *gin.Context)
	GetPersonalizedFeed(c *gin.Context)
	GetRelatedBlogs(c *gin.Context)
}
