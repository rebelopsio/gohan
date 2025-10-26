package services

import (
	"context"
	"os"

	"github.com/rebelopsio/gohan/internal/domain/installation"
)

// ConfigurationMerger implements installation.ConfigurationMerger
// Handles merging of installation configurations
type ConfigurationMerger struct{}

// NewConfigurationMerger creates a new configuration merger
func NewConfigurationMerger() *ConfigurationMerger {
	return &ConfigurationMerger{}
}

// MergeConfigurations implements installation.ConfigurationMerger
// Merges existing and new configurations, preserving user settings
func (c *ConfigurationMerger) MergeConfigurations(
	ctx context.Context,
	existing, new installation.InstallationConfiguration,
) (installation.InstallationConfiguration, error) {
	// Build component map from existing to track what we have
	existingComponents := make(map[installation.ComponentName]installation.ComponentSelection)
	for _, comp := range existing.Components() {
		existingComponents[comp.Component()] = comp
	}

	// Start with all new components (which will have newer versions)
	mergedComponents := make([]installation.ComponentSelection, 0)
	newComponentMap := make(map[installation.ComponentName]bool)

	for _, comp := range new.Components() {
		mergedComponents = append(mergedComponents, comp)
		newComponentMap[comp.Component()] = true
	}

	// Add any existing components that weren't in new
	for _, comp := range existing.Components() {
		if !newComponentMap[comp.Component()] {
			mergedComponents = append(mergedComponents, comp)
		}
	}

	// Determine GPU support - prefer existing if new doesn't have it
	var gpuSupport *installation.GPUSupport
	if new.HasGPUSupport() {
		gpuSupport = new.GPUSupport()
	} else if existing.HasGPUSupport() {
		gpuSupport = existing.GPUSupport()
	}

	// Use new disk space if available, otherwise use existing
	diskSpace := new.DiskSpace()
	if diskSpace.Available() == 0 {
		diskSpace = existing.DiskSpace()
	}

	// Preserve merge flag from existing
	mergeExistingConfig := existing.MergeExistingConfig()

	// Create merged configuration
	merged, err := installation.NewInstallationConfiguration(
		mergedComponents,
		gpuSupport,
		diskSpace,
		mergeExistingConfig,
	)
	if err != nil {
		return installation.InstallationConfiguration{}, err
	}

	return merged, nil
}

// ShouldBackupExisting implements installation.ConfigurationMerger
// Determines if existing configuration should be backed up
func (c *ConfigurationMerger) ShouldBackupExisting(ctx context.Context, path string) (bool, error) {
	if path == "" {
		return false, nil
	}

	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// Only backup if it's a regular file, not a directory
	if info.IsDir() {
		return false, nil
	}

	// File exists and is regular file - should backup
	return true, nil
}
