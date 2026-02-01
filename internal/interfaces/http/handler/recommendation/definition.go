package recommendation

import "github.com/gin-gonic/gin"

type RecommendationHandler interface {
	GetPopularTags(c *gin.Context)
	UpdateInterests(c *gin.Context)
	GetPersonalizedFeed(c *gin.Context)
	GetRelatedBlogs(c *gin.Context)
}
