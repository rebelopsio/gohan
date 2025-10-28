package theme

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTheme(t *testing.T) {
	tests := []struct {
		name        string
		themeName   ThemeName
		metadata    ThemeMetadata
		colorScheme ColorScheme
		wantErr     bool
		checkResult func(*testing.T, *Theme)
	}{
		{
			name:      "valid theme creation",
			themeName: ThemeMocha,
			metadata: ThemeMetadata{
				displayName: "Catppuccin Mocha",
				author:      "Catppuccin",
				description: "Warm dark theme",
				variant:     ThemeVariantDark,
			},
			colorScheme: ColorScheme{
				base:    Color("#1e1e2e"),
				surface: Color("#313244"),
				text:    Color("#cdd6f4"),
			},
			wantErr: false,
			checkResult: func(t *testing.T, theme *Theme) {
				assert.Equal(t, ThemeMocha, theme.Name())
				assert.Equal(t, "Catppuccin Mocha", theme.DisplayName())
				assert.Equal(t, "Catppuccin", theme.Author())
				assert.Equal(t, ThemeVariantDark, theme.Variant())
				assert.NotZero(t, theme.CreatedAt())
			},
		},
		{
			name:      "empty theme name",
			themeName: ThemeName(""),
			metadata: ThemeMetadata{
				displayName: "Test Theme",
				author:      "Test Author",
				variant:     ThemeVariantDark,
			},
			colorScheme: ColorScheme{},
			wantErr:     true,
		},
		{
			name:      "empty display name",
			themeName: ThemeMocha,
			metadata: ThemeMetadata{
				displayName: "",
				author:      "Catppuccin",
				variant:     ThemeVariantDark,
			},
			colorScheme: ColorScheme{},
			wantErr:     true,
		},
		{
			name:      "empty author",
			themeName: ThemeMocha,
			metadata: ThemeMetadata{
				displayName: "Catppuccin Mocha",
				author:      "",
				variant:     ThemeVariantDark,
			},
			colorScheme: ColorScheme{},
			wantErr:     true,
		},
		{
			name:      "invalid variant",
			themeName: ThemeMocha,
			metadata: ThemeMetadata{
				displayName: "Catppuccin Mocha",
				author:      "Catppuccin",
				variant:     ThemeVariant(""),
			},
			colorScheme: ColorScheme{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme, err := NewTheme(tt.themeName, tt.metadata, tt.colorScheme)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, theme)
			} else {
				require.NoError(t, err)
				require.NotNil(t, theme)
				if tt.checkResult != nil {
					tt.checkResult(t, theme)
				}
			}
		})
	}
}

func TestThemeConstants(t *testing.T) {
	t.Run("theme names are defined", func(t *testing.T) {
		assert.Equal(t, ThemeName("mocha"), ThemeMocha)
		assert.Equal(t, ThemeName("latte"), ThemeLatte)
		assert.Equal(t, ThemeName("frappe"), ThemeFrappe)
		assert.Equal(t, ThemeName("macchiato"), ThemeMacchiato)
		assert.Equal(t, ThemeName("gohan"), ThemeGohan)
	})

	t.Run("theme variants are defined", func(t *testing.T) {
		assert.Equal(t, ThemeVariant("dark"), ThemeVariantDark)
		assert.Equal(t, ThemeVariant("light"), ThemeVariantLight)
	})
}

func TestColorScheme(t *testing.T) {
	t.Run("valid color scheme", func(t *testing.T) {
		cs := ColorScheme{
			base:    Color("#1e1e2e"),
			surface: Color("#313244"),
			text:    Color("#cdd6f4"),
			mauve:   Color("#cba6f7"),
			red:     Color("#f38ba8"),
			green:   Color("#a6e3a1"),
		}

		assert.Equal(t, Color("#1e1e2e"), cs.Base())
		assert.Equal(t, Color("#313244"), cs.Surface())
		assert.Equal(t, Color("#cdd6f4"), cs.Text())
		assert.Equal(t, Color("#cba6f7"), cs.Mauve())
		assert.Equal(t, Color("#f38ba8"), cs.Red())
		assert.Equal(t, Color("#a6e3a1"), cs.Green())
	})

	t.Run("color validation", func(t *testing.T) {
		tests := []struct {
			name    string
			color   Color
			wantErr bool
		}{
			{
				name:    "valid hex color",
				color:   Color("#1e1e2e"),
				wantErr: false,
			},
			{
				name:    "valid hex color uppercase",
				color:   Color("#1E1E2E"),
				wantErr: false,
			},
			{
				name:    "invalid hex color - missing #",
				color:   Color("1e1e2e"),
				wantErr: true,
			},
			{
				name:    "invalid hex color - too short",
				color:   Color("#1e1e2"),
				wantErr: true,
			},
			{
				name:    "invalid hex color - invalid characters",
				color:   Color("#zzzzzz"),
				wantErr: true,
			},
			{
				name:    "empty color",
				color:   Color(""),
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.color.Validate()
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestThemeMetadata(t *testing.T) {
	t.Run("valid metadata", func(t *testing.T) {
		meta := ThemeMetadata{
			displayName: "Catppuccin Mocha",
			author:      "Catppuccin",
			description: "Warm dark theme",
			variant:     ThemeVariantDark,
			previewURL:  "https://example.com/preview.png",
		}

		assert.Equal(t, "Catppuccin Mocha", meta.DisplayName())
		assert.Equal(t, "Catppuccin", meta.Author())
		assert.Equal(t, "Warm dark theme", meta.Description())
		assert.Equal(t, ThemeVariantDark, meta.Variant())
		assert.Equal(t, "https://example.com/preview.png", meta.PreviewURL())
	})

	t.Run("variant is light or dark", func(t *testing.T) {
		darkMeta := ThemeMetadata{
			displayName: "Dark Theme",
			author:      "Test",
			variant:     ThemeVariantDark,
		}
		assert.True(t, darkMeta.IsDark())
		assert.False(t, darkMeta.IsLight())

		lightMeta := ThemeMetadata{
			displayName: "Light Theme",
			author:      "Test",
			variant:     ThemeVariantLight,
		}
		assert.False(t, lightMeta.IsDark())
		assert.True(t, lightMeta.IsLight())
	})
}

func TestTheme_Equality(t *testing.T) {
	now := time.Now()

	theme1, err := NewTheme(
		ThemeMocha,
		ThemeMetadata{
			displayName: "Catppuccin Mocha",
			author:      "Catppuccin",
			variant:     ThemeVariantDark,
		},
		ColorScheme{
			base: Color("#1e1e2e"),
			text: Color("#cdd6f4"),
		},
	)
	require.NoError(t, err)
	theme1.createdAt = now

	theme2, err := NewTheme(
		ThemeMocha,
		ThemeMetadata{
			displayName: "Catppuccin Mocha",
			author:      "Catppuccin",
			variant:     ThemeVariantDark,
		},
		ColorScheme{
			base: Color("#1e1e2e"),
			text: Color("#cdd6f4"),
		},
	)
	require.NoError(t, err)
	theme2.createdAt = now

	theme3, err := NewTheme(
		ThemeLatte,
		ThemeMetadata{
			displayName: "Catppuccin Latte",
			author:      "Catppuccin",
			variant:     ThemeVariantLight,
		},
		ColorScheme{
			base: Color("#eff1f5"),
			text: Color("#4c4f69"),
		},
	)
	require.NoError(t, err)

	t.Run("same themes are equal", func(t *testing.T) {
		assert.True(t, theme1.Equals(theme2))
		assert.True(t, theme2.Equals(theme1))
	})

	t.Run("different themes are not equal", func(t *testing.T) {
		assert.False(t, theme1.Equals(theme3))
		assert.False(t, theme3.Equals(theme1))
	})
}

func TestTheme_ColorScheme(t *testing.T) {
	colorScheme := ColorScheme{
		base:    Color("#1e1e2e"),
		surface: Color("#313244"),
		text:    Color("#cdd6f4"),
		mauve:   Color("#cba6f7"),
	}

	theme, err := NewTheme(
		ThemeMocha,
		ThemeMetadata{
			displayName: "Catppuccin Mocha",
			author:      "Catppuccin",
			variant:     ThemeVariantDark,
		},
		colorScheme,
	)
	require.NoError(t, err)

	t.Run("color scheme is accessible", func(t *testing.T) {
		cs := theme.ColorScheme()
		assert.Equal(t, Color("#1e1e2e"), cs.Base())
		assert.Equal(t, Color("#313244"), cs.Surface())
		assert.Equal(t, Color("#cdd6f4"), cs.Text())
		assert.Equal(t, Color("#cba6f7"), cs.Mauve())
	})
}
