package repository_test

import (
	"context"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/configuration"
	"github.com/rebelopsio/gohan/internal/domain/installation"
	configrepo "github.com/rebelopsio/gohan/internal/infrastructure/configuration/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create a valid template
func createTestTemplate(t *testing.T, name string) *configuration.ConfigurationTemplate {
	// Create manifest
	compSel, err := installation.NewComponentSelection(installation.ComponentHyprland, "0.32.0", nil)
	require.NoError(t, err)
	manifest, err := configuration.NewConfigurationManifest(
		[]installation.ComponentSelection{compSel},
		1000000000,
		false,
	)
	require.NoError(t, err)

	// Create metadata
	metadata, err := configuration.NewConfigurationMetadata(
		name,
		"Test configuration",
		"Test Author",
		[]string{"test"},
		configuration.CategoryDevelopment,
	)
	require.NoError(t, err)

	// Create template
	template, err := configuration.NewConfigurationTemplate(metadata, manifest)
	require.NoError(t, err)

	return template
}

func TestMemoryRepository_Save(t *testing.T) {
	t.Run("saves template successfully", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		template := createTestTemplate(t, "Test Config")

		err := repo.Save(ctx, template)
		require.NoError(t, err)

		// Verify it was saved
		found, err := repo.FindByID(ctx, template.ID())
		require.NoError(t, err)
		assert.Equal(t, template.ID(), found.ID())
	})

	t.Run("updates existing template", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		template := createTestTemplate(t, "Test Config")

		// Save first time
		err := repo.Save(ctx, template)
		require.NoError(t, err)

		// Save again (update)
		err = repo.Save(ctx, template)
		require.NoError(t, err)

		// Should still have only one template
		templates, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Len(t, templates, 1)
	})

	t.Run("saves multiple templates", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		template1 := createTestTemplate(t, "Config 1")
		template2 := createTestTemplate(t, "Config 2")

		err := repo.Save(ctx, template1)
		require.NoError(t, err)

		err = repo.Save(ctx, template2)
		require.NoError(t, err)

		templates, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Len(t, templates, 2)
	})
}

func TestMemoryRepository_FindByID(t *testing.T) {
	t.Run("finds existing template", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		template := createTestTemplate(t, "Test Config")
		err := repo.Save(ctx, template)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, template.ID())
		require.NoError(t, err)
		assert.Equal(t, template.ID(), found.ID())
		assert.Equal(t, template.Metadata().Name().String(), found.Metadata().Name().String())
	})

	t.Run("returns error for non-existent template", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		_, err := repo.FindByID(ctx, "non-existent-id")
		assert.ErrorIs(t, err, configuration.ErrTemplateNotFound)
	})

	t.Run("finds correct template among multiple", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		template1 := createTestTemplate(t, "Config 1")
		template2 := createTestTemplate(t, "Config 2")

		err := repo.Save(ctx, template1)
		require.NoError(t, err)
		err = repo.Save(ctx, template2)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, template2.ID())
		require.NoError(t, err)
		assert.Equal(t, template2.ID(), found.ID())
	})
}

func TestMemoryRepository_FindByName(t *testing.T) {
	t.Run("finds template by name", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		template := createTestTemplate(t, "My Dev Config")
		err := repo.Save(ctx, template)
		require.NoError(t, err)

		found, err := repo.FindByName(ctx, "My Dev Config")
		require.NoError(t, err)
		assert.Equal(t, template.ID(), found.ID())
	})

	t.Run("returns error for non-existent name", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		_, err := repo.FindByName(ctx, "Non-existent Config")
		assert.ErrorIs(t, err, configuration.ErrTemplateNotFound)
	})
}

func TestMemoryRepository_ExistsByName(t *testing.T) {
	t.Run("returns true for existing name", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		template := createTestTemplate(t, "Existing Config")
		err := repo.Save(ctx, template)
		require.NoError(t, err)

		exists, err := repo.ExistsByName(ctx, "Existing Config")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns false for non-existent name", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		exists, err := repo.ExistsByName(ctx, "Non-existent Config")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestMemoryRepository_List(t *testing.T) {
	t.Run("lists all templates", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		template1 := createTestTemplate(t, "Config 1")
		template2 := createTestTemplate(t, "Config 2")
		template3 := createTestTemplate(t, "Config 3")

		err := repo.Save(ctx, template1)
		require.NoError(t, err)
		err = repo.Save(ctx, template2)
		require.NoError(t, err)
		err = repo.Save(ctx, template3)
		require.NoError(t, err)

		templates, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Len(t, templates, 3)
	})

	t.Run("returns empty list when no templates", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		templates, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, templates)
	})

	t.Run("returns templates ordered by creation date", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		template1 := createTestTemplate(t, "Config 1")
		template2 := createTestTemplate(t, "Config 2")

		// Save in order
		err := repo.Save(ctx, template1)
		require.NoError(t, err)
		err = repo.Save(ctx, template2)
		require.NoError(t, err)

		templates, err := repo.List(ctx)
		require.NoError(t, err)

		// Most recent first
		assert.Equal(t, template2.ID(), templates[0].ID())
		assert.Equal(t, template1.ID(), templates[1].ID())
	})
}

func TestMemoryRepository_ListByCategory(t *testing.T) {
	t.Run("filters templates by category", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		// Create templates with different categories
		template1 := createTestTemplate(t, "Dev Config")
		template2 := createTestTemplate(t, "Prod Config")

		err := repo.Save(ctx, template1)
		require.NoError(t, err)
		err = repo.Save(ctx, template2)
		require.NoError(t, err)

		// Both are development category
		templates, err := repo.ListByCategory(ctx, configuration.CategoryDevelopment)
		require.NoError(t, err)
		assert.Len(t, templates, 2)

		// None are production
		templates, err = repo.ListByCategory(ctx, configuration.CategoryProduction)
		require.NoError(t, err)
		assert.Empty(t, templates)
	})
}

func TestMemoryRepository_ListByTag(t *testing.T) {
	t.Run("filters templates by tag", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		template := createTestTemplate(t, "Tagged Config")
		err := repo.Save(ctx, template)
		require.NoError(t, err)

		// Should find by existing tag
		templates, err := repo.ListByTag(ctx, "test")
		require.NoError(t, err)
		assert.Len(t, templates, 1)

		// Should not find non-existent tag
		templates, err = repo.ListByTag(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Empty(t, templates)
	})
}

func TestMemoryRepository_Delete(t *testing.T) {
	t.Run("deletes existing template", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		template := createTestTemplate(t, "Test Config")
		err := repo.Save(ctx, template)
		require.NoError(t, err)

		err = repo.Delete(ctx, template.ID())
		require.NoError(t, err)

		// Should no longer exist
		_, err = repo.FindByID(ctx, template.ID())
		assert.ErrorIs(t, err, configuration.ErrTemplateNotFound)
	})

	t.Run("returns error when deleting non-existent template", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		err := repo.Delete(ctx, "non-existent-id")
		assert.ErrorIs(t, err, configuration.ErrTemplateNotFound)
	})
}

func TestMemoryRepository_ThreadSafety(t *testing.T) {
	t.Run("concurrent saves are thread-safe", func(t *testing.T) {
		repo := configrepo.NewMemoryRepository()
		ctx := context.Background()

		done := make(chan bool)

		// Concurrent saves
		for i := 0; i < 10; i++ {
			go func(n int) {
				template := createTestTemplate(t, "Config "+string(rune(n)))
				_ = repo.Save(ctx, template)
				done <- true
			}(i)
		}

		// Wait for all to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		templates, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Len(t, templates, 10)
	})
}
