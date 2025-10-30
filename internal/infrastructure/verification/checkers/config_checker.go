package checkers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebelopsio/gohan/internal/domain/verification"
)

// ConfigChecker verifies configuration files
type ConfigChecker struct {
	configDir string
}

// NewConfigChecker creates a new config checker
func NewConfigChecker() *ConfigChecker {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config")
	return &ConfigChecker{configDir: configDir}
}

// Name returns the checker name
func (c *ConfigChecker) Name() string {
	return "Configuration Files"
}

// Component returns the component being checked
func (c *ConfigChecker) Component() verification.ComponentName {
	return verification.ComponentConfiguration
}

// Check verifies essential configuration files exist
func (c *ConfigChecker) Check(ctx context.Context) verification.CheckResult {
	essentialFiles := []string{
		filepath.Join(c.configDir, "hypr/hyprland.conf"),
	}

	missingFiles := []string{}
	readableFiles := 0

	for _, file := range essentialFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			missingFiles = append(missingFiles, filepath.Base(file))
		} else if err == nil {
			// File exists and is readable
			readableFiles++
		}
	}

	if len(missingFiles) > 0 {
		return verification.NewCheckResult(
			verification.ComponentConfiguration,
			verification.StatusFail,
			verification.SeverityCritical,
			"Essential configuration files are missing",
			append([]string{"Missing files:"}, missingFiles...),
			[]string{
				"Run installation: gohan install",
				"Or apply theme: gohan theme set mocha",
			},
		)
	}

	return verification.NewCheckResult(
		verification.ComponentConfiguration,
		verification.StatusPass,
		verification.SeverityLow,
		"Configuration files are present",
		[]string{
			fmt.Sprintf("Checked %d essential files", len(essentialFiles)),
			fmt.Sprintf("All %d files are readable", readableFiles),
		},
		nil,
	)
}
