package installation

import (
	"fmt"
	"strings"
)

// ComponentSelection represents a user's choice to install a specific component
type ComponentSelection struct {
	component   ComponentName
	version     string
	packageInfo *PackageInfo
}

// NewComponentSelection creates a new component selection value object
// Version is required and will be trimmed of whitespace
// PackageInfo is optional (nil is valid)
func NewComponentSelection(component ComponentName, version string, packageInfo *PackageInfo) (ComponentSelection, error) {
	// Trim and validate version
	version = strings.TrimSpace(version)
	if version == "" {
		return ComponentSelection{}, ErrInvalidComponentSelection
	}

	return ComponentSelection{
		component:   component,
		version:     version,
		packageInfo: packageInfo,
	}, nil
}

// Component returns the component name
func (c ComponentSelection) Component() ComponentName {
	return c.component
}

// Version returns the component version
func (c ComponentSelection) Version() string {
	return c.version
}

// PackageInfo returns the package info if available
func (c ComponentSelection) PackageInfo() *PackageInfo {
	return c.packageInfo
}

// HasPackageInfo returns true if package info is available
func (c ComponentSelection) HasPackageInfo() bool {
	return c.packageInfo != nil
}

// IsCore returns true if this is the core Hyprland component
func (c ComponentSelection) IsCore() bool {
	return c.component.IsCore()
}

// IsDriver returns true if this is a GPU driver component
func (c ComponentSelection) IsDriver() bool {
	return c.component.IsDriver()
}

// EstimatedSizeBytes returns the estimated size in bytes
// Returns 0 if package info is not available
func (c ComponentSelection) EstimatedSizeBytes() uint64 {
	if c.packageInfo == nil {
		return 0
	}
	return c.packageInfo.SizeBytes()
}

// String returns human-readable representation
func (c ComponentSelection) String() string {
	if c.packageInfo != nil {
		return fmt.Sprintf("%s v%s (%s)",
			c.component.String(), c.version, c.packageInfo.String())
	}
	return fmt.Sprintf("%s v%s",
		c.component.String(), c.version)
}
