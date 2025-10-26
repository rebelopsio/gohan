package installation

import "fmt"

const (
	GB = 1024 * 1024 * 1024
	MB = 1024 * 1024
)

// DiskSpace represents disk space requirements for installation
type DiskSpace struct {
	available uint64
	required  uint64
}

// NewDiskSpace creates a new disk space value object
// Returns an error if available space is insufficient for requirements
func NewDiskSpace(available uint64, required uint64) (DiskSpace, error) {
	if available < required {
		return DiskSpace{}, ErrInsufficientDiskSpace
	}

	return DiskSpace{
		available: available,
		required:  required,
	}, nil
}

// Available returns available bytes
func (d DiskSpace) Available() uint64 {
	return d.available
}

// Required returns required bytes for installation
func (d DiskSpace) Required() uint64 {
	return d.required
}

// IsSufficient returns true if available space meets requirements
func (d DiskSpace) IsSufficient() bool {
	return d.available >= d.required
}

// RemainingAfterInstall returns space that will remain after installation
func (d DiskSpace) RemainingAfterInstall() uint64 {
	if d.available >= d.required {
		return d.available - d.required
	}
	return 0
}

// AvailableGB returns available space in gigabytes
func (d DiskSpace) AvailableGB() float64 {
	return float64(d.available) / float64(GB)
}

// RequiredGB returns required space in gigabytes
func (d DiskSpace) RequiredGB() float64 {
	return float64(d.required) / float64(GB)
}

// RemainingAfterInstallGB returns remaining space after installation in gigabytes
func (d DiskSpace) RemainingAfterInstallGB() float64 {
	return float64(d.RemainingAfterInstall()) / float64(GB)
}

// String returns human-readable representation
func (d DiskSpace) String() string {
	return fmt.Sprintf("%.2f GB available, %.2f GB required, %.2f GB remaining",
		d.AvailableGB(), d.RequiredGB(), d.RemainingAfterInstallGB())
}
