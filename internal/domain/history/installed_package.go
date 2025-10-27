package history

import (
	"fmt"
	"strings"
)

// InstalledPackage is a value object representing a package that was installed
// This is a simplified snapshot for historical purposes (not the full entity)
type InstalledPackage struct {
	name      string
	version   string
	sizeBytes uint64
}

// NewInstalledPackage creates an installed package value object
func NewInstalledPackage(name, version string, sizeBytes uint64) (InstalledPackage, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return InstalledPackage{}, ErrInvalidPackageName
	}

	version = strings.TrimSpace(version)
	if version == "" {
		return InstalledPackage{}, ErrInvalidPackageVersion
	}

	return InstalledPackage{
		name:      name,
		version:   version,
		sizeBytes: sizeBytes,
	}, nil
}

// Name returns the package name
func (p InstalledPackage) Name() string {
	return p.name
}

// Version returns the package version
func (p InstalledPackage) Version() string {
	return p.version
}

// SizeBytes returns the package size in bytes
func (p InstalledPackage) SizeBytes() uint64 {
	return p.sizeBytes
}

// SizeMB returns the package size in megabytes
func (p InstalledPackage) SizeMB() float64 {
	return float64(p.sizeBytes) / (1024 * 1024)
}

// Equals checks if two packages are equal (by name and version)
func (p InstalledPackage) Equals(other InstalledPackage) bool {
	return p.name == other.name && p.version == other.version
}

// String returns human-readable representation
func (p InstalledPackage) String() string {
	return fmt.Sprintf("%s v%s (%.2f MB)", p.name, p.version, p.SizeMB())
}
