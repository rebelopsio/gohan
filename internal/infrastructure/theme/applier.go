package theme

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebelopsio/gohan/internal/domain/theme"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/configservice"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/templates"
)

// ComponentConfiguration defines how a theme applies to a component
type ComponentConfiguration struct {
	Component    string // Component name (hyprland, waybar, kitty, rofi)
	TemplatePath string // Path to template file
	TargetPath   string // Where to deploy the config
	BackupBefore bool   // Whether to backup before overwriting
}

// ThemeApplierImpl implements the ThemeApplier interface using ConfigDeployer
type ThemeApplierImpl struct {
	configDeployer    *configservice.ConfigDeployer
	componentReloader *ComponentReloader
}

// NewThemeApplier creates a new theme applier
func NewThemeApplier(configDeployer *configservice.ConfigDeployer) *ThemeApplierImpl {
	executor := NewSystemCommandExecutor()
	return &ThemeApplierImpl{
		configDeployer:    configDeployer,
		componentReloader: NewComponentReloader(executor),
	}
}

// NewThemeApplierWithReloader creates a new theme applier with a custom reloader
func NewThemeApplierWithReloader(configDeployer *configservice.ConfigDeployer, reloader *ComponentReloader) *ThemeApplierImpl {
	return &ThemeApplierImpl{
		configDeployer:    configDeployer,
		componentReloader: reloader,
	}
}

// ApplyTheme applies a theme to the system by updating configuration files
func (ta *ThemeApplierImpl) ApplyTheme(ctx context.Context, th *theme.Theme) error {
	// Convert theme to template variables
	vars := ThemeToTemplateVars(th)

	// Get system variables (home, config_dir, etc.)
	systemVars, err := templates.CollectSystemVars()
	if err != nil {
		return fmt.Errorf("failed to collect system variables: %w", err)
	}

	// Merge theme vars with system vars
	for k, v := range systemVars {
		vars[k] = v
	}

	// Get component configurations
	componentConfigs := GetComponentConfigurations()

	// Convert to ConfigurationFile format
	configFiles := make([]configservice.ConfigurationFile, 0, len(componentConfigs))
	for _, compCfg := range componentConfigs {
		// Check if template file exists
		if _, err := os.Stat(compCfg.TemplatePath); err != nil {
			// Skip if template doesn't exist (optional component)
			continue
		}

		configFiles = append(configFiles, configservice.ConfigurationFile{
			SourceTemplate: compCfg.TemplatePath,
			TargetPath:     compCfg.TargetPath,
			Permissions:    0644,
			BackupBefore:   compCfg.BackupBefore,
		})
	}

	// Deploy configurations without progress channel (synchronous)
	if len(configFiles) > 0 {
		progressChan := make(chan configservice.DeploymentProgress, len(configFiles)*3)
		done := make(chan error, 1)

		go func() {
			done <- ta.configDeployer.DeployConfigurations(ctx, configFiles, vars, progressChan)
			close(progressChan)
		}()

		// Consume progress updates
		for range progressChan {
			// Silently consume for now
		}

		if err := <-done; err != nil {
			return fmt.Errorf("failed to deploy configurations: %w", err)
		}
	}

	// Reload components to apply changes
	if ta.componentReloader != nil {
		if err := ta.componentReloader.ReloadAll(ctx); err != nil {
			// Log warning but don't fail the operation
			fmt.Printf("Warning: some components failed to reload: %v\n", err)
		}
	}

	return nil
}

// ThemeToTemplateVars converts a theme to template variables
func ThemeToTemplateVars(th *theme.Theme) templates.TemplateVars {
	cs := th.ColorScheme()

	vars := templates.TemplateVars{
		// Theme metadata
		"theme_name":         string(th.Name()),
		"theme_display_name": th.DisplayName(),
		"theme_variant":      string(th.Variant()),

		// Base colors
		"theme_base":    cs.Base().String(),
		"theme_surface": cs.Surface().String(),
		"theme_overlay": cs.Overlay().String(),
		"theme_text":    cs.Text().String(),
		"theme_subtext": cs.Subtext().String(),

		// Accent colors
		"theme_rosewater": cs.Rosewater().String(),
		"theme_flamingo":  cs.Flamingo().String(),
		"theme_pink":      cs.Pink().String(),
		"theme_mauve":     cs.Mauve().String(),
		"theme_red":       cs.Red().String(),
		"theme_maroon":    cs.Maroon().String(),
		"theme_peach":     cs.Peach().String(),
		"theme_yellow":    cs.Yellow().String(),
		"theme_green":     cs.Green().String(),
		"theme_teal":      cs.Teal().String(),
		"theme_sky":       cs.Sky().String(),
		"theme_sapphire":  cs.Sapphire().String(),
		"theme_blue":      cs.Blue().String(),
		"theme_lavender":  cs.Lavender().String(),
	}

	return vars
}

// getProjectRoot finds the project root by looking for go.mod
func getProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// Fallback to current directory
	return "."
}

// GetComponentConfigurations returns the list of component configurations
func GetComponentConfigurations() []ComponentConfiguration {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config")
	projectRoot := getProjectRoot()

	return []ComponentConfiguration{
		{
			Component:    "hyprland",
			TemplatePath: filepath.Join(projectRoot, "templates/hyprland/hyprland.conf.tmpl"),
			TargetPath:   filepath.Join(configDir, "hypr", "hyprland.conf"),
			BackupBefore: true,
		},
		{
			Component:    "waybar",
			TemplatePath: filepath.Join(projectRoot, "templates/waybar/style.css.tmpl"),
			TargetPath:   filepath.Join(configDir, "waybar", "style.css"),
			BackupBefore: true,
		},
		{
			Component:    "kitty",
			TemplatePath: filepath.Join(projectRoot, "templates/kitty/kitty.conf.tmpl"),
			TargetPath:   filepath.Join(configDir, "kitty", "kitty.conf"),
			BackupBefore: true,
		},
		{
			Component:    "rofi",
			TemplatePath: filepath.Join(projectRoot, "templates/rofi/config.rasi.tmpl"),
			TargetPath:   filepath.Join(configDir, "rofi", "config.rasi"),
			BackupBefore: true,
		},
		{
			Component:    "mako",
			TemplatePath: filepath.Join(projectRoot, "templates/mako/config.tmpl"),
			TargetPath:   filepath.Join(configDir, "mako", "config"),
			BackupBefore: true,
		},
		{
			Component:    "alacritty",
			TemplatePath: filepath.Join(projectRoot, "templates/alacritty/alacritty.toml.tmpl"),
			TargetPath:   filepath.Join(configDir, "alacritty", "alacritty.toml"),
			BackupBefore: true,
		},
		{
			Component:    "fuzzel",
			TemplatePath: filepath.Join(projectRoot, "templates/fuzzel/fuzzel.ini.tmpl"),
			TargetPath:   filepath.Join(configDir, "fuzzel", "fuzzel.ini"),
			BackupBefore: true,
		},
	}
}
