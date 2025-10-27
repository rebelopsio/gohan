package repository_test

import (
	"context"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemorySessionRepository_Save(t *testing.T) {
	t.Run("saves session successfully", func(t *testing.T) {
		repo := repository.NewMemorySessionRepository()
		ctx := context.Background()

		session := createTestSession(t)

		err := repo.Save(ctx, session)

		require.NoError(t, err)
		assert.Equal(t, 1, repo.Count())
	})

	t.Run("updates existing session", func(t *testing.T) {
		repo := repository.NewMemorySessionRepository()
		ctx := context.Background()

		session := createTestSession(t)

		// Save once
		err := repo.Save(ctx, session)
		require.NoError(t, err)

		// Save again (update)
		err = repo.Save(ctx, session)
		require.NoError(t, err)

		// Should still only have one session
		assert.Equal(t, 1, repo.Count())
	})

	t.Run("saves multiple sessions", func(t *testing.T) {
		repo := repository.NewMemorySessionRepository()
		ctx := context.Background()

		session1 := createTestSession(t)
		session2 := createTestSession(t)

		err := repo.Save(ctx, session1)
		require.NoError(t, err)

		err = repo.Save(ctx, session2)
		require.NoError(t, err)

		assert.Equal(t, 2, repo.Count())
	})
}

func TestMemorySessionRepository_FindByID(t *testing.T) {
	t.Run("finds existing session", func(t *testing.T) {
		repo := repository.NewMemorySessionRepository()
		ctx := context.Background()

		session := createTestSession(t)
		err := repo.Save(ctx, session)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, session.ID())

		require.NoError(t, err)
		assert.Equal(t, session.ID(), found.ID())
	})

	t.Run("returns error for non-existent session", func(t *testing.T) {
		repo := repository.NewMemorySessionRepository()
		ctx := context.Background()

		_, err := repo.FindByID(ctx, "non-existent")

		assert.ErrorIs(t, err, installation.ErrSessionNotFound)
	})

	t.Run("finds correct session among multiple", func(t *testing.T) {
		repo := repository.NewMemorySessionRepository()
		ctx := context.Background()

		session1 := createTestSession(t)
		session2 := createTestSession(t)

		err := repo.Save(ctx, session1)
		require.NoError(t, err)
		err = repo.Save(ctx, session2)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, session2.ID())

		require.NoError(t, err)
		assert.Equal(t, session2.ID(), found.ID())
	})
}

func TestMemorySessionRepository_Clear(t *testing.T) {
	t.Run("clears all sessions", func(t *testing.T) {
		repo := repository.NewMemorySessionRepository()
		ctx := context.Background()

		session1 := createTestSession(t)
		session2 := createTestSession(t)

		err := repo.Save(ctx, session1)
		require.NoError(t, err)
		err = repo.Save(ctx, session2)
		require.NoError(t, err)

		assert.Equal(t, 2, repo.Count())

		repo.Clear()

		assert.Equal(t, 0, repo.Count())
	})
}

func TestMemorySessionRepository_ThreadSafety(t *testing.T) {
	t.Run("concurrent saves are thread-safe", func(t *testing.T) {
		repo := repository.NewMemorySessionRepository()
		ctx := context.Background()

		const goroutines = 10
		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				session := createTestSession(t)
				repo.Save(ctx, session)
				done <- true
			}()
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}

		assert.Equal(t, goroutines, repo.Count())
	})

	t.Run("concurrent reads are thread-safe", func(t *testing.T) {
		repo := repository.NewMemorySessionRepository()
		ctx := context.Background()

		session := createTestSession(t)
		err := repo.Save(ctx, session)
		require.NoError(t, err)

		const goroutines = 10
		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				_, err := repo.FindByID(ctx, session.ID())
				assert.NoError(t, err)
				done <- true
			}()
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})
}

// Helper function to create a test session
func createTestSession(t *testing.T) *installation.InstallationSession {
	t.Helper()

	components := createTestComponents(t)
	diskSpace, err := installation.NewDiskSpace(
		100*uint64(installation.GB),
		10*uint64(installation.GB),
	)
	require.NoError(t, err)

	config, err := installation.NewInstallationConfiguration(
		components,
		nil,
		diskSpace,
		false,
	)
	require.NoError(t, err)

	session, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	return session
}

func createTestComponents(t *testing.T) []installation.ComponentSelection {
	t.Helper()

	pkg, err := installation.NewPackageInfo(
		"hyprland",
		"0.35.0",
		50*uint64(installation.MB),
		nil,
	)
	require.NoError(t, err)

	component, err := installation.NewComponentSelection(
		installation.ComponentHyprland,
		"0.35.0",
		&pkg,
	)
	require.NoError(t, err)

	return []installation.ComponentSelection{component}
}
