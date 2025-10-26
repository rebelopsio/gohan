package installation

import (
	"fmt"
	"strings"
)

// GPUSupport represents GPU-specific configuration for installation
type GPUSupport struct {
	vendor          string
	requiresDriver  bool
	driverComponent ComponentName
}

// NewGPUSupport creates a new GPU support value object
// Vendor is required and will be normalized (trimmed and lowercased)
// If requiresDriver is true, driverComponent must be specified and match the vendor
func NewGPUSupport(vendor string, requiresDriver bool, driverComponent ComponentName) (GPUSupport, error) {
	// Normalize vendor
	vendor = strings.ToLower(strings.TrimSpace(vendor))
	if vendor == "" {
		return GPUSupport{}, ErrInvalidGPUSupport
	}

	// Validate driver requirements
	if requiresDriver {
		// Driver component must be specified
		if driverComponent == "" {
			return GPUSupport{}, ErrInvalidGPUSupport
		}

		// Driver component must actually be a driver
		if !driverComponent.IsDriver() {
			return GPUSupport{}, ErrInvalidGPUSupport
		}

		// Driver must match vendor
		if !driverMatchesVendor(vendor, driverComponent) {
			return GPUSupport{}, ErrInvalidGPUSupport
		}
	}

	return GPUSupport{
		vendor:          vendor,
		requiresDriver:  requiresDriver,
		driverComponent: driverComponent,
	}, nil
}

// driverMatchesVendor checks if driver component matches GPU vendor
func driverMatchesVendor(vendor string, driver ComponentName) bool {
	switch vendor {
	case "amd":
		return driver == ComponentAMDDriver
	case "nvidia":
		return driver == ComponentNVIDIADriver
	case "intel":
		return driver == ComponentIntelDriver
	default:
		return false
	}
}

// Vendor returns the GPU vendor (normalized)
func (g GPUSupport) Vendor() string {
	return g.vendor
}

// RequiresDriver returns true if driver installation is required
func (g GPUSupport) RequiresDriver() bool {
	return g.requiresDriver
}

// DriverComponent returns the required driver component
// Returns empty ComponentName if no driver is required
func (g GPUSupport) DriverComponent() ComponentName {
	return g.driverComponent
}

// IsAMD returns true for AMD GPUs
func (g GPUSupport) IsAMD() bool {
	return g.vendor == "amd"
}

// IsNVIDIA returns true for NVIDIA GPUs
func (g GPUSupport) IsNVIDIA() bool {
	return g.vendor == "nvidia"
}

// IsIntel returns true for Intel GPUs
func (g GPUSupport) IsIntel() bool {
	return g.vendor == "intel"
}

// RequiresProprietary returns true if proprietary drivers are needed
// Currently only NVIDIA requires proprietary drivers
func (g GPUSupport) RequiresProprietary() bool {
	return g.IsNVIDIA()
}

// String returns human-readable representation
func (g GPUSupport) String() string {
	if g.requiresDriver {
		return fmt.Sprintf("%s GPU (driver: %s)",
			g.vendor, g.driverComponent.String())
	}
	return fmt.Sprintf("%s GPU (no driver required)", g.vendor)
}
