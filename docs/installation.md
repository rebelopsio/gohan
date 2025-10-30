# Installation Guide

This guide will walk you through installing Gohan and setting up your Hyprland environment.

## Prerequisites

!!! info "System Requirements"
    - **OS**: Debian Sid (unstable) or Trixie (testing) **only**
    - **Architecture**: x86_64 or ARM64
    - **Go**: Version 1.21 or higher (for building from source)
    - **Git**: For cloning the repository

!!! warning "Ubuntu and Debian Stable Not Supported"
    - **Ubuntu**: Not supported. Use official Hyprland installation methods for Ubuntu.
    - **Debian Bookworm (stable)**: Not supported. Hyprland requires newer packages only available in Sid or Trixie.
    - **Reason**: Gohan targets cutting-edge Hyprland features that require the latest packages from Debian's rolling branches.

!!! note "Optional Dependencies"
    These components are optional but enable full theme management features:

    - **Hyprland**: For desktop environment features
    - **Waybar**: For status bar theming
    - **Kitty**: For terminal theming
    - **Rofi/Fuzzel**: For launcher theming
    - **Mako**: For notification theming
    - **Alacritty**: Alternative terminal theming

## Installation Methods

=== "Build from Source"

    !!! tip "Recommended for developers"
        Building from source ensures you have the latest code and can customize the build.

    ```bash
    # Clone the repository
    git clone https://github.com/rebelopsio/gohan.git
    cd gohan

    # Build the binary
    go build -o gohan ./cmd/gohan

    # Optional: Install system-wide
    sudo cp gohan /usr/local/bin/

    # Verify installation
    gohan version
    ```

=== "Pre-built Binary"

    !!! success "Quick installation"
        Download and install a pre-compiled binary in seconds.

    ```bash
    # Download latest release
    wget https://github.com/rebelopsio/gohan/releases/latest/download/gohan-linux-amd64

    # Make executable
    chmod +x gohan-linux-amd64
    sudo mv gohan-linux-amd64 /usr/local/bin/gohan

    # Verify installation
    gohan version
    ```

=== "Go Install"

    !!! info "For Go developers"
        Install directly using Go's package manager.

    ```bash
    # Install directly with go install
    go install github.com/rebelopsio/gohan/cmd/gohan@latest

    # Verify installation
    gohan version
    ```

## Post-Installation Setup

### 1. Initialize Configuration

!!! note "Automatic Configuration"
    Gohan creates its configuration directory automatically on first run.

```bash
# Initialize (happens automatically on first command)
gohan theme list

# Configuration location:
# ~/.config/gohan/
```

### 2. Verify Installation

!!! tip "Preflight Checks"
    Run preflight checks to ensure your system is ready:

```bash
# Check system requirements
gohan preflight check
```

This will verify:

- Operating system compatibility
- Required packages
- File permissions
- Hyprland installation (if applicable)

### 3. Initial Theme Setup

!!! success "First Theme"
    Set your first theme to initialize the theme system:

```bash
# List available themes
gohan theme list

# Apply a theme
gohan theme set mocha

# Verify theme was applied
gohan theme show
```

## Directory Structure

After installation, Gohan uses these directories:

```
~/.config/gohan/
├── theme-state.json          # Current theme state
├── theme-history.json         # Theme change history
└── backups/                   # Configuration backups
    └── <timestamp>/          # Timestamped backup directories

# System configurations (created by theme application)
~/.config/
├── hypr/hyprland.conf        # Hyprland config
├── waybar/style.css          # Waybar styling
├── kitty/kitty.conf          # Kitty terminal
├── rofi/config.rasi          # Rofi launcher
├── mako/config               # Mako notifications
├── alacritty/alacritty.toml  # Alacritty terminal
└── fuzzel/fuzzel.ini         # Fuzzel launcher
```

## Upgrading

### Upgrade from Source

```bash
# Navigate to repository
cd gohan

# Pull latest changes
git pull origin main

# Rebuild
go build -o gohan ./cmd/gohan

# Replace existing binary
sudo cp gohan /usr/local/bin/
```

### Upgrade Pre-built Binary

```bash
# Download latest release
wget https://github.com/rebelopsio/gohan/releases/latest/download/gohan-linux-amd64

# Replace existing binary
chmod +x gohan-linux-amd64
sudo mv gohan-linux-amd64 /usr/local/bin/gohan

# Verify upgrade
gohan version
```

## Uninstallation

### Remove Gohan Binary

```bash
# Remove system-wide installation
sudo rm /usr/local/bin/gohan

# Or remove from Go bin
rm ~/go/bin/gohan
```

### Remove Configuration (Optional)

```bash
# Remove Gohan configuration directory
rm -rf ~/.config/gohan

# Note: This removes:
# - Theme state and history
# - All configuration backups
```

### Restore Original Configurations

```bash
# Gohan creates backups before modifying files
# To restore from a backup:
gohan config rollback <backup-id>

# Or manually restore from backup directory
cp ~/.config/gohan/backups/<timestamp>/hyprland.conf ~/.config/hypr/hyprland.conf
```

## Troubleshooting Installation

!!! warning "Common Issues"
    Here are solutions to common installation problems.

??? failure "command not found: gohan"
    **Problem**: Gohan binary not in PATH

    **Solution**:
    ```bash
    # If installed to /usr/local/bin, ensure it's in PATH
    echo $PATH | grep /usr/local/bin

    # If not, add to ~/.bashrc or ~/.zshrc:
    export PATH="/usr/local/bin:$PATH"
    source ~/.bashrc  # or ~/.zshrc
    ```

??? failure "permission denied"
    **Problem**: Binary not executable

    **Solution**:
    ```bash
    chmod +x /path/to/gohan
    ```

??? failure "Build fails with Go errors"
    **Problem**: Go version too old

    **Solution**:
    ```bash
    # Check Go version
    go version

    # Should be 1.21 or higher
    # Install newer Go from https://go.dev/dl/
    ```

??? failure "Missing dependencies"
    **Problem**: Template system requires certain packages

    **Solution**:
    ```bash
    # Install required development tools
    sudo apt-get update
    sudo apt-get install build-essential git

    # For full Hyprland support
    gohan install hyprland-complete
    ```

## Next Steps

- **[Theme Management](theme-management.md)** - Learn about themes
- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions
- **[Development Guide](development.md)** - Contributing and architecture
