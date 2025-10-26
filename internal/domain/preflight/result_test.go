package preflight_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
	"github.com/stretchr/testify/assert"
)

func TestNewValidationResult(t *testing.T) {
	guidance := preflight.NewUserGuidance(
		"Test failed",
		"Reason for failure",
		[]string{"Fix step 1"},
		"https://docs.example.com",
	)

	result := preflight.NewValidationResult(
		preflight.RequirementDebianVersion,
		preflight.StatusFail,
		preflight.SeverityCritical,
		"bookworm",
		"sid or trixie",
		guidance,
	)

	assert.NotEmpty(t, result.ID())
	assert.Equal(t, preflight.RequirementDebianVersion, result.RequirementName())
	assert.Equal(t, preflight.StatusFail, result.Status())
	assert.Equal(t, preflight.SeverityCritical, result.Severity())
	assert.Equal(t, "bookworm", result.ActualValue())
	assert.Equal(t, "sid or trixie", result.ExpectedValue())
	assert.Equal(t, guidance.Message(), result.Guidance().Message())
	assert.False(t, result.DetectedAt().IsZero())
}

func TestValidationResult_IsBlocking(t *testing.T) {
	tests := []struct {
		name        string
		status      preflight.ValidationStatus
		severity    preflight.Severity
		wantBlocking bool
	}{
		{
			name:        "critical failure blocks",
			status:      preflight.StatusFail,
			severity:    preflight.SeverityCritical,
			wantBlocking: true,
		},
		{
			name:        "high severity failure blocks",
			status:      preflight.StatusFail,
			severity:    preflight.SeverityHigh,
			wantBlocking: true,
		},
		{
			name:        "medium severity failure does not block",
			status:      preflight.StatusFail,
			severity:    preflight.SeverityMedium,
			wantBlocking: false,
		},
		{
			name:        "low severity failure does not block",
			status:      preflight.StatusFail,
			severity:    preflight.SeverityLow,
			wantBlocking: false,
		},
		{
			name:        "warning status does not block",
			status:      preflight.StatusWarning,
			severity:    preflight.SeverityCritical,
			wantBlocking: false,
		},
		{
			name:        "pass status does not block",
			status:      preflight.StatusPass,
			severity:    preflight.SeverityCritical,
			wantBlocking: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preflight.NewValidationResult(
				preflight.RequirementDebianVersion,
				tt.status,
				tt.severity,
				"actual",
				"expected",
				preflight.NewUserGuidance("", "", nil, ""),
			)

			assert.Equal(t, tt.wantBlocking, result.IsBlocking())
		})
	}
}

func TestValidationResult_IsWarning(t *testing.T) {
	tests := []struct {
		name       string
		status     preflight.ValidationStatus
		severity   preflight.Severity
		wantWarning bool
	}{
		{
			name:       "warning status is warning",
			status:     preflight.StatusWarning,
			severity:   preflight.SeverityCritical,
			wantWarning: true,
		},
		{
			name:       "medium severity failure is warning",
			status:     preflight.StatusFail,
			severity:   preflight.SeverityMedium,
			wantWarning: true,
		},
		{
			name:       "low severity failure is warning",
			status:     preflight.StatusFail,
			severity:   preflight.SeverityLow,
			wantWarning: true,
		},
		{
			name:       "critical failure is not warning",
			status:     preflight.StatusFail,
			severity:   preflight.SeverityCritical,
			wantWarning: false,
		},
		{
			name:       "high severity failure is not warning",
			status:     preflight.StatusFail,
			severity:   preflight.SeverityHigh,
			wantWarning: false,
		},
		{
			name:       "pass status is not warning",
			status:     preflight.StatusPass,
			severity:   preflight.SeverityMedium,
			wantWarning: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preflight.NewValidationResult(
				preflight.RequirementGPUSupport,
				tt.status,
				tt.severity,
				"actual",
				"expected",
				preflight.NewUserGuidance("", "", nil, ""),
			)

			assert.Equal(t, tt.wantWarning, result.IsWarning())
		})
	}
}

func TestValidationResult_IsPassing(t *testing.T) {
	tests := []struct {
		name        string
		status      preflight.ValidationStatus
		wantPassing bool
	}{
		{
			name:        "pass status is passing",
			status:      preflight.StatusPass,
			wantPassing: true,
		},
		{
			name:        "fail status is not passing",
			status:      preflight.StatusFail,
			wantPassing: false,
		},
		{
			name:        "warning status is not passing",
			status:      preflight.StatusWarning,
			wantPassing: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preflight.NewValidationResult(
				preflight.RequirementDiskSpace,
				tt.status,
				preflight.SeverityMedium,
				"actual",
				"expected",
				preflight.NewUserGuidance("", "", nil, ""),
			)

			assert.Equal(t, tt.wantPassing, result.IsPassing())
		})
	}
}

func TestValidationResult_FormatMessage(t *testing.T) {
	tests := []struct {
		name             string
		requirementName  preflight.RequirementName
		status           preflight.ValidationStatus
		guidanceMessage  string
		wantContains     []string
	}{
		{
			name:            "pass status shows checkmark",
			requirementName: preflight.RequirementDebianVersion,
			status:          preflight.StatusPass,
			guidanceMessage: "",
			wantContains:    []string{"✓", "debian_version", "Valid"},
		},
		{
			name:            "fail status shows cross",
			requirementName: preflight.RequirementDiskSpace,
			status:          preflight.StatusFail,
			guidanceMessage: "Insufficient disk space",
			wantContains:    []string{"✗", "disk_space", "Insufficient disk space"},
		},
		{
			name:            "warning status shows warning symbol",
			requirementName: preflight.RequirementGPUSupport,
			status:          preflight.StatusWarning,
			guidanceMessage: "NVIDIA requires configuration",
			wantContains:    []string{"⚠", "gpu_support", "NVIDIA requires configuration"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			guidance := preflight.NewUserGuidance(tt.guidanceMessage, "", nil, "")
			result := preflight.NewValidationResult(
				tt.requirementName,
				tt.status,
				preflight.SeverityMedium,
				"actual",
				"expected",
				guidance,
			)

			message := result.FormatMessage()
			for _, want := range tt.wantContains {
				assert.Contains(t, message, want)
			}
		})
	}
}

func TestValidationResult_BlockingAndWarningMutualExclusivity(t *testing.T) {
	// A result cannot be both blocking and warning

	t.Run("critical failure is blocking but not warning", func(t *testing.T) {
		result := preflight.NewValidationResult(
			preflight.RequirementDebianVersion,
			preflight.StatusFail,
			preflight.SeverityCritical,
			"bookworm",
			"sid",
			preflight.NewUserGuidance("", "", nil, ""),
		)

		assert.True(t, result.IsBlocking())
		assert.False(t, result.IsWarning())
		assert.False(t, result.IsPassing())
	})

	t.Run("medium failure is warning but not blocking", func(t *testing.T) {
		result := preflight.NewValidationResult(
			preflight.RequirementGPUSupport,
			preflight.StatusFail,
			preflight.SeverityMedium,
			"nvidia",
			"any",
			preflight.NewUserGuidance("", "", nil, ""),
		)

		assert.False(t, result.IsBlocking())
		assert.True(t, result.IsWarning())
		assert.False(t, result.IsPassing())
	})

	t.Run("pass is neither blocking nor warning", func(t *testing.T) {
		result := preflight.NewValidationResult(
			preflight.RequirementInternet,
			preflight.StatusPass,
			preflight.SeverityLow,
			"connected",
			"connected",
			preflight.NewUserGuidance("", "", nil, ""),
		)

		assert.False(t, result.IsBlocking())
		assert.False(t, result.IsWarning())
		assert.True(t, result.IsPassing())
	})
}

func TestValidationResult_RealWorldScenarios(t *testing.T) {
	t.Run("Debian Bookworm blocks installation", func(t *testing.T) {
		version, _ := preflight.NewDebianVersion("bookworm", "12")
		guidance := preflight.NewUserGuidance(
			"Debian Bookworm is not supported",
			"Hyprland requires newer packages",
			[]string{"Upgrade to Debian Sid or Trixie"},
			"https://gohan.sh/docs/supported-versions",
		)

		result := preflight.NewValidationResult(
			preflight.RequirementDebianVersion,
			preflight.StatusFail,
			preflight.SeverityCritical,
			version,
			"sid or trixie",
			guidance,
		)

		assert.True(t, result.IsBlocking())
		assert.Equal(t, version, result.ActualValue())
	})

	t.Run("NVIDIA GPU generates warning", func(t *testing.T) {
		gpu, _ := preflight.NewGPUType(preflight.GPUVendorNVIDIA, "RTX 4090", "10de:2684")
		guidance := preflight.NewUserGuidance(
			"NVIDIA GPU requires additional configuration",
			"Proprietary drivers and Wayland setup needed",
			[]string{"Install nvidia-driver", "Configure WLR_DRM_DEVICES"},
			"https://gohan.sh/docs/nvidia",
		)

		result := preflight.NewValidationResult(
			preflight.RequirementGPUSupport,
			preflight.StatusWarning,
			preflight.SeverityMedium,
			gpu,
			"any GPU",
			guidance,
		)

		assert.False(t, result.IsBlocking())
		assert.True(t, result.IsWarning())
		assert.Equal(t, gpu, result.ActualValue())
	})

	t.Run("Insufficient disk space blocks installation", func(t *testing.T) {
		diskSpace, _ := preflight.NewDiskSpace(5*preflight.GB, 100*preflight.GB, "/")
		guidance := preflight.NewUserGuidance(
			"Insufficient disk space",
			"Need 10GB, have 5GB",
			[]string{"Free up space", "Run apt clean"},
			"https://gohan.sh/docs/disk-space",
		)

		result := preflight.NewValidationResult(
			preflight.RequirementDiskSpace,
			preflight.StatusFail,
			preflight.SeverityHigh,
			diskSpace,
			"10GB",
			guidance,
		)

		assert.True(t, result.IsBlocking())
		assert.Equal(t, diskSpace, result.ActualValue())
	})
}
