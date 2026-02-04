package plan

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import "github.com/gin-gonic/gin"

type PlanHandler interface {
	// UpsertPlans creates or updates subscription plans for the current author
	UpsertPlans(c *gin.Context)

	// GetAuthorPlans retrieves all plans for a specific author
	GetAuthorPlans(c *gin.Context)

	// AssignTagToTier assigns a tag to a subscription tier for the current author
	AssignTagToTier(c *gin.Context)

	// UnassignTagFromTier removes tier requirement from a tag for the current author
	UnassignTagFromTier(c *gin.Context)

	// GetAuthorTagTiers retrieves all tag-tier mappings for the current author
	GetAuthorTagTiers(c *gin.Context)

	// CheckBlogAccess checks if a user can access a specific blog
	CheckBlogAccess(c *gin.Context)
}
