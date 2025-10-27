# Phase 2: Package Definitions

## Overview

Phase 2 defines all Debian packages required for a complete Hyprland desktop environment. This includes comprehensive package definitions, installation profiles (minimal, recommended, full), and GPU-specific packages.

## Key Findings

### Debian Version Compatibility

**Critical Discovery**: Hyprland was removed from Debian 13 "Trixie" in 2025 due to versioning concerns and upstream lag.

- **Debian Sid**: Hyprland 0.41.2+ds-1.3 available
- **Debian Trixie**: Hyprland **removed** from the release

**Impact**: Gohan **requires Debian Sid** for full Hyprland support.

### Package Availability Summary

| Package | Debian Sid | Debian Trixie | Notes |
|---------|------------|---------------|-------|
| hyprland | ✅ Yes | ❌ No | Core compositor |
| xdg-desktop-portal-hyprland | ✅ Yes | ❌ No | Desktop portal |
| waybar | ✅ Yes | ✅ Yes | Status bar |
| fuzzel | ✅ Yes | ✅ Yes | Launcher (Wayland-native) |
| mako-notifier | ✅ Yes | ✅ Yes | Notifications |
| kitty | ✅ Yes | ✅ Yes | Terminal emulator (GPU-accelerated) |
| alacritty | ✅ Yes | ✅ Yes | Alternative terminal |
| foot | ✅ Yes | ✅ Yes | Lightweight terminal |
| swaylock | ✅ Yes | ✅ Yes | Screen locker |
| swayidle | ✅ Yes | ✅ Yes | Idle manager |
| swaybg | ✅ Yes | ✅ Yes | Wallpaper daemon |
| grim | ✅ Yes | ✅ Yes | Screenshots |
| slurp | ✅ Yes | ✅ Yes | Region selector |
| wl-clipboard | ✅ Yes | ✅ Yes | Clipboard utilities |
| hypridle | ❌ No | ❌ No | Use swayidle instead |
| hyprlock | ❌ No | ❌ No | Use swaylock instead |

## Package Organization

### Package Groups

Packages are organized into logical groups:

1. **Core** - Essential Hyprland components
2. **Essential** - Must-have desktop tools
3. **Utilities** - Additional tools (screenshots, clipboard, etc.)
4. **GPU** - Graphics driver packages
5. **Fonts** - Font packages
6. **Desktop** - Desktop integration (polkit, xdg-utils, Qt support)
7. **Development** - Development tools (future)

### Installation Profiles

Three curated profiles for different use cases:

#### 1. Minimal Profile (18 packages)

**Purpose**: Bare minimum for a functional Hyprland desktop

**Includes**:
- Core Hyprland (hyprland, xdg-desktop-portal-hyprland)
- Essential Wayland tools (waybar, fuzzel, mako-notifier, foot)
- Lock/idle (swaylock, swayidle)
- Wallpaper (swaybg)
- Desktop integration (polkit-gnome, xdg-utils, qt5-wayland)
- Basic utilities (grim, slurp, wl-clipboard)
- Essential fonts (fonts-jetbrains-mono, fonts-font-awesome)

**Use Case**: Minimal resource usage, experienced users, custom setups

#### 2. Recommended Profile (28+ packages)

**Purpose**: Complete desktop experience with common tools

**Includes**: Minimal +
- Clipboard history (cliphist)
- System controls (brightnessctl, playerctl, pavucontrol)
- Network/Bluetooth (network-manager-gnome, blueman)
- Power menu (wlogout)
- Qt6 support
- Additional fonts (fonts-noto, fonts-noto-color-emoji)
- File manager (nautilus)

**Use Case**: Default installation, most users

#### 3. Full Profile (30+ packages)

**Purpose**: Everything including alternatives and extras

**Includes**: Recommended +
- Alternative terminal (alacritty)
- Calculator (gnome-calculator)
- Mesa drivers (mesa-vulkan-drivers, libgl1-mesa-dri)

**Use Case**: Maximum features, power users, testing

### GPU Profiles

GPU-specific packages for different vendors:

#### NVIDIA Profile
```
nvidia-driver
nvidia-vulkan-icd
nvidia-settings
```

#### AMD Profile
```
xserver-xorg-video-amdgpu
firmware-amd-graphics
mesa-vulkan-drivers
libgl1-mesa-dri
```

#### Intel Profile
```
mesa-vulkan-drivers
libgl1-mesa-dri
```

## Design Decisions

### 1. Tool Alternatives

**Original Plan** → **Actual Implementation** (Reason)

- ~~hyprlock~~ → **swaylock** (hyprlock not in Debian repos)
- ~~hypridle~~ → **swayidle** (hypridle not in Debian repos)
- Default terminal: **kitty** (GPU-accelerated, feature-rich)
  - Alternative: **alacritty** (also GPU-accelerated)
  - Alternative: **foot** (lightweight, Wayland-native)

### 2. Wayland-Native Tools

All selected tools are Wayland-native:
- ✅ Fuzzel (Wayland launcher) instead of Rofi (X11/XWayland)
- ✅ Foot (Wayland terminal) instead of Kitty
- ✅ Grim/Slurp (Wayland screenshots) native
- ✅ Mako (Wayland notifications) native

### 3. Debian Sid Requirement

**Decision**: Gohan targets **Debian Sid only** for initial release

**Rationale**:
1. Hyprland unavailable in Trixie
2. Sid is "unstable" but rolling (like Arch)
3. Most Hyprland users comfortable with bleeding-edge
4. Future: Add Trixie support when/if Hyprland returns

## Code Structure

### Domain Layer

**`package_definitions.go`** - 800+ lines
- `PackageDefinition` struct - Individual package metadata
- `PackageGroup` enum - Package categorization
- `AllPackageDefinitions` - Complete package catalog (60+ packages)
- Helper functions:
  - `GetPackagesByComponent()` - Get packages for a component
  - `GetPackagesByGroup()` - Get packages in a group
  - `GetRequiredPackages()` - Get all required packages
  - `GetPackagesForDebianVersion()` - Filter by Debian version
  - `IsPackageAvailable()` - Check package availability

**`package_profiles.go`** - 300+ lines
- `InstallationProfile` struct - Named package sets
- `ProfileType` enum - Profile identifiers
- Profile generators:
  - `GetMinimalProfile()` - Minimal installation
  - `GetRecommendedProfile()` - Recommended installation
  - `GetFullProfile()` - Full installation
  - `GetGPUProfile(vendor)` - GPU-specific packages
- Utilities:
  - `CombineProfiles()` - Merge multiple profiles
  - `FilterPackagesForDebianVersion()` - Filter by Debian version
  - `GetUnavailablePackages()` - Find unavailable packages

### Test Coverage

**`package_definitions_test.go`** - 300+ lines
- Component filtering tests
- Group filtering tests
- Required packages validation
- Debian version availability tests
- Alternative packages verification
- Package structure validation

**`package_profiles_test.go`** - 300+ lines
- Profile content tests
- Profile hierarchy tests (minimal ⊂ recommended ⊂ full)
- GPU profile tests
- Profile combination tests
- Debian version filtering tests
- Package count validation

**Test Results**: All 30+ tests passing ✅

## Usage Examples

### Get Packages for Installation

```go
// Get recommended profile
profile := installation.GetRecommendedProfile()

// Add GPU support
gpuProfile := installation.GetGPUProfile("nvidia")
combined := installation.CombineProfiles(profile, gpuProfile)

// Filter for Debian Sid
filtered := installation.FilterPackagesForDebianVersion(combined, "sid")

// Install packages
for _, pkg := range filtered.Packages {
    // Install pkg using apt
}
```

### Check Package Availability

```go
if installation.IsPackageAvailable("hyprland", "sid") {
    // Safe to install on Debian Sid
}

if !installation.IsPackageAvailable("hyprland", "trixie") {
    // Cannot install on Trixie
}
```

### Get Required Packages Only

```go
required := installation.GetRequiredPackages()
// Install only essential packages
```

## Integration with Existing Code

### Updated Components

**`types.go`** - Component definitions updated:
- ~~`ComponentRofi`~~ → `ComponentFuzzel`
- ~~`ComponentKitty`~~ → `ComponentGhostty`
- Added: `ComponentHypridle`, `ComponentMako`, `ComponentSwaybg`

**`types_test.go`** - Tests updated:
- Fixed component references
- All tests passing

### Next Steps (Phase 3)

Now that package definitions are complete, Phase 3 will implement:

1. **Package Installation Service**
   - Use APT to install packages
   - Handle dependencies
   - Progress reporting

2. **Configuration Deployment**
   - Copy templates to user directories
   - Template variable substitution
   - Backup existing configs

3. **Post-Install Tasks**
   - Environment variables
   - Display manager integration
   - Systemd services

## Package Statistics

- **Total Packages Defined**: 60+
- **Minimal Profile**: 18 packages
- **Recommended Profile**: 28+ packages
- **Full Profile**: 30+ packages
- **GPU Packages (NVIDIA)**: 3 packages
- **GPU Packages (AMD)**: 4 packages
- **GPU Packages (Intel)**: 2 packages

## Debian Version Requirements

**Gohan Requires**: Debian Sid (Unstable)

**Why Not Trixie**:
- Hyprland removed from Trixie (2025)
- xdg-desktop-portal-hyprland unavailable
- hyprland-backgrounds unavailable

**Future Consideration**: If Hyprland returns to Trixie, Gohan can support both versions using the version filtering functions already implemented.

## Key Takeaways

1. ✅ **Complete Package Catalog**: 60+ packages defined with metadata
2. ✅ **Three Installation Profiles**: Minimal, Recommended, Full
3. ✅ **GPU Support**: NVIDIA, AMD, Intel profiles
4. ✅ **Debian Version Aware**: Filters packages by availability
5. ✅ **Wayland-Native Stack**: All tools native to Wayland
6. ✅ **Comprehensive Tests**: 30+ tests, 100% passing
7. ⚠️ **Debian Sid Only**: Hyprland unavailable in Trixie
8. ✅ **Production Ready**: Code ready for Phase 3 integration

## References

- [Debian Packages - Hyprland](https://packages.debian.org/sid/main/hyprland)
- [Debian Packages - Fuzzel](https://packages.debian.org/sid/fuzzel)
- [Debian Packages - Foot](https://packages.debian.org/sid/foot)
- [Hyprland Removed from Trixie](https://linuxiac.com/hyprland-will-not-be-included-in-debian-13-trixie/)
- [JaKooLit/Debian-Hyprland](https://github.com/JaKooLit/Debian-Hyprland) - Community installer reference
