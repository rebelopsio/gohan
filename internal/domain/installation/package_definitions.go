package installation

// PackageDefinition represents a Debian package with its metadata
type PackageDefinition struct {
	Name         string
	Component    ComponentName
	Group        PackageGroup
	DebianSid    bool // Available in Debian Sid
	DebianTrixie bool // Available in Debian Trixie
	Required     bool // Required for minimal installation
	Description  string
	Alternatives []string // Alternative package names
}

// PackageGroup categorizes packages by their role
type PackageGroup string

const (
	GroupCore       PackageGroup = "core"        // Core Hyprland components
	GroupEssential  PackageGroup = "essential"   // Essential desktop tools
	GroupUtilities  PackageGroup = "utilities"   // Additional utilities
	GroupGPU        PackageGroup = "gpu"         // GPU drivers
	GroupFonts      PackageGroup = "fonts"       // Font packages
	GroupDesktop    PackageGroup = "desktop"     // Desktop integration
	GroupDevelopment PackageGroup = "development" // Development tools
)

// AllPackageDefinitions returns all available package definitions for Debian
var AllPackageDefinitions = []PackageDefinition{
	// ========================================================================
	// CORE HYPRLAND STACK
	// ========================================================================
	{
		Name:         "hyprland",
		Component:    ComponentHyprland,
		Group:        GroupCore,
		DebianSid:    true,
		DebianTrixie: false, // Removed from Trixie in 2025
		Required:     true,
		Description:  "Dynamic tiling Wayland compositor",
	},
	{
		Name:         "xdg-desktop-portal-hyprland",
		Component:    ComponentHyprland,
		Group:        GroupCore,
		DebianSid:    true,
		DebianTrixie: false,
		Required:     true,
		Description:  "xdg-desktop-portal backend for Hyprland",
	},
	{
		Name:         "hyprland-backgrounds",
		Component:    ComponentHyprland,
		Group:        GroupCore,
		DebianSid:    true,
		DebianTrixie: false,
		Required:     false,
		Description:  "Default backgrounds for Hyprland",
	},

	// ========================================================================
	// ESSENTIAL WAYLAND TOOLS
	// ========================================================================
	{
		Name:         "waybar",
		Component:    ComponentWaybar,
		Group:        GroupEssential,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Highly customizable Wayland bar for Sway and Wlroots based compositors",
	},
	{
		Name:         "fuzzel",
		Component:    ComponentFuzzel,
		Group:        GroupEssential,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Wayland-native application launcher",
	},
	{
		Name:         "mako-notifier",
		Component:    ComponentMako,
		Group:        GroupEssential,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Lightweight Wayland notification daemon",
	},
	{
		Name:         "swaybg",
		Component:    ComponentSwaybg,
		Group:        GroupEssential,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Wallpaper tool for Wayland compositors",
		Alternatives: []string{"hyprpaper"},
	},
	{
		Name:         "swaylock",
		Component:    ComponentHyprlock,
		Group:        GroupEssential,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Screen locker for Wayland",
		Alternatives: []string{"hyprlock"},
	},
	{
		Name:         "swayidle",
		Component:    ComponentHypridle,
		Group:        GroupEssential,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Idle management daemon for Wayland",
		Alternatives: []string{"hypridle"},
	},

	// ========================================================================
	// TERMINAL EMULATORS
	// ========================================================================
	{
		Name:         "kitty",
		Component:    ComponentKitty,
		Group:        GroupEssential,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Fast, GPU-based terminal emulator",
		Alternatives: []string{"alacritty", "foot"},
	},
	{
		Name:         "kitty-terminfo",
		Component:    ComponentKitty,
		Group:        GroupEssential,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Terminfo configuration for kitty",
	},
	{
		Name:         "alacritty",
		Component:    ComponentKitty,
		Group:        GroupEssential,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "GPU-accelerated terminal emulator",
		Alternatives: []string{"kitty", "foot"},
	},
	{
		Name:         "foot",
		Component:    ComponentKitty,
		Group:        GroupEssential,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Fast, lightweight Wayland terminal emulator",
		Alternatives: []string{"kitty", "alacritty"},
	},

	// ========================================================================
	// UTILITIES
	// ========================================================================
	{
		Name:         "grim",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Screenshot utility for Wayland",
	},
	{
		Name:         "slurp",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Region selector for Wayland",
	},
	{
		Name:         "wl-clipboard",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Command-line copy/paste utilities for Wayland",
	},
	{
		Name:         "cliphist",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Clipboard history manager for Wayland",
	},
	{
		Name:         "brightnessctl",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Screen and keyboard backlight control",
	},
	{
		Name:         "playerctl",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Media player controller for MPRIS",
	},
	{
		Name:         "pavucontrol",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "PulseAudio volume control",
	},
	{
		Name:         "network-manager-gnome",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Network connection manager (GUI)",
	},
	{
		Name:         "blueman",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "GTK+ Bluetooth manager",
	},
	{
		Name:         "wlogout",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Wayland-based logout menu",
	},
	{
		Name:         "hyprpicker",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    false, // May need to check availability
		DebianTrixie: false,
		Required:     false,
		Description:  "Color picker for Hyprland",
	},

	// ========================================================================
	// DESKTOP INTEGRATION
	// ========================================================================
	{
		Name:         "polkit-gnome",
		Component:    "",
		Group:        GroupDesktop,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "PolicyKit authentication agent",
		Alternatives: []string{"polkit-kde-agent-1"},
	},
	{
		Name:         "xdg-utils",
		Component:    "",
		Group:        GroupDesktop,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Desktop integration utilities",
	},
	{
		Name:         "qt5-wayland",
		Component:    "",
		Group:        GroupDesktop,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Qt5 Wayland platform plugin",
		Alternatives: []string{"qt6-wayland"},
	},
	{
		Name:         "qt6-wayland",
		Component:    "",
		Group:        GroupDesktop,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Qt6 Wayland platform plugin",
	},

	// ========================================================================
	// GPU DRIVERS
	// ========================================================================
	{
		Name:         "mesa-vulkan-drivers",
		Component:    ComponentIntelDriver,
		Group:        GroupGPU,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Mesa Vulkan graphics drivers (Intel, AMD)",
	},
	{
		Name:         "libgl1-mesa-dri",
		Component:    ComponentIntelDriver,
		Group:        GroupGPU,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Free Mesa DRI implementation (Intel, AMD)",
	},
	{
		Name:         "xserver-xorg-video-amdgpu",
		Component:    ComponentAMDDriver,
		Group:        GroupGPU,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "X.Org X server -- AMDGPU display driver",
	},
	{
		Name:         "firmware-amd-graphics",
		Component:    ComponentAMDDriver,
		Group:        GroupGPU,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Binary firmware for AMD graphics cards",
	},
	{
		Name:         "nvidia-driver",
		Component:    ComponentNVIDIADriver,
		Group:        GroupGPU,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "NVIDIA metapackage for proprietary driver",
	},
	{
		Name:         "nvidia-vulkan-icd",
		Component:    ComponentNVIDIADriver,
		Group:        GroupGPU,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "NVIDIA Vulkan installable client driver",
	},
	{
		Name:         "nvidia-settings",
		Component:    ComponentNVIDIADriver,
		Group:        GroupGPU,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Tool for configuring the NVIDIA graphics driver",
	},

	// ========================================================================
	// FONTS
	// ========================================================================
	{
		Name:         "fonts-noto",
		Component:    "",
		Group:        GroupFonts,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Noto font family with support for many languages",
	},
	{
		Name:         "fonts-noto-color-emoji",
		Component:    "",
		Group:        GroupFonts,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "Color emoji font from Google",
	},
	{
		Name:         "fonts-font-awesome",
		Component:    "",
		Group:        GroupFonts,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Iconic font designed for use with Bootstrap",
	},
	{
		Name:         "fonts-jetbrains-mono",
		Component:    "",
		Group:        GroupFonts,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     true,
		Description:  "Typeface for developers from JetBrains",
	},

	// ========================================================================
	// FILE MANAGERS & APPLICATIONS
	// ========================================================================
	{
		Name:         "nautilus",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "GNOME file manager",
		Alternatives: []string{"thunar", "nemo", "dolphin"},
	},
	{
		Name:         "gnome-calculator",
		Component:    "",
		Group:        GroupUtilities,
		DebianSid:    true,
		DebianTrixie: true,
		Required:     false,
		Description:  "GNOME desktop calculator",
	},
}

// GetPackagesByComponent returns all package definitions for a component
func GetPackagesByComponent(component ComponentName) []PackageDefinition {
	var packages []PackageDefinition
	for _, pkg := range AllPackageDefinitions {
		if pkg.Component == component {
			packages = append(packages, pkg)
		}
	}
	return packages
}

// GetPackagesByGroup returns all package definitions in a group
func GetPackagesByGroup(group PackageGroup) []PackageDefinition {
	var packages []PackageDefinition
	for _, pkg := range AllPackageDefinitions {
		if pkg.Group == group {
			packages = append(packages, pkg)
		}
	}
	return packages
}

// GetRequiredPackages returns all required package definitions
func GetRequiredPackages() []PackageDefinition {
	var packages []PackageDefinition
	for _, pkg := range AllPackageDefinitions {
		if pkg.Required {
			packages = append(packages, pkg)
		}
	}
	return packages
}

// GetPackagesForDebianVersion returns packages available for a Debian version
func GetPackagesForDebianVersion(version string) []PackageDefinition {
	var packages []PackageDefinition
	for _, pkg := range AllPackageDefinitions {
		if version == "sid" && pkg.DebianSid {
			packages = append(packages, pkg)
		} else if version == "trixie" && pkg.DebianTrixie {
			packages = append(packages, pkg)
		}
	}
	return packages
}

// IsPackageAvailable checks if a package is available for a Debian version
func IsPackageAvailable(packageName, debianVersion string) bool {
	for _, pkg := range AllPackageDefinitions {
		if pkg.Name == packageName {
			if debianVersion == "sid" {
				return pkg.DebianSid
			} else if debianVersion == "trixie" {
				return pkg.DebianTrixie
			}
		}
	}
	return false
}
