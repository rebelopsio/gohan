package installation_test

import (
	"context"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for testing

type MockConflictResolver struct {
	mock.Mock
}

func (m *MockConflictResolver) DetectConflicts(ctx context.Context, components []installation.ComponentSelection) ([]installation.PackageConflict, error) {
	args := m.Called(ctx, components)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]installation.PackageConflict), args.Error(1)
}

func (m *MockConflictResolver) ResolveConflict(ctx context.Context, conflict installation.PackageConflict, strategy installation.ResolutionAction) error {
	args := m.Called(ctx, conflict, strategy)
	return args.Error(0)
}

type MockProgressEstimator struct {
	mock.Mock
}

func (m *MockProgressEstimator) EstimateRemainingTime(
	currentPhase installation.InstallationStatus,
	percentComplete int,
	elapsedTime time.Duration,
) time.Duration {
	args := m.Called(currentPhase, percentComplete, elapsedTime)
	return args.Get(0).(time.Duration)
}

func (m *MockProgressEstimator) CalculatePhaseProgress(
	phase installation.InstallationStatus,
	totalItems, completedItems int,
) int {
	args := m.Called(phase, totalItems, completedItems)
	return args.Get(0).(int)
}

type MockConfigurationMerger struct {
	mock.Mock
}

func (m *MockConfigurationMerger) MergeConfigurations(
	ctx context.Context,
	existing, new installation.InstallationConfiguration,
) (installation.InstallationConfiguration, error) {
	args := m.Called(ctx, existing, new)
	if args.Get(0) == nil {
		return installation.InstallationConfiguration{}, args.Error(1)
	}
	return args.Get(0).(installation.InstallationConfiguration), args.Error(1)
}

func (m *MockConfigurationMerger) ShouldBackupExisting(
	ctx context.Context,
	path string,
) (bool, error) {
	args := m.Called(ctx, path)
	return args.Get(0).(bool), args.Error(1)
}

// Tests for ConflictResolver interface

func TestConflictResolver_DetectConflicts(t *testing.T) {
	t.Run("detects no conflicts for compatible components", func(t *testing.T) {
		resolver := new(MockConflictResolver)
		ctx := context.Background()

		components := []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.20"),
		}

		resolver.On("DetectConflicts", ctx, components).Return([]installation.PackageConflict(nil), nil)

		conflicts, err := resolver.DetectConflicts(ctx, components)

		assert.NoError(t, err)
		assert.Empty(t, conflicts)
		resolver.AssertExpectations(t)
	})

	t.Run("detects conflicts between components", func(t *testing.T) {
		resolver := new(MockConflictResolver)
		ctx := context.Background()

		components := []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
		}

		expectedConflict := mustCreatePackageConflict(t, "hyprland", "hyprland-git")
		expectedConflicts := []installation.PackageConflict{expectedConflict}

		resolver.On("DetectConflicts", ctx, components).Return(expectedConflicts, nil)

		conflicts, err := resolver.DetectConflicts(ctx, components)

		assert.NoError(t, err)
		assert.Len(t, conflicts, 1)
		assert.Equal(t, "hyprland", conflicts[0].PackageName())
		resolver.AssertExpectations(t)
	})
}

func TestConflictResolver_ResolveConflict(t *testing.T) {
	t.Run("resolves conflict with remove strategy", func(t *testing.T) {
		resolver := new(MockConflictResolver)
		ctx := context.Background()

		conflict := mustCreatePackageConflict(t, "hyprland", "hyprland-git")
		strategy := installation.ActionRemove

		resolver.On("ResolveConflict", ctx, conflict, strategy).Return(nil)

		err := resolver.ResolveConflict(ctx, conflict, strategy)

		assert.NoError(t, err)
		resolver.AssertExpectations(t)
	})

	t.Run("resolves conflict with skip strategy", func(t *testing.T) {
		resolver := new(MockConflictResolver)
		ctx := context.Background()

		conflict := mustCreatePackageConflict(t, "package1", "package2")
		strategy := installation.ActionSkip

		resolver.On("ResolveConflict", ctx, conflict, strategy).Return(nil)

		err := resolver.ResolveConflict(ctx, conflict, strategy)

		assert.NoError(t, err)
		resolver.AssertExpectations(t)
	})
}

// Tests for ProgressEstimator interface

func TestProgressEstimator_EstimateRemainingTime(t *testing.T) {
	t.Run("estimates remaining time based on current progress", func(t *testing.T) {
		estimator := new(MockProgressEstimator)

		currentPhase := installation.StatusInstalling
		percentComplete := 50
		elapsedTime := 60 * time.Second
		expectedRemaining := 60 * time.Second // 50% done, 50% remaining

		estimator.On("EstimateRemainingTime", currentPhase, percentComplete, elapsedTime).
			Return(expectedRemaining)

		remaining := estimator.EstimateRemainingTime(currentPhase, percentComplete, elapsedTime)

		assert.Equal(t, expectedRemaining, remaining)
		estimator.AssertExpectations(t)
	})

	t.Run("estimates longer time for early phases", func(t *testing.T) {
		estimator := new(MockProgressEstimator)

		currentPhase := installation.StatusPreparation
		percentComplete := 10
		elapsedTime := 10 * time.Second
		expectedRemaining := 90 * time.Second

		estimator.On("EstimateRemainingTime", currentPhase, percentComplete, elapsedTime).
			Return(expectedRemaining)

		remaining := estimator.EstimateRemainingTime(currentPhase, percentComplete, elapsedTime)

		assert.Equal(t, expectedRemaining, remaining)
		estimator.AssertExpectations(t)
	})

	t.Run("returns zero when complete", func(t *testing.T) {
		estimator := new(MockProgressEstimator)

		currentPhase := installation.StatusCompleted
		percentComplete := 100
		elapsedTime := 120 * time.Second

		estimator.On("EstimateRemainingTime", currentPhase, percentComplete, elapsedTime).
			Return(time.Duration(0))

		remaining := estimator.EstimateRemainingTime(currentPhase, percentComplete, elapsedTime)

		assert.Equal(t, time.Duration(0), remaining)
		estimator.AssertExpectations(t)
	})
}

func TestProgressEstimator_CalculatePhaseProgress(t *testing.T) {
	t.Run("calculates progress percentage for phase", func(t *testing.T) {
		estimator := new(MockProgressEstimator)

		phase := installation.StatusInstalling
		totalItems := 10
		completedItems := 5
		expectedPercent := 50

		estimator.On("CalculatePhaseProgress", phase, totalItems, completedItems).
			Return(expectedPercent)

		percent := estimator.CalculatePhaseProgress(phase, totalItems, completedItems)

		assert.Equal(t, expectedPercent, percent)
		estimator.AssertExpectations(t)
	})

	t.Run("returns 100 when all items complete", func(t *testing.T) {
		estimator := new(MockProgressEstimator)

		phase := installation.StatusInstalling
		totalItems := 10
		completedItems := 10

		estimator.On("CalculatePhaseProgress", phase, totalItems, completedItems).
			Return(100)

		percent := estimator.CalculatePhaseProgress(phase, totalItems, completedItems)

		assert.Equal(t, 100, percent)
		estimator.AssertExpectations(t)
	})

	t.Run("returns 0 when no items complete", func(t *testing.T) {
		estimator := new(MockProgressEstimator)

		phase := installation.StatusDownloading
		totalItems := 5
		completedItems := 0

		estimator.On("CalculatePhaseProgress", phase, totalItems, completedItems).
			Return(0)

		percent := estimator.CalculatePhaseProgress(phase, totalItems, completedItems)

		assert.Equal(t, 0, percent)
		estimator.AssertExpectations(t)
	})
}

// Tests for ConfigurationMerger interface

func TestConfigurationMerger_MergeConfigurations(t *testing.T) {
	t.Run("merges configurations preserving user settings", func(t *testing.T) {
		merger := new(MockConfigurationMerger)
		ctx := context.Background()

		existing := mustCreateConfiguration(t, []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.34.0"),
		})

		new := mustCreateConfiguration(t, []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.20"),
		})

		expectedMerged := mustCreateConfiguration(t, []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.20"),
		})

		merger.On("MergeConfigurations", ctx, existing, new).Return(expectedMerged, nil)

		merged, err := merger.MergeConfigurations(ctx, existing, new)

		assert.NoError(t, err)
		assert.Equal(t, 2, merged.ComponentCount())
		merger.AssertExpectations(t)
	})

	t.Run("preserves GPU support from existing", func(t *testing.T) {
		merger := new(MockConfigurationMerger)
		ctx := context.Background()

		gpuSupport := mustCreateGPUSupport(t, "amd", true, installation.ComponentAMDDriver)

		existing := mustCreateConfigurationWithGPU(t,
			[]installation.ComponentSelection{
				mustCreateComponentSelection(t, installation.ComponentHyprland, "0.34.0"),
			},
			gpuSupport,
		)

		new := mustCreateConfiguration(t, []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
		})

		expectedMerged := mustCreateConfigurationWithGPU(t,
			[]installation.ComponentSelection{
				mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			},
			gpuSupport,
		)

		merger.On("MergeConfigurations", ctx, existing, new).Return(expectedMerged, nil)

		merged, err := merger.MergeConfigurations(ctx, existing, new)

		assert.NoError(t, err)
		assert.True(t, merged.HasGPUSupport())
		merger.AssertExpectations(t)
	})
}

func TestConfigurationMerger_ShouldBackupExisting(t *testing.T) {
	t.Run("returns true when existing configuration exists", func(t *testing.T) {
		merger := new(MockConfigurationMerger)
		ctx := context.Background()

		path := "/home/user/.config/hypr/hyprland.conf"

		merger.On("ShouldBackupExisting", ctx, path).Return(true, nil)

		shouldBackup, err := merger.ShouldBackupExisting(ctx, path)

		assert.NoError(t, err)
		assert.True(t, shouldBackup)
		merger.AssertExpectations(t)
	})

	t.Run("returns false when no existing configuration", func(t *testing.T) {
		merger := new(MockConfigurationMerger)
		ctx := context.Background()

		path := "/home/newuser/.config/hypr/hyprland.conf"

		merger.On("ShouldBackupExisting", ctx, path).Return(false, nil)

		shouldBackup, err := merger.ShouldBackupExisting(ctx, path)

		assert.NoError(t, err)
		assert.False(t, shouldBackup)
		merger.AssertExpectations(t)
	})
}

// Helper functions for test data

func mustCreatePackageConflict(t *testing.T, pkg, conflictingPkg string) installation.PackageConflict {
	t.Helper()
	conflict, err := installation.NewPackageConflict(pkg, conflictingPkg, "version conflict")
	if err != nil {
		t.Fatalf("Failed to create package conflict: %v", err)
	}
	return conflict
}

func mustCreateConfigurationWithGPU(
	t *testing.T,
	components []installation.ComponentSelection,
	gpuSupport *installation.GPUSupport,
) installation.InstallationConfiguration {
	t.Helper()
	config, err := installation.NewInstallationConfiguration(
		components,
		gpuSupport,
		mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
		false,
	)
	if err != nil {
		t.Fatalf("Failed to create configuration: %v", err)
	}
	return config
}
