package theme

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewThemeRegistry(t *testing.T) {
	t.Run("creates empty registry", func(t *testing.T) {
		registry := NewThemeRegistry()
		assert.NotNil(t, registry)

		themes := registry.ListAll()
		assert.Empty(t, themes)
	})
}

func TestThemeRegistry_Register(t *testing.T) {
	tests := []struct {
		name    string
		theme   *Theme
		wantErr bool
	}{
		{
			name: "register valid theme",
			theme: createTestTheme(t, ThemeMocha, ThemeMetadata{
				displayName: "Catppuccin Mocha",
				author:      "Catppuccin",
				variant:     ThemeVariantDark,
			}),
			wantErr: false,
		},
		{
			name:    "register nil theme",
			theme:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewThemeRegistry()
			err := registry.Register(tt.theme)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify theme was registered
				themes := registry.ListAll()
				assert.Len(t, themes, 1)
			}
		})
	}
}

func TestThemeRegistry_RegisterDuplicate(t *testing.T) {
	registry := NewThemeRegistry()

	theme1 := createTestTheme(t, ThemeMocha, ThemeMetadata{
		displayName: "Catppuccin Mocha",
		author:      "Catppuccin",
		variant:     ThemeVariantDark,
	})

	theme2 := createTestTheme(t, ThemeMocha, ThemeMetadata{
		displayName: "Catppuccin Mocha Updated",
		author:      "Catppuccin",
		variant:     ThemeVariantDark,
	})

	// Register first theme
	err := registry.Register(theme1)
	require.NoError(t, err)

	// Attempt to register duplicate
	err = registry.Register(theme2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")

	// Verify only one theme exists
	themes := registry.ListAll()
	assert.Len(t, themes, 1)
}

func TestThemeRegistry_FindByName(t *testing.T) {
	registry := NewThemeRegistry()

	mocha := createTestTheme(t, ThemeMocha, ThemeMetadata{
		displayName: "Catppuccin Mocha",
		author:      "Catppuccin",
		variant:     ThemeVariantDark,
	})

	latte := createTestTheme(t, ThemeLatte, ThemeMetadata{
		displayName: "Catppuccin Latte",
		author:      "Catppuccin",
		variant:     ThemeVariantLight,
	})

	err := registry.Register(mocha)
	require.NoError(t, err)
	err = registry.Register(latte)
	require.NoError(t, err)

	t.Run("find existing theme", func(t *testing.T) {
		theme, err := registry.FindByName(context.Background(), ThemeMocha)
		require.NoError(t, err)
		assert.Equal(t, ThemeMocha, theme.Name())
		assert.Equal(t, "Catppuccin Mocha", theme.DisplayName())
	})

	t.Run("find non-existent theme", func(t *testing.T) {
		_, err := registry.FindByName(context.Background(), ThemeName("nonexistent"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestThemeRegistry_ListAll(t *testing.T) {
	registry := NewThemeRegistry()

	mocha := createTestTheme(t, ThemeMocha, ThemeMetadata{
		displayName: "Catppuccin Mocha",
		author:      "Catppuccin",
		variant:     ThemeVariantDark,
	})

	latte := createTestTheme(t, ThemeLatte, ThemeMetadata{
		displayName: "Catppuccin Latte",
		author:      "Catppuccin",
		variant:     ThemeVariantLight,
	})

	frappe := createTestTheme(t, ThemeFrappe, ThemeMetadata{
		displayName: "Catppuccin Frappe",
		author:      "Catppuccin",
		variant:     ThemeVariantDark,
	})

	err := registry.Register(mocha)
	require.NoError(t, err)
	err = registry.Register(latte)
	require.NoError(t, err)
	err = registry.Register(frappe)
	require.NoError(t, err)

	t.Run("list all registered themes", func(t *testing.T) {
		themes := registry.ListAll()
		assert.Len(t, themes, 3)

		names := make([]ThemeName, len(themes))
		for i, theme := range themes {
			names[i] = theme.Name()
		}
		assert.Contains(t, names, ThemeMocha)
		assert.Contains(t, names, ThemeLatte)
		assert.Contains(t, names, ThemeFrappe)
	})
}

func TestThemeRegistry_SetActive(t *testing.T) {
	registry := NewThemeRegistry()

	mocha := createTestTheme(t, ThemeMocha, ThemeMetadata{
		displayName: "Catppuccin Mocha",
		author:      "Catppuccin",
		variant:     ThemeVariantDark,
	})

	latte := createTestTheme(t, ThemeLatte, ThemeMetadata{
		displayName: "Catppuccin Latte",
		author:      "Catppuccin",
		variant:     ThemeVariantLight,
	})

	err := registry.Register(mocha)
	require.NoError(t, err)
	err = registry.Register(latte)
	require.NoError(t, err)

	t.Run("set active theme", func(t *testing.T) {
		ctx := context.Background()

		err := registry.SetActive(ctx, ThemeMocha)
		assert.NoError(t, err)

		active, err := registry.GetActive(ctx)
		require.NoError(t, err)
		assert.Equal(t, ThemeMocha, active.Name())
	})

	t.Run("set non-existent theme as active", func(t *testing.T) {
		err := registry.SetActive(context.Background(), ThemeName("nonexistent"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("change active theme", func(t *testing.T) {
		ctx := context.Background()

		// Set mocha as active
		err := registry.SetActive(ctx, ThemeMocha)
		require.NoError(t, err)

		// Change to latte
		err = registry.SetActive(ctx, ThemeLatte)
		require.NoError(t, err)

		// Verify latte is now active
		active, err := registry.GetActive(ctx)
		require.NoError(t, err)
		assert.Equal(t, ThemeLatte, active.Name())
	})
}

func TestThemeRegistry_GetActive(t *testing.T) {
	t.Run("get active when none set - returns mocha as default", func(t *testing.T) {
		registry := NewThemeRegistry()

		mocha := createTestTheme(t, ThemeMocha, ThemeMetadata{
			displayName: "Catppuccin Mocha",
			author:      "Catppuccin",
			variant:     ThemeVariantDark,
		})

		err := registry.Register(mocha)
		require.NoError(t, err)

		active, err := registry.GetActive(context.Background())
		require.NoError(t, err)
		assert.Equal(t, ThemeMocha, active.Name())
	})

	t.Run("get active when set", func(t *testing.T) {
		registry := NewThemeRegistry()

		mocha := createTestTheme(t, ThemeMocha, ThemeMetadata{
			displayName: "Catppuccin Mocha",
			author:      "Catppuccin",
			variant:     ThemeVariantDark,
		})

		latte := createTestTheme(t, ThemeLatte, ThemeMetadata{
			displayName: "Catppuccin Latte",
			author:      "Catppuccin",
			variant:     ThemeVariantLight,
		})

		err := registry.Register(mocha)
		require.NoError(t, err)
		err = registry.Register(latte)
		require.NoError(t, err)

		ctx := context.Background()
		err = registry.SetActive(ctx, ThemeLatte)
		require.NoError(t, err)

		active, err := registry.GetActive(ctx)
		require.NoError(t, err)
		assert.Equal(t, ThemeLatte, active.Name())
	})

	t.Run("get active when no themes registered", func(t *testing.T) {
		registry := NewThemeRegistry()

		_, err := registry.GetActive(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no themes")
	})
}

func TestThemeRegistry_ListByVariant(t *testing.T) {
	registry := NewThemeRegistry()

	// Register dark themes
	mocha := createTestTheme(t, ThemeMocha, ThemeMetadata{
		displayName: "Catppuccin Mocha",
		author:      "Catppuccin",
		variant:     ThemeVariantDark,
	})

	frappe := createTestTheme(t, ThemeFrappe, ThemeMetadata{
		displayName: "Catppuccin Frappe",
		author:      "Catppuccin",
		variant:     ThemeVariantDark,
	})

	macchiato := createTestTheme(t, ThemeMacchiato, ThemeMetadata{
		displayName: "Catppuccin Macchiato",
		author:      "Catppuccin",
		variant:     ThemeVariantDark,
	})

	gohan := createTestTheme(t, ThemeGohan, ThemeMetadata{
		displayName: "Gohan",
		author:      "Gohan Team",
		variant:     ThemeVariantDark,
	})

	// Register light theme
	latte := createTestTheme(t, ThemeLatte, ThemeMetadata{
		displayName: "Catppuccin Latte",
		author:      "Catppuccin",
		variant:     ThemeVariantLight,
	})

	err := registry.Register(mocha)
	require.NoError(t, err)
	err = registry.Register(frappe)
	require.NoError(t, err)
	err = registry.Register(macchiato)
	require.NoError(t, err)
	err = registry.Register(gohan)
	require.NoError(t, err)
	err = registry.Register(latte)
	require.NoError(t, err)

	t.Run("list dark themes", func(t *testing.T) {
		darkThemes := registry.ListByVariant(ThemeVariantDark)
		assert.Len(t, darkThemes, 4)

		for _, theme := range darkThemes {
			assert.Equal(t, ThemeVariantDark, theme.Variant())
		}
	})

	t.Run("list light themes", func(t *testing.T) {
		lightThemes := registry.ListByVariant(ThemeVariantLight)
		assert.Len(t, lightThemes, 1)
		assert.Equal(t, ThemeLatte, lightThemes[0].Name())
	})

	t.Run("list invalid variant", func(t *testing.T) {
		themes := registry.ListByVariant(ThemeVariant("invalid"))
		assert.Empty(t, themes)
	})
}

// Helper function to create test themes
func createTestTheme(t *testing.T, name ThemeName, metadata ThemeMetadata) *Theme {
	t.Helper()

	colorScheme := ColorScheme{
		base:    Color("#1e1e2e"),
		surface: Color("#313244"),
		text:    Color("#cdd6f4"),
		mauve:   Color("#cba6f7"),
		red:     Color("#f38ba8"),
		green:   Color("#a6e3a1"),
	}

	theme, err := NewTheme(name, metadata, colorScheme)
	require.NoError(t, err)

	return theme
}
