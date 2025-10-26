package installation_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstalledComponent(t *testing.T) {
	tests := []struct {
		name        string
		component   installation.ComponentName
		version     string
		packageInfo *installation.PackageInfo
		wantErr     bool
		errType     error
	}{
		{
			name:      "valid component with package info",
			component: installation.ComponentHyprland,
			version:   "0.35.0",
			packageInfo: func() *installation.PackageInfo {
				pkg, _ := installation.NewPackageInfo("hyprland", "0.35.0", 50*installation.MB, nil)
				return &pkg
			}(),
			wantErr: false,
		},
		{
			name:        "valid component without package info",
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
			comp, err := installation.NewInstalledComponent(
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
			assert.NotEmpty(t, comp.ID(), "Component should have a unique ID")
			assert.Equal(t, tt.component, comp.Component())
			assert.Equal(t, tt.version, comp.Version())
			assert.False(t, comp.InstalledAt().IsZero(), "InstalledAt should be set")
		})
	}
}

func TestInstalledComponent_Identity(t *testing.T) {
	t.Run("each component has unique ID", func(t *testing.T) {
		comp1, err := installation.NewInstalledComponent(
			installation.ComponentHyprland,
			"0.35.0",
			nil,
		)
		require.NoError(t, err)

		comp2, err := installation.NewInstalledComponent(
			installation.ComponentHyprland,
			"0.35.0",
			nil,
		)
		require.NoError(t, err)

		assert.NotEqual(t, comp1.ID(), comp2.ID(),
			"Different installations should have different IDs")
	})

	t.Run("ID is not empty", func(t *testing.T) {
		comp, err := installation.NewInstalledComponent(
			installation.ComponentWaybar,
			"0.9.20",
			nil,
		)
		require.NoError(t, err)

		assert.NotEmpty(t, comp.ID())
	})
}

func TestInstalledComponent_InstalledAt(t *testing.T) {
	before := time.Now()

	comp, err := installation.NewInstalledComponent(
		installation.ComponentHyprland,
		"0.35.0",
		nil,
	)
	require.NoError(t, err)

	after := time.Now()

	assert.False(t, comp.InstalledAt().Before(before),
		"InstalledAt should be at or after creation start")
	assert.False(t, comp.InstalledAt().After(after),
		"InstalledAt should be at or before creation end")
}

func TestInstalledComponent_IsCore(t *testing.T) {
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
			comp, err := installation.NewInstalledComponent(
				tt.component,
				"1.0.0",
				nil,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.want, comp.IsCore())
		})
	}
}

func TestInstalledComponent_IsDriver(t *testing.T) {
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
			name:      "Hyprland is not a driver",
			component: installation.ComponentHyprland,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp, err := installation.NewInstalledComponent(
				tt.component,
				"1.0.0",
				nil,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.want, comp.IsDriver())
		})
	}
}

func TestInstalledComponent_MarkAsVerified(t *testing.T) {
	comp, err := installation.NewInstalledComponent(
		installation.ComponentHyprland,
		"0.35.0",
		nil,
	)
	require.NoError(t, err)

	// Initially not verified
	assert.False(t, comp.IsVerified(), "New component should not be verified")

	// Mark as verified
	comp.MarkAsVerified()

	// Now should be verified
	assert.True(t, comp.IsVerified(), "Component should be marked as verified")
	assert.False(t, comp.VerifiedAt().IsZero(), "VerifiedAt should be set")
}

func TestInstalledComponent_Age(t *testing.T) {
	comp, err := installation.NewInstalledComponent(
		installation.ComponentWaybar,
		"0.9.20",
		nil,
	)
	require.NoError(t, err)

	// Wait a tiny bit
	time.Sleep(10 * time.Millisecond)

	age := comp.Age()
	assert.Greater(t, age, time.Duration(0), "Age should be positive")
	assert.Less(t, age, 1*time.Second, "Age should be less than 1 second in test")
}

func TestInstalledComponent_String(t *testing.T) {
	t.Run("with package info", func(t *testing.T) {
		pkg, _ := installation.NewPackageInfo("hyprland", "0.35.0", 50*installation.MB, nil)
		comp, err := installation.NewInstalledComponent(
			installation.ComponentHyprland,
			"0.35.0",
			&pkg,
		)
		require.NoError(t, err)

		str := comp.String()
		assert.Contains(t, str, "hyprland")
		assert.Contains(t, str, "0.35.0")
	})

	t.Run("without package info", func(t *testing.T) {
		comp, err := installation.NewInstalledComponent(
			installation.ComponentWaybar,
			"0.9.20",
			nil,
		)
		require.NoError(t, err)

		str := comp.String()
		assert.Contains(t, str, "waybar")
		assert.Contains(t, str, "0.9.20")
	})
}

func TestInstalledComponent_EntityBehavior(t *testing.T) {
	t.Run("entities with same ID are equal", func(t *testing.T) {
		comp1, err := installation.NewInstalledComponent(
			installation.ComponentHyprland,
			"0.35.0",
			nil,
		)
		require.NoError(t, err)

		// Simulate persistence and retrieval (same ID)
		comp2 := comp1

		assert.Equal(t, comp1.ID(), comp2.ID(),
			"Entities with same ID should be considered equal")
	})

	t.Run("entities with different IDs are not equal", func(t *testing.T) {
		comp1, err := installation.NewInstalledComponent(
			installation.ComponentHyprland,
			"0.35.0",
			nil,
		)
		require.NoError(t, err)

		comp2, err := installation.NewInstalledComponent(
			installation.ComponentHyprland,
			"0.35.0",
			nil,
		)
		require.NoError(t, err)

		assert.NotEqual(t, comp1.ID(), comp2.ID(),
			"Different installations should have different IDs")
	})

	t.Run("entity state can change (mutability)", func(t *testing.T) {
		comp, err := installation.NewInstalledComponent(
			installation.ComponentHyprland,
			"0.35.0",
			nil,
		)
		require.NoError(t, err)

		// Entity state changes
		id := comp.ID()
		assert.False(t, comp.IsVerified())

		comp.MarkAsVerified()

		// ID stays same but state changed
		assert.Equal(t, id, comp.ID())
		assert.True(t, comp.IsVerified())
	})
}

func TestInstalledComponent_HasPackageInfo(t *testing.T) {
	t.Run("with package info", func(t *testing.T) {
		pkg, _ := installation.NewPackageInfo("hyprland", "0.35.0", 50*installation.MB, nil)
		comp, err := installation.NewInstalledComponent(
			installation.ComponentHyprland,
			"0.35.0",
			&pkg,
		)
		require.NoError(t, err)

		assert.True(t, comp.HasPackageInfo())
		assert.NotNil(t, comp.PackageInfo())
	})

	t.Run("without package info", func(t *testing.T) {
		comp, err := installation.NewInstalledComponent(
			installation.ComponentWaybar,
			"0.9.20",
			nil,
		)
		require.NoError(t, err)

		assert.False(t, comp.HasPackageInfo())
		assert.Nil(t, comp.PackageInfo())
	})
}
