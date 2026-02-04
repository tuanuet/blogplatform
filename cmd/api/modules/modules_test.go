package modules

import (
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"

	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/domain/service"
	planhandler "github.com/aiagent/internal/interfaces/http/handler/plan"
)

// TestDIResolution verifies that all new subscription plan dependencies
// can be resolved through the Fx dependency injection container
func TestDIResolution(t *testing.T) {
	// Create a test application with all modules
	app := fxtest.New(
		t,
		fx.Supply(
			// Mock database connection for testing
			&gorm.DB{},
		),
		RepositoryModule,
		DomainServiceModule,
		HandlerModule,
		fx.Invoke(func(
			// Verify repositories are provided
			planRepo repository.SubscriptionPlanRepository,
			tagTierRepo repository.TagTierMappingRepository,

			// Verify services are provided
			planService service.PlanManagementService,
			tagService service.TagTierService,
			accessService service.ContentAccessService,

			// Verify handler is provided
			planHandler planhandler.PlanHandler,
		) {
			// If this runs successfully, all dependencies are resolved
			if planRepo == nil {
				t.Error("SubscriptionPlanRepository is nil")
			}
			if tagTierRepo == nil {
				t.Error("TagTierMappingRepository is nil")
			}
			if planService == nil {
				t.Error("PlanManagementService is nil")
			}
			if tagService == nil {
				t.Error("TagTierService is nil")
			}
			if accessService == nil {
				t.Error("ContentAccessService is nil")
			}
			if planHandler == nil {
				t.Error("PlanHandler is nil")
			}
		}),
	)

	defer app.RequireStart().RequireStop()
}
