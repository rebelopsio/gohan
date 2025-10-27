package installation

// InstallationProfile defines a set of packages for different installation types
type InstallationProfile struct {
	Name        string
	Description string
	Packages    []string
}

// ProfileType identifies different installation profile types
type ProfileType string

const (
	ProfileMinimal     ProfileType = "minimal"     // Bare minimum for functional Hyprland
	ProfileRecommended ProfileType = "recommended" // Recommended setup with common tools
	ProfileFull        ProfileType = "full"        // Complete setup with all features
)

// GetMinimalProfile returns the minimal installation profile
// This includes only the core components needed for a functional Hyprland desktop
func GetMinimalProfile() InstallationProfile {
	return InstallationProfile{
		Name:        "Minimal",
		Description: "Bare minimum for a functional Hyprland desktop",
		Packages: []string{
			// Core Hyprland
			"hyprland",
			"xdg-desktop-portal-hyprland",

			// Essential Wayland tools
			"waybar",
			"fuzzel",
			"mako-notifier",
			"kitty",         // Terminal
			"kitty-terminfo", // Terminal info

			// Lock and idle management
			"swaylock",
			"swayidle",

			// Wallpaper
			"swaybg",

			// Desktop integration
			"polkit-gnome",
			"xdg-utils",
			"qt5-wayland",

			// Basic utilities
			"grim",       // Screenshots
			"slurp",      // Region select
			"wl-clipboard", // Clipboard

			// Essential fonts
			"fonts-jetbrains-mono",
			"fonts-font-awesome",
		},
	}
}

// GetRecommendedProfile returns the recommended installation profile
// This includes minimal + commonly used tools and utilities
func GetRecommendedProfile() InstallationProfile {
	minimal := GetMinimalProfile()

	additionalPackages := []string{
		// Additional backgrounds
		"hyprland-backgrounds",

		// Clipboard history
		"cliphist",

		// System controls
		"brightnessctl", // Brightness
		"playerctl",     // Media control
		"pavucontrol",   // Audio control

		// Network and Bluetooth
		"network-manager-gnome",
		"blueman",

		// Power menu
		"wlogout",

		// Qt6 support
		"qt6-wayland",

		// Fonts
		"fonts-noto",
		"fonts-noto-color-emoji",

		// File manager
		"nautilus",
	}

	return InstallationProfile{
		Name:        "Recommended",
		Description: "Recommended setup with common tools and utilities",
		Packages:    append(minimal.Packages, additionalPackages...),
	}
}

// GetFullProfile returns the full installation profile
// This includes everything: recommended + optional tools + alternative apps
func GetFullProfile() InstallationProfile {
	recommended := GetRecommendedProfile()

	additionalPackages := []string{
		// Alternative terminal
		"alacritty",

		// Calculator
		"gnome-calculator",

		// Mesa drivers (for Intel/AMD)
		"mesa-vulkan-drivers",
		"libgl1-mesa-dri",
	}

	return InstallationProfile{
		Name:        "Full",
		Description: "Complete setup with all features and alternatives",
		Packages:    append(recommended.Packages, additionalPackages...),
	}
}

// GetGPUProfile returns packages for specific GPU vendor
func GetGPUProfile(vendor string) InstallationProfile {
	switch vendor {
	case "nvidia":
		return InstallationProfile{
			Name:        "NVIDIA GPU Support",
			Description: "NVIDIA proprietary driver and tools",
			Packages: []string{
				"nvidia-driver",
				"nvidia-vulkan-icd",
				"nvidia-settings",
			},
		}
	case "amd":
		return InstallationProfile{
			Name:        "AMD GPU Support",
			Description: "AMD open-source drivers and firmware",
			Packages: []string{
				"xserver-xorg-video-amdgpu",
				"firmware-amd-graphics",
				"mesa-vulkan-drivers",
				"libgl1-mesa-dri",
			},
		}
	case "intel":
		return InstallationProfile{
			Name:        "Intel GPU Support",
			Description: "Intel open-source drivers",
			Packages: []string{
				"mesa-vulkan-drivers",
				"libgl1-mesa-dri",
			},
		}
	default:
		return InstallationProfile{
			Name:        "Generic GPU Support",
			Description: "Generic Mesa drivers",
			Packages: []string{
				"mesa-vulkan-drivers",
				"libgl1-mesa-dri",
			},
		}
	}
}

// GetProfileByType returns an installation profile by its type
func GetProfileByType(profileType ProfileType) InstallationProfile {
	switch profileType {
	case ProfileMinimal:
		return GetMinimalProfile()
	case ProfileRecommended:
		return GetRecommendedProfile()
	case ProfileFull:
		return GetFullProfile()
	default:
		return GetRecommendedProfile()
	}
}

// GetAllProfiles returns all available installation profiles
func GetAllProfiles() []InstallationProfile {
	return []InstallationProfile{
		GetMinimalProfile(),
		GetRecommendedProfile(),
		GetFullProfile(),
	}
}

// CombineProfiles merges multiple profiles into one
func CombineProfiles(profiles ...InstallationProfile) InstallationProfile {
	packageSet := make(map[string]bool)
	var allPackages []string

	for _, profile := range profiles {
		for _, pkg := range profile.Packages {
			if !packageSet[pkg] {
				packageSet[pkg] = true
				allPackages = append(allPackages, pkg)
			}
		}
	}

	return InstallationProfile{
		Name:        "Combined Profile",
		Description: "Combination of multiple profiles",
		Packages:    allPackages,
	}
}

// FilterPackagesForDebianVersion filters profile packages by Debian version availability
func FilterPackagesForDebianVersion(profile InstallationProfile, debianVersion string) InstallationProfile {
	var availablePackages []string

	for _, pkg := range profile.Packages {
		if IsPackageAvailable(pkg, debianVersion) {
			availablePackages = append(availablePackages, pkg)
		}
	}

	return InstallationProfile{
		Name:        profile.Name + " (filtered for " + debianVersion + ")",
		Description: profile.Description,
		Packages:    availablePackages,
	}
}

// GetUnavailablePackages returns packages from profile not available in Debian version
func GetUnavailablePackages(profile InstallationProfile, debianVersion string) []string {
	var unavailable []string

	for _, pkg := range profile.Packages {
		if !IsPackageAvailable(pkg, debianVersion) {
			unavailable = append(unavailable, pkg)
		}
	}

	return unavailable
}
