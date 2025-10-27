package configuration_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/configuration"
	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create a valid manifest
func createValidManifest(t *testing.T) configuration.ConfigurationManifest {
	manifest, err := configuration.NewConfigurationManifest(
		[]installation.ComponentSelection{
			createComponentSelection(t, installation.ComponentHyprland, "0.32.0"),
			createComponentSelection(t, installation.ComponentWaybar, "0.9.0"),
		},
		2000000000,
		false,
	)
	require.NoError(t, err)
	return manifest
}

// Helper to create valid metadata
func createValidMetadata(t *testing.T, name string) configuration.ConfigurationMetadata {
	metadata, err := configuration.NewConfigurationMetadata(
		name,
		"Test configuration",
		"Test Author",
		[]string{"test", "dev"},
		configuration.CategoryDevelopment,
	)
	require.NoError(t, err)
	return metadata
}

func TestNewConfigurationTemplate(t *testing.T) {
	tests := []struct {
		name     string
		metadata configuration.ConfigurationMetadata
		manifest configuration.ConfigurationManifest
		wantErr  bool
	}{
		{
			name:     "valid template",
			metadata: createValidMetadata(t, "Test Config"),
			manifest: createValidManifest(t),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := configuration.NewConfigurationTemplate(
				tt.metadata,
				tt.manifest,
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, template.ID())
			assert.Equal(t, tt.metadata.Name().String(), template.Metadata().Name().String())
			assert.Equal(t, tt.manifest.ComponentCount(), template.Manifest().ComponentCount())
			assert.False(t, template.CreatedAt().IsZero())
			assert.Equal(t, 1, template.Version())
		})
	}
}

func TestConfigurationTemplate_ID(t *testing.T) {
	template1, err := configuration.NewConfigurationTemplate(
		createValidMetadata(t, "Config 1"),
		createValidManifest(t),
	)
	require.NoError(t, err)

	template2, err := configuration.NewConfigurationTemplate(
		createValidMetadata(t, "Config 2"),
		createValidManifest(t),
	)
	require.NoError(t, err)

	// Each template should have unique ID
	assert.NotEqual(t, template1.ID(), template2.ID())
	assert.NotEmpty(t, template1.ID())
	assert.NotEmpty(t, template2.ID())
}

func TestConfigurationTemplate_Metadata(t *testing.T) {
	metadata := createValidMetadata(t, "Test Config")
	template, err := configuration.NewConfigurationTemplate(
		metadata,
		createValidManifest(t),
	)
	require.NoError(t, err)

	retrievedMetadata := template.Metadata()
	assert.Equal(t, metadata.Name().String(), retrievedMetadata.Name().String())
	assert.Equal(t, metadata.Description(), retrievedMetadata.Description())
	assert.Equal(t, metadata.Author(), retrievedMetadata.Author())
	assert.Equal(t, metadata.Category(), retrievedMetadata.Category())
}

func TestConfigurationTemplate_Manifest(t *testing.T) {
	manifest := createValidManifest(t)
	template, err := configuration.NewConfigurationTemplate(
		createValidMetadata(t, "Test Config"),
		manifest,
	)
	require.NoError(t, err)

	retrievedManifest := template.Manifest()
	assert.Equal(t, manifest.ComponentCount(), retrievedManifest.ComponentCount())
	assert.Equal(t, manifest.DiskRequiredBytes(), retrievedManifest.DiskRequiredBytes())
	assert.Equal(t, manifest.GPURequired(), retrievedManifest.GPURequired())
}

func TestConfigurationTemplate_CreatedAt(t *testing.T) {
	before := time.Now()
	template, err := configuration.NewConfigurationTemplate(
		createValidMetadata(t, "Test Config"),
		createValidManifest(t),
	)
	require.NoError(t, err)
	after := time.Now()

	createdAt := template.CreatedAt()
	assert.False(t, createdAt.IsZero())
	assert.True(t, createdAt.After(before) || createdAt.Equal(before))
	assert.True(t, createdAt.Before(after) || createdAt.Equal(after))
}

func TestConfigurationTemplate_Version(t *testing.T) {
	t.Run("new template starts at version 1", func(t *testing.T) {
		template, err := configuration.NewConfigurationTemplate(
			createValidMetadata(t, "Test Config"),
			createValidManifest(t),
		)
		require.NoError(t, err)

		assert.Equal(t, 1, template.Version())
	})
}

func TestConfigurationTemplate_Age(t *testing.T) {
	template, err := configuration.NewConfigurationTemplate(
		createValidMetadata(t, "Test Config"),
		createValidManifest(t),
	)
	require.NoError(t, err)

	// Small delay to ensure age is measurable
	time.Sleep(10 * time.Millisecond)

	age := template.Age()
	assert.True(t, age > 0)
	assert.True(t, age < time.Second) // Should be very recent
}

func TestConfigurationTemplate_String(t *testing.T) {
	template, err := configuration.NewConfigurationTemplate(
		createValidMetadata(t, "My Dev Config"),
		createValidManifest(t),
	)
	require.NoError(t, err)

	str := template.String()
	assert.Contains(t, str, "My Dev Config")
	assert.Contains(t, str, "Development")
	assert.NotEmpty(t, str)
}
