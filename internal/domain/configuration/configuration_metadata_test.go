package configuration_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigurationMetadata(t *testing.T) {
	tests := []struct {
		name        string
		configName  string
		description string
		author      string
		tags        []string
		category    configuration.ConfigurationCategory
		wantErr     error
	}{
		{
			name:        "valid metadata with all fields",
			configName:  "Dev Stack",
			description: "Development environment setup",
			author:      "John Doe",
			tags:        []string{"development", "hyprland"},
			category:    configuration.CategoryDevelopment,
			wantErr:     nil,
		},
		{
			name:        "valid metadata minimal",
			configName:  "Minimal",
			description: "Minimal setup",
			author:      "",
			tags:        nil,
			category:    configuration.CategoryCustom,
			wantErr:     nil,
		},
		{
			name:        "empty description is valid",
			configName:  "No Desc",
			description: "",
			author:      "Jane",
			tags:        []string{"test"},
			category:    configuration.CategoryTesting,
			wantErr:     nil,
		},
		{
			name:        "description too long",
			configName:  "Test",
			description: string(make([]byte, 1001)),
			author:      "Test",
			tags:        nil,
			category:    configuration.CategoryCustom,
			wantErr:     configuration.ErrDescriptionTooLong,
		},
		{
			name:        "too many tags",
			configName:  "Test",
			description: "Test",
			author:      "Test",
			tags:        []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
			category:    configuration.CategoryCustom,
			wantErr:     configuration.ErrTooManyTags,
		},
		{
			name:        "invalid tag - empty",
			configName:  "Test",
			description: "Test",
			author:      "Test",
			tags:        []string{"valid", ""},
			category:    configuration.CategoryCustom,
			wantErr:     configuration.ErrInvalidTag,
		},
		{
			name:        "invalid configuration name",
			configName:  "",
			description: "Test",
			author:      "Test",
			tags:        nil,
			category:    configuration.CategoryCustom,
			wantErr:     configuration.ErrInvalidConfigurationName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := configuration.NewConfigurationMetadata(
				tt.configName,
				tt.description,
				tt.author,
				tt.tags,
				tt.category,
			)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.configName, metadata.Name().String())
			assert.Equal(t, tt.description, metadata.Description())
			assert.Equal(t, tt.author, metadata.Author())
			assert.Equal(t, tt.category, metadata.Category())
		})
	}
}

func TestConfigurationMetadata_Tags(t *testing.T) {
	t.Run("tags are defensively copied", func(t *testing.T) {
		originalTags := []string{"tag1", "tag2"}
		metadata, err := configuration.NewConfigurationMetadata(
			"Test",
			"Description",
			"Author",
			originalTags,
			configuration.CategoryCustom,
		)
		require.NoError(t, err)

		// Modify original
		originalTags[0] = "modified"

		// Metadata should be unchanged
		tags := metadata.Tags()
		assert.Equal(t, "tag1", tags[0])

		// Modify returned tags
		tags[1] = "modified"

		// Get tags again, should be unchanged
		tags2 := metadata.Tags()
		assert.Equal(t, "tag2", tags2[1])
	})

	t.Run("duplicate tags are removed", func(t *testing.T) {
		metadata, err := configuration.NewConfigurationMetadata(
			"Test",
			"Description",
			"Author",
			[]string{"tag1", "tag2", "tag1", "tag3", "tag2"},
			configuration.CategoryCustom,
		)
		require.NoError(t, err)

		tags := metadata.Tags()
		assert.Len(t, tags, 3)
		assert.Contains(t, tags, "tag1")
		assert.Contains(t, tags, "tag2")
		assert.Contains(t, tags, "tag3")
	})
}

func TestConfigurationCategory(t *testing.T) {
	tests := []struct {
		category configuration.ConfigurationCategory
		want     string
	}{
		{configuration.CategoryDevelopment, "Development"},
		{configuration.CategoryProduction, "Production"},
		{configuration.CategoryTesting, "Testing"},
		{configuration.CategoryCustom, "Custom"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.category.String())
		})
	}
}
