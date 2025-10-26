package installation

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SystemSnapshot is an entity representing the system state before installation
// Entities have identity and can change state over time
type SystemSnapshot struct {
	id               string
	createdAt        time.Time
	path             string
	diskSpace        DiskSpace
	packages         []string
	corrupted        bool
	corruptionReason string
}

// NewSystemSnapshot creates a new system snapshot entity
// Path is required and represents where the snapshot is stored
// Packages slice is defensively copied
func NewSystemSnapshot(path string, diskSpace DiskSpace, packages []string) (*SystemSnapshot, error) {
	// Validate path
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, ErrSnapshotInvalid
	}

	// Generate unique ID
	id := uuid.New().String()

	// Defensive copy of packages
	var pkgs []string
	if packages != nil {
		pkgs = make([]string, len(packages))
		copy(pkgs, packages)
	}

	return &SystemSnapshot{
		id:        id,
		createdAt: time.Now(),
		path:      path,
		diskSpace: diskSpace,
		packages:  pkgs,
		corrupted: false,
	}, nil
}

// ID returns the unique identifier for this snapshot
// Entities are identified by their ID, not their attributes
func (s *SystemSnapshot) ID() string {
	return s.id
}

// CreatedAt returns when the snapshot was created
func (s *SystemSnapshot) CreatedAt() time.Time {
	return s.createdAt
}

// Path returns the filesystem path where the snapshot is stored
func (s *SystemSnapshot) Path() string {
	return s.path
}

// DiskSpace returns the disk space state at snapshot time
func (s *SystemSnapshot) DiskSpace() DiskSpace {
	return s.diskSpace
}

// Packages returns a defensive copy of the package list
func (s *SystemSnapshot) Packages() []string {
	if s.packages == nil {
		return nil
	}
	pkgs := make([]string, len(s.packages))
	copy(pkgs, s.packages)
	return pkgs
}

// PackageCount returns the number of packages in the snapshot
func (s *SystemSnapshot) PackageCount() int {
	return len(s.packages)
}

// IsCorrupted returns true if the snapshot has been marked as corrupted
func (s *SystemSnapshot) IsCorrupted() bool {
	return s.corrupted
}

// CorruptionReason returns the reason why the snapshot was marked as corrupted
func (s *SystemSnapshot) CorruptionReason() string {
	return s.corruptionReason
}

// MarkAsCorrupted marks the snapshot as corrupted with a reason
// This is a state mutation - entities can change state
func (s *SystemSnapshot) MarkAsCorrupted(reason string) {
	s.corrupted = true
	s.corruptionReason = reason
}

// IsValid returns true if the snapshot is not corrupted
func (s *SystemSnapshot) IsValid() bool {
	return !s.corrupted
}

// Age returns how long ago the snapshot was created
func (s *SystemSnapshot) Age() time.Duration {
	return time.Since(s.createdAt)
}

// String returns human-readable representation
func (s *SystemSnapshot) String() string {
	status := "valid"
	if s.corrupted {
		status = fmt.Sprintf("corrupted: %s", s.corruptionReason)
	}
	return fmt.Sprintf("Snapshot %s (%d packages, %s)",
		s.path, len(s.packages), status)
}
