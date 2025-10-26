# Gohan: Building an Omakase Hyprland Installer for Debian Sid/Trixie

**Executive Summary:** Debian Sid/Trixie provides the ideal foundation for an Omarchy-inspired Hyprland installer. With official repository support, bleeding-edge dependencies, and a proven community installer ecosystem, Debian enables the polished experience that Ubuntu cannot deliver. This document outlines the technical architecture for **Gohan** (ご飯) - an opinionated, beautiful Hyprland setup tool that transforms Debian Sid into a production-ready tiling compositor environment.

---

## The Opportunity: Why Debian Sid Changes Everything

Omarchy succeeded by providing zero-configuration transformation of Arch Linux into a beautiful Hyprland environment. The challenge has been bringing that experience to Debian-based systems. **Debian Sid/Trixie fundamentally changes the equation** - Hyprland has official repository support, near-Arch package freshness, and active maintainer commitment that Ubuntu completely lacks.

### The Debian Advantage

**Official Repository Support:** Hyprland is now available in Debian Sid's official repository and can be installed with `sudo apt install hyprland`. This is transformative - Debian maintainers actively package and update Hyprland, unlike Ubuntu where users are completely on their own.

**Active Maintenance:** The Debian package tracker shows ongoing development with recent commits updating Hyprland to version 0.47.2 in March 2025, with regular updates to dependencies like libaquamarine and hyprutils. The packaging is alive and evolving.

**Bleeding-Edge Dependencies:** Debian Sid (unstable) updates every 6 hours with the latest packages, providing the modern libraries Hyprland requires without the years-long freeze of Ubuntu LTS. Ubuntu stable is based on a snapshot of Debian unstable but then freezes those packages, making Sid actually newer than Ubuntu between releases.

**Proven Installer Foundation:** JaKooLit's Debian-Hyprland installer is mature and actively maintained for both Trixie and Sid, building Hyprland 0.51.1 from source with automated updates. The groundwork exists - Gohan's value is adding polish, UX, and Omarchy-inspired integration on top of this proven foundation.

### The Honest Challenges

**Sid is Forever Unstable:** Debian Sid is explicitly "forever unstable" with no release-like quality assurance or integration testing. It's "rolling development" not "rolling release." Users must accept occasional breakage as the cost of bleeding-edge packages.

**Security Model:** Sid exclusively gets security updates through package maintainers, not from Debian's Security Team. This requires users to monitor and apply updates promptly.

**User Sophistication Required:** Recommended practices include using `apt upgrade` instead of `apt full-upgrade` to avoid unwanted package removal, installing apt-listbugs and apt-listchanges to catch issues early, and timing updates carefully. Sid demands Linux comfort.

**Trixie as Middle Ground:** Trixie (testing) provides 10 days of community testing before packages arrive from Sid, offering more stability at the cost of slightly older packages. However, Testing doesn't get security fixes as fast as Stable, and sometimes things break until fixed upstream in Sid.

### The Market Fit

Hyprland itself requires Linux sophistication - tiling compositor knowledge, terminal comfort, configuration editing. Sid users already meet this bar. The target user profile aligns perfectly: experienced Linux users wanting cutting-edge desktop environments without maintaining Arch. Gohan's value proposition: **"Take Debian Sid and transform it into a polished Hyprland experience with Omarchy-level integration."**

---

## What Makes Omarchy Successful (And How Gohan Adapts)

Omarchy achieved 15.7k+ stars by solving configuration paralysis with opinionated defaults. Understanding its architecture reveals patterns Gohan should emulate and adapt for Debian.

### Omarchy's Winning Patterns

**Two-Stage Bootstrap Architecture:** Omarchy separates bootstrap (`boot.sh`) from main installation (`install.sh`) which orchestrates five sequential phases: preflight, packaging, configuration, login, and finalization. This modularity enables independent section execution and clear error isolation.

**System vs User Separation:** Omarchy-managed files live in `~/.local/share/omarchy/` as a Git repository, while user configurations in `~/.config/` override defaults and remain under user control. This balance allows curated updates without destroying customization.

**Comprehensive Theme System:** Omarchy unifies aesthetics across 8+ components with symbolic link-based theme switching enabling one-hotkey changes, includes 12+ pre-configured themes, and provides a GUI theme builder (Aether).

**Migration System:** Timestamp-based migration scripts ensure sequential execution with state tracking to prevent duplicate runs, handling package transitions and breaking changes automatically.

**Hardware Integration:** Automatic GPU detection, fingerprint sensor support, Apple display brightness control, and Intel Mac compatibility demonstrate production-grade polish.

### Gohan's Debian-Specific Adaptations

Where Omarchy provides single-command zero-configuration (viable on officially-supported Arch), Gohan provides **guided installation with Debian-aware intelligence**:

- **Repo-First Strategy:** Try official Debian packages before building from source
- **Sid/Trixie Detection:** Automatic branch detection with appropriate warnings
- **Safety-First Approach:** Comprehensive backup and rollback for Sid's occasional breakage
- **Update Intelligence:** Check compatibility before applying updates
- **Package Source Tracking:** Remember which components came from repos vs source builds

---

## Hyprland on Debian: Current State and Best Practices

### Repository Status

Hyprland recently made it into the Sid repository and can be installed with `sudo apt install hyprland`. Even though Hyprland is in the Trixie repos, it is still recommended to install from Sid, as some dependencies in the Trixie repo are outdated.

Hyprland 0.41.2 is currently in Sid with migration to Testing blocked by dependency issues, showing active packaging work but also the reality of Debian's testing migration process.

### Installation Approaches

**Official Packages (Preferred):**
```bash
sudo apt install hyprland
```
Provides Hyprland itself from Debian repos, reducing build complexity significantly.

**Source Building (When Necessary):**
Most core ecosystem packages are downloaded and built from source since there are no pre-built binaries yet for Debian. The installation takes 30-60+ minutes due to compilation.

Required source builds typically include:
- xdg-desktop-portal-hyprland
- hyprlock
- hypridle  
- Latest waybar (if repo version too old)
- hyprpaper/swww (wallpaper daemons)

**Prerequisites:** Edit `/etc/apt/sources.list` to enable deb-src lines for source packaging, and ensure non-free drivers are available for NVIDIA GPU users.

### Essential Ecosystem Components

**Portal Configuration:** xdg-desktop-portal-hyprland requires building from source with proper configuration file generation. This enables screen sharing, file pickers, and other desktop integration.

**Display Managers:** Recommend SDDM or GDM - other login managers may not work or launch Hyprland properly. Hyprland can be launched through TTY by typing `Hyprland`.

**Audio Setup:** The installer scripts don't setup audio. PipeWire is recommended: `sudo apt install -y pipewire`.

**Network Management:** network-manager-gnome has been removed from automated install packages because it restarts NetworkManager during installation causing issues. Install manually after boot: `sudo apt install network-manager-gnome`.

### Common Issues and Solutions

**NVIDIA Complications:** NVIDIA users may get stuck on SDDM login. Fix involves pressing Ctrl+Alt+F2, logging in, finding GPU ID via `lspci -nn` and `ls /dev/dri/by-path`, then adding `env = WLR_DRM_DEVICES,/dev/dri/cardX` to Hyprland's environment variables config.

**Rofi Scaling Problems:** Rofi issues (scaling, unexplained scaling) are common when rofi is already installed. Fix by uninstalling rofi and installing rofi-wayland: `sudo apt autoremove rofi` then install rofi-wayland.

**Package Updates:** Re-running the install script auto-updates all packages as it's configured to pull latest package builds. This provides a convenient update path.

---

## Learning from Successful Hyprland Projects

Three major Hyprland dotfile projects demonstrate architectural patterns Gohan should adopt.

### prasanthrangan/hyprdots (HyDE)

**Modularity Through Control Files:** Control file system uses pipe-delimited `.lst` files for fine-grained package and configuration management. `custom_hypr.lst` defines packages with dependency checks, `restore_cfg.lst` specifies overwrite flags and backup behavior.

**Trust Through Transparency:** The backup system auto-creates backups before overwriting any file - users trust it because it never destroys data.

**Community Extension:** Theme management through hyde-cli enables dynamic switching with a community gallery ecosystem.

**Lesson:** Control files separate data from code, enabling extension without programming.

### end-4/dots-hyprland (illogical-impulse)

**Maximum Features:** AI integration, auto-generated Material Design colors from wallpapers (wallust), advanced UI with workspace groups.

**Transparent Installation:** Shows every command before execution, building user trust through visibility.

**Multi-Distro Support:** Distribution-specific PKGBUILD systems enable multi-distro support (Arch, OpenSUSE, Fedora experimentally) with clear dependency tracking.

**Configuration Layering:** Separates defaults (`.config/hypr/hyprland/`) from user overrides (`.config/hypr/custom/`), preserving customization across updates.

**Lesson:** Transparency and modularity enable ambitious features without overwhelming users.

### mylinuxforwork/dotfiles (ML4W)

**GUI Management Tools:** ML4W Settings App (GTK4) and ML4W Welcome App lower the barrier dramatically - users configure themes, wallpapers, and settings visually.

**Comprehensive Onboarding:** YouTube tutorials and documentation create onboarding paths for those less comfortable with terminal-only workflows.

**Lesson:** GUI tools democratize access - not everyone wants to live in terminal-only workflows.

### Common Success Patterns

✅ Comprehensive backup before any changes  
✅ Modular organization enabling independent testing  
✅ Theme systems providing cohesive aesthetics  
✅ Hardware detection (especially NVIDIA)  
✅ Clear update mechanisms respecting user modifications  

❌ Destructive operations without backups  
❌ Monolithic scripts that fail entirely if one component breaks  
❌ No rollback capability  
❌ Unclear error messages  
❌ Conflicting file ownership between system and user

---

## Dotfile Management Strategies for Gohan

Gohan shouldn't force a specific dotfile manager, but should provide clear integration paths for popular tools.

### chezmoi (Recommended for Power Users)

**Strengths:** Go-based single binary, powerful templating with `text/template`, password manager integration, machine-specific configurations via conditionals.

**Integration Pattern:**
```bash
# After Gohan installs Hyprland
gohan export --format chezmoi > ~/.local/share/chezmoi/hyprland.tmpl

# Or import existing chezmoi configs
gohan import --source chezmoi
```

### GNU Stow (Recommended for Simplicity)

**Strengths:** Symlink simplicity, transparent operation, no configuration files.

**Integration Pattern:**
```bash
# Gohan can organize configs for stow
gohan export --format stow
# Creates ~/.dotfiles/hypr/, ~/.dotfiles/waybar/, etc.

# User then runs
stow -d ~/.dotfiles -t ~ hypr waybar
```

### yadm (Git-Familiar Users)

**Strengths:** Git wrapper with alternate files for machine differences, built-in encryption.

**Integration Pattern:**
```bash
# Gohan installs to ~/.config as normal
# User initializes yadm after
yadm init
yadm add ~/.config/hypr ~/.config/waybar
yadm commit -m "Initial Hyprland setup from Gohan"
```

### Gohan's Approach

**Philosophy:** Install Hyprland first with opinionated defaults in a system-managed location. Offer dotfile setup as step two with tool selection. Provide templates for all major tools but don't force a specific approach.

**Template Generation:**
```go
// gohan export --format chezmoi
func ExportChezmoi() error {
    configs := []string{
        "~/.config/hypr/hyprland.conf",
        "~/.config/waybar/config.jsonc",
    }
    
    for _, cfg := range configs {
        // Convert to chezmoi template with conditionals
        tmpl := convertToChezmoi(cfg)
        write(tmpl, chezmoiDir(cfg))
    }
}
```

---

## Go and Charmbracelet: Building Polished TUIs

The Charmbracelet ecosystem provides production-ready components for beautiful terminal interfaces, with Huh specifically designed for installer wizards.

### Huh: The Primary Tool

**Why Huh for Gohan:** Huh runs standalone or integrates into Bubble Tea apps, provides built-in field types (Input, Select, MultiSelect, Confirm, Text), supports validation, includes themes, and organizes fields into groups for multi-step workflows.

**Accessibility:** Accessible mode replaces TUI with standard prompts for screen readers - critical for inclusivity.

**Example Gohan Wizard:**
```go
form := huh.NewForm(
    huh.NewGroup(
        huh.NewNote().
            Title("Gohan - Omakase Hyprland for Debian").
            Description("Detected: Debian Sid on AMD GPU"),
        
        huh.NewConfirm().
            Title("Use official Hyprland from Debian repos?").
            Description("Version 0.41.2 available").
            Value(&useRepoVersion),
    ),
    huh.NewGroup(
        huh.NewSelect[string]().
            Title("Display Manager").
            Options(
                huh.NewOption("SDDM (Recommended)", "sddm"),
                huh.NewOption("GDM (GNOME)", "gdm"),
                huh.NewOption("TTY Launch", "tty"),
            ).
            Value(&displayManager),
    ),
    huh.NewGroup(
        huh.NewMultiSelect[string]().
            Title("Ecosystem Components").
            Options(
                huh.NewOption("Waybar (Status bar)", "waybar"),
                huh.NewOption("Rofi (Launcher)", "rofi"),
                huh.NewOption("Mako (Notifications)", "mako"),
                huh.NewOption("Swaylock (Screen lock)", "swaylock"),
            ).
            Value(&components),
    ),
)
```

### Bubble Tea: Custom Progress Views

**When to Use:** For installation progress beyond Huh forms - real-time package building, progress bars, spinners.

**Model-View-Update Architecture:**
```go
type installModel struct {
    phase       Phase
    progress    float64
    currentPkg  string
    logs        []string
    spinner     spinner.Model
}

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case PhaseCompleteMsg:
        m.phase = msg.NextPhase
        m.progress = msg.Progress
    case PackageBuildMsg:
        m.currentPkg = msg.Package
        m.logs = append(m.logs, msg.Output)
    }
    return m, nil
}

func (m installModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        titleStyle.Render("Installing Gohan"),
        progressBar.ViewAs(m.progress),
        fmt.Sprintf("Building: %s", m.currentPkg),
        logView.Render(m.logs),
    )
}
```

### Template Management

**Embedding Templates:**
```go
//go:embed templates/*
var templates embed.FS

type TemplateEngine struct {
    templates *template.Template
}

func NewTemplateEngine() (*TemplateEngine, error) {
    tmpl, err := template.ParseFS(templates, "templates/**/*.tmpl")
    if err != nil {
        return nil, err
    }
    return &TemplateEngine{templates: tmpl}, nil
}
```

**Debian-Specific Templates:**
```
templates/hyprland/hyprland.conf.tmpl
templates/hyprland/nvidia.conf.tmpl
templates/waybar/config.jsonc.tmpl
templates/themes/catppuccin/
templates/system/apt-sources.list.tmpl
```

### Build Tags for Future Multi-Distro

```go
//go:build linux && debian

package packages

type DebianManager struct {
    release DebianRelease
}

func Detect() Manager {
    if isDebian() {
        return &DebianManager{release: detectRelease()}
    }
    // Future: Arch, Fedora, etc.
}
```

---

## Technical Architecture for Gohan

### Core Abstractions

**Package Manager Interface:**
```go
type PackageManager interface {
    Install(pkg string) error
    Remove(pkg string) error
    Update(pkg string) error
    IsInstalled(pkg string) bool
    GetVersion(pkg string) (string, error)
}

type PackageSource string

const (
    DebianRepo  PackageSource = "debian_repo"
    BuildSource PackageSource = "build_source"
    Binary      PackageSource = "binary"
)

// Strategy: Try repo first, fallback to build
func (m *DebianManager) Install(pkg string) error {
    if version, available := m.checkDebianRepo(pkg); available {
        if m.isVersionAcceptable(pkg, version) {
            return m.aptInstall(pkg)
        }
    }
    return m.buildFromSource(pkg)
}
```

**Configuration Template Engine:**
```go
type TemplateData struct {
    System   SystemInfo
    Hardware HardwareInfo
    Theme    ThemeConfig
    User     UserPreferences
}

type SystemInfo struct {
    OS       string
    Release  DebianRelease
    Arch     string
}

type HardwareInfo struct {
    GPU      GPUConfig
    Monitors []MonitorConfig
}

type GPUConfig struct {
    Vendor   string
    IsNVIDIA bool
    IsAMD    bool
    IsIntel  bool
    EnvVars  map[string]string
}

func (e *TemplateEngine) GenerateHyprlandConfig(data *TemplateData) error {
    content, err := e.Render("hyprland/hyprland.conf.tmpl", data)
    if err != nil {
        return err
    }
    
    dest := expandPath("~/.config/hypr/hyprland.conf")
    return e.WriteAtomic(dest, content)
}
```

**Backup and Checkpoint System:**
```go
type CheckpointManager struct {
    checkpointDir string
}

type Checkpoint struct {
    ID        string
    Phase     Phase
    Timestamp time.Time
    Paths     []BackupPath
    Packages  []PackageSnapshot
}

func (m *CheckpointManager) Create(phase Phase) (*Checkpoint, error) {
    cp := &Checkpoint{
        ID:        generateID(),
        Phase:     phase,
        Timestamp: time.Now(),
    }
    
    // Backup configs
    configPaths := []string{
        "~/.config/hypr",
        "~/.config/waybar",
        "~/.zshrc",
    }
    
    for _, path := range configPaths {
        if err := m.backupPath(cp, path); err != nil {
            return nil, err
        }
    }
    
    // Snapshot installed packages
    cp.Packages = m.snapshotPackages()
    
    return cp, m.save(cp)
}

func (m *CheckpointManager) Rollback(checkpointID string) error {
    cp, err := m.load(checkpointID)
    if err != nil {
        return err
    }
    
    // Restore configs
    for _, path := range cp.Paths {
        if err := m.restorePath(path); err != nil {
            return err
        }
    }
    
    // Rollback packages (best effort)
    for _, pkg := range cp.Packages {
        m.rollbackPackage(pkg)
    }
    
    return nil
}
```

### Phase-Based Installation

```go
type Phase string

const (
    Preflight   Phase = "preflight"
    Backup      Phase = "backup"
    Repository  Phase = "repository"
    Packages    Phase = "packages"
    Config      Phase = "config"
    PostInstall Phase = "post_install"
    Verify      Phase = "verify"
)

type Orchestrator struct {
    state      *State
    backup     *CheckpointManager
    packages   PackageManager
    templates  *TemplateEngine
    hardware   *HardwareDetector
}

func (o *Orchestrator) Execute() error {
    phases := []Phase{
        Preflight,
        Backup,
        Repository,
        Packages,
        Config,
        PostInstall,
        Verify,
    }
    
    for _, phase := range phases {
        // Create checkpoint before phase
        cp, err := o.backup.Create(phase)
        if err != nil {
            return fmt.Errorf("checkpoint failed: %w", err)
        }
        
        // Execute phase
        if err := o.executePhase(phase); err != nil {
            // Rollback on failure
            o.backup.Rollback(cp.ID)
            return fmt.Errorf("phase %s failed: %w", phase, err)
        }
        
        o.state.RecordSuccess(phase)
    }
    
    return nil
}

func (o *Orchestrator) executePhase(phase Phase) error {
    switch phase {
    case Preflight:
        return o.runPreflight()
    case Backup:
        return o.runBackup()
    case Repository:
        return o.setupRepositories()
    case Packages:
        return o.installPackages()
    case Config:
        return o.generateConfigs()
    case PostInstall:
        return o.runPostInstall()
    case Verify:
        return o.verifyInstallation()
    }
    return nil
}
```

### Preflight Checks

```go
func (o *Orchestrator) runPreflight() error {
    checks := []PreflightCheck{
        {
            Name: "Debian Version",
            Check: func() error {
                release, err := detectDebianRelease()
                if err != nil {
                    return err
                }
                if release == Bookworm {
                    return errors.New("Debian Bookworm not supported - use Sid or Trixie")
                }
                return nil
            },
        },
        {
            Name: "Disk Space",
            Check: func() error {
                available := getAvailableSpace("/")
                required := 10 * GB
                if available < required {
                    return fmt.Errorf("need %d GB, have %d GB", required/GB, available/GB)
                }
                return nil
            },
        },
        {
            Name: "Internet Connection",
            Check: func() error {
                return checkConnectivity("debian.org")
            },
        },
        {
            Name: "deb-src Enabled",
            Check: func() error {
                return checkDebSrcEnabled()
            },
        },
    }
    
    for _, check := range checks {
        if err := check.Check(); err != nil {
            return fmt.Errorf("%s: %w", check.Name, err)
        }
    }
    
    return nil
}
```

---

## Gohan Implementation Blueprint

### Project Structure

```
gohan/
├── cmd/
│   └── gohan/
│       └── main.go
├── internal/
│   ├── cli/              # Command implementations
│   │   ├── init.go
│   │   ├── theme.go
│   │   ├── update.go
│   │   ├── doctor.go
│   │   └── rollback.go
│   ├── orchestrator/     # Installation phases
│   │   ├── phase.go
│   │   ├── state.go
│   │   └── checkpoint.go
│   ├── system/           # OS/hardware detection
│   │   ├── detect.go
│   │   ├── debian.go
│   │   └── hardware.go
│   ├── packages/         # Package management
│   │   ├── manager.go
│   │   ├── apt.go
│   │   └── builder.go
│   ├── config/           # Config generation
│   │   ├── generator.go
│   │   ├── template.go
│   │   └── validator.go
│   ├── backup/           # Backup/rollback
│   │   ├── manager.go
│   │   ├── snapshot.go
│   │   └── restore.go
│   ├── theme/            # Theme system
│   │   ├── manager.go
│   │   ├── registry.go
│   │   └── builder.go
│   └── ui/               # TUI components
│       ├── wizard.go
│       ├── progress.go
│       └── spinner.go
├── templates/            # Embedded configs
│   ├── hyprland/
│   ├── waybar/
│   └── themes/
├── scripts/              # Shell helpers
│   ├── bootstrap.sh
│   └── post-install/
├── docs/
│   ├── README.md
│   ├── DEBIAN.md
│   └── CONTRIBUTING.md
├── go.mod
└── go.sum
```

### Installation Flow

**Phase 1: Preflight**
- Detect Debian Sid/Trixie (fail on Bookworm)
- Check GPU (NVIDIA detection critical)
- Verify disk space (10GB+)
- Confirm internet connectivity
- Validate deb-src enabled

**Phase 2: Backup**
- Create timestamped checkpoint
- Backup ~/.config/hypr, ~/.config/waybar, etc.
- Snapshot current package state
- Save metadata

**Phase 3: Repository Setup**
- Update package lists
- Enable non-free for NVIDIA if needed
- Add any required third-party repos

**Phase 4: Package Installation**
- Try Hyprland from Debian repos first
- Build xdg-desktop-portal-hyprland from source
- Build/install waybar, rofi-wayland, mako
- Install display manager (SDDM/GDM)
- Install audio (PipeWire)
- Install fonts and themes

**Phase 5: Configuration**
- Generate Hyprland config from templates
- Apply GPU-specific settings (NVIDIA env vars)
- Configure waybar, rofi, mako
- Set up theme
- Create desktop entries

**Phase 6: Post-Install**
- Configure display manager
- Set up shell (zsh with theme)
- Install network-manager-gnome
- Generate wallpaper cache
- Enable services

**Phase 7: Verify**
- Check Hyprland binary exists
- Verify portal configuration
- Test theme files present
- Validate display manager setup

### Command Suite

```bash
# Primary commands
gohan init              # Installation wizard
gohan theme list        # List themes
gohan theme set <name>  # Switch theme
gohan update            # Update system
gohan rollback [id]     # Rollback to checkpoint
gohan doctor            # Health check

# Advanced commands
gohan config edit       # Safe config editing
gohan backup create     # Manual backup
gohan backup list       # List backups
gohan export <format>   # Export to dotfile manager
gohan import <source>   # Import existing configs
```

### Theme System

```go
type Theme struct {
    ID          string
    Name        string
    Description string
    Author      string
    Preview     string
    Components  ThemeComponents
}

type ThemeComponents struct {
    Hyprland HyprlandTheme
    Waybar   WaybarTheme
    Rofi     RofiTheme
    Terminal TerminalTheme
    GTK      GTKTheme
}

type ThemeManager struct {
    activeTheme  string
    registry     *ThemeRegistry
    configWriter *ConfigWriter
}

func (m *ThemeManager) Switch(themeID string) error {
    theme, err := m.registry.Get(themeID)
    if err != nil {
        return err
    }
    
    // Write theme configs
    if err := m.applyTheme(theme); err != nil {
        return err
    }
    
    // Reload components
    return m.reloadComponents()
}

func (m *ThemeManager) reloadComponents() error {
    // Hyprland reload
    exec.Command("hyprctl", "reload").Run()
    
    // Waybar restart
    exec.Command("killall", "waybar").Run()
    time.Sleep(100 * time.Millisecond)
    exec.Command("waybar", "&").Start()
    
    // Notify user
    exec.Command("notify-send", "Gohan", "Theme applied").Run()
    
    return nil
}
```

### Doctor Command

```go
func runDoctor() error {
    checks := []HealthCheck{
        {
            Name:        "Hyprland Version",
            Check:       checkHyprlandVersion,
            Fix:         suggestHyprlandUpdate,
        },
        {
            Name:        "GPU Driver",
            Check:       checkGPUDriver,
            Fix:         suggestDriverInstall,
        },
        {
            Name:        "Portal Configuration",
            Check:       checkPortalConfig,
            Fix:         regeneratePortalConfig,
        },
        {
            Name:        "Audio Setup",
            Check:       checkPipeWire,
            Fix:         installPipeWire,
        },
        {
            Name:        "Waybar Running",
            Check:       checkWaybarProcess,
            Fix:         restartWaybar,
        },
        {
            Name:        "Theme Integrity",
            Check:       validateThemeFiles,
            Fix:         repairTheme,
        },
    }
    
    results := make([]CheckResult, 0)
    
    for _, check := range checks {
        result := check.Check()
        results = append(results, result)
        
        // Display with color
        displayResult(result)
        
        // Offer fix if available
        if !result.Passed && check.Fix != nil {
            if askToFix() {
                check.Fix()
            }
        }
    }
    
    return nil
}
```

---

## Testing Strategy

### Unit Tests
```bash
go test ./internal/...
```

Focus on business logic:
- Package manager detection
- Template rendering
- Backup/restore logic
- Configuration validation

### Integration Tests (Docker)

```dockerfile
# tests/docker/sid.Dockerfile
FROM debian:sid

RUN apt-get update && apt-get install -y \
    git build-essential curl sudo

# Add test user
RUN useradd -m -s /bin/bash tester && \
    echo "tester ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

USER tester
WORKDIR /home/tester

COPY . /home/tester/gohan
WORKDIR /home/tester/gohan

RUN make test-integration
```

```yaml
# docker-compose.yml
version: '3'
services:
  sid-test:
    build:
      context: .
      dockerfile: tests/docker/sid.Dockerfile
    volumes:
      - .:/gohan
  
  trixie-test:
    build:
      context: .
      dockerfile: tests/docker/trixie.Dockerfile
    volumes:
      - .:/gohan
```

### End-to-End Tests (VM)

```bash
# tests/vm/test-install.sh
#!/bin/bash

# Create fresh Debian Sid VM
vagrant up debian-sid

# Copy gohan
vagrant ssh -c "curl -sSL https://gohan.sh | sh"

# Run installation
vagrant ssh -c "gohan init --preset tests/presets/default.yaml"

# Verify Hyprland launches
vagrant ssh -c "Hyprland --version"

# Cleanup
vagrant destroy -f
```

### CI/CD Pipeline

```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        debian: [sid, trixie]
    
    steps:
      - uses: actions/checkout@v3
      
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - name: Unit Tests
        run: make test
      
      - name: Integration Tests
        run: |
          docker-compose build ${{ matrix.debian }}-test
          docker-compose run ${{ matrix.debian }}-test
      
      - name: Build Binary
        run: make build
      
      - name: Upload Artifact
        uses: actions/upload-artifact@v3
        with:
          name: gohan-${{ matrix.debian }}
          path: bin/gohan
```

---

## Distribution & Release

### Quick Install Script

```bash
# install.sh
#!/bin/bash
set -e

# Detect OS
if [ ! -f /etc/debian_version ]; then
    echo "Error: Debian required"
    exit 1
fi

# Detect release
release=$(grep VERSION_CODENAME /etc/os-release | cut -d= -f2)
if [ "$release" = "bookworm" ]; then
    echo "Error: Debian Bookworm not supported. Use Sid or Trixie."
    exit 1
fi

# Download latest release
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
esac

URL="https://github.com/user/gohan/releases/latest/download/gohan-linux-$ARCH"

echo "Downloading Gohan..."
curl -fsSL "$URL" -o /tmp/gohan
chmod +x /tmp/gohan
sudo mv /tmp/gohan /usr/local/bin/gohan

echo "Gohan installed! Run 'gohan init' to begin."
```

### Release Process (goreleaser)

```yaml
# .goreleaser.yaml
project_name: gohan

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}

archives:
  - format: tar.gz
    name_template: >-
      {{ .ProjectName }}_
      {{- .Version }}_
      {{- .Os }}_
      {{- .Arch }}

checksum:
  name_template: 'checksums.txt'

release:
  github:
    owner: user
    name: gohan
  
  name_template: "v{{.Version}}"
```

---

## Success Metrics & Roadmap

### MVP Success (Phase 1: 2-3 weeks)
- [ ] Successful install on clean Debian Sid VM
- [ ] Hyprland launches and is usable
- [ ] Backup/rollback works
- [ ] < 30 minute installation time
- [ ] 5 core themes included
- [ ] `gohan doctor` catches common issues

### Feature Complete (Phase 2: 4-6 weeks)
- [ ] Theme switching works seamlessly
- [ ] 10+ themes available
- [ ] GPU detection (NVIDIA/AMD/Intel)
- [ ] Display manager integration (SDDM/GDM)
- [ ] Update system (`gohan update`)
- [ ] Documentation covers 90% of use cases
- [ ] 100+ users successfully installed

### Community Growth (Phase 3: Ongoing)
- [ ] 20+ community themes
- [ ] Plugin system for extensions
- [ ] 500+ GitHub stars
- [ ] Active Discord/discussions
- [ ] YouTube tutorials
- [ ] Trixie support stabilized
- [ ] Consider Arch support (future)

---

## Critical Debian Sid Considerations

### User Expectation Setting

Gohan's messaging must be crystal clear about Sid's nature:

> **Gohan targets Debian Sid - the rolling development branch.**  
> 
> ✓ You get the latest Hyprland and ecosystem tools  
> ✓ Official Debian package support  
> ✓ Active maintainer community  
> 
> ⚠ Expect occasional breakage (this is Sid!)  
> ⚠ Security updates come via package maintainers, not Security Team  
> ⚠ You need Linux comfort for troubleshooting  
> 
> **Trixie (Testing) users:** More stability, slightly older packages, same install process.

### Recommended User Practices

Gohan should guide users to Sid best practices:

```bash
# Setup apt-listbugs and apt-listchanges
sudo apt install apt-listbugs apt-listchanges

# Use apt upgrade (not full-upgrade) by default
alias apt-upgrade='sudo apt update && sudo apt upgrade'

# Create backups before major updates
gohan backup create "pre-update-$(date +%Y%m%d)"
```

### Update Strategy

```go
func (c *UpdateCommand) Run() error {
    // Check for breaking changes
    breaking, err := c.checkBreakingChanges()
    if err != nil {
        return err
    }
    
    if len(breaking) > 0 {
        fmt.Println("⚠ Breaking changes detected:")
        for _, change := range breaking {
            fmt.Printf("  - %s\n", change)
        }
        if !askConfirm("Continue anyway?") {
            return nil
        }
    }
    
    // Create checkpoint
    cp, err := c.backup.Create("pre-update")
    if err != nil {
        return err
    }
    
    // Run updates
    if err := c.runUpdates(); err != nil {
        fmt.Println("Update failed, rolling back...")
        return c.backup.Rollback(cp.ID)
    }
    
    return nil
}
```

---

## Conclusion

**Gohan brings Omarchy's opinionated polish to Debian Sid** by leveraging official repository support, bleeding-edge dependencies, and proven community installers. The technical architecture balances automation with safety, providing comprehensive backup/rollback for Sid's occasional instability while delivering a beautiful, cohesive Hyprland experience.

**Key differentiators:**
- Repo-first package strategy (unlike Ubuntu's build-everything approach)
- Sid-aware update intelligence
- Comprehensive checkpoint system
- Go + Charmbracelet TUI polish
- Omarchy-inspired theme system
- Clear user expectation management

**The path forward:** Build MVP on proven JaKooLit foundation, add Go-based orchestration and UX, create comprehensive theme system, enable community extension, and maintain obsessive focus on backup/rollback as the safety net for Sid's inherent instability.

Gohan's target user already runs Sid - they understand the tradeoffs. The value proposition: **transform raw Sid into a polished, beautiful Hyprland environment with Omarchy-level integration, backed by comprehensive safety systems for when Sid inevitably breaks.**
