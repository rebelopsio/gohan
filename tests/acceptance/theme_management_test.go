//go:build acceptance
// +build acceptance

package acceptance

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestThemeManagement covers scenarios from theme-management.feature
func TestThemeManagement(t *testing.T) {
	ctx := context.Background()

	t.Run("List available themes", func(t *testing.T) {
		// Given the theme system is initialized
		// And the following themes are available: mocha, latte, frappe, macchiato, gohan
		themeService := setupThemeService(t)

		// When I view available themes
		themes, err := themeService.ListThemes(ctx)

		// Then I should see 5 themes
		require.NoError(t, err)
		assert.Len(t, themes, 5)

		// And each theme should have a name
		for _, theme := range themes {
			assert.NotEmpty(t, theme.Name)
			// And each theme should indicate if it's suitable for day or night use
			assert.Contains(t, []string{"light", "dark"}, theme.Variant)
			// And each theme should show its creator
			assert.NotEmpty(t, theme.Author)
		}
	})

	t.Run("Identify active theme", func(t *testing.T) {
		// Given the theme system is initialized
		// And the "mocha" theme is active
		themeService := setupThemeService(t)
		err := themeService.SetActiveTheme(ctx, "mocha")
		require.NoError(t, err)

		// When I request a list of available themes
		themes, err := themeService.ListThemes(ctx)
		require.NoError(t, err)

		// Then the "mocha" theme should be marked as active
		var mochaFound bool
		var activeCount int
		for _, theme := range themes {
			if theme.Name == "mocha" {
				mochaFound = true
				assert.True(t, theme.IsActive, "mocha should be marked as active")
			}
			if theme.IsActive {
				activeCount++
			}
		}
		assert.True(t, mochaFound, "mocha theme should be in the list")
		// And all other themes should not be marked as active
		assert.Equal(t, 1, activeCount, "only one theme should be active")
	})

	t.Run("View theme information", func(t *testing.T) {
		// Given the theme system is initialized
		themeService := setupThemeService(t)

		// When I view the "latte" theme
		theme, err := themeService.GetTheme(ctx, "latte")
		require.NoError(t, err)

		// Then I should see it is a light theme
		assert.Equal(t, "light", theme.Variant)
		// And I should see it was created by "Catppuccin"
		assert.Equal(t, "Catppuccin", theme.Author)
		// And I should see a preview of its colors
		assert.NotEmpty(t, theme.ColorScheme)
	})

	t.Run("Attempt to get non-existent theme", func(t *testing.T) {
		// Given the theme system is initialized
		themeService := setupThemeService(t)

		// When I request details for the "nonexistent" theme
		_, err := themeService.GetTheme(ctx, "nonexistent")

		// Then I should receive an error
		require.Error(t, err)
		// And the error should indicate the theme was not found
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Find themes suitable for nighttime use", func(t *testing.T) {
		// Given the theme system is initialized
		themeService := setupThemeService(t)

		// When I look for dark themes
		darkThemes, err := themeService.ListThemesByVariant(ctx, "dark")
		require.NoError(t, err)

		// Then I should see 4 themes
		assert.Len(t, darkThemes, 4)
		// And all themes should be suitable for low-light environments
		for _, theme := range darkThemes {
			assert.Equal(t, "dark", theme.Variant)
		}
	})

	t.Run("Find themes suitable for daytime use", func(t *testing.T) {
		// Given the theme system is initialized
		themeService := setupThemeService(t)

		// When I look for light themes
		lightThemes, err := themeService.ListThemesByVariant(ctx, "light")
		require.NoError(t, err)

		// Then I should see 1 theme
		require.Len(t, lightThemes, 1)
		// And the theme should be "latte"
		assert.Equal(t, "latte", lightThemes[0].Name)
	})

	t.Run("Get active theme information", func(t *testing.T) {
		// Given the theme system is initialized
		// And the "frappe" theme is active
		themeService := setupThemeService(t)
		err := themeService.SetActiveTheme(ctx, "frappe")
		require.NoError(t, err)

		// When I request the active theme
		activeTheme, err := themeService.GetActiveTheme(ctx)
		require.NoError(t, err)

		// Then I should receive the "frappe" theme
		assert.Equal(t, "frappe", activeTheme.Name)
		// And it should be marked as active
		assert.True(t, activeTheme.IsActive)
	})

	t.Run("System has default theme when none set", func(t *testing.T) {
		// Given no theme has been explicitly set
		themeService := setupThemeService(t)

		// When I request the active theme
		activeTheme, err := themeService.GetActiveTheme(ctx)
		require.NoError(t, err)

		// Then I should receive the "mocha" theme
		assert.Equal(t, "mocha", activeTheme.Name)
		// And it should be marked as active
		assert.True(t, activeTheme.IsActive)
	})
}

// setupThemeService initializes a theme service for testing
// This will be implemented as we build the domain
func setupThemeService(t *testing.T) ThemeService {
	t.Helper()
	// TODO: Implement once we have the domain models
	t.Skip("Theme service not yet implemented")
	return nil
}

// ThemeService defines the interface for theme operations
// This interface will guide our domain implementation
type ThemeService interface {
	ListThemes(ctx context.Context) ([]ThemeInfo, error)
	ListThemesByVariant(ctx context.Context, variant string) ([]ThemeInfo, error)
	GetTheme(ctx context.Context, name string) (ThemeInfo, error)
	GetActiveTheme(ctx context.Context) (ThemeInfo, error)
	SetActiveTheme(ctx context.Context, name string) error
}

// ThemeInfo represents theme information returned to users
type ThemeInfo struct {
	Name        string
	DisplayName string
	Author      string
	Description string
	Variant     string // "dark" or "light"
	IsActive    bool
	ColorScheme map[string]string
}
