package services_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/services"
	"github.com/stretchr/testify/assert"
)

func TestNewProgressEstimator(t *testing.T) {
	t.Run("creates progress estimator", func(t *testing.T) {
		estimator := services.NewProgressEstimator()

		assert.NotNil(t, estimator)
	})

	t.Run("implements ProgressEstimator interface", func(t *testing.T) {
		estimator := services.NewProgressEstimator()

		var _ installation.ProgressEstimator = estimator
	})
}

func TestProgressEstimator_EstimateRemainingTime(t *testing.T) {
	estimator := services.NewProgressEstimator()

	tests := []struct {
		name            string
		currentPhase    installation.InstallationStatus
		percentComplete int
		elapsedTime     time.Duration
		validateResult  func(*testing.T, time.Duration)
	}{
		{
			name:            "50% complete estimates equal time remaining",
			currentPhase:    installation.StatusInstalling,
			percentComplete: 50,
			elapsedTime:     60 * time.Second,
			validateResult: func(t *testing.T, remaining time.Duration) {
				// At 50%, should estimate approximately equal time remaining
				assert.InDelta(t, 60.0, remaining.Seconds(), 10.0)
			},
		},
		{
			name:            "75% complete estimates 1/3 time remaining",
			currentPhase:    installation.StatusInstalling,
			percentComplete: 75,
			elapsedTime:     90 * time.Second,
			validateResult: func(t *testing.T, remaining time.Duration) {
				// At 75%, should estimate approximately 1/3 time remaining
				assert.InDelta(t, 30.0, remaining.Seconds(), 10.0)
			},
		},
		{
			name:            "10% complete estimates long time remaining",
			currentPhase:    installation.StatusPreparation,
			percentComplete: 10,
			elapsedTime:     10 * time.Second,
			validateResult: func(t *testing.T, remaining time.Duration) {
				// At 10%, should estimate approximately 9x time remaining
				assert.Greater(t, remaining.Seconds(), 80.0)
			},
		},
		{
			name:            "100% complete returns zero",
			currentPhase:    installation.StatusCompleted,
			percentComplete: 100,
			elapsedTime:     120 * time.Second,
			validateResult: func(t *testing.T, remaining time.Duration) {
				assert.Equal(t, time.Duration(0), remaining)
			},
		},
		{
			name:            "0% complete handles gracefully",
			currentPhase:    installation.StatusPending,
			percentComplete: 0,
			elapsedTime:     1 * time.Second,
			validateResult: func(t *testing.T, remaining time.Duration) {
				// Should return a reasonable estimate, not infinite
				assert.Greater(t, remaining, time.Duration(0))
				assert.Less(t, remaining, 24*time.Hour) // Sanity check
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining := estimator.EstimateRemainingTime(
				tt.currentPhase,
				tt.percentComplete,
				tt.elapsedTime,
			)

			tt.validateResult(t, remaining)
		})
	}
}

func TestProgressEstimator_CalculatePhaseProgress(t *testing.T) {
	estimator := services.NewProgressEstimator()

	tests := []struct {
		name           string
		phase          installation.InstallationStatus
		totalItems     int
		completedItems int
		expectedPercent int
	}{
		{
			name:            "0 of 10 items",
			phase:           installation.StatusInstalling,
			totalItems:      10,
			completedItems:  0,
			expectedPercent: 0,
		},
		{
			name:            "5 of 10 items",
			phase:           installation.StatusInstalling,
			totalItems:      10,
			completedItems:  5,
			expectedPercent: 50,
		},
		{
			name:            "10 of 10 items",
			phase:           installation.StatusInstalling,
			totalItems:      10,
			completedItems:  10,
			expectedPercent: 100,
		},
		{
			name:            "1 of 3 items",
			phase:           installation.StatusDownloading,
			totalItems:      3,
			completedItems:  1,
			expectedPercent: 33,
		},
		{
			name:            "2 of 3 items",
			phase:           installation.StatusDownloading,
			totalItems:      3,
			completedItems:  2,
			expectedPercent: 66,
		},
		{
			name:            "handles zero total items",
			phase:           installation.StatusInstalling,
			totalItems:      0,
			completedItems:  0,
			expectedPercent: 0,
		},
		{
			name:            "handles completed > total",
			phase:           installation.StatusInstalling,
			totalItems:      5,
			completedItems:  7,
			expectedPercent: 100, // Caps at 100%
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			percent := estimator.CalculatePhaseProgress(
				tt.phase,
				tt.totalItems,
				tt.completedItems,
			)

			assert.Equal(t, tt.expectedPercent, percent)
		})
	}
}

func TestProgressEstimator_PhaseWeighting(t *testing.T) {
	estimator := services.NewProgressEstimator()

	t.Run("different phases have different weights when no progress", func(t *testing.T) {
		// Phase weights are used when percent is 0
		elapsedTime := 30 * time.Second

		preparationTime := estimator.EstimateRemainingTime(
			installation.StatusPreparation,
			0, // No progress yet
			elapsedTime,
		)

		installingTime := estimator.EstimateRemainingTime(
			installation.StatusInstalling,
			0, // No progress yet
			elapsedTime,
		)

		// Different phases have different weights, so estimates should differ
		// (This test validates that phase weighting exists, not specific values)
		assert.NotEqual(t, preparationTime, installingTime)
	})

	t.Run("same progress in different phases estimates same time", func(t *testing.T) {
		// Linear estimation is used when there's progress
		elapsedTime := 30 * time.Second

		preparationTime := estimator.EstimateRemainingTime(
			installation.StatusPreparation,
			25,
			elapsedTime,
		)

		installingTime := estimator.EstimateRemainingTime(
			installation.StatusInstalling,
			25,
			elapsedTime,
		)

		// Same percent progress should give same linear estimate
		assert.Equal(t, preparationTime, installingTime)
	})
}

func TestProgressEstimator_ConsistentCalculations(t *testing.T) {
	estimator := services.NewProgressEstimator()

	t.Run("same inputs produce same outputs", func(t *testing.T) {
		phase := installation.StatusInstalling
		percent := 42
		elapsed := 60 * time.Second

		result1 := estimator.EstimateRemainingTime(phase, percent, elapsed)
		result2 := estimator.EstimateRemainingTime(phase, percent, elapsed)

		assert.Equal(t, result1, result2, "Estimator should be deterministic")
	})

	t.Run("progress calculation is consistent", func(t *testing.T) {
		phase := installation.StatusDownloading
		total := 10
		completed := 7

		result1 := estimator.CalculatePhaseProgress(phase, total, completed)
		result2 := estimator.CalculatePhaseProgress(phase, total, completed)

		assert.Equal(t, result1, result2, "Progress calculation should be deterministic")
	})
}
