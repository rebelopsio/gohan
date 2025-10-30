package postinstall

import (
	"context"

	"github.com/rebelopsio/gohan/internal/domain/postinstall"
)

// NetworkInstaller handles network manager installation
type NetworkInstaller struct {
	packageMgr PackageManager
	serviceMgr ServiceManager
}

// NewNetworkInstaller creates a new network installer
func NewNetworkInstaller(packageMgr PackageManager, serviceMgr ServiceManager) *NetworkInstaller {
	return &NetworkInstaller{
		packageMgr: packageMgr,
		serviceMgr: serviceMgr,
	}
}

// Name returns the installer name
func (i *NetworkInstaller) Name() string {
	return "Network Manager"
}

// Component returns the component type
func (i *NetworkInstaller) Component() postinstall.ComponentType {
	return postinstall.ComponentNetwork
}

// Install performs the installation
func (i *NetworkInstaller) Install(ctx context.Context) (postinstall.ComponentResult, error) {
	result := postinstall.NewComponentResult(
		postinstall.ComponentNetwork,
		postinstall.StatusInProgress,
		"Installing Network Manager",
	)

	details := []string{}

	// Required packages
	packages := []string{
		"network-manager",
		"network-manager-gnome",
	}

	// Install packages
	for _, pkg := range packages {
		installed, err := i.packageMgr.IsInstalled(ctx, pkg)
		if err != nil {
			return postinstall.NewComponentResultWithError(
				postinstall.ComponentNetwork,
				"Failed to check package installation",
				err,
			), err
		}

		if !installed {
			if err := i.packageMgr.Install(ctx, pkg); err != nil {
				return postinstall.NewComponentResultWithError(
					postinstall.ComponentNetwork,
					"Failed to install "+pkg,
					err,
				), err
			}
			details = append(details, pkg+" installed")
		} else {
			details = append(details, pkg+" already installed")
		}
	}

	// Enable and start NetworkManager service
	if err := i.serviceMgr.Enable(ctx, "NetworkManager"); err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentNetwork,
			"Failed to enable NetworkManager service",
			err,
		), err
	}
	details = append(details, "NetworkManager service enabled")

	if err := i.serviceMgr.Start(ctx, "NetworkManager"); err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentNetwork,
			"Failed to start NetworkManager service",
			err,
		), err
	}
	details = append(details, "NetworkManager service started")

	return result.
		WithDetails(details...).
		Complete(postinstall.StatusCompleted), nil
}

// Verify checks if network manager is properly configured
func (i *NetworkInstaller) Verify(ctx context.Context) (bool, error) {
	// Check if NetworkManager is enabled
	return i.serviceMgr.IsEnabled(ctx, "NetworkManager")
}

// Rollback reverts the installation
func (i *NetworkInstaller) Rollback(ctx context.Context) error {
	// Just disable the service, don't uninstall
	return i.serviceMgr.Disable(ctx, "NetworkManager")
}
