# Backup & Restore

Gohan provides comprehensive backup and restore capabilities to protect your Hyprland configurations. Every configuration change automatically creates a timestamped backup, and you can manually create, restore, and manage backups anytime.

## Overview

**Key Features:**

- ✅ **Automatic Backups** - Created before every config deployment
- ✅ **Timestamped Snapshots** - Never lose track of when backups were created
- ✅ **One-Command Restore** - Rollback to any previous state instantly
- ✅ **Selective Restore** - Restore specific files or entire snapshots
- ✅ **Cleanup Tools** - Manage old backups and free disk space
- ✅ **Backup Metadata** - Track what changed and when

## Quick Start

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

## Backup Commands

### List Backups

View all available backups:

```bash
gohan backup list
```

**Example output:**

```
Configuration Backups

┌─────────────────┬─────────────────────┬──────────┬──────────────┐
│ Backup ID       │ Created             │ Files    │ Size         │
├─────────────────┼─────────────────────┼──────────┼──────────────┤
│ 20241030_153022 │ 2024-10-30 15:30:22 │ 5 files  │ 24.5 KB      │
│ 20241030_120000 │ 2024-10-30 12:00:00 │ 5 files  │ 24.3 KB      │
│ 20241029_183015 │ 2024-10-29 18:30:15 │ 4 files  │ 18.2 KB      │
│ 20241029_090000 │ 2024-10-29 09:00:00 │ 5 files  │ 24.1 KB      │
└─────────────────┴─────────────────────┴──────────┴──────────────┘

Total: 4 backups (91.1 KB)
Backup location: ~/.local/share/gohan/backups
```

**Flags:**

| Flag | Description | Example |
|------|-------------|---------|
| `--json` | Output in JSON format | `gohan backup list --json` |
| `--limit` | Limit number of results | `gohan backup list --limit 10` |
| `--sort` | Sort by date (asc/desc) | `gohan backup list --sort desc` |

### Create Backup

Create a manual backup:

```bash
gohan backup create
```

**Example output:**

```
Creating configuration backup...

📦 Backup created successfully!

Backup ID:   20241030_153022
Path:        ~/.local/share/gohan/backups/20241030_153022
Files:       5 backed up
Size:        24.5 KB

Backed up files:
  ✓ .config/hypr/hyprland.conf
  ✓ .config/waybar/config.jsonc
  ✓ .config/waybar/style.css
  ✓ .config/kitty/kitty.conf
  ✓ .config/fuzzel/fuzzel.ini
```

**Flags:**

| Flag | Description | Example |
|------|-------------|---------|
| `--description` | Add description | `gohan backup create --description "Before major changes"` |
| `--components` | Backup specific components | `gohan backup create --components hyprland,waybar` |

### Restore Backup

Restore configurations from a backup:

```bash
gohan backup restore <backup-id>
```

**Example:**

```bash
gohan backup restore 20241030_120000
```

**Output:**

```
Restoring backup: 20241030_120000

⚠️  This will overwrite your current configurations!

Files to restore:
  • ~/.config/hypr/hyprland.conf
  • ~/.config/waybar/config.jsonc
  • ~/.config/waybar/style.css
  • ~/.config/kitty/kitty.conf
  • ~/.config/fuzzel/fuzzel.ini

Continue? (y/N): y

Creating safety backup of current state...
Safety backup: 20241030_153500

Restoring files...
  ✓ Restored .config/hypr/hyprland.conf
  ✓ Restored .config/waybar/config.jsonc
  ✓ Restored .config/waybar/style.css
  ✓ Restored .config/kitty/kitty.conf
  ✓ Restored .config/fuzzel/fuzzel.ini

✅ Restore completed successfully!

Reloading Hyprland...
✓ Hyprland reloaded

Restarting Waybar...
✓ Waybar restarted
```

**Flags:**

| Flag | Description | Example |
|------|-------------|---------|
| `--force` | Skip confirmation prompt | `gohan backup restore <id> --force` |
| `--no-reload` | Don't reload services | `gohan backup restore <id> --no-reload` |
| `--components` | Restore specific files only | `gohan backup restore <id> --components waybar` |

### Delete Backup

Remove a specific backup:

```bash
gohan backup delete <backup-id>
```

**Example:**

```bash
gohan backup delete 20241029_090000
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--force` | Skip confirmation |

### Cleanup Backups

Remove old backups to free disk space:

```bash
gohan backup cleanup [flags]
```

**Examples:**

```bash
# Remove backups older than 30 days
gohan backup cleanup --older-than 30d

# Keep only the 10 most recent backups
gohan backup cleanup --keep 10

# Remove all but the last 5
gohan backup cleanup --keep 5

# Dry run to see what would be deleted
gohan backup cleanup --older-than 30d --dry-run
```

**Flags:**

| Flag | Description | Example |
|------|-------------|---------|
| `--older-than` | Remove backups older than duration | `--older-than 30d` |
| `--keep` | Keep N most recent backups | `--keep 10` |
| `--dry-run` | Preview without deleting | `--dry-run` |
| `--force` | Skip confirmation | `--force` |

**Duration formats:**
- `7d` - 7 days
- `2w` - 2 weeks
- `3m` - 3 months
- `1y` - 1 year

### Show Backup Details

View detailed information about a backup:

```bash
gohan backup show <backup-id>
```

**Example output:**

```
Backup Details: 20241030_120000

Created:     2024-10-30 12:00:00
Location:    ~/.local/share/gohan/backups/20241030_120000
Total Size:  24.3 KB
Files:       5

Files in backup:
┌───────────────────────────────────────┬──────────┬─────────────────────┐
│ File                                   │ Size     │ Modified            │
├───────────────────────────────────────┼──────────┼─────────────────────┤
│ .config/hypr/hyprland.conf            │ 4.2 KB   │ 2024-10-30 11:59:58 │
│ .config/waybar/config.jsonc           │ 3.9 KB   │ 2024-10-30 11:59:58 │
│ .config/waybar/style.css              │ 2.8 KB   │ 2024-10-30 11:59:58 │
│ .config/kitty/kitty.conf              │ 4.0 KB   │ 2024-10-30 11:59:58 │
│ .config/fuzzel/fuzzel.ini             │ 0.6 KB   │ 2024-10-30 11:59:58 │
└───────────────────────────────────────┴──────────┴─────────────────────┘

Metadata:
  Created by:  gohan config deploy
  Theme:       mocha (Catppuccin Mocha)
  Components:  hyprland, waybar, kitty, fuzzel
```

## Automatic Backups

Backups are automatically created when:

1. **Configuration Deployment**
   ```bash
   gohan config deploy
   # ✓ Backup created: 20241030_153022
   ```

2. **Theme Changes**
   ```bash
   gohan theme set latte
   # ✓ Backup created before theme application
   ```

3. **Manual Trigger**
   ```bash
   gohan backup create
   # ✓ Manual backup: 20241030_154500
   ```

### Disable Automatic Backups

```bash
# Skip backup during deployment
gohan config deploy --skip-backup

# Skip backup during theme change
gohan theme set mocha --skip-backup
```

!!! warning "Use with Caution"
    Skipping backups means you cannot rollback changes. Only use `--skip-backup` when you're absolutely certain about the changes.

## Backup Storage

### Directory Structure

```
~/.local/share/gohan/backups/
├── 20241030_153022/
│   ├── .config/
│   │   ├── hypr/
│   │   │   └── hyprland.conf
│   │   ├── waybar/
│   │   │   ├── config.jsonc
│   │   │   └── style.css
│   │   ├── kitty/
│   │   │   └── kitty.conf
│   │   └── fuzzel/
│   │       └── fuzzel.ini
│   └── metadata.json
├── 20241030_120000/
│   ├── .config/
│   │   └── [files]
│   └── metadata.json
└── [other backups]/
```

### Metadata File

Each backup includes a `metadata.json` file:

```json
{
  "backup_id": "20241030_153022",
  "created_at": "2024-10-30T15:30:22Z",
  "created_by": "gohan config deploy",
  "description": "Automatic backup before deployment",
  "theme": {
    "name": "mocha",
    "display_name": "Catppuccin Mocha"
  },
  "components": ["hyprland", "waybar", "kitty", "fuzzel"],
  "files": [
    {
      "path": ".config/hypr/hyprland.conf",
      "size": 4294,
      "checksum": "sha256:abc123..."
    }
  ],
  "total_size": 25088,
  "version": "1.0"
}
```

## Use Cases

### Daily Workflow

Backup before experimenting:

```bash
# Morning: Create a working baseline
gohan backup create --description "Working config before experiments"

# Experiment with themes
gohan theme set latte
gohan theme set frappe

# Revert to morning baseline if needed
gohan backup list
gohan backup restore 20241030_090000
```

### Before Major Changes

```bash
# Create annotated backup
gohan backup create --description "Before Hyprland 0.42 update"

# Make changes
sudo apt update && sudo apt upgrade hyprland

# Test
hyprctl reload

# Restore if problems occur
gohan backup restore <backup-id>
```

### Regular Maintenance

Set up a cleanup schedule:

```bash
# Weekly: Keep last 2 months of backups
gohan backup cleanup --older-than 60d

# Monthly: Keep last 20 backups
gohan backup cleanup --keep 20
```

### Disaster Recovery

Complete system restore:

```bash
# After fresh install or major corruption
gohan backup list

# Restore latest working backup
gohan backup restore <latest-backup-id>

# Verify
gohan doctor
```

## Advanced Features

### Selective Restore

Restore only specific components:

```bash
# Restore only Waybar config
gohan backup restore 20241030_120000 --components waybar

# Restore multiple components
gohan backup restore 20241030_120000 --components hyprland,kitty
```

### Backup Comparison

Compare two backups (coming soon):

```bash
# Not yet implemented
gohan backup diff 20241030_120000 20241030_153022
```

### Export/Import Backups

Transfer backups between systems (coming soon):

```bash
# Export backup as archive
gohan backup export 20241030_120000 --output ~/backup.tar.gz

# Import on another system
gohan backup import ~/backup.tar.gz
```

## Troubleshooting

### Backup Failed - No Disk Space

**Problem:** `Error: no space left on device`

**Solution:**
```bash
# Check disk usage
df -h ~/.local/share/gohan

# Clean up old backups
gohan backup cleanup --older-than 30d

# Check remaining space
df -h ~/.local/share/gohan
```

### Cannot Restore - Checksum Mismatch

**Problem:** `Error: backup corrupted (checksum mismatch)`

**Solution:**
```bash
# List available backups
gohan backup list

# Try a different backup
gohan backup restore <different-backup-id>

# If all backups are corrupted, redeploy
gohan config deploy
```

### Backup Directory Not Found

**Problem:** `Error: backup directory does not exist`

**Solution:**
```bash
# Create backup directory
mkdir -p ~/.local/share/gohan/backups

# Verify permissions
ls -la ~/.local/share/gohan/

# Create a new backup
gohan backup create
```

### Restore Didn't Reload Services

**Problem:** Changes restored but Hyprland didn't reload

**Solution:**
```bash
# Manually reload
hyprctl reload

# Restart Waybar
pkill waybar && waybar &

# Or restore with explicit reload
gohan backup restore <backup-id>
```

## Best Practices

### 1. Regular Manual Backups

Create manual snapshots before risky operations:

```bash
# Before system updates
gohan backup create --description "Before system update"

# Before experimenting
gohan backup create --description "Working baseline"
```

### 2. Use Descriptive Names

Help future-you understand what the backup contains:

```bash
gohan backup create --description "Perfect Mocha setup for work"
gohan backup create --description "Gaming config with performance tweaks"
```

### 3. Regular Cleanup

Prevent backup directory bloat:

```bash
# Add to cron (monthly cleanup)
0 0 1 * * gohan backup cleanup --older-than 60d
```

### 4. Test Restores

Periodically verify backups work:

```bash
# Test restore (creates safety backup first)
gohan backup restore <backup-id>

# Verify everything works
gohan doctor

# Roll forward to latest if needed
gohan backup list
gohan backup restore <latest>
```

### 5. Document Important Backups

Keep notes about significant backups:

```bash
gohan backup create --description "Pre-production setup - DO NOT DELETE"
```

## Backup Size Management

### Typical Backup Sizes

| Component | Typical Size |
|-----------|-------------|
| Hyprland config | 3-5 KB |
| Waybar config + style | 5-8 KB |
| Kitty config | 3-4 KB |
| Fuzzel config | 0.5-1 KB |
| **Total per backup** | **~15-25 KB** |

### Disk Space Estimates

| Backups | Approximate Size |
|---------|-----------------|
| 10 backups | ~200 KB |
| 50 backups | ~1 MB |
| 100 backups | ~2 MB |
| 365 backups (daily for 1 year) | ~8 MB |

!!! tip "Space-Efficient"
    Gohan backups are extremely space-efficient. Even with daily backups for a year, you'll only use about 8 MB of disk space.

## Related Documentation

- [Configuration Management](configuration-management.md) - Config deployment
- [Theme Management](theme-management.md) - Theme switching
- [Troubleshooting](troubleshooting.md) - Common issues
- [CLI Reference](cli-reference.md) - All commands
