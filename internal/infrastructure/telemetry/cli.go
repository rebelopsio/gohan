package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// CLIMetrics holds metrics for CLI commands
type CLIMetrics struct {
	commandDuration metric.Float64Histogram
	commandCounter  metric.Int64Counter
	commandErrors   metric.Int64Counter
}

// NewCLIMetrics creates a new CLI metrics instance
func NewCLIMetrics() (*CLIMetrics, error) {
	meter := otel.Meter("gohan-cli")

	commandDuration, err := meter.Float64Histogram(
		"cli.command.duration",
		metric.WithDescription("Duration of CLI command execution"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create command duration histogram: %w", err)
	}

	commandCounter, err := meter.Int64Counter(
		"cli.command.count",
		metric.WithDescription("Number of CLI commands executed"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create command counter: %w", err)
	}

	commandErrors, err := meter.Int64Counter(
		"cli.command.errors",
		metric.WithDescription("Number of CLI command errors"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create command error counter: %w", err)
	}

	return &CLIMetrics{
		commandDuration: commandDuration,
		commandCounter:  commandCounter,
		commandErrors:   commandErrors,
	}, nil
}

// RecordCommand records metrics for a CLI command execution
func (m *CLIMetrics) RecordCommand(ctx context.Context, command string, duration time.Duration, err error) {
	attrs := []attribute.KeyValue{
		attribute.String("command", command),
	}

	// Record duration
	m.commandDuration.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))

	// Record count
	if err != nil {
		attrs = append(attrs, attribute.String("status", "error"))
		m.commandErrors.Add(ctx, 1, metric.WithAttributes(attrs...))
	} else {
		attrs = append(attrs, attribute.String("status", "success"))
	}

	m.commandCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// TraceCommand creates a span for a CLI command and returns a cleanup function
// Usage:
//
//	ctx, end := TraceCommand(ctx, "gohan init")
//	defer end(err)
//	// ... command execution
func TraceCommand(ctx context.Context, command string) (context.Context, func(error)) {
	tracer := otel.Tracer("gohan-cli")
	ctx, span := tracer.Start(ctx, command,
		trace.WithAttributes(
			attribute.String("command", command),
		),
	)

	startTime := time.Now()

	return ctx, func(err error) {
		duration := time.Since(startTime)

		span.SetAttributes(
			attribute.Int64("duration_ms", duration.Milliseconds()),
		)

		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		} else {
			span.SetStatus(codes.Ok, "")
		}

		span.End()
	}
}

// InitCLITelemetry initializes telemetry for CLI commands with immediate flushing
// This is optimized for short-lived CLI executions
func InitCLITelemetry(cfg Config) (*Provider, *CLIMetrics, error) {
	if !cfg.Enabled {
		return &Provider{}, nil, nil
	}

	// Create provider with immediate flushing for CLI
	provider, err := NewProvider(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize telemetry provider: %w", err)
	}

	// Create CLI metrics
	metrics, err := NewCLIMetrics()
	if err != nil {
		return provider, nil, fmt.Errorf("failed to initialize CLI metrics: %w", err)
	}

	return provider, metrics, nil
}

// CLICommand wraps a CLI command execution with telemetry
// Usage:
//
//	err := telemetry.CLICommand(ctx, provider, metrics, "gohan init", func(ctx context.Context) error {
//	    // ... command implementation
//	    return nil
//	})
func CLICommand(
	ctx context.Context,
	provider *Provider,
	metrics *CLIMetrics,
	command string,
	fn func(context.Context) error,
) error {
	// Create span for command
	ctx, endTrace := TraceCommand(ctx, command)

	// Record execution time
	startTime := time.Now()

	// Execute command
	err := fn(ctx)

	// Calculate duration
	duration := time.Since(startTime)

	// End trace
	endTrace(err)

	// Record metrics if available
	if metrics != nil {
		metrics.RecordCommand(ctx, command, duration, err)
	}

	// Ensure all telemetry is flushed before CLI exits
	if provider != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if shutdownErr := provider.Shutdown(shutdownCtx); shutdownErr != nil {
			// Don't fail the command if telemetry shutdown fails
			// Just log it (in production, use proper logger)
			fmt.Printf("Warning: failed to shutdown telemetry: %v\n", shutdownErr)
		}
	}

	return err
}
