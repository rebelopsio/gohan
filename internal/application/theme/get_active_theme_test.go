package theme

import (
	"context"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetActiveThemeUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		setupRegistry func() theme.ThemeRegistry
		wantTheme     string
		wantErr       bool
	}{
		{
			name: "get default active theme (mocha)",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				return registry
			},
			wantTheme: "mocha",
			wantErr:   false,
		},
		{
			name: "get explicitly set active theme",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				err = registry.SetActive(context.Background(), theme.ThemeLatte)
				require.NoError(t, err)
				return registry
			},
			wantTheme: "latte",
			wantErr:   false,
		},
		{
			name: "no themes registered",
			setupRegistry: func() theme.ThemeRegistry {
				return theme.NewThemeRegistry()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := tt.setupRegistry()
			useCase := NewGetActiveThemeUseCase(registry)

			themeInfo, err := useCase.Execute(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantTheme, themeInfo.Name)
				assert.True(t, themeInfo.IsActive)
				assert.NotEmpty(t, themeInfo.DisplayName)
				assert.NotEmpty(t, themeInfo.ColorScheme)
			}
		})
	}
}
