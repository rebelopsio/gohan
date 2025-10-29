# Gohan - Hyprland Environment Manager

Welcome to Gohan documentation! Gohan is a comprehensive tool for managing your Hyprland desktop environment with automated installation, beautiful theming, and seamless configuration management.

!!! success "What is Gohan?"
    Gohan is your all-in-one Hyprland environment manager that brings beautiful themes, automated installation, and powerful configuration management to your fingertips.

## 🚀 Quick Start

```bash title="Install and run Gohan"
# Install Gohan
git clone https://github.com/rebelopsio/gohan.git
cd gohan
go build -o gohan ./cmd/gohan

# Apply a theme
./gohan theme set latte

# List available themes
./gohan theme list
```

## 📚 Documentation

### Getting Started
- **[Installation](installation.md)** - Complete installation guide

### Features
- **[Theme Management](theme-management.md)** - Browse, apply, and customize themes

### Reference
- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions
- **[Development Guide](development.md)** - Contributing and architecture

## ✨ Features

### 🎨 Theme Management
- **5 Built-in Themes**: Catppuccin Mocha, Latte, Frappe, Macchiato, and Gohan
- **Live Preview**: See colors before applying
- **One-Command Application**: Updates all components instantly
- **Theme History & Rollback**: Undo theme changes anytime
- **Auto-Reload**: Components refresh automatically

### 📦 Automated Installation
- **Preflight Validation**: Checks system requirements
- **Dependency Management**: Handles all package installations
- **Progress Tracking**: Real-time installation status
- **Backup System**: Automatic configuration backups
- **Rollback Support**: Revert changes if needed

### ⚙️ Configuration Management
- **Template System**: Jinja2-style template processing
- **7 Component Support**: Hyprland, Waybar, Kitty, Rofi, Mako, Alacritty, Fuzzel
- **Variable Injection**: Automatic theme color integration
- **Atomic Deployments**: Safe file updates
- **Backup Before Modify**: Never lose your configs

### 🔄 Hot Reload
- **Automatic Hyprland Reload**: Seamless config updates
- **Waybar Restart**: Instant bar theme updates
- **No Manual Intervention**: Everything happens automatically

## 🎯 Use Cases

=== "New Users"

    !!! tip "Complete Setup"
        Get started with Hyprland in minutes with our automated installer.

    ```bash
    # Complete Hyprland setup in one command
    gohan install hyprland-complete

    # Apply your favorite theme
    gohan theme set mocha
    ```

=== "Theme Enthusiasts"

    !!! info "Theme Management"
        Browse, preview, and apply themes with instant hot-reload.

    ```bash
    # Browse themes
    gohan theme list

    # Preview before applying
    gohan theme preview latte

    # Apply and auto-reload
    gohan theme set latte

    # Rollback if you change your mind
    gohan theme rollback
    ```

=== "Power Users"

    !!! success "Advanced Features"
        Customize everything with templates and manage your configuration history.

    ```bash
    # Custom configurations with templates
    gohan config deploy --template custom.tmpl

    # View installation history
    gohan history list

    # Rollback installations
    gohan config rollback <backup-id>
    ```

## 🏗️ Architecture

Gohan is built with clean architecture principles:

```
┌─────────────────────────────────────────┐
│          CLI Commands                    │
├─────────────────────────────────────────┤
│      Application Use Cases               │
├─────────────────────────────────────────┤
│        Domain Models                     │
├─────────────────────────────────────────┤
│      Infrastructure Layer                │
│  (Templates, Packages, Backup)           │
└─────────────────────────────────────────┘
```

- **BDD-Driven**: Features defined with Gherkin scenarios
- **TDD Implementation**: Comprehensive test coverage
- **Dependency Injection**: Clean, testable code
- **Interface-Based**: Flexible and extensible

## 🤝 Contributing

We welcome contributions! See the [Development Guide](development.md) for:
- Setting up your development environment
- Running tests
- Code standards
- Submitting pull requests

## 📄 License

Gohan is open source software.

## 🔗 Links

- **GitHub**: https://github.com/rebelopsio/gohan
- **Issues**: https://github.com/rebelopsio/gohan/issues
- **Discussions**: https://github.com/rebelopsio/gohan/discussions
