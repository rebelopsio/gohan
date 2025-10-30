package theme

import (
	"context"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThemeToTemplateVars(t *testing.T) {
	// Create a test theme
	themeObj := createTestTheme(t)

	// Convert to template vars
	vars := ThemeToTemplateVars(themeObj)

	t.Run("includes all base colors", func(t *testing.T) {
		assert.NotEmpty(t, vars["theme_base"])
		assert.NotEmpty(t, vars["theme_surface"])
		assert.NotEmpty(t, vars["theme_overlay"])
		assert.NotEmpty(t, vars["theme_text"])
		assert.NotEmpty(t, vars["theme_subtext"])
	})

	t.Run("includes all accent colors", func(t *testing.T) {
		assert.NotEmpty(t, vars["theme_rosewater"])
		assert.NotEmpty(t, vars["theme_flamingo"])
		assert.NotEmpty(t, vars["theme_pink"])
		assert.NotEmpty(t, vars["theme_mauve"])
		assert.NotEmpty(t, vars["theme_red"])
		assert.NotEmpty(t, vars["theme_maroon"])
		assert.NotEmpty(t, vars["theme_peach"])
		assert.NotEmpty(t, vars["theme_yellow"])
		assert.NotEmpty(t, vars["theme_green"])
		assert.NotEmpty(t, vars["theme_teal"])
		assert.NotEmpty(t, vars["theme_sky"])
		assert.NotEmpty(t, vars["theme_sapphire"])
		assert.NotEmpty(t, vars["theme_blue"])
		assert.NotEmpty(t, vars["theme_lavender"])
	})

	t.Run("includes theme metadata", func(t *testing.T) {
		assert.Equal(t, "mocha", vars["theme_name"])
		assert.Equal(t, "Catppuccin Mocha", vars["theme_display_name"])
		assert.Equal(t, "dark", vars["theme_variant"])
	})

	t.Run("color values are hex codes", func(t *testing.T) {
		assert.Regexp(t, "^#[0-9A-Fa-f]{6}$", vars["theme_base"])
		assert.Regexp(t, "^#[0-9A-Fa-f]{6}$", vars["theme_text"])
		assert.Regexp(t, "^#[0-9A-Fa-f]{6}$", vars["theme_mauve"])
	})
}

func TestGetComponentConfigurations(t *testing.T) {
	configs := GetComponentConfigurations()

	t.Run("includes all expected components", func(t *testing.T) {
		componentNames := make([]string, len(configs))
		for i, cfg := range configs {
			componentNames[i] = cfg.Component
		}

		assert.Contains(t, componentNames, "hyprland")
		assert.Contains(t, componentNames, "waybar")
		assert.Contains(t, componentNames, "kitty")
		assert.Contains(t, componentNames, "rofi")
	})

	t.Run("each component has valid configuration", func(t *testing.T) {
		for _, cfg := range configs {
			assert.NotEmpty(t, cfg.Component, "component name should not be empty")
			assert.NotEmpty(t, cfg.TemplatePath, "template path should not be empty")
			assert.NotEmpty(t, cfg.TargetPath, "target path should not be empty")
			assert.True(t, cfg.BackupBefore, "should backup before overwriting")
		}
	})
}

// Helper function to create a test theme
func createTestTheme(t *testing.T) *theme.Theme {
	t.Helper()

	// Use a standard theme from the domain
	registry := theme.NewThemeRegistry()
	err := theme.InitializeStandardThemes(registry)
	require.NoError(t, err)

	themeObj, err := registry.FindByName(context.Background(), theme.ThemeMocha)
	require.NoError(t, err)

	return themeObj
}
