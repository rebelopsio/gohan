package checkers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebelopsio/gohan/internal/domain/verification"
)

// ThemeChecker verifies theme installation
type ThemeChecker struct {
	configDir string
}

// NewThemeChecker creates a new theme checker
func NewThemeChecker() *ThemeChecker {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config")
	return &ThemeChecker{configDir: configDir}
}

// Name returns the checker name
func (c *ThemeChecker) Name() string {
	return "Theme Configuration"
}

// Component returns the component being checked
func (c *ThemeChecker) Component() verification.ComponentName {
	return verification.ComponentTheme
}

// Check verifies theme is configured
func (c *ThemeChecker) Check(ctx context.Context) verification.CheckResult {
	// Check if theme state file exists
	homeDir, _ := os.UserHomeDir()
	themeStateFile := filepath.Join(homeDir, ".config/gohan/theme-state.json")

	if _, err := os.Stat(themeStateFile); os.IsNotExist(err) {
		return verification.NewCheckResult(
			verification.ComponentTheme,
			verification.StatusWarning,
			verification.SeverityMedium,
			"No theme has been applied yet",
			[]string{
				"Theme state file not found",
				fmt.Sprintf("Expected at: %s", themeStateFile),
			},
			[]string{
				"Apply a theme: gohan theme set mocha",
				"List available themes: gohan theme list",
			},
		)
	}

	// Check key theme files exist
	themeFiles := []string{
		filepath.Join(c.configDir, "hypr/hyprland.conf"),
		filepath.Join(c.configDir, "waybar/style.css"),
	}

	missingFiles := []string{}
	for _, file := range themeFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			missingFiles = append(missingFiles, file)
		}
	}

	if len(missingFiles) > 0 {
		return verification.NewCheckResult(
			verification.ComponentTheme,
			verification.StatusWarning,
			verification.SeverityMedium,
			"Some theme configuration files are missing",
			append([]string{"Missing files:"}, missingFiles...),
			[]string{
				"Reapply theme: gohan theme set <theme-name>",
				"Run installation: gohan install",
			},
		)
	}

	return verification.NewCheckResult(
		verification.ComponentTheme,
		verification.StatusPass,
		verification.SeverityLow,
		"Theme configuration is present",
		[]string{
			fmt.Sprintf("Theme state: %s", themeStateFile),
			"Configuration files found",
		},
		nil,
	)
}
