package history

import "errors"

// Domain errors for history management
var (
	// Record errors
	ErrRecordNotFound        = errors.New("installation record not found")
	ErrInvalidRecordID       = errors.New("record ID is invalid")
	ErrInvalidSessionID      = errors.New("session ID is invalid")
	ErrSessionNotComplete    = errors.New("session is not in terminal state")
	ErrInvalidRecordedTime   = errors.New("recorded time is invalid")
	ErrMissingFailureDetails = errors.New("failed record must have failure details")
	ErrNoPackagesInstalled   = errors.New("successful record must have installed packages")

	// Metadata errors
	ErrInvalidTimestamp = errors.New("timestamp is invalid")
	ErrInvalidTimeRange = errors.New("time range is invalid (end before start)")
	ErrNoPackages       = errors.New("metadata must have at least one package")

	// Outcome errors
	ErrInvalidOutcome = errors.New("installation outcome is invalid")

	// Package errors
	ErrInvalidPackageName    = errors.New("package name is invalid")
	ErrInvalidPackageVersion = errors.New("package version is invalid")
	ErrInvalidVersion        = errors.New("version is invalid")

	// Failure errors
	ErrInvalidFailureReason = errors.New("failure reason is required")

	// System context errors
	ErrInvalidSystemContext = errors.New("system context is invalid")

	// Period errors
	ErrInvalidPeriod = errors.New("installation period is invalid")

	// Retention errors
	ErrInvalidRetentionPeriod = errors.New("retention period must be at least 1 day")
	ErrRetentionPeriodTooLong = errors.New("retention period exceeds maximum")

	// Archive errors
	ErrInvalidArchive               = errors.New("archive is invalid")
	ErrIncompatibleArchiveVersion   = errors.New("archive version is incompatible")
	ErrCorruptedArchive             = errors.New("archive data is corrupted")
	ErrArchiveSerializationFailed   = errors.New("failed to serialize archive")
	ErrArchiveDeserializationFailed = errors.New("failed to deserialize archive")
)

// Constants for validation
const (
	DefaultRetentionDays = 90
	MaxRetentionDays     = 3650 // 10 years
	ArchiveVersion       = "1.0"
)
