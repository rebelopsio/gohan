package installation

import "time"

// InstallationStatus represents the current state of an installation session
type InstallationStatus string

const (
	StatusPending     InstallationStatus = "pending"      // Not yet started
	StatusPreparation InstallationStatus = "preparation"  // Taking snapshot, checking space
	StatusDownloading InstallationStatus = "downloading"  // Downloading packages
	StatusInstalling  InstallationStatus = "installing"   // Installing packages
	StatusConfiguring InstallationStatus = "configuring"  // Setting up configs
	StatusVerifying   InstallationStatus = "verifying"    // Checking installation
	StatusCompleted   InstallationStatus = "completed"    // Successfully finished
	StatusFailed      InstallationStatus = "failed"       // Failed with error
	StatusRollingBack InstallationStatus = "rolling_back" // Restoring previous state
	StatusRolledBack  InstallationStatus = "rolled_back"  // Rollback completed
)

// InstallationPhase represents distinct steps in the installation process
type InstallationPhase string

const (
	PhaseSnapshot        InstallationPhase = "snapshot"          // Creating system snapshot
	PhaseDiskCheck       InstallationPhase = "disk_check"        // Verifying disk space
	PhaseConflictCheck   InstallationPhase = "conflict_check"    // Checking package conflicts
	PhaseBackup          InstallationPhase = "backup"            // Backing up configs
	PhaseDownload        InstallationPhase = "download"          // Downloading packages
	PhaseInstallCore     InstallationPhase = "install_core"      // Installing Hyprland core
	PhaseInstallGPU      InstallationPhase = "install_gpu"       // Installing GPU support
	PhaseInstallOptional InstallationPhase = "install_optional"  // Installing optional components
	PhaseConfiguration   InstallationPhase = "configuration"     // Setting up configs
	PhaseDisplayManager  InstallationPhase = "display_manager"   // Configuring login screen
	PhaseVerification    InstallationPhase = "verification"      // Verifying installation
)

// ComponentName identifies specific components that can be installed
type ComponentName string

const (
	ComponentHyprland      ComponentName = "hyprland"        // Core compositor
	ComponentHyprpaper     ComponentName = "hyprpaper"       // Wallpaper utility
	ComponentHyprlock      ComponentName = "hyprlock"        // Screen locker
	ComponentWaybar        ComponentName = "waybar"          // Status bar
	ComponentRofi          ComponentName = "rofi"            // Application launcher
	ComponentKitty         ComponentName = "kitty"           // Terminal emulator
	ComponentDefaultConfig ComponentName = "default_config"  // Default configuration files
	ComponentAMDDriver     ComponentName = "amd_driver"      // AMD GPU drivers
	ComponentNVIDIADriver  ComponentName = "nvidia_driver"   // NVIDIA GPU drivers
	ComponentIntelDriver   ComponentName = "intel_driver"    // Intel GPU drivers
)

// ResolutionAction defines how to resolve package conflicts
type ResolutionAction string

const (
	ActionRemove  ResolutionAction = "remove"  // Remove conflicting package
	ActionReplace ResolutionAction = "replace" // Replace with compatible version
	ActionSkip    ResolutionAction = "skip"    // Skip installation of conflicting component
	ActionAbort   ResolutionAction = "abort"   // Cancel installation
)

// DomainEvent is the base interface for all domain events
type DomainEvent interface {
	OccurredAt() time.Time
	EventType() string
}

// IsCore returns true if the component is required for Hyprland
func (c ComponentName) IsCore() bool {
	return c == ComponentHyprland
}

// IsDriver returns true if the component is a GPU driver
func (c ComponentName) IsDriver() bool {
	return c == ComponentAMDDriver ||
		c == ComponentNVIDIADriver ||
		c == ComponentIntelDriver
}

// String returns the string representation of ComponentName
func (c ComponentName) String() string {
	return string(c)
}

// String returns the string representation of InstallationStatus
func (s InstallationStatus) String() string {
	return string(s)
}

// String returns the string representation of InstallationPhase
func (p InstallationPhase) String() string {
	return string(p)
}

// String returns the string representation of ResolutionAction
func (a ResolutionAction) String() string {
	return string(a)
}

// IsTerminal returns true if this is a final state (completed, failed, rolled back)
func (s InstallationStatus) IsTerminal() bool {
	return s == StatusCompleted || s == StatusFailed || s == StatusRolledBack
}

// CanTransitionTo checks if transitioning to the new status is valid
func (s InstallationStatus) CanTransitionTo(newStatus InstallationStatus) bool {
	// Cannot transition from terminal states
	if s.IsTerminal() {
		return false
	}

	// Can always transition to failed or rolling back
	if newStatus == StatusFailed || newStatus == StatusRollingBack {
		return true
	}

	// Define valid transitions
	validTransitions := map[InstallationStatus][]InstallationStatus{
		StatusPending:     {StatusPreparation},
		StatusPreparation: {StatusDownloading, StatusInstalling},
		StatusDownloading: {StatusInstalling},
		StatusInstalling:  {StatusConfiguring},
		StatusConfiguring: {StatusVerifying},
		StatusVerifying:   {StatusCompleted},
		StatusRollingBack: {StatusRolledBack},
	}

	allowed, exists := validTransitions[s]
	if !exists {
		return false
	}

	for _, allowedStatus := range allowed {
		if newStatus == allowedStatus {
			return true
		}
	}

	return false
}
