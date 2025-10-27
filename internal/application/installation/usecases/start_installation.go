package usecases

import (
	"context"
	"fmt"

	"github.com/rebelopsio/gohan/internal/application/installation/dto"
	"github.com/rebelopsio/gohan/internal/domain/installation"
)

// StartInstallationUseCase handles starting a new installation session
type StartInstallationUseCase struct{
	sessionRepo installation.InstallationSessionRepository
}

// NewStartInstallationUseCase creates a new start installation use case
func NewStartInstallationUseCase(sessionRepo installation.InstallationSessionRepository) *StartInstallationUseCase {
	return &StartInstallationUseCase{
		sessionRepo: sessionRepo,
	}
}

// Execute starts a new installation session
func (u *StartInstallationUseCase) Execute(ctx context.Context, request dto.InstallationRequest) (*dto.InstallationResponse, error) {
	// Validate request
	if len(request.Components) == 0 {
		return nil, fmt.Errorf("at least one component required: %w", installation.ErrInvalidConfiguration)
	}

	// Convert DTOs to domain objects
	components, err := u.convertComponents(request.Components)
	if err != nil {
		return nil, err
	}

	// Convert GPU support if provided
	var gpuSupport *installation.GPUSupport
	if request.GPU != nil {
		gpu, err := u.convertGPUSupport(request.GPU)
		if err != nil {
			return nil, err
		}
		gpuSupport = &gpu
	}

	// Create disk space
	diskSpace, err := installation.NewDiskSpace(request.AvailableSpace, request.RequiredSpace)
	if err != nil {
		return nil, err
	}

	// Create installation configuration
	config, err := installation.NewInstallationConfiguration(
		components,
		gpuSupport,
		diskSpace,
		request.MergeExistingConfig,
	)
	if err != nil {
		return nil, err
	}

	// Create installation session
	session, err := installation.NewInstallationSession(config)
	if err != nil {
		return nil, err
	}

	// Save session to repository
	if err := u.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	// Build response
	response := &dto.InstallationResponse{
		SessionID:      session.ID(),
		Status:         session.Status().String(),
		Message:        "Installation session created successfully",
		StartedAt:      session.StartedAt().Format("2006-01-02T15:04:05Z07:00"),
		ComponentCount: config.ComponentCount(),
	}

	return response, nil
}

// convertComponents converts DTO components to domain component selections
func (u *StartInstallationUseCase) convertComponents(dtoComponents []dto.ComponentRequest) ([]installation.ComponentSelection, error) {
	components := make([]installation.ComponentSelection, 0, len(dtoComponents))

	for _, comp := range dtoComponents {
		if comp.Name == "" {
			return nil, fmt.Errorf("component name cannot be empty: %w", installation.ErrInvalidComponentSelection)
		}
		if comp.Version == "" {
			return nil, fmt.Errorf("component version cannot be empty for %s: %w", comp.Name, installation.ErrInvalidComponentSelection)
		}

		componentName := ConvertComponentName(comp.Name)

		// Create package info if size is provided
		var packageInfo *installation.PackageInfo
		if comp.SizeBytes > 0 {
			pkgName := comp.PackageName
			if pkgName == "" {
				pkgName = comp.Name
			}

			pkg, err := installation.NewPackageInfo(pkgName, comp.Version, comp.SizeBytes, nil)
			if err == nil {
				packageInfo = &pkg
			}
		}

		selection, err := installation.NewComponentSelection(componentName, comp.Version, packageInfo)
		if err != nil {
			return nil, fmt.Errorf("invalid component %s: %w", comp.Name, err)
		}

		components = append(components, selection)
	}

	return components, nil
}

// convertGPUSupport converts DTO GPU request to domain GPU support
func (u *StartInstallationUseCase) convertGPUSupport(dtoGPU *dto.GPURequest) (installation.GPUSupport, error) {
	driverComponent := installation.ComponentName("")
	if dtoGPU.RequiresDriver && dtoGPU.DriverName != "" {
		driverComponent = ConvertComponentName(dtoGPU.DriverName)
	}

	return installation.NewGPUSupport(dtoGPU.Vendor, dtoGPU.RequiresDriver, driverComponent)
}

// ConvertComponentName converts a string component name to ComponentName enum
func ConvertComponentName(name string) installation.ComponentName {
	switch name {
	case "hyprland":
		return installation.ComponentHyprland
	case "hyprpaper":
		return installation.ComponentHyprpaper
	case "hyprlock":
		return installation.ComponentHyprlock
	case "waybar":
		return installation.ComponentWaybar
	case "rofi":
		return installation.ComponentRofi
	case "kitty":
		return installation.ComponentKitty
	case "default_config":
		return installation.ComponentDefaultConfig
	case "amd_driver":
		return installation.ComponentAMDDriver
	case "nvidia_driver":
		return installation.ComponentNVIDIADriver
	case "intel_driver":
		return installation.ComponentIntelDriver
	default:
		// Return as-is for unknown components
		return installation.ComponentName(name)
	}
}
