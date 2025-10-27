package usecases_test

import (
	"context"
	"testing"

	"github.com/rebelopsio/gohan/internal/application/installation/usecases"
	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCancelInstallationUseCase_Execute(t *testing.T) {
	t.Run("successfully cancels a pending installation", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewCancelInstallationUseCase(sessionRepo)
		ctx := context.Background()

		// Create a pending session
		components, err := createTestComponents()
		require.NoError(t, err)

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

		err = sessionRepo.Save(ctx, session)
		require.NoError(t, err)

		// Cancel the session
		err = useCase.Execute(ctx, session.ID())

		require.NoError(t, err)

		// Verify session was cancelled
		cancelledSession, err := sessionRepo.FindByID(ctx, session.ID())
		require.NoError(t, err)
		assert.Equal(t, installation.StatusFailed, cancelledSession.Status())
		assert.Contains(t, cancelledSession.FailureReason(), "cancelled")
	})

	t.Run("returns error when session not found", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewCancelInstallationUseCase(sessionRepo)
		ctx := context.Background()

		err := useCase.Execute(ctx, "nonexistent")

		assert.Error(t, err)
		assert.ErrorIs(t, err, installation.ErrSessionNotFound)
	})

	t.Run("returns error when trying to cancel already failed session", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewCancelInstallationUseCase(sessionRepo)
		ctx := context.Background()

		// Create and fail a session
		components, err := createTestComponents()
		require.NoError(t, err)

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

		// Fail the session
		err = session.Fail("some error")
		require.NoError(t, err)

		err = sessionRepo.Save(ctx, session)
		require.NoError(t, err)

		// Try to cancel already failed session
		err = useCase.Execute(ctx, session.ID())

		assert.Error(t, err)
	})
}
