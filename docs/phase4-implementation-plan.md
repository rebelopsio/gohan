# Phase 4: Theme System - Implementation Plan

## Overview

Phase 4 implements a comprehensive theme management system inspired by Omarchy's successful theme switching architecture. This phase enables users to switch between pre-configured themes with a single command, applying consistent styling across all Hyprland ecosystem components.

## Goals

1. **Theme Management**: Domain model for themes with metadata and discovery
2. **Theme Switching**: Apply themes across Hyprland, Waybar, Kitty, etc.
3. **CLI Integration**: User-friendly commands for theme operations
4. **Safety First**: Backup before switching, rollback on failure
5. **Hot Reload**: Apply themes without logout/restart where possible

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ PHASE 4: THEME SYSTEM                                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────────┐     ┌──────────────────┐            │
│  │  Theme Domain    │     │  Theme Service   │            │
│  │  - Theme         │────▶│  - List themes   │            │
│  │  - ThemeMetadata │     │  - Apply theme   │            │
│  │  - ThemeRegistry │     │  - Preview theme │            │
│  │  - ColorScheme   │     │  - Get active    │            │
│  └──────────────────┘     └──────────────────┘            │
│           │                         │                      │
│           │                         ▼                      │
│           │              ┌──────────────────┐             │
│           │              │  Theme Applier   │             │
│           │              │  - Component map │             │
│           └─────────────▶│  - Hot reload    │             │
│                          │  - Rollback      │             │
│                          └──────────────────┘             │
│                                   │                        │
│                                   ▼                        │
│                    ┌──────────────────────────┐           │
│                    │  Uses Phase 3 Components │           │
│                    │  - ConfigDeployer        │           │
│                    │  - TemplateEngine        │           │
│                    │  - BackupService         │           │
│                    └──────────────────────────┘           │
│                                                             │
│  ┌─────────────────────────────────────────┐              │
│  │  CLI Commands                           │              │
│  │  - gohan theme list                     │              │
│  │  - gohan theme set <name>               │              │
│  │  - gohan theme preview <name>           │              │
│  │  - gohan theme show                     │              │
│  └─────────────────────────────────────────┘              │
└─────────────────────────────────────────────────────────────┘
```

## Core Themes (Catppuccin Family)

1. **Mocha** (Dark) - Current default
2. **Latte** (Light)
3. **Frappe** (Dark, muted)
4. **Macchiato** (Dark, warmer)
5. **Gohan** (Custom brand theme)

Each theme includes color definitions for:
- Hyprland (borders, shadows, colors)
- Waybar (background, foreground, modules)
- Kitty (terminal colors)
- Rofi/Fuzzel (menu styling)

## Component Breakdown

### Phase 4.1: Theme Domain Models

**Location:** `internal/domain/theme/`

**Domain Types:**

```go
// Theme represents a complete visual theme
type Theme struct {
    name        ThemeName
    metadata    ThemeMetadata
    colorScheme ColorScheme
    createdAt   time.Time
}

type ThemeName string

const (
    ThemeMocha      ThemeName = "mocha"
    ThemeLatte      ThemeName = "latte"
    ThemeFrappe     ThemeName = "frappe"
    ThemeMacchiato  ThemeName = "macchiato"
    ThemeGohan      ThemeName = "gohan"
)

// ThemeMetadata contains theme information
type ThemeMetadata struct {
    displayName string
    author      string
    description string
    variant     ThemeVariant // "dark" or "light"
    previewURL  string
}

type ThemeVariant string

const (
    ThemeVariantDark  ThemeVariant = "dark"
    ThemeVariantLight ThemeVariant = "light"
)

// ColorScheme defines all colors in a theme
type ColorScheme struct {
    // Base colors
    base        Color
    surface     Color
    overlay     Color
    text        Color
    subtext     Color

    // Accent colors
    rosewater   Color
    flamingo    Color
    pink        Color
    mauve       Color
    red         Color
    maroon      Color
    peach       Color
    yellow      Color
    green       Color
    teal        Color
    sky         Color
    sapphire    Color
    blue        Color
    lavender    Color
}

type Color string // Hex color code

// ThemeRegistry manages available themes
type ThemeRegistry interface {
    Register(theme *Theme) error
    FindByName(name ThemeName) (*Theme, error)
    ListAll() []*Theme
    GetActive() (*Theme, error)
    SetActive(name ThemeName) error
}
```

**Business Rules:**
1. Theme names must be unique
2. Color schemes must include all required colors
3. Only one theme can be active at a time
4. System must have at least one theme (default to Mocha)

### Phase 4.2: Theme Service

**Location:** `internal/application/theme/`

**Use Cases:**

```go
// ListThemesUseCase lists all available themes
type ListThemesUseCase struct {
    registry ThemeRegistry
}

func (uc *ListThemesUseCase) Execute(ctx context.Context) ([]*dto.ThemeInfo, error)

// ApplyThemeUseCase applies a theme system-wide
type ApplyThemeUseCase struct {
    registry       ThemeRegistry
    themeApplier   ThemeApplier
    backupService  BackupService
}

func (uc *ApplyThemeUseCase) Execute(
    ctx context.Context,
    themeName string,
    progressChan chan<- ThemeProgress,
) error

// GetActiveThemeUseCase returns current active theme
type GetActiveThemeUseCase struct {
    registry ThemeRegistry
}

func (uc *GetActiveThemeUseCase) Execute(ctx context.Context) (*dto.ThemeInfo, error)

// PreviewThemeUseCase generates preview of theme without applying
type PreviewThemeUseCase struct {
    registry ThemeRegistry
}

func (uc *PreviewThemeUseCase) Execute(
    ctx context.Context,
    themeName string,
) (*dto.ThemePreview, error)
```

**DTO Types:**

```go
type ThemeInfo struct {
    Name        string
    DisplayName string
    Author      string
    Description string
    Variant     string // "dark" or "light"
    IsActive    bool
    PreviewURL  string
}

type ThemeProgress struct {
    Component       string
    Status          string // "started", "applying", "completed", "failed"
    PercentComplete float64
    Error           error
}

type ThemePreview struct {
    Name        string
    ColorScheme map[string]string // color name -> hex value
    Preview     string            // ASCII art or emoji visualization
}
```

### Phase 4.3: Theme Applier

**Location:** `internal/infrastructure/theme/`

**Implementation:**

```go
// ThemeApplier applies themes to system components
type ThemeApplier struct {
    configDeployer *configservice.ConfigDeployer
    templateEngine *templates.TemplateEngine
    backupService  *backup.BackupService
}

// ComponentConfig maps theme to component-specific config
type ComponentConfig struct {
    Component    string // "hyprland", "waybar", "kitty", etc.
    TemplatePath string
    TargetPath   string
    HotReloadCmd []string // Command to reload without restart
}

func (a *ThemeApplier) Apply(
    ctx context.Context,
    theme *theme.Theme,
    progressChan chan<- ThemeProgress,
) error

func (a *ThemeApplier) Rollback(ctx context.Context, backupID string) error

func (a *ThemeApplier) GetComponentConfigs() []ComponentConfig
```

**Components Affected:**
1. Hyprland colors and borders
2. Waybar styling
3. Kitty terminal colors
4. Rofi/Fuzzel theme
5. Hyprlock colors
6. Future: mako, dunst notifications

### Phase 4.4: CLI Commands

**Location:** `internal/cli/cmd/`

**Commands:**

```bash
# List all available themes
gohan theme list
gohan theme list --format json

# Show current active theme
gohan theme show
gohan theme show --verbose

# Apply a theme
gohan theme set mocha
gohan theme set latte --no-backup
gohan theme set frappe --preview

# Preview theme without applying
gohan theme preview macchiato
gohan theme preview gohan --colors

# Rollback to previous theme
gohan theme rollback
gohan theme rollback --backup-id <id>
```

**Output Examples:**

```
$ gohan theme list

Available Themes:

  ● mocha       Catppuccin Mocha (Dark)      [ACTIVE]
    latte       Catppuccin Latte (Light)
    frappe      Catppuccin Frappe (Dark)
    macchiato   Catppuccin Macchiato (Dark)
    gohan       Gohan Default (Dark)

Use 'gohan theme set <name>' to switch themes

$ gohan theme set latte

Switching to theme: Latte
Creating backup...                    [✓]
Applying to Hyprland...              [✓]
Applying to Waybar...                [✓]
Applying to Kitty...                 [✓]
Applying to Rofi...                  [✓]
Reloading components...              [✓]

Theme 'latte' applied successfully!
Backup ID: 2025-10-27_190000
Use 'gohan theme rollback' to revert if needed.
```

## Test Strategy

### BDD Features

**Feature:** Theme Listing
```gherkin
Scenario: User lists available themes
  Given the theme system is initialized
  When I run "gohan theme list"
  Then I should see at least 5 themes
  And one theme should be marked as active
  And each theme should show its variant (dark/light)
```

**Feature:** Theme Switching
```gherkin
Scenario: User switches to a different theme
  Given I am using the "mocha" theme
  When I run "gohan theme set latte"
  Then the active theme should be "latte"
  And all configuration files should be updated
  And a backup should be created
  And components should be reloaded
```

**Feature:** Theme Rollback
```gherkin
Scenario: User rolls back after theme change
  Given I switched from "mocha" to "latte"
  And a backup was created
  When I run "gohan theme rollback"
  Then the active theme should be "mocha"
  And all configurations should be restored
```

### Unit Tests (TDD)

**Phase 4.1 Tests:**
- Theme creation and validation
- ThemeMetadata validation
- ColorScheme parsing
- ThemeRegistry operations (register, find, list, set active)

**Phase 4.2 Tests:**
- ListThemesUseCase execution
- ApplyThemeUseCase with progress reporting
- GetActiveThemeUseCase
- PreviewThemeUseCase color rendering

**Phase 4.3 Tests:**
- ThemeApplier component mapping
- Hot reload command execution
- Rollback functionality
- Template variable generation from ColorScheme

**Phase 4.4 Tests:**
- CLI command parsing
- Output formatting (table, JSON)
- Error handling and user guidance

## Implementation Phases

### Phase 4.1: Domain Models (TDD)
**Time:** ~30 minutes
- Create theme domain types
- Implement ThemeRegistry with in-memory storage
- Define 5 core themes with color schemes
- Write 15+ unit tests

**Commit:** `feat(phase4.1): implement theme domain models`

### Phase 4.2: Theme Use Cases (TDD)
**Time:** ~45 minutes
- Implement ListThemesUseCase
- Implement ApplyThemeUseCase
- Implement GetActiveThemeUseCase
- Implement PreviewThemeUseCase
- Write 20+ unit tests

**Commit:** `feat(phase4.2): implement theme use cases`

### Phase 4.3: Theme Applier (TDD)
**Time:** ~45 minutes
- Create ThemeApplier infrastructure
- Map themes to component configs
- Integrate with ConfigDeployer from Phase 3
- Implement hot reload
- Write 15+ unit tests

**Commit:** `feat(phase4.3): implement theme applier`

### Phase 4.4: CLI Commands (TDD)
**Time:** ~30 minutes
- Implement `theme list` command
- Implement `theme set` command
- Implement `theme show` command
- Implement `theme preview` command
- Wire into dependency injection container
- Write CLI integration tests

**Commit:** `feat(phase4.4): implement theme CLI commands`

### Phase 4.5: Documentation & Testing
**Time:** ~30 minutes
- Update README with theme commands
- Create theme documentation
- Run full test suite
- Manual testing
- Update phase4-status.md

**Commit:** `docs(phase4): complete Phase 4 documentation`

## Success Criteria

- [ ] 5 themes registered and discoverable
- [ ] Theme switching works end-to-end
- [ ] Backups created before theme changes
- [ ] Hot reload works for supported components
- [ ] CLI commands functional and user-friendly
- [ ] 60+ unit tests passing
- [ ] Full integration test passing
- [ ] Documentation complete

## Dependencies

**From Phase 3:**
- ✅ ConfigDeployer (configuration deployment)
- ✅ TemplateEngine (variable substitution)
- ✅ BackupService (backup/restore)

**New Templates Needed:**
- Theme-specific template variants for each component
- Template variables for all ColorScheme colors

## Technical Considerations

1. **Color Format:** Use hex colors (#RRGGBB) consistently
2. **Hot Reload:** Not all components support hot reload, may need restart
3. **Symlinks:** Consider Omarchy's symlink approach for theme switching
4. **Preview:** ASCII art or emoji-based color preview for CLI
5. **Persistence:** Store active theme in ~/.config/gohan/theme.json
6. **Validation:** Ensure theme changes are atomic (all-or-nothing)

## Timeline

**Total Estimated Time:** ~3 hours following BDD → ATDD → TDD

**Breakdown:**
- BDD Feature Files: 15 minutes
- BDD Expert Review: 10 minutes
- ATDD Conversion: 20 minutes
- Phase 4.1: 30 minutes
- Phase 4.2: 45 minutes
- Phase 4.3: 45 minutes
- Phase 4.4: 30 minutes
- Documentation: 30 minutes
- Buffer: 15 minutes

---

**Author:** Claude Code (claude-sonnet-4-5)
**Date:** 2025-10-27
**Status:** Planning Phase
