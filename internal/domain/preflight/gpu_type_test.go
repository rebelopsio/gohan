package preflight_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGPUType(t *testing.T) {
	tests := []struct {
		name    string
		vendor  preflight.GPUVendor
		model   string
		pciID   string
		wantErr bool
		errType error
	}{
		{
			name:    "valid AMD GPU",
			vendor:  preflight.GPUVendorAMD,
			model:   "Radeon RX 7900 XTX",
			pciID:   "1002:744c",
			wantErr: false,
		},
		{
			name:    "valid NVIDIA GPU",
			vendor:  preflight.GPUVendorNVIDIA,
			model:   "GeForce RTX 4090",
			pciID:   "10de:2684",
			wantErr: false,
		},
		{
			name:    "valid Intel GPU",
			vendor:  preflight.GPUVendorIntel,
			model:   "UHD Graphics 770",
			pciID:   "8086:4680",
			wantErr: false,
		},
		{
			name:    "empty vendor",
			vendor:  "",
			model:   "Some GPU",
			pciID:   "0000:0000",
			wantErr: true,
			errType: preflight.ErrInvalidGPU,
		},
		{
			name:    "GPU with empty model",
			vendor:  preflight.GPUVendorAMD,
			model:   "",
			pciID:   "1002:744c",
			wantErr: false,
		},
		{
			name:    "GPU with whitespace model",
			vendor:  preflight.GPUVendorNVIDIA,
			model:   "  GeForce RTX 4090  ",
			pciID:   "10de:2684",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gpu, err := preflight.NewGPUType(tt.vendor, tt.model, tt.pciID)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.vendor, gpu.Vendor())
			assert.Equal(t, tt.pciID, gpu.PCIID())
		})
	}
}

func TestGPUType_IsNVIDIA(t *testing.T) {
	nvidiaGPU, _ := preflight.NewGPUType(preflight.GPUVendorNVIDIA, "RTX 4090", "10de:2684")
	amdGPU, _ := preflight.NewGPUType(preflight.GPUVendorAMD, "RX 7900", "1002:744c")
	intelGPU, _ := preflight.NewGPUType(preflight.GPUVendorIntel, "UHD 770", "8086:4680")

	assert.True(t, nvidiaGPU.IsNVIDIA())
	assert.False(t, amdGPU.IsNVIDIA())
	assert.False(t, intelGPU.IsNVIDIA())
}

func TestGPUType_IsAMD(t *testing.T) {
	nvidiaGPU, _ := preflight.NewGPUType(preflight.GPUVendorNVIDIA, "RTX 4090", "10de:2684")
	amdGPU, _ := preflight.NewGPUType(preflight.GPUVendorAMD, "RX 7900", "1002:744c")
	intelGPU, _ := preflight.NewGPUType(preflight.GPUVendorIntel, "UHD 770", "8086:4680")

	assert.False(t, nvidiaGPU.IsAMD())
	assert.True(t, amdGPU.IsAMD())
	assert.False(t, intelGPU.IsAMD())
}

func TestGPUType_IsIntel(t *testing.T) {
	nvidiaGPU, _ := preflight.NewGPUType(preflight.GPUVendorNVIDIA, "RTX 4090", "10de:2684")
	amdGPU, _ := preflight.NewGPUType(preflight.GPUVendorAMD, "RX 7900", "1002:744c")
	intelGPU, _ := preflight.NewGPUType(preflight.GPUVendorIntel, "UHD 770", "8086:4680")

	assert.False(t, nvidiaGPU.IsIntel())
	assert.False(t, amdGPU.IsIntel())
	assert.True(t, intelGPU.IsIntel())
}

func TestGPUType_RequiresProprietaryDriver(t *testing.T) {
	tests := []struct {
		name                     string
		vendor                   preflight.GPUVendor
		wantProprietaryDriver    bool
	}{
		{
			name:                     "NVIDIA requires proprietary driver",
			vendor:                   preflight.GPUVendorNVIDIA,
			wantProprietaryDriver:    true,
		},
		{
			name:                     "AMD does not require proprietary driver",
			vendor:                   preflight.GPUVendorAMD,
			wantProprietaryDriver:    false,
		},
		{
			name:                     "Intel does not require proprietary driver",
			vendor:                   preflight.GPUVendorIntel,
			wantProprietaryDriver:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gpu, _ := preflight.NewGPUType(tt.vendor, "Test Model", "0000:0000")
			assert.Equal(t, tt.wantProprietaryDriver, gpu.RequiresProprietaryDriver())
		})
	}
}

func TestGPUType_RequiresSpecialConfiguration(t *testing.T) {
	tests := []struct {
		name                        string
		vendor                      preflight.GPUVendor
		wantSpecialConfiguration    bool
	}{
		{
			name:                        "NVIDIA requires special configuration",
			vendor:                      preflight.GPUVendorNVIDIA,
			wantSpecialConfiguration:    true,
		},
		{
			name:                        "AMD does not require special configuration",
			vendor:                      preflight.GPUVendorAMD,
			wantSpecialConfiguration:    false,
		},
		{
			name:                        "Intel does not require special configuration",
			vendor:                      preflight.GPUVendorIntel,
			wantSpecialConfiguration:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gpu, _ := preflight.NewGPUType(tt.vendor, "Test Model", "0000:0000")
			assert.Equal(t, tt.wantSpecialConfiguration, gpu.RequiresSpecialConfiguration())
		})
	}
}

func TestGPUType_String(t *testing.T) {
	tests := []struct {
		name       string
		vendor     preflight.GPUVendor
		model      string
		wantString string
	}{
		{
			name:       "GPU with model",
			vendor:     preflight.GPUVendorNVIDIA,
			model:      "GeForce RTX 4090",
			wantString: "nvidia GeForce RTX 4090",
		},
		{
			name:       "GPU without model",
			vendor:     preflight.GPUVendorAMD,
			model:      "",
			wantString: "amd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gpu, _ := preflight.NewGPUType(tt.vendor, tt.model, "0000:0000")
			assert.Equal(t, tt.wantString, gpu.String())
		})
	}
}
