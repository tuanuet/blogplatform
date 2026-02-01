package dto

import "time"

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string        `json:"status"`
	Timestamp time.Time     `json:"timestamp"`
	Services  ServiceHealth `json:"services"`
	Version   string        `json:"version,omitempty"`
}

// ServiceHealth represents individual service health statuses
type ServiceHealth struct {
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

// HealthStatus constants
const (
	HealthStatusHealthy   = "healthy"
	HealthStatusUnhealthy = "unhealthy"
	HealthStatusDegraded  = "degraded"

	ServiceStatusConnected    = "connected"
	ServiceStatusDisconnected = "disconnected"
)
