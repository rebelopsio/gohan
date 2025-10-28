package theme

import (
	"context"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListThemesUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		setupRegistry func() theme.ThemeRegistry
		wantCount     int
		checkResult   func(*testing.T, []ThemeInfo)
	}{
		{
			name: "list all themes",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				return registry
			},
			wantCount: 5,
			checkResult: func(t *testing.T, themes []ThemeInfo) {
				names := make([]string, len(themes))
				for i, th := range themes {
					names[i] = th.Name
				}
				assert.Contains(t, names, "mocha")
				assert.Contains(t, names, "latte")
				assert.Contains(t, names, "frappe")
				assert.Contains(t, names, "macchiato")
				assert.Contains(t, names, "gohan")
			},
		},
		{
			name: "empty registry",
			setupRegistry: func() theme.ThemeRegistry {
				return theme.NewThemeRegistry()
			},
			wantCount: 0,
		},
		{
			name: "active theme is marked",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				// Mocha is default active
				return registry
			},
			wantCount: 5,
			checkResult: func(t *testing.T, themes []ThemeInfo) {
				var activeCount int
				var mochaActive bool
				for _, th := range themes {
					if th.IsActive {
						activeCount++
						if th.Name == "mocha" {
							mochaActive = true
						}
					}
				}
				assert.Equal(t, 1, activeCount, "exactly one theme should be active")
				assert.True(t, mochaActive, "mocha should be the active theme")
			},
		},
		{
			name: "theme info contains all fields",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				return registry
			},
			wantCount: 5,
			checkResult: func(t *testing.T, themes []ThemeInfo) {
				for _, th := range themes {
					assert.NotEmpty(t, th.Name, "theme should have name")
					assert.NotEmpty(t, th.DisplayName, "theme should have display name")
					assert.NotEmpty(t, th.Author, "theme should have author")
					assert.Contains(t, []string{"dark", "light"}, th.Variant)
					assert.NotEmpty(t, th.ColorScheme, "theme should have color scheme")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := tt.setupRegistry()
			useCase := NewListThemesUseCase(registry)

			themes, err := useCase.Execute(context.Background())

			require.NoError(t, err)
			assert.Len(t, themes, tt.wantCount)

			if tt.checkResult != nil {
				tt.checkResult(t, themes)
			}
		})
	}
}

func TestListThemesUseCase_ColorSchemeMapping(t *testing.T) {
	registry := theme.NewThemeRegistry()
	err := theme.InitializeStandardThemes(registry)
	require.NoError(t, err)

	useCase := NewListThemesUseCase(registry)
	themes, err := useCase.Execute(context.Background())
	require.NoError(t, err)

	t.Run("color scheme is mapped correctly", func(t *testing.T) {
		// Find mocha theme
		var mochaInfo *ThemeInfo
		for i := range themes {
			if themes[i].Name == "mocha" {
				mochaInfo = &themes[i]
				break
			}
		}
		require.NotNil(t, mochaInfo)

		// Check key colors are present
		assert.Contains(t, mochaInfo.ColorScheme, "base")
		assert.Contains(t, mochaInfo.ColorScheme, "text")
		assert.Contains(t, mochaInfo.ColorScheme, "mauve")
		assert.Contains(t, mochaInfo.ColorScheme, "red")
		assert.Contains(t, mochaInfo.ColorScheme, "green")

		// Verify color format
		assert.Regexp(t, "^#[0-9A-Fa-f]{6}$", mochaInfo.ColorScheme["base"])
	})
}
