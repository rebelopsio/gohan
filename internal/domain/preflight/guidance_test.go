package preflight_test

import (
	"strings"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
	"github.com/stretchr/testify/assert"
)

func TestNewUserGuidance(t *testing.T) {
	message := "Debian Bookworm is not supported"
	reason := "Hyprland requires newer package versions"
	steps := []string{"Upgrade to Debian Sid", "Or use Debian Trixie"}
	docURL := "https://gohan.sh/docs/supported-versions"

	guidance := preflight.NewUserGuidance(message, reason, steps, docURL)

	assert.Equal(t, message, guidance.Message())
	assert.Equal(t, reason, guidance.Reason())
	assert.Equal(t, steps, guidance.ActionableSteps())
	assert.Equal(t, docURL, guidance.DocumentationURL())
}

func TestUserGuidance_HasSteps(t *testing.T) {
	tests := []struct {
		name      string
		steps     []string
		wantSteps bool
	}{
		{
			name:      "has steps",
			steps:     []string{"Step 1", "Step 2"},
			wantSteps: true,
		},
		{
			name:      "no steps",
			steps:     []string{},
			wantSteps: false,
		},
		{
			name:      "nil steps",
			steps:     nil,
			wantSteps: false,
		},
		{
			name:      "single step",
			steps:     []string{"Only step"},
			wantSteps: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			guidance := preflight.NewUserGuidance(
				"Test message",
				"Test reason",
				tt.steps,
				"",
			)
			assert.Equal(t, tt.wantSteps, guidance.HasSteps())
		})
	}
}

func TestUserGuidance_Format(t *testing.T) {
	t.Run("complete guidance with all fields", func(t *testing.T) {
		guidance := preflight.NewUserGuidance(
			"Insufficient disk space",
			"Installation requires at least 10GB",
			[]string{
				"Free up space by removing unused packages",
				"Run 'sudo apt clean' to clear package cache",
				"Consider moving data to external storage",
			},
			"https://gohan.sh/docs/disk-requirements",
		)

		formatted := guidance.Format()

		// Check all components are present
		assert.Contains(t, formatted, "Insufficient disk space")
		assert.Contains(t, formatted, "Reason: Installation requires at least 10GB")
		assert.Contains(t, formatted, "How to fix:")
		assert.Contains(t, formatted, "1. Free up space by removing unused packages")
		assert.Contains(t, formatted, "2. Run 'sudo apt clean' to clear package cache")
		assert.Contains(t, formatted, "3. Consider moving data to external storage")
		assert.Contains(t, formatted, "Learn more: https://gohan.sh/docs/disk-requirements")
	})

	t.Run("guidance without reason", func(t *testing.T) {
		guidance := preflight.NewUserGuidance(
			"Problem occurred",
			"",
			[]string{"Fix step"},
			"",
		)

		formatted := guidance.Format()

		assert.Contains(t, formatted, "Problem occurred")
		assert.NotContains(t, formatted, "Reason:")
		assert.Contains(t, formatted, "How to fix:")
	})

	t.Run("guidance without steps", func(t *testing.T) {
		guidance := preflight.NewUserGuidance(
			"Information message",
			"Just FYI",
			[]string{},
			"",
		)

		formatted := guidance.Format()

		assert.Contains(t, formatted, "Information message")
		assert.Contains(t, formatted, "Reason: Just FYI")
		assert.NotContains(t, formatted, "How to fix:")
	})

	t.Run("guidance without documentation URL", func(t *testing.T) {
		guidance := preflight.NewUserGuidance(
			"Message",
			"Reason",
			[]string{"Step"},
			"",
		)

		formatted := guidance.Format()

		assert.Contains(t, formatted, "Message")
		assert.NotContains(t, formatted, "Learn more:")
	})

	t.Run("minimal guidance with only message", func(t *testing.T) {
		guidance := preflight.NewUserGuidance(
			"Simple message",
			"",
			[]string{},
			"",
		)

		formatted := guidance.Format()

		assert.Contains(t, formatted, "Simple message")
		assert.NotContains(t, formatted, "Reason:")
		assert.NotContains(t, formatted, "How to fix:")
		assert.NotContains(t, formatted, "Learn more:")
	})
}

func TestUserGuidance_FormatNumberedSteps(t *testing.T) {
	guidance := preflight.NewUserGuidance(
		"Test",
		"",
		[]string{
			"First step",
			"Second step",
			"Third step",
		},
		"",
	)

	formatted := guidance.Format()

	// Verify steps are numbered correctly
	assert.Contains(t, formatted, "1. First step")
	assert.Contains(t, formatted, "2. Second step")
	assert.Contains(t, formatted, "3. Third step")

	// Verify steps are in correct order
	firstIndex := strings.Index(formatted, "1. First step")
	secondIndex := strings.Index(formatted, "2. Second step")
	thirdIndex := strings.Index(formatted, "3. Third step")

	assert.Less(t, firstIndex, secondIndex)
	assert.Less(t, secondIndex, thirdIndex)
}

func TestUserGuidance_RealWorldExamples(t *testing.T) {
	t.Run("Debian version error", func(t *testing.T) {
		guidance := preflight.NewUserGuidance(
			"Debian Bookworm is not supported",
			"Hyprland requires newer library versions only available in Sid or Trixie",
			[]string{
				"Upgrade to Debian Sid: https://wiki.debian.org/DebianUnstable",
				"Or switch to Debian Trixie (testing)",
				"Backup your system before upgrading",
			},
			"https://gohan.sh/docs/supported-versions",
		)

		assert.Equal(t, "Debian Bookworm is not supported", guidance.Message())
		assert.True(t, guidance.HasSteps())
		assert.Equal(t, 3, len(guidance.ActionableSteps()))
	})

	t.Run("NVIDIA warning", func(t *testing.T) {
		guidance := preflight.NewUserGuidance(
			"NVIDIA GPU detected - additional configuration required",
			"NVIDIA GPUs require proprietary drivers and special Wayland configuration",
			[]string{
				"Install nvidia-driver from non-free repository",
				"Configure WLR_DRM_DEVICES environment variable",
				"Add required kernel parameters",
			},
			"https://gohan.sh/docs/nvidia-setup",
		)

		assert.Contains(t, guidance.Message(), "NVIDIA")
		assert.True(t, guidance.HasSteps())
		formatted := guidance.Format()
		assert.Contains(t, formatted, "nvidia-driver")
	})

	t.Run("Disk space error", func(t *testing.T) {
		guidance := preflight.NewUserGuidance(
			"Insufficient disk space available",
			"Installation requires at least 10GB free space, but only 5.2GB is available",
			[]string{
				"Remove unused packages: sudo apt autoremove",
				"Clean package cache: sudo apt clean",
				"Check for large files: du -sh /* | sort -h",
			},
			"https://gohan.sh/docs/disk-requirements",
		)

		assert.Contains(t, guidance.Reason(), "10GB")
		assert.Contains(t, guidance.Reason(), "5.2GB")
		assert.True(t, guidance.HasSteps())
	})
}
