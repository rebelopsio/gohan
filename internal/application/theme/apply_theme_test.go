package theme

import (
	"context"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyThemeUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		themeName     string
		setupRegistry func() theme.ThemeRegistry
		wantErr       bool
		checkResult   func(*testing.T, *ApplyThemeResult)
	}{
		{
			name:      "apply valid theme",
			themeName: "latte",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				return registry
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *ApplyThemeResult) {
				assert.True(t, result.Success)
				assert.Equal(t, "latte", result.ThemeName)
				assert.NotEmpty(t, result.Message)
			},
		},
		{
			name:      "apply non-existent theme",
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
			name:      "apply already active theme",
			themeName: "mocha",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				// Mocha is default active
				return registry
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *ApplyThemeResult) {
				assert.True(t, result.Success)
				assert.Contains(t, result.Message, "already active")
				assert.Empty(t, result.AffectedComponents) // No changes needed
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := tt.setupRegistry()
			useCase := NewApplyThemeUseCase(registry, nil, nil, nil) // nil applier, stateStore, historyStore for now

			result, err := useCase.Execute(context.Background(), tt.themeName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}
		})
	}
}

func TestApplyThemeUseCase_SetsActiveTheme(t *testing.T) {
	registry := theme.NewThemeRegistry()
	err := theme.InitializeStandardThemes(registry)
	require.NoError(t, err)

	useCase := NewApplyThemeUseCase(registry, nil, nil, nil)

	t.Run("sets theme as active after application", func(t *testing.T) {
		ctx := context.Background()

		// Initially mocha is active
		active, err := registry.GetActive(ctx)
		require.NoError(t, err)
		assert.Equal(t, theme.ThemeMocha, active.Name())

		// Apply latte
		result, err := useCase.Execute(ctx, "latte")
		require.NoError(t, err)
		assert.True(t, result.Success)

		// Verify latte is now active
		active, err = registry.GetActive(ctx)
		require.NoError(t, err)
		assert.Equal(t, theme.ThemeLatte, active.Name())
	})
}
