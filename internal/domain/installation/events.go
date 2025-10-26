package installation

import (
	"time"
)

// InstallationStartedEvent signals that an installation session has begun
type InstallationStartedEvent struct {
	occurredAt time.Time
	sessionID  string
}

// NewInstallationStartedEvent creates a new installation started event
func NewInstallationStartedEvent(sessionID string) InstallationStartedEvent {
	return InstallationStartedEvent{
		occurredAt: time.Now(),
		sessionID:  sessionID,
	}
}

func (e InstallationStartedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e InstallationStartedEvent) EventType() string {
	return "installation.started"
}

func (e InstallationStartedEvent) SessionID() string {
	return e.sessionID
}

// InstallationProgressUpdatedEvent signals progress in the installation
type InstallationProgressUpdatedEvent struct {
	occurredAt      time.Time
	sessionID       string
	currentPhase    InstallationStatus
	percentComplete int
	message         string
}

// NewInstallationProgressUpdatedEvent creates a new progress update event
func NewInstallationProgressUpdatedEvent(
	sessionID string,
	currentPhase InstallationStatus,
	percentComplete int,
	message string,
) InstallationProgressUpdatedEvent {
	return InstallationProgressUpdatedEvent{
		occurredAt:      time.Now(),
		sessionID:       sessionID,
		currentPhase:    currentPhase,
		percentComplete: percentComplete,
		message:         message,
	}
}

func (e InstallationProgressUpdatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e InstallationProgressUpdatedEvent) EventType() string {
	return "installation.progress.updated"
}

func (e InstallationProgressUpdatedEvent) SessionID() string {
	return e.sessionID
}

func (e InstallationProgressUpdatedEvent) CurrentPhase() InstallationStatus {
	return e.currentPhase
}

func (e InstallationProgressUpdatedEvent) PercentComplete() int {
	return e.percentComplete
}

func (e InstallationProgressUpdatedEvent) Message() string {
	return e.message
}

// PhaseCompletedEvent signals completion of an installation phase
type PhaseCompletedEvent struct {
	occurredAt time.Time
	sessionID  string
	phase      InstallationStatus
	duration   time.Duration
}

// NewPhaseCompletedEvent creates a new phase completed event
func NewPhaseCompletedEvent(sessionID string, phase InstallationStatus, duration time.Duration) PhaseCompletedEvent {
	return PhaseCompletedEvent{
		occurredAt: time.Now(),
		sessionID:  sessionID,
		phase:      phase,
		duration:   duration,
	}
}

func (e PhaseCompletedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e PhaseCompletedEvent) EventType() string {
	return "installation.phase.completed"
}

func (e PhaseCompletedEvent) SessionID() string {
	return e.sessionID
}

func (e PhaseCompletedEvent) Phase() InstallationStatus {
	return e.phase
}

func (e PhaseCompletedEvent) Duration() time.Duration {
	return e.duration
}

// ComponentInstalledEvent signals successful installation of a component
type ComponentInstalledEvent struct {
	occurredAt time.Time
	sessionID  string
	component  ComponentName
	version    string
}

// NewComponentInstalledEvent creates a new component installed event
func NewComponentInstalledEvent(sessionID string, component ComponentName, version string) ComponentInstalledEvent {
	return ComponentInstalledEvent{
		occurredAt: time.Now(),
		sessionID:  sessionID,
		component:  component,
		version:    version,
	}
}

func (e ComponentInstalledEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e ComponentInstalledEvent) EventType() string {
	return "installation.component.installed"
}

func (e ComponentInstalledEvent) SessionID() string {
	return e.sessionID
}

func (e ComponentInstalledEvent) Component() ComponentName {
	return e.component
}

func (e ComponentInstalledEvent) Version() string {
	return e.version
}

// InstallationCompletedEvent signals successful completion of installation
type InstallationCompletedEvent struct {
	occurredAt          time.Time
	sessionID           string
	duration            time.Duration
	componentsInstalled int
}

// NewInstallationCompletedEvent creates a new installation completed event
func NewInstallationCompletedEvent(
	sessionID string,
	duration time.Duration,
	componentsInstalled int,
) InstallationCompletedEvent {
	return InstallationCompletedEvent{
		occurredAt:          time.Now(),
		sessionID:           sessionID,
		duration:            duration,
		componentsInstalled: componentsInstalled,
	}
}

func (e InstallationCompletedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e InstallationCompletedEvent) EventType() string {
	return "installation.completed"
}

func (e InstallationCompletedEvent) SessionID() string {
	return e.sessionID
}

func (e InstallationCompletedEvent) Duration() time.Duration {
	return e.duration
}

func (e InstallationCompletedEvent) ComponentsInstalled() int {
	return e.componentsInstalled
}

// InstallationFailedEvent signals that installation has failed
type InstallationFailedEvent struct {
	occurredAt  time.Time
	sessionID   string
	phase       InstallationStatus
	reason      string
	recoverable bool
}

// NewInstallationFailedEvent creates a new installation failed event
func NewInstallationFailedEvent(
	sessionID string,
	phase InstallationStatus,
	reason string,
	recoverable bool,
) InstallationFailedEvent {
	return InstallationFailedEvent{
		occurredAt:  time.Now(),
		sessionID:   sessionID,
		phase:       phase,
		reason:      reason,
		recoverable: recoverable,
	}
}

func (e InstallationFailedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e InstallationFailedEvent) EventType() string {
	return "installation.failed"
}

func (e InstallationFailedEvent) SessionID() string {
	return e.sessionID
}

func (e InstallationFailedEvent) Phase() InstallationStatus {
	return e.phase
}

func (e InstallationFailedEvent) Reason() string {
	return e.reason
}

func (e InstallationFailedEvent) IsRecoverable() bool {
	return e.recoverable
}

// RollbackStartedEvent signals that rollback has begun
type RollbackStartedEvent struct {
	occurredAt time.Time
	sessionID  string
	snapshotID string
	reason     string
}

// NewRollbackStartedEvent creates a new rollback started event
func NewRollbackStartedEvent(sessionID, snapshotID, reason string) RollbackStartedEvent {
	return RollbackStartedEvent{
		occurredAt: time.Now(),
		sessionID:  sessionID,
		snapshotID: snapshotID,
		reason:     reason,
	}
}

func (e RollbackStartedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e RollbackStartedEvent) EventType() string {
	return "installation.rollback.started"
}

func (e RollbackStartedEvent) SessionID() string {
	return e.sessionID
}

func (e RollbackStartedEvent) SnapshotID() string {
	return e.snapshotID
}

func (e RollbackStartedEvent) Reason() string {
	return e.reason
}

// RollbackCompletedEvent signals that rollback has completed
type RollbackCompletedEvent struct {
	occurredAt time.Time
	sessionID  string
	snapshotID string
	duration   time.Duration
	success    bool
}

// NewRollbackCompletedEvent creates a new rollback completed event
func NewRollbackCompletedEvent(
	sessionID, snapshotID string,
	duration time.Duration,
	success bool,
) RollbackCompletedEvent {
	return RollbackCompletedEvent{
		occurredAt: time.Now(),
		sessionID:  sessionID,
		snapshotID: snapshotID,
		duration:   duration,
		success:    success,
	}
}

func (e RollbackCompletedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e RollbackCompletedEvent) EventType() string {
	return "installation.rollback.completed"
}

func (e RollbackCompletedEvent) SessionID() string {
	return e.sessionID
}

func (e RollbackCompletedEvent) SnapshotID() string {
	return e.snapshotID
}

func (e RollbackCompletedEvent) Duration() time.Duration {
	return e.duration
}

func (e RollbackCompletedEvent) Success() bool {
	return e.success
}

// ConflictDetectedEvent signals a package conflict was detected
type ConflictDetectedEvent struct {
	occurredAt         time.Time
	sessionID          string
	packageName        string
	conflictingPackage string
	severity           string
}

// NewConflictDetectedEvent creates a new conflict detected event
func NewConflictDetectedEvent(
	sessionID, packageName, conflictingPackage, severity string,
) ConflictDetectedEvent {
	return ConflictDetectedEvent{
		occurredAt:         time.Now(),
		sessionID:          sessionID,
		packageName:        packageName,
		conflictingPackage: conflictingPackage,
		severity:           severity,
	}
}

func (e ConflictDetectedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e ConflictDetectedEvent) EventType() string {
	return "installation.conflict.detected"
}

func (e ConflictDetectedEvent) SessionID() string {
	return e.sessionID
}

func (e ConflictDetectedEvent) PackageName() string {
	return e.packageName
}

func (e ConflictDetectedEvent) ConflictingPackage() string {
	return e.conflictingPackage
}

func (e ConflictDetectedEvent) Severity() string {
	return e.severity
}

// BackupCreatedEvent signals a backup was successfully created
type BackupCreatedEvent struct {
	occurredAt    time.Time
	sessionID     string
	backupPath    string
	filesBackedUp int
	totalSize     int64
}

// NewBackupCreatedEvent creates a new backup created event
func NewBackupCreatedEvent(
	sessionID, backupPath string,
	filesBackedUp int,
	totalSize int64,
) BackupCreatedEvent {
	return BackupCreatedEvent{
		occurredAt:    time.Now(),
		sessionID:     sessionID,
		backupPath:    backupPath,
		filesBackedUp: filesBackedUp,
		totalSize:     totalSize,
	}
}

func (e BackupCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e BackupCreatedEvent) EventType() string {
	return "installation.backup.created"
}

func (e BackupCreatedEvent) SessionID() string {
	return e.sessionID
}

func (e BackupCreatedEvent) BackupPath() string {
	return e.backupPath
}

func (e BackupCreatedEvent) FilesBackedUp() int {
	return e.filesBackedUp
}

func (e BackupCreatedEvent) TotalSize() int64 {
	return e.totalSize
}

// DiskSpaceInsufficientEvent signals insufficient disk space
type DiskSpaceInsufficientEvent struct {
	occurredAt     time.Time
	sessionID      string
	requiredBytes  int64
	availableBytes int64
	path           string
}

// NewDiskSpaceInsufficientEvent creates a new disk space insufficient event
func NewDiskSpaceInsufficientEvent(
	sessionID string,
	requiredBytes, availableBytes int64,
	path string,
) DiskSpaceInsufficientEvent {
	return DiskSpaceInsufficientEvent{
		occurredAt:     time.Now(),
		sessionID:      sessionID,
		requiredBytes:  requiredBytes,
		availableBytes: availableBytes,
		path:           path,
	}
}

func (e DiskSpaceInsufficientEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e DiskSpaceInsufficientEvent) EventType() string {
	return "installation.disk.insufficient"
}

func (e DiskSpaceInsufficientEvent) SessionID() string {
	return e.sessionID
}

func (e DiskSpaceInsufficientEvent) RequiredBytes() int64 {
	return e.requiredBytes
}

func (e DiskSpaceInsufficientEvent) AvailableBytes() int64 {
	return e.availableBytes
}

func (e DiskSpaceInsufficientEvent) Path() string {
	return e.path
}

// NetworkInterruptionEvent signals a network interruption occurred
type NetworkInterruptionEvent struct {
	occurredAt   time.Time
	sessionID    string
	operation    string
	retryable    bool
	errorMessage string
}

// NewNetworkInterruptionEvent creates a new network interruption event
func NewNetworkInterruptionEvent(
	sessionID, operation string,
	retryable bool,
	errorMessage string,
) NetworkInterruptionEvent {
	return NetworkInterruptionEvent{
		occurredAt:   time.Now(),
		sessionID:    sessionID,
		operation:    operation,
		retryable:    retryable,
		errorMessage: errorMessage,
	}
}

func (e NetworkInterruptionEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e NetworkInterruptionEvent) EventType() string {
	return "installation.network.interrupted"
}

func (e NetworkInterruptionEvent) SessionID() string {
	return e.sessionID
}

func (e NetworkInterruptionEvent) Operation() string {
	return e.operation
}

func (e NetworkInterruptionEvent) IsRetryable() bool {
	return e.retryable
}

func (e NetworkInterruptionEvent) ErrorMessage() string {
	return e.errorMessage
}
