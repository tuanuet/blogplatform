package telemetry

import (
	"github.com/aiagent/boilerplate/internal/infrastructure/config"
	"github.com/aiagent/boilerplate/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// NewTracerProvider creates a new TracerProvider with a stdout exporter
func NewTracerProvider(cfg *config.Config) (*sdktrace.TracerProvider, error) {
	if !cfg.Telemetry.Enabled {
		logger.Info("Telemetry is disabled")
		return nil, nil
	}

	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceName(cfg.Telemetry.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set global provider
	otel.SetTracerProvider(tp)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	logger.Info("Telemetry enabled with stdout exporter", map[string]interface{}{
		"service": cfg.Telemetry.ServiceName,
	})

	return tp, nil
}
