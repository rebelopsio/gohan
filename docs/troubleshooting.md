# Troubleshooting Guide

This guide covers common issues and their solutions when using Gohan.

## Table of Contents
- [Installation Issues](#installation-issues)
- [Theme Issues](#theme-issues)
- [Configuration Issues](#configuration-issues)
- [Component Reload Issues](#component-reload-issues)
- [General Issues](#general-issues)

## Installation Issues

### "command not found: gohan"

**Symptoms**: Running `gohan` returns "command not found"

**Cause**: Gohan binary is not in your PATH

**Solutions**:

1. **Check installation location**:
   ```bash
   which gohan
   # If no output, gohan is not in PATH
   ```

2. **Add to PATH** (if installed to custom location):
   ```bash
   # Add to ~/.bashrc or ~/.zshrc
   export PATH="/path/to/gohan:$PATH"
   source ~/.bashrc  # or ~/.zshrc
   ```

3. **Install system-wide**:
   ```bash
   sudo cp gohan /usr/local/bin/
   ```

### "permission denied" when running gohan

**Symptoms**: `bash: ./gohan: Permission denied`

**Cause**: Binary is not executable

**Solution**:
```bash
chmod +x gohan
```

### Build fails with Go errors

**Symptoms**: `go build` fails with compilation errors

**Causes & Solutions**:

1. **Go version too old**:
   ```bash
   # Check version
   go version
   # Should be 1.21 or higher

   # Upgrade Go from https://go.dev/dl/
   ```

2. **Missing dependencies**:
   ```bash
   # Download dependencies
   go mod download

   # Tidy dependencies
   go mod tidy
   ```

3. **Build cache issues**:
   ```bash
   # Clean build cache
   go clean -cache -modcache

   # Rebuild
   go build -o gohan ./cmd/gohan
   ```

## Theme Issues

### Theme doesn't apply to all components

**Symptoms**: After applying a theme, some components still show old colors

**Causes & Solutions**:

1. **Component not running**:
   ```bash
   # Check if Waybar is running
   pgrep waybar

   # Start if not running
   waybar &
   ```

2. **Manual reload needed**:
   ```bash
   # Reload Hyprland
   hyprctl reload

   # Restart Waybar
   killall waybar && waybar &

   # Restart Kitty (open new terminal)
   ```

3. **Configuration file not updated**:
   ```bash
   # Check if config file was modified
   ls -la ~/.config/waybar/style.css

   # View recent modifications
   stat ~/.config/waybar/style.css
   ```

### Colors look wrong or washed out

**Symptoms**: Theme colors appear incorrect or faded

**Causes & Solutions**:

1. **Terminal doesn't support true color**:
   ```bash
   # Check color support
   echo $COLORTERM
   # Should output "truecolor" or "24bit"

   # For Kitty
   export COLORTERM=truecolor

   # Add to ~/.bashrc or ~/.zshrc
   ```

2. **Old configuration cached**:
   ```bash
   # Clear Kitty cache
   rm -rf ~/.cache/kitty

   # Restart Kitty
   ```

### "Theme not found" error

**Symptoms**: `gohan theme set xyz` returns "theme not found"

**Cause**: Invalid theme name

**Solution**:
```bash
# List available themes
gohan theme list

# Use exact theme name
gohan theme set mocha
```

### Theme rollback doesn't work

**Symptoms**: "No theme history available" when trying to rollback

**Causes & Solutions**:

1. **No theme changes made**:
   ```bash
   # Apply a theme first
   gohan theme set mocha
   gohan theme set latte

   # Now rollback works
   gohan theme rollback
   ```

2. **History file corrupted**:
   ```bash
   # Check history file
   cat ~/.config/gohan/theme-history.json

   # If corrupted, remove and rebuild
   rm ~/.config/gohan/theme-history.json
   gohan theme set mocha  # Starts fresh history
   ```

3. **Only one theme in history**:
   - Rollback requires at least 2 themes in history
   - Apply another theme to build history

## Configuration Issues

### Configuration file not created

**Symptoms**: Template deployed but target file doesn't exist

**Causes & Solutions**:

1. **Directory doesn't exist**:
   ```bash
   # Create target directory
   mkdir -p ~/.config/waybar

   # Reapply theme
   gohan theme set mocha
   ```

2. **Permission issues**:
   ```bash
   # Check directory permissions
   ls -la ~/.config/

   # Fix permissions
   chmod 755 ~/.config/waybar
   ```

### Backup creation fails

**Symptoms**: Warning about backup failure when applying theme

**Causes & Solutions**:

1. **Backup directory doesn't exist**:
   ```bash
   # Create backup directory
   mkdir -p ~/.config/gohan/backups

   # Set permissions
   chmod 755 ~/.config/gohan/backups
   ```

2. **Disk space full**:
   ```bash
   # Check disk space
   df -h

   # Clean old backups if needed
   rm -rf ~/.config/gohan/backups/old-*
   ```

### Template variables not replaced

**Symptoms**: Config files contain `{{.theme_base}}` instead of actual colors

**Causes & Solutions**:

1. **Template syntax error**:
   ```bash
   # Check template file for syntax errors
   # Variables must be: {{.variable_name}}
   # Not: {{ .variable_name }} (extra spaces)
   ```

2. **Missing template file**:
   ```bash
   # Check template exists
   ls -la templates/waybar/style.css.tmpl

   # If missing, template can't be processed
   ```

## Component Reload Issues

### Hyprland doesn't reload after theme change

**Symptoms**: Hyprland still shows old theme after applying new theme

**Causes & Solutions**:

1. **hyprctl not available**:
   ```bash
   # Check if hyprctl is installed
   which hyprctl

   # If not installed, Hyprland may not be running
   # or not in PATH
   ```

2. **Manual reload needed**:
   ```bash
   # Reload Hyprland manually
   hyprctl reload

   # Or restart Hyprland session
   ```

### Waybar doesn't restart

**Symptoms**: Waybar still shows old theme

**Causes & Solutions**:

1. **Waybar not running**:
   ```bash
   # Check if running
   pgrep waybar

   # Start Waybar
   waybar &
   ```

2. **Waybar crashes on restart**:
   ```bash
   # Check Waybar logs
   journalctl --user -u waybar

   # Or run Waybar in foreground to see errors
   killall waybar
   waybar
   ```

3. **Configuration syntax error**:
   ```bash
   # Validate Waybar config
   waybar --config ~/.config/waybar/config \
          --style ~/.config/waybar/style.css

   # Fix any syntax errors reported
   ```

### "failed to reload Hyprland" warning

**Symptoms**: Warning message when applying theme

**Cause**: hyprctl command failed

**Solutions**:

1. **Hyprland not running**:
   - This is expected if not using Hyprland
   - Warning can be safely ignored

2. **Permission issues**:
   ```bash
   # Check Hyprland socket permissions
   ls -la /tmp/hypr/
   ```

## General Issues

### State file corruption

**Symptoms**: Errors reading theme state, inconsistent behavior

**Solution**:
```bash
# Remove corrupt state files
rm ~/.config/gohan/theme-state.json
rm ~/.config/gohan/theme-history.json

# Reapply theme to rebuild state
gohan theme set mocha
```

### Slow theme application

**Symptoms**: Theme takes a long time to apply

**Causes & Solutions**:

1. **Large number of backups**:
   ```bash
   # Check backup directory size
   du -sh ~/.config/gohan/backups

   # Clean old backups
   find ~/.config/gohan/backups -type d -mtime +30 -exec rm -rf {} +
   ```

2. **Slow template processing**:
   - Usually fast, but complex templates may take longer
   - Consider simplifying custom templates

### Configuration not persisting across reboots

**Symptoms**: Theme resets after system restart

**Causes & Solutions**:

1. **State file not readable**:
   ```bash
   # Check file exists and is readable
   cat ~/.config/gohan/theme-state.json

   # If missing or unreadable, reapply theme
   gohan theme set mocha
   ```

2. **Components reset their configs**:
   - Some components may override configs on startup
   - Check component-specific startup configs

### "Failed to deploy configurations" error

**Symptoms**: Theme application fails with deployment error

**Causes & Solutions**:

1. **Template not found**:
   ```bash
   # Check templates exist
   ls -la templates/

   # Reinstall Gohan if templates are missing
   ```

2. **Target directory permissions**:
   ```bash
   # Check ~/.config permissions
   ls -la ~/.config

   # Fix permissions
   chmod 755 ~/.config
   chmod 755 ~/.config/hypr ~/.config/waybar
   ```

3. **Disk space full**:
   ```bash
   # Check disk space
   df -h
   ```

## Getting Help

If you've tried the solutions above and still have issues:

1. **Check GitHub Issues**: https://github.com/rebelopsio/gohan/issues
   - Search existing issues
   - Someone may have encountered the same problem

2. **Enable Verbose Output**:
   ```bash
   gohan --verbose theme set mocha
   ```

3. **Collect Debug Information**:
   ```bash
   # System information
   uname -a
   go version
   gohan version

   # Component versions
   hyprctl version
   waybar --version
   kitty --version

   # File permissions
   ls -la ~/.config/gohan
   ls -la ~/.config/hypr
   ls -la ~/.config/waybar

   # Recent logs
   journalctl --user -n 50
   ```

4. **Create an Issue**:
   - Include the debug information above
   - Describe what you were trying to do
   - Include exact error messages
   - Mention your OS and versions

5. **Join Discussions**: https://github.com/rebelopsio/gohan/discussions
   - Ask questions
   - Share experiences
   - Get community help

## Prevention Tips

1. **Keep Backups**: Gohan creates backups, but maintain your own for important configs
2. **Test Changes**: Use `gohan theme preview` before applying
3. **Update Regularly**: Keep Gohan updated for bug fixes
4. **Read Logs**: Check logs when things go wrong
5. **Use Rollback**: Don't hesitate to rollback if something breaks

## Next Steps

- **[Theme Management](theme-management.md)** - Learn about theme features
- **[Installation](installation.md)** - Install Gohan
- **[Development Guide](development.md)** - Report bugs and contribute
