package postinstall

import (
	"context"

	"github.com/rebelopsio/gohan/internal/infrastructure/installation/packagemanager"
)

// APTPackageManagerAdapter adapts the existing APTManager to the PackageManager interface
type APTPackageManagerAdapter struct {
	aptManager *packagemanager.APTManager
}

// NewAPTPackageManagerAdapter creates a new adapter
func NewAPTPackageManagerAdapter() *APTPackageManagerAdapter {
	return &APTPackageManagerAdapter{
		aptManager: packagemanager.NewAPTManager(),
	}
}

// Install installs one or more packages
func (a *APTPackageManagerAdapter) Install(ctx context.Context, packages ...string) error {
	// Use InstallPackages without progress channel
	return a.aptManager.InstallPackages(ctx, packages, nil)
}

// IsInstalled checks if a package is installed
func (a *APTPackageManagerAdapter) IsInstalled(ctx context.Context, pkg string) (bool, error) {
	return a.aptManager.IsPackageInstalled(ctx, pkg)
}

// Remove removes one or more packages
func (a *APTPackageManagerAdapter) Remove(ctx context.Context, packages ...string) error {
	// Remove packages one by one
	for _, pkg := range packages {
		if err := a.aptManager.RemovePackage(ctx, pkg); err != nil {
			return err
		}
	}
	return nil
}
