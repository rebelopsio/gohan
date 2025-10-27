package installation_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPackagesByComponent(t *testing.T) {
	tests := []struct {
		name          string
		component     installation.ComponentName
		minExpected   int
		checkPackages []string
	}{
		{
			name:          "Hyprland component",
			component:     installation.ComponentHyprland,
			minExpected:   1,
			checkPackages: []string{"hyprland", "xdg-desktop-portal-hyprland"},
		},
		{
			name:          "Waybar component",
			component:     installation.ComponentWaybar,
			minExpected:   1,
			checkPackages: []string{"waybar"},
		},
		{
			name:          "Fuzzel component",
			component:     installation.ComponentFuzzel,
			minExpected:   1,
			checkPackages: []string{"fuzzel"},
		},
		{
			name:          "Terminal component (Kitty alternatives)",
			component:     installation.ComponentKitty,
			minExpected:   1,
			checkPackages: []string{"kitty", "kitty-terminfo"},
		},
		{
			name:        "AMD Driver component",
			component:   installation.ComponentAMDDriver,
			minExpected: 1,
			checkPackages: []string{
				"xserver-xorg-video-amdgpu",
				"firmware-amd-graphics",
			},
		},
		{
			name:        "NVIDIA Driver component",
			component:   installation.ComponentNVIDIADriver,
			minExpected: 1,
			checkPackages: []string{
				"nvidia-driver",
				"nvidia-vulkan-icd",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packages := installation.GetPackagesByComponent(tt.component)

			assert.GreaterOrEqual(t, len(packages), tt.minExpected,
				"Expected at least %d packages for component %s", tt.minExpected, tt.component)

			// Check that expected packages are present
			packageNames := make(map[string]bool)
			for _, pkg := range packages {
				packageNames[pkg.Name] = true
			}

			for _, expectedPkg := range tt.checkPackages {
				assert.True(t, packageNames[expectedPkg],
					"Expected package %s for component %s", expectedPkg, tt.component)
			}
		})
	}
}

func TestGetPackagesByGroup(t *testing.T) {
	tests := []struct {
		name         string
		group        installation.PackageGroup
		minExpected  int
		shouldInclude []string
	}{
		{
			name:        "Core packages",
			group:       installation.GroupCore,
			minExpected: 1,
			shouldInclude: []string{
				"hyprland",
				"xdg-desktop-portal-hyprland",
			},
		},
		{
			name:        "Essential packages",
			group:       installation.GroupEssential,
			minExpected: 5,
			shouldInclude: []string{
				"waybar",
				"fuzzel",
				"mako-notifier",
				"foot",
			},
		},
		{
			name:        "Utilities",
			group:       installation.GroupUtilities,
			minExpected: 3,
			shouldInclude: []string{
				"grim",
				"slurp",
				"wl-clipboard",
			},
		},
		{
			name:        "GPU drivers",
			group:       installation.GroupGPU,
			minExpected: 3,
			shouldInclude: []string{
				"mesa-vulkan-drivers",
				"nvidia-driver",
			},
		},
		{
			name:        "Fonts",
			group:       installation.GroupFonts,
			minExpected: 2,
			shouldInclude: []string{
				"fonts-jetbrains-mono",
				"fonts-font-awesome",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packages := installation.GetPackagesByGroup(tt.group)

			assert.GreaterOrEqual(t, len(packages), tt.minExpected,
				"Expected at least %d packages in group %s", tt.minExpected, tt.group)

			packageNames := make(map[string]bool)
			for _, pkg := range packages {
				packageNames[pkg.Name] = true
			}

			for _, expectedPkg := range tt.shouldInclude {
				assert.True(t, packageNames[expectedPkg],
					"Expected package %s in group %s", expectedPkg, tt.group)
			}
		})
	}
}

func TestGetRequiredPackages(t *testing.T) {
	packages := installation.GetRequiredPackages()

	assert.NotEmpty(t, packages, "Should have required packages")

	// All required packages should be marked as required
	for _, pkg := range packages {
		assert.True(t, pkg.Required,
			"Package %s in required list but Required=false", pkg.Name)
	}

	// Check for essential required packages
	requiredNames := make(map[string]bool)
	for _, pkg := range packages {
		requiredNames[pkg.Name] = true
	}

	essentialPackages := []string{
		"hyprland",
		"waybar",
		"fuzzel",
		"kitty",
	}

	for _, essential := range essentialPackages {
		assert.True(t, requiredNames[essential],
			"Essential package %s should be in required list", essential)
	}
}

func TestGetPackagesForDebianVersion(t *testing.T) {
	tests := []struct {
		name               string
		version            string
		shouldInclude      []string
		shouldExclude      []string
		minExpectedPackages int
	}{
		{
			name:    "Debian Sid",
			version: "sid",
			shouldInclude: []string{
				"hyprland",
				"waybar",
				"fuzzel",
				"foot",
			},
			shouldExclude: []string{}, // Most packages available in Sid
			minExpectedPackages: 20,
		},
		{
			name:    "Debian Trixie",
			version: "trixie",
			shouldInclude: []string{
				"waybar",
				"fuzzel",
				"foot",
				"grim",
				"slurp",
			},
			shouldExclude: []string{
				"hyprland", // Removed from Trixie
				"xdg-desktop-portal-hyprland",
			},
			minExpectedPackages: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packages := installation.GetPackagesForDebianVersion(tt.version)

			assert.GreaterOrEqual(t, len(packages), tt.minExpectedPackages,
				"Expected at least %d packages for %s", tt.minExpectedPackages, tt.version)

			packageNames := make(map[string]bool)
			for _, pkg := range packages {
				packageNames[pkg.Name] = true
			}

			// Check inclusions
			for _, expectedPkg := range tt.shouldInclude {
				assert.True(t, packageNames[expectedPkg],
					"Package %s should be available in %s", expectedPkg, tt.version)
			}

			// Check exclusions
			for _, excludedPkg := range tt.shouldExclude {
				assert.False(t, packageNames[excludedPkg],
					"Package %s should NOT be available in %s", excludedPkg, tt.version)
			}
		})
	}
}

func TestIsPackageAvailable(t *testing.T) {
	tests := []struct {
		name          string
		packageName   string
		debianVersion string
		expected      bool
	}{
		{
			name:          "Hyprland available in Sid",
			packageName:   "hyprland",
			debianVersion: "sid",
			expected:      true,
		},
		{
			name:          "Hyprland NOT available in Trixie",
			packageName:   "hyprland",
			debianVersion: "trixie",
			expected:      false,
		},
		{
			name:          "Fuzzel available in Sid",
			packageName:   "fuzzel",
			debianVersion: "sid",
			expected:      true,
		},
		{
			name:          "Fuzzel available in Trixie",
			packageName:   "fuzzel",
			debianVersion: "trixie",
			expected:      true,
		},
		{
			name:          "Waybar available in both",
			packageName:   "waybar",
			debianVersion: "sid",
			expected:      true,
		},
		{
			name:          "Non-existent package",
			packageName:   "nonexistent-package",
			debianVersion: "sid",
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := installation.IsPackageAvailable(tt.packageName, tt.debianVersion)
			assert.Equal(t, tt.expected, result,
				"Package %s availability in %s", tt.packageName, tt.debianVersion)
		})
	}
}

func TestPackageDefinitionStructure(t *testing.T) {
	t.Run("All packages have required fields", func(t *testing.T) {
		for _, pkg := range installation.AllPackageDefinitions {
			assert.NotEmpty(t, pkg.Name, "Package must have a name")
			assert.NotEmpty(t, pkg.Group, "Package must have a group")
			assert.NotEmpty(t, pkg.Description, "Package must have a description")
		}
	})

	t.Run("Required packages are in appropriate groups", func(t *testing.T) {
		for _, pkg := range installation.AllPackageDefinitions {
			if pkg.Required {
				// Required packages can be in core, essential, desktop, utilities, or fonts groups
				// Utilities includes essential tools like screenshot tools
				// Fonts includes fonts needed for proper display
				assert.True(t,
					pkg.Group == installation.GroupCore ||
						pkg.Group == installation.GroupEssential ||
						pkg.Group == installation.GroupDesktop ||
						pkg.Group == installation.GroupUtilities ||
						pkg.Group == installation.GroupFonts,
					"Required package %s should be in core, essential, desktop, utilities, or fonts group, got %s",
					pkg.Name, pkg.Group)
			}
		}
	})

	t.Run("GPU driver packages are in GPU group", func(t *testing.T) {
		gpuPackages := installation.GetPackagesByGroup(installation.GroupGPU)
		for _, pkg := range gpuPackages {
			assert.True(t,
				pkg.Component == installation.ComponentAMDDriver ||
					pkg.Component == installation.ComponentNVIDIADriver ||
					pkg.Component == installation.ComponentIntelDriver ||
					pkg.Component == "",
				"GPU group package should have a GPU component or no component")
		}
	})

	t.Run("Packages unavailable in Trixie have Sid alternative", func(t *testing.T) {
		for _, pkg := range installation.AllPackageDefinitions {
			if !pkg.DebianTrixie && pkg.Required {
				assert.True(t, pkg.DebianSid,
					"Required package %s not in Trixie should be in Sid", pkg.Name)
			}
		}
	})
}

func TestAlternativePackages(t *testing.T) {
	t.Run("Terminal alternatives exist", func(t *testing.T) {
		terminals := installation.GetPackagesByComponent(installation.ComponentKitty)
		require.NotEmpty(t, terminals, "Should have terminal packages")

		// At least one should be available
		found := false
		for _, pkg := range terminals {
			if pkg.DebianSid {
				found = true
				break
			}
		}
		assert.True(t, found, "At least one terminal should be available in Sid")

		// Kitty should be in the list
		foundKitty := false
		for _, pkg := range terminals {
			if pkg.Name == "kitty" {
				foundKitty = true
				assert.True(t, pkg.DebianSid, "kitty should be in Sid")
			}
		}
		assert.True(t, foundKitty, "kitty should be available as terminal option")
	})

	t.Run("Screen locker alternatives", func(t *testing.T) {
		lockers := installation.GetPackagesByComponent(installation.ComponentHyprlock)
		require.NotEmpty(t, lockers, "Should have screen locker packages")

		// swaylock should be available as alternative
		foundSwaylock := false
		for _, pkg := range lockers {
			if pkg.Name == "swaylock" {
				foundSwaylock = true
				assert.True(t, pkg.DebianSid, "swaylock should be in Sid")
			}
		}
		assert.True(t, foundSwaylock, "swaylock should be available as alternative")
	})
}
