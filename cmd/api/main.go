package main

import (
	"github.com/aiagent/cmd/api/modules"
	"go.uber.org/fx"
)

// @title Go Boilerplate API
// @version 1.0
// @description Clean Architecture Go API with GORM, Gin, PostgreSQL, Redis
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the access token.

func main() {
	app := fx.New(
		// Configuration module (loads config, initializes logger & validator)
		modules.ConfigModule,

		// Telemetry module (OpenTelemetry)
		modules.TelemetryModule,

		// Database module (PostgreSQL + Redis with lifecycle hooks)
		modules.DatabaseModule,

		// Repository module (all repository implementations)
		modules.RepositoryModule,

		// Domain Service module (pure domain logic)
		modules.DomainServiceModule,

		// Use Case module (application logic)
		modules.UseCaseModule,

		// Handler module (all HTTP handlers)
		modules.HandlerModule,

		// HTTP module (router, gin engine, HTTP server with lifecycle hooks)
		modules.HTTPModule,

		// Scheduler module (background jobs and cron tasks)
		modules.SchedulerModule,

		// FX options
		fx.NopLogger, // Suppress FX's default logger (we use our own)
	)

	app.Run()
}
