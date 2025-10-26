package preflight_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDiskSpace(t *testing.T) {
	tests := []struct {
		name      string
		available uint64
		total     uint64
		path      string
		wantErr   bool
		errType   error
		wantPath  string
	}{
		{
			name:      "valid disk space",
			available: 50 * preflight.GB,
			total:     100 * preflight.GB,
			path:      "/",
			wantErr:   false,
			wantPath:  "/",
		},
		{
			name:      "disk space with custom path",
			available: 100 * preflight.GB,
			total:     500 * preflight.GB,
			path:      "/home",
			wantErr:   false,
			wantPath:  "/home",
		},
		{
			name:      "disk space with empty path defaults to root",
			available: 20 * preflight.GB,
			total:     50 * preflight.GB,
			path:      "",
			wantErr:   false,
			wantPath:  "/",
		},
		{
			name:      "available exceeds total",
			available: 100 * preflight.GB,
			total:     50 * preflight.GB,
			path:      "/",
			wantErr:   true,
			errType:   preflight.ErrInvalidDiskSpace,
		},
		{
			name:      "available equals total",
			available: 100 * preflight.GB,
			total:     100 * preflight.GB,
			path:      "/",
			wantErr:   false,
			wantPath:  "/",
		},
		{
			name:      "zero available",
			available: 0,
			total:     100 * preflight.GB,
			path:      "/",
			wantErr:   false,
			wantPath:  "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diskSpace, err := preflight.NewDiskSpace(tt.available, tt.total, tt.path)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.available, diskSpace.Available())
			assert.Equal(t, tt.total, diskSpace.Total())
			assert.Equal(t, tt.wantPath, diskSpace.Path())
		})
	}
}

func TestDiskSpace_MeetsMinimum(t *testing.T) {
	tests := []struct {
		name             string
		availableGB      uint64
		totalGB          uint64
		requiredGB       uint64
		wantMeetsMinimum bool
	}{
		{
			name:             "meets minimum exactly",
			availableGB:      10,
			totalGB:          50,
			requiredGB:       10,
			wantMeetsMinimum: true,
		},
		{
			name:             "exceeds minimum",
			availableGB:      20,
			totalGB:          50,
			requiredGB:       10,
			wantMeetsMinimum: true,
		},
		{
			name:             "below minimum",
			availableGB:      5,
			totalGB:          50,
			requiredGB:       10,
			wantMeetsMinimum: false,
		},
		{
			name:             "zero available below minimum",
			availableGB:      0,
			totalGB:          50,
			requiredGB:       10,
			wantMeetsMinimum: false,
		},
		{
			name:             "large disk exceeds minimum",
			availableGB:      500,
			totalGB:          1000,
			requiredGB:       10,
			wantMeetsMinimum: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diskSpace, err := preflight.NewDiskSpace(
				tt.availableGB*preflight.GB,
				tt.totalGB*preflight.GB,
				"/",
			)
			require.NoError(t, err)

			assert.Equal(t, tt.wantMeetsMinimum, diskSpace.MeetsMinimum(tt.requiredGB))
		})
	}
}

func TestDiskSpace_AvailableGB(t *testing.T) {
	tests := []struct {
		name            string
		availableBytes  uint64
		wantAvailableGB float64
	}{
		{
			name:            "exactly 10 GB",
			availableBytes:  10 * preflight.GB,
			wantAvailableGB: 10.0,
		},
		{
			name:            "exactly 50 GB",
			availableBytes:  50 * preflight.GB,
			wantAvailableGB: 50.0,
		},
		{
			name:            "fractional GB",
			availableBytes:  10*preflight.GB + 512*preflight.MB,
			wantAvailableGB: 10.5,
		},
		{
			name:            "zero bytes",
			availableBytes:  0,
			wantAvailableGB: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diskSpace, err := preflight.NewDiskSpace(
				tt.availableBytes,
				100*preflight.GB,
				"/",
			)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantAvailableGB, diskSpace.AvailableGB(), 0.01)
		})
	}
}

func TestDiskSpace_TotalGB(t *testing.T) {
	tests := []struct {
		name       string
		totalBytes uint64
		wantTotalGB float64
	}{
		{
			name:       "100 GB disk",
			totalBytes: 100 * preflight.GB,
			wantTotalGB: 100.0,
		},
		{
			name:       "500 GB disk",
			totalBytes: 500 * preflight.GB,
			wantTotalGB: 500.0,
		},
		{
			name:       "1 TB disk",
			totalBytes: 1000 * preflight.GB,
			wantTotalGB: 1000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diskSpace, err := preflight.NewDiskSpace(
				50*preflight.GB,
				tt.totalBytes,
				"/",
			)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantTotalGB, diskSpace.TotalGB(), 0.01)
		})
	}
}

func TestDiskSpace_UsagePercent(t *testing.T) {
	tests := []struct {
		name             string
		availableGB      uint64
		totalGB          uint64
		wantUsagePercent float64
	}{
		{
			name:             "50% used",
			availableGB:      50,
			totalGB:          100,
			wantUsagePercent: 50.0,
		},
		{
			name:             "75% used",
			availableGB:      25,
			totalGB:          100,
			wantUsagePercent: 75.0,
		},
		{
			name:             "10% used",
			availableGB:      90,
			totalGB:          100,
			wantUsagePercent: 10.0,
		},
		{
			name:             "100% used (0 available)",
			availableGB:      0,
			totalGB:          100,
			wantUsagePercent: 100.0,
		},
		{
			name:             "0% used (all available)",
			availableGB:      100,
			totalGB:          100,
			wantUsagePercent: 0.0,
		},
		{
			name:             "zero total returns 0",
			availableGB:      0,
			totalGB:          0,
			wantUsagePercent: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diskSpace, err := preflight.NewDiskSpace(
				tt.availableGB*preflight.GB,
				tt.totalGB*preflight.GB,
				"/",
			)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantUsagePercent, diskSpace.UsagePercent(), 0.01)
		})
	}
}

func TestDiskSpace_String(t *testing.T) {
	diskSpace, err := preflight.NewDiskSpace(
		50*preflight.GB,
		100*preflight.GB,
		"/",
	)
	require.NoError(t, err)

	str := diskSpace.String()
	assert.Contains(t, str, "50.00 GB available")
	assert.Contains(t, str, "100.00 GB total")
	assert.Contains(t, str, "50.0% used")
}
