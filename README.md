# Gohan - Hyprland Environment Manager

<div align="center">

**Your all-in-one Hyprland desktop environment manager**

Beautiful themes • Automated installation • Seamless configuration management

[Documentation](https://rebelopsio.github.io/gohan) • [Features](#features) • [Quick Start](#quick-start) • [Contributing](#contributing)

</div>

---

## 🌟 Overview

Gohan is a comprehensive tool for managing your Hyprland desktop environment on Debian systems. It combines automated installation, beautiful theming with hot-reload capabilities, and powerful configuration management into a single, elegant CLI tool.

### Why Gohan?

- **🎨 Beautiful Themes**: 5 built-in Catppuccin themes with instant hot-reload
- **⚡ One-Command Setup**: Complete Hyprland installation with all dependencies
- **🔄 Smart Configuration**: Template-based configs with automatic theme integration
- **💾 Safe Updates**: Automatic backups before every change
- **🔧 Interactive TUI**: Beautiful terminal interfaces for guided workflows
- **📦 Dependency-Free**: Single binary with no external requirements

---

## ✨ Features

### 🎨 Theme Management

Switch between gorgeous themes with automatic hot-reload:

```bash
# List available themes
gohan theme list

# Preview colors before applying
gohan theme preview mocha

# Apply theme instantly (auto-reloads Hyprland & Waybar)
gohan theme set mocha

# Rollback to previous theme
gohan theme rollback
```

**Available Themes:**
- 🌙 **Mocha** - Warm dark theme (default)
- ☀️ **Latte** - Soft light theme
- 🌆 **Frappe** - Cool dark theme
- 🌸 **Macchiato** - Purple-tinted dark theme
- 🚀 **Gohan** - Custom branded theme

### 📦 Automated Installation

Complete Hyprland setup with preflight validation:

```bash
# Interactive installation with wizard
gohan install hyprland-complete

# Run preflight checks first
gohan preflight run

# View installation progress in real-time
gohan install --progress
```

**Includes:**
- Hyprland compositor
- Waybar status bar
- Kitty terminal
- Rofi/Fuzzel launcher
- Mako notifications
- All required dependencies

### ⚙️ Configuration Management

Deploy and manage configuration files:

```bash
# Deploy all configurations
gohan config deploy

# Deploy specific components
gohan config deploy --components hyprland,waybar

# Preview before deploying
gohan config deploy --dry-run

# List available components
gohan config list
```

**Supported Components:**
- Hyprland core configuration
- Waybar (status bar)
- Kitty (terminal)
- Rofi (launcher)
- Fuzzel (launcher)
- Mako (notifications)
- Alacritty (terminal)

### 💾 Backup & Restore

Automatic backups with easy restoration:

```bash
# List all backups
gohan backup list

# Create manual backup
gohan backup create

# Restore from backup
gohan backup restore <backup-id>

# Clean up old backups
gohan backup cleanup --older-than 30d
```

### 📊 System Health

Monitor and verify your Hyprland setup:

```bash
# Run health checks
gohan doctor

# Check installation status
gohan status

# View installation history
gohan history list
```

---

## 🚀 Quick Start

### Prerequisites

- **OS**: Debian Sid (unstable) or Trixie (testing)
- **Architecture**: amd64
- **Permissions**: sudo access for installation

> ⚠️ **Important**: Ubuntu and Debian Bookworm are NOT supported. Gohan requires cutting-edge packages from Debian's rolling branches.

### Installation

```bash
# Clone the repository
git clone https://github.com/rebelopsio/gohan.git
cd gohan

# Build from source
go build -o gohan ./cmd/gohan

# Move to PATH (optional)
sudo mv gohan /usr/local/bin/

# Initialize Gohan
gohan init

# Verify installation
gohan version
```

### First Run

```bash
# Run preflight checks
gohan preflight run

# Install Hyprland and components
gohan install hyprland-complete

# Apply your favorite theme
gohan theme set mocha

# Verify everything works
gohan doctor
```

---

## 🏗️ Architecture

Gohan is built with **Clean Architecture** principles and **Test-Driven Development**:

```
┌─────────────────────────────────────────────┐
│           CLI Layer (Cobra)                  │
│  Interactive TUI (Bubble Tea + Lipgloss)     │
├─────────────────────────────────────────────┤
│        Application Use Cases                 │
│  (Installation, Theme, Config, Backup)       │
├─────────────────────────────────────────────┤
│         Domain Models & Logic                │
│  (Preflight, Packages, Themes, Sessions)     │
├─────────────────────────────────────────────┤
│       Infrastructure Layer                   │
│  APT, Templates, Backup, History, Services   │
└─────────────────────────────────────────────┘
```

**Key Technologies:**
- **Language**: Go 1.23+
- **CLI Framework**: Cobra
- **TUI**: Bubble Tea, Lipgloss, Bubbles
- **Testing**: Testify, BDD with Gherkin
- **Templates**: Custom Jinja2-style engine
- **Storage**: SQLite (BBolt)
- **Package Management**: APT integration

---

## 📖 Documentation

Full documentation is available at **[rebelopsio.github.io/gohan](https://rebelopsio.github.io/gohan)**

**Guides:**
- [Installation Guide](https://rebelopsio.github.io/gohan/installation) - Complete setup instructions
- [Theme Management](https://rebelopsio.github.io/gohan/theme-management) - Theme switching and customization
- [Troubleshooting](https://rebelopsio.github.io/gohan/troubleshooting) - Common issues and solutions
- [Development Guide](https://rebelopsio.github.io/gohan/development) - Contributing and architecture

---

## 🎯 Use Cases

### For New Users
Get started with Hyprland in minutes:
```bash
gohan install hyprland-complete
gohan theme set latte
```

### For Theme Enthusiasts
Experiment with different color schemes:
```bash
# Try different themes
gohan theme preview mocha
gohan theme set mocha

# Don't like it? Rollback instantly
gohan theme rollback
```

### For Power Users
Full control over your environment:
```bash
# Custom configuration deployment
gohan config deploy --components hyprland,waybar

# Track installation history
gohan history list

# Restore previous states
gohan backup restore <backup-id>
```

---

## 🧪 Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./tests/integration/...

# Run specific test
go test -v ./internal/application/configuration/...
```

### BDD Features

Gohan follows **Behavior-Driven Development**. Features are defined in Gherkin:

```gherkin
Feature: Theme Management
  As a Hyprland user
  I want to switch between themes
  So that I can customize my desktop appearance

  Scenario: Apply new theme
    Given I have Hyprland installed
    When I run "gohan theme set mocha"
    Then the theme should be applied
    And Hyprland should reload automatically
```

See `docs/features/*.feature` for all feature definitions.

### Project Structure

```
gohan/
├── cmd/gohan/           # Application entrypoint
├── internal/
│   ├── cli/             # Cobra CLI commands
│   ├── tui/             # Bubble Tea TUI components
│   ├── application/     # Use cases and orchestration
│   ├── domain/          # Business logic and models
│   └── infrastructure/  # External integrations
├── templates/           # Configuration file templates
├── tests/
│   └── integration/     # Integration tests
└── docs/                # MkDocs documentation
```

---

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Write tests**: Follow TDD approach
4. **Commit changes**: `git commit -m 'Add amazing feature'`
5. **Push to branch**: `git push origin feature/amazing-feature`
6. **Open a Pull Request**

**Development Guidelines:**
- Write tests first (TDD)
- Follow Clean Architecture
- Document with Gherkin scenarios
- Use conventional commits
- Run `go fmt` and `go vet`

See [CONTRIBUTING.md](docs/development.md) for detailed guidelines.

---

## 📋 Roadmap

- [x] Phase 1: Core installation engine
- [x] Phase 2: Package management
- [x] Phase 3: Repository setup
- [x] Phase 4: Configuration deployment
- [x] Phase 5: Theme management
- [x] Phase 6: Backup & restore
- [ ] Phase 7: Post-installation automation
- [ ] Phase 8: Web dashboard
- [ ] Phase 9: Plugin system
- [ ] Phase 10: Multi-distro support

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- **[Hyprland](https://hyprland.org/)** - The amazing Wayland compositor
- **[Catppuccin](https://github.com/catppuccin)** - Beautiful color schemes
- **[Charm Bracelet](https://charm.sh/)** - Excellent TUI libraries
- **Debian Community** - For the solid foundation

---

## 📞 Support

- **Documentation**: [rebelopsio.github.io/gohan](https://rebelopsio.github.io/gohan)
- **Issues**: [GitHub Issues](https://github.com/rebelopsio/gohan/issues)
- **Discussions**: [GitHub Discussions](https://github.com/rebelopsio/gohan/discussions)

---

<div align="center">

**Made with ❤️ by the RebelOps team**

[⭐ Star us on GitHub](https://github.com/rebelopsio/gohan) • [📖 Read the docs](https://rebelopsio.github.io/gohan) • [🐛 Report a bug](https://github.com/rebelopsio/gohan/issues)

</div>
