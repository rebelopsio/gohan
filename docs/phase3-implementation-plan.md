# Phase 3: Package Installation & Configuration Deployment

## Implementation Plan

### Overview

Phase 3 brings Gohan to life by implementing actual package installation and configuration deployment. This phase transforms our package definitions and templates into a working system.

### Goals

1. ✅ Install Debian packages using APT
2. ✅ Deploy configuration files to user directories
3. ✅ Backup existing configurations before overwriting
4. ✅ Handle template variable substitution
5. ✅ Provide progress reporting during installation
6. ✅ Implement error handling and rollback capabilities

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ PHASE 3 ARCHITECTURE                                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐      ┌──────────────┐                   │
│  │ Installation │      │ Configuration│                   │
│  │ Use Case     │─────▶│ Deployment   │                   │
│  └──────────────┘      │ Use Case     │                   │
│         │              └──────────────┘                   │
│         │                      │                           │
│         ▼                      ▼                           │
│  ┌──────────────┐      ┌──────────────┐                   │
│  │ APT Package  │      │ Config File  │                   │
│  │ Manager      │      │ Service      │                   │
│  └──────────────┘      └──────────────┘                   │
│         │                      │                           │
│         │                      ├─────▶ Template Engine    │
│         │                      ├─────▶ Backup Service     │
│         │                      └─────▶ File Copier        │
│         │                                                  │
│         ▼                                                  │
│  ┌──────────────┐                                         │
│  │   APT CLI    │                                         │
│  └──────────────┘                                         │
└─────────────────────────────────────────────────────────────┘
```

## Component Breakdown

### 1. Enhanced APT Package Manager

**Location**: `internal/infrastructure/installation/packagemanager/`

**What exists**: Basic APT manager with install/remove/check
**What we need**: Batch installation with progress reporting

**New Methods**:
```go
// Install multiple packages with progress reporting
InstallPackages(ctx context.Context, packages []string, progressChan chan<- Progress) error

// Install package profile (minimal/recommended/full)
InstallProfile(ctx context.Context, profile installation.InstallationProfile, progressChan chan<- Progress) error

// Check if all packages in list are installed
ArePackagesInstalled(ctx context.Context, packages []string) (map[string]bool, error)
```

**Files to Create**:
- ✅ Update `apt_manager.go` with batch installation
- ✅ Create `apt_manager_progress.go` for progress reporting
- ✅ Update tests in `apt_manager_test.go`

### 2. Configuration Deployment Service

**Location**: `internal/infrastructure/installation/configservice/`

**Purpose**: Deploy template configurations to user directories

**Key Responsibilities**:
- Copy configuration files from templates to target locations
- Handle directory creation
- Set appropriate file permissions
- Report progress

**Interface**:
```go
type ConfigurationDeployer interface {
    DeployConfigurations(ctx context.Context, configs []ConfigurationFile) error
    GetTargetPath(template string, username string) (string, error)
}
```

**Files to Create**:
- ✅ `config_deployer.go` - Main deployment logic
- ✅ `config_deployer_test.go` - Comprehensive tests
- ✅ `file_copier.go` - Safe file copying with verification

### 3. Backup Service

**Location**: `internal/infrastructure/installation/backup/`

**Purpose**: Backup existing configurations before overwriting

**Key Features**:
- Create timestamped backups
- Support selective backup (only files that will be overwritten)
- Provide restore capability
- Track backup metadata

**Interface**:
```go
type BackupService interface {
    BackupFile(ctx context.Context, filePath string) (backupPath string, error)
    BackupDirectory(ctx context.Context, dirPath string) (backupPath string, error)
    RestoreBackup(ctx context.Context, backupPath string) error
    ListBackups(ctx context.Context) ([]BackupMetadata, error)
}
```

**Backup Structure**:
```
~/.config/gohan/backups/
├── 2025-10-27_143022/          # Timestamp
│   ├── manifest.json           # What was backed up
│   ├── hypr/                   # Original configs
│   │   ├── hyprland.conf
│   │   └── bindings.conf
│   ├── waybar/
│   └── kitty/
```

**Files to Create**:
- ✅ `backup_service.go` - Backup creation logic
- ✅ `backup_metadata.go` - Metadata tracking
- ✅ `backup_service_test.go` - Tests
- ✅ `restore.go` - Restore functionality

### 4. Template Engine

**Location**: `internal/infrastructure/installation/templates/`

**Purpose**: Replace variables in configuration files

**Variables Supported**:
```
{{username}}       → Current user's username
{{home}}           → User's home directory
{{config_dir}}     → ~/.config
{{hostname}}       → System hostname
{{display}}        → Primary display name
{{resolution}}     → Primary display resolution
```

**Interface**:
```go
type TemplateEngine interface {
    ProcessTemplate(content string, vars TemplateVars) (string, error)
    ProcessFile(srcPath, dstPath string, vars TemplateVars) error
}

type TemplateVars struct {
    Username   string
    Home       string
    ConfigDir  string
    Hostname   string
    Display    string
    Resolution string
}
```

**Files to Create**:
- ✅ `template_engine.go` - Variable substitution logic
- ✅ `template_vars.go` - Variable collection from system
- ✅ `template_engine_test.go` - Tests

### 5. Domain Models

**Location**: `internal/domain/configuration/`

**New Value Objects**:

**ConfigurationFile**:
```go
type ConfigurationFile struct {
    TemplatePath string       // Source template path
    TargetPath   string       // Where to deploy
    Permissions  os.FileMode  // File permissions
    Owner        string       // File owner
    Backup       bool         // Whether to backup before overwriting
}
```

**DeploymentResult**:
```go
type DeploymentResult struct {
    FilePath       string
    Success        bool
    BackupPath     string
    Error          error
    ProcessedLines int
}
```

**Files to Create**:
- ✅ `configuration_file.go` - ConfigurationFile value object
- ✅ `deployment_result.go` - Result tracking
- ✅ Tests for both

## Implementation Phases

### Phase 3.1: APT Package Manager Enhancement

**Goal**: Install packages with progress reporting

**Tasks**:
1. ✅ Update `apt_manager.go` with batch installation
2. ✅ Add progress reporting via channels
3. ✅ Handle package installation failures gracefully
4. ✅ Implement transaction-like behavior (all or nothing)
5. ✅ Add comprehensive error messages
6. ✅ Update tests

**Acceptance Criteria**:
- Can install list of packages
- Progress reported for each package
- Failures don't leave system in broken state
- All tests pass

### Phase 3.2: Template Engine

**Goal**: Variable substitution in config files

**Tasks**:
1. ✅ Implement basic variable substitution
2. ✅ Collect system variables (username, home, etc.)
3. ✅ Handle missing variables gracefully
4. ✅ Support nested variables ({{config_dir}}/hypr → /home/user/.config/hypr)
5. ✅ Add tests with various edge cases

**Acceptance Criteria**:
- All variables replaced correctly
- Handles missing variables
- Works with real config files
- All tests pass

### Phase 3.3: Backup Service

**Goal**: Safe backup of existing configurations

**Tasks**:
1. ✅ Implement backup directory creation
2. ✅ Create backup manifest (JSON)
3. ✅ Copy files to backup location
4. ✅ Implement restore functionality
5. ✅ Add cleanup for old backups (keep last N)
6. ✅ Add tests

**Acceptance Criteria**:
- Creates timestamped backups
- Manifest tracks all backed up files
- Restore works correctly
- Old backups cleaned up
- All tests pass

### Phase 3.4: Configuration Deployment Service

**Goal**: Deploy templates to user directories

**Tasks**:
1. ✅ Implement file copying with backup
2. ✅ Create target directories if missing
3. ✅ Set appropriate permissions
4. ✅ Process templates through template engine
5. ✅ Report deployment progress
6. ✅ Handle errors (file in use, permissions, etc.)
7. ✅ Add comprehensive tests

**Acceptance Criteria**:
- Copies all config files correctly
- Creates necessary directories
- Backups created before overwriting
- Templates processed
- Progress reported
- All tests pass

### Phase 3.5: Integration

**Goal**: Wire everything into installation use case

**Tasks**:
1. ✅ Update `execute_installation.go` use case
2. ✅ Add configuration deployment phase
3. ✅ Integrate progress reporting (15-80% packages, 80-95% configs)
4. ✅ Handle errors and rollback
5. ✅ Update container for dependency injection
6. ✅ Update integration tests

**Acceptance Criteria**:
- Full installation works end-to-end
- Packages installed correctly
- Configs deployed correctly
- Progress reporting accurate
- Rollback works on failure
- All tests pass

## Progress Allocation

Total installation progress: 0-100%

```
0-15%    Preflight checks (existing)
15-20%   Update APT cache
20-75%   Install packages (55% total)
         - Each package gets equal share
75-80%   Verify package installation
80-90%   Deploy configurations (10% total)
         - Template processing
         - File copying
90-95%   Set permissions and ownership
95-100%  Final verification
```

## Error Handling Strategy

### Package Installation Errors

1. **Network failure**: Retry up to 3 times
2. **Package not found**: Report error, offer alternative
3. **Dependency conflict**: Attempt resolution, or abort
4. **Disk space**: Check before install, abort if insufficient
5. **Permission denied**: Report error, request sudo

### Configuration Deployment Errors

1. **File in use**: Warn user, offer to retry
2. **Permission denied**: Report error, request proper permissions
3. **Disk full**: Abort, rollback
4. **Backup failed**: Warn, offer to proceed anyway or abort

### Rollback Strategy

1. **Package installation failed**:
   - Remove successfully installed packages
   - Restore system to pre-installation state

2. **Configuration deployment failed**:
   - Restore backed up configurations
   - Keep successfully installed packages
   - Report partial success state

## Testing Strategy

### Unit Tests

- Template engine with various inputs
- Backup service with different file structures
- Config deployer with mocked file system
- APT manager with dry-run mode

### Integration Tests

- Full package installation (dry-run)
- Configuration deployment with real files (in test directory)
- Backup and restore cycle
- End-to-end installation simulation

### Manual Testing

- Install on fresh Debian Sid VM
- Test with existing configurations
- Test rollback scenarios
- Test with different user permissions

## Security Considerations

1. **File Permissions**: Ensure configs have appropriate permissions (0644 for files, 0755 for directories)
2. **Backup Security**: Backups should have same permissions as originals
3. **Template Injection**: Validate template variables to prevent injection
4. **Sudo Usage**: Minimize sudo usage, only for package installation
5. **File Overwrite**: Always backup before overwriting

## Performance Considerations

1. **Parallel Package Installation**: Install independent packages in parallel
2. **Buffered File I/O**: Use buffered reading/writing for large files
3. **Progress Reporting**: Batch progress updates to avoid overhead
4. **Backup Optimization**: Only backup files that will be overwritten

## File Mapping

### Template → Target Mapping

```
templates/hyprland/*.conf     → ~/.config/hypr/*.conf
templates/waybar/*            → ~/.config/waybar/*
templates/kitty/kitty.conf    → ~/.config/kitty/kitty.conf
templates/fuzzel/fuzzel.ini   → ~/.config/fuzzel/fuzzel.ini
```

## Success Criteria

Phase 3 is complete when:

✅ All packages from selected profile are installed
✅ All configuration files are deployed correctly
✅ Template variables are substituted properly
✅ Existing configurations are backed up safely
✅ Progress is reported accurately throughout
✅ Errors are handled gracefully with rollback
✅ All unit tests pass (100+ tests)
✅ Integration tests pass
✅ Manual testing on Debian Sid successful
✅ Documentation is complete

## Timeline Estimate

- **3.1 APT Enhancement**: 1 session
- **3.2 Template Engine**: 1 session
- **3.3 Backup Service**: 1 session
- **3.4 Config Deployment**: 1 session
- **3.5 Integration**: 1 session
- **Testing & Polish**: 1 session

**Total**: 6 development sessions

## Next Steps After Phase 3

**Phase 4: Post-Install Setup**
- Environment variable configuration
- Display manager integration (GDM/SDDM)
- Systemd service management
- User session configuration

**Phase 5: Polish & UX**
- First-run welcome screen
- System tray integration
- Wallpaper selection
- Theme customization
