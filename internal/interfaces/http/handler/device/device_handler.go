package device

import (
	"net/http"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/notification"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeviceHandler interface {
	RegisterDevice(c *gin.Context)
	UnregisterDevice(c *gin.Context)
}

type deviceHandler struct {
	usecase notification.NotificationUseCase
}

func NewDeviceHandler(uc notification.NotificationUseCase) DeviceHandler {
	return &deviceHandler{
		usecase: uc,
	}
}

// getUserID extracts the user ID from the gin context
func getUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.UUID{}, false
	}

	uid, ok := userID.(uuid.UUID)
	return uid, ok
}

func (h *deviceHandler) RegisterDevice(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.RegisterDeviceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.usecase.RegisterDeviceToken(c.Request.Context(), userID, req); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Device registered successfully"})
}

func (h *deviceHandler) UnregisterDevice(c *gin.Context) {
	// For now, return success as the usecase doesn't have an UnregisterDeviceToken method
	// This can be implemented later when the repository supports it
	response.Success(c, http.StatusNoContent, nil)
}
