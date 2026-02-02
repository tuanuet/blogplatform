package usecase

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"time"

	"github.com/aiagent/internal/application/dto"
	domainService "github.com/aiagent/internal/domain/service"
)

// HealthUseCase handles health check application logic
type HealthUseCase interface {
	Check(ctx context.Context) *dto.HealthResponse
}

type healthUseCase struct {
	systemSvc domainService.SystemService
}

// NewHealthUseCase creates a new health use case
func NewHealthUseCase(systemSvc domainService.SystemService) HealthUseCase {
	return &healthUseCase{
		systemSvc: systemSvc,
	}
}

func (uc *healthUseCase) Check(ctx context.Context) *dto.HealthResponse {
	status, services := uc.systemSvc.CheckHealth(ctx)

	resp := &dto.HealthResponse{
		Status:    string(status),
		Timestamp: time.Now(),
		Services: dto.ServiceHealth{
			Database: dto.ServiceStatusDisconnected,
			Redis:    dto.ServiceStatusDisconnected,
		},
	}

	if services["database"] {
		resp.Services.Database = dto.ServiceStatusConnected
	}
	if services["redis"] {
		resp.Services.Redis = dto.ServiceStatusConnected
	}

	return resp
}
