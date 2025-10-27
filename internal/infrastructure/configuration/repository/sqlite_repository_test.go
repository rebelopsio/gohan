package repository_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/configuration"
	configrepo "github.com/rebelopsio/gohan/internal/infrastructure/configuration/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create a test SQLite database
func setupTestSQLiteDB(t *testing.T) *configrepo.SQLiteRepository {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	repo, err := configrepo.NewSQLiteRepository(dbPath)
	require.NoError(t, err)

	return repo
}

func TestSQLiteRepository_Save(t *testing.T) {
	t.Run("saves template successfully", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

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
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()
		template := createTestTemplate(t, "Test Config")

		// Save first time
		err := repo.Save(ctx, template)
		require.NoError(t, err)

		// Save again (update)
		err = repo.Save(ctx, template)
		require.NoError(t, err)

		// Should still have only one template
		count, err := repo.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestSQLiteRepository_FindByID(t *testing.T) {
	t.Run("finds existing template", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()
		template := createTestTemplate(t, "Test Config")

		err := repo.Save(ctx, template)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, template.ID())
		require.NoError(t, err)
		assert.Equal(t, template.ID(), found.ID())
		assert.Equal(t, template.Metadata().Name().String(), found.Metadata().Name().String())
		assert.Equal(t, template.Version(), found.Version())
	})

	t.Run("returns error for non-existent template", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()

		_, err := repo.FindByID(ctx, "non-existent-id")
		assert.Error(t, err)
	})

	t.Run("reconstructs template with all fields", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()
		template := createTestTemplate(t, "Full Config")

		err := repo.Save(ctx, template)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, template.ID())
		require.NoError(t, err)

		// Verify metadata
		assert.Equal(t, template.Metadata().Name().String(), found.Metadata().Name().String())
		assert.Equal(t, template.Metadata().Description(), found.Metadata().Description())
		assert.Equal(t, template.Metadata().Author(), found.Metadata().Author())
		assert.Equal(t, template.Metadata().Category(), found.Metadata().Category())

		// Verify manifest
		assert.Equal(t, template.Manifest().ComponentCount(), found.Manifest().ComponentCount())
		assert.Equal(t, template.Manifest().DiskRequiredBytes(), found.Manifest().DiskRequiredBytes())
		assert.Equal(t, template.Manifest().GPURequired(), found.Manifest().GPURequired())

		// Verify timestamps
		assert.Equal(t, template.CreatedAt().Unix(), found.CreatedAt().Unix())
		assert.Equal(t, template.Version(), found.Version())
	})
}

func TestSQLiteRepository_FindByName(t *testing.T) {
	t.Run("finds template by name", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()
		template := createTestTemplate(t, "My Dev Config")

		err := repo.Save(ctx, template)
		require.NoError(t, err)

		found, err := repo.FindByName(ctx, "My Dev Config")
		require.NoError(t, err)
		assert.Equal(t, template.ID(), found.ID())
	})

	t.Run("returns error for non-existent name", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()

		_, err := repo.FindByName(ctx, "Non-existent Config")
		assert.Error(t, err)
	})
}

func TestSQLiteRepository_ExistsByName(t *testing.T) {
	t.Run("returns true for existing name", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()
		template := createTestTemplate(t, "Existing Config")

		err := repo.Save(ctx, template)
		require.NoError(t, err)

		exists, err := repo.ExistsByName(ctx, "Existing Config")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns false for non-existent name", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()

		exists, err := repo.ExistsByName(ctx, "Non-existent Config")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestSQLiteRepository_List(t *testing.T) {
	t.Run("lists all templates", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

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

	t.Run("returns empty list when no templates", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()

		templates, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, templates)
	})
}

func TestSQLiteRepository_ListByCategory(t *testing.T) {
	t.Run("filters templates by category", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()

		template1 := createTestTemplate(t, "Dev Config")
		template2 := createTestTemplate(t, "Another Dev")

		err := repo.Save(ctx, template1)
		require.NoError(t, err)
		err = repo.Save(ctx, template2)
		require.NoError(t, err)

		templates, err := repo.ListByCategory(ctx, configuration.CategoryDevelopment)
		require.NoError(t, err)
		assert.Len(t, templates, 2)

		templates, err = repo.ListByCategory(ctx, configuration.CategoryProduction)
		require.NoError(t, err)
		assert.Empty(t, templates)
	})
}

func TestSQLiteRepository_ListByTag(t *testing.T) {
	t.Run("filters templates by tag", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()

		template := createTestTemplate(t, "Tagged Config")

		err := repo.Save(ctx, template)
		require.NoError(t, err)

		templates, err := repo.ListByTag(ctx, "test")
		require.NoError(t, err)
		assert.Len(t, templates, 1)

		templates, err = repo.ListByTag(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Empty(t, templates)
	})
}

func TestSQLiteRepository_Delete(t *testing.T) {
	t.Run("deletes existing template", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()
		template := createTestTemplate(t, "Test Config")

		err := repo.Save(ctx, template)
		require.NoError(t, err)

		err = repo.Delete(ctx, template.ID())
		require.NoError(t, err)

		_, err = repo.FindByID(ctx, template.ID())
		assert.Error(t, err)
	})

	t.Run("returns error when deleting non-existent template", func(t *testing.T) {
		repo := setupTestSQLiteDB(t)
		defer repo.Close()

		ctx := context.Background()

		err := repo.Delete(ctx, "non-existent-id")
		assert.Error(t, err)
	})
}
