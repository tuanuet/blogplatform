package handler

import (
	"net/http"

	"github.com/aiagent/boilerplate/internal/application/usecase"
	"github.com/aiagent/boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	healthUseCase usecase.HealthUseCase
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(healthUseCase usecase.HealthUseCase) *HealthHandler {
	return &HealthHandler{
		healthUseCase: healthUseCase,
	}
}

// Check godoc
// @Summary Health Check
// @Description Check the health status of all services
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /api/v1/health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	health := h.healthUseCase.Check(c.Request.Context())
	response.Success(c, http.StatusOK, health)
}

// Ping godoc
// @Summary Ping
// @Description Simple ping endpoint for load balancer
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func (h *HealthHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
