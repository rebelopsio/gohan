package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/application/history/services"
	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/rebelopsio/gohan/internal/infrastructure/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistoryRecordingService_RecordSuccessfulInstallation(t *testing.T) {
	repo := memory.NewHistoryRepository()
	service := services.NewHistoryRecordingService(repo)
	ctx := context.Background()

	// Create a completed installation session
	session := createCompletedSession(t)

	// Record the installation
	recordID, err := service.RecordInstallation(ctx, session)
	require.NoError(t, err)
	assert.True(t, recordID.IsValid())

	// Verify record was saved
	record, err := repo.FindByID(ctx, recordID)
	require.NoError(t, err)

	// Verify record details
	assert.Equal(t, session.ID(), record.SessionID())
	assert.True(t, record.WasSuccessful())
	assert.Equal(t, "hyprland", record.PackageName())
	assert.False(t, record.HasFailureDetails())
	assert.Greater(t, record.PackageCount(), 0)
}

func TestHistoryRecordingService_RecordFailedInstallation(t *testing.T) {
	repo := memory.NewHistoryRepository()
	service := services.NewHistoryRecordingService(repo)
	ctx := context.Background()

	// Create a failed installation session
	session := createFailedSession(t)

	// Record the installation
	recordID, err := service.RecordInstallation(ctx, session)
	require.NoError(t, err)
	assert.True(t, recordID.IsValid())

	// Verify record was saved
	record, err := repo.FindByID(ctx, recordID)
	require.NoError(t, err)

	// Verify record details
	assert.Equal(t, session.ID(), record.SessionID())
	assert.True(t, record.WasFailed())
	assert.True(t, record.HasFailureDetails())

	// Check failure details
	failureDetails := record.FailureDetails()
	require.NotNil(t, failureDetails)
	assert.Contains(t, failureDetails.Reason(), "Package conflict")
}

func TestHistoryRecordingService_RecordIncompleteSession(t *testing.T) {
	repo := memory.NewHistoryRepository()
	service := services.NewHistoryRecordingService(repo)
	ctx := context.Background()

	// Create an incomplete session (not finished)
	session := createPendingSession(t)

	// Attempt to record - should fail
	_, err := service.RecordInstallation(ctx, session)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot record incomplete session")
}

func TestHistoryRecordingService_CapturesSystemContext(t *testing.T) {
	repo := memory.NewHistoryRepository()
	service := services.NewHistoryRecordingService(repo)
	ctx := context.Background()

	session := createCompletedSession(t)

	recordID, err := service.RecordInstallation(ctx, session)
	require.NoError(t, err)

	record, err := repo.FindByID(ctx, recordID)
	require.NoError(t, err)

	// Verify system context was captured
	sysCtx := record.SystemContext()
	assert.NotEmpty(t, sysCtx.OSVersion())
	// Hostname may be empty in tests, so just check it's accessible
	_ = sysCtx.Hostname()
}

func TestHistoryRecordingService_CapturesInstalledPackages(t *testing.T) {
	repo := memory.NewHistoryRepository()
	service := services.NewHistoryRecordingService(repo)
	ctx := context.Background()

	session := createCompletedSessionWithMultipleComponents(t)

	recordID, err := service.RecordInstallation(ctx, session)
	require.NoError(t, err)

	record, err := repo.FindByID(ctx, recordID)
	require.NoError(t, err)

	// Verify multiple packages were captured
	assert.GreaterOrEqual(t, record.PackageCount(), 2)

	metadata := record.Metadata()
	packages := metadata.InstalledPackages()
	assert.GreaterOrEqual(t, len(packages), 2)

	// Verify package details
	for _, pkg := range packages {
		assert.NotEmpty(t, pkg.Name())
		assert.NotEmpty(t, pkg.Version())
		assert.Greater(t, pkg.SizeBytes(), uint64(0))
	}
}

func TestHistoryRecordingService_CaptureDuration(t *testing.T) {
	repo := memory.NewHistoryRepository()
	service := services.NewHistoryRecordingService(repo)
	ctx := context.Background()

	session := createCompletedSession(t)

	recordID, err := service.RecordInstallation(ctx, session)
	require.NoError(t, err)

	record, err := repo.FindByID(ctx, recordID)
	require.NoError(t, err)

	// Verify duration is captured (can be 0 in fast tests)
	duration := record.Duration()
	assert.GreaterOrEqual(t, duration, time.Duration(0))

	metadata := record.Metadata()
	assert.True(t, metadata.CompletedAt().After(metadata.InstalledAt()) || metadata.CompletedAt().Equal(metadata.InstalledAt()))
}

// Helper functions to create test sessions

func createCompletedSession(t *testing.T) *installation.InstallationSession {
	// Create config
	components := []installation.ComponentSelection{
		mustCreateComponent(t, installation.ComponentHyprland, "0.45.0", 15728640),
	}
	diskSpace, err := installation.NewDiskSpace(21474836480, 0) // 20 GB available
	require.NoError(t, err)

	config, err := installation.NewInstallationConfiguration(
		components,
		nil, // No GPU support specified
		diskSpace,
		false,
	)
	require.NoError(t, err)

	// Create session
	session, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	// Create snapshot
	snapshot, err := installation.NewSystemSnapshot(
		"/tmp/snapshots",
		diskSpace,
		[]string{},
	)
	require.NoError(t, err)

	// Progress through installation
	require.NoError(t, session.StartPreparation(snapshot))
	require.NoError(t, session.StartInstalling())

	// Add installed component
	installedComp, err := installation.NewInstalledComponent(
		installation.ComponentHyprland,
		"0.45.0",
		nil,
	)
	require.NoError(t, err)
	require.NoError(t, session.AddInstalledComponent(installedComp))

	require.NoError(t, session.StartConfiguring())
	require.NoError(t, session.StartVerifying())
	require.NoError(t, session.Complete())

	return session
}

func createFailedSession(t *testing.T) *installation.InstallationSession {
	components := []installation.ComponentSelection{
		mustCreateComponent(t, installation.ComponentHyprland, "0.45.0", 15728640),
	}
	diskSpace, err := installation.NewDiskSpace(21474836480, 0)
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

	snapshot, err := installation.NewSystemSnapshot(
		"/tmp/snapshots",
		diskSpace,
		[]string{},
	)
	require.NoError(t, err)

	require.NoError(t, session.StartPreparation(snapshot))
	require.NoError(t, session.StartInstalling())

	// Fail the session
	require.NoError(t, session.Fail("Package conflict detected"))

	return session
}

func createPendingSession(t *testing.T) *installation.InstallationSession {
	components := []installation.ComponentSelection{
		mustCreateComponent(t, installation.ComponentHyprland, "0.45.0", 15728640),
	}
	diskSpace, err := installation.NewDiskSpace(21474836480, 0)
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

func createCompletedSessionWithMultipleComponents(t *testing.T) *installation.InstallationSession {
	components := []installation.ComponentSelection{
		mustCreateComponent(t, installation.ComponentHyprland, "0.45.0", 15728640),
		mustCreateComponent(t, installation.ComponentWaybar, "0.10.4", 5242880),
	}
	diskSpace, err := installation.NewDiskSpace(21474836480, 0)
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

	snapshot, err := installation.NewSystemSnapshot(
		"/tmp/snapshots",
		diskSpace,
		[]string{},
	)
	require.NoError(t, err)

	require.NoError(t, session.StartPreparation(snapshot))
	require.NoError(t, session.StartInstalling())

	// Add multiple installed components
	comp1, err := installation.NewInstalledComponent(
		installation.ComponentHyprland,
		"0.45.0",
		nil,
	)
	require.NoError(t, err)
	require.NoError(t, session.AddInstalledComponent(comp1))

	comp2, err := installation.NewInstalledComponent(
		installation.ComponentWaybar,
		"0.10.4",
		nil,
	)
	require.NoError(t, err)
	require.NoError(t, session.AddInstalledComponent(comp2))

	require.NoError(t, session.StartConfiguring())
	require.NoError(t, session.StartVerifying())
	require.NoError(t, session.Complete())

	return session
}

func mustCreateComponent(t *testing.T, component installation.ComponentName, version string, sizeBytes uint64) installation.ComponentSelection {
	pkgInfo, err := installation.NewPackageInfo(
		string(component),
		version,
		sizeBytes,
		[]string{},
	)
	require.NoError(t, err)

	comp, err := installation.NewComponentSelection(
		component,
		version,
		&pkgInfo,
	)
	require.NoError(t, err)

	return comp
}
