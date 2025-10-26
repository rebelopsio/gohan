package services_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigurationMerger(t *testing.T) {
	t.Run("creates configuration merger", func(t *testing.T) {
		merger := services.NewConfigurationMerger()

		assert.NotNil(t, merger)
	})

	t.Run("implements ConfigurationMerger interface", func(t *testing.T) {
		merger := services.NewConfigurationMerger()

		var _ installation.ConfigurationMerger = merger
	})
}

func TestConfigurationMerger_MergeConfigurations(t *testing.T) {
	merger := services.NewConfigurationMerger()
	ctx := context.Background()

	t.Run("merges configurations with new components", func(t *testing.T) {
		existing := mustCreateConfiguration(t, []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.34.0"),
		})

		new := mustCreateConfiguration(t, []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.20"),
		})

		merged, err := merger.MergeConfigurations(ctx, existing, new)

		require.NoError(t, err)
		assert.Equal(t, 2, merged.ComponentCount())
		// Should use new version of Hyprland
		components := merged.Components()
		assert.Equal(t, "0.35.0", components[0].Version())
	})

	t.Run("preserves GPU support from existing", func(t *testing.T) {
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

		merged, err := merger.MergeConfigurations(ctx, existing, new)

		require.NoError(t, err)
		assert.True(t, merged.HasGPUSupport())
	})

	t.Run("adds components from new that aren't in existing", func(t *testing.T) {
		existing := mustCreateConfiguration(t, []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.34.0"),
		})

		new := mustCreateConfiguration(t, []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.20"),
			mustCreateComponentSelection(t, installation.ComponentRofi, "1.7.5"),
		})

		merged, err := merger.MergeConfigurations(ctx, existing, new)

		require.NoError(t, err)
		assert.Equal(t, 3, merged.ComponentCount())
	})

	t.Run("uses newer versions", func(t *testing.T) {
		existing := mustCreateConfiguration(t, []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.34.0"),
			mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.19"),
		})

		new := mustCreateConfiguration(t, []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.20"),
		})

		merged, err := merger.MergeConfigurations(ctx, existing, new)

		require.NoError(t, err)
		components := merged.Components()
		assert.Equal(t, 2, len(components))
		// Should use newer versions
		for _, comp := range components {
			if comp.Component() == installation.ComponentHyprland {
				assert.Equal(t, "0.35.0", comp.Version())
			}
			if comp.Component() == installation.ComponentWaybar {
				assert.Equal(t, "0.9.20", comp.Version())
			}
		}
	})

	t.Run("preserves merge existing config flag", func(t *testing.T) {
		// Create with merge flag true
		components := []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.34.0"),
		}
		existing, err := installation.NewInstallationConfiguration(
			components,
			nil,
			mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			true, // mergeExistingConfig = true
		)
		require.NoError(t, err)

		new := mustCreateConfiguration(t, []installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
		})

		merged, err := merger.MergeConfigurations(ctx, existing, new)

		require.NoError(t, err)
		assert.True(t, merged.MergeExistingConfig())
	})
}

func TestConfigurationMerger_ShouldBackupExisting(t *testing.T) {
	merger := services.NewConfigurationMerger()
	ctx := context.Background()

	t.Run("returns true when file exists", func(t *testing.T) {
		// Create a temporary file
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "hyprland.conf")
		err := os.WriteFile(configPath, []byte("# test config"), 0644)
		require.NoError(t, err)

		shouldBackup, err := merger.ShouldBackupExisting(ctx, configPath)

		require.NoError(t, err)
		assert.True(t, shouldBackup)
	})

	t.Run("returns false when file does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "nonexistent.conf")

		shouldBackup, err := merger.ShouldBackupExisting(ctx, configPath)

		require.NoError(t, err)
		assert.False(t, shouldBackup)
	})

	t.Run("returns false for empty path", func(t *testing.T) {
		shouldBackup, err := merger.ShouldBackupExisting(ctx, "")

		require.NoError(t, err)
		assert.False(t, shouldBackup)
	})

	t.Run("handles directory instead of file", func(t *testing.T) {
		tmpDir := t.TempDir()

		shouldBackup, err := merger.ShouldBackupExisting(ctx, tmpDir)

		require.NoError(t, err)
		// Directory should not be backed up as a config file
		assert.False(t, shouldBackup)
	})
}

// Helper functions

func mustCreateConfiguration(t *testing.T, components []installation.ComponentSelection) installation.InstallationConfiguration {
	t.Helper()
	config, err := installation.NewInstallationConfiguration(
		components,
		nil,
		mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
		false,
	)
	require.NoError(t, err)
	return config
}

func mustCreateComponentSelection(t *testing.T, component installation.ComponentName, version string) installation.ComponentSelection {
	t.Helper()
	selection, err := installation.NewComponentSelection(component, version, nil)
	require.NoError(t, err)
	return selection
}

func mustCreateDiskSpace(t *testing.T, availableBytes, requiredBytes int64) installation.DiskSpace {
	t.Helper()
	diskSpace, err := installation.NewDiskSpace(uint64(availableBytes), uint64(requiredBytes))
	require.NoError(t, err)
	return diskSpace
}

func mustCreateGPUSupport(t *testing.T, vendor string, requiresDriver bool, driverComponent installation.ComponentName) *installation.GPUSupport {
	t.Helper()
	gpu, err := installation.NewGPUSupport(vendor, requiresDriver, driverComponent)
	require.NoError(t, err)
	return &gpu
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
	require.NoError(t, err)
	return config
}
