package preflight

import (
	"fmt"
	"strings"
)

const (
	GB = 1024 * 1024 * 1024
	MB = 1024 * 1024
)

// DiskSpace represents available disk space
type DiskSpace struct {
	available uint64
	total     uint64
	path      string
}

// NewDiskSpace creates a new disk space value object
func NewDiskSpace(available uint64, total uint64, path string) (DiskSpace, error) {
	if available > total {
		return DiskSpace{}, ErrInvalidDiskSpace
	}

	// Trim whitespace and default to root if empty
	path = strings.TrimSpace(path)
	if path == "" {
		path = "/"
	}

	return DiskSpace{
		available: available,
		total:     total,
		path:      path,
	}, nil
}

// Available returns available bytes
func (d DiskSpace) Available() uint64 {
	return d.available
}

// Total returns total bytes
func (d DiskSpace) Total() uint64 {
	return d.total
}

// Path returns the filesystem path
func (d DiskSpace) Path() string {
	return d.path
}

// MeetsMinimum checks if available space meets requirement
func (d DiskSpace) MeetsMinimum(requiredGB uint64) bool {
	required := requiredGB * GB
	return d.available >= required
}

// AvailableGB returns available space in gigabytes
func (d DiskSpace) AvailableGB() float64 {
	return float64(d.available) / float64(GB)
}

// TotalGB returns total space in gigabytes
func (d DiskSpace) TotalGB() float64 {
	return float64(d.total) / float64(GB)
}

// UsagePercent returns disk usage as percentage
func (d DiskSpace) UsagePercent() float64 {
	if d.total == 0 {
		return 0
	}
	used := d.total - d.available
	return (float64(used) / float64(d.total)) * 100
}

// String returns human-readable representation
func (d DiskSpace) String() string {
	return fmt.Sprintf("%.2f GB available / %.2f GB total (%.1f%% used)",
		d.AvailableGB(), d.TotalGB(), d.UsagePercent())
}
