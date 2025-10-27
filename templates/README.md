# Gohan Configuration Templates

## Overview

This directory contains opinionated, production-ready configuration templates for a complete Hyprland desktop environment on Debian. These configurations are adapted from [Omarchy](https://github.com/rebelopsio/omarchy) but tailored for **Debian Sid**.

**⚠️ IMPORTANT**: Gohan requires **Debian Sid (Unstable)** because Hyprland was removed from Debian 13 "Trixie" in 2025. See [Phase 2 Documentation](../docs/phase2-package-definitions.md) for details.

## What's Included

### 🪟 Hyprland (Compositor)
- **hyprland.conf** - Main configuration with modular sourcing
- **bindings.conf** - Comprehensive keybindings (150+ shortcuts)
- **looknfeel.conf** - Appearance, animations, decorations
- **input.conf** - Keyboard, mouse, touchpad, gestures
- **monitors.conf** - Display configuration template
- **autostart.conf** - Essential services and applications
- **hyprlock.conf** - Lock screen with modern UI
- **hypridle.conf** - Idle management and power saving

### 📊 Waybar (Status Bar)
- **config.jsonc** - Module configuration
- **style.css** - Catppuccin Mocha theme

### 💻 Kitty (Terminal)
- **kitty.conf** - GPU-accelerated terminal configuration with Catppuccin theme

### 🚀 Fuzzel (Application Launcher)
- **fuzzel.ini** - Wayland-native launcher with Catppuccin theme

## Design Philosophy

### Modular & Customizable
- Hyprland config is split into logical modules
- Easy to override individual parts
- User customizations don't conflict with defaults

### Sane Defaults
- Works out of the box
- Familiar keybindings (SUPER key as main modifier)
- Beautiful Catppuccin Mocha color scheme

### Debian-Native
- Uses standard Debian packages
- No Omarchy-specific scripts (yet)
- Standard Linux tools (grim, slurp, wl-clipboard)

## Keybinding Highlights

### Essential
- `SUPER + Return` - Terminal
- `SUPER + SPACE` - Application launcher
- `SUPER + W` / `SUPER + Q` - Close window
- `SUPER + F` - Fullscreen
- `SUPER + T` - Toggle floating
- `SUPER + L` - Lock screen
- `SUPER + ESCAPE` - Power menu

### Window Management
- `SUPER + Arrow Keys` - Move focus
- `SUPER + SHIFT + Arrow Keys` - Move window
- `SUPER + CTRL + Arrow Keys` - Resize window
- `SUPER + 1-9` - Switch workspace
- `SUPER + SHIFT + 1-9` - Move window to workspace

### Media & System
- Media keys work out of the box
- `PRINT` - Screenshot
- `SHIFT + PRINT` - Screenshot selection
- `CTRL + PRINT` - Screenshot to clipboard

## Dependencies

### Core Hyprland Stack
```bash
hyprland
hyprlock
hypridle
xdg-desktop-portal-hyprland
```

### Essential Tools
```bash
waybar              # Status bar
kitty               # Terminal (GPU-accelerated)
fuzzel              # Application launcher (Wayland-native)
mako                # Notifications
swaybg              # Wallpaper
```

### Utilities
```bash
grim slurp          # Screenshots
wl-clipboard        # Clipboard manager
cliphist            # Clipboard history
brightnessctl       # Brightness control
playerctl           # Media control
pavucontrol         # Audio control
```

### System Integration
```bash
polkit-gnome        # Authentication
network-manager-gnome  # Network management
blueman             # Bluetooth
```

### Fonts
```bash
fonts-jetbrains-mono  # Or any Nerd Font
fonts-font-awesome    # Icons
```

## Configuration Layers

Gohan uses a simple layered approach:

```
1. System Templates (/usr/share/gohan/templates/)
   ↓
2. User Config (~/.config/)
   - hypr/
   - waybar/
   - kitty/
   - fuzzel/
```

Users can override any setting by editing files in `~/.config/`.

## Color Scheme

All configurations use **Catppuccin Mocha** for consistency:

- Background: `#1e1e2e`
- Foreground: `#cdd6f4`
- Accent Blue: `#89b4fa`
- Accent Pink: `#f5c2e7`
- Accent Red: `#f38ba8`
- Accent Green: `#a6e3a1`

## What's Different from Omarchy?

### Omarchy ✨
- Custom scripts and helpers
- Walker launcher (custom)
- UWSM session manager
- Omarchy-specific branding
- Arch Linux focused

### Gohan 🥋
- Standard Linux tools
- Fuzzel launcher (Wayland-native)
- Kitty terminal (GPU-accelerated)
- Direct Hyprland startup
- Debian branding (future)
- Debian Sid focused

## Next Steps

The infrastructure is complete! What we still need:

### Phase 2: Complete Package Lists
- [ ] Define all dependencies for each component
- [ ] Create package groups (minimal, recommended, full)
- [ ] Handle Debian Sid vs Trixie differences

### Phase 3: Configuration Deployment
- [ ] Implement config file copying service
- [ ] Add backup mechanism
- [ ] Template variable substitution (username, paths)

### Phase 4: Post-Install Setup
- [ ] Environment variables configuration
- [ ] Display manager integration (GDM/SDDM)
- [ ] Systemd service management

### Phase 5: Polish
- [ ] Custom Gohan wallpapers
- [ ] Fuzzel integration scripts (power menu, etc.)
- [ ] First-run welcome screen

## Testing Locally

To test these configs on a running system:

```bash
# Backup existing configs
cp -r ~/.config/hypr ~/.config/hypr.backup

# Copy Gohan configs
cp -r templates/hyprland/* ~/.config/hypr/
cp -r templates/waybar/* ~/.config/waybar/
cp -r templates/kitty/* ~/.config/kitty/
cp -r templates/fuzzel/* ~/.config/fuzzel/

# Reload Hyprland
hyprctl reload
```

## Contributing

When modifying configs:
1. Keep them modular
2. Comment everything
3. Test on a fresh Debian Sid VM
4. Ensure no Omarchy-specific dependencies

## References

- [Hyprland Wiki](https://wiki.hyprland.org/)
- [Omarchy Repository](https://github.com/rebelopsio/omarchy)
- [Catppuccin Theme](https://github.com/catppuccin/catppuccin)
- [Waybar Wiki](https://github.com/Alexays/Waybar/wiki)
