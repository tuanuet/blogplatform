package notification

import (
	"net/http"
	"strconv"

	"github.com/aiagent/internal/application/dto"
	notificationuc "github.com/aiagent/internal/application/usecase/notification"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler interface {
	List(c *gin.Context)
	GetUnreadCount(c *gin.Context)
	MarkAsRead(c *gin.Context)
	MarkAllAsRead(c *gin.Context)
	GetPreferences(c *gin.Context)
	UpdatePreferences(c *gin.Context)
	RegisterDeviceToken(c *gin.Context)
}

type notificationHandler struct {
	usecase notificationuc.NotificationUseCase
}

func NewNotificationHandler(uc notificationuc.NotificationUseCase) NotificationHandler {
	return &notificationHandler{
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

func (h *notificationHandler) List(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.usecase.List(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.SuccessWithMeta(c, result.Data, &response.Meta{
		Page:       result.Page,
		PageSize:   result.PageSize,
		Total:      result.Total,
		TotalPages: result.TotalPages,
	})
}

func (h *notificationHandler) GetUnreadCount(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	count, err := h.usecase.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"count": count})
}

func (h *notificationHandler) MarkAsRead(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	notificationIDStr := c.Param("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid notification ID")
		return
	}

	if err := h.usecase.MarkAsRead(c.Request.Context(), userID, notificationID); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

func (h *notificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	if err := h.usecase.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

func (h *notificationHandler) GetPreferences(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	preferences, err := h.usecase.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"preferences": preferences})
}

func (h *notificationHandler) UpdatePreferences(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.usecase.UpdatePreferences(c.Request.Context(), userID, req); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

func (h *notificationHandler) RegisterDeviceToken(c *gin.Context) {
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

	response.Success(c, http.StatusOK, nil)
}
