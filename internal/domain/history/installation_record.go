package history

import (
	"strings"
	"time"
)

// InstallationRecord is an aggregate root representing an immutable historical fact about an installation
type InstallationRecord struct {
	id             RecordID
	sessionID      string
	outcome        InstallationOutcome
	metadata       InstallationMetadata
	systemContext  SystemContext
	failureDetails *FailureDetails
	recordedAt     time.Time
}

// NewInstallationRecord creates a new installation record
func NewInstallationRecord(
	sessionID string,
	outcome InstallationOutcome,
	metadata InstallationMetadata,
	systemContext SystemContext,
	failureDetails *FailureDetails,
	recordedAt time.Time,
) (InstallationRecord, error) {
	// Generate unique ID
	id, err := NewRecordID()
	if err != nil {
		return InstallationRecord{}, err
	}

	return reconstructInstallationRecord(
		id,
		sessionID,
		outcome,
		metadata,
		systemContext,
		failureDetails,
		recordedAt,
	)
}

// ReconstructInstallationRecord reconstructs a record with a specific ID (for persistence)
func ReconstructInstallationRecord(
	id RecordID,
	sessionID string,
	outcome InstallationOutcome,
	metadata InstallationMetadata,
	systemContext SystemContext,
	failureDetails *FailureDetails,
	recordedAt time.Time,
) (InstallationRecord, error) {
	return reconstructInstallationRecord(
		id,
		sessionID,
		outcome,
		metadata,
		systemContext,
		failureDetails,
		recordedAt,
	)
}

// reconstructInstallationRecord is the internal constructor
func reconstructInstallationRecord(
	id RecordID,
	sessionID string,
	outcome InstallationOutcome,
	metadata InstallationMetadata,
	systemContext SystemContext,
	failureDetails *FailureDetails,
	recordedAt time.Time,
) (InstallationRecord, error) {

	// Trim and validate session ID
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return InstallationRecord{}, ErrInvalidSessionID
	}

	// Validate recorded time
	if recordedAt.IsZero() {
		return InstallationRecord{}, ErrInvalidRecordedTime
	}

	// Validate business rules
	if outcome.IsFailed() && failureDetails == nil {
		return InstallationRecord{}, ErrMissingFailureDetails
	}

	if outcome.IsSuccessful() && metadata.PackageCount() == 0 {
		return InstallationRecord{}, ErrNoPackagesInstalled
	}

	return InstallationRecord{
		id:             id,
		sessionID:      sessionID,
		outcome:        outcome,
		metadata:       metadata,
		systemContext:  systemContext,
		failureDetails: failureDetails,
		recordedAt:     recordedAt,
	}, nil
}

// ID returns the record's unique identifier
func (r InstallationRecord) ID() RecordID {
	return r.id
}

// SessionID returns the original installation session ID
func (r InstallationRecord) SessionID() string {
	return r.sessionID
}

// Outcome returns the installation outcome
func (r InstallationRecord) Outcome() InstallationOutcome {
	return r.outcome
}

// Metadata returns the installation metadata
func (r InstallationRecord) Metadata() InstallationMetadata {
	return r.metadata
}

// SystemContext returns the system context at installation time
func (r InstallationRecord) SystemContext() SystemContext {
	return r.systemContext
}

// FailureDetails returns failure details if present
func (r InstallationRecord) FailureDetails() *FailureDetails {
	return r.failureDetails
}

// RecordedAt returns when the record was created
func (r InstallationRecord) RecordedAt() time.Time {
	return r.recordedAt
}

// WasSuccessful returns true if installation was successful
func (r InstallationRecord) WasSuccessful() bool {
	return r.outcome.IsSuccessful()
}

// WasFailed returns true if installation failed
func (r InstallationRecord) WasFailed() bool {
	return r.outcome.IsFailed()
}

// WasRolledBack returns true if installation was rolled back
func (r InstallationRecord) WasRolledBack() bool {
	return r.outcome.IsRolledBack()
}

// HasFailureDetails returns true if failure details are present
func (r InstallationRecord) HasFailureDetails() bool {
	return r.failureDetails != nil
}

// PackageName returns the target package name (convenience method)
func (r InstallationRecord) PackageName() string {
	return r.metadata.PackageName()
}

// TargetVersion returns the target version (convenience method)
func (r InstallationRecord) TargetVersion() string {
	return r.metadata.TargetVersion()
}

// InstalledAt returns when installation started (convenience method)
func (r InstallationRecord) InstalledAt() time.Time {
	return r.metadata.InstalledAt()
}

// Duration returns the installation duration (convenience method)
func (r InstallationRecord) Duration() time.Duration {
	return time.Duration(r.metadata.DurationMs()) * time.Millisecond
}

// PackageCount returns the number of installed packages (convenience method)
func (r InstallationRecord) PackageCount() int {
	return r.metadata.PackageCount()
}
