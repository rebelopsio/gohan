package history_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstallationRecord(t *testing.T) {
	sessionID := "session-123"
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := time.Date(2025, 10, 26, 14, 32, 30, 0, time.UTC)
	recordedAt := time.Date(2025, 10, 26, 14, 32, 35, 0, time.UTC)

	pkg1, _ := history.NewInstalledPackage("hyprland", "0.45.0", 15728640)
	packages := []history.InstalledPackage{pkg1}

	successMetadata, _ := history.NewInstallationMetadata(
		"hyprland",
		"0.45.0",
		installedAt,
		completedAt,
		packages,
	)

	failedMetadata, _ := history.NewInstallationMetadata(
		"kitty",
		"0.36.0",
		installedAt,
		completedAt,
		nil,
	)

	systemCtx, _ := history.NewSystemContext(
		"Debian GNU/Linux 13 (trixie)",
		"6.1.0-13-amd64",
		"1.0.0",
		"myserver",
	)

	failureDetails, _ := history.NewFailureDetails(
		"Package dependency conflict",
		completedAt,
		"package_installation",
		"DEP_CONFLICT",
	)

	successOutcome, _ := history.NewInstallationOutcome("success")
	failedOutcome, _ := history.NewInstallationOutcome("failed")

	tests := []struct {
		name           string
		sessionID      string
		outcome        history.InstallationOutcome
		metadata       history.InstallationMetadata
		systemContext  history.SystemContext
		failureDetails *history.FailureDetails
		recordedAt     time.Time
		wantErr        error
	}{
		{
			name:           "valid successful record",
			sessionID:      sessionID,
			outcome:        successOutcome,
			metadata:       successMetadata,
			systemContext:  systemCtx,
			failureDetails: nil,
			recordedAt:     recordedAt,
			wantErr:        nil,
		},
		{
			name:           "valid failed record with failure details",
			sessionID:      sessionID,
			outcome:        failedOutcome,
			metadata:       failedMetadata,
			systemContext:  systemCtx,
			failureDetails: &failureDetails,
			recordedAt:     recordedAt,
			wantErr:        nil,
		},
		{
			name:           "empty session ID",
			sessionID:      "",
			outcome:        successOutcome,
			metadata:       successMetadata,
			systemContext:  systemCtx,
			failureDetails: nil,
			recordedAt:     recordedAt,
			wantErr:        history.ErrInvalidSessionID,
		},
		{
			name:           "whitespace-only session ID",
			sessionID:      "   ",
			outcome:        successOutcome,
			metadata:       successMetadata,
			systemContext:  systemCtx,
			failureDetails: nil,
			recordedAt:     recordedAt,
			wantErr:        history.ErrInvalidSessionID,
		},
		{
			name:           "zero recordedAt timestamp",
			sessionID:      sessionID,
			outcome:        successOutcome,
			metadata:       successMetadata,
			systemContext:  systemCtx,
			failureDetails: nil,
			recordedAt:     time.Time{},
			wantErr:        history.ErrInvalidRecordedTime,
		},
		{
			name:           "failed record without failure details",
			sessionID:      sessionID,
			outcome:        failedOutcome,
			metadata:       failedMetadata,
			systemContext:  systemCtx,
			failureDetails: nil,
			recordedAt:     recordedAt,
			wantErr:        history.ErrMissingFailureDetails,
		},
		{
			name:           "successful record with zero packages",
			sessionID:      sessionID,
			outcome:        successOutcome,
			metadata:       failedMetadata, // has zero packages
			systemContext:  systemCtx,
			failureDetails: nil,
			recordedAt:     recordedAt,
			wantErr:        history.ErrNoPackagesInstalled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record, err := history.NewInstallationRecord(
				tt.sessionID,
				tt.outcome,
				tt.metadata,
				tt.systemContext,
				tt.failureDetails,
				tt.recordedAt,
			)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.True(t, record.ID().IsValid())
				assert.NotEmpty(t, record.SessionID())
			}
		})
	}
}

func TestInstallationRecord_Accessors(t *testing.T) {
	sessionID := "session-123"
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := time.Date(2025, 10, 26, 14, 32, 30, 0, time.UTC)
	recordedAt := time.Date(2025, 10, 26, 14, 32, 35, 0, time.UTC)

	pkg1, _ := history.NewInstalledPackage("hyprland", "0.45.0", 15728640)
	metadata, _ := history.NewInstallationMetadata(
		"hyprland",
		"0.45.0",
		installedAt,
		completedAt,
		[]history.InstalledPackage{pkg1},
	)

	systemCtx, _ := history.NewSystemContext(
		"Debian GNU/Linux 13",
		"6.1.0-13-amd64",
		"1.0.0",
		"myserver",
	)

	outcome, _ := history.NewInstallationOutcome("success")

	record, err := history.NewInstallationRecord(
		sessionID,
		outcome,
		metadata,
		systemCtx,
		nil,
		recordedAt,
	)
	require.NoError(t, err)

	assert.True(t, record.ID().IsValid())
	assert.Equal(t, sessionID, record.SessionID())
	assert.Equal(t, outcome, record.Outcome())
	assert.Equal(t, metadata, record.Metadata())
	assert.Equal(t, systemCtx, record.SystemContext())
	assert.Equal(t, recordedAt, record.RecordedAt())
	assert.Nil(t, record.FailureDetails())
}

func TestInstallationRecord_WasSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		outcome  string
		expected bool
	}{
		{"success outcome", "success", true},
		{"failed outcome", "failed", false},
		{"rolled_back outcome", "rolled_back", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := createTestRecord(t, tt.outcome, nil, 1)
			assert.Equal(t, tt.expected, record.WasSuccessful())
		})
	}
}

func TestInstallationRecord_WasFailed(t *testing.T) {
	tests := []struct {
		name     string
		outcome  string
		expected bool
	}{
		{"success outcome", "success", false},
		{"failed outcome", "failed", true},
		{"rolled_back outcome", "rolled_back", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			failureDetails, _ := history.NewFailureDetails(
				"error",
				time.Now(),
				"phase",
				"ERR",
			)
			pkgCount := 0
			if tt.outcome == "success" {
				pkgCount = 1
			}
			record := createTestRecord(t, tt.outcome, &failureDetails, pkgCount)
			assert.Equal(t, tt.expected, record.WasFailed())
		})
	}
}

func TestInstallationRecord_WasRolledBack(t *testing.T) {
	tests := []struct {
		name     string
		outcome  string
		expected bool
	}{
		{"success outcome", "success", false},
		{"failed outcome", "failed", false},
		{"rolled_back outcome", "rolled_back", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			failureDetails, _ := history.NewFailureDetails(
				"error",
				time.Now(),
				"phase",
				"ERR",
			)
			pkgCount := 0
			if tt.outcome == "success" {
				pkgCount = 1
			}
			record := createTestRecord(t, tt.outcome, &failureDetails, pkgCount)
			assert.Equal(t, tt.expected, record.WasRolledBack())
		})
	}
}

func TestInstallationRecord_HasFailureDetails(t *testing.T) {
	t.Run("record with failure details", func(t *testing.T) {
		failureDetails, _ := history.NewFailureDetails(
			"error",
			time.Now(),
			"phase",
			"ERR",
		)
		record := createTestRecord(t, "failed", &failureDetails, 0)
		assert.True(t, record.HasFailureDetails())
	})

	t.Run("record without failure details", func(t *testing.T) {
		record := createTestRecord(t, "success", nil, 1)
		assert.False(t, record.HasFailureDetails())
	})
}

func TestInstallationRecord_PackageName(t *testing.T) {
	record := createTestRecord(t, "success", nil, 1)
	assert.Equal(t, "test-package", record.PackageName())
}

func TestInstallationRecord_TargetVersion(t *testing.T) {
	record := createTestRecord(t, "success", nil, 1)
	assert.Equal(t, "1.0.0", record.TargetVersion())
}

func TestInstallationRecord_InstalledAt(t *testing.T) {
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	record := createTestRecordWithTime(t, "success", nil, 1, installedAt)
	assert.Equal(t, installedAt, record.InstalledAt())
}

func TestInstallationRecord_Duration(t *testing.T) {
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := installedAt.Add(2 * time.Minute)

	pkg1, _ := history.NewInstalledPackage("test-pkg", "1.0.0", 1024)
	metadata, _ := history.NewInstallationMetadata(
		"test-pkg",
		"1.0.0",
		installedAt,
		completedAt,
		[]history.InstalledPackage{pkg1},
	)

	systemCtx, _ := history.NewSystemContext("OS", "", "", "")
	outcome, _ := history.NewInstallationOutcome("success")

	record, err := history.NewInstallationRecord(
		"session-123",
		outcome,
		metadata,
		systemCtx,
		nil,
		time.Now(),
	)
	require.NoError(t, err)

	assert.Equal(t, 2*time.Minute, record.Duration())
}

func TestInstallationRecord_PackageCount(t *testing.T) {
	t.Run("with multiple packages", func(t *testing.T) {
		record := createTestRecord(t, "success", nil, 3)
		assert.Equal(t, 3, record.PackageCount())
	})

	t.Run("with no packages (failed)", func(t *testing.T) {
		failureDetails, _ := history.NewFailureDetails(
			"error",
			time.Now(),
			"phase",
			"ERR",
		)
		record := createTestRecord(t, "failed", &failureDetails, 0)
		assert.Equal(t, 0, record.PackageCount())
	})
}

func TestInstallationRecord_TrimsSessionID(t *testing.T) {
	installedAt := time.Now()
	completedAt := installedAt.Add(time.Minute)

	pkg1, _ := history.NewInstalledPackage("test-pkg", "1.0.0", 1024)
	metadata, _ := history.NewInstallationMetadata(
		"test-pkg",
		"1.0.0",
		installedAt,
		completedAt,
		[]history.InstalledPackage{pkg1},
	)

	systemCtx, _ := history.NewSystemContext("OS", "", "", "")
	outcome, _ := history.NewInstallationOutcome("success")

	record, err := history.NewInstallationRecord(
		"  session-123  ",
		outcome,
		metadata,
		systemCtx,
		nil,
		time.Now(),
	)
	require.NoError(t, err)

	assert.Equal(t, "session-123", record.SessionID())
}

// Helper function to create test record
func createTestRecord(t *testing.T, outcomeStr string, failureDetails *history.FailureDetails, packageCount int) history.InstallationRecord {
	return createTestRecordWithTime(t, outcomeStr, failureDetails, packageCount, time.Now())
}

// Helper function to create test record with specific time
func createTestRecordWithTime(t *testing.T, outcomeStr string, failureDetails *history.FailureDetails, packageCount int, installedAt time.Time) history.InstallationRecord {
	completedAt := installedAt.Add(time.Minute)

	var packages []history.InstalledPackage
	for i := 0; i < packageCount; i++ {
		pkg, _ := history.NewInstalledPackage("test-package", "1.0.0", 1024)
		packages = append(packages, pkg)
	}

	metadata, _ := history.NewInstallationMetadata(
		"test-package",
		"1.0.0",
		installedAt,
		completedAt,
		packages,
	)

	systemCtx, _ := history.NewSystemContext("OS", "", "", "")
	outcome, _ := history.NewInstallationOutcome(outcomeStr)

	// Auto-create failure details for failed outcomes if not provided
	if outcomeStr == "failed" && failureDetails == nil {
		fd, _ := history.NewFailureDetails("test failure", completedAt, "test", "ERR")
		failureDetails = &fd
	}

	record, err := history.NewInstallationRecord(
		"session-123",
		outcome,
		metadata,
		systemCtx,
		failureDetails,
		time.Now(),
	)
	require.NoError(t, err)

	return record
}
