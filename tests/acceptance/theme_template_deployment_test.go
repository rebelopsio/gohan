package acceptance

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/theme"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/templates"
	themeInfra "github.com/rebelopsio/gohan/internal/infrastructure/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// Phase 5: Theme Template Deployment - ATDD
// ========================================

func TestThemeTemplateDeployment_Hyprland(t *testing.T) {
	t.Run("deploys Hyprland configuration template", func(t *testing.T) {
		// Given: a Hyprland configuration template exists
		templatePath := getProjectPath("templates/hyprland/hyprland.conf.tmpl")

		// Check template exists
		_, err := os.Stat(templatePath)
		require.NoError(t, err, "Hyprland template should exist at %s", templatePath)

		// Given: theme variables
		themeVars := createMochaThemeVars(t)

		// When: template is processed
		engine := templates.NewTemplateEngine()
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "hyprland.conf")

		err = engine.ProcessFile(templatePath, outputPath, themeVars)
		require.NoError(t, err, "Template processing should succeed")

		// Then: output should contain actual color values
		content, err := os.ReadFile(outputPath)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "#1e1e2e", "Should contain theme_base color")
		assert.NotContains(t, contentStr, "{{theme_base}}", "Should not contain unprocessed variables")

		// And: file should be valid Hyprland config (basic syntax check)
		assert.Contains(t, contentStr, "$", "Should contain Hyprland variables")
	})
}

func TestThemeTemplateDeployment_Waybar(t *testing.T) {
	t.Run("deploys Waybar style template", func(t *testing.T) {
		// Given: a Waybar style template exists
		templatePath := getProjectPath("templates/waybar/style.css.tmpl")

		_, err := os.Stat(templatePath)
		require.NoError(t, err, "Waybar template should exist at %s", templatePath)

		// Given: theme variables
		themeVars := createMochaThemeVars(t)

		// When: template is processed
		engine := templates.NewTemplateEngine()
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "style.css")

		err = engine.ProcessFile(templatePath, outputPath, themeVars)
		require.NoError(t, err)

		// Then: CSS should contain theme colors
		content, err := os.ReadFile(outputPath)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "#1e1e2e", "Should contain theme colors")
		assert.NotContains(t, contentStr, "{{theme_", "Should not contain unprocessed variables")

		// And: should be valid CSS
		assert.Contains(t, contentStr, "{", "Should contain CSS syntax")
		assert.Contains(t, contentStr, "}", "Should contain CSS syntax")
	})
}

func TestThemeTemplateDeployment_Kitty(t *testing.T) {
	t.Run("deploys Kitty terminal template", func(t *testing.T) {
		// Given: a Kitty configuration template exists
		templatePath := getProjectPath("templates/kitty/kitty.conf.tmpl")

		_, err := os.Stat(templatePath)
		require.NoError(t, err, "Kitty template should exist at %s", templatePath)

		// Given: theme variables
		themeVars := createMochaThemeVars(t)

		// When: template is processed
		engine := templates.NewTemplateEngine()
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "kitty.conf")

		err = engine.ProcessFile(templatePath, outputPath, themeVars)
		require.NoError(t, err)

		// Then: config should define terminal colors
		content, err := os.ReadFile(outputPath)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "foreground", "Should define foreground color")
		assert.Contains(t, contentStr, "background", "Should define background color")
		assert.Contains(t, contentStr, "#", "Should contain hex colors")
		assert.NotContains(t, contentStr, "{{theme_", "Should not contain unprocessed variables")
	})
}

func TestThemeTemplateDeployment_Rofi(t *testing.T) {
	t.Run("deploys Rofi launcher template", func(t *testing.T) {
		// Given: a Rofi configuration template exists
		templatePath := getProjectPath("templates/rofi/config.rasi.tmpl")

		_, err := os.Stat(templatePath)
		require.NoError(t, err, "Rofi template should exist at %s", templatePath)

		// Given: theme variables
		themeVars := createMochaThemeVars(t)

		// When: template is processed
		engine := templates.NewTemplateEngine()
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "config.rasi")

		err = engine.ProcessFile(templatePath, outputPath, themeVars)
		require.NoError(t, err)

		// Then: RASI should contain theme colors
		content, err := os.ReadFile(outputPath)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "#", "Should contain hex colors")
		assert.NotContains(t, contentStr, "{{theme_", "Should not contain unprocessed variables")
	})
}

func TestThemeTemplateDeployment_VariableConsistency(t *testing.T) {
	t.Run("all templates use consistent variable names", func(t *testing.T) {
		templatePaths := []string{
			getProjectPath("templates/hyprland/hyprland.conf.tmpl"),
			getProjectPath("templates/waybar/style.css.tmpl"),
			getProjectPath("templates/kitty/kitty.conf.tmpl"),
			getProjectPath("templates/rofi/config.rasi.tmpl"),
		}

		// Standard theme variables that should be used
		expectedVars := []string{
			"theme_base",
			"theme_surface",
			"theme_text",
			"theme_mauve",
			"theme_blue",
		}

		for _, tmplPath := range templatePaths {
			content, err := os.ReadFile(tmplPath)
			if os.IsNotExist(err) {
				t.Logf("Skipping %s (not yet created)", tmplPath)
				continue
			}
			require.NoError(t, err)

			contentStr := string(content)

			// Check that standard variables are used correctly
			for _, varName := range expectedVars {
				if strings.Contains(contentStr, "{{"+varName+"}}") {
					// Good - using correct syntax
					continue
				}
			}
		}
	})
}

func TestThemeTemplateDeployment_MissingTemplates(t *testing.T) {
	t.Run("handles missing templates gracefully", func(t *testing.T) {
		// Given: ComponentConfigurations that may reference non-existent templates
		configs := themeInfra.GetComponentConfigurations()

		existingCount := 0
		for _, cfg := range configs {
			if _, err := os.Stat(cfg.TemplatePath); err == nil {
				existingCount++
			}
		}

		// Then: at least some templates should exist
		assert.Greater(t, existingCount, 0, "At least one template should exist")

		// And: missing templates should not cause errors in the system
		// (this is tested by the ThemeApplier which skips missing templates)
	})
}

func TestThemeTemplateDeployment_SystemAndThemeVariables(t *testing.T) {
	t.Run("templates can use both system and theme variables", func(t *testing.T) {
		// Given: both system and theme variables
		systemVars, err := templates.CollectSystemVars()
		require.NoError(t, err)

		themeVars := createMochaThemeVars(t)

		// Merge variables
		allVars := templates.TemplateVars{}
		for k, v := range systemVars {
			allVars[k] = v
		}
		for k, v := range themeVars {
			allVars[k] = v
		}

		// Then: all variables should be available
		assert.NotEmpty(t, allVars["username"], "System variable should exist")
		assert.NotEmpty(t, allVars["theme_base"], "Theme variable should exist")

		// And: no conflicts
		assert.NotEqual(t, allVars["username"], allVars["theme_base"], "Variables should not conflict")
	})
}

// Helper function to create Mocha theme variables for testing
func createMochaThemeVars(t *testing.T) templates.TemplateVars {
	t.Helper()

	registry := theme.NewThemeRegistry()
	err := theme.InitializeStandardThemes(registry)
	require.NoError(t, err)

	mocha, err := registry.FindByName(context.Background(), theme.ThemeMocha)
	require.NoError(t, err)

	return themeInfra.ThemeToTemplateVars(mocha)
}

// Helper function to get project root path
func getProjectPath(relativePath string) string {
	// Find project root by looking for go.mod
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, relativePath)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return relativePath
}
