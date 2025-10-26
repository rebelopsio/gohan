# Installation Domain Model

## Overview
The installation domain handles the complete lifecycle of installing Hyprland, from initial setup through completion, including progress tracking, error recovery, configuration preservation, and GPU-specific configurations.

## Ubiquitous Language

### Core Concepts
- **Installation Session** - A complete attempt to install Hyprland
- **System Snapshot** - Captured state of the system before changes
- **Component** - Individual piece of software to install (Hyprland, Waybar, etc.)
- **GPU Support** - Configuration specific to detected graphics hardware
- **Installation Phase** - Distinct step in the installation process
- **Rollback** - Restoring system to pre-installation state
- **Package Conflict** - When existing packages prevent installation
- **Configuration Merge** - Combining new configs with existing user configs
- **Backup** - Preserved copy of configuration files
- **Progress** - Current state and completion percentage of installation
- **Verification** - Checking that installation was successful

## Domain Structure

### Aggregate Root
- **InstallationSession** - Coordinates entire installation lifecycle

### Entities
- **SystemSnapshot** - Pre-installation system state
- **InstalledComponent** - Successfully installed component

### Value Objects
- **InstallationConfiguration** - What to install and how
- **GPUSupport** - GPU-specific requirements
- **ComponentSelection** - Selected component with version
- **InstallationProgress** - Current progress state
- **DiskSpace** - Disk space measurement
- **PackageInfo** - Package metadata
- **ConfigurationFile** - Config file information
- **PackageConflict** - Conflict information
- **ConflictResolution** - Resolution strategy
- **BackupLocation** - Backup storage location

### Domain Services
- **InstallationOrchestrator** - Coordinates installation phases
- **ConflictResolver** - Resolves package conflicts
- **ProgressEstimator** - Calculates time remaining
- **ConfigurationMerger** - Merges configurations

### Key Invariants

1. **Installation Session Integrity**
   - Cannot start if disk space insufficient
   - Must create system snapshot before making changes
   - Failed sessions must rollback or mark as requiring manual cleanup

2. **Component Dependencies**
   - Core Hyprland component always required
   - GPU drivers must match detected hardware
   - Dependencies automatically included

3. **Configuration Preservation**
   - Existing user configs never deleted without backup
   - Backups must be verified before proceeding
   - Backup location must be accessible for restoration

4. **Progress Tracking**
   - Progress only moves forward (0-100%)
   - Phases execute in defined order
   - Time estimates adjust based on actual performance

5. **Rollback Safety**
   - System snapshot must be valid and complete
   - Verify rollback success before marking complete
   - Partial rollbacks are not allowed

## Domain Events

- InstallationStartedEvent
- InstallationProgressUpdatedEvent
- PhaseCompletedEvent
- ComponentInstalledEvent
- InstallationCompletedEvent
- InstallationFailedEvent
- RollbackStartedEvent
- RollbackCompletedEvent
- ConflictDetectedEvent
- BackupCreatedEvent
- DiskSpaceInsufficientEvent
- NetworkInterruptionEvent

## References
- Feature File: docs/features/hyprland-installation.feature
- Preflight Domain: internal/domain/preflight/
