package installation_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstallationConfiguration(t *testing.T) {
	tests := []struct {
		name              string
		components        []installation.ComponentSelection
		gpuSupport        *installation.GPUSupport
		diskSpace         installation.DiskSpace
		mergeExistingConf bool
		wantErr           bool
		errType           error
	}{
		{
			name: "valid configuration with all fields",
			components: []installation.ComponentSelection{
				mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
				mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.20"),
			},
			gpuSupport:        mustCreateGPUSupport(t, "amd", true, installation.ComponentAMDDriver),
			diskSpace:         mustCreateDiskSpace(t, 50*installation.GB, 10*installation.GB),
			mergeExistingConf: true,
			wantErr:           false,
		},
		{
			name: "valid configuration with minimal fields",
			components: []installation.ComponentSelection{
				mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			},
			gpuSupport:        nil,
			diskSpace:         mustCreateDiskSpace(t, 20*installation.GB, 10*installation.GB),
			mergeExistingConf: false,
			wantErr:           false,
		},
		{
			name:              "empty components list",
			components:        []installation.ComponentSelection{},
			gpuSupport:        nil,
			diskSpace:         mustCreateDiskSpace(t, 20*installation.GB, 10*installation.GB),
			mergeExistingConf: false,
			wantErr:           true,
			errType:           installation.ErrInvalidConfiguration,
		},
		{
			name: "missing Hyprland core component",
			components: []installation.ComponentSelection{
				mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.20"),
			},
			gpuSupport:        nil,
			diskSpace:         mustCreateDiskSpace(t, 20*installation.GB, 10*installation.GB),
			mergeExistingConf: false,
			wantErr:           true,
			errType:           installation.ErrInvalidConfiguration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := installation.NewInstallationConfiguration(
				tt.components,
				tt.gpuSupport,
				tt.diskSpace,
				tt.mergeExistingConf,
			)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, len(tt.components), len(config.Components()))
			assert.Equal(t, tt.mergeExistingConf, config.MergeExistingConfig())
			assert.Equal(t, tt.diskSpace.Available(), config.DiskSpace().Available())
		})
	}
}

func TestInstallationConfiguration_HasCoreComponent(t *testing.T) {
	t.Run("with Hyprland", func(t *testing.T) {
		config, err := installation.NewInstallationConfiguration(
			[]installation.ComponentSelection{
				mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			},
			nil,
			mustCreateDiskSpace(t, 20*installation.GB, 10*installation.GB),
			false,
		)
		require.NoError(t, err)

		assert.True(t, config.HasCoreComponent())
	})
}

func TestInstallationConfiguration_ComponentCount(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{
			name:  "single component",
			count: 1,
		},
		{
			name:  "multiple components",
			count: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			components := make([]installation.ComponentSelection, tt.count)
			components[0] = mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0")
			for i := 1; i < tt.count; i++ {
				components[i] = mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.20")
			}

			config, err := installation.NewInstallationConfiguration(
				components,
				nil,
				mustCreateDiskSpace(t, 20*installation.GB, 10*installation.GB),
				false,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.count, config.ComponentCount())
		})
	}
}

func TestInstallationConfiguration_HasGPUSupport(t *testing.T) {
	t.Run("with GPU support", func(t *testing.T) {
		gpu := mustCreateGPUSupport(t, "amd", true, installation.ComponentAMDDriver)
		config, err := installation.NewInstallationConfiguration(
			[]installation.ComponentSelection{
				mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			},
			gpu,
			mustCreateDiskSpace(t, 20*installation.GB, 10*installation.GB),
			false,
		)
		require.NoError(t, err)

		assert.True(t, config.HasGPUSupport())
		assert.NotNil(t, config.GPUSupport())
	})

	t.Run("without GPU support", func(t *testing.T) {
		config, err := installation.NewInstallationConfiguration(
			[]installation.ComponentSelection{
				mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			},
			nil,
			mustCreateDiskSpace(t, 20*installation.GB, 10*installation.GB),
			false,
		)
		require.NoError(t, err)

		assert.False(t, config.HasGPUSupport())
		assert.Nil(t, config.GPUSupport())
	})
}

func TestInstallationConfiguration_TotalEstimatedSize(t *testing.T) {
	t.Run("with package info", func(t *testing.T) {
		pkg1, _ := installation.NewPackageInfo("hyprland", "0.35.0", 50*installation.MB, nil)
		pkg2, _ := installation.NewPackageInfo("waybar", "0.9.20", 10*installation.MB, nil)

		comp1, _ := installation.NewComponentSelection(installation.ComponentHyprland, "0.35.0", &pkg1)
		comp2, _ := installation.NewComponentSelection(installation.ComponentWaybar, "0.9.20", &pkg2)

		config, err := installation.NewInstallationConfiguration(
			[]installation.ComponentSelection{comp1, comp2},
			nil,
			mustCreateDiskSpace(t, 100*installation.GB, 60*installation.MB),
			false,
		)
		require.NoError(t, err)

		assert.Equal(t, uint64(60*installation.MB), config.TotalEstimatedSizeBytes())
	})

	t.Run("without package info", func(t *testing.T) {
		config, err := installation.NewInstallationConfiguration(
			[]installation.ComponentSelection{
				mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			},
			nil,
			mustCreateDiskSpace(t, 20*installation.GB, 10*installation.GB),
			false,
		)
		require.NoError(t, err)

		assert.Equal(t, uint64(0), config.TotalEstimatedSizeBytes())
	})
}

func TestInstallationConfiguration_String(t *testing.T) {
	config, err := installation.NewInstallationConfiguration(
		[]installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.20"),
		},
		mustCreateGPUSupport(t, "amd", true, installation.ComponentAMDDriver),
		mustCreateDiskSpace(t, 50*installation.GB, 10*installation.GB),
		true,
	)
	require.NoError(t, err)

	str := config.String()
	assert.Contains(t, str, "2 components")
	assert.Contains(t, str, "GPU: amd")
}

func TestInstallationConfiguration_ValueObjectImmutability(t *testing.T) {
	// Components slice should be copied to prevent external modification
	components := []installation.ComponentSelection{
		mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
	}

	config, err := installation.NewInstallationConfiguration(
		components,
		nil,
		mustCreateDiskSpace(t, 20*installation.GB, 10*installation.GB),
		false,
	)
	require.NoError(t, err)

	// Modify original slice
	components[0] = mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.20")

	// Config should not be affected
	assert.Equal(t, installation.ComponentHyprland, config.Components()[0].Component())
}

func TestInstallationConfiguration_DuplicateComponents(t *testing.T) {
	// Should allow duplicate components (different versions)
	config, err := installation.NewInstallationConfiguration(
		[]installation.ComponentSelection{
			mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.20"),
			mustCreateComponentSelection(t, installation.ComponentWaybar, "0.9.21"), // Duplicate
		},
		nil,
		mustCreateDiskSpace(t, 20*installation.GB, 10*installation.GB),
		false,
	)

	// For now, allow duplicates (business logic can handle deduplication)
	require.NoError(t, err)
	assert.Equal(t, 3, config.ComponentCount())
}

// Helper functions for creating valid value objects in tests
func mustCreateComponentSelection(t *testing.T, component installation.ComponentName, version string) installation.ComponentSelection {
	t.Helper()
	sel, err := installation.NewComponentSelection(component, version, nil)
	require.NoError(t, err)
	return sel
}

func mustCreateGPUSupport(t *testing.T, vendor string, requiresDriver bool, driver installation.ComponentName) *installation.GPUSupport {
	t.Helper()
	gpu, err := installation.NewGPUSupport(vendor, requiresDriver, driver)
	require.NoError(t, err)
	return &gpu
}

func mustCreateDiskSpace(t *testing.T, available, required uint64) installation.DiskSpace {
	t.Helper()
	ds, err := installation.NewDiskSpace(available, required)
	require.NoError(t, err)
	return ds
}
