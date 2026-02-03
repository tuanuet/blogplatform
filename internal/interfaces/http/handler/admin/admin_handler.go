package admin

import (
	"net/http"

	"github.com/aiagent/internal/application/usecase/admin"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
)

type AdminHandler interface {
	GetDashboardStats(c *gin.Context)
}

type adminHandler struct {
	useCase admin.AdminUseCase
}

func NewAdminHandler(useCase admin.AdminUseCase) AdminHandler {
	return &adminHandler{
		useCase: useCase,
	}
}

func (h *adminHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.useCase.GetDashboardStats(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, stats)
}
