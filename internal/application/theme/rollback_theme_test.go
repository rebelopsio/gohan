package theme

import (
	"context"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/theme"
	themeInfra "github.com/rebelopsio/gohan/internal/infrastructure/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// Rollback Theme Use Case - TDD Tests
// ========================================

func TestRollbackThemeUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		setupRegistry func() theme.ThemeRegistry
		setupHistory  func(themeInfra.ThemeHistoryStore)
		wantErr       bool
		errorType     error
		checkResult   func(*testing.T, *RollbackThemeResult)
	}{
		{
			name: "rollback to previous theme",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				// Set latte as active
				err = registry.SetActive(context.Background(), theme.ThemeLatte)
				require.NoError(t, err)
				return registry
			},
			setupHistory: func(store themeInfra.ThemeHistoryStore) {
				// History: latte (current), mocha (previous)
				ctx := context.Background()
				_ = store.Add(ctx, theme.ThemeMocha)
				_ = store.Add(ctx, theme.ThemeLatte)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *RollbackThemeResult) {
				assert.True(t, result.Success)
				assert.Equal(t, theme.ThemeMocha, result.RestoredTheme)
				assert.Equal(t, theme.ThemeLatte, result.PreviousTheme)
				assert.NotEmpty(t, result.Message)
			},
		},
		{
			name: "error when no history available",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				return registry
			},
			setupHistory: func(store themeInfra.ThemeHistoryStore) {
				// No history
			},
			wantErr:   true,
			errorType: themeInfra.ErrNoThemeHistory,
		},
		{
			name: "error with only one theme in history",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				return registry
			},
			setupHistory: func(store themeInfra.ThemeHistoryStore) {
				// Only current theme, no previous
				ctx := context.Background()
				_ = store.Add(ctx, theme.ThemeLatte)
			},
			wantErr:   true,
			errorType: themeInfra.ErrNoThemeHistory,
		},
		{
			name: "rollback updates active theme",
			setupRegistry: func() theme.ThemeRegistry {
				registry := theme.NewThemeRegistry()
				err := theme.InitializeStandardThemes(registry)
				require.NoError(t, err)
				// Set frappe as active
				err = registry.SetActive(context.Background(), theme.ThemeFrappe)
				require.NoError(t, err)
				return registry
			},
			setupHistory: func(store themeInfra.ThemeHistoryStore) {
				// History: frappe (current), latte (previous), mocha (older)
				ctx := context.Background()
				_ = store.Add(ctx, theme.ThemeMocha)
				_ = store.Add(ctx, theme.ThemeLatte)
				_ = store.Add(ctx, theme.ThemeFrappe)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *RollbackThemeResult) {
				assert.True(t, result.Success)
				assert.Equal(t, theme.ThemeLatte, result.RestoredTheme)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for history file
			tmpDir := t.TempDir()
			historyFile := tmpDir + "/theme-history.json"
			historyStore := themeInfra.NewFileThemeHistoryStore(historyFile)

			// Setup
			registry := tt.setupRegistry()
			if tt.setupHistory != nil {
				tt.setupHistory(historyStore)
			}

			// Create use case
			useCase := NewRollbackThemeUseCase(registry, historyStore, nil, nil)

			// Execute
			result, err := useCase.Execute(context.Background())

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
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

func TestRollbackThemeUseCase_RemovesHistoryEntry(t *testing.T) {
	t.Run("removes most recent history entry after rollback", func(t *testing.T) {
		ctx := context.Background()

		// Setup registry
		registry := theme.NewThemeRegistry()
		err := theme.InitializeStandardThemes(registry)
		require.NoError(t, err)
		err = registry.SetActive(ctx, theme.ThemeLatte)
		require.NoError(t, err)

		// Setup history
		tmpDir := t.TempDir()
		historyFile := tmpDir + "/theme-history.json"
		historyStore := themeInfra.NewFileThemeHistoryStore(historyFile)
		_ = historyStore.Add(ctx, theme.ThemeMocha)
		_ = historyStore.Add(ctx, theme.ThemeLatte)

		// Verify initial history
		history, err := historyStore.GetHistory(ctx)
		require.NoError(t, err)
		assert.Len(t, history, 2)

		// Execute rollback
		useCase := NewRollbackThemeUseCase(registry, historyStore, nil, nil)
		_, err = useCase.Execute(ctx)
		require.NoError(t, err)

		// Verify history was updated (most recent removed)
		history, err = historyStore.GetHistory(ctx)
		require.NoError(t, err)
		assert.Len(t, history, 1)
		assert.Equal(t, theme.ThemeMocha, history[0])
	})
}

func TestRollbackThemeUseCase_SequentialRollbacks(t *testing.T) {
	t.Run("allows multiple sequential rollbacks", func(t *testing.T) {
		ctx := context.Background()

		// Setup registry
		registry := theme.NewThemeRegistry()
		err := theme.InitializeStandardThemes(registry)
		require.NoError(t, err)
		err = registry.SetActive(ctx, theme.ThemeMacchiato)
		require.NoError(t, err)

		// Setup history: macchiato (current), frappe, latte, mocha
		tmpDir := t.TempDir()
		historyFile := tmpDir + "/theme-history.json"
		historyStore := themeInfra.NewFileThemeHistoryStore(historyFile)
		_ = historyStore.Add(ctx, theme.ThemeMocha)
		_ = historyStore.Add(ctx, theme.ThemeLatte)
		_ = historyStore.Add(ctx, theme.ThemeFrappe)
		_ = historyStore.Add(ctx, theme.ThemeMacchiato)

		useCase := NewRollbackThemeUseCase(registry, historyStore, nil, nil)

		// First rollback: macchiato -> frappe
		result, err := useCase.Execute(ctx)
		require.NoError(t, err)
		assert.Equal(t, theme.ThemeFrappe, result.RestoredTheme)

		active, err := registry.GetActive(ctx)
		require.NoError(t, err)
		assert.Equal(t, theme.ThemeFrappe, active.Name())

		// Second rollback: frappe -> latte
		result, err = useCase.Execute(ctx)
		require.NoError(t, err)
		assert.Equal(t, theme.ThemeLatte, result.RestoredTheme)

		active, err = registry.GetActive(ctx)
		require.NoError(t, err)
		assert.Equal(t, theme.ThemeLatte, active.Name())

		// Third rollback: latte -> mocha
		result, err = useCase.Execute(ctx)
		require.NoError(t, err)
		assert.Equal(t, theme.ThemeMocha, result.RestoredTheme)

		active, err = registry.GetActive(ctx)
		require.NoError(t, err)
		assert.Equal(t, theme.ThemeMocha, active.Name())

		// Fourth rollback should fail (no more history)
		_, err = useCase.Execute(ctx)
		assert.Error(t, err)
		assert.ErrorIs(t, err, themeInfra.ErrNoThemeHistory)
	})
}
