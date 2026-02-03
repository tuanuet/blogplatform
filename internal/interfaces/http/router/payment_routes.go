package router

import (
	"github.com/aiagent/internal/interfaces/http/handler/payment"
	"github.com/gin-gonic/gin"
)

// RegisterPaymentRoutes registers payment and webhook routes
func RegisterPaymentRoutes(v1 *gin.RouterGroup, paymentH payment.PaymentHandler, webhookH payment.WebhookHandler, sessionAuth gin.HandlerFunc) {
	// Payment routes (authenticated)
	payments := v1.Group("/payments", sessionAuth)
	{
		payments.POST("", paymentH.CreatePayment)
	}

	// Webhook routes (public)
	webhooks := v1.Group("/webhooks")
	{
		webhooks.POST("/sepay", webhookH.HandleSePayWebhook)
	}
}
