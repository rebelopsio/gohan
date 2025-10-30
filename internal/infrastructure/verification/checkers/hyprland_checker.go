package checkers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rebelopsio/gohan/internal/domain/verification"
)

// HyprlandChecker verifies Hyprland installation
type HyprlandChecker struct{}

// NewHyprlandChecker creates a new Hyprland checker
func NewHyprlandChecker() *HyprlandChecker {
	return &HyprlandChecker{}
}

// Name returns the checker name
func (c *HyprlandChecker) Name() string {
	return "Hyprland Binary"
}

// Component returns the component being checked
func (c *HyprlandChecker) Component() verification.ComponentName {
	return verification.ComponentHyprland
}

// Check verifies Hyprland installation
func (c *HyprlandChecker) Check(ctx context.Context) verification.CheckResult {
	// Check if Hyprland binary exists in PATH
	path, err := exec.LookPath("Hyprland")
	if err != nil {
		return verification.NewCheckResult(
			verification.ComponentHyprland,
			verification.StatusFail,
			verification.SeverityCritical,
			"Hyprland binary not found in PATH",
			[]string{
				"Hyprland executable could not be located",
				"This indicates Hyprland is not installed or not in PATH",
			},
			[]string{
				"Install Hyprland: sudo apt install hyprland",
				"Or run: gohan install hyprland",
				"Ensure /usr/bin is in your PATH",
			},
		)
	}

	// Check if binary is executable
	info, err := os.Stat(path)
	if err != nil {
		return verification.NewCheckResult(
			verification.ComponentHyprland,
			verification.StatusFail,
			verification.SeverityCritical,
			"Cannot access Hyprland binary",
			[]string{
				fmt.Sprintf("Found at: %s", path),
				fmt.Sprintf("Error: %v", err),
			},
			[]string{
				"Check file permissions",
				"Reinstall Hyprland if corrupted",
			},
		)
	}

	if info.Mode()&0111 == 0 {
		return verification.NewCheckResult(
			verification.ComponentHyprland,
			verification.StatusFail,
			verification.SeverityHigh,
			"Hyprland binary is not executable",
			[]string{
				fmt.Sprintf("Found at: %s", path),
				fmt.Sprintf("Permissions: %v", info.Mode()),
			},
			[]string{
				fmt.Sprintf("Fix permissions: chmod +x %s", path),
			},
		)
	}

	// Try to get version
	cmd := exec.CommandContext(ctx, "Hyprland", "--version")
	output, err := cmd.Output()
	version := "unknown"
	if err == nil {
		version = strings.TrimSpace(string(output))
	}

	return verification.NewCheckResult(
		verification.ComponentHyprland,
		verification.StatusPass,
		verification.SeverityLow,
		"Hyprland is installed and executable",
		[]string{
			fmt.Sprintf("Location: %s", path),
			fmt.Sprintf("Version: %s", version),
		},
		nil,
	)
}
