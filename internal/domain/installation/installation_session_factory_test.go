package installation_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a valid configuration for testing
func createValidConfig(t *testing.T) installation.InstallationConfiguration {
	// Create component selection
	compSel, err := installation.NewComponentSelection(
		installation.ComponentHyprland,
		"0.32.0",
		nil,
	)
	require.NoError(t, err)

	// Create disk space (available >= required)
	diskSpace, err := installation.NewDiskSpace(
		500000000, // available
		100000000, // required
	)
	require.NoError(t, err)

	// Create configuration
	config, err := installation.NewInstallationConfiguration(
		[]installation.ComponentSelection{compSel},
		nil,
		diskSpace,
		false,
	)
	require.NoError(t, err)
	return config
}

func TestReconstructInstallationSession(t *testing.T) {
	t.Run("reconstructs completed session successfully", func(t *testing.T) {
		// Arrange
		id := "test-session-123"
		config := createValidConfig(t)
		status := installation.StatusCompleted

		// Create disk space for snapshot
		diskSpace, err := installation.NewDiskSpace(500000000, 100000000)
		require.NoError(t, err)

		// Create system snapshot
		snapshot, err := installation.NewSystemSnapshot(
			"/tmp/snapshot",
			diskSpace,
			[]string{"hyprland", "waybar"},
		)
		require.NoError(t, err)

		// Create installed component
		installedComp, err := installation.NewInstalledComponent(
			installation.ComponentHyprland,
			"0.32.0",
			nil,
		)
		require.NoError(t, err)

		components := []*installation.InstalledComponent{installedComp}
		startedAt := time.Now().Add(-1 * time.Hour)
		completedAt := time.Now()
		failureReason := ""

		// Act
		session, err := installation.ReconstructInstallationSession(
			id,
			config,
			status,
			snapshot,
			components,
			startedAt,
			completedAt,
			failureReason,
		)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, id, session.ID())
		assert.Equal(t, status, session.Status())
		assert.Equal(t, snapshot, session.Snapshot())
		assert.Len(t, session.InstalledComponents(), 1)
		assert.Equal(t, startedAt, session.StartedAt())
		assert.Equal(t, completedAt, session.CompletedAt())
		assert.Equal(t, failureReason, session.FailureReason())
	})

	t.Run("reconstructs failed session successfully", func(t *testing.T) {
		// Arrange
		id := "failed-session-456"
		config := createValidConfig(t)
		status := installation.StatusFailed
		var snapshot *installation.SystemSnapshot = nil
		components := []*installation.InstalledComponent{}
		startedAt := time.Now().Add(-30 * time.Minute)
		completedAt := time.Now()
		failureReason := "package installation failed: network timeout"

		// Act
		session, err := installation.ReconstructInstallationSession(
			id,
			config,
			status,
			snapshot,
			components,
			startedAt,
			completedAt,
			failureReason,
		)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, id, session.ID())
		assert.Equal(t, status, session.Status())
		assert.Equal(t, failureReason, session.FailureReason())
		assert.Nil(t, session.Snapshot())
		assert.Empty(t, session.InstalledComponents())
	})

	t.Run("reconstructs pending session successfully", func(t *testing.T) {
		// Arrange
		id := "pending-session-789"
		config := createValidConfig(t)
		status := installation.StatusPending
		var snapshot *installation.SystemSnapshot = nil
		components := []*installation.InstalledComponent{}
		startedAt := time.Now()
		completedAt := time.Time{} // Zero time for pending
		failureReason := ""

		// Act
		session, err := installation.ReconstructInstallationSession(
			id,
			config,
			status,
			snapshot,
			components,
			startedAt,
			completedAt,
			failureReason,
		)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, id, session.ID())
		assert.Equal(t, status, session.Status())
		assert.True(t, session.CompletedAt().IsZero())
	})

	t.Run("rejects empty session ID", func(t *testing.T) {
		// Arrange
		id := ""
		config := createValidConfig(t)
		status := installation.StatusPending
		startedAt := time.Now()

		// Act
		session, err := installation.ReconstructInstallationSession(
			id,
			config,
			status,
			nil,
			nil,
			startedAt,
			time.Time{},
			"",
		)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "session ID cannot be empty")
	})

	t.Run("rejects zero started time", func(t *testing.T) {
		// Arrange
		id := "test-session"
		config := createValidConfig(t)
		status := installation.StatusPending
		startedAt := time.Time{} // Zero time - invalid

		// Act
		session, err := installation.ReconstructInstallationSession(
			id,
			config,
			status,
			nil,
			nil,
			startedAt,
			time.Time{},
			"",
		)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "started time cannot be zero")
	})

	t.Run("rejects completed session without completed time", func(t *testing.T) {
		// Arrange
		id := "test-session"
		config := createValidConfig(t)
		status := installation.StatusCompleted
		startedAt := time.Now()
		completedAt := time.Time{} // Zero time but status is completed

		// Act
		session, err := installation.ReconstructInstallationSession(
			id,
			config,
			status,
			nil,
			nil,
			startedAt,
			completedAt,
			"",
		)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "completed session must have completed time")
	})

	t.Run("rejects failed session without failure reason", func(t *testing.T) {
		// Arrange
		id := "test-session"
		config := createValidConfig(t)
		status := installation.StatusFailed
		startedAt := time.Now()
		completedAt := time.Now()
		failureReason := "" // Empty but status is failed

		// Act
		session, err := installation.ReconstructInstallationSession(
			id,
			config,
			status,
			nil,
			nil,
			startedAt,
			completedAt,
			failureReason,
		)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "failed session must have failure reason")
	})
}
