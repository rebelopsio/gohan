# CLI Reference

Complete reference for all Gohan CLI commands, flags, and options.

## Global Flags

These flags are available for all commands:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--api-url` | | API server URL | `http://localhost:8080` |
| `--verbose` | `-v` | Verbose output | `false` |
| `--help` | `-h` | Show help | |

**Example:**
```bash
gohan --verbose theme set mocha
gohan -v install hyprland-complete
```

---

## Commands

### `gohan`

Root command - shows help and version.

```bash
gohan [command]
```

### `gohan version`

Print version information.

```bash
gohan version
```

**Output:**
```
Gohan v1.0.0
Build: 2024-10-30T12:00:00Z
Commit: abc1234
Go version: go1.23.0
```

---

## Installation Commands

### `gohan init`

Initialize Gohan configuration and databases.

```bash
gohan init [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--force` | Overwrite existing configuration | `false` |
| `--config-dir` | Configuration directory | `~/.config/gohan` |

**Example:**
```bash
# Standard initialization
gohan init

# Force reinitialization
gohan init --force

# Custom config directory
gohan init --config-dir ~/my-gohan-config
```

**Output:**
```
Initializing Gohan...

✓ Created config directory: ~/.config/gohan
✓ Created database: ~/.local/share/gohan/gohan.db
✓ Initialized theme system
✓ Created backup directory

Gohan initialized successfully!
```

---

### `gohan install`

Install Hyprland and components.

```bash
gohan install <profile> [flags]
```

**Profiles:**

| Profile | Description |
|---------|-------------|
| `hyprland-complete` | Complete Hyprland setup (recommended) |
| `hyprland-minimal` | Minimal Hyprland installation |
| `components` | Individual components only |

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--dry-run` | Preview installation without making changes | `false` |
| `--no-confirm` | Skip confirmation prompts | `false` |
| `--skip-preflight` | Skip preflight checks | `false` |
| `--progress` | Show installation progress | `true` |

**Examples:**
```bash
# Complete installation (recommended)
gohan install hyprland-complete

# Minimal installation
gohan install hyprland-minimal

# Dry run to preview
gohan install hyprland-complete --dry-run

# Non-interactive installation
gohan install hyprland-complete --no-confirm
```

---

### `gohan preflight`

Run preflight checks before installation.

```bash
gohan preflight <subcommand> [flags]
```

**Subcommands:**

#### `gohan preflight run`

Run all preflight checks:

```bash
gohan preflight run [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--fix` | Attempt to fix issues automatically | `false` |
| `--strict` | Fail on warnings | `false` |
| `--json` | Output in JSON format | `false` |

**Examples:**
```bash
# Run checks
gohan preflight run

# Auto-fix issues
gohan preflight run --fix

# Strict mode (fail on warnings)
gohan preflight run --strict
```

#### `gohan preflight list`

List all preflight checks:

```bash
gohan preflight list
```

---

### `gohan check`

Quick alias for `gohan preflight run`.

```bash
gohan check [flags]
```

Same flags as `gohan preflight run`.

---

### `gohan post-install`

Run post-installation configuration.

```bash
gohan post-install [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--display-manager` | Display manager to configure | prompt |
| `--skip-sddm` | Skip SDDM configuration | `false` |
| `--auto-start` | Enable auto-start | `true` |

**Display Managers:**
- `sddm` - Simple Desktop Display Manager (recommended)
- `gdm` - GNOME Display Manager
- `tty` - TTY login (manual start)
- `none` - Skip display manager setup

**Examples:**
```bash
# Interactive post-install
gohan post-install

# Configure SDDM
gohan post-install --display-manager sddm

# TTY login only
gohan post-install --display-manager tty

# Skip display manager
gohan post-install --display-manager none
```

---

## Configuration Commands

### `gohan config`

Manage Hyprland configurations.

```bash
gohan config <subcommand> [flags]
```

**Subcommands:**

#### `gohan config deploy`

Deploy configuration files:

```bash
gohan config deploy [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--components` | Components to deploy (comma-separated) | all |
| `--dry-run` | Preview without deploying | `false` |
| `--force` | Skip confirmation prompts | `false` |
| `--skip-backup` | Don't create backup | `false` |
| `--progress` | Show progress | `false` |

**Examples:**
```bash
# Deploy all configurations
gohan config deploy

# Deploy specific components
gohan config deploy --components hyprland,waybar

# Preview deployment
gohan config deploy --dry-run

# Force deployment
gohan config deploy --force

# Deploy without backup (dangerous!)
gohan config deploy --skip-backup
```

#### `gohan config list`

List available configuration components:

```bash
gohan config list
```

---

## Theme Commands

### `gohan theme`

Manage visual themes.

```bash
gohan theme <subcommand> [flags]
```

**Subcommands:**

#### `gohan theme list`

List all available themes:

```bash
gohan theme list [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--json` | Output in JSON format | `false` |
| `--show-colors` | Display color swatches | `true` |

**Example:**
```bash
gohan theme list
```

**Output:**
```
Available Themes

🌙 mocha - Catppuccin Mocha (dark) [ACTIVE]
   Warm dark theme with muted colors
   Base: #1e1e2e  Text: #cdd6f4

☀️  latte - Catppuccin Latte (light)
   Soft light theme for daytime use
   Base: #eff1f5  Text: #4c4f69

🌆 frappe - Catppuccin Frappe (dark)
   Cool dark theme with blue tones
   Base: #303446  Text: #c6d0f5
```

#### `gohan theme set`

Apply a theme:

```bash
gohan theme set <theme-name> [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--no-reload` | Don't reload Hyprland/Waybar | `false` |
| `--skip-backup` | Don't create backup | `false` |
| `--force` | Skip confirmation | `false` |

**Examples:**
```bash
# Apply theme with auto-reload
gohan theme set mocha

# Apply without reloading
gohan theme set latte --no-reload

# Apply without backup
gohan theme set frappe --skip-backup
```

#### `gohan theme pick`

Launch interactive theme picker:

```bash
gohan theme pick
```

**Features:**
- Browse themes with arrow keys or vim bindings (j/k)
- Live color preview for each theme
- Theme descriptions and icons
- Current active theme indicator
- Apply theme instantly with Enter

**Keyboard Shortcuts:**

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` / `Space` | Select and apply theme |
| `?` | Toggle help |
| `q` / `Esc` | Quit without applying |

**Example:**
```bash
gohan theme pick
```

**Interface:**
```
🎨 Theme Selector
Choose a theme for your Hyprland environment

▶ 🌙 mocha - Catppuccin Mocha [ACTIVE]
    Warm dark theme with muted colors - perfect for evening work

  ☀️  latte - Catppuccin Latte
  🌆 frappe - Catppuccin Frappe
  🌸 macchiato - Catppuccin Macchiato

┌─────────────────────────────────────────┐
│ Color Preview:                          │
│                                          │
│   base:     ████ #1e1e2e                │
│   surface:  ████ #313244                │
│   text:     ████ #cdd6f4                │
│   mauve:    ████ #cba6f7                │
│   pink:     ████ #f5c2e7                │
│   blue:     ████ #89b4fa                │
│   green:    ████ #a6e3a1                │
│   yellow:   ████ #f9e2af                │
│   red:      ████ #f38ba8                │
└─────────────────────────────────────────┘

Press ? for help
```

#### `gohan theme preview`

Preview theme colors:

```bash
gohan theme preview <theme-name>
```

**Example:**
```bash
gohan theme preview mocha
```

**Output:**
```
Theme Preview: Catppuccin Mocha

Base Colors:
  Base:     ███ #1e1e2e
  Surface:  ███ #313244
  Overlay:  ███ #45475a

Text Colors:
  Text:     ███ #cdd6f4
  Subtext:  ███ #bac2de

Accent Colors:
  Mauve:    ███ #cba6f7
  Pink:     ███ #f5c2e7
  Blue:     ███ #89b4fa
  ...
```

#### `gohan theme show`

Show current theme:

```bash
gohan theme show
```

**Output:**
```
Current Theme: mocha

Name:         Catppuccin Mocha
Variant:      dark
Applied:      2024-10-30 12:00:00
Components:   hyprland, waybar, kitty, fuzzel
```

#### `gohan theme rollback`

Rollback to previous theme:

```bash
gohan theme rollback [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--no-reload` | Don't reload services | `false` |

**Example:**
```bash
gohan theme rollback
```

---

## Backup Commands

### `gohan backup`

Manage configuration backups.

```bash
gohan backup <subcommand> [flags]
```

**Subcommands:**

#### `gohan backup list`

List all backups:

```bash
gohan backup list [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--json` | Output in JSON format | `false` |
| `--limit` | Limit results | unlimited |
| `--sort` | Sort order (asc/desc) | `desc` |

**Example:**
```bash
gohan backup list --limit 10
```

#### `gohan backup create`

Create a manual backup:

```bash
gohan backup create [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--description` | Backup description | auto-generated |
| `--components` | Components to backup | all |

**Examples:**
```bash
# Create backup with description
gohan backup create --description "Before major changes"

# Backup specific components
gohan backup create --components hyprland,waybar
```

#### `gohan backup restore`

Restore from backup:

```bash
gohan backup restore <backup-id> [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--force` | Skip confirmation | `false` |
| `--no-reload` | Don't reload services | `false` |
| `--components` | Restore specific files only | all |

**Examples:**
```bash
# Restore backup
gohan backup restore 20241030_120000

# Force restore without confirmation
gohan backup restore 20241030_120000 --force

# Restore only Waybar
gohan backup restore 20241030_120000 --components waybar
```

#### `gohan backup delete`

Delete a backup:

```bash
gohan backup delete <backup-id> [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--force` | Skip confirmation | `false` |

**Example:**
```bash
gohan backup delete 20241029_090000 --force
```

#### `gohan backup cleanup`

Clean up old backups:

```bash
gohan backup cleanup [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--older-than` | Remove backups older than duration | |
| `--keep` | Keep N most recent backups | |
| `--dry-run` | Preview without deleting | `false` |
| `--force` | Skip confirmation | `false` |

**Examples:**
```bash
# Remove backups older than 30 days
gohan backup cleanup --older-than 30d

# Keep only 10 most recent
gohan backup cleanup --keep 10

# Preview cleanup
gohan backup cleanup --older-than 30d --dry-run
```

#### `gohan backup show`

Show backup details:

```bash
gohan backup show <backup-id>
```

---

## System Commands

### `gohan doctor`

Run system health checks:

```bash
gohan doctor [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--fix` | Attempt to fix issues | `false` |
| `--json` | Output in JSON format | `false` |

**Example:**
```bash
# Run health checks
gohan doctor

# Fix issues
gohan doctor --fix
```

**Output:**
```
Running system health checks...

✓ Hyprland is installed and running
✓ Waybar is configured correctly
✓ Theme system is operational
✓ Backups are enabled
✓ Configuration files are valid
⚠ Display manager: not configured

Overall Health: Good (1 warning)
```

---

### `gohan status`

Get installation status:

```bash
gohan status [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--json` | Output in JSON format | `false` |
| `--verbose` | Show detailed information | `false` |

**Example:**
```bash
gohan status
```

**Output:**
```
Gohan Installation Status

Core:
  ✓ Hyprland      v0.42.0
  ✓ Waybar        v0.10.3
  ✓ Kitty         v0.32.0

Theme:
  Active:  mocha (Catppuccin Mocha)
  Applied: 2024-10-30 12:00:00

Configuration:
  Deployed: 5 components
  Last update: 2024-10-30 12:00:00

Backups:
  Total: 4 backups
  Latest: 20241030_153022
```

---

### `gohan history`

View installation history:

```bash
gohan history <subcommand> [flags]
```

**Subcommands:**

#### `gohan history list`

List installation history:

```bash
gohan history list [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--limit` | Limit results | `20` |
| `--json` | Output in JSON format | `false` |

**Example:**
```bash
gohan history list --limit 10
```

#### `gohan history show`

Show detailed history entry:

```bash
gohan history show <id>
```

---

### `gohan repo`

Manage Debian repositories:

```bash
gohan repo <subcommand> [flags]
```

**Subcommands:**

#### `gohan repo setup`

Setup Debian repositories:

```bash
gohan repo setup [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--branch` | Debian branch | auto-detect |
| `--dry-run` | Preview changes | `false` |

**Example:**
```bash
gohan repo setup
```

#### `gohan repo verify`

Verify repository configuration:

```bash
gohan repo verify
```

---

### `gohan server`

Start the API server:

```bash
gohan server [flags]
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--port` | Server port | `8080` |
| `--host` | Server host | `localhost` |
| `--tls` | Enable TLS | `false` |
| `--cert` | TLS certificate file | |
| `--key` | TLS key file | |

**Example:**
```bash
# Start server
gohan server

# Custom port
gohan server --port 3000

# With TLS
gohan server --tls --cert server.crt --key server.key
```

---

## Exit Codes

Gohan uses standard exit codes:

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error |
| `2` | Misuse of command |
| `3` | Preflight checks failed |
| `4` | Installation failed |
| `5` | Configuration error |
| `6` | Backup/restore failed |
| `130` | Interrupted by user (Ctrl+C) |

**Example:**
```bash
gohan preflight run
echo $?  # Check exit code
```

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GOHAN_CONFIG_DIR` | Configuration directory | `~/.config/gohan` |
| `GOHAN_DATA_DIR` | Data directory | `~/.local/share/gohan` |
| `GOHAN_CACHE_DIR` | Cache directory | `~/.cache/gohan` |
| `GOHAN_LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `GOHAN_NO_COLOR` | Disable colored output | `false` |

**Example:**
```bash
# Enable debug logging
export GOHAN_LOG_LEVEL=debug
gohan install hyprland-complete

# Disable colors
export GOHAN_NO_COLOR=true
gohan theme list
```

---

## Shell Completion

Generate shell completion scripts:

```bash
# Bash
gohan completion bash > /etc/bash_completion.d/gohan

# Zsh
gohan completion zsh > /usr/local/share/zsh/site-functions/_gohan

# Fish
gohan completion fish > ~/.config/fish/completions/gohan.fish

# PowerShell
gohan completion powershell > gohan.ps1
```

---

## Configuration Files

### Main Configuration

`~/.config/gohan/config.yaml`:

```yaml
# Gohan configuration

theme:
  active: mocha
  auto_reload: true

backup:
  enabled: true
  auto_cleanup: true
  keep_days: 60

installation:
  auto_confirm: false
  run_preflight: true

logging:
  level: info
  file: ~/.local/share/gohan/gohan.log
```

### Database Location

SQLite database: `~/.local/share/gohan/gohan.db`

---

## Scripting Examples

### Automated Installation

```bash
#!/bin/bash
# Automated Hyprland setup

set -e

# Initialize
gohan init --force

# Run preflight checks
if ! gohan preflight run --strict; then
    echo "Preflight checks failed"
    exit 1
fi

# Install
gohan install hyprland-complete --no-confirm

# Configure
gohan post-install --display-manager sddm --auto-start

# Apply theme
gohan theme set mocha

# Verify
gohan doctor
```

### Backup Before Updates

```bash
#!/bin/bash
# Create backup before system updates

gohan backup create --description "Pre-update $(date +%Y-%m-%d)"

sudo apt update && sudo apt upgrade -y

# Verify everything still works
if ! gohan doctor; then
    echo "Issues detected, restoring backup..."
    BACKUP_ID=$(gohan backup list --limit 1 --json | jq -r '.[0].id')
    gohan backup restore "$BACKUP_ID" --force
fi
```

---

## Related Documentation

- [Installation Guide](installation.md) - Setup instructions
- [Configuration Management](configuration-management.md) - Config deployment
- [Backup & Restore](backup-restore.md) - Backup management
- [Theme Management](theme-management.md) - Theme system
- [Troubleshooting](troubleshooting.md) - Common issues
