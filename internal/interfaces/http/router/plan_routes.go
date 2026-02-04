package router

import (
	"github.com/aiagent/internal/interfaces/http/handler/plan"
	"github.com/gin-gonic/gin"
)

// RegisterPlanRoutes registers plan-related routes
// Public routes: GET /authors/:authorId/plans, GET /blogs/:blogId/access
// Protected routes: All /authors/me/* endpoints require session authentication
func RegisterPlanRoutes(v1 *gin.RouterGroup, planH plan.PlanHandler, sessionAuth gin.HandlerFunc) {
	// Authors group
	authors := v1.Group("/authors")
	{
		// Public endpoint: Get author's subscription plans
		authors.GET("/:authorId/plans", planH.GetAuthorPlans)

		// Protected endpoints: Current author's plan management
		authorsMe := authors.Group("/me")
		authorsMe.Use(sessionAuth)
		{
			authorsMe.POST("/plans", planH.UpsertPlans)
			authorsMe.POST("/tags/:tagId/tier", planH.AssignTagToTier)
			authorsMe.DELETE("/tags/:tagId/tier", planH.UnassignTagFromTier)
			authorsMe.GET("/tag-tiers", planH.GetAuthorTagTiers)
		}
	}

	// Blogs group - public endpoint for access checking
	blogs := v1.Group("/blogs")
	{
		blogs.GET("/:blogId/access", planH.CheckBlogAccess)
	}
}
