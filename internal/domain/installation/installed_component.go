package installation

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// InstalledComponent is an entity representing a successfully installed component
// Entities have identity and can change state over time
type InstalledComponent struct {
	id          string
	component   ComponentName
	version     string
	packageInfo *PackageInfo
	installedAt time.Time
	verified    bool
	verifiedAt  time.Time
}

// NewInstalledComponent creates a new installed component entity
// Version is required and will be trimmed of whitespace
// PackageInfo is optional
func NewInstalledComponent(component ComponentName, version string, packageInfo *PackageInfo) (*InstalledComponent, error) {
	// Validate version
	version = strings.TrimSpace(version)
	if version == "" {
		return nil, ErrInvalidComponentSelection
	}

	// Generate unique ID
	id := uuid.New().String()

	return &InstalledComponent{
		id:          id,
		component:   component,
		version:     version,
		packageInfo: packageInfo,
		installedAt: time.Now(),
		verified:    false,
	}, nil
}

// ID returns the unique identifier for this installed component
// Entities are identified by their ID, not their attributes
func (c *InstalledComponent) ID() string {
	return c.id
}

// Component returns the component name
func (c *InstalledComponent) Component() ComponentName {
	return c.component
}

// Version returns the installed version
func (c *InstalledComponent) Version() string {
	return c.version
}

// PackageInfo returns the package info if available
func (c *InstalledComponent) PackageInfo() *PackageInfo {
	return c.packageInfo
}

// HasPackageInfo returns true if package info is available
func (c *InstalledComponent) HasPackageInfo() bool {
	return c.packageInfo != nil
}

// InstalledAt returns when the component was installed
func (c *InstalledComponent) InstalledAt() time.Time {
	return c.installedAt
}

// IsCore returns true if this is the core Hyprland component
func (c *InstalledComponent) IsCore() bool {
	return c.component.IsCore()
}

// IsDriver returns true if this is a GPU driver component
func (c *InstalledComponent) IsDriver() bool {
	return c.component.IsDriver()
}

// IsVerified returns true if the installation has been verified
func (c *InstalledComponent) IsVerified() bool {
	return c.verified
}

// VerifiedAt returns when the component was verified
// Returns zero time if not yet verified
func (c *InstalledComponent) VerifiedAt() time.Time {
	return c.verifiedAt
}

// MarkAsVerified marks the component as verified
// This is a state mutation - entities can change state
func (c *InstalledComponent) MarkAsVerified() {
	c.verified = true
	c.verifiedAt = time.Now()
}

// Age returns how long ago the component was installed
func (c *InstalledComponent) Age() time.Duration {
	return time.Since(c.installedAt)
}

// String returns human-readable representation
func (c *InstalledComponent) String() string {
	status := "installed"
	if c.verified {
		status = "verified"
	}

	if c.packageInfo != nil {
		return fmt.Sprintf("%s v%s (%s, %s)",
			c.component.String(), c.version, status, c.packageInfo.String())
	}
	return fmt.Sprintf("%s v%s (%s)",
		c.component.String(), c.version, status)
}
