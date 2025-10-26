package installation

import "fmt"

// InstallationConfiguration represents the complete configuration for an installation
type InstallationConfiguration struct {
	components        []ComponentSelection
	gpuSupport        *GPUSupport
	diskSpace         DiskSpace
	mergeExistingConf bool
}

// NewInstallationConfiguration creates a new installation configuration value object
// Validates that:
// - At least one component is selected
// - The core Hyprland component is included
// - Components slice is defensively copied for immutability
func NewInstallationConfiguration(
	components []ComponentSelection,
	gpuSupport *GPUSupport,
	diskSpace DiskSpace,
	mergeExistingConfig bool,
) (InstallationConfiguration, error) {
	// Must have at least one component
	if len(components) == 0 {
		return InstallationConfiguration{}, ErrInvalidConfiguration
	}

	// Must include the core Hyprland component
	hasCoreComponent := false
	for _, comp := range components {
		if comp.IsCore() {
			hasCoreComponent = true
			break
		}
	}

	if !hasCoreComponent {
		return InstallationConfiguration{}, ErrInvalidConfiguration
	}

	// Defensive copy of components slice
	componentsCopy := make([]ComponentSelection, len(components))
	copy(componentsCopy, components)

	return InstallationConfiguration{
		components:        componentsCopy,
		gpuSupport:        gpuSupport,
		diskSpace:         diskSpace,
		mergeExistingConf: mergeExistingConfig,
	}, nil
}

// Components returns a defensive copy of the components slice
func (c InstallationConfiguration) Components() []ComponentSelection {
	components := make([]ComponentSelection, len(c.components))
	copy(components, c.components)
	return components
}

// ComponentCount returns the number of components to install
func (c InstallationConfiguration) ComponentCount() int {
	return len(c.components)
}

// HasCoreComponent returns true if Hyprland core is included
func (c InstallationConfiguration) HasCoreComponent() bool {
	for _, comp := range c.components {
		if comp.IsCore() {
			return true
		}
	}
	return false
}

// GPUSupport returns the GPU support configuration if available
func (c InstallationConfiguration) GPUSupport() *GPUSupport {
	return c.gpuSupport
}

// HasGPUSupport returns true if GPU support is configured
func (c InstallationConfiguration) HasGPUSupport() bool {
	return c.gpuSupport != nil
}

// DiskSpace returns the disk space configuration
func (c InstallationConfiguration) DiskSpace() DiskSpace {
	return c.diskSpace
}

// MergeExistingConfig returns true if existing configuration should be merged
func (c InstallationConfiguration) MergeExistingConfig() bool {
	return c.mergeExistingConf
}

// TotalEstimatedSizeBytes returns the sum of all component sizes
// Returns 0 if components don't have package info
func (c InstallationConfiguration) TotalEstimatedSizeBytes() uint64 {
	var total uint64
	for _, comp := range c.components {
		total += comp.EstimatedSizeBytes()
	}
	return total
}

// TotalEstimatedSizeMB returns the total estimated size in megabytes
func (c InstallationConfiguration) TotalEstimatedSizeMB() float64 {
	return float64(c.TotalEstimatedSizeBytes()) / float64(MB)
}

// String returns human-readable representation
func (c InstallationConfiguration) String() string {
	gpuInfo := "no GPU config"
	if c.gpuSupport != nil {
		gpuInfo = fmt.Sprintf("GPU: %s", c.gpuSupport.Vendor())
	}

	mergeInfo := ""
	if c.mergeExistingConf {
		mergeInfo = ", merge existing"
	}

	return fmt.Sprintf("Installation: %d components, %s%s",
		len(c.components), gpuInfo, mergeInfo)
}
