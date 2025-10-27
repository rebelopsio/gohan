package configuration_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/configuration"
	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create component selections for testing
func createComponentSelection(t *testing.T, component installation.ComponentName, version string) installation.ComponentSelection {
	cs, err := installation.NewComponentSelection(component, version, nil)
	require.NoError(t, err)
	return cs
}

func TestNewConfigurationManifest(t *testing.T) {
	tests := []struct {
		name              string
		components        []installation.ComponentSelection
		diskRequiredBytes uint64
		gpuRequired       bool
		wantErr           error
	}{
		{
			name: "valid manifest with core component",
			components: []installation.ComponentSelection{
				createComponentSelection(t, installation.ComponentHyprland, "0.32.0"),
			},
			diskRequiredBytes: 1000000000,
			gpuRequired:       false,
			wantErr:           nil,
		},
		{
			name: "valid manifest with multiple components",
			components: []installation.ComponentSelection{
				createComponentSelection(t, installation.ComponentHyprland, "0.32.0"),
				createComponentSelection(t, installation.ComponentWaybar, "0.9.0"),
			},
			diskRequiredBytes: 2000000000,
			gpuRequired:       true,
			wantErr:           nil,
		},
		{
			name:              "empty components list",
			components:        []installation.ComponentSelection{},
			diskRequiredBytes: 1000000000,
			gpuRequired:       false,
			wantErr:           configuration.ErrNoComponents,
		},
		{
			name:              "nil components list",
			components:        nil,
			diskRequiredBytes: 1000000000,
			gpuRequired:       false,
			wantErr:           configuration.ErrNoComponents,
		},
		{
			name: "missing core hyprland component",
			components: []installation.ComponentSelection{
				createComponentSelection(t, installation.ComponentWaybar, "0.9.0"),
			},
			diskRequiredBytes: 1000000000,
			gpuRequired:       false,
			wantErr:           configuration.ErrMissingCoreComponent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest, err := configuration.NewConfigurationManifest(
				tt.components,
				tt.diskRequiredBytes,
				tt.gpuRequired,
			)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Len(t, manifest.Components(), len(tt.components))
			assert.Equal(t, tt.diskRequiredBytes, manifest.DiskRequiredBytes())
			assert.Equal(t, tt.gpuRequired, manifest.GPURequired())
		})
	}
}

func TestConfigurationManifest_Components(t *testing.T) {
	t.Run("components are defensively copied", func(t *testing.T) {
		components := []installation.ComponentSelection{
			createComponentSelection(t, installation.ComponentHyprland, "0.32.0"),
		}

		manifest, err := configuration.NewConfigurationManifest(components, 1000000000, false)
		require.NoError(t, err)

		// Get components
		returnedComponents := manifest.Components()

		// Should be a copy, not same reference
		assert.NotSame(t, &components, &returnedComponents)
		assert.Len(t, returnedComponents, 1)

		// Modifying returned slice shouldn't affect manifest
		returnedComponents = append(returnedComponents, createComponentSelection(t, installation.ComponentWaybar, "0.9.0"))

		// Get again, should still have only 1
		components2 := manifest.Components()
		assert.Len(t, components2, 1)
	})

	t.Run("has core component", func(t *testing.T) {
		manifest, err := configuration.NewConfigurationManifest(
			[]installation.ComponentSelection{
				createComponentSelection(t, installation.ComponentHyprland, "0.32.0"),
				createComponentSelection(t, installation.ComponentWaybar, "0.9.0"),
			},
			1000000000,
			false,
		)
		require.NoError(t, err)

		assert.True(t, manifest.HasCoreComponent())
	})
}

func TestConfigurationManifest_ComponentCount(t *testing.T) {
	manifest, err := configuration.NewConfigurationManifest(
		[]installation.ComponentSelection{
			createComponentSelection(t, installation.ComponentHyprland, "0.32.0"),
			createComponentSelection(t, installation.ComponentWaybar, "0.9.0"),
			createComponentSelection(t, installation.ComponentRofi, "1.3.0"),
		},
		1000000000,
		false,
	)
	require.NoError(t, err)

	assert.Equal(t, 3, manifest.ComponentCount())
}

func TestConfigurationManifest_DiskRequirements(t *testing.T) {
	manifest, err := configuration.NewConfigurationManifest(
		[]installation.ComponentSelection{
			createComponentSelection(t, installation.ComponentHyprland, "0.32.0"),
		},
		5000000000, // 5GB
		false,
	)
	require.NoError(t, err)

	assert.Equal(t, uint64(5000000000), manifest.DiskRequiredBytes())
	// Should provide GB conversion
	assert.InDelta(t, 4.66, manifest.DiskRequiredGB(), 0.01)
}
