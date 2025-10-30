package postinstall

import "context"

// ComponentInstaller defines the interface for installing/configuring components
type ComponentInstaller interface {
	// Name returns the installer name
	Name() string

	// Component returns the component type this installer handles
	Component() ComponentType

	// Install performs the installation/configuration
	Install(ctx context.Context) (ComponentResult, error)

	// Verify checks if the component is properly installed
	Verify(ctx context.Context) (bool, error)

	// Rollback reverts the installation
	Rollback(ctx context.Context) error
}

// SetupOrchestrator coordinates the installation of multiple components
type SetupOrchestrator struct {
	installers []ComponentInstaller
}

// NewSetupOrchestrator creates a new setup orchestrator
func NewSetupOrchestrator(installers []ComponentInstaller) *SetupOrchestrator {
	return &SetupOrchestrator{
		installers: installers,
	}
}

// RunSetup executes all installers
func (o *SetupOrchestrator) RunSetup(ctx context.Context) *SetupSession {
	session := NewSetupSession()

	for _, installer := range o.installers {
		result, err := installer.Install(ctx)

		if err != nil {
			result = NewComponentResultWithError(
				installer.Component(),
				"Installation failed",
				err,
			)
		}

		session.AddResult(result)

		// Add rollback action if successful
		if result.IsSuccess() {
			session.AddRollbackAction(RollbackAction{
				Component:   installer.Component(),
				Description: "Rollback " + installer.Name(),
				UndoFunc:    func() error { return installer.Rollback(ctx) },
			})
		}

		// Stop on failure if it's critical
		if result.IsFailure() {
			break
		}
	}

	session.Complete()
	return session
}

// RunSetupWithProgress executes all installers with progress callbacks
func (o *SetupOrchestrator) RunSetupWithProgress(
	ctx context.Context,
	progressFn func(installerName string, result ComponentResult),
) *SetupSession {
	session := NewSetupSession()

	for _, installer := range o.installers {
		result, err := installer.Install(ctx)

		if err != nil {
			result = NewComponentResultWithError(
				installer.Component(),
				"Installation failed",
				err,
			)
		}

		session.AddResult(result)

		// Call progress callback
		if progressFn != nil {
			progressFn(installer.Name(), result)
		}

		// Add rollback action if successful
		if result.IsSuccess() {
			session.AddRollbackAction(RollbackAction{
				Component:   installer.Component(),
				Description: "Rollback " + installer.Name(),
				UndoFunc:    func() error { return installer.Rollback(ctx) },
			})
		}

		// Stop on failure if it's critical
		if result.IsFailure() {
			break
		}
	}

	session.Complete()
	return session
}
