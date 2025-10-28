package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeStandardThemes(t *testing.T) {
	registry := NewThemeRegistry()

	err := InitializeStandardThemes(registry)
	require.NoError(t, err)

	t.Run("all 5 themes are registered", func(t *testing.T) {
		themes := registry.ListAll()
		assert.Len(t, themes, 5)

		themeNames := make([]ThemeName, len(themes))
		for i, theme := range themes {
			themeNames[i] = theme.Name()
		}

		assert.Contains(t, themeNames, ThemeMocha)
		assert.Contains(t, themeNames, ThemeLatte)
		assert.Contains(t, themeNames, ThemeFrappe)
		assert.Contains(t, themeNames, ThemeMacchiato)
		assert.Contains(t, themeNames, ThemeGohan)
	})

	t.Run("mocha is the default active theme", func(t *testing.T) {
		active, err := registry.GetActive(nil)
		require.NoError(t, err)
		assert.Equal(t, ThemeMocha, active.Name())
	})

	t.Run("4 dark themes", func(t *testing.T) {
		darkThemes := registry.ListByVariant(ThemeVariantDark)
		assert.Len(t, darkThemes, 4)
	})

	t.Run("1 light theme", func(t *testing.T) {
		lightThemes := registry.ListByVariant(ThemeVariantLight)
		assert.Len(t, lightThemes, 1)
		assert.Equal(t, ThemeLatte, lightThemes[0].Name())
	})
}

func TestMochaTheme(t *testing.T) {
	theme := createMochaTheme()

	t.Run("has correct metadata", func(t *testing.T) {
		assert.Equal(t, ThemeMocha, theme.Name())
		assert.Equal(t, "Catppuccin Mocha", theme.DisplayName())
		assert.Equal(t, "Catppuccin", theme.Author())
		assert.Equal(t, ThemeVariantDark, theme.Variant())
		assert.True(t, theme.IsDark())
		assert.False(t, theme.IsLight())
	})

	t.Run("has complete color scheme", func(t *testing.T) {
		cs := theme.ColorScheme()

		// Verify base colors
		assert.Equal(t, Color("#1e1e2e"), cs.Base())
		assert.Equal(t, Color("#313244"), cs.Surface())
		assert.Equal(t, Color("#cdd6f4"), cs.Text())

		// Verify some accent colors
		assert.Equal(t, Color("#cba6f7"), cs.Mauve())
		assert.Equal(t, Color("#f38ba8"), cs.Red())
		assert.Equal(t, Color("#a6e3a1"), cs.Green())
		assert.Equal(t, Color("#89b4fa"), cs.Blue())
	})

	t.Run("all colors are valid hex codes", func(t *testing.T) {
		cs := theme.ColorScheme()

		colors := []Color{
			cs.Base(), cs.Surface(), cs.Overlay(), cs.Text(), cs.Subtext(),
			cs.Rosewater(), cs.Flamingo(), cs.Pink(), cs.Mauve(),
			cs.Red(), cs.Maroon(), cs.Peach(), cs.Yellow(),
			cs.Green(), cs.Teal(), cs.Sky(), cs.Sapphire(),
			cs.Blue(), cs.Lavender(),
		}

		for _, color := range colors {
			if color != "" {
				err := color.Validate()
				assert.NoError(t, err, "Color %s should be valid", color)
			}
		}
	})
}

func TestLatteTheme(t *testing.T) {
	theme := createLatteTheme()

	t.Run("has correct metadata", func(t *testing.T) {
		assert.Equal(t, ThemeLatte, theme.Name())
		assert.Equal(t, "Catppuccin Latte", theme.DisplayName())
		assert.Equal(t, "Catppuccin", theme.Author())
		assert.Equal(t, ThemeVariantLight, theme.Variant())
		assert.False(t, theme.IsDark())
		assert.True(t, theme.IsLight())
	})

	t.Run("has light color scheme", func(t *testing.T) {
		cs := theme.ColorScheme()

		// Light theme should have light base
		assert.Equal(t, Color("#eff1f5"), cs.Base())
		// Dark text on light background
		assert.Equal(t, Color("#4c4f69"), cs.Text())
	})
}

func TestFrappeTheme(t *testing.T) {
	theme := createFrappeTheme()

	t.Run("has correct metadata", func(t *testing.T) {
		assert.Equal(t, ThemeFrappe, theme.Name())
		assert.Equal(t, "Catppuccin Frappe", theme.DisplayName())
		assert.Equal(t, "Catppuccin", theme.Author())
		assert.Equal(t, ThemeVariantDark, theme.Variant())
	})

	t.Run("has distinct colors from mocha", func(t *testing.T) {
		mocha := createMochaTheme()

		// Frappe should have different base color than Mocha
		assert.NotEqual(t, mocha.ColorScheme().Base(), theme.ColorScheme().Base())
	})
}

func TestMacchiatoTheme(t *testing.T) {
	theme := createMacchiatoTheme()

	t.Run("has correct metadata", func(t *testing.T) {
		assert.Equal(t, ThemeMacchiato, theme.Name())
		assert.Equal(t, "Catppuccin Macchiato", theme.DisplayName())
		assert.Equal(t, "Catppuccin", theme.Author())
		assert.Equal(t, ThemeVariantDark, theme.Variant())
	})
}

func TestGohanTheme(t *testing.T) {
	theme := createGohanTheme()

	t.Run("has correct metadata", func(t *testing.T) {
		assert.Equal(t, ThemeGohan, theme.Name())
		assert.Equal(t, "Gohan", theme.DisplayName())
		assert.Equal(t, "Gohan Team", theme.Author())
		assert.Equal(t, ThemeVariantDark, theme.Variant())
	})

	t.Run("has complete color scheme", func(t *testing.T) {
		cs := theme.ColorScheme()

		// Verify it has all colors defined
		assert.NotEmpty(t, cs.Base())
		assert.NotEmpty(t, cs.Text())
		assert.NotEmpty(t, cs.Mauve())
	})
}

func TestAllThemesAreUnique(t *testing.T) {
	themes := []*Theme{
		createMochaTheme(),
		createLatteTheme(),
		createFrappeTheme(),
		createMacchiatoTheme(),
		createGohanTheme(),
	}

	t.Run("all themes have unique names", func(t *testing.T) {
		names := make(map[ThemeName]bool)
		for _, theme := range themes {
			assert.False(t, names[theme.Name()], "Duplicate theme name: %s", theme.Name())
			names[theme.Name()] = true
		}
	})

	t.Run("catppuccin themes have unique base colors", func(t *testing.T) {
		// Catppuccin themes (Mocha, Frappe, Macchiato) should have distinct base colors
		// Gohan may share colors with others as it's a custom brand theme
		catppuccinThemes := []*Theme{
			createMochaTheme(),
			createFrappeTheme(),
			createMacchiatoTheme(),
		}

		baseColors := make(map[Color]bool)
		for _, theme := range catppuccinThemes {
			base := theme.ColorScheme().Base()
			assert.False(t, baseColors[base], "Duplicate base color in Catppuccin themes: %s", base)
			baseColors[base] = true
		}
	})
}
