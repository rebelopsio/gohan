package configuration

import "github.com/rebelopsio/gohan/internal/domain/installation"

const (
	GB = 1024 * 1024 * 1024
)

// ConfigurationManifest is a value object representing the immutable core definition
// of what should be installed
type ConfigurationManifest struct {
	components        []installation.ComponentSelection
	diskRequiredBytes uint64
	gpuRequired       bool
}

// NewConfigurationManifest creates a new configuration manifest value object
// Validates that at least one component is specified and core Hyprland is included
// Components slice is defensively copied for immutability
func NewConfigurationManifest(
	components []installation.ComponentSelection,
	diskRequiredBytes uint64,
	gpuRequired bool,
) (ConfigurationManifest, error) {
	// Must have at least one component
	if len(components) == 0 {
		return ConfigurationManifest{}, ErrNoComponents
	}

	// Must include the core Hyprland component if any Hyprland-related components
	hasCoreComponent := false
	for _, comp := range components {
		if comp.IsCore() {
			hasCoreComponent = true
			break
		}
	}

	if !hasCoreComponent {
		return ConfigurationManifest{}, ErrMissingCoreComponent
	}

	// Defensive copy of components slice
	componentsCopy := make([]installation.ComponentSelection, len(components))
	copy(componentsCopy, components)

	return ConfigurationManifest{
		components:        componentsCopy,
		diskRequiredBytes: diskRequiredBytes,
		gpuRequired:       gpuRequired,
	}, nil
}

// Components returns a defensive copy of the components slice
func (m ConfigurationManifest) Components() []installation.ComponentSelection {
	components := make([]installation.ComponentSelection, len(m.components))
	copy(components, m.components)
	return components
}

// ComponentCount returns the number of components in the manifest
func (m ConfigurationManifest) ComponentCount() int {
	return len(m.components)
}

// HasCoreComponent returns true if the manifest includes the core Hyprland component
func (m ConfigurationManifest) HasCoreComponent() bool {
	for _, comp := range m.components {
		if comp.IsCore() {
			return true
		}
	}
	return false
}

// DiskRequiredBytes returns the required disk space in bytes
func (m ConfigurationManifest) DiskRequiredBytes() uint64 {
	return m.diskRequiredBytes
}

// DiskRequiredGB returns the required disk space in gigabytes
func (m ConfigurationManifest) DiskRequiredGB() float64 {
	return float64(m.diskRequiredBytes) / float64(GB)
}

// GPURequired returns true if GPU support is required
func (m ConfigurationManifest) GPURequired() bool {
	return m.gpuRequired
}
