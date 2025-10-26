package installation_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGPUSupport(t *testing.T) {
	tests := []struct {
		name            string
		vendor          string
		requiresDriver  bool
		driverComponent installation.ComponentName
		wantErr         bool
		errType         error
	}{
		{
			name:            "AMD GPU support",
			vendor:          "amd",
			requiresDriver:  true,
			driverComponent: installation.ComponentAMDDriver,
			wantErr:         false,
		},
		{
			name:            "NVIDIA GPU support",
			vendor:          "nvidia",
			requiresDriver:  true,
			driverComponent: installation.ComponentNVIDIADriver,
			wantErr:         false,
		},
		{
			name:            "Intel GPU support",
			vendor:          "intel",
			requiresDriver:  true,
			driverComponent: installation.ComponentIntelDriver,
			wantErr:         false,
		},
		{
			name:            "no driver required",
			vendor:          "amd",
			requiresDriver:  false,
			driverComponent: "",
			wantErr:         false,
		},
		{
			name:            "empty vendor",
			vendor:          "",
			requiresDriver:  true,
			driverComponent: installation.ComponentAMDDriver,
			wantErr:         true,
			errType:         installation.ErrInvalidGPUSupport,
		},
		{
			name:            "whitespace vendor",
			vendor:          "   ",
			requiresDriver:  true,
			driverComponent: installation.ComponentAMDDriver,
			wantErr:         true,
			errType:         installation.ErrInvalidGPUSupport,
		},
		{
			name:            "requires driver but no component specified",
			vendor:          "amd",
			requiresDriver:  true,
			driverComponent: "",
			wantErr:         true,
			errType:         installation.ErrInvalidGPUSupport,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gpu, err := installation.NewGPUSupport(
				tt.vendor,
				tt.requiresDriver,
				tt.driverComponent,
			)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.vendor, gpu.Vendor())
			assert.Equal(t, tt.requiresDriver, gpu.RequiresDriver())
			if tt.requiresDriver {
				assert.Equal(t, tt.driverComponent, gpu.DriverComponent())
			}
		})
	}
}

func TestGPUSupport_IsAMD(t *testing.T) {
	tests := []struct {
		name   string
		vendor string
		want   bool
	}{
		{
			name:   "AMD GPU",
			vendor: "amd",
			want:   true,
		},
		{
			name:   "NVIDIA GPU",
			vendor: "nvidia",
			want:   false,
		},
		{
			name:   "Intel GPU",
			vendor: "intel",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gpu, err := installation.NewGPUSupport(tt.vendor, false, "")
			require.NoError(t, err)

			assert.Equal(t, tt.want, gpu.IsAMD())
		})
	}
}

func TestGPUSupport_IsNVIDIA(t *testing.T) {
	tests := []struct {
		name   string
		vendor string
		want   bool
	}{
		{
			name:   "NVIDIA GPU",
			vendor: "nvidia",
			want:   true,
		},
		{
			name:   "AMD GPU",
			vendor: "amd",
			want:   false,
		},
		{
			name:   "Intel GPU",
			vendor: "intel",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gpu, err := installation.NewGPUSupport(tt.vendor, false, "")
			require.NoError(t, err)

			assert.Equal(t, tt.want, gpu.IsNVIDIA())
		})
	}
}

func TestGPUSupport_IsIntel(t *testing.T) {
	tests := []struct {
		name   string
		vendor string
		want   bool
	}{
		{
			name:   "Intel GPU",
			vendor: "intel",
			want:   true,
		},
		{
			name:   "AMD GPU",
			vendor: "amd",
			want:   false,
		},
		{
			name:   "NVIDIA GPU",
			vendor: "nvidia",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gpu, err := installation.NewGPUSupport(tt.vendor, false, "")
			require.NoError(t, err)

			assert.Equal(t, tt.want, gpu.IsIntel())
		})
	}
}

func TestGPUSupport_RequiresProprietary(t *testing.T) {
	tests := []struct {
		name   string
		vendor string
		want   bool
	}{
		{
			name:   "NVIDIA requires proprietary",
			vendor: "nvidia",
			want:   true,
		},
		{
			name:   "AMD does not require proprietary",
			vendor: "amd",
			want:   false,
		},
		{
			name:   "Intel does not require proprietary",
			vendor: "intel",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gpu, err := installation.NewGPUSupport(tt.vendor, false, "")
			require.NoError(t, err)

			assert.Equal(t, tt.want, gpu.RequiresProprietary())
		})
	}
}

func TestGPUSupport_String(t *testing.T) {
	t.Run("with driver", func(t *testing.T) {
		gpu, err := installation.NewGPUSupport(
			"amd",
			true,
			installation.ComponentAMDDriver,
		)
		require.NoError(t, err)

		str := gpu.String()
		assert.Contains(t, str, "amd")
		assert.Contains(t, str, "amd_driver")
	})

	t.Run("without driver", func(t *testing.T) {
		gpu, err := installation.NewGPUSupport("amd", false, "")
		require.NoError(t, err)

		str := gpu.String()
		assert.Contains(t, str, "amd")
		assert.Contains(t, str, "no driver")
	})
}

func TestGPUSupport_ValueObjectImmutability(t *testing.T) {
	// Value objects should be immutable
	gpu1, err := installation.NewGPUSupport("amd", true, installation.ComponentAMDDriver)
	require.NoError(t, err)

	gpu2, err := installation.NewGPUSupport("amd", true, installation.ComponentAMDDriver)
	require.NoError(t, err)

	// Same values should be equal
	assert.Equal(t, gpu1.Vendor(), gpu2.Vendor())
	assert.Equal(t, gpu1.RequiresDriver(), gpu2.RequiresDriver())
	assert.Equal(t, gpu1.DriverComponent(), gpu2.DriverComponent())
}

func TestGPUSupport_EdgeCases(t *testing.T) {
	t.Run("vendor trimming and lowercasing", func(t *testing.T) {
		gpu, err := installation.NewGPUSupport("  AMD  ", false, "")
		require.NoError(t, err)
		assert.Equal(t, "amd", gpu.Vendor())
		assert.True(t, gpu.IsAMD())
	})

	t.Run("case insensitive vendor matching", func(t *testing.T) {
		gpu, err := installation.NewGPUSupport("NVIDIA", false, "")
		require.NoError(t, err)
		assert.Equal(t, "nvidia", gpu.Vendor())
		assert.True(t, gpu.IsNVIDIA())
	})

	t.Run("driver component must be a driver", func(t *testing.T) {
		// Should not allow non-driver components
		_, err := installation.NewGPUSupport(
			"amd",
			true,
			installation.ComponentHyprland, // Not a driver
		)
		assert.ErrorIs(t, err, installation.ErrInvalidGPUSupport)
	})
}

func TestGPUSupport_DriverValidation(t *testing.T) {
	tests := []struct {
		name            string
		vendor          string
		driverComponent installation.ComponentName
		shouldMatch     bool
	}{
		{
			name:            "AMD vendor with AMD driver",
			vendor:          "amd",
			driverComponent: installation.ComponentAMDDriver,
			shouldMatch:     true,
		},
		{
			name:            "NVIDIA vendor with NVIDIA driver",
			vendor:          "nvidia",
			driverComponent: installation.ComponentNVIDIADriver,
			shouldMatch:     true,
		},
		{
			name:            "Intel vendor with Intel driver",
			vendor:          "intel",
			driverComponent: installation.ComponentIntelDriver,
			shouldMatch:     true,
		},
		{
			name:            "AMD vendor with NVIDIA driver (mismatch)",
			vendor:          "amd",
			driverComponent: installation.ComponentNVIDIADriver,
			shouldMatch:     false,
		},
		{
			name:            "NVIDIA vendor with AMD driver (mismatch)",
			vendor:          "nvidia",
			driverComponent: installation.ComponentAMDDriver,
			shouldMatch:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := installation.NewGPUSupport(
				tt.vendor,
				true,
				tt.driverComponent,
			)

			if tt.shouldMatch {
				assert.NoError(t, err, "Should allow matching vendor and driver")
			} else {
				assert.ErrorIs(t, err, installation.ErrInvalidGPUSupport,
					"Should reject mismatched vendor and driver")
			}
		})
	}
}
