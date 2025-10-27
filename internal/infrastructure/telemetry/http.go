package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// HTTPMetrics holds metrics for HTTP server
type HTTPMetrics struct {
	requestDuration   metric.Float64Histogram
	requestCounter    metric.Int64Counter
	activeConnections metric.Int64UpDownCounter
	requestSize       metric.Int64Histogram
	responseSize      metric.Int64Histogram
}

// NewHTTPMetrics creates a new HTTP metrics instance
func NewHTTPMetrics() (*HTTPMetrics, error) {
	meter := otel.Meter("gohan-http")

	requestDuration, err := meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithDescription("Duration of HTTP server requests"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request duration histogram: %w", err)
	}

	requestCounter, err := meter.Int64Counter(
		"http.server.request.count",
		metric.WithDescription("Number of HTTP requests"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request counter: %w", err)
	}

	activeConnections, err := meter.Int64UpDownCounter(
		"http.server.active_connections",
		metric.WithDescription("Number of active HTTP connections"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create active connections counter: %w", err)
	}

	requestSize, err := meter.Int64Histogram(
		"http.server.request.size",
		metric.WithDescription("Size of HTTP request bodies"),
		metric.WithUnit("bytes"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request size histogram: %w", err)
	}

	responseSize, err := meter.Int64Histogram(
		"http.server.response.size",
		metric.WithDescription("Size of HTTP response bodies"),
		metric.WithUnit("bytes"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create response size histogram: %w", err)
	}

	return &HTTPMetrics{
		requestDuration:   requestDuration,
		requestCounter:    requestCounter,
		activeConnections: activeConnections,
		requestSize:       requestSize,
		responseSize:      responseSize,
	}, nil
}

// RecordRequest records metrics for an HTTP request
func (m *HTTPMetrics) RecordRequest(
	ctx context.Context,
	method string,
	path string,
	statusCode int,
	duration time.Duration,
	requestSize int64,
	responseSize int64,
) {
	attrs := []attribute.KeyValue{
		attribute.String("http.method", method),
		attribute.String("http.route", path),
		attribute.Int("http.status_code", statusCode),
	}

	// Record duration
	m.requestDuration.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))

	// Record count
	m.requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))

	// Record sizes
	if requestSize > 0 {
		m.requestSize.Record(ctx, requestSize, metric.WithAttributes(attrs...))
	}
	if responseSize > 0 {
		m.responseSize.Record(ctx, responseSize, metric.WithAttributes(attrs...))
	}
}

// IncrementActiveConnections increments the active connections counter
func (m *HTTPMetrics) IncrementActiveConnections(ctx context.Context) {
	m.activeConnections.Add(ctx, 1)
}

// DecrementActiveConnections decrements the active connections counter
func (m *HTTPMetrics) DecrementActiveConnections(ctx context.Context) {
	m.activeConnections.Add(ctx, -1)
}
