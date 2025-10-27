package installation_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/stretchr/testify/assert"
)

func TestGetMinimalProfile(t *testing.T) {
	profile := installation.GetMinimalProfile()

	assert.Equal(t, "Minimal", profile.Name)
	assert.NotEmpty(t, profile.Description)
	assert.NotEmpty(t, profile.Packages)

	// Check for essential packages
	packageSet := toSet(profile.Packages)

	essentialPackages := []string{
		"hyprland",
		"waybar",
		"fuzzel",
		"kitty",
		"grim",
		"slurp",
		"wl-clipboard",
	}

	for _, pkg := range essentialPackages {
		assert.True(t, packageSet[pkg],
			"Minimal profile should include %s", pkg)
	}

	// Should not include optional packages
	optionalPackages := []string{
		"alacritty",
		"gnome-calculator",
		"hyprland-backgrounds",
	}

	for _, pkg := range optionalPackages {
		assert.False(t, packageSet[pkg],
			"Minimal profile should not include optional package %s", pkg)
	}
}

func TestGetRecommendedProfile(t *testing.T) {
	minimal := installation.GetMinimalProfile()
	recommended := installation.GetRecommendedProfile()

	assert.Equal(t, "Recommended", recommended.Name)
	assert.NotEmpty(t, recommended.Description)
	assert.Greater(t, len(recommended.Packages), len(minimal.Packages),
		"Recommended should have more packages than minimal")

	// All minimal packages should be in recommended
	minimalSet := toSet(minimal.Packages)
	recommendedSet := toSet(recommended.Packages)

	for _, pkg := range minimal.Packages {
		assert.True(t, recommendedSet[pkg],
			"Recommended should include all minimal packages, missing %s", pkg)
	}

	// Should include additional utilities
	additionalPackages := []string{
		"cliphist",
		"brightnessctl",
		"playerctl",
		"pavucontrol",
	}

	for _, pkg := range additionalPackages {
		assert.True(t, recommendedSet[pkg],
			"Recommended should include %s", pkg)
		assert.False(t, minimalSet[pkg],
			"Package %s should not be in minimal", pkg)
	}
}

func TestGetFullProfile(t *testing.T) {
	recommended := installation.GetRecommendedProfile()
	full := installation.GetFullProfile()

	assert.Equal(t, "Full", full.Name)
	assert.NotEmpty(t, full.Description)
	assert.Greater(t, len(full.Packages), len(recommended.Packages),
		"Full should have more packages than recommended")

	// All recommended packages should be in full
	recommendedSet := toSet(recommended.Packages)
	fullSet := toSet(full.Packages)

	for _, pkg := range recommended.Packages {
		assert.True(t, fullSet[pkg],
			"Full should include all recommended packages, missing %s", pkg)
	}

	// Should include alternatives
	additionalPackages := []string{
		"alacritty", // Alternative terminal
	}

	for _, pkg := range additionalPackages {
		assert.True(t, fullSet[pkg],
			"Full should include %s", pkg)
		assert.False(t, recommendedSet[pkg],
			"Package %s should not be in recommended", pkg)
	}
}

func TestGetGPUProfile(t *testing.T) {
	tests := []struct {
		name             string
		vendor           string
		expectedName     string
		expectedPackages []string
		minPackages      int
	}{
		{
			name:         "NVIDIA profile",
			vendor:       "nvidia",
			expectedName: "NVIDIA GPU Support",
			expectedPackages: []string{
				"nvidia-driver",
				"nvidia-vulkan-icd",
				"nvidia-settings",
			},
			minPackages: 3,
		},
		{
			name:         "AMD profile",
			vendor:       "amd",
			expectedName: "AMD GPU Support",
			expectedPackages: []string{
				"xserver-xorg-video-amdgpu",
				"firmware-amd-graphics",
				"mesa-vulkan-drivers",
			},
			minPackages: 4,
		},
		{
			name:         "Intel profile",
			vendor:       "intel",
			expectedName: "Intel GPU Support",
			expectedPackages: []string{
				"mesa-vulkan-drivers",
				"libgl1-mesa-dri",
			},
			minPackages: 2,
		},
		{
			name:         "Generic profile",
			vendor:       "unknown",
			expectedName: "Generic GPU Support",
			expectedPackages: []string{
				"mesa-vulkan-drivers",
			},
			minPackages: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := installation.GetGPUProfile(tt.vendor)

			assert.Equal(t, tt.expectedName, profile.Name)
			assert.NotEmpty(t, profile.Description)
			assert.GreaterOrEqual(t, len(profile.Packages), tt.minPackages)

			packageSet := toSet(profile.Packages)
			for _, expectedPkg := range tt.expectedPackages {
				assert.True(t, packageSet[expectedPkg],
					"%s profile should include %s", tt.vendor, expectedPkg)
			}
		})
	}
}

func TestGetProfileByType(t *testing.T) {
	tests := []struct {
		name         string
		profileType  installation.ProfileType
		expectedName string
	}{
		{
			name:         "Get minimal profile",
			profileType:  installation.ProfileMinimal,
			expectedName: "Minimal",
		},
		{
			name:         "Get recommended profile",
			profileType:  installation.ProfileRecommended,
			expectedName: "Recommended",
		},
		{
			name:         "Get full profile",
			profileType:  installation.ProfileFull,
			expectedName: "Full",
		},
		{
			name:         "Default to recommended",
			profileType:  installation.ProfileType("invalid"),
			expectedName: "Recommended",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := installation.GetProfileByType(tt.profileType)
			assert.Equal(t, tt.expectedName, profile.Name)
		})
	}
}

func TestGetAllProfiles(t *testing.T) {
	profiles := installation.GetAllProfiles()

	assert.Len(t, profiles, 3, "Should have exactly 3 profiles")

	names := make(map[string]bool)
	for _, profile := range profiles {
		names[profile.Name] = true
	}

	assert.True(t, names["Minimal"], "Should include Minimal profile")
	assert.True(t, names["Recommended"], "Should include Recommended profile")
	assert.True(t, names["Full"], "Should include Full profile")
}

func TestCombineProfiles(t *testing.T) {
	t.Run("Combine minimal and GPU profile", func(t *testing.T) {
		minimal := installation.GetMinimalProfile()
		nvidiaGPU := installation.GetGPUProfile("nvidia")

		combined := installation.CombineProfiles(minimal, nvidiaGPU)

		assert.NotEmpty(t, combined.Packages)

		// Should have packages from both profiles
		packageSet := toSet(combined.Packages)

		// Check minimal packages
		assert.True(t, packageSet["hyprland"])
		assert.True(t, packageSet["waybar"])

		// Check NVIDIA packages
		assert.True(t, packageSet["nvidia-driver"])
		assert.True(t, packageSet["nvidia-vulkan-icd"])
	})

	t.Run("No duplicate packages", func(t *testing.T) {
		// Create two profiles with overlapping packages
		profile1 := installation.InstallationProfile{
			Packages: []string{"pkg1", "pkg2", "pkg3"},
		}
		profile2 := installation.InstallationProfile{
			Packages: []string{"pkg2", "pkg3", "pkg4"},
		}

		combined := installation.CombineProfiles(profile1, profile2)

		// Should have unique packages only
		packageSet := toSet(combined.Packages)
		assert.Len(t, packageSet, 4, "Should have 4 unique packages")
		assert.True(t, packageSet["pkg1"])
		assert.True(t, packageSet["pkg2"])
		assert.True(t, packageSet["pkg3"])
		assert.True(t, packageSet["pkg4"])
	})
}

func TestFilterPackagesForDebianVersion(t *testing.T) {
	t.Run("Filter for Debian Sid", func(t *testing.T) {
		minimal := installation.GetMinimalProfile()
		filtered := installation.FilterPackagesForDebianVersion(minimal, "sid")

		// All packages in minimal should be available in Sid
		assert.Equal(t, len(minimal.Packages), len(filtered.Packages),
			"All minimal packages should be available in Sid")

		assert.Contains(t, filtered.Name, "sid")
	})

	t.Run("Filter for Debian Trixie", func(t *testing.T) {
		minimal := installation.GetMinimalProfile()
		filtered := installation.FilterPackagesForDebianVersion(minimal, "trixie")

		// hyprland should be filtered out for Trixie
		packageSet := toSet(filtered.Packages)
		assert.False(t, packageSet["hyprland"],
			"hyprland should not be in Trixie filtered profile")

		// But other packages should remain
		assert.True(t, packageSet["waybar"])
		assert.True(t, packageSet["fuzzel"])
	})
}

func TestGetUnavailablePackages(t *testing.T) {
	t.Run("No unavailable packages in Sid", func(t *testing.T) {
		minimal := installation.GetMinimalProfile()
		unavailable := installation.GetUnavailablePackages(minimal, "sid")

		assert.Empty(t, unavailable,
			"All minimal packages should be available in Sid")
	})

	t.Run("Hyprland unavailable in Trixie", func(t *testing.T) {
		minimal := installation.GetMinimalProfile()
		unavailable := installation.GetUnavailablePackages(minimal, "trixie")

		assert.NotEmpty(t, unavailable,
			"Some packages should be unavailable in Trixie")

		unavailableSet := toSet(unavailable)
		assert.True(t, unavailableSet["hyprland"],
			"hyprland should be in unavailable list for Trixie")
	})
}

func TestProfilePackageCounts(t *testing.T) {
	minimal := installation.GetMinimalProfile()
	recommended := installation.GetRecommendedProfile()
	full := installation.GetFullProfile()

	assert.GreaterOrEqual(t, len(minimal.Packages), 15,
		"Minimal should have at least 15 packages")
	assert.GreaterOrEqual(t, len(recommended.Packages), 20,
		"Recommended should have at least 20 packages")
	assert.GreaterOrEqual(t, len(full.Packages), 25,
		"Full should have at least 25 packages")

	assert.Less(t, len(minimal.Packages), len(recommended.Packages),
		"Minimal < Recommended")
	assert.Less(t, len(recommended.Packages), len(full.Packages),
		"Recommended < Full")
}

// Helper function to convert slice to set
func toSet(slice []string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range slice {
		set[item] = true
	}
	return set
}
