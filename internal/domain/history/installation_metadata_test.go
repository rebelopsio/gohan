package history_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstallationMetadata(t *testing.T) {
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := time.Date(2025, 10, 26, 14, 32, 30, 0, time.UTC)

	pkg1, _ := history.NewInstalledPackage("hyprland", "0.45.0", 15728640)
	pkg2, _ := history.NewInstalledPackage("waybar", "0.10.4", 5242880)
	packages := []history.InstalledPackage{pkg1, pkg2}

	tests := []struct {
		name          string
		packageName   string
		targetVersion string
		installedAt   time.Time
		completedAt   time.Time
		packages      []history.InstalledPackage
		wantErr       error
	}{
		{
			name:          "valid metadata with all fields",
			packageName:   "hyprland",
			targetVersion: "0.45.0",
			installedAt:   installedAt,
			completedAt:   completedAt,
			packages:      packages,
			wantErr:       nil,
		},
		{
			name:          "valid metadata with empty packages list",
			packageName:   "kitty",
			targetVersion: "0.36.0",
			installedAt:   installedAt,
			completedAt:   completedAt,
			packages:      []history.InstalledPackage{},
			wantErr:       nil,
		},
		{
			name:          "valid metadata with nil packages list",
			packageName:   "kitty",
			targetVersion: "0.36.0",
			installedAt:   installedAt,
			completedAt:   completedAt,
			packages:      nil,
			wantErr:       nil,
		},
		{
			name:          "with whitespace trimmed",
			packageName:   "  hyprland  ",
			targetVersion: "  0.45.0  ",
			installedAt:   installedAt,
			completedAt:   completedAt,
			packages:      packages,
			wantErr:       nil,
		},
		{
			name:          "empty package name",
			packageName:   "",
			targetVersion: "0.45.0",
			installedAt:   installedAt,
			completedAt:   completedAt,
			packages:      packages,
			wantErr:       history.ErrInvalidPackageName,
		},
		{
			name:          "whitespace-only package name",
			packageName:   "   ",
			targetVersion: "0.45.0",
			installedAt:   installedAt,
			completedAt:   completedAt,
			packages:      packages,
			wantErr:       history.ErrInvalidPackageName,
		},
		{
			name:          "empty target version",
			packageName:   "hyprland",
			targetVersion: "",
			installedAt:   installedAt,
			completedAt:   completedAt,
			packages:      packages,
			wantErr:       history.ErrInvalidVersion,
		},
		{
			name:          "zero installedAt timestamp",
			packageName:   "hyprland",
			targetVersion: "0.45.0",
			installedAt:   time.Time{},
			completedAt:   completedAt,
			packages:      packages,
			wantErr:       history.ErrInvalidTimestamp,
		},
		{
			name:          "zero completedAt timestamp",
			packageName:   "hyprland",
			targetVersion: "0.45.0",
			installedAt:   installedAt,
			completedAt:   time.Time{},
			packages:      packages,
			wantErr:       history.ErrInvalidTimestamp,
		},
		{
			name:          "completedAt before installedAt",
			packageName:   "hyprland",
			targetVersion: "0.45.0",
			installedAt:   completedAt,
			completedAt:   installedAt,
			packages:      packages,
			wantErr:       history.ErrInvalidTimeRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := history.NewInstallationMetadata(
				tt.packageName,
				tt.targetVersion,
				tt.installedAt,
				tt.completedAt,
				tt.packages,
			)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, metadata.PackageName())
				assert.NotEmpty(t, metadata.TargetVersion())
			}
		})
	}
}

func TestInstallationMetadata_Accessors(t *testing.T) {
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := time.Date(2025, 10, 26, 14, 32, 30, 0, time.UTC)

	pkg1, _ := history.NewInstalledPackage("hyprland", "0.45.0", 15728640)
	pkg2, _ := history.NewInstalledPackage("waybar", "0.10.4", 5242880)
	packages := []history.InstalledPackage{pkg1, pkg2}

	metadata, err := history.NewInstallationMetadata(
		"hyprland-meta",
		"0.45.0",
		installedAt,
		completedAt,
		packages,
	)
	require.NoError(t, err)

	assert.Equal(t, "hyprland-meta", metadata.PackageName())
	assert.Equal(t, "0.45.0", metadata.TargetVersion())
	assert.Equal(t, installedAt, metadata.InstalledAt())
	assert.Equal(t, completedAt, metadata.CompletedAt())
	assert.Len(t, metadata.InstalledPackages(), 2)
}

func TestInstallationMetadata_DurationMs(t *testing.T) {
	tests := []struct {
		name        string
		installedAt time.Time
		completedAt time.Time
		expectedMs  int64
	}{
		{
			name:        "2.5 seconds duration",
			installedAt: time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC),
			completedAt: time.Date(2025, 10, 26, 14, 30, 2, 500000000, time.UTC),
			expectedMs:  2500,
		},
		{
			name:        "1 minute duration",
			installedAt: time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC),
			completedAt: time.Date(2025, 10, 26, 14, 31, 0, 0, time.UTC),
			expectedMs:  60000,
		},
		{
			name:        "zero duration (same time)",
			installedAt: time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC),
			completedAt: time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC),
			expectedMs:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := history.NewInstallationMetadata(
				"hyprland",
				"0.45.0",
				tt.installedAt,
				tt.completedAt,
				nil,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMs, metadata.DurationMs())
		})
	}
}

func TestInstallationMetadata_DurationSeconds(t *testing.T) {
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := time.Date(2025, 10, 26, 14, 32, 30, 500000000, time.UTC)

	metadata, err := history.NewInstallationMetadata(
		"hyprland",
		"0.45.0",
		installedAt,
		completedAt,
		nil,
	)
	require.NoError(t, err)

	// 150.5 seconds
	assert.InDelta(t, 150.5, metadata.DurationSeconds(), 0.01)
}

func TestInstallationMetadata_PackageCount(t *testing.T) {
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := time.Date(2025, 10, 26, 14, 32, 30, 0, time.UTC)

	t.Run("with packages", func(t *testing.T) {
		pkg1, _ := history.NewInstalledPackage("hyprland", "0.45.0", 15728640)
		pkg2, _ := history.NewInstalledPackage("waybar", "0.10.4", 5242880)
		pkg3, _ := history.NewInstalledPackage("kitty", "0.36.0", 8388608)
		packages := []history.InstalledPackage{pkg1, pkg2, pkg3}

		metadata, err := history.NewInstallationMetadata(
			"hyprland-meta",
			"0.45.0",
			installedAt,
			completedAt,
			packages,
		)
		require.NoError(t, err)

		assert.Equal(t, 3, metadata.PackageCount())
	})

	t.Run("with empty packages", func(t *testing.T) {
		metadata, err := history.NewInstallationMetadata(
			"hyprland-meta",
			"0.45.0",
			installedAt,
			completedAt,
			[]history.InstalledPackage{},
		)
		require.NoError(t, err)

		assert.Equal(t, 0, metadata.PackageCount())
	})
}

func TestInstallationMetadata_TotalSizeBytes(t *testing.T) {
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := time.Date(2025, 10, 26, 14, 32, 30, 0, time.UTC)

	pkg1, _ := history.NewInstalledPackage("hyprland", "0.45.0", 15728640)   // 15 MB
	pkg2, _ := history.NewInstalledPackage("waybar", "0.10.4", 5242880)      // 5 MB
	pkg3, _ := history.NewInstalledPackage("kitty", "0.36.0", 8388608)       // 8 MB
	packages := []history.InstalledPackage{pkg1, pkg2, pkg3}

	metadata, err := history.NewInstallationMetadata(
		"hyprland-meta",
		"0.45.0",
		installedAt,
		completedAt,
		packages,
	)
	require.NoError(t, err)

	// 15 + 5 + 8 = 28 MB
	assert.Equal(t, uint64(29360128), metadata.TotalSizeBytes())
}

func TestInstallationMetadata_TotalSizeMB(t *testing.T) {
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := time.Date(2025, 10, 26, 14, 32, 30, 0, time.UTC)

	pkg1, _ := history.NewInstalledPackage("hyprland", "0.45.0", 15728640)
	pkg2, _ := history.NewInstalledPackage("waybar", "0.10.4", 5242880)
	packages := []history.InstalledPackage{pkg1, pkg2}

	metadata, err := history.NewInstallationMetadata(
		"hyprland-meta",
		"0.45.0",
		installedAt,
		completedAt,
		packages,
	)
	require.NoError(t, err)

	// ~20 MB
	assert.InDelta(t, 20.0, metadata.TotalSizeMB(), 0.1)
}

func TestInstallationMetadata_TrimsWhitespace(t *testing.T) {
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := time.Date(2025, 10, 26, 14, 32, 30, 0, time.UTC)

	metadata, err := history.NewInstallationMetadata(
		"  hyprland  ",
		"  0.45.0  ",
		installedAt,
		completedAt,
		nil,
	)
	require.NoError(t, err)

	assert.Equal(t, "hyprland", metadata.PackageName())
	assert.Equal(t, "0.45.0", metadata.TargetVersion())
}

func TestInstallationMetadata_DefensiveCopy(t *testing.T) {
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := time.Date(2025, 10, 26, 14, 32, 30, 0, time.UTC)

	pkg1, _ := history.NewInstalledPackage("hyprland", "0.45.0", 15728640)
	packages := []history.InstalledPackage{pkg1}

	metadata, err := history.NewInstallationMetadata(
		"hyprland-meta",
		"0.45.0",
		installedAt,
		completedAt,
		packages,
	)
	require.NoError(t, err)

	// Modify original slice
	packages = append(packages, history.InstalledPackage{})

	// Metadata should be unaffected
	assert.Equal(t, 1, metadata.PackageCount())
	assert.Len(t, metadata.InstalledPackages(), 1)
}

func TestInstallationMetadata_HasPackage(t *testing.T) {
	installedAt := time.Date(2025, 10, 26, 14, 30, 0, 0, time.UTC)
	completedAt := time.Date(2025, 10, 26, 14, 32, 30, 0, time.UTC)

	pkg1, _ := history.NewInstalledPackage("hyprland", "0.45.0", 15728640)
	pkg2, _ := history.NewInstalledPackage("waybar", "0.10.4", 5242880)
	packages := []history.InstalledPackage{pkg1, pkg2}

	metadata, err := history.NewInstallationMetadata(
		"hyprland-meta",
		"0.45.0",
		installedAt,
		completedAt,
		packages,
	)
	require.NoError(t, err)

	assert.True(t, metadata.HasPackage("hyprland"))
	assert.True(t, metadata.HasPackage("waybar"))
	assert.False(t, metadata.HasPackage("kitty"))
	assert.False(t, metadata.HasPackage(""))
}
