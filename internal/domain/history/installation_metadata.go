package history

import (
	"strings"
	"time"
)

// InstallationMetadata is a value object capturing what was installed, when, and for how long
type InstallationMetadata struct {
	packageName    string
	targetVersion  string
	installedAt    time.Time
	completedAt    time.Time
	installedPkgs  []InstalledPackage
}

// NewInstallationMetadata creates installation metadata value object
func NewInstallationMetadata(
	packageName string,
	targetVersion string,
	installedAt time.Time,
	completedAt time.Time,
	packages []InstalledPackage,
) (InstallationMetadata, error) {
	// Trim and validate package name
	packageName = strings.TrimSpace(packageName)
	if packageName == "" {
		return InstallationMetadata{}, ErrInvalidPackageName
	}

	// Trim and validate target version
	targetVersion = strings.TrimSpace(targetVersion)
	if targetVersion == "" {
		return InstallationMetadata{}, ErrInvalidVersion
	}

	// Validate timestamps
	if installedAt.IsZero() {
		return InstallationMetadata{}, ErrInvalidTimestamp
	}
	if completedAt.IsZero() {
		return InstallationMetadata{}, ErrInvalidTimestamp
	}

	// Validate time ordering
	if completedAt.Before(installedAt) {
		return InstallationMetadata{}, ErrInvalidTimeRange
	}

	// Defensive copy of packages slice
	var pkgsCopy []InstalledPackage
	if packages != nil {
		pkgsCopy = make([]InstalledPackage, len(packages))
		copy(pkgsCopy, packages)
	} else {
		pkgsCopy = []InstalledPackage{}
	}

	return InstallationMetadata{
		packageName:    packageName,
		targetVersion:  targetVersion,
		installedAt:    installedAt,
		completedAt:    completedAt,
		installedPkgs:  pkgsCopy,
	}, nil
}

// PackageName returns the target package name
func (m InstallationMetadata) PackageName() string {
	return m.packageName
}

// TargetVersion returns the target version
func (m InstallationMetadata) TargetVersion() string {
	return m.targetVersion
}

// InstalledAt returns when installation started
func (m InstallationMetadata) InstalledAt() time.Time {
	return m.installedAt
}

// CompletedAt returns when installation completed
func (m InstallationMetadata) CompletedAt() time.Time {
	return m.completedAt
}

// InstalledPackages returns a defensive copy of installed packages
func (m InstallationMetadata) InstalledPackages() []InstalledPackage {
	pkgs := make([]InstalledPackage, len(m.installedPkgs))
	copy(pkgs, m.installedPkgs)
	return pkgs
}

// DurationMs returns installation duration in milliseconds
func (m InstallationMetadata) DurationMs() int64 {
	duration := m.completedAt.Sub(m.installedAt)
	return duration.Milliseconds()
}

// DurationSeconds returns installation duration in seconds
func (m InstallationMetadata) DurationSeconds() float64 {
	duration := m.completedAt.Sub(m.installedAt)
	return duration.Seconds()
}

// PackageCount returns the number of installed packages
func (m InstallationMetadata) PackageCount() int {
	return len(m.installedPkgs)
}

// TotalSizeBytes returns total size of all installed packages in bytes
func (m InstallationMetadata) TotalSizeBytes() uint64 {
	var total uint64
	for _, pkg := range m.installedPkgs {
		total += pkg.SizeBytes()
	}
	return total
}

// TotalSizeMB returns total size of all installed packages in MB
func (m InstallationMetadata) TotalSizeMB() float64 {
	return float64(m.TotalSizeBytes()) / (1024 * 1024)
}

// HasPackage returns true if a package with the given name was installed
func (m InstallationMetadata) HasPackage(name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}

	for _, pkg := range m.installedPkgs {
		if pkg.Name() == name {
			return true
		}
	}
	return false
}
