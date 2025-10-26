package installation_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewComponentSelection(t *testing.T) {
	tests := []struct {
		name        string
		component   installation.ComponentName
		version     string
		packageInfo *installation.PackageInfo
		wantErr     bool
		errType     error
	}{
		{
			name:      "valid selection with package info",
			component: installation.ComponentHyprland,
			version:   "0.35.0",
			packageInfo: func() *installation.PackageInfo {
				pkg, _ := installation.NewPackageInfo("hyprland", "0.35.0", 50*installation.MB, nil)
				return &pkg
			}(),
			wantErr: false,
		},
		{
			name:        "valid selection without package info",
			component:   installation.ComponentWaybar,
			version:     "0.9.20",
			packageInfo: nil,
			wantErr:     false,
		},
		{
			name:        "empty version",
			component:   installation.ComponentHyprland,
			version:     "",
			packageInfo: nil,
			wantErr:     true,
			errType:     installation.ErrInvalidComponentSelection,
		},
		{
			name:        "whitespace version",
			component:   installation.ComponentHyprland,
			version:     "   ",
			packageInfo: nil,
			wantErr:     true,
			errType:     installation.ErrInvalidComponentSelection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selection, err := installation.NewComponentSelection(
				tt.component,
				tt.version,
				tt.packageInfo,
			)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.component, selection.Component())
			assert.Equal(t, tt.version, selection.Version())
			if tt.packageInfo != nil {
				assert.True(t, selection.HasPackageInfo())
				assert.NotNil(t, selection.PackageInfo())
			} else {
				assert.False(t, selection.HasPackageInfo())
				assert.Nil(t, selection.PackageInfo())
			}
		})
	}
}

func TestComponentSelection_IsCore(t *testing.T) {
	tests := []struct {
		name      string
		component installation.ComponentName
		want      bool
	}{
		{
			name:      "Hyprland is core",
			component: installation.ComponentHyprland,
			want:      true,
		},
		{
			name:      "Waybar is not core",
			component: installation.ComponentWaybar,
			want:      false,
		},
		{
			name:      "AMD driver is not core",
			component: installation.ComponentAMDDriver,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selection, err := installation.NewComponentSelection(
				tt.component,
				"1.0.0",
				nil,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.want, selection.IsCore())
		})
	}
}

func TestComponentSelection_IsDriver(t *testing.T) {
	tests := []struct {
		name      string
		component installation.ComponentName
		want      bool
	}{
		{
			name:      "AMD driver is a driver",
			component: installation.ComponentAMDDriver,
			want:      true,
		},
		{
			name:      "NVIDIA driver is a driver",
			component: installation.ComponentNVIDIADriver,
			want:      true,
		},
		{
			name:      "Intel driver is a driver",
			component: installation.ComponentIntelDriver,
			want:      true,
		},
		{
			name:      "Hyprland is not a driver",
			component: installation.ComponentHyprland,
			want:      false,
		},
		{
			name:      "Waybar is not a driver",
			component: installation.ComponentWaybar,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selection, err := installation.NewComponentSelection(
				tt.component,
				"1.0.0",
				nil,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.want, selection.IsDriver())
		})
	}
}

func TestComponentSelection_EstimatedSizeBytes(t *testing.T) {
	tests := []struct {
		name        string
		packageInfo *installation.PackageInfo
		want        uint64
	}{
		{
			name: "with package info",
			packageInfo: func() *installation.PackageInfo {
				pkg, _ := installation.NewPackageInfo("hyprland", "0.35.0", 50*installation.MB, nil)
				return &pkg
			}(),
			want: 50 * installation.MB,
		},
		{
			name:        "without package info",
			packageInfo: nil,
			want:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selection, err := installation.NewComponentSelection(
				installation.ComponentHyprland,
				"1.0.0",
				tt.packageInfo,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.want, selection.EstimatedSizeBytes())
		})
	}
}

func TestComponentSelection_String(t *testing.T) {
	t.Run("with package info", func(t *testing.T) {
		pkg, err := installation.NewPackageInfo("hyprland", "0.35.0", 50*installation.MB, nil)
		require.NoError(t, err)

		selection, err := installation.NewComponentSelection(
			installation.ComponentHyprland,
			"0.35.0",
			&pkg,
		)
		require.NoError(t, err)

		str := selection.String()
		assert.Contains(t, str, "hyprland")
		assert.Contains(t, str, "0.35.0")
	})

	t.Run("without package info", func(t *testing.T) {
		selection, err := installation.NewComponentSelection(
			installation.ComponentWaybar,
			"0.9.20",
			nil,
		)
		require.NoError(t, err)

		str := selection.String()
		assert.Contains(t, str, "waybar")
		assert.Contains(t, str, "0.9.20")
	})
}

func TestComponentSelection_ValueObjectImmutability(t *testing.T) {
	// Value objects should be immutable
	pkg1, err := installation.NewPackageInfo("pkg1", "1.0.0", 10*installation.MB, nil)
	require.NoError(t, err)

	selection1, err := installation.NewComponentSelection(
		installation.ComponentHyprland,
		"1.0.0",
		&pkg1,
	)
	require.NoError(t, err)

	pkg2, err := installation.NewPackageInfo("pkg1", "1.0.0", 10*installation.MB, nil)
	require.NoError(t, err)

	selection2, err := installation.NewComponentSelection(
		installation.ComponentHyprland,
		"1.0.0",
		&pkg2,
	)
	require.NoError(t, err)

	// Same values should be equal
	assert.Equal(t, selection1.Component(), selection2.Component())
	assert.Equal(t, selection1.Version(), selection2.Version())
	assert.Equal(t, selection1.HasPackageInfo(), selection2.HasPackageInfo())
}

func TestComponentSelection_EdgeCases(t *testing.T) {
	t.Run("version trimming", func(t *testing.T) {
		selection, err := installation.NewComponentSelection(
			installation.ComponentHyprland,
			"  1.0.0  ",
			nil,
		)
		require.NoError(t, err)
		assert.Equal(t, "1.0.0", selection.Version())
	})

	t.Run("nil package info is valid", func(t *testing.T) {
		selection, err := installation.NewComponentSelection(
			installation.ComponentHyprland,
			"1.0.0",
			nil,
		)
		require.NoError(t, err)
		assert.False(t, selection.HasPackageInfo())
		assert.Nil(t, selection.PackageInfo())
		assert.Equal(t, uint64(0), selection.EstimatedSizeBytes())
	})
}
