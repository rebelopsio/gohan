package repository

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create a test database
func setupTestDB(t *testing.T) *SQLiteSimpleSessionRepository {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	repo, err := NewSQLiteSimpleSessionRepository(dbPath)
	require.NoError(t, err)

	return repo
}

// Helper to create a valid test session
func createTestSession(t *testing.T) *installation.InstallationSession {
	// Create component selection
	compSel, err := installation.NewComponentSelection(
		installation.ComponentHyprland,
		"0.32.0",
		nil,
	)
	require.NoError(t, err)

	// Create disk space
	diskSpace, err := installation.NewDiskSpace(500000000, 100000000)
	require.NoError(t, err)

	// Create configuration
	config, err := installation.NewInstallationConfiguration(
		[]installation.ComponentSelection{compSel},
		nil,
		diskSpace,
		false,
	)
	require.NoError(t, err)

	// Create session
	session, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	return session
}

func TestSQLiteSimpleSessionRepository_FindByID(t *testing.T) {
	t.Run("finds existing session", func(t *testing.T) {
		// Arrange
		repo := setupTestDB(t)
		defer repo.Close()

		session := createTestSession(t)
		ctx := context.Background()

		// Save session
		err := repo.Save(ctx, session)
		require.NoError(t, err)

		// Act
		found, err := repo.FindByID(ctx, session.ID())

		// Assert
		require.NoError(t, err)
		assert.Equal(t, session.ID(), found.ID())
		assert.Equal(t, session.Status(), found.Status())
		assert.Equal(t, session.StartedAt().Unix(), found.StartedAt().Unix())
	})

	t.Run("returns error for non-existent session", func(t *testing.T) {
		// Arrange
		repo := setupTestDB(t)
		defer repo.Close()

		ctx := context.Background()

		// Act
		found, err := repo.FindByID(ctx, "non-existent-id")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, found)
	})

	t.Run("finds session with snapshot and components", func(t *testing.T) {
		// Arrange
		repo := setupTestDB(t)
		defer repo.Close()

		session := createTestSession(t)
		ctx := context.Background()

		// Create and attach snapshot
		diskSpace, err := installation.NewDiskSpace(500000000, 100000000)
		require.NoError(t, err)

		snapshot, err := installation.NewSystemSnapshot(
			"/tmp/snapshot",
			diskSpace,
			[]string{"pkg1", "pkg2"},
		)
		require.NoError(t, err)

		err = session.StartPreparation(snapshot)
		require.NoError(t, err)

		// Start installing and add component
		err = session.StartInstalling()
		require.NoError(t, err)

		component, err := installation.NewInstalledComponent(
			installation.ComponentHyprland,
			"0.32.0",
			nil,
		)
		require.NoError(t, err)

		err = session.AddInstalledComponent(component)
		require.NoError(t, err)

		// Save session
		err = repo.Save(ctx, session)
		require.NoError(t, err)

		// Act
		found, err := repo.FindByID(ctx, session.ID())

		// Assert
		require.NoError(t, err)
		assert.Equal(t, session.ID(), found.ID())
		assert.Equal(t, session.Status(), found.Status())
		assert.NotNil(t, found.Snapshot())
		assert.Len(t, found.InstalledComponents(), 1)
	})
}

func TestSQLiteSimpleSessionRepository_List(t *testing.T) {
	t.Run("lists all sessions", func(t *testing.T) {
		// Arrange
		repo := setupTestDB(t)
		defer repo.Close()

		ctx := context.Background()

		// Create and save multiple sessions
		session1 := createTestSession(t)
		session2 := createTestSession(t)

		err := repo.Save(ctx, session1)
		require.NoError(t, err)

		err = repo.Save(ctx, session2)
		require.NoError(t, err)

		// Act
		sessions, err := repo.List(ctx)

		// Assert
		require.NoError(t, err)
		assert.Len(t, sessions, 2)

		// Check IDs are present
		ids := make(map[string]bool)
		for _, s := range sessions {
			ids[s.ID()] = true
		}
		assert.True(t, ids[session1.ID()])
		assert.True(t, ids[session2.ID()])
	})

	t.Run("returns empty list when no sessions", func(t *testing.T) {
		// Arrange
		repo := setupTestDB(t)
		defer repo.Close()

		ctx := context.Background()

		// Act
		sessions, err := repo.List(ctx)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, sessions)
	})

	t.Run("lists sessions with different statuses", func(t *testing.T) {
		// Arrange
		repo := setupTestDB(t)
		defer repo.Close()

		ctx := context.Background()

		// Create pending session
		session1 := createTestSession(t)

		// Create completed session
		session2 := createTestSession(t)
		diskSpace, err := installation.NewDiskSpace(500000000, 100000000)
		require.NoError(t, err)
		snapshot, err := installation.NewSystemSnapshot("/tmp/snapshot", diskSpace, nil)
		require.NoError(t, err)
		err = session2.StartPreparation(snapshot)
		require.NoError(t, err)
		err = session2.StartInstalling()
		require.NoError(t, err)
		component, err := installation.NewInstalledComponent(installation.ComponentHyprland, "0.32.0", nil)
		require.NoError(t, err)
		err = session2.AddInstalledComponent(component)
		require.NoError(t, err)
		err = session2.StartConfiguring()
		require.NoError(t, err)
		err = session2.StartVerifying()
		require.NoError(t, err)
		err = session2.Complete()
		require.NoError(t, err)

		// Create failed session
		session3 := createTestSession(t)
		err = session3.Fail("test failure")
		require.NoError(t, err)

		// Save all
		err = repo.Save(ctx, session1)
		require.NoError(t, err)
		err = repo.Save(ctx, session2)
		require.NoError(t, err)
		err = repo.Save(ctx, session3)
		require.NoError(t, err)

		// Act
		sessions, err := repo.List(ctx)

		// Assert
		require.NoError(t, err)
		assert.Len(t, sessions, 3)

		// Verify statuses are preserved
		statuses := make(map[installation.InstallationStatus]int)
		for _, s := range sessions {
			statuses[s.Status()]++
		}

		assert.Equal(t, 1, statuses[installation.StatusPending])
		assert.Equal(t, 1, statuses[installation.StatusCompleted])
		assert.Equal(t, 1, statuses[installation.StatusFailed])
	})
}
