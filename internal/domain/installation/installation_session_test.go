package installation_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstallationSession(t *testing.T) {
	config := mustCreateConfiguration(t, []installation.ComponentSelection{
		mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
	})

	session, err := installation.NewInstallationSession(config)

	require.NoError(t, err)
	assert.NotEmpty(t, session.ID(), "Session should have a unique ID")
	assert.Equal(t, installation.StatusPending, session.Status())
	assert.False(t, session.StartedAt().IsZero())
	assert.True(t, session.CompletedAt().IsZero(), "CompletedAt should not be set initially")
	assert.NotNil(t, session.Configuration())
}

func TestInstallationSession_Identity(t *testing.T) {
	config := mustCreateConfiguration(t, []installation.ComponentSelection{
		mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
	})

	session1, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	session2, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	assert.NotEqual(t, session1.ID(), session2.ID(),
		"Different sessions should have different IDs")
}

func TestInstallationSession_StartPreparation(t *testing.T) {
	config := mustCreateConfiguration(t, []installation.ComponentSelection{
		mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
	})

	session, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	// Create and attach snapshot
	snapshot, err := installation.NewSystemSnapshot(
		"/var/backup/test",
		mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
		[]string{"pkg1"},
	)
	require.NoError(t, err)

	// Start preparation
	err = session.StartPreparation(snapshot)
	require.NoError(t, err)

	assert.Equal(t, installation.StatusPreparation, session.Status())
	assert.NotNil(t, session.Snapshot())
	assert.Equal(t, snapshot.ID(), session.Snapshot().ID())
}

func TestInstallationSession_StartInstalling(t *testing.T) {
	config := mustCreateConfiguration(t, []installation.ComponentSelection{
		mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
	})

	session, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	snapshot, _ := installation.NewSystemSnapshot(
		"/var/backup/test",
		mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
		nil,
	)
	session.StartPreparation(snapshot)

	// Start installing
	err = session.StartInstalling()
	require.NoError(t, err)

	assert.Equal(t, installation.StatusInstalling, session.Status())
}

func TestInstallationSession_AddInstalledComponent(t *testing.T) {
	config := mustCreateConfiguration(t, []installation.ComponentSelection{
		mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
	})

	session, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	snapshot, _ := installation.NewSystemSnapshot("/var/backup/test",
		mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB), nil)
	session.StartPreparation(snapshot)
	session.StartInstalling()

	// Add installed component
	component, _ := installation.NewInstalledComponent(
		installation.ComponentHyprland,
		"0.35.0",
		nil,
	)

	err = session.AddInstalledComponent(component)
	require.NoError(t, err)

	components := session.InstalledComponents()
	assert.Len(t, components, 1)
	assert.Equal(t, component.ID(), components[0].ID())
}

func TestInstallationSession_Complete(t *testing.T) {
	config := mustCreateConfiguration(t, []installation.ComponentSelection{
		mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
	})

	session, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	snapshot, _ := installation.NewSystemSnapshot("/var/backup/test",
		mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB), nil)
	session.StartPreparation(snapshot)
	session.StartInstalling()

	component, _ := installation.NewInstalledComponent(
		installation.ComponentHyprland,
		"0.35.0",
		nil,
	)
	session.AddInstalledComponent(component)

	// Go through proper state transitions
	session.StartConfiguring()
	session.StartVerifying()

	// Complete installation
	err = session.Complete()
	require.NoError(t, err)

	assert.Equal(t, installation.StatusCompleted, session.Status())
	assert.False(t, session.CompletedAt().IsZero(), "CompletedAt should be set")
}

func TestInstallationSession_Fail(t *testing.T) {
	config := mustCreateConfiguration(t, []installation.ComponentSelection{
		mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
	})

	session, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	reason := "network error during download"
	err = session.Fail(reason)
	require.NoError(t, err)

	assert.Equal(t, installation.StatusFailed, session.Status())
	assert.Equal(t, reason, session.FailureReason())
	assert.False(t, session.CompletedAt().IsZero(), "CompletedAt should be set on failure")
}

func TestInstallationSession_StateTransitions(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*installation.InstallationSession)
		action        func(*installation.InstallationSession) error
		expectSuccess bool
	}{
		{
			name:  "can start preparation from pending",
			setup: func(s *installation.InstallationSession) {},
			action: func(s *installation.InstallationSession) error {
				snapshot, _ := installation.NewSystemSnapshot("/path",
					mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB), nil)
				return s.StartPreparation(snapshot)
			},
			expectSuccess: true,
		},
		{
			name: "can start installing from preparation",
			setup: func(s *installation.InstallationSession) {
				snapshot, _ := installation.NewSystemSnapshot("/path",
					mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB), nil)
				s.StartPreparation(snapshot)
			},
			action: func(s *installation.InstallationSession) error {
				return s.StartInstalling()
			},
			expectSuccess: true,
		},
		{
			name: "cannot complete without installed components",
			setup: func(s *installation.InstallationSession) {
				snapshot, _ := installation.NewSystemSnapshot("/path",
					mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB), nil)
				s.StartPreparation(snapshot)
				s.StartInstalling()
			},
			action: func(s *installation.InstallationSession) error {
				return s.Complete()
			},
			expectSuccess: false,
		},
		{
			name: "can fail from any state",
			setup: func(s *installation.InstallationSession) {
				snapshot, _ := installation.NewSystemSnapshot("/path",
					mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB), nil)
				s.StartPreparation(snapshot)
			},
			action: func(s *installation.InstallationSession) error {
				return s.Fail("test failure")
			},
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := mustCreateConfiguration(t, []installation.ComponentSelection{
				mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
			})

			session, err := installation.NewInstallationSession(config)
			require.NoError(t, err)

			tt.setup(session)
			err = tt.action(session)

			if tt.expectSuccess {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestInstallationSession_IsInProgress(t *testing.T) {
	config := mustCreateConfiguration(t, []installation.ComponentSelection{
		mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
	})

	session, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	// Pending is not in progress
	assert.False(t, session.IsInProgress())

	// Preparation is in progress
	snapshot, _ := installation.NewSystemSnapshot("/path",
		mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB), nil)
	session.StartPreparation(snapshot)
	assert.True(t, session.IsInProgress())

	// Installing is in progress
	session.StartInstalling()
	assert.True(t, session.IsInProgress())

	// Go through proper state transitions to complete
	component, _ := installation.NewInstalledComponent(
		installation.ComponentHyprland, "0.35.0", nil)
	session.AddInstalledComponent(component)
	session.StartConfiguring()
	session.StartVerifying()
	session.Complete()

	// Completed is not in progress
	assert.False(t, session.IsInProgress())
}

func TestInstallationSession_Duration(t *testing.T) {
	config := mustCreateConfiguration(t, []installation.ComponentSelection{
		mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
	})

	session, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	// Wait a tiny bit
	time.Sleep(10 * time.Millisecond)

	duration := session.Duration()
	assert.Greater(t, duration, time.Duration(0), "Duration should be positive")
}

func TestInstallationSession_String(t *testing.T) {
	config := mustCreateConfiguration(t, []installation.ComponentSelection{
		mustCreateComponentSelection(t, installation.ComponentHyprland, "0.35.0"),
	})

	session, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	str := session.String()
	assert.Contains(t, str, "pending")
	assert.Contains(t, str, "1 components")
}

// Helper function to create valid configuration
func mustCreateConfiguration(t *testing.T, components []installation.ComponentSelection) installation.InstallationConfiguration {
	t.Helper()
	config, err := installation.NewInstallationConfiguration(
		components,
		nil,
		mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
		false,
	)
	require.NoError(t, err)
	return config
}
