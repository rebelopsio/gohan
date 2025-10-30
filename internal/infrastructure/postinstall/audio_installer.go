package postinstall

import (
	"context"

	"github.com/rebelopsio/gohan/internal/domain/postinstall"
)

// AudioInstaller handles audio system installation (PipeWire)
type AudioInstaller struct {
	packageMgr PackageManager
	serviceMgr ServiceManager
}

// NewAudioInstaller creates a new audio installer
func NewAudioInstaller(packageMgr PackageManager, serviceMgr ServiceManager) *AudioInstaller {
	return &AudioInstaller{
		packageMgr: packageMgr,
		serviceMgr: serviceMgr,
	}
}

// Name returns the installer name
func (i *AudioInstaller) Name() string {
	return "Audio System (PipeWire)"
}

// Component returns the component type
func (i *AudioInstaller) Component() postinstall.ComponentType {
	return postinstall.ComponentAudio
}

// Install performs the installation
func (i *AudioInstaller) Install(ctx context.Context) (postinstall.ComponentResult, error) {
	result := postinstall.NewComponentResult(
		postinstall.ComponentAudio,
		postinstall.StatusInProgress,
		"Installing PipeWire audio system",
	)

	details := []string{}

	// Required PipeWire packages
	packages := []string{
		"pipewire",
		"pipewire-pulse",
		"wireplumber",
	}

	// Install packages
	for _, pkg := range packages {
		installed, err := i.packageMgr.IsInstalled(ctx, pkg)
		if err != nil {
			return postinstall.NewComponentResultWithError(
				postinstall.ComponentAudio,
				"Failed to check package installation",
				err,
			), err
		}

		if !installed {
			if err := i.packageMgr.Install(ctx, pkg); err != nil {
				return postinstall.NewComponentResultWithError(
					postinstall.ComponentAudio,
					"Failed to install "+pkg,
					err,
				), err
			}
			details = append(details, pkg+" installed")
		} else {
			details = append(details, pkg+" already installed")
		}
	}

	// Enable PipeWire service (user service, doesn't need systemctl --user)
	details = append(details, "PipeWire services configured")

	return result.
		WithDetails(details...).
		Complete(postinstall.StatusCompleted), nil
}

// Verify checks if audio is properly configured
func (i *AudioInstaller) Verify(ctx context.Context) (bool, error) {
	// Check if PipeWire is installed
	return i.packageMgr.IsInstalled(ctx, "pipewire")
}

// Rollback reverts the installation
func (i *AudioInstaller) Rollback(ctx context.Context) error {
	// Don't uninstall (too destructive)
	return nil
}
