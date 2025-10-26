package installation_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSystemSnapshot(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		diskSpace    installation.DiskSpace
		packages     []string
		wantErr      bool
		errType      error
	}{
		{
			name:      "valid snapshot with packages",
			path:      "/var/backup/snapshot-001",
			diskSpace: mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			packages:  []string{"package1", "package2", "package3"},
			wantErr:   false,
		},
		{
			name:      "valid snapshot without packages",
			path:      "/var/backup/snapshot-002",
			diskSpace: mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			packages:  nil,
			wantErr:   false,
		},
		{
			name:      "empty path",
			path:      "",
			diskSpace: mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			packages:  []string{"package1"},
			wantErr:   true,
			errType:   installation.ErrSnapshotInvalid,
		},
		{
			name:      "whitespace path",
			path:      "   ",
			diskSpace: mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			packages:  []string{"package1"},
			wantErr:   true,
			errType:   installation.ErrSnapshotInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot, err := installation.NewSystemSnapshot(
				tt.path,
				tt.diskSpace,
				tt.packages,
			)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, snapshot.ID(), "Snapshot should have a unique ID")
			assert.Equal(t, tt.path, snapshot.Path())
			assert.False(t, snapshot.CreatedAt().IsZero(), "CreatedAt should be set")
			assert.Equal(t, len(tt.packages), snapshot.PackageCount())
		})
	}
}

func TestSystemSnapshot_Identity(t *testing.T) {
	t.Run("each snapshot has unique ID", func(t *testing.T) {
		snapshot1, err := installation.NewSystemSnapshot(
			"/path1",
			mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			nil,
		)
		require.NoError(t, err)

		snapshot2, err := installation.NewSystemSnapshot(
			"/path2",
			mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			nil,
		)
		require.NoError(t, err)

		assert.NotEqual(t, snapshot1.ID(), snapshot2.ID(),
			"Different snapshots should have different IDs")
	})

	t.Run("ID is not empty", func(t *testing.T) {
		snapshot, err := installation.NewSystemSnapshot(
			"/path",
			mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			nil,
		)
		require.NoError(t, err)

		assert.NotEmpty(t, snapshot.ID())
	})
}

func TestSystemSnapshot_CreatedAt(t *testing.T) {
	before := time.Now()

	snapshot, err := installation.NewSystemSnapshot(
		"/path",
		mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
		nil,
	)
	require.NoError(t, err)

	after := time.Now()

	assert.False(t, snapshot.CreatedAt().Before(before),
		"CreatedAt should be at or after creation start")
	assert.False(t, snapshot.CreatedAt().After(after),
		"CreatedAt should be at or before creation end")
}

func TestSystemSnapshot_PackageCount(t *testing.T) {
	tests := []struct {
		name      string
		packages  []string
		wantCount int
	}{
		{
			name:      "multiple packages",
			packages:  []string{"pkg1", "pkg2", "pkg3"},
			wantCount: 3,
		},
		{
			name:      "single package",
			packages:  []string{"pkg1"},
			wantCount: 1,
		},
		{
			name:      "no packages",
			packages:  nil,
			wantCount: 0,
		},
		{
			name:      "empty slice",
			packages:  []string{},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot, err := installation.NewSystemSnapshot(
				"/path",
				mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
				tt.packages,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCount, snapshot.PackageCount())
		})
	}
}

func TestSystemSnapshot_Packages(t *testing.T) {
	t.Run("returns defensive copy", func(t *testing.T) {
		original := []string{"pkg1", "pkg2", "pkg3"}

		snapshot, err := installation.NewSystemSnapshot(
			"/path",
			mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			original,
		)
		require.NoError(t, err)

		// Modify original
		original[0] = "modified"

		// Snapshot should not be affected
		packages := snapshot.Packages()
		assert.Equal(t, "pkg1", packages[0], "Snapshot should not be affected by original modification")
	})

	t.Run("returned slice modification doesn't affect snapshot", func(t *testing.T) {
		snapshot, err := installation.NewSystemSnapshot(
			"/path",
			mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			[]string{"pkg1", "pkg2"},
		)
		require.NoError(t, err)

		// Get and modify returned slice
		packages := snapshot.Packages()
		packages[0] = "modified"

		// Get again - should be unchanged
		packagesAgain := snapshot.Packages()
		assert.Equal(t, "pkg1", packagesAgain[0], "Snapshot should not be affected by returned slice modification")
	})
}

func TestSystemSnapshot_MarkAsCorrupted(t *testing.T) {
	snapshot, err := installation.NewSystemSnapshot(
		"/path",
		mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
		[]string{"pkg1"},
	)
	require.NoError(t, err)

	// Initially not corrupted
	assert.False(t, snapshot.IsCorrupted(), "New snapshot should not be corrupted")

	// Mark as corrupted with reason
	reason := "checksum mismatch"
	snapshot.MarkAsCorrupted(reason)

	// Now should be corrupted
	assert.True(t, snapshot.IsCorrupted(), "Snapshot should be marked as corrupted")
	assert.Equal(t, reason, snapshot.CorruptionReason())
}

func TestSystemSnapshot_IsValid(t *testing.T) {
	t.Run("valid snapshot", func(t *testing.T) {
		snapshot, err := installation.NewSystemSnapshot(
			"/path",
			mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			[]string{"pkg1"},
		)
		require.NoError(t, err)

		assert.True(t, snapshot.IsValid(), "Uncorrupted snapshot should be valid")
	})

	t.Run("corrupted snapshot is not valid", func(t *testing.T) {
		snapshot, err := installation.NewSystemSnapshot(
			"/path",
			mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			[]string{"pkg1"},
		)
		require.NoError(t, err)

		snapshot.MarkAsCorrupted("test reason")

		assert.False(t, snapshot.IsValid(), "Corrupted snapshot should not be valid")
	})
}

func TestSystemSnapshot_Age(t *testing.T) {
	snapshot, err := installation.NewSystemSnapshot(
		"/path",
		mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
		nil,
	)
	require.NoError(t, err)

	// Wait a tiny bit
	time.Sleep(10 * time.Millisecond)

	age := snapshot.Age()
	assert.Greater(t, age, time.Duration(0), "Age should be positive")
	assert.Less(t, age, 1*time.Second, "Age should be less than 1 second in test")
}

func TestSystemSnapshot_String(t *testing.T) {
	snapshot, err := installation.NewSystemSnapshot(
		"/var/backup/snapshot-001",
		mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
		[]string{"pkg1", "pkg2", "pkg3"},
	)
	require.NoError(t, err)

	str := snapshot.String()
	assert.Contains(t, str, "snapshot-001")
	assert.Contains(t, str, "3 packages")
}

func TestSystemSnapshot_EntityBehavior(t *testing.T) {
	t.Run("entities with same ID are equal", func(t *testing.T) {
		snapshot1, err := installation.NewSystemSnapshot(
			"/path",
			mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			nil,
		)
		require.NoError(t, err)

		// Simulate persistence and retrieval (same ID)
		snapshot2 := snapshot1

		assert.Equal(t, snapshot1.ID(), snapshot2.ID(),
			"Entities with same ID should be considered equal")
	})

	t.Run("entities with different IDs are not equal", func(t *testing.T) {
		snapshot1, err := installation.NewSystemSnapshot(
			"/path1",
			mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			nil,
		)
		require.NoError(t, err)

		snapshot2, err := installation.NewSystemSnapshot(
			"/path2",
			mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			nil,
		)
		require.NoError(t, err)

		assert.NotEqual(t, snapshot1.ID(), snapshot2.ID(),
			"Different entities should have different IDs")
	})

	t.Run("entity state can change (mutability)", func(t *testing.T) {
		snapshot, err := installation.NewSystemSnapshot(
			"/path",
			mustCreateDiskSpace(t, 100*installation.GB, 10*installation.GB),
			nil,
		)
		require.NoError(t, err)

		// Entity state changes
		id := snapshot.ID()
		assert.False(t, snapshot.IsCorrupted())

		snapshot.MarkAsCorrupted("test")

		// ID stays same but state changed
		assert.Equal(t, id, snapshot.ID())
		assert.True(t, snapshot.IsCorrupted())
	})
}
