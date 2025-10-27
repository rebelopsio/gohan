package packagemanager

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rebelopsio/gohan/internal/domain/installation"
)

// APTManager implements package management operations using APT
// Implements installation.ConflictResolver interface
type APTManager struct {
	dryRun bool
}

// PackageInfo contains information about a package
type PackageInfo struct {
	Name         string
	Version      string
	Architecture string
	Description  string
}

// PackageProgress represents progress for a single package installation
type PackageProgress struct {
	PackageName    string
	Status         string // "started", "installing", "completed", "failed"
	PercentComplete float64
	Error          error
}

// ProgressStatus constants for package installation
const (
	StatusStarted    = "started"
	StatusInstalling = "installing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

// NewAPTManager creates a new APT package manager
func NewAPTManager() *APTManager {
	return &APTManager{
		dryRun: false,
	}
}

// NewAPTManagerDryRun creates a new APT manager in dry-run mode (for testing)
func NewAPTManagerDryRun() *APTManager {
	return &APTManager{
		dryRun: true,
	}
}

// DetectConflicts implements installation.ConflictResolver
// Checks for package conflicts using APT
func (a *APTManager) DetectConflicts(ctx context.Context, components []installation.ComponentSelection) ([]installation.PackageConflict, error) {
	var conflicts []installation.PackageConflict

	// For each component, check if there are conflicting packages
	for _, comp := range components {
		packageName := componentToPackageName(comp.Component())

		// Check for conflicts using dpkg
		cmd := exec.CommandContext(ctx, "dpkg", "-s", packageName)
		output, err := cmd.CombinedOutput()

		if err == nil {
			// Package exists, check for conflicts
			conflicts = append(conflicts, a.checkPackageConflicts(ctx, packageName, string(output))...)
		}
	}

	return conflicts, nil
}

// ResolveConflict implements installation.ConflictResolver
// Applies a resolution strategy to a package conflict
func (a *APTManager) ResolveConflict(ctx context.Context, conflict installation.PackageConflict, strategy installation.ResolutionAction) error {
	switch strategy {
	case installation.ActionRemove:
		// Remove the conflicting package
		return a.RemovePackage(ctx, conflict.ConflictingPackage())

	case installation.ActionSkip:
		// Skip installation - no action needed
		return nil

	case installation.ActionReplace:
		// Remove old, will install new later
		return a.RemovePackage(ctx, conflict.ConflictingPackage())

	case installation.ActionAbort:
		return fmt.Errorf("installation aborted due to package conflict: %s", conflict.String())

	default:
		return fmt.Errorf("unknown resolution action: %s", strategy)
	}
}

// InstallPackage installs a package using APT
func (a *APTManager) InstallPackage(ctx context.Context, packageName, version string) error {
	if packageName == "" {
		return errors.New("package name cannot be empty")
	}

	if a.dryRun {
		return nil
	}

	// If version is specified, append it to package name
	fullPackageName := packageName
	if version != "" {
		fullPackageName = fmt.Sprintf("%s=%s", packageName, version)
	}

	cmd := exec.CommandContext(ctx, "apt-get", "install", "-y", fullPackageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install package %s: %w\nOutput: %s", fullPackageName, err, string(output))
	}

	return nil
}

// RemovePackage removes a package using APT
func (a *APTManager) RemovePackage(ctx context.Context, packageName string) error {
	if packageName == "" {
		return errors.New("package name cannot be empty")
	}

	if a.dryRun {
		return nil
	}

	cmd := exec.CommandContext(ctx, "apt-get", "remove", "-y", packageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove package %s: %w\nOutput: %s", packageName, err, string(output))
	}

	return nil
}

// IsPackageInstalled checks if a package is currently installed
func (a *APTManager) IsPackageInstalled(ctx context.Context, packageName string) (bool, error) {
	if packageName == "" {
		return false, errors.New("package name cannot be empty")
	}

	cmd := exec.CommandContext(ctx, "dpkg-query", "-W", "-f=${Status}", packageName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Package not found
		return false, nil
	}

	status := strings.TrimSpace(string(output))
	return strings.Contains(status, "install ok installed"), nil
}

// UpdatePackageCache updates the APT package cache
func (a *APTManager) UpdatePackageCache(ctx context.Context) error {
	if a.dryRun {
		return nil
	}

	cmd := exec.CommandContext(ctx, "apt-get", "update")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update package cache: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// GetPackageInfo retrieves information about a package
func (a *APTManager) GetPackageInfo(ctx context.Context, packageName string) (*PackageInfo, error) {
	if packageName == "" {
		return nil, errors.New("package name cannot be empty")
	}

	cmd := exec.CommandContext(ctx, "dpkg-query", "-W", "-f=${Package}|${Version}|${Architecture}|${Description}", packageName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("package not found: %s", packageName)
	}

	parts := strings.Split(strings.TrimSpace(string(output)), "|")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid package info format")
	}

	return &PackageInfo{
		Name:         parts[0],
		Version:      parts[1],
		Architecture: parts[2],
		Description:  parts[3],
	}, nil
}

// Helper function to map component names to package names
func componentToPackageName(component installation.ComponentName) string {
	// This is a simplified mapping - real implementation would be more sophisticated
	switch component {
	case installation.ComponentHyprland:
		return "hyprland"
	case installation.ComponentHyprpaper:
		return "hyprpaper"
	case installation.ComponentHyprlock:
		return "hyprlock"
	case installation.ComponentWaybar:
		return "waybar"
	case installation.ComponentFuzzel:
		return "rofi"
	case installation.ComponentKitty:
		return "kitty"
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

// Helper to check for conflicts in package metadata
func (a *APTManager) checkPackageConflicts(ctx context.Context, packageName, dpkgOutput string) []installation.PackageConflict {
	var conflicts []installation.PackageConflict

	// Parse dpkg output for Conflicts field
	lines := strings.Split(dpkgOutput, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Conflicts:") {
			conflictsList := strings.TrimPrefix(line, "Conflicts:")
			conflictsList = strings.TrimSpace(conflictsList)

			// Parse comma-separated conflicts
			conflictPackages := strings.Split(conflictsList, ",")
			for _, conflictPkg := range conflictPackages {
				conflictPkg = strings.TrimSpace(conflictPkg)
				// Remove version constraints
				conflictPkg = strings.Split(conflictPkg, " ")[0]

				if conflictPkg != "" {
					conflict, err := installation.NewPackageConflict(
						packageName,
						conflictPkg,
						"package conflict declared in metadata",
					)
					if err == nil {
						conflicts = append(conflicts, conflict)
					}
				}
			}
		}
	}

	return conflicts
}

// ========================================
// Phase 3.1: Batch Installation Methods
// ========================================

// ArePackagesInstalled checks if multiple packages are installed
// Returns a map of package name to installation status
func (a *APTManager) ArePackagesInstalled(ctx context.Context, packages []string) (map[string]bool, error) {
	result := make(map[string]bool)

	for _, pkg := range packages {
		installed, err := a.IsPackageInstalled(ctx, pkg)
		if err != nil {
			return nil, fmt.Errorf("failed to check package %s: %w", pkg, err)
		}
		result[pkg] = installed
	}

	return result, nil
}

// InstallPackages installs multiple packages with progress reporting
// Progress is reported via the progressChan for each package
func (a *APTManager) InstallPackages(ctx context.Context, packages []string, progressChan chan<- PackageProgress) error {
	if len(packages) == 0 {
		return nil
	}

	// Check context before starting
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error before installation: %w", err)
	}

	totalPackages := len(packages)

	for i, pkg := range packages {
		// Check context before each package
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled during installation: %w", err)
		}

		percentComplete := float64(i) / float64(totalPackages) * 100

		// Report started
		if progressChan != nil {
			progressChan <- PackageProgress{
				PackageName:     pkg,
				Status:          StatusStarted,
				PercentComplete: percentComplete,
			}
		}

		// Report installing
		if progressChan != nil {
			progressChan <- PackageProgress{
				PackageName:     pkg,
				Status:          StatusInstalling,
				PercentComplete: percentComplete + (50.0 / float64(totalPackages)),
			}
		}

		// Install the package
		err := a.InstallPackage(ctx, pkg, "")
		if err != nil {
			// Report failure
			if progressChan != nil {
				progressChan <- PackageProgress{
					PackageName:     pkg,
					Status:          StatusFailed,
					PercentComplete: percentComplete,
					Error:           err,
				}
			}
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}

		// Report completed
		if progressChan != nil {
			progressChan <- PackageProgress{
				PackageName:     pkg,
				Status:          StatusCompleted,
				PercentComplete: float64(i+1) / float64(totalPackages) * 100,
			}
		}
	}

	return nil
}

// InstallProfile installs packages for a specific profile (minimal/recommended/full)
func (a *APTManager) InstallProfile(ctx context.Context, profileName string, progressChan chan<- PackageProgress) error {
	if profileName == "" {
		return errors.New("profile name cannot be empty")
	}

	// Get packages for the profile
	var packages []string
	switch profileName {
	case "minimal":
		packages = installation.GetMinimalProfile().Packages
	case "recommended":
		packages = installation.GetRecommendedProfile().Packages
	case "full":
		packages = installation.GetFullProfile().Packages
	default:
		return fmt.Errorf("unknown profile: %s", profileName)
	}

	// Install packages with progress reporting
	return a.InstallPackages(ctx, packages, progressChan)
}
