package payment

import (
	"net/http"
	"strings"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/application/usecase/payment"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
)

// WebhookHandler defines the interface for webhook handlers
type WebhookHandler interface {
	HandleSePayWebhook(c *gin.Context)
}

type webhookHandler struct {
	processWebhookUseCase payment.ProcessWebhookUseCase
	sepayAPIKey           string
}

// NewWebhookHandler creates a new WebhookHandler instance
func NewWebhookHandler(processWebhookUseCase payment.ProcessWebhookUseCase, sepayAPIKey string) WebhookHandler {
	return &webhookHandler{
		processWebhookUseCase: processWebhookUseCase,
		sepayAPIKey:           sepayAPIKey,
	}
}

// HandleSePayWebhook handles POST /api/v1/webhooks/sepay
// @Summary Handle SePay webhook callback
// @Description Processes incoming webhook from SePay payment gateway
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer <token>"
// @Param X-SePay-API-Key header string false "Alternative API Key"
// @Param body body dto.ProcessWebhookRequest true "Webhook request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/webhooks/sepay [post]
func (h *webhookHandler) HandleSePayWebhook(c *gin.Context) {
	// Verify API key
	// Check X-SePay-API-Key header (Legacy/Alternative)
	apiKey := c.GetHeader("X-SePay-API-Key")

	// Check Authorization header (Bearer token)
	if apiKey == "" {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if apiKey != h.sepayAPIKey {
		response.Unauthorized(c, "invalid or missing API key")
		return
	}

	// Parse request body
	var req dto.ProcessWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid JSON payload")
		return
	}

	// Process webhook
	tx, err := h.processWebhookUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, map[string]interface{}{
		"transactionId": tx.ID,
		"status":        tx.Status,
	})
}
