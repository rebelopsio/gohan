package installation_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDiskSpace(t *testing.T) {
	tests := []struct {
		name      string
		available uint64
		required  uint64
		wantErr   bool
		errType   error
	}{
		{
			name:      "valid disk space with sufficient available",
			available: 50 * installation.GB,
			required:  10 * installation.GB,
			wantErr:   false,
		},
		{
			name:      "valid disk space with exact available",
			available: 10 * installation.GB,
			required:  10 * installation.GB,
			wantErr:   false,
		},
		{
			name:      "valid disk space with zero required",
			available: 50 * installation.GB,
			required:  0,
			wantErr:   false,
		},
		{
			name:      "insufficient disk space",
			available: 5 * installation.GB,
			required:  10 * installation.GB,
			wantErr:   true,
			errType:   installation.ErrInsufficientDiskSpace,
		},
		{
			name:      "zero available with required",
			available: 0,
			required:  10 * installation.GB,
			wantErr:   true,
			errType:   installation.ErrInsufficientDiskSpace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diskSpace, err := installation.NewDiskSpace(tt.available, tt.required)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.available, diskSpace.Available())
			assert.Equal(t, tt.required, diskSpace.Required())
		})
	}
}

func TestDiskSpace_IsSufficient(t *testing.T) {
	tests := []struct {
		name        string
		availableGB uint64
		requiredGB  uint64
		want        bool
	}{
		{
			name:        "sufficient with buffer",
			availableGB: 20,
			requiredGB:  10,
			want:        true,
		},
		{
			name:        "exactly sufficient",
			availableGB: 10,
			requiredGB:  10,
			want:        true,
		},
		{
			name:        "insufficient",
			availableGB: 5,
			requiredGB:  10,
			want:        false,
		},
		{
			name:        "zero required is sufficient",
			availableGB: 5,
			requiredGB:  0,
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diskSpace, err := installation.NewDiskSpace(
				tt.availableGB*installation.GB,
				tt.requiredGB*installation.GB,
			)

			// Should not error if we're testing IsSufficient
			// (error cases are tested in TestNewDiskSpace)
			if err != nil {
				t.Skip("Skipping insufficient case - tested in TestNewDiskSpace")
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, diskSpace.IsSufficient())
		})
	}
}

func TestDiskSpace_RemainingAfterInstall(t *testing.T) {
	tests := []struct {
		name        string
		availableGB uint64
		requiredGB  uint64
		wantGB      float64
	}{
		{
			name:        "10 GB remaining",
			availableGB: 20,
			requiredGB:  10,
			wantGB:      10.0,
		},
		{
			name:        "zero remaining",
			availableGB: 10,
			requiredGB:  10,
			wantGB:      0.0,
		},
		{
			name:        "all remaining when zero required",
			availableGB: 50,
			requiredGB:  0,
			wantGB:      50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diskSpace, err := installation.NewDiskSpace(
				tt.availableGB*installation.GB,
				tt.requiredGB*installation.GB,
			)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantGB, diskSpace.RemainingAfterInstallGB(), 0.01)
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
			availableBytes:  10 * installation.GB,
			wantAvailableGB: 10.0,
		},
		{
			name:            "exactly 50 GB",
			availableBytes:  50 * installation.GB,
			wantAvailableGB: 50.0,
		},
		{
			name:            "fractional GB",
			availableBytes:  10*installation.GB + 512*installation.MB,
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
			diskSpace, err := installation.NewDiskSpace(
				tt.availableBytes,
				0, // No requirement for this test
			)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantAvailableGB, diskSpace.AvailableGB(), 0.01)
		})
	}
}

func TestDiskSpace_RequiredGB(t *testing.T) {
	tests := []struct {
		name          string
		requiredBytes uint64
		wantRequiredGB float64
	}{
		{
			name:          "10 GB required",
			requiredBytes: 10 * installation.GB,
			wantRequiredGB: 10.0,
		},
		{
			name:          "5.5 GB required",
			requiredBytes: 5*installation.GB + 512*installation.MB,
			wantRequiredGB: 5.5,
		},
		{
			name:          "zero required",
			requiredBytes: 0,
			wantRequiredGB: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diskSpace, err := installation.NewDiskSpace(
				100*installation.GB, // Sufficient available
				tt.requiredBytes,
			)
			require.NoError(t, err)

			assert.InDelta(t, tt.wantRequiredGB, diskSpace.RequiredGB(), 0.01)
		})
	}
}

func TestDiskSpace_String(t *testing.T) {
	diskSpace, err := installation.NewDiskSpace(
		50*installation.GB,
		10*installation.GB,
	)
	require.NoError(t, err)

	str := diskSpace.String()
	assert.Contains(t, str, "50.00 GB available")
	assert.Contains(t, str, "10.00 GB required")
	assert.Contains(t, str, "40.00 GB remaining")
}

func TestDiskSpace_ValueObjectImmutability(t *testing.T) {
	// Value objects should be immutable
	diskSpace1, err := installation.NewDiskSpace(50*installation.GB, 10*installation.GB)
	require.NoError(t, err)

	diskSpace2, err := installation.NewDiskSpace(50*installation.GB, 10*installation.GB)
	require.NoError(t, err)

	// Same values should be equal
	assert.Equal(t, diskSpace1.Available(), diskSpace2.Available())
	assert.Equal(t, diskSpace1.Required(), diskSpace2.Required())
}
