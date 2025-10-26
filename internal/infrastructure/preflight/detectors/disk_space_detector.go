package detectors

import (
	"context"
	"syscall"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
)

// SystemDiskSpaceDetector implements preflight.DiskSpaceDetector using syscall.Statfs
type SystemDiskSpaceDetector struct{}

// NewSystemDiskSpaceDetector creates a new disk space detector
func NewSystemDiskSpaceDetector() *SystemDiskSpaceDetector {
	return &SystemDiskSpaceDetector{}
}

// DetectAvailableSpace checks disk space at path
func (d *SystemDiskSpaceDetector) DetectAvailableSpace(ctx context.Context, path string) (preflight.DiskSpace, error) {
	if path == "" {
		path = "/"
	}

	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return preflight.DiskSpace{}, err
	}

	// Calculate available and total bytes
	// Bavail is available to non-root users
	// Blocks is total data blocks
	// Bsize is block size
	available := stat.Bavail * uint64(stat.Bsize)
	total := stat.Blocks * uint64(stat.Bsize)

	return preflight.NewDiskSpace(available, total, path)
}
