# Theme Management

Gohan provides a powerful theming system that allows you to customize the appearance of your entire Hyprland environment with a single command.

## Table of Contents
- [Quick Start](#quick-start)
- [Available Themes](#available-themes)
- [Applying Themes](#applying-themes)
- [Theme Preview](#theme-preview)
- [Theme History & Rollback](#theme-history-rollback)
- [Customization](#customization)
- [Component Support](#component-support)

## Quick Start

```bash
# List all available themes
gohan theme list

# Preview a theme
gohan theme preview latte

# Apply a theme
gohan theme set latte

# Show current theme
gohan theme show

# Rollback to previous theme
gohan theme rollback
```

## Available Themes

Gohan includes 5 beautiful built-in themes based on the Catppuccin color palette:

### 🌙 Mocha (Dark)
**Default theme** - Soothing pastel colors on a dark background
```bash
gohan theme set mocha
```
- Perfect for night-time coding
- Easy on the eyes
- Professional appearance

### ☕ Latte (Light)
Pastel colors on a light background
```bash
gohan theme set latte
```
- Ideal for daytime use
- High contrast
- Clean and minimal

### 🥤 Frappe (Dark)
Muted pastels on a dark surface
```bash
gohan theme set frappe
```
- Balanced contrast
- Soft colors
- Great for extended use

### 🥛 Macchiato (Dark)
Rich pastels on a dark base
```bash
gohan theme set macchiato
```
- Vibrant accents
- Deep background
- Modern aesthetic

### 🍚 Gohan (Dark)
Custom theme with unique color scheme
```bash
gohan theme set gohan
```
- Project signature theme
- Distinctive palette
- Optimized for Hyprland

## Applying Themes

### Basic Application

The simplest way to apply a theme:

```bash
gohan theme set <theme-name>
```

**Example:**
```bash
$ gohan theme set latte
Applying theme 'latte'...

✓ Successfully applied theme 'Catppuccin Latte'

Theme applied successfully!
Updated configuration files:
  - Hyprland configuration
  - Waybar configuration
  - Kitty terminal colors
  - Rofi/Fuzzel theme
  - Mako notifications
  - Alacritty terminal
  - Fuzzel launcher

Backups have been created. Use 'gohan theme rollback' to restore previous theme.
```

### What Happens When You Apply a Theme?

1. **Backup Creation**: Current configurations are backed up
2. **Template Processing**: Theme templates are processed with new colors
3. **File Deployment**: New configurations are deployed atomically
4. **Component Reload**: Hyprland and Waybar automatically reload
5. **State Saving**: Theme state is persisted to disk
6. **History Recording**: Theme change is added to history

### Auto-Reload Feature

When you apply a theme, Gohan automatically:
- **Reloads Hyprland**: Uses `hyprctl reload` to apply new config
- **Restarts Waybar**: Kills and restarts Waybar with new theme
- **No Manual Intervention**: Everything happens seamlessly

## Theme Preview

Preview a theme before applying it:

```bash
gohan theme preview <theme-name>
```

**Example:**
```bash
$ gohan theme preview mocha

Theme: Catppuccin Mocha
Variant: dark
Author: Catppuccin

Color Scheme:
  Base:      #1e1e2e
  Surface:   #313244
  Overlay:   #45475a
  Text:      #cdd6f4
  Subtext:   #bac2de

  Rosewater: #f5e0dc
  Flamingo:  #f2cdcd
  Pink:      #f5c2e7
  Mauve:     #cba6f7
  Red:       #f38ba8
  Maroon:    #eba0ac
  Peach:     #fab387
  Yellow:    #f9e2af
  Green:     #a6e3a1
  Teal:      #94e2d5
  Sky:       #89dceb
  Sapphire:  #74c7ec
  Blue:      #89b4fa
  Lavender:  #b4befe
```

## Theme History & Rollback

### View Theme History

See your theme change history:

```bash
gohan theme list
```

The active theme is marked with ●:

```
Available Themes:
  ● mocha      Catppuccin Mocha      (dark)   [Catppuccin]
    latte      Catppuccin Latte      (light)  [Catppuccin]
    frappe     Catppuccin Frappe     (dark)   [Catppuccin]
    macchiato  Catppuccin Macchiato  (dark)   [Catppuccin]
    gohan      Gohan                 (dark)   [Gohan Team]
```

### Rollback to Previous Theme

Undo your last theme change:

```bash
gohan theme rollback
```

**Example:**
```bash
$ gohan theme rollback
Rolling back to previous theme...

✓ Successfully rolled back from 'latte' to 'mocha'

Theme rolled back successfully!
Restored theme: mocha
Previous theme: latte

Configuration files have been updated.
Use 'gohan theme rollback' again to continue rolling back through history.
```

### Sequential Rollbacks

You can rollback multiple times to step through your theme history:

```bash
# Apply several themes
gohan theme set mocha
gohan theme set latte
gohan theme set frappe

# Roll back one by one
gohan theme rollback  # frappe -> latte
gohan theme rollback  # latte -> mocha
gohan theme rollback  # mocha -> (previous theme)
```

### History Limits

- Gohan maintains up to **10 theme changes** in history
- Oldest entries are automatically removed
- History persists across restarts
- History is stored in `~/.config/gohan/theme-history.json`

### Error Handling

If no history is available:

```bash
$ gohan theme rollback

✗ No theme history available
You need to change themes at least once before you can rollback.
```

## Show Current Theme

View detailed information about the active theme:

```bash
gohan theme show
```

**Example:**
```bash
$ gohan theme show

Active Theme: Catppuccin Mocha

  Name:    mocha
  Author:  Catppuccin
  Variant: dark
  Description: Soothing pastel theme for the high-spirited!
```

**With verbose output:**
```bash
$ gohan theme show --verbose

Active Theme: Catppuccin Mocha

  Name:    mocha
  Author:  Catppuccin
  Variant: dark
  Description: Soothing pastel theme for the high-spirited!

Colors:
  base:      #1e1e2e
  surface:   #313244
  overlay:   #45475a
  text:      #cdd6f4
  subtext:   #bac2de
  rosewater: #f5e0dc
  flamingo:  #f2cdcd
  pink:      #f5c2e7
  mauve:     #cba6f7
  red:       #f38ba8
  maroon:    #eba0ac
  peach:     #fab387
  yellow:    #f9e2af
  green:     #a6e3a1
  teal:      #94e2d5
  sky:       #89dceb
  sapphire:  #74c7ec
  blue:      #89b4fa
  lavender:  #b4befe
```

## Component Support

Gohan applies themes to all supported components automatically:

### Hyprland
- Window borders and shadows
- Workspace colors
- General UI elements

### Waybar
- Bar background and text
- Module colors
- Hover states

### Kitty Terminal
- Background and foreground
- Cursor colors
- 16-color palette
- Bold/italic variants

### Rofi Launcher
- Window background
- Text and selection colors
- Border styling

### Mako Notifications
- Background and text
- Urgency levels (low/normal/high)
- Border colors
- Progress bar

### Alacritty Terminal
- Complete color scheme
- Normal and bright colors
- Cursor styling

### Fuzzel Launcher
- Background and text
- Match highlighting
- Selection colors
- Border styling

## Customization

### Template Locations

Theme templates are located in:
```
templates/
├── hyprland/hyprland.conf.tmpl
├── waybar/style.css.tmpl
├── kitty/kitty.conf.tmpl
├── rofi/config.rasi.tmpl
├── mako/config.tmpl
├── alacritty/alacritty.toml.tmpl
└── fuzzel/fuzzel.ini.tmpl
```

### Available Template Variables

All templates have access to theme color variables:

```
{{.theme_name}}          - Theme identifier
{{.theme_display_name}}  - Human-readable name
{{.theme_variant}}       - "dark" or "light"

{{.theme_base}}          - Background color
{{.theme_surface}}       - Surface color
{{.theme_overlay}}       - Overlay color
{{.theme_text}}          - Primary text
{{.theme_subtext}}       - Secondary text

{{.theme_rosewater}}     - Rosewater accent
{{.theme_flamingo}}      - Flamingo accent
{{.theme_pink}}          - Pink accent
{{.theme_mauve}}         - Mauve accent
{{.theme_red}}           - Red accent
{{.theme_maroon}}        - Maroon accent
{{.theme_peach}}         - Peach accent
{{.theme_yellow}}        - Yellow accent
{{.theme_green}}         - Green accent
{{.theme_teal}}          - Teal accent
{{.theme_sky}}           - Sky accent
{{.theme_sapphire}}      - Sapphire accent
{{.theme_blue}}          - Blue accent
{{.theme_lavender}}      - Lavender accent
```

### Creating Custom Templates

You can customize templates to match your preferences:

1. **Copy the template**:
   ```bash
   cp templates/waybar/style.css.tmpl templates/waybar/style.css.tmpl.custom
   ```

2. **Edit with theme variables**:
   ```css
   /* Custom Waybar styling */
   window#waybar {
       background-color: {{.theme_base}};
       color: {{.theme_text}};
       border-top: 2px solid {{.theme_mauve}};
   }

   #workspaces button.active {
       background-color: {{.theme_blue}};
       color: {{.theme_base}};
   }
   ```

3. **Apply your custom template**:
   ```bash
   gohan config deploy --template templates/waybar/style.css.tmpl.custom \
       --target ~/.config/waybar/style.css
   ```

## Persistence

Theme state is automatically persisted to disk:

```bash
# State location
~/.config/gohan/theme-state.json

# Example content:
{
  "theme_name": "mocha",
  "theme_variant": "dark",
  "set_at": "2025-10-28T15:30:00Z"
}
```

Theme settings persist across:
- Terminal restarts
- System reboots
- Gohan updates

## Troubleshooting

### Theme doesn't apply to all components

**Problem**: Some components still show old theme

**Solution**:
```bash
# Manually reload components
hyprctl reload
killall waybar && waybar &
```

### Colors look wrong

**Problem**: Terminal doesn't support true color

**Solution**:
```bash
# Check terminal color support
echo $COLORTERM  # Should output "truecolor" or "24bit"

# For Kitty/Alacritty, this should work automatically
# For other terminals, enable true color support
```

### Rollback doesn't work

**Problem**: "No theme history available"

**Solution**:
- History only tracks theme changes made with `gohan theme set`
- Apply at least one theme to build history
- Check `~/.config/gohan/theme-history.json` exists

## Best Practices

1. **Preview Before Applying**: Use `gohan theme preview` to see colors first
2. **Backup Important Configs**: While Gohan creates backups, maintain your own
3. **Use Rollback**: Don't hesitate to rollback if you don't like a theme
4. **Customize Templates**: Make themes your own with custom templates
5. **Check Component Support**: Ensure your components are supported

## Next Steps

- **[Installation](installation.md)** - Install Gohan
- **[Troubleshooting](troubleshooting.md)** - Solve common issues
- **[Development Guide](development.md)** - Contributing and architecture
