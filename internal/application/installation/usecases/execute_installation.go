package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/rebelopsio/gohan/internal/application/installation/dto"
	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/rebelopsio/gohan/internal/domain/installation"
)

// PackageManager defines the interface for installing packages
type PackageManager interface {
	InstallPackage(ctx context.Context, packageName, version string) error
	IsPackageInstalled(ctx context.Context, packageName string) (bool, error)
}

// HistoryRecorder defines the interface for recording installation history
type HistoryRecorder interface {
	RecordInstallation(ctx context.Context, session *installation.InstallationSession) (history.RecordID, error)
}

// ProgressCallback is called during installation to report progress
type ProgressCallback func(phase string, percent int, message string, componentsInstalled, componentsTotal int)

// ExecuteInstallationUseCase handles executing an installation session
type ExecuteInstallationUseCase struct {
	sessionRepo       installation.InstallationSessionRepository
	conflictResolver  installation.ConflictResolver
	progressEstimator installation.ProgressEstimator
	configMerger      installation.ConfigurationMerger
	packageManager    PackageManager
	historyRecorder   HistoryRecorder
}

// NewExecuteInstallationUseCase creates a new execute installation use case
func NewExecuteInstallationUseCase(
	sessionRepo installation.InstallationSessionRepository,
	conflictResolver installation.ConflictResolver,
	progressEstimator installation.ProgressEstimator,
	configMerger installation.ConfigurationMerger,
	packageManager PackageManager,
	historyRecorder HistoryRecorder,
) *ExecuteInstallationUseCase {
	return &ExecuteInstallationUseCase{
		sessionRepo:       sessionRepo,
		conflictResolver:  conflictResolver,
		progressEstimator: progressEstimator,
		configMerger:      configMerger,
		packageManager:    packageManager,
		historyRecorder:   historyRecorder,
	}
}

// Execute executes an installation session
// The progressCallback parameter is optional and will be called with progress updates
func (u *ExecuteInstallationUseCase) Execute(ctx context.Context, sessionID string, progressCallback ProgressCallback) (*dto.InstallationProgressResponse, error) {
	// Retrieve the session
	session, err := u.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Get total components for progress reporting
	totalComponents := len(session.Configuration().Components())

	// Report initial progress
	if progressCallback != nil {
		progressCallback("Starting Preparation", 0, "Initializing installation session", 0, totalComponents)
	}

	// Create a system snapshot before starting
	// TODO: In real implementation, capture actual system state
	config := session.Configuration()

	if progressCallback != nil {
		progressCallback("Creating Snapshot", 5, "Creating system snapshot", 0, totalComponents)
	}

	snapshot, err := installation.NewSystemSnapshot(
		"/var/lib/gohan/snapshots",
		config.DiskSpace(),
		[]string{}, // Packages would be captured from actual system
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create system snapshot: %w", err)
	}

	// Start the installation process
	if err := session.StartPreparation(snapshot); err != nil {
		// Session might already be in progress, continue
		if err != installation.ErrInvalidStateTransition {
			return nil, err
		}
	}

	// Save updated session state
	if err := u.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session state: %w", err)
	}

	// Detect conflicts
	if progressCallback != nil {
		progressCallback("Checking Requirements", 10, "Detecting package conflicts", 0, totalComponents)
	}

	conflicts, err := u.conflictResolver.DetectConflicts(ctx, config.Components())
	if err != nil {
		return u.handleInstallationError(ctx, session, fmt.Sprintf("conflict detection failed: %v", err))
	}

	// Resolve conflicts if any
	if len(conflicts) > 0 {
		if progressCallback != nil {
			progressCallback("Resolving Conflicts", 15, fmt.Sprintf("Resolving %d package conflicts", len(conflicts)), 0, totalComponents)
		}

		for _, conflict := range conflicts {
			// Default strategy: remove conflicting package
			if err := u.conflictResolver.ResolveConflict(ctx, conflict, installation.ActionRemove); err != nil {
				return u.handleInstallationError(ctx, session, fmt.Sprintf("conflict resolution failed: %v", err))
			}
		}
	}

	// Start installing phase
	if err := session.StartInstalling(); err != nil {
		if err != installation.ErrInvalidStateTransition {
			return u.handleInstallationError(ctx, session, fmt.Sprintf("failed to start installing: %v", err))
		}
	}

	if err := u.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session state: %w", err)
	}

	// Install each component
	components := config.Components()
	for i, comp := range components {
		// Extract package name and version
		packageName := componentToPackageName(comp.Component())
		version := comp.Version()

		// Calculate progress percentage (20-80% range for installations)
		// Each component gets equal portion of the 60% range
		baseProgress := 20
		progressRange := 60
		componentProgress := baseProgress + (progressRange * i / len(components))

		if progressCallback != nil {
			progressCallback(
				"Installing Components",
				componentProgress,
				fmt.Sprintf("Installing %s (%d/%d)", packageName, i+1, len(components)),
				i,
				totalComponents,
			)
		}

		// Install the package
		if err := u.packageManager.InstallPackage(ctx, packageName, version); err != nil {
			return u.handleInstallationError(ctx, session, fmt.Sprintf("failed to install %s: %v", packageName, err))
		}

		// Create installed component
		var pkgInfo *installation.PackageInfo
		if comp.PackageInfo() != nil {
			pkg := *comp.PackageInfo()
			pkgInfo = &pkg
		}

		installedComp, err := installation.NewInstalledComponent(
			comp.Component(),
			version,
			pkgInfo,
		)
		if err != nil {
			return u.handleInstallationError(ctx, session, fmt.Sprintf("failed to create installed component: %v", err))
		}

		// Add to session
		if err := session.AddInstalledComponent(installedComp); err != nil {
			return u.handleInstallationError(ctx, session, fmt.Sprintf("failed to add installed component: %v", err))
		}

		// Calculate progress
		progress := u.progressEstimator.CalculatePhaseProgress(
			installation.StatusInstalling,
			len(components),
			i+1,
		)

		// Emit progress event (could be published to event bus in real implementation)
		_ = installation.NewInstallationProgressUpdatedEvent(
			session.ID(),
			installation.StatusInstalling,
			progress,
			fmt.Sprintf("Installed %s", comp.Component()),
		)

		// Report completion of this component
		if progressCallback != nil {
			progressCallback(
				"Installing Components",
				baseProgress + (progressRange * (i+1) / len(components)),
				fmt.Sprintf("Installed %s successfully", packageName),
				i+1,
				totalComponents,
			)
		}

		// Save progress
		if err := u.sessionRepo.Save(ctx, session); err != nil {
			return nil, fmt.Errorf("failed to save session state: %w", err)
		}
	}

	// Move to configuring phase
	if progressCallback != nil {
		progressCallback("Configuring", 80, "Applying configuration files", len(components), totalComponents)
	}

	if err := session.StartConfiguring(); err != nil {
		if err != installation.ErrInvalidStateTransition {
			return u.handleInstallationError(ctx, session, fmt.Sprintf("failed to start configuring: %v", err))
		}
	}

	if err := u.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session state: %w", err)
	}

	// Move to verifying phase
	if progressCallback != nil {
		progressCallback("Verifying", 90, "Verifying installation", len(components), totalComponents)
	}

	if err := session.StartVerifying(); err != nil {
		if err != installation.ErrInvalidStateTransition {
			return u.handleInstallationError(ctx, session, fmt.Sprintf("failed to start verifying: %v", err))
		}
	}

	if err := u.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session state: %w", err)
	}

	// Complete the installation
	if progressCallback != nil {
		progressCallback("Finalizing", 95, "Cleaning up temporary files", len(components), totalComponents)
	}

	if err := session.Complete(); err != nil {
		return u.handleInstallationError(ctx, session, fmt.Sprintf("failed to complete installation: %v", err))
	}

	// Calculate duration
	duration := time.Since(session.StartedAt())

	// Emit completion event (could be published to event bus in real implementation)
	_ = installation.NewInstallationCompletedEvent(
		session.ID(),
		duration,
		len(session.InstalledComponents()),
	)

	if err := u.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session state: %w", err)
	}

	// Record installation to history
	if u.historyRecorder != nil {
		if _, err := u.historyRecorder.RecordInstallation(ctx, session); err != nil {
			// Log error but don't fail the installation
			// History recording is not critical to installation success
			fmt.Printf("Warning: failed to record installation to history: %v\n", err)
		}
	}

	// Calculate elapsed time
	elapsedTime := time.Since(session.StartedAt())
	estimatedRemaining := u.progressEstimator.EstimateRemainingTime(
		session.Status(),
		100,
		elapsedTime,
	)

	// Build response
	response := &dto.InstallationProgressResponse{
		SessionID:           session.ID(),
		Status:              session.Status().String(),
		CurrentPhase:        session.Status().String(),
		PercentComplete:     100,
		Message:             "Installation completed successfully",
		EstimatedRemaining:  estimatedRemaining.String(),
		ComponentsInstalled: len(session.InstalledComponents()),
		ComponentsTotal:     len(components),
	}

	return response, nil
}

// handleInstallationError marks the session as failed and returns an error response
func (u *ExecuteInstallationUseCase) handleInstallationError(
	ctx context.Context,
	session *installation.InstallationSession,
	errorMessage string,
) (*dto.InstallationProgressResponse, error) {
	// Mark session as failed
	_ = session.Fail(errorMessage)

	// Emit failure event (could be published to event bus in real implementation)
	_ = installation.NewInstallationFailedEvent(
		session.ID(),
		session.Status(),
		errorMessage,
		false, // not recoverable by default
	)

	// Save failed state
	_ = u.sessionRepo.Save(ctx, session)

	// Record failed installation to history
	if u.historyRecorder != nil {
		if _, err := u.historyRecorder.RecordInstallation(ctx, session); err != nil {
			// Log error but don't fail the installation
			fmt.Printf("Warning: failed to record failed installation to history: %v\n", err)
		}
	}

	// Return response (not an error, but a failed installation)
	response := &dto.InstallationProgressResponse{
		SessionID:           session.ID(),
		Status:              "failed",
		CurrentPhase:        session.Status().String(),
		PercentComplete:     0,
		Message:             errorMessage,
		EstimatedRemaining:  "0s",
		ComponentsInstalled: len(session.InstalledComponents()),
		ComponentsTotal:     len(session.Configuration().Components()),
	}

	return response, nil
}

// componentToPackageName converts component name to package name
// This is a simple mapping - could be externalized to configuration
func componentToPackageName(component installation.ComponentName) string {
	switch component {
	case installation.ComponentHyprland:
		return "hyprland"
	case installation.ComponentHyprpaper:
		return "hyprpaper"
	case installation.ComponentHyprlock:
		return "hyprlock"
	case installation.ComponentWaybar:
		return "waybar"
	case installation.ComponentRofi:
		return "rofi"
	case installation.ComponentKitty:
		return "kitty"
	case installation.ComponentDefaultConfig:
		return "gohan-default-config"
	case installation.ComponentAMDDriver:
		return "xserver-xorg-video-amdgpu"
	case installation.ComponentNVIDIADriver:
		return "nvidia-driver"
	case installation.ComponentIntelDriver:
		return "xserver-xorg-video-intel"
	default:
		return string(component)
	}
}
