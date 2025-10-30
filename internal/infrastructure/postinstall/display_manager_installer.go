package postinstall

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebelopsio/gohan/internal/domain/postinstall"
)

// DisplayManagerInstaller handles display manager installation
type DisplayManagerInstaller struct {
	dm          postinstall.DisplayManager
	packageMgr  PackageManager
	serviceMgr  ServiceManager
	homeDir     string
	previousDM  string
}

// PackageManager interface for installing packages
type PackageManager interface {
	Install(ctx context.Context, packages ...string) error
	IsInstalled(ctx context.Context, pkg string) (bool, error)
	Remove(ctx context.Context, packages ...string) error
}

// ServiceManager interface for managing systemd services
type ServiceManager interface {
	Enable(ctx context.Context, service string) error
	Disable(ctx context.Context, service string) error
	Start(ctx context.Context, service string) error
	Stop(ctx context.Context, service string) error
	IsEnabled(ctx context.Context, service string) (bool, error)
	IsActive(ctx context.Context, service string) (bool, error)
}

// NewDisplayManagerInstaller creates a new display manager installer
func NewDisplayManagerInstaller(
	dm postinstall.DisplayManager,
	packageMgr PackageManager,
	serviceMgr ServiceManager,
) *DisplayManagerInstaller {
	homeDir, _ := os.UserHomeDir()
	return &DisplayManagerInstaller{
		dm:         dm,
		packageMgr: packageMgr,
		serviceMgr: serviceMgr,
		homeDir:    homeDir,
	}
}

// Name returns the installer name
func (i *DisplayManagerInstaller) Name() string {
	return fmt.Sprintf("Display Manager (%s)", i.dm)
}

// Component returns the component type
func (i *DisplayManagerInstaller) Component() postinstall.ComponentType {
	return postinstall.ComponentDisplayManager
}

// Install performs the installation
func (i *DisplayManagerInstaller) Install(ctx context.Context) (postinstall.ComponentResult, error) {
	result := postinstall.NewComponentResult(
		postinstall.ComponentDisplayManager,
		postinstall.StatusInProgress,
		"Installing display manager",
	)

	switch i.dm {
	case postinstall.DisplayManagerSDDM:
		return i.installSDDM(ctx, result)
	case postinstall.DisplayManagerGDM:
		return i.installGDM(ctx, result)
	case postinstall.DisplayManagerTTY:
		return i.installTTY(ctx, result)
	case postinstall.DisplayManagerNone:
		return result.
			WithDetails("No display manager selected").
			Complete(postinstall.StatusSkipped), nil
	default:
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentDisplayManager,
			"Invalid display manager",
			fmt.Errorf("unknown display manager: %s", i.dm),
		), nil
	}
}

func (i *DisplayManagerInstaller) installSDDM(ctx context.Context, result postinstall.ComponentResult) (postinstall.ComponentResult, error) {
	// Check if already installed
	installed, err := i.packageMgr.IsInstalled(ctx, "sddm")
	if err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentDisplayManager,
			"Failed to check SDDM installation status",
			err,
		), err
	}

	details := []string{}

	// Install if not present
	if !installed {
		if err := i.packageMgr.Install(ctx, "sddm"); err != nil {
			return postinstall.NewComponentResultWithError(
				postinstall.ComponentDisplayManager,
				"Failed to install SDDM",
				err,
			), err
		}
		details = append(details, "SDDM package installed")
	} else {
		details = append(details, "SDDM already installed")
	}

	// Create Hyprland session file
	sessionFile := "/usr/share/wayland-sessions/hyprland.desktop"
	if err := i.createHyprlandSession(sessionFile); err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentDisplayManager,
			"Failed to create Hyprland session file",
			err,
		), err
	}
	details = append(details, "Hyprland session file created")

	// Enable SDDM service
	if err := i.serviceMgr.Enable(ctx, "sddm"); err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentDisplayManager,
			"Failed to enable SDDM service",
			err,
		), err
	}
	details = append(details, "SDDM service enabled")

	return result.
		WithDetails(details...).
		Complete(postinstall.StatusCompleted), nil
}

func (i *DisplayManagerInstaller) installGDM(ctx context.Context, result postinstall.ComponentResult) (postinstall.ComponentResult, error) {
	// Check if already installed
	installed, err := i.packageMgr.IsInstalled(ctx, "gdm3")
	if err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentDisplayManager,
			"Failed to check GDM installation status",
			err,
		), err
	}

	details := []string{}

	// Install if not present
	if !installed {
		if err := i.packageMgr.Install(ctx, "gdm3"); err != nil {
			return postinstall.NewComponentResultWithError(
				postinstall.ComponentDisplayManager,
				"Failed to install GDM",
				err,
			), err
		}
		details = append(details, "GDM package installed")
	} else {
		details = append(details, "GDM already installed")
	}

	// Create Hyprland session file
	sessionFile := "/usr/share/wayland-sessions/hyprland.desktop"
	if err := i.createHyprlandSession(sessionFile); err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentDisplayManager,
			"Failed to create Hyprland session file",
			err,
		), err
	}
	details = append(details, "Hyprland session file created")

	// Enable GDM service
	if err := i.serviceMgr.Enable(ctx, "gdm"); err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentDisplayManager,
			"Failed to enable GDM service",
			err,
		), err
	}
	details = append(details, "GDM service enabled")

	return result.
		WithDetails(details...).
		Complete(postinstall.StatusCompleted), nil
}

func (i *DisplayManagerInstaller) installTTY(ctx context.Context, result postinstall.ComponentResult) (postinstall.ComponentResult, error) {
	// Create launch script in user's home directory
	scriptPath := filepath.Join(i.homeDir, ".local/bin/start-hyprland")

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(scriptPath), 0755); err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentDisplayManager,
			"Failed to create bin directory",
			err,
		), err
	}

	// Create launch script
	script := `#!/bin/bash
# Hyprland launch script
exec Hyprland
`

	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentDisplayManager,
			"Failed to create launch script",
			err,
		), err
	}

	details := []string{
		"TTY launch configured",
		fmt.Sprintf("Launch script: %s", scriptPath),
		"To start Hyprland, run: start-hyprland",
	}

	return result.
		WithDetails(details...).
		Complete(postinstall.StatusCompleted), nil
}

func (i *DisplayManagerInstaller) createHyprlandSession(path string) error {
	sessionContent := `[Desktop Entry]
Name=Hyprland
Comment=An intelligent dynamic tiling Wayland compositor
Exec=Hyprland
Type=Application
`

	return os.WriteFile(path, []byte(sessionContent), 0644)
}

// Verify checks if the display manager is properly configured
func (i *DisplayManagerInstaller) Verify(ctx context.Context) (bool, error) {
	switch i.dm {
	case postinstall.DisplayManagerSDDM:
		return i.serviceMgr.IsEnabled(ctx, "sddm")
	case postinstall.DisplayManagerGDM:
		return i.serviceMgr.IsEnabled(ctx, "gdm")
	case postinstall.DisplayManagerTTY:
		scriptPath := filepath.Join(i.homeDir, ".local/bin/start-hyprland")
		_, err := os.Stat(scriptPath)
		return err == nil, nil
	default:
		return true, nil // None or unknown, consider verified
	}
}

// Rollback reverts the installation
func (i *DisplayManagerInstaller) Rollback(ctx context.Context) error {
	switch i.dm {
	case postinstall.DisplayManagerSDDM:
		// Disable service but don't uninstall (user may want to keep it)
		return i.serviceMgr.Disable(ctx, "sddm")
	case postinstall.DisplayManagerGDM:
		return i.serviceMgr.Disable(ctx, "gdm")
	case postinstall.DisplayManagerTTY:
		scriptPath := filepath.Join(i.homeDir, ".local/bin/start-hyprland")
		return os.Remove(scriptPath)
	default:
		return nil
	}
}
