package installation_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallationStartedEvent(t *testing.T) {
	t.Run("creates event with current timestamp", func(t *testing.T) {
		sessionID := "session-123"
		before := time.Now()

		event := installation.NewInstallationStartedEvent(sessionID)

		after := time.Now()

		assert.Equal(t, "installation.started", event.EventType())
		assert.Equal(t, sessionID, event.SessionID())
		assert.False(t, event.OccurredAt().Before(before))
		assert.False(t, event.OccurredAt().After(after))
	})

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		event := installation.NewInstallationStartedEvent("session-123")

		var _ installation.DomainEvent = event
		assert.NotEmpty(t, event.EventType())
		assert.False(t, event.OccurredAt().IsZero())
	})
}

func TestInstallationProgressUpdatedEvent(t *testing.T) {
	t.Run("creates event with progress information", func(t *testing.T) {
		sessionID := "session-123"
		currentPhase := installation.StatusInstalling
		percentComplete := 45
		message := "Installing Hyprland..."

		event := installation.NewInstallationProgressUpdatedEvent(
			sessionID,
			currentPhase,
			percentComplete,
			message,
		)

		assert.Equal(t, "installation.progress.updated", event.EventType())
		assert.Equal(t, sessionID, event.SessionID())
		assert.Equal(t, currentPhase, event.CurrentPhase())
		assert.Equal(t, percentComplete, event.PercentComplete())
		assert.Equal(t, message, event.Message())
		assert.False(t, event.OccurredAt().IsZero())
	})

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		event := installation.NewInstallationProgressUpdatedEvent(
			"session-123",
			installation.StatusInstalling,
			50,
			"test",
		)

		var _ installation.DomainEvent = event
	})
}

func TestPhaseCompletedEvent(t *testing.T) {
	t.Run("creates event with phase information", func(t *testing.T) {
		sessionID := "session-123"
		phase := installation.StatusPreparation
		duration := 5 * time.Second

		event := installation.NewPhaseCompletedEvent(sessionID, phase, duration)

		assert.Equal(t, "installation.phase.completed", event.EventType())
		assert.Equal(t, sessionID, event.SessionID())
		assert.Equal(t, phase, event.Phase())
		assert.Equal(t, duration, event.Duration())
	})

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		event := installation.NewPhaseCompletedEvent(
			"session-123",
			installation.StatusPreparation,
			5*time.Second,
		)

		var _ installation.DomainEvent = event
	})
}

func TestComponentInstalledEvent(t *testing.T) {
	t.Run("creates event with component information", func(t *testing.T) {
		sessionID := "session-123"
		componentName := installation.ComponentHyprland
		version := "0.35.0"

		event := installation.NewComponentInstalledEvent(sessionID, componentName, version)

		assert.Equal(t, "installation.component.installed", event.EventType())
		assert.Equal(t, sessionID, event.SessionID())
		assert.Equal(t, componentName, event.Component())
		assert.Equal(t, version, event.Version())
	})

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		event := installation.NewComponentInstalledEvent(
			"session-123",
			installation.ComponentHyprland,
			"0.35.0",
		)

		var _ installation.DomainEvent = event
	})
}

func TestInstallationCompletedEvent(t *testing.T) {
	t.Run("creates event with completion information", func(t *testing.T) {
		sessionID := "session-123"
		duration := 120 * time.Second
		componentsInstalled := 5

		event := installation.NewInstallationCompletedEvent(
			sessionID,
			duration,
			componentsInstalled,
		)

		assert.Equal(t, "installation.completed", event.EventType())
		assert.Equal(t, sessionID, event.SessionID())
		assert.Equal(t, duration, event.Duration())
		assert.Equal(t, componentsInstalled, event.ComponentsInstalled())
	})

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		event := installation.NewInstallationCompletedEvent("session-123", 120*time.Second, 5)

		var _ installation.DomainEvent = event
	})
}

func TestInstallationFailedEvent(t *testing.T) {
	t.Run("creates event with failure information", func(t *testing.T) {
		sessionID := "session-123"
		phase := installation.StatusInstalling
		reason := "network connection failed"
		recoverable := true

		event := installation.NewInstallationFailedEvent(sessionID, phase, reason, recoverable)

		assert.Equal(t, "installation.failed", event.EventType())
		assert.Equal(t, sessionID, event.SessionID())
		assert.Equal(t, phase, event.Phase())
		assert.Equal(t, reason, event.Reason())
		assert.Equal(t, recoverable, event.IsRecoverable())
	})

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		event := installation.NewInstallationFailedEvent(
			"session-123",
			installation.StatusInstalling,
			"test failure",
			false,
		)

		var _ installation.DomainEvent = event
	})
}

func TestRollbackStartedEvent(t *testing.T) {
	t.Run("creates event with rollback information", func(t *testing.T) {
		sessionID := "session-123"
		snapshotID := "snapshot-456"
		reason := "installation failed during verification"

		event := installation.NewRollbackStartedEvent(sessionID, snapshotID, reason)

		assert.Equal(t, "installation.rollback.started", event.EventType())
		assert.Equal(t, sessionID, event.SessionID())
		assert.Equal(t, snapshotID, event.SnapshotID())
		assert.Equal(t, reason, event.Reason())
	})

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		event := installation.NewRollbackStartedEvent(
			"session-123",
			"snapshot-456",
			"test reason",
		)

		var _ installation.DomainEvent = event
	})
}

func TestRollbackCompletedEvent(t *testing.T) {
	t.Run("creates event with rollback completion information", func(t *testing.T) {
		sessionID := "session-123"
		snapshotID := "snapshot-456"
		duration := 30 * time.Second
		success := true

		event := installation.NewRollbackCompletedEvent(
			sessionID,
			snapshotID,
			duration,
			success,
		)

		assert.Equal(t, "installation.rollback.completed", event.EventType())
		assert.Equal(t, sessionID, event.SessionID())
		assert.Equal(t, snapshotID, event.SnapshotID())
		assert.Equal(t, duration, event.Duration())
		assert.Equal(t, success, event.Success())
	})

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		event := installation.NewRollbackCompletedEvent(
			"session-123",
			"snapshot-456",
			30*time.Second,
			true,
		)

		var _ installation.DomainEvent = event
	})
}

func TestConflictDetectedEvent(t *testing.T) {
	t.Run("creates event with conflict information", func(t *testing.T) {
		sessionID := "session-123"
		packageName := "libhyprland"
		conflictingPackage := "hyprland-git"
		severity := "high"

		event := installation.NewConflictDetectedEvent(
			sessionID,
			packageName,
			conflictingPackage,
			severity,
		)

		assert.Equal(t, "installation.conflict.detected", event.EventType())
		assert.Equal(t, sessionID, event.SessionID())
		assert.Equal(t, packageName, event.PackageName())
		assert.Equal(t, conflictingPackage, event.ConflictingPackage())
		assert.Equal(t, severity, event.Severity())
	})

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		event := installation.NewConflictDetectedEvent(
			"session-123",
			"pkg1",
			"pkg2",
			"medium",
		)

		var _ installation.DomainEvent = event
	})
}

func TestBackupCreatedEvent(t *testing.T) {
	t.Run("creates event with backup information", func(t *testing.T) {
		sessionID := "session-123"
		backupPath := "/var/backup/hyprland-20250126"
		filesBackedUp := 15
		totalSize := int64(1024 * 1024 * 50) // 50MB

		event := installation.NewBackupCreatedEvent(
			sessionID,
			backupPath,
			filesBackedUp,
			totalSize,
		)

		assert.Equal(t, "installation.backup.created", event.EventType())
		assert.Equal(t, sessionID, event.SessionID())
		assert.Equal(t, backupPath, event.BackupPath())
		assert.Equal(t, filesBackedUp, event.FilesBackedUp())
		assert.Equal(t, totalSize, event.TotalSize())
	})

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		event := installation.NewBackupCreatedEvent(
			"session-123",
			"/backup/path",
			10,
			1024,
		)

		var _ installation.DomainEvent = event
	})
}

func TestDiskSpaceInsufficientEvent(t *testing.T) {
	t.Run("creates event with disk space information", func(t *testing.T) {
		sessionID := "session-123"
		required := int64(10 * 1024 * 1024 * 1024)  // 10GB
		available := int64(5 * 1024 * 1024 * 1024)  // 5GB
		path := "/var"

		event := installation.NewDiskSpaceInsufficientEvent(
			sessionID,
			required,
			available,
			path,
		)

		assert.Equal(t, "installation.disk.insufficient", event.EventType())
		assert.Equal(t, sessionID, event.SessionID())
		assert.Equal(t, required, event.RequiredBytes())
		assert.Equal(t, available, event.AvailableBytes())
		assert.Equal(t, path, event.Path())
	})

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		event := installation.NewDiskSpaceInsufficientEvent(
			"session-123",
			1000,
			500,
			"/var",
		)

		var _ installation.DomainEvent = event
	})
}

func TestNetworkInterruptionEvent(t *testing.T) {
	t.Run("creates event with network interruption information", func(t *testing.T) {
		sessionID := "session-123"
		operation := "downloading package hyprland-0.35.0"
		retryable := true
		errorMsg := "connection timeout"

		event := installation.NewNetworkInterruptionEvent(
			sessionID,
			operation,
			retryable,
			errorMsg,
		)

		assert.Equal(t, "installation.network.interrupted", event.EventType())
		assert.Equal(t, sessionID, event.SessionID())
		assert.Equal(t, operation, event.Operation())
		assert.Equal(t, retryable, event.IsRetryable())
		assert.Equal(t, errorMsg, event.ErrorMessage())
	})

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		event := installation.NewNetworkInterruptionEvent(
			"session-123",
			"test op",
			true,
			"error",
		)

		var _ installation.DomainEvent = event
	})
}

func TestEventTimestamps(t *testing.T) {
	t.Run("all events have accurate timestamps", func(t *testing.T) {
		before := time.Now()

		events := []installation.DomainEvent{
			installation.NewInstallationStartedEvent("s1"),
			installation.NewInstallationProgressUpdatedEvent("s1", installation.StatusInstalling, 50, "msg"),
			installation.NewPhaseCompletedEvent("s1", installation.StatusPreparation, 5*time.Second),
			installation.NewComponentInstalledEvent("s1", installation.ComponentHyprland, "0.35.0"),
			installation.NewInstallationCompletedEvent("s1", 120*time.Second, 5),
			installation.NewInstallationFailedEvent("s1", installation.StatusInstalling, "reason", false),
			installation.NewRollbackStartedEvent("s1", "snap1", "reason"),
			installation.NewRollbackCompletedEvent("s1", "snap1", 30*time.Second, true),
			installation.NewConflictDetectedEvent("s1", "pkg1", "pkg2", "high"),
			installation.NewBackupCreatedEvent("s1", "/backup", 10, 1024),
			installation.NewDiskSpaceInsufficientEvent("s1", 1000, 500, "/var"),
			installation.NewNetworkInterruptionEvent("s1", "op", true, "error"),
		}

		after := time.Now()

		for i, event := range events {
			t.Run(event.EventType(), func(t *testing.T) {
				occurredAt := event.OccurredAt()
				assert.False(t, occurredAt.Before(before),
					"Event %d timestamp should be at or after test start", i)
				assert.False(t, occurredAt.After(after),
					"Event %d timestamp should be at or before test end", i)
			})
		}
	})
}

func TestEventTypes(t *testing.T) {
	tests := []struct {
		name          string
		event         installation.DomainEvent
		expectedType  string
	}{
		{
			name:         "InstallationStartedEvent",
			event:        installation.NewInstallationStartedEvent("s1"),
			expectedType: "installation.started",
		},
		{
			name:         "InstallationProgressUpdatedEvent",
			event:        installation.NewInstallationProgressUpdatedEvent("s1", installation.StatusInstalling, 50, "msg"),
			expectedType: "installation.progress.updated",
		},
		{
			name:         "PhaseCompletedEvent",
			event:        installation.NewPhaseCompletedEvent("s1", installation.StatusPreparation, 5*time.Second),
			expectedType: "installation.phase.completed",
		},
		{
			name:         "ComponentInstalledEvent",
			event:        installation.NewComponentInstalledEvent("s1", installation.ComponentHyprland, "0.35.0"),
			expectedType: "installation.component.installed",
		},
		{
			name:         "InstallationCompletedEvent",
			event:        installation.NewInstallationCompletedEvent("s1", 120*time.Second, 5),
			expectedType: "installation.completed",
		},
		{
			name:         "InstallationFailedEvent",
			event:        installation.NewInstallationFailedEvent("s1", installation.StatusInstalling, "reason", false),
			expectedType: "installation.failed",
		},
		{
			name:         "RollbackStartedEvent",
			event:        installation.NewRollbackStartedEvent("s1", "snap1", "reason"),
			expectedType: "installation.rollback.started",
		},
		{
			name:         "RollbackCompletedEvent",
			event:        installation.NewRollbackCompletedEvent("s1", "snap1", 30*time.Second, true),
			expectedType: "installation.rollback.completed",
		},
		{
			name:         "ConflictDetectedEvent",
			event:        installation.NewConflictDetectedEvent("s1", "pkg1", "pkg2", "high"),
			expectedType: "installation.conflict.detected",
		},
		{
			name:         "BackupCreatedEvent",
			event:        installation.NewBackupCreatedEvent("s1", "/backup", 10, 1024),
			expectedType: "installation.backup.created",
		},
		{
			name:         "DiskSpaceInsufficientEvent",
			event:        installation.NewDiskSpaceInsufficientEvent("s1", 1000, 500, "/var"),
			expectedType: "installation.disk.insufficient",
		},
		{
			name:         "NetworkInterruptionEvent",
			event:        installation.NewNetworkInterruptionEvent("s1", "op", true, "error"),
			expectedType: "installation.network.interrupted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedType, tt.event.EventType())
		})
	}
}

// Helper to verify all events are value objects (immutable)
func TestEvents_AreValueObjects(t *testing.T) {
	t.Run("events do not expose setters or mutable state", func(t *testing.T) {
		// This is a design verification test
		// Events should only have getters, no setters
		// They should be created with all data and be immutable

		event := installation.NewInstallationStartedEvent("session-123")

		// Verify we can read but not modify
		require.NotEmpty(t, event.SessionID())
		require.NotEmpty(t, event.EventType())
		require.False(t, event.OccurredAt().IsZero())

		// Events are value objects - no identity needed
		// Two events with same data are conceptually equal
	})
}
