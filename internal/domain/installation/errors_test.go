package installation

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors_AreSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		// Value Object validation errors
		{"ErrInvalidDiskSpace", ErrInvalidDiskSpace},
		{"ErrInvalidPackageInfo", ErrInvalidPackageInfo},
		{"ErrInvalidComponentSelection", ErrInvalidComponentSelection},
		{"ErrInvalidGPUSupport", ErrInvalidGPUSupport},
		{"ErrInvalidConfiguration", ErrInvalidConfiguration},
		{"ErrInvalidProgress", ErrInvalidProgress},

		// Installation Session errors
		{"ErrInsufficientDiskSpace", ErrInsufficientDiskSpace},
		{"ErrPackageConflict", ErrPackageConflict},
		{"ErrNetworkInterruption", ErrNetworkInterruption},
		{"ErrInstallationFailed", ErrInstallationFailed},
		{"ErrRollbackFailed", ErrRollbackFailed},
		{"ErrInvalidStateTransition", ErrInvalidStateTransition},
		{"ErrSessionNotStarted", ErrSessionNotStarted},
		{"ErrSessionAlreadyComplete", ErrSessionAlreadyComplete},

		// Component errors
		{"ErrComponentNotFound", ErrComponentNotFound},
		{"ErrComponentAlreadyExists", ErrComponentAlreadyExists},
		{"ErrDependencyMissing", ErrDependencyMissing},
		{"ErrCoreComponentRequired", ErrCoreComponentRequired},

		// Snapshot errors
		{"ErrSnapshotCreationFailed", ErrSnapshotCreationFailed},
		{"ErrSnapshotRestorationFailed", ErrSnapshotRestorationFailed},
		{"ErrSnapshotInvalid", ErrSnapshotInvalid},
		{"ErrSnapshotNotFound", ErrSnapshotNotFound},

		// Configuration errors
		{"ErrConfigurationMergeFailed", ErrConfigurationMergeFailed},
		{"ErrConfigurationBackupFailed", ErrConfigurationBackupFailed},
		{"ErrConfigurationInvalid", ErrConfigurationInvalid},

		// Repository errors
		{"ErrSessionNotFound", ErrSessionNotFound},
		{"ErrSnapshotSaveFailed", ErrSnapshotSaveFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify error is not nil
			assert.NotNil(t, tt.err)

			// Verify error has a message
			assert.NotEmpty(t, tt.err.Error())

			// Verify error can be compared with errors.Is
			assert.True(t, errors.Is(tt.err, tt.err))
		})
	}
}

func TestErrors_CanBeWrapped(t *testing.T) {
	tests := []struct {
		name        string
		baseErr     error
		wrappedErr  error
		shouldMatch bool
	}{
		{
			name:        "wrapped ErrInsufficientDiskSpace",
			baseErr:     ErrInsufficientDiskSpace,
			wrappedErr:  fmt.Errorf("checking disk space: %w", ErrInsufficientDiskSpace),
			shouldMatch: true,
		},
		{
			name:        "wrapped ErrPackageConflict",
			baseErr:     ErrPackageConflict,
			wrappedErr:  fmt.Errorf("installing hyprland: %w", ErrPackageConflict),
			shouldMatch: true,
		},
		{
			name:        "wrapped ErrNetworkInterruption",
			baseErr:     ErrNetworkInterruption,
			wrappedErr:  fmt.Errorf("downloading packages: %w", ErrNetworkInterruption),
			shouldMatch: true,
		},
		{
			name:        "different errors don't match",
			baseErr:     ErrInsufficientDiskSpace,
			wrappedErr:  fmt.Errorf("some error: %w", ErrPackageConflict),
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := errors.Is(tt.wrappedErr, tt.baseErr)
			assert.Equal(t, tt.shouldMatch, matched)
		})
	}
}

func TestErrors_HaveDescriptiveMessages(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedSubstr string
	}{
		{
			name:           "ErrInsufficientDiskSpace",
			err:            ErrInsufficientDiskSpace,
			expectedSubstr: "insufficient disk space",
		},
		{
			name:           "ErrPackageConflict",
			err:            ErrPackageConflict,
			expectedSubstr: "package conflict",
		},
		{
			name:           "ErrSnapshotCreationFailed",
			err:            ErrSnapshotCreationFailed,
			expectedSubstr: "snapshot creation",
		},
		{
			name:           "ErrInvalidStateTransition",
			err:            ErrInvalidStateTransition,
			expectedSubstr: "invalid state transition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, tt.err.Error(), tt.expectedSubstr)
		})
	}
}

func TestErrors_AreDistinct(t *testing.T) {
	// Ensure all errors are distinct sentinel values
	allErrors := []error{
		ErrInvalidDiskSpace,
		ErrInvalidPackageInfo,
		ErrInvalidComponentSelection,
		ErrInvalidGPUSupport,
		ErrInvalidConfiguration,
		ErrInvalidProgress,
		ErrInsufficientDiskSpace,
		ErrPackageConflict,
		ErrNetworkInterruption,
		ErrInstallationFailed,
		ErrRollbackFailed,
		ErrInvalidStateTransition,
		ErrSessionNotStarted,
		ErrSessionAlreadyComplete,
		ErrComponentNotFound,
		ErrComponentAlreadyExists,
		ErrDependencyMissing,
		ErrCoreComponentRequired,
		ErrSnapshotCreationFailed,
		ErrSnapshotRestorationFailed,
		ErrSnapshotInvalid,
		ErrSnapshotNotFound,
		ErrConfigurationMergeFailed,
		ErrConfigurationBackupFailed,
		ErrConfigurationInvalid,
		ErrSessionNotFound,
		ErrSnapshotSaveFailed,
	}

	// Check that no two errors are the same
	for i, err1 := range allErrors {
		for j, err2 := range allErrors {
			if i != j {
				assert.NotEqual(t, err1, err2, "Errors at index %d and %d should be distinct", i, j)
			}
		}
	}
}
