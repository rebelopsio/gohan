package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// InstallationMetrics holds metrics for installation operations
type InstallationMetrics struct {
	installationDuration    metric.Float64Histogram
	installationCounter     metric.Int64Counter
	packagesInstalled       metric.Int64Counter
	installationSize        metric.Int64Histogram
	activeSessions          metric.Int64UpDownCounter
	sessionsByStatus        metric.Int64ObservableGauge
	packageInstallDuration  metric.Float64Histogram
	configurationValidation metric.Int64Counter
}

// NewInstallationMetrics creates a new installation metrics instance
func NewInstallationMetrics() (*InstallationMetrics, error) {
	meter := otel.Meter("gohan-installation")

	installationDuration, err := meter.Float64Histogram(
		"installation.duration",
		metric.WithDescription("Duration of installation operations"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create installation duration histogram: %w", err)
	}

	installationCounter, err := meter.Int64Counter(
		"installation.count",
		metric.WithDescription("Number of installations"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create installation counter: %w", err)
	}

	packagesInstalled, err := meter.Int64Counter(
		"installation.packages.count",
		metric.WithDescription("Number of packages installed"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create packages installed counter: %w", err)
	}

	installationSize, err := meter.Int64Histogram(
		"installation.size",
		metric.WithDescription("Total size of installed packages"),
		metric.WithUnit("bytes"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create installation size histogram: %w", err)
	}

	activeSessions, err := meter.Int64UpDownCounter(
		"installation.sessions.active",
		metric.WithDescription("Number of active installation sessions"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create active sessions counter: %w", err)
	}

	packageInstallDuration, err := meter.Float64Histogram(
		"installation.package.duration",
		metric.WithDescription("Duration of individual package installation"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create package install duration histogram: %w", err)
	}

	configurationValidation, err := meter.Int64Counter(
		"installation.configuration.validation",
		metric.WithDescription("Configuration validation results"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create configuration validation counter: %w", err)
	}

	return &InstallationMetrics{
		installationDuration:    installationDuration,
		installationCounter:     installationCounter,
		packagesInstalled:       packagesInstalled,
		installationSize:        installationSize,
		activeSessions:          activeSessions,
		packageInstallDuration:  packageInstallDuration,
		configurationValidation: configurationValidation,
	}, nil
}

// RecordInstallation records metrics for a complete installation
func (m *InstallationMetrics) RecordInstallation(
	ctx context.Context,
	status string,
	duration time.Duration,
	packageCount int,
	totalSize int64,
) {
	attrs := []attribute.KeyValue{
		attribute.String("status", status),
	}

	// Record duration
	m.installationDuration.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))

	// Record installation count
	m.installationCounter.Add(ctx, 1, metric.WithAttributes(attrs...))

	// Record packages installed
	if packageCount > 0 {
		m.packagesInstalled.Add(ctx, int64(packageCount), metric.WithAttributes(attrs...))
	}

	// Record installation size
	if totalSize > 0 {
		m.installationSize.Record(ctx, totalSize, metric.WithAttributes(attrs...))
	}
}

// RecordPackageInstallation records metrics for a single package installation
func (m *InstallationMetrics) RecordPackageInstallation(
	ctx context.Context,
	packageName string,
	duration time.Duration,
	success bool,
) {
	status := "success"
	if !success {
		status = "failure"
	}

	attrs := []attribute.KeyValue{
		attribute.String("package", packageName),
		attribute.String("status", status),
	}

	m.packageInstallDuration.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))
}

// RecordConfigurationValidation records configuration validation results
func (m *InstallationMetrics) RecordConfigurationValidation(
	ctx context.Context,
	valid bool,
	validationType string,
) {
	status := "valid"
	if !valid {
		status = "invalid"
	}

	attrs := []attribute.KeyValue{
		attribute.String("type", validationType),
		attribute.String("status", status),
	}

	m.configurationValidation.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// IncrementActiveSessions increments the active sessions counter
func (m *InstallationMetrics) IncrementActiveSessions(ctx context.Context) {
	m.activeSessions.Add(ctx, 1)
}

// DecrementActiveSessions decrements the active sessions counter
func (m *InstallationMetrics) DecrementActiveSessions(ctx context.Context) {
	m.activeSessions.Add(ctx, -1)
}
