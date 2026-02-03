package payment

import (
	"context"
	"net/http"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreatePaymentUseCase defines the interface for payment creation use case
type CreatePaymentUseCase interface {
	Execute(ctx context.Context, req dto.CreatePaymentRequest) (*dto.CreatePaymentResponse, error)
}

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler interface {
	CreatePayment(c *gin.Context)
}

type paymentHandler struct {
	createPaymentUseCase CreatePaymentUseCase
}

// NewPaymentHandler creates a new PaymentHandler instance
func NewPaymentHandler(createPaymentUseCase CreatePaymentUseCase) PaymentHandler {
	return &paymentHandler{
		createPaymentUseCase: createPaymentUseCase,
	}
}

// CreatePayment handles POST /api/v1/payments
// @Summary Create a new payment
// @Description Initiates a new payment transaction
// @Tags Payments
// @Accept json
// @Produce json
// @Param body body dto.CreatePaymentRequest true "Payment request"
// @Success 201 {object} dto.CreatePaymentResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Security Bearer
// @Router /api/v1/payments [post]
func (h *paymentHandler) CreatePayment(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Override UserID from JWT context for security
	uid, ok := userIDVal.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "invalid user ID")
		return
	}
	req.UserID = uid.String()

	resp, err := h.createPaymentUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, resp)
}
