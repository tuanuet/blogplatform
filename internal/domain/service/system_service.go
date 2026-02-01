package service

import (
	"context"

	"github.com/aiagent/boilerplate/internal/domain/repository"
)

// System status constants
type SystemStatus string

const (
	SystemStatusHealthy   SystemStatus = "healthy"
	SystemStatusDegraded  SystemStatus = "degraded"
	SystemStatusUnhealthy SystemStatus = "unhealthy"
)

// SystemService handles system health domain logic
type SystemService interface {
	// CheckHealth checks the health of the system
	CheckHealth(ctx context.Context) (status SystemStatus, services map[string]bool)
}

type systemService struct {
	repo repository.SystemRepository
}

// NewSystemService creates a new system domain service
func NewSystemService(repo repository.SystemRepository) SystemService {
	return &systemService{
		repo: repo,
	}
}

func (s *systemService) CheckHealth(ctx context.Context) (SystemStatus, map[string]bool) {
	services := make(map[string]bool)

	// Check Database
	dbErr := s.repo.CheckDatabase(ctx)
	services["database"] = dbErr == nil

	// Check Redis
	redisErr := s.repo.CheckRedis(ctx)
	services["redis"] = redisErr == nil

	// Determine overall status
	status := SystemStatusUnhealthy
	if services["database"] && services["redis"] {
		status = SystemStatusHealthy
	} else if services["database"] {
		status = SystemStatusDegraded
	}

	return status, services
}
