package installation_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPackageInfo(t *testing.T) {
	tests := []struct {
		name         string
		packageName  string
		version      string
		sizeBytes    uint64
		dependencies []string
		wantErr      bool
		errType      error
	}{
		{
			name:         "valid package with dependencies",
			packageName:  "hyprland",
			version:      "0.35.0",
			sizeBytes:    50 * installation.MB,
			dependencies: []string{"wayland", "wlroots"},
			wantErr:      false,
		},
		{
			name:         "valid package without dependencies",
			packageName:  "waybar",
			version:      "0.9.20",
			sizeBytes:    10 * installation.MB,
			dependencies: nil,
			wantErr:      false,
		},
		{
			name:         "empty package name",
			packageName:  "",
			version:      "1.0.0",
			sizeBytes:    1 * installation.MB,
			dependencies: nil,
			wantErr:      true,
			errType:      installation.ErrInvalidPackageInfo,
		},
		{
			name:         "whitespace package name",
			packageName:  "   ",
			version:      "1.0.0",
			sizeBytes:    1 * installation.MB,
			dependencies: nil,
			wantErr:      true,
			errType:      installation.ErrInvalidPackageInfo,
		},
		{
			name:         "empty version",
			packageName:  "hyprland",
			version:      "",
			sizeBytes:    1 * installation.MB,
			dependencies: nil,
			wantErr:      true,
			errType:      installation.ErrInvalidPackageInfo,
		},
		{
			name:         "zero size",
			packageName:  "hyprland",
			version:      "1.0.0",
			sizeBytes:    0,
			dependencies: nil,
			wantErr:      false, // Zero size is valid (metadata packages)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := installation.NewPackageInfo(
				tt.packageName,
				tt.version,
				tt.sizeBytes,
				tt.dependencies,
			)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.packageName, pkg.Name())
			assert.Equal(t, tt.version, pkg.Version())
			assert.Equal(t, tt.sizeBytes, pkg.SizeBytes())
			assert.Equal(t, tt.dependencies, pkg.Dependencies())
		})
	}
}

func TestPackageInfo_SizeMB(t *testing.T) {
	tests := []struct {
		name      string
		sizeBytes uint64
		wantMB    float64
	}{
		{
			name:      "exactly 10 MB",
			sizeBytes: 10 * installation.MB,
			wantMB:    10.0,
		},
		{
			name:      "exactly 50 MB",
			sizeBytes: 50 * installation.MB,
			wantMB:    50.0,
		},
		{
			name:      "fractional MB",
			sizeBytes: 10*installation.MB + 512*1024, // 10.5 MB
			wantMB:    10.5,
		},
		{
			name:      "zero bytes",
			sizeBytes: 0,
			wantMB:    0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := installation.NewPackageInfo(
				"test-package",
				"1.0.0",
				tt.sizeBytes,
				nil,
			)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantMB, pkg.SizeMB(), 0.01)
		})
	}
}

func TestPackageInfo_HasDependencies(t *testing.T) {
	tests := []struct {
		name         string
		dependencies []string
		want         bool
	}{
		{
			name:         "has dependencies",
			dependencies: []string{"dep1", "dep2"},
			want:         true,
		},
		{
			name:         "no dependencies (nil)",
			dependencies: nil,
			want:         false,
		},
		{
			name:         "no dependencies (empty slice)",
			dependencies: []string{},
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := installation.NewPackageInfo(
				"test-package",
				"1.0.0",
				1*installation.MB,
				tt.dependencies,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.want, pkg.HasDependencies())
		})
	}
}

func TestPackageInfo_DependencyCount(t *testing.T) {
	tests := []struct {
		name         string
		dependencies []string
		wantCount    int
	}{
		{
			name:         "multiple dependencies",
			dependencies: []string{"dep1", "dep2", "dep3"},
			wantCount:    3,
		},
		{
			name:         "single dependency",
			dependencies: []string{"dep1"},
			wantCount:    1,
		},
		{
			name:         "no dependencies",
			dependencies: nil,
			wantCount:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := installation.NewPackageInfo(
				"test-package",
				"1.0.0",
				1*installation.MB,
				tt.dependencies,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.wantCount, pkg.DependencyCount())
		})
	}
}

func TestPackageInfo_String(t *testing.T) {
	pkg, err := installation.NewPackageInfo(
		"hyprland",
		"0.35.0",
		50*installation.MB,
		[]string{"wayland", "wlroots"},
	)
	require.NoError(t, err)

	str := pkg.String()
	assert.Contains(t, str, "hyprland")
	assert.Contains(t, str, "0.35.0")
	assert.Contains(t, str, "50.00 MB")
}

func TestPackageInfo_ValueObjectImmutability(t *testing.T) {
	// Value objects should be immutable
	deps1 := []string{"dep1", "dep2"}
	pkg1, err := installation.NewPackageInfo("package1", "1.0.0", 10*installation.MB, deps1)
	require.NoError(t, err)

	deps2 := []string{"dep1", "dep2"}
	pkg2, err := installation.NewPackageInfo("package1", "1.0.0", 10*installation.MB, deps2)
	require.NoError(t, err)

	// Same values should be equal
	assert.Equal(t, pkg1.Name(), pkg2.Name())
	assert.Equal(t, pkg1.Version(), pkg2.Version())
	assert.Equal(t, pkg1.SizeBytes(), pkg2.SizeBytes())
	assert.Equal(t, pkg1.Dependencies(), pkg2.Dependencies())

	// Modifying original slice should not affect value object
	deps1[0] = "modified"
	assert.NotEqual(t, deps1[0], pkg1.Dependencies()[0])
}

func TestPackageInfo_EdgeCases(t *testing.T) {
	t.Run("package name trimming", func(t *testing.T) {
		// Package names should be trimmed of whitespace
		pkg, err := installation.NewPackageInfo(
			"  hyprland  ",
			"1.0.0",
			1*installation.MB,
			nil,
		)
		require.NoError(t, err)
		assert.Equal(t, "hyprland", pkg.Name())
	})

	t.Run("version trimming", func(t *testing.T) {
		// Versions should be trimmed of whitespace
		pkg, err := installation.NewPackageInfo(
			"hyprland",
			"  1.0.0  ",
			1*installation.MB,
			nil,
		)
		require.NoError(t, err)
		assert.Equal(t, "1.0.0", pkg.Version())
	})

	t.Run("large package size", func(t *testing.T) {
		// Should handle large packages (multiple GB)
		pkg, err := installation.NewPackageInfo(
			"large-package",
			"1.0.0",
			5*installation.GB,
			nil,
		)
		require.NoError(t, err)
		assert.Equal(t, float64(5*1024), pkg.SizeMB()) // 5 GB = 5120 MB
	})
}
