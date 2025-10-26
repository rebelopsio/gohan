package services

import (
	"time"

	"github.com/rebelopsio/gohan/internal/domain/installation"
)

// ProgressEstimator implements installation.ProgressEstimator
// Calculates installation progress and time estimates
type ProgressEstimator struct {
	phaseWeights map[installation.InstallationStatus]float64
}

// NewProgressEstimator creates a new progress estimator
func NewProgressEstimator() *ProgressEstimator {
	return &ProgressEstimator{
		phaseWeights: map[installation.InstallationStatus]float64{
			installation.StatusPending:     0.01, // Very fast
			installation.StatusPreparation: 0.05, // Quick snapshot
			installation.StatusDownloading: 0.30, // Can be slow depending on network
			installation.StatusInstalling:  0.40, // Most time-consuming
			installation.StatusConfiguring: 0.15, // Moderate
			installation.StatusVerifying:   0.09, // Quick validation
			installation.StatusCompleted:   0.00, // Done
		},
	}
}

// EstimateRemainingTime implements installation.ProgressEstimator
// Calculates remaining time based on current progress
func (p *ProgressEstimator) EstimateRemainingTime(
	currentPhase installation.InstallationStatus,
	percentComplete int,
	elapsedTime time.Duration,
) time.Duration {
	// Handle edge cases
	if percentComplete >= 100 || currentPhase == installation.StatusCompleted {
		return 0
	}

	if percentComplete <= 0 {
		// Use phase weights to estimate total time
		return p.estimateFromPhaseWeights(currentPhase, elapsedTime)
	}

	// Simple linear estimation: remaining = (elapsed / percentComplete) * (100 - percentComplete)
	progressFraction := float64(percentComplete) / 100.0
	remainingFraction := 1.0 - progressFraction

	if progressFraction == 0 {
		return 0
	}

	totalEstimatedTime := float64(elapsedTime) / progressFraction
	remainingTime := totalEstimatedTime * remainingFraction

	return time.Duration(remainingTime)
}

// CalculatePhaseProgress implements installation.ProgressEstimator
// Calculates percentage complete for a specific phase
func (p *ProgressEstimator) CalculatePhaseProgress(
	phase installation.InstallationStatus,
	totalItems, completedItems int,
) int {
	// Handle edge cases
	if totalItems <= 0 {
		return 0
	}

	if completedItems >= totalItems {
		return 100
	}

	if completedItems <= 0 {
		return 0
	}

	// Calculate percentage
	percent := (completedItems * 100) / totalItems
	if percent > 100 {
		return 100
	}

	return percent
}

// Helper function to get phase weight
func (p *ProgressEstimator) getPhaseWeight(phase installation.InstallationStatus) float64 {
	weight, exists := p.phaseWeights[phase]
	if !exists {
		return 0.20 // Default weight
	}
	return weight
}

// Helper function to estimate time from phase weights when no progress yet
func (p *ProgressEstimator) estimateFromPhaseWeights(
	currentPhase installation.InstallationStatus,
	elapsedTime time.Duration,
) time.Duration {
	currentWeight := p.getPhaseWeight(currentPhase)
	if currentWeight == 0 {
		// Avoid division by zero
		return 5 * time.Minute // Default estimate
	}

	// Estimate total time based on how long current phase took
	// and what fraction of total work it represents
	estimatedTotal := float64(elapsedTime) / currentWeight

	// Subtract elapsed to get remaining
	remaining := estimatedTotal - float64(elapsedTime)
	if remaining < 0 {
		remaining = 0
	}

	return time.Duration(remaining)
}
