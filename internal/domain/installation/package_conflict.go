package installation

import (
	"fmt"
	"strings"
)

// PackageConflict is a value object representing a package conflict
// Value objects are immutable and identified by their attributes
type PackageConflict struct {
	packageName        string
	conflictingPackage string
	reason             string
}

// NewPackageConflict creates a new package conflict value object
// Both package names and reason are required
func NewPackageConflict(packageName, conflictingPackage, reason string) (PackageConflict, error) {
	packageName = strings.TrimSpace(packageName)
	conflictingPackage = strings.TrimSpace(conflictingPackage)
	reason = strings.TrimSpace(reason)

	if packageName == "" || conflictingPackage == "" {
		return PackageConflict{}, ErrInvalidConfiguration
	}

	if reason == "" {
		reason = "package conflict detected"
	}

	return PackageConflict{
		packageName:        packageName,
		conflictingPackage: conflictingPackage,
		reason:             reason,
	}, nil
}

// PackageName returns the name of the package being installed
func (p PackageConflict) PackageName() string {
	return p.packageName
}

// ConflictingPackage returns the name of the conflicting package
func (p PackageConflict) ConflictingPackage() string {
	return p.conflictingPackage
}

// Reason returns the reason for the conflict
func (p PackageConflict) Reason() string {
	return p.reason
}

// String returns human-readable representation
func (p PackageConflict) String() string {
	return fmt.Sprintf("Package conflict: %s conflicts with %s (%s)",
		p.packageName, p.conflictingPackage, p.reason)
}

// ConflictResolution is a value object representing how a conflict was resolved
type ConflictResolution struct {
	conflict PackageConflict
	action   ResolutionAction
	applied  bool
}

// NewConflictResolution creates a new conflict resolution value object
func NewConflictResolution(conflict PackageConflict, action ResolutionAction) ConflictResolution {
	return ConflictResolution{
		conflict: conflict,
		action:   action,
		applied:  false,
	}
}

// Conflict returns the package conflict being resolved
func (r ConflictResolution) Conflict() PackageConflict {
	return r.conflict
}

// Action returns the resolution action
func (r ConflictResolution) Action() ResolutionAction {
	return r.action
}

// IsApplied returns true if the resolution has been applied
func (r ConflictResolution) IsApplied() bool {
	return r.applied
}

// MarkAsApplied marks the resolution as applied
// Returns a new ConflictResolution with applied set to true
// Value objects are immutable, so we return a new instance
func (r ConflictResolution) MarkAsApplied() ConflictResolution {
	return ConflictResolution{
		conflict: r.conflict,
		action:   r.action,
		applied:  true,
	}
}

// String returns human-readable representation
func (r ConflictResolution) String() string {
	status := "pending"
	if r.applied {
		status = "applied"
	}
	return fmt.Sprintf("Resolution for %s: %s (%s)",
		r.conflict.PackageName(), r.action.String(), status)
}
