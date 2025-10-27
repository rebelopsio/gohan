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

func TestListInstallationsUseCase_Execute(t *testing.T) {
	t.Run("returns empty list when no sessions exist", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewListInstallationsUseCase(sessionRepo)
		ctx := context.Background()

		response, err := useCase.Execute(ctx)

		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.Empty(t, response.Sessions)
		assert.Equal(t, 0, response.TotalCount)
	})

	t.Run("returns all sessions", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewListInstallationsUseCase(sessionRepo)
		ctx := context.Background()

		// Create multiple sessions
		components1, err := createTestComponents()
		require.NoError(t, err)

		diskSpace1, err := installation.NewDiskSpace(
			100*uint64(installation.GB),
			10*uint64(installation.GB),
		)
		require.NoError(t, err)

		config1, err := installation.NewInstallationConfiguration(
			components1,
			nil,
			diskSpace1,
			false,
		)
		require.NoError(t, err)

		session1, err := installation.NewInstallationSession(config1)
		require.NoError(t, err)

		err = sessionRepo.Save(ctx, session1)
		require.NoError(t, err)

		// Create second session
		components2, err := createTestComponents()
		require.NoError(t, err)

		diskSpace2, err := installation.NewDiskSpace(
			100*uint64(installation.GB),
			10*uint64(installation.GB),
		)
		require.NoError(t, err)

		config2, err := installation.NewInstallationConfiguration(
			components2,
			nil,
			diskSpace2,
			false,
		)
		require.NoError(t, err)

		session2, err := installation.NewInstallationSession(config2)
		require.NoError(t, err)

		err = sessionRepo.Save(ctx, session2)
		require.NoError(t, err)

		// List all sessions
		response, err := useCase.Execute(ctx)

		require.NoError(t, err)
		assert.Equal(t, 2, response.TotalCount)
		assert.Len(t, response.Sessions, 2)

		// Verify session IDs are present
		sessionIDs := []string{response.Sessions[0].SessionID, response.Sessions[1].SessionID}
		assert.Contains(t, sessionIDs, session1.ID())
		assert.Contains(t, sessionIDs, session2.ID())
	})

	t.Run("includes session status in response", func(t *testing.T) {
		sessionRepo := repository.NewMemorySessionRepository()
		useCase := usecases.NewListInstallationsUseCase(sessionRepo)
		ctx := context.Background()

		// Create a session
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

		// List sessions
		response, err := useCase.Execute(ctx)

		require.NoError(t, err)
		assert.Equal(t, 1, response.TotalCount)
		assert.Len(t, response.Sessions, 1)

		// Verify session details
		sessionInfo := response.Sessions[0]
		assert.Equal(t, session.ID(), sessionInfo.SessionID)
		assert.Equal(t, "pending", sessionInfo.Status)
		assert.Equal(t, 1, sessionInfo.ComponentsTotal)
	})
}
