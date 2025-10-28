package theme

import (
	"context"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreviewThemeUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		themeName     string
		setupRegistry func() theme.ThemeRegistry
		wantErr       bool
		checkResult   func(*testing.T, *ThemePreview)
	}{
		{
			name:      "preview existing theme",
			themeName: "mocha",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				return registry
			},
			wantErr: false,
			checkResult: func(t *testing.T, preview *ThemePreview) {
				assert.Equal(t, "mocha", preview.Name)
				assert.Equal(t, "Catppuccin Mocha", preview.DisplayName)
				assert.Equal(t, "Catppuccin", preview.Author)
				assert.Equal(t, "dark", preview.Variant)
				assert.NotEmpty(t, preview.ColorScheme)
				assert.NotEmpty(t, preview.PreviewText)
			},
		},
		{
			name:      "preview light theme",
			themeName: "latte",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				return registry
			},
			wantErr: false,
			checkResult: func(t *testing.T, preview *ThemePreview) {
				assert.Equal(t, "latte", preview.Name)
				assert.Equal(t, "light", preview.Variant)
			},
		},
		{
			name:      "preview non-existent theme",
			themeName: "nonexistent",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				return registry
			},
			wantErr: true,
		},
		{
			name:      "preview has color scheme",
			themeName: "mocha",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				return registry
			},
			wantErr: false,
			checkResult: func(t *testing.T, preview *ThemePreview) {
				assert.Contains(t, preview.ColorScheme, "base")
				assert.Contains(t, preview.ColorScheme, "text")
				assert.Contains(t, preview.ColorScheme, "mauve")
				assert.Regexp(t, "^#[0-9A-Fa-f]{6}$", preview.ColorScheme["base"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := tt.setupRegistry()
			useCase := NewPreviewThemeUseCase(registry)

			preview, err := useCase.Execute(context.Background(), tt.themeName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, preview)
				if tt.checkResult != nil {
					tt.checkResult(t, preview)
				}
			}
		})
	}
}

func TestPreviewThemeUseCase_PreviewText(t *testing.T) {
	registry := theme.NewThemeRegistry()
	err := theme.InitializeStandardThemes(registry)
	require.NoError(t, err)

	useCase := NewPreviewThemeUseCase(registry)

	t.Run("dark theme preview text", func(t *testing.T) {
		preview, err := useCase.Execute(context.Background(), "mocha")
		require.NoError(t, err)

		assert.Contains(t, preview.PreviewText, "dark")
		assert.NotEmpty(t, preview.PreviewText)
	})

	t.Run("light theme preview text", func(t *testing.T) {
		preview, err := useCase.Execute(context.Background(), "latte")
		require.NoError(t, err)

		assert.Contains(t, preview.PreviewText, "light")
		assert.NotEmpty(t, preview.PreviewText)
	})
}
