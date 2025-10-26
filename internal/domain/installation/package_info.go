package installation

import (
	"fmt"
	"strings"
)

// PackageInfo represents metadata about a package to be installed
type PackageInfo struct {
	name         string
	version      string
	sizeBytes    uint64
	dependencies []string
}

// NewPackageInfo creates a new package info value object
// Package name and version are required and will be trimmed of whitespace
// Dependencies slice is copied to ensure immutability
func NewPackageInfo(name, version string, sizeBytes uint64, dependencies []string) (PackageInfo, error) {
	// Trim and validate name
	name = strings.TrimSpace(name)
	if name == "" {
		return PackageInfo{}, ErrInvalidPackageInfo
	}

	// Trim and validate version
	version = strings.TrimSpace(version)
	if version == "" {
		return PackageInfo{}, ErrInvalidPackageInfo
	}

	// Copy dependencies to ensure immutability
	var deps []string
	if dependencies != nil {
		deps = make([]string, len(dependencies))
		copy(deps, dependencies)
	}

	return PackageInfo{
		name:         name,
		version:      version,
		sizeBytes:    sizeBytes,
		dependencies: deps,
	}, nil
}

// Name returns the package name
func (p PackageInfo) Name() string {
	return p.name
}

// Version returns the package version
func (p PackageInfo) Version() string {
	return p.version
}

// SizeBytes returns the package size in bytes
func (p PackageInfo) SizeBytes() uint64 {
	return p.sizeBytes
}

// SizeMB returns the package size in megabytes
func (p PackageInfo) SizeMB() float64 {
	return float64(p.sizeBytes) / float64(MB)
}

// Dependencies returns a copy of the dependencies slice
func (p PackageInfo) Dependencies() []string {
	if p.dependencies == nil {
		return nil
	}
	// Return a copy to preserve immutability
	deps := make([]string, len(p.dependencies))
	copy(deps, p.dependencies)
	return deps
}

// HasDependencies returns true if the package has dependencies
func (p PackageInfo) HasDependencies() bool {
	return len(p.dependencies) > 0
}

// DependencyCount returns the number of dependencies
func (p PackageInfo) DependencyCount() int {
	return len(p.dependencies)
}

// String returns human-readable representation
func (p PackageInfo) String() string {
	if p.HasDependencies() {
		return fmt.Sprintf("%s v%s (%.2f MB, %d dependencies)",
			p.name, p.version, p.SizeMB(), p.DependencyCount())
	}
	return fmt.Sprintf("%s v%s (%.2f MB)",
		p.name, p.version, p.SizeMB())
}
