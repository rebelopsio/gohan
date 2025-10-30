package postinstall

import (
	"context"

	"github.com/rebelopsio/gohan/internal/domain/postinstall"
)

// PostInstallRequest contains parameters for post-installation
type PostInstallRequest struct {
	DisplayManager  postinstall.DisplayManager
	Shell           postinstall.Shell
	ShellTheme      string
	SetupAudio      bool
	SetupNetwork    bool
	Services        []string
	WallpaperDir    string
	ShowProgress    bool
}

// PostInstallResponse contains post-installation results
type PostInstallResponse struct {
	SessionID                 string
	OverallStatus             string
	DisplayManagerConfigured  bool
	ShellConfigured           bool
	AudioConfigured           bool
	NetworkManagerConfigured  bool
	ServicesEnabled           []string
	ServicesFailed            []string
	WallpaperCacheGenerated   bool
	ComponentResults          []ComponentResultDTO
	Recommendations           []string
	DurationMs                int64
}

// ComponentResultDTO represents a component result for display
type ComponentResultDTO struct {
	Component string
	Status    string
	Message   string
	Details   []string
	Error     string
}

// ProgressCallback is called for each component installation
type ProgressCallback func(installerName string, result ComponentResultDTO)

// Installers aggregates all post-installation installers
type Installers struct {
	DisplayManagerInstaller postinstall.ComponentInstaller
	ShellInstaller          postinstall.ComponentInstaller
	AudioInstaller          postinstall.ComponentInstaller
	NetworkInstaller        postinstall.ComponentInstaller
	WallpaperGenerator      postinstall.ComponentInstaller
}

// PostInstallUseCase coordinates post-installation setup
type PostInstallUseCase struct {
	installers Installers
}

// NewPostInstallUseCase creates a new use case instance
func NewPostInstallUseCase(installers Installers) *PostInstallUseCase {
	return &PostInstallUseCase{
		installers: installers,
	}
}

// Execute runs post-installation setup
func (uc *PostInstallUseCase) Execute(ctx context.Context, req PostInstallRequest) (*PostInstallResponse, error) {
	// Create list of installers to run
	installerList := uc.getInstallers(req)

	// Create orchestrator
	orchestrator := postinstall.NewSetupOrchestrator(installerList)

	// Execute setup
	session := orchestrator.RunSetup(ctx)

	// Convert to response
	return uc.buildResponse(session), nil
}

// ExecuteWithProgress runs setup with progress callbacks
func (uc *PostInstallUseCase) ExecuteWithProgress(
	ctx context.Context,
	req PostInstallRequest,
	progressFn ProgressCallback,
) (*PostInstallResponse, error) {
	// Create list of installers to run
	installerList := uc.getInstallers(req)

	// Create orchestrator
	orchestrator := postinstall.NewSetupOrchestrator(installerList)

	// Execute with progress
	session := orchestrator.RunSetupWithProgress(ctx, func(name string, result postinstall.ComponentResult) {
		if progressFn != nil {
			progressFn(name, uc.convertResult(result))
		}
	})

	return uc.buildResponse(session), nil
}

// Rollback reverts post-installation changes
func (uc *PostInstallUseCase) Rollback(ctx context.Context, sessionID string) error {
	// TODO: Implement session persistence and rollback
	// For now, this is a placeholder
	return nil
}

func (uc *PostInstallUseCase) getInstallers(req PostInstallRequest) []postinstall.ComponentInstaller {
	var installers []postinstall.ComponentInstaller

	// Display manager
	if req.DisplayManager != "" && req.DisplayManager != postinstall.DisplayManagerNone {
		if uc.installers.DisplayManagerInstaller != nil {
			installers = append(installers, uc.installers.DisplayManagerInstaller)
		}
	}

	// Shell
	if req.Shell != "" {
		if uc.installers.ShellInstaller != nil {
			installers = append(installers, uc.installers.ShellInstaller)
		}
	}

	// Audio
	if req.SetupAudio {
		if uc.installers.AudioInstaller != nil {
			installers = append(installers, uc.installers.AudioInstaller)
		}
	}

	// Network
	if req.SetupNetwork {
		if uc.installers.NetworkInstaller != nil {
			installers = append(installers, uc.installers.NetworkInstaller)
		}
	}

	// Wallpaper cache
	if req.WallpaperDir != "" {
		if uc.installers.WallpaperGenerator != nil {
			installers = append(installers, uc.installers.WallpaperGenerator)
		}
	}

	return installers
}

func (uc *PostInstallUseCase) buildResponse(session *postinstall.SetupSession) *PostInstallResponse {
	results := session.Results()

	response := &PostInstallResponse{
		SessionID:                session.ID(),
		OverallStatus:            session.OverallStatus().String(),
		DisplayManagerConfigured: false,
		ShellConfigured:          false,
		AudioConfigured:          false,
		NetworkManagerConfigured: false,
		ServicesEnabled:          []string{},
		ServicesFailed:           []string{},
		WallpaperCacheGenerated:  false,
		ComponentResults:         make([]ComponentResultDTO, 0),
		Recommendations:          []string{},
		DurationMs:               session.Duration().Milliseconds(),
	}

	// Process results
	for _, result := range results {
		response.ComponentResults = append(response.ComponentResults, uc.convertResult(result))

		// Update specific flags
		if result.IsSuccess() {
			switch result.Component() {
			case postinstall.ComponentDisplayManager:
				response.DisplayManagerConfigured = true
			case postinstall.ComponentShell:
				response.ShellConfigured = true
			case postinstall.ComponentAudio:
				response.AudioConfigured = true
			case postinstall.ComponentNetwork:
				response.NetworkManagerConfigured = true
			case postinstall.ComponentWallpaper:
				response.WallpaperCacheGenerated = true
			}
		}

		// Collect failures
		if result.IsFailure() {
			errorMsg := result.Message()
			if result.Error() != nil {
				errorMsg += ": " + result.Error().Error()
			}
			response.ServicesFailed = append(response.ServicesFailed, errorMsg)
		}
	}

	// Generate recommendations
	if session.HasFailures() {
		response.Recommendations = append(response.Recommendations,
			"Some components failed to install. Check the error messages above.")
		response.Recommendations = append(response.Recommendations,
			"You can try running post-installation again after resolving the issues.")
	}

	return response
}

func (uc *PostInstallUseCase) convertResult(result postinstall.ComponentResult) ComponentResultDTO {
	errorMsg := ""
	if result.Error() != nil {
		errorMsg = result.Error().Error()
	}

	return ComponentResultDTO{
		Component: result.Component().String(),
		Status:    result.Status().String(),
		Message:   result.Message(),
		Details:   result.Details(),
		Error:     errorMsg,
	}
}
