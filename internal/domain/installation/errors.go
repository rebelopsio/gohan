package installation

import "errors"

var (
	// Value Object validation errors
	ErrInvalidDiskSpace          = errors.New("invalid disk space value")
	ErrInvalidPackageInfo        = errors.New("invalid package information")
	ErrInvalidComponentSelection = errors.New("invalid component selection")
	ErrInvalidGPUSupport         = errors.New("invalid GPU support configuration")
	ErrInvalidConfiguration      = errors.New("invalid installation configuration")
	ErrInvalidProgress           = errors.New("invalid installation progress")

	// Installation Session errors
	ErrInsufficientDiskSpace  = errors.New("insufficient disk space for installation")
	ErrPackageConflict        = errors.New("package conflict detected")
	ErrNetworkInterruption    = errors.New("network connection interrupted")
	ErrInstallationFailed     = errors.New("installation failed")
	ErrRollbackFailed         = errors.New("rollback operation failed")
	ErrInvalidStateTransition = errors.New("invalid state transition")
	ErrSessionNotStarted      = errors.New("installation session not started")
	ErrSessionAlreadyComplete = errors.New("installation session already completed")

	// Component errors
	ErrComponentNotFound      = errors.New("component not found")
	ErrComponentAlreadyExists = errors.New("component already installed")
	ErrDependencyMissing      = errors.New("required dependency missing")
	ErrCoreComponentRequired  = errors.New("core component cannot be removed")

	// Snapshot errors
	ErrSnapshotCreationFailed   = errors.New("system snapshot creation failed")
	ErrSnapshotRestorationFailed = errors.New("system snapshot restoration failed")
	ErrSnapshotInvalid          = errors.New("system snapshot is invalid or corrupted")
	ErrSnapshotNotFound         = errors.New("system snapshot not found")

	// Configuration errors
	ErrConfigurationMergeFailed = errors.New("configuration merge failed")
	ErrConfigurationBackupFailed = errors.New("configuration backup failed")
	ErrConfigurationInvalid     = errors.New("configuration is invalid")

	// Repository errors
	ErrSessionNotFound = errors.New("installation session not found")
	ErrSnapshotSaveFailed = errors.New("failed to save system snapshot")
)
