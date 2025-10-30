package theme

import (
	"context"
	"fmt"
	"os/exec"
)

// CommandExecutor defines the interface for executing system commands
type CommandExecutor interface {
	Execute(ctx context.Context, command string, args ...string) error
}

// ComponentReloader handles reloading of system components after theme changes
type ComponentReloader struct {
	executor CommandExecutor
	dryRun   bool
}

// NewComponentReloader creates a new component reloader with default executor
func NewComponentReloader(executor CommandExecutor) *ComponentReloader {
	return &ComponentReloader{
		executor: executor,
		dryRun:   false,
	}
}

// NewComponentReloaderWithOptions creates a new component reloader with options
func NewComponentReloaderWithOptions(executor CommandExecutor, dryRun bool) *ComponentReloader {
	return &ComponentReloader{
		executor: executor,
		dryRun:   dryRun,
	}
}

// ReloadHyprland reloads Hyprland configuration
func (r *ComponentReloader) ReloadHyprland(ctx context.Context) error {
	if r.dryRun {
		fmt.Println("[DRY-RUN] Would reload Hyprland configuration")
		return nil
	}

	if err := r.executor.Execute(ctx, "hyprctl", "reload"); err != nil {
		return fmt.Errorf("failed to reload Hyprland: %w", err)
	}

	return nil
}

// ReloadWaybar restarts Waybar
func (r *ComponentReloader) ReloadWaybar(ctx context.Context) error {
	if r.dryRun {
		fmt.Println("[DRY-RUN] Would restart Waybar")
		return nil
	}

	// Kill existing waybar process (ignore error if not running)
	_ = r.executor.Execute(ctx, "killall", "waybar")

	// Start waybar in background
	if err := r.executor.Execute(ctx, "waybar", "&"); err != nil {
		// Don't fail if waybar can't start - user might not have it installed
		fmt.Printf("Warning: failed to restart Waybar: %v\n", err)
	}

	return nil
}

// ReloadAll reloads all components in the correct order
func (r *ComponentReloader) ReloadAll(ctx context.Context) error {
	if r.dryRun {
		fmt.Println("[DRY-RUN] Would reload all components")
		return nil
	}

	// Reload Hyprland first
	if err := r.ReloadHyprland(ctx); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	// Then reload Waybar
	if err := r.ReloadWaybar(ctx); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	return nil
}

// SystemCommandExecutor executes commands using os/exec
type SystemCommandExecutor struct{}

// NewSystemCommandExecutor creates a new system command executor
func NewSystemCommandExecutor() *SystemCommandExecutor {
	return &SystemCommandExecutor{}
}

// Execute runs a system command
func (e *SystemCommandExecutor) Execute(ctx context.Context, command string, args ...string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	return cmd.Run()
}
