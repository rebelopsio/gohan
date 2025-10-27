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

func TestGetInstallationStatusUseCase_Execute(t *testing.T) {
	t.Run("successfully get status for pending session", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewGetInstallationStatusUseCase(sessionRepo)
		ctx := context.Background()

		// Create components and disk space
		components, err := createTestComponents()
		require.NoError(t, err)

		diskSpace, err := installation.NewDiskSpace(
			100*uint64(installation.GB),
			10*uint64(installation.GB),
		)
		require.NoError(t, err)

		// Create configuration
		config, err := installation.NewInstallationConfiguration(
			components,
			nil,
			diskSpace,
			false,
		)
		require.NoError(t, err)

		// Create a session
		session, err := installation.NewInstallationSession(config)
		require.NoError(t, err)

		err = sessionRepo.Save(ctx, session)
		require.NoError(t, err)

		// Get status
		response, err := useCase.Execute(ctx, session.ID())

		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, session.ID(), response.SessionID)
		assert.Equal(t, "pending", response.Status)
		assert.Equal(t, "pending", response.CurrentPhase)
		assert.Equal(t, 0, response.ComponentsInstalled)
		assert.Equal(t, 1, response.ComponentsTotal)
	})

	t.Run("session not found", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewGetInstallationStatusUseCase(sessionRepo)
		ctx := context.Background()

		_, err := useCase.Execute(ctx, "nonexistent")

		assert.Error(t, err)
		assert.ErrorIs(t, err, installation.ErrSessionNotFound)
	})
}
