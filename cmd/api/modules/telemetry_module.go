package modules

import (
	"context"

	"github.com/aiagent/boilerplate/internal/infrastructure/telemetry"
	"github.com/aiagent/boilerplate/pkg/logger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

// TelemetryModule provides telemetry services
var TelemetryModule = fx.Module("telemetry",
	fx.Provide(telemetry.NewTracerProvider),
	fx.Invoke(registerTelemetryLifecycle),
)

func registerTelemetryLifecycle(lc fx.Lifecycle, tp *sdktrace.TracerProvider) {
	if tp == nil {
		return
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down TracerProvider...")
			if err := tp.Shutdown(ctx); err != nil {
				logger.Error("Error shutting down TracerProvider", err, nil)
				return err
			}
			return nil
		},
	})
}
