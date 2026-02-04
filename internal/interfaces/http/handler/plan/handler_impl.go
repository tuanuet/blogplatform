package plan

import (
	"net/http"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type planHandler struct {
	planService   service.PlanManagementService
	tagService    service.TagTierService
	accessService service.ContentAccessService
}

// NewPlanHandler creates a new PlanHandler instance
func NewPlanHandler(
	planService service.PlanManagementService,
	tagService service.TagTierService,
	accessService service.ContentAccessService,
) PlanHandler {
	return &planHandler{
		planService:   planService,
		tagService:    tagService,
		accessService: accessService,
	}
}

// UpsertPlans godoc
// @Summary Create or update subscription plans
// @Tags Plans
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.UpsertPlansRequest true "Plans to upsert"
// @Success 200 {object} dto.UpsertPlansResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/authors/me/plans [post]
func (h *planHandler) UpsertPlans(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "invalid user ID")
		return
	}

	var req dto.UpsertPlansRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Manual validation for tier values (Gin's binding doesn't trigger validate tags)
	validTiers := map[string]bool{"BRONZE": true, "SILVER": true, "GOLD": true}
	for _, plan := range req.Plans {
		if !validTiers[plan.Tier] {
			response.BadRequest(c, "invalid tier: must be BRONZE, SILVER, or GOLD")
			return
		}
		// Validate price is non-negative
		if plan.Price.IsNegative() {
			response.BadRequest(c, "price must be greater than or equal to 0")
			return
		}
	}

	// Convert DTO to service DTO
	plans := make([]service.CreatePlanDTO, len(req.Plans))
	for i, p := range req.Plans {
		plans[i] = service.CreatePlanDTO{
			Tier:         entity.SubscriptionTier(p.Tier),
			Price:        p.Price,
			Name:         p.Name,
			Description:  p.Description,
			DurationDays: 30, // Default duration
		}
	}

	createdPlans, warnings, err := h.planService.UpsertPlans(c.Request.Context(), userID, plans)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Convert entity plans to DTO response
	planResponses := make([]dto.PlanResponse, len(createdPlans))
	for i, p := range createdPlans {
		planResponses[i] = dto.PlanResponse{
			ID:           p.ID,
			Tier:         string(p.Tier),
			Price:        p.Price,
			DurationDays: p.DurationDays,
			Name:         p.Name,
			Description:  p.Description,
			IsActive:     p.IsActive,
			CreatedAt:    p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	resp := dto.UpsertPlansResponse{
		Plans:    planResponses,
		Warnings: warnings,
	}

	response.Success(c, http.StatusOK, resp)
}

// GetAuthorPlans godoc
// @Summary Get author's subscription plans
// @Description Retrieve all pricing plans for a specific author
// @Tags Plans
// @Accept json
// @Produce json
// @Param authorId path string true "Author ID"
// @Success 200 {object} dto.GetAuthorPlansResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/authors/{authorId}/plans [get]
func (h *planHandler) GetAuthorPlans(c *gin.Context) {
	authorID, err := uuid.Parse(c.Param("authorId"))
	if err != nil {
		response.BadRequest(c, "invalid author ID")
		return
	}

	plansWithTags, err := h.planService.GetAuthorPlans(c.Request.Context(), authorID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Convert service DTO to API response
	planResponses := make([]dto.PlanWithTagsResponse, len(plansWithTags))
	for i, p := range plansWithTags {
		planResponses[i] = dto.PlanWithTagsResponse{
			Tier:         string(p.Plan.Tier),
			Price:        p.Plan.Price,
			DurationDays: p.Plan.DurationDays,
			Name:         p.Plan.Name,
			Description:  p.Plan.Description,
			TagCount:     p.TagCount,
			Tags:         p.Tags,
		}
	}

	resp := dto.GetAuthorPlansResponse{
		AuthorID: authorID,
		Plans:    planResponses,
	}

	response.Success(c, http.StatusOK, resp)
}

// AssignTagToTier godoc
// @Summary Assign a tag to a subscription tier
// @Description Set which subscription tier is required to access content with this tag
// @Tags Plans
// @Accept json
// @Produce json
// @Security Bearer
// @Param tagId path string true "Tag ID"
// @Param request body dto.AssignTagTierRequest true "Required tier"
// @Success 200 {object} dto.AssignTagTierResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/authors/me/tags/{tagId}/tier [post]
func (h *planHandler) AssignTagToTier(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "invalid user ID")
		return
	}

	tagID, err := uuid.Parse(c.Param("tagId"))
	if err != nil {
		response.BadRequest(c, "invalid tag ID")
		return
	}

	var req dto.AssignTagTierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Manual validation for tier values
	validTiers := map[string]bool{"FREE": true, "BRONZE": true, "SILVER": true, "GOLD": true}
	if !validTiers[req.RequiredTier] {
		response.BadRequest(c, "invalid tier: must be FREE, BRONZE, SILVER, or GOLD")
		return
	}

	mapping, affectedCount, err := h.tagService.AssignTagToTier(
		c.Request.Context(),
		userID,
		tagID,
		entity.SubscriptionTier(req.RequiredTier),
	)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	resp := dto.AssignTagTierResponse{
		TagID:              mapping.TagID,
		TagName:            "", // Tag name would need to be fetched from tag service
		RequiredTier:       string(mapping.RequiredTier),
		AffectedBlogsCount: affectedCount,
	}

	response.Success(c, http.StatusOK, resp)
}

// UnassignTagFromTier godoc
// @Summary Remove tier requirement from a tag
// @Description Remove tier requirement, making content with this tag accessible to all subscribers
// @Tags Plans
// @Accept json
// @Produce json
// @Security Bearer
// @Param tagId path string true "Tag ID"
// @Success 200 {object} dto.UnassignTagTierResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/authors/me/tags/{tagId}/tier [delete]
func (h *planHandler) UnassignTagFromTier(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "invalid user ID")
		return
	}

	tagID, err := uuid.Parse(c.Param("tagId"))
	if err != nil {
		response.BadRequest(c, "invalid tag ID")
		return
	}

	affectedCount, err := h.tagService.UnassignTagFromTier(
		c.Request.Context(),
		userID,
		tagID,
	)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	resp := dto.UnassignTagTierResponse{
		Message:            "Tag tier assignment removed successfully",
		AffectedBlogsCount: affectedCount,
	}

	response.Success(c, http.StatusOK, resp)
}

// GetAuthorTagTiers godoc
// @Summary Get author's tag-tier mappings
// @Description Retrieve all tag-tier mappings for the current author
// @Tags Plans
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.GetTagTiersResponse
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/authors/me/tag-tiers [get]
func (h *planHandler) GetAuthorTagTiers(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "invalid user ID")
		return
	}

	mappings, err := h.tagService.GetAuthorTagTiers(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Convert service DTO to API response
	mappingResponses := make([]dto.TagTierMappingResponse, len(mappings))
	for i, m := range mappings {
		mappingResponses[i] = dto.TagTierMappingResponse{
			TagID:        m.Mapping.TagID,
			TagName:      m.TagName,
			RequiredTier: string(m.Mapping.RequiredTier),
			BlogCount:    m.BlogCount,
		}
	}

	resp := dto.GetTagTiersResponse{
		Mappings: mappingResponses,
	}

	response.Success(c, http.StatusOK, resp)
}

// CheckBlogAccess godoc
// @Summary Check blog access
// @Description Check if a user can access a specific blog post
// @Tags Plans
// @Accept json
// @Produce json
// @Param blogId path string true "Blog ID"
// @Success 200 {object} dto.CheckBlogAccessResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/blogs/{blogId}/access [get]
func (h *planHandler) CheckBlogAccess(c *gin.Context) {
	blogID, err := uuid.Parse(c.Param("blogId"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	// Get user ID if authenticated
	var userID *uuid.UUID
	if userIDVal, exists := c.Get("userID"); exists {
		if uid, ok := userIDVal.(uuid.UUID); ok {
			userID = &uid
		}
	}

	result, err := h.accessService.CheckBlogAccess(c.Request.Context(), blogID, userID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Convert upgrade options to DTO
	upgradeOptions := make([]dto.UpgradeOption, len(result.UpgradeOptions))
	for i, uo := range result.UpgradeOptions {
		price, _ := decimal.NewFromString(uo.Price)
		upgradeOptions[i] = dto.UpgradeOption{
			Tier:         string(uo.Tier),
			Price:        price,
			DurationDays: uo.DurationDays,
			PlanID:       uo.PlanID,
		}
	}

	resp := dto.CheckBlogAccessResponse{
		Accessible:     result.Accessible,
		UserTier:       string(result.UserTier),
		RequiredTier:   string(result.RequiredTier),
		Reason:         result.Reason,
		UpgradeOptions: upgradeOptions,
	}

	response.Success(c, http.StatusOK, resp)
}
