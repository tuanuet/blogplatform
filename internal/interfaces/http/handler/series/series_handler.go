package series

import (
	"math"
	"net/http"

	"github.com/aiagent/internal/application/dto"
	seriesUsecase "github.com/aiagent/internal/application/usecase/series"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type seriesHandler struct {
	seriesUseCase seriesUsecase.SeriesUseCase
}

// NewSeriesHandler creates a new instance of SeriesHandler
func NewSeriesHandler(seriesUseCase seriesUsecase.SeriesUseCase) SeriesHandler {
	return &seriesHandler{
		seriesUseCase: seriesUseCase,
	}
}

// Create godoc
// @Summary Create a new series
// @Description Create a new series
// @Tags Series
// @Accept json
// @Produce json
// @Param request body dto.CreateSeriesRequest true "Series data"
// @Success 201 {object} dto.SeriesResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Security Bearer
// @Router /api/v1/series [post]
func (h *seriesHandler) Create(c *gin.Context) {
	authorID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	var req dto.CreateSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	series, err := h.seriesUseCase.CreateSeries(c.Request.Context(), authorID.(uuid.UUID), &req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, series)
}

// Update godoc
// @Summary Update a series
// @Description Update an existing series (author only)
// @Tags Series
// @Accept json
// @Produce json
// @Param id path string true "Series ID"
// @Param request body dto.UpdateSeriesRequest true "Series data"
// @Success 200 {object} dto.SeriesResponse
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/series/{id} [put]
func (h *seriesHandler) Update(c *gin.Context) {
	authorID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid series ID")
		return
	}

	var req dto.UpdateSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	series, err := h.seriesUseCase.UpdateSeries(c.Request.Context(), authorID.(uuid.UUID), id, &req)
	if err != nil {
		// Basic error matching, could be refined
		if err.Error() == "unauthorized: you are not the author of this series" {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, series)
}

// Delete godoc
// @Summary Delete a series
// @Description Delete a series (author only)
// @Tags Series
// @Accept json
// @Produce json
// @Param id path string true "Series ID"
// @Success 204
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/series/{id} [delete]
func (h *seriesHandler) Delete(c *gin.Context) {
	authorID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid series ID")
		return
	}

	if err := h.seriesUseCase.DeleteSeries(c.Request.Context(), authorID.(uuid.UUID), id); err != nil {
		if err.Error() == "unauthorized: you are not the author of this series" {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetByID godoc
// @Summary Get series by ID
// @Description Get a series by its ID
// @Tags Series
// @Accept json
// @Produce json
// @Param id path string true "Series ID"
// @Success 200 {object} dto.SeriesResponse
// @Failure 404 {object} response.Response
// @Router /api/v1/series/{id} [get]
func (h *seriesHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid series ID")
		return
	}

	series, err := h.seriesUseCase.GetSeriesByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "series not found")
		return
	}

	response.Success(c, http.StatusOK, series)
}

// GetBySlug godoc
// @Summary Get series by slug
// @Description Get a series by its slug
// @Tags Series
// @Accept json
// @Produce json
// @Param slug path string true "Series Slug"
// @Success 200 {object} dto.SeriesResponse
// @Failure 404 {object} response.Response
// @Router /api/v1/series/slug/{slug} [get]
func (h *seriesHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "slug is required")
		return
	}

	series, err := h.seriesUseCase.GetSeriesBySlug(c.Request.Context(), slug)
	if err != nil {
		response.NotFound(c, "series not found")
		return
	}

	response.Success(c, http.StatusOK, series)
}

// List godoc
// @Summary List series
// @Description List series with optional filters
// @Tags Series
// @Accept json
// @Produce json
// @Param authorId query string false "Filter by author ID"
// @Param search query string false "Search in title or description"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/series [get]
func (h *seriesHandler) List(c *gin.Context) {
	var params dto.SeriesFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}

	seriesList, total, err := h.seriesUseCase.ListSeries(c.Request.Context(), &params)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))

	response.SuccessWithMeta(c, seriesList, &response.Meta{
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// AddBlog godoc
// @Summary Add blog to series
// @Description Add a blog to a series (author only)
// @Tags Series
// @Accept json
// @Produce json
// @Param id path string true "Series ID"
// @Param request body dto.AddBlogToSeriesRequest true "Blog ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/series/{id}/blogs [post]
func (h *seriesHandler) AddBlog(c *gin.Context) {
	authorID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	seriesID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid series ID")
		return
	}

	var req dto.AddBlogToSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.seriesUseCase.AddBlogToSeries(c.Request.Context(), authorID.(uuid.UUID), seriesID, req.BlogID); err != nil {
		if err.Error() == "unauthorized: you are not the author of this series" {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "blog added to series"})
}

// RemoveBlog godoc
// @Summary Remove blog from series
// @Description Remove a blog from a series (author only)
// @Tags Series
// @Accept json
// @Produce json
// @Param id path string true "Series ID"
// @Param blogId path string true "Blog ID"
// @Success 204
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Security Bearer
// @Router /api/v1/series/{id}/blogs/{blogId} [delete]
func (h *seriesHandler) RemoveBlog(c *gin.Context) {
	authorID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	seriesID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid series ID")
		return
	}

	blogID, err := uuid.Parse(c.Param("blogId"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	if err := h.seriesUseCase.RemoveBlogFromSeries(c.Request.Context(), authorID.(uuid.UUID), seriesID, blogID); err != nil {
		if err.Error() == "unauthorized: you are not the author of this series" {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetHighlightedSeries godoc
// @Summary Get highlighted series
// @Description Get a list of highlighted series
// @Tags Series
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/series/highlighted [get]
func (h *seriesHandler) GetHighlightedSeries(c *gin.Context) {
	series, err := h.seriesUseCase.GetHighlightedSeries(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, series)
}
