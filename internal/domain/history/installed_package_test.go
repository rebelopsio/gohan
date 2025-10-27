package history_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstalledPackage(t *testing.T) {
	tests := []struct {
		name      string
		pkgName   string
		version   string
		sizeBytes uint64
		wantErr   error
	}{
		{
			name:      "valid package",
			pkgName:   "hyprland",
			version:   "0.35.0",
			sizeBytes: 1024000,
			wantErr:   nil,
		},
		{
			name:      "package with zero size",
			pkgName:   "waybar",
			version:   "1.0.0",
			sizeBytes: 0,
			wantErr:   nil,
		},
		{
			name:      "package with whitespace in name",
			pkgName:   "  hyprland  ",
			version:   "0.35.0",
			sizeBytes: 1024000,
			wantErr:   nil,
		},
		{
			name:      "empty package name",
			pkgName:   "",
			version:   "0.35.0",
			sizeBytes: 1024000,
			wantErr:   history.ErrInvalidPackageName,
		},
		{
			name:      "whitespace-only package name",
			pkgName:   "   ",
			version:   "0.35.0",
			sizeBytes: 1024000,
			wantErr:   history.ErrInvalidPackageName,
		},
		{
			name:      "empty version",
			pkgName:   "hyprland",
			version:   "",
			sizeBytes: 1024000,
			wantErr:   history.ErrInvalidPackageVersion,
		},
		{
			name:      "whitespace-only version",
			pkgName:   "hyprland",
			version:   "   ",
			sizeBytes: 1024000,
			wantErr:   history.ErrInvalidPackageVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := history.NewInstalledPackage(tt.pkgName, tt.version, tt.sizeBytes)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, pkg.Name())
				assert.NotEmpty(t, pkg.Version())
			}
		})
	}
}

func TestInstalledPackage_Accessors(t *testing.T) {
	pkg, err := history.NewInstalledPackage("hyprland", "0.35.0", 10485760) // 10 MB
	require.NoError(t, err)

	assert.Equal(t, "hyprland", pkg.Name())
	assert.Equal(t, "0.35.0", pkg.Version())
	assert.Equal(t, uint64(10485760), pkg.SizeBytes())
}

func TestInstalledPackage_SizeMB(t *testing.T) {
	tests := []struct {
		name       string
		sizeBytes  uint64
		expectedMB float64
	}{
		{
			name:       "10 MB",
			sizeBytes:  10485760,
			expectedMB: 10.0,
		},
		{
			name:       "1.5 MB",
			sizeBytes:  1572864,
			expectedMB: 1.5,
		},
		{
			name:       "zero bytes",
			sizeBytes:  0,
			expectedMB: 0.0,
		},
		{
			name:       "less than 1 MB",
			sizeBytes:  524288,
			expectedMB: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := history.NewInstalledPackage("test", "1.0.0", tt.sizeBytes)
			require.NoError(t, err)

			assert.InDelta(t, tt.expectedMB, pkg.SizeMB(), 0.01)
		})
	}
}

func TestInstalledPackage_Equals(t *testing.T) {
	pkg1, err := history.NewInstalledPackage("hyprland", "0.35.0", 1024000)
	require.NoError(t, err)

	pkg2, err := history.NewInstalledPackage("hyprland", "0.35.0", 2048000) // Different size
	require.NoError(t, err)

	pkg3, err := history.NewInstalledPackage("hyprland", "0.36.0", 1024000) // Different version
	require.NoError(t, err)

	pkg4, err := history.NewInstalledPackage("waybar", "0.35.0", 1024000) // Different name
	require.NoError(t, err)

	t.Run("same name and version equals regardless of size", func(t *testing.T) {
		assert.True(t, pkg1.Equals(pkg2))
	})

	t.Run("different version not equal", func(t *testing.T) {
		assert.False(t, pkg1.Equals(pkg3))
	})

	t.Run("different name not equal", func(t *testing.T) {
		assert.False(t, pkg1.Equals(pkg4))
	})
}

func TestInstalledPackage_String(t *testing.T) {
	pkg, err := history.NewInstalledPackage("hyprland", "0.35.0", 10485760)
	require.NoError(t, err)

	str := pkg.String()
	assert.Contains(t, str, "hyprland")
	assert.Contains(t, str, "0.35.0")
	assert.Contains(t, str, "10.00")
}

func TestInstalledPackage_TrimsWhitespace(t *testing.T) {
	pkg, err := history.NewInstalledPackage("  hyprland  ", "  0.35.0  ", 1024000)
	require.NoError(t, err)

	assert.Equal(t, "hyprland", pkg.Name())
	assert.Equal(t, "0.35.0", pkg.Version())
}
