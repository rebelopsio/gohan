package installation

import (
	"context"
	"time"
)

// ConflictResolver is a domain service for detecting and resolving package conflicts
// This is defined as an interface (port) for hexagonal architecture
// Implementations will be provided by the infrastructure layer
type ConflictResolver interface {
	// DetectConflicts checks for package conflicts in the selected components
	DetectConflicts(ctx context.Context, components []ComponentSelection) ([]PackageConflict, error)

	// ResolveConflict applies a resolution strategy to a conflict
	ResolveConflict(ctx context.Context, conflict PackageConflict, strategy ResolutionAction) error
}

// ProgressEstimator is a domain service for calculating installation progress and estimates
// Contains domain logic for time estimation based on installation phases
type ProgressEstimator interface {
	// EstimateRemainingTime calculates how much time is left based on current progress
	EstimateRemainingTime(
		currentPhase InstallationStatus,
		percentComplete int,
		elapsedTime time.Duration,
	) time.Duration

	// CalculatePhaseProgress calculates the percentage complete for a specific phase
	CalculatePhaseProgress(
		phase InstallationStatus,
		totalItems, completedItems int,
	) int
}

// ConfigurationMerger is a domain service for merging installation configurations
// Handles the logic of combining new configurations with existing ones
type ConfigurationMerger interface {
	// MergeConfigurations combines existing and new configurations
	// Preserves user settings while applying new defaults
	MergeConfigurations(
		ctx context.Context,
		existing, new InstallationConfiguration,
	) (InstallationConfiguration, error)

	// ShouldBackupExisting determines if existing configuration should be backed up
	ShouldBackupExisting(ctx context.Context, path string) (bool, error)
}
