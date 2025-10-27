package configuration_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReconstructConfigurationTemplate(t *testing.T) {
	t.Run("reconstructs template successfully", func(t *testing.T) {
		// Arrange
		id := "test-template-123"
		metadata := createValidMetadata(t, "Test Config")
		manifest := createValidManifest(t)
		createdAt := time.Now().Add(-24 * time.Hour)
		version := 1

		// Act
		template, err := configuration.ReconstructConfigurationTemplate(
			id,
			metadata,
			manifest,
			createdAt,
			version,
		)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, id, template.ID())
		assert.Equal(t, metadata.Name().String(), template.Metadata().Name().String())
		assert.Equal(t, manifest.ComponentCount(), template.Manifest().ComponentCount())
		assert.Equal(t, createdAt.Unix(), template.CreatedAt().Unix())
		assert.Equal(t, version, template.Version())
	})

	t.Run("reconstructs template with higher version", func(t *testing.T) {
		// Arrange
		id := "test-template-456"
		metadata := createValidMetadata(t, "Test Config v5")
		manifest := createValidManifest(t)
		createdAt := time.Now().Add(-48 * time.Hour)
		version := 5

		// Act
		template, err := configuration.ReconstructConfigurationTemplate(
			id,
			metadata,
			manifest,
			createdAt,
			version,
		)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, version, template.Version())
	})

	t.Run("rejects empty ID", func(t *testing.T) {
		// Arrange
		id := ""
		metadata := createValidMetadata(t, "Test Config")
		manifest := createValidManifest(t)
		createdAt := time.Now()
		version := 1

		// Act
		template, err := configuration.ReconstructConfigurationTemplate(
			id,
			metadata,
			manifest,
			createdAt,
			version,
		)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, template)
		assert.Contains(t, err.Error(), "ID cannot be empty")
	})

	t.Run("rejects zero created time", func(t *testing.T) {
		// Arrange
		id := "test-template-789"
		metadata := createValidMetadata(t, "Test Config")
		manifest := createValidManifest(t)
		createdAt := time.Time{} // Zero time
		version := 1

		// Act
		template, err := configuration.ReconstructConfigurationTemplate(
			id,
			metadata,
			manifest,
			createdAt,
			version,
		)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, template)
		assert.Contains(t, err.Error(), "created time cannot be zero")
	})

	t.Run("rejects invalid version", func(t *testing.T) {
		// Arrange
		id := "test-template-999"
		metadata := createValidMetadata(t, "Test Config")
		manifest := createValidManifest(t)
		createdAt := time.Now()
		version := 0 // Invalid version

		// Act
		template, err := configuration.ReconstructConfigurationTemplate(
			id,
			metadata,
			manifest,
			createdAt,
			version,
		)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, template)
		assert.Contains(t, err.Error(), "version must be positive")
	})
}
