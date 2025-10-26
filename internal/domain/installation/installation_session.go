package installation

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// InstallationSession is the aggregate root for installation domain
// It coordinates the entire installation lifecycle and enforces invariants
type InstallationSession struct {
	id                   string
	configuration        InstallationConfiguration
	status               InstallationStatus
	snapshot             *SystemSnapshot
	installedComponents  []*InstalledComponent
	startedAt            time.Time
	completedAt          time.Time
	failureReason        string
}

// NewInstallationSession creates a new installation session aggregate root
func NewInstallationSession(configuration InstallationConfiguration) (*InstallationSession, error) {
	id := uuid.New().String()

	return &InstallationSession{
		id:                  id,
		configuration:       configuration,
		status:              StatusPending,
		installedComponents: make([]*InstalledComponent, 0),
		startedAt:           time.Now(),
	}, nil
}

// ID returns the unique identifier for this session
func (s *InstallationSession) ID() string {
	return s.id
}

// Configuration returns the installation configuration
func (s *InstallationSession) Configuration() InstallationConfiguration {
	return s.configuration
}

// Status returns the current installation status
func (s *InstallationSession) Status() InstallationStatus {
	return s.status
}

// Snapshot returns the system snapshot if available
func (s *InstallationSession) Snapshot() *SystemSnapshot {
	return s.snapshot
}

// InstalledComponents returns a defensive copy of installed components
func (s *InstallationSession) InstalledComponents() []*InstalledComponent {
	components := make([]*InstalledComponent, len(s.installedComponents))
	copy(components, s.installedComponents)
	return components
}

// StartedAt returns when the session was created
func (s *InstallationSession) StartedAt() time.Time {
	return s.startedAt
}

// CompletedAt returns when the session completed (success or failure)
// Returns zero time if not yet completed
func (s *InstallationSession) CompletedAt() time.Time {
	return s.completedAt
}

// FailureReason returns why the installation failed
// Empty string if not failed
func (s *InstallationSession) FailureReason() string {
	return s.failureReason
}

// StartPreparation transitions to preparation phase and attaches snapshot
func (s *InstallationSession) StartPreparation(snapshot *SystemSnapshot) error {
	if !s.status.CanTransitionTo(StatusPreparation) {
		return ErrInvalidStateTransition
	}

	if snapshot == nil {
		return ErrSnapshotInvalid
	}

	s.status = StatusPreparation
	s.snapshot = snapshot
	return nil
}

// StartInstalling transitions to installing phase
func (s *InstallationSession) StartInstalling() error {
	if !s.status.CanTransitionTo(StatusInstalling) {
		return ErrInvalidStateTransition
	}

	s.status = StatusInstalling
	return nil
}

// AddInstalledComponent adds a successfully installed component
// Can only be called during installation phase
func (s *InstallationSession) AddInstalledComponent(component *InstalledComponent) error {
	if s.status != StatusInstalling && s.status != StatusConfiguring {
		return fmt.Errorf("can only add components during installation: %w", ErrSessionNotStarted)
	}

	if component == nil {
		return ErrComponentNotFound
	}

	s.installedComponents = append(s.installedComponents, component)
	return nil
}

// StartConfiguring transitions to configuring phase
func (s *InstallationSession) StartConfiguring() error {
	if !s.status.CanTransitionTo(StatusConfiguring) {
		return ErrInvalidStateTransition
	}

	s.status = StatusConfiguring
	return nil
}

// StartVerifying transitions to verifying phase
func (s *InstallationSession) StartVerifying() error {
	if !s.status.CanTransitionTo(StatusVerifying) {
		return ErrInvalidStateTransition
	}

	s.status = StatusVerifying
	return nil
}

// Complete marks the installation as successfully completed
// Enforces that at least one component was installed
func (s *InstallationSession) Complete() error {
	// Must have installed at least one component
	if len(s.installedComponents) == 0 {
		return fmt.Errorf("cannot complete without installed components: %w", ErrInstallationFailed)
	}

	if !s.status.CanTransitionTo(StatusCompleted) {
		return ErrInvalidStateTransition
	}

	s.status = StatusCompleted
	s.completedAt = time.Now()
	return nil
}

// Fail marks the installation as failed with a reason
func (s *InstallationSession) Fail(reason string) error {
	if s.status.IsTerminal() {
		return ErrSessionAlreadyComplete
	}

	s.status = StatusFailed
	s.failureReason = reason
	s.completedAt = time.Now()
	return nil
}

// IsInProgress returns true if installation is actively running
func (s *InstallationSession) IsInProgress() bool {
	return s.status == StatusPreparation ||
		s.status == StatusDownloading ||
		s.status == StatusInstalling ||
		s.status == StatusConfiguring ||
		s.status == StatusVerifying
}

// IsCompleted returns true if installation finished successfully
func (s *InstallationSession) IsCompleted() bool {
	return s.status == StatusCompleted
}

// IsFailed returns true if installation failed
func (s *InstallationSession) IsFailed() bool {
	return s.status == StatusFailed
}

// Duration returns how long the session has been running
// If completed, returns total duration. Otherwise, duration so far.
func (s *InstallationSession) Duration() time.Duration {
	if !s.completedAt.IsZero() {
		return s.completedAt.Sub(s.startedAt)
	}
	return time.Since(s.startedAt)
}

// String returns human-readable representation
func (s *InstallationSession) String() string {
	componentsInfo := fmt.Sprintf("%d components", s.configuration.ComponentCount())
	if len(s.installedComponents) > 0 {
		componentsInfo = fmt.Sprintf("%d/%d installed",
			len(s.installedComponents), s.configuration.ComponentCount())
	}

	return fmt.Sprintf("Installation session %s (%s, %s)",
		s.id[:8], s.status.String(), componentsInfo)
}
