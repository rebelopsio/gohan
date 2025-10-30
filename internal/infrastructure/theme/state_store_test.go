package theme_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/theme"
	themeInfra "github.com/rebelopsio/gohan/internal/infrastructure/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// Theme State Store - TDD Unit Tests
// ========================================

func TestFileThemeStateStore_Save(t *testing.T) {
	tests := []struct {
		name      string
		themeName theme.ThemeName
		wantErr   bool
	}{
		{
			name:      "saves theme state successfully",
			themeName: theme.ThemeLatte,
			wantErr:   false,
		},
		{
			name:      "saves different theme",
			themeName: theme.ThemeFrappe,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			stateFile := filepath.Join(tmpDir, "theme-state.json")

			store := themeInfra.NewFileThemeStateStore(stateFile)
			ctx := context.Background()

			state := &themeInfra.ThemeState{
				ThemeName:   tt.themeName,
				ThemeVariant: theme.ThemeVariantDark,
				SetAt:       time.Now(),
			}

			err := store.Save(ctx, state)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify file exists
				_, statErr := os.Stat(stateFile)
				require.NoError(t, statErr)
			}
		})
	}
}

func TestFileThemeStateStore_Load(t *testing.T) {
	t.Run("loads existing state", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "theme-state.json")

		store := themeInfra.NewFileThemeStateStore(stateFile)
		ctx := context.Background()

		// Save state first
		originalState := &themeInfra.ThemeState{
			ThemeName:   theme.ThemeMocha,
			ThemeVariant: theme.ThemeVariantDark,
			SetAt:       time.Now(),
		}

		err := store.Save(ctx, originalState)
		require.NoError(t, err)

		// Load state
		loadedState, err := store.Load(ctx)
		require.NoError(t, err)
		require.NotNil(t, loadedState)

		assert.Equal(t, originalState.ThemeName, loadedState.ThemeName)
		assert.Equal(t, originalState.ThemeVariant, loadedState.ThemeVariant)
		assert.WithinDuration(t, originalState.SetAt, loadedState.SetAt, time.Second)
	})

	t.Run("returns error when state file doesn't exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "nonexistent.json")

		store := themeInfra.NewFileThemeStateStore(stateFile)
		ctx := context.Background()

		loadedState, err := store.Load(ctx)
		require.Error(t, err)
		assert.Nil(t, loadedState)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("returns error for corrupted state file", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "corrupted.json")

		// Write invalid JSON
		err := os.WriteFile(stateFile, []byte("invalid json {{{"), 0644)
		require.NoError(t, err)

		store := themeInfra.NewFileThemeStateStore(stateFile)
		ctx := context.Background()

		loadedState, err := store.Load(ctx)
		require.Error(t, err)
		assert.Nil(t, loadedState)
	})
}

func TestFileThemeStateStore_Exists(t *testing.T) {
	t.Run("returns true when state exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "theme-state.json")

		store := themeInfra.NewFileThemeStateStore(stateFile)
		ctx := context.Background()

		// Save state
		state := &themeInfra.ThemeState{
			ThemeName:   theme.ThemeLatte,
			ThemeVariant: theme.ThemeVariantLight,
			SetAt:       time.Now(),
		}

		err := store.Save(ctx, state)
		require.NoError(t, err)

		// Check exists
		exists, err := store.Exists(ctx)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns false when state doesn't exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "nonexistent.json")

		store := themeInfra.NewFileThemeStateStore(stateFile)
		ctx := context.Background()

		exists, err := store.Exists(ctx)
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestFileThemeStateStore_RoundTrip(t *testing.T) {
	t.Run("save and load preserves all data", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "theme-state.json")

		store := themeInfra.NewFileThemeStateStore(stateFile)
		ctx := context.Background()

		testCases := []struct {
			themeName   theme.ThemeName
			themeVariant theme.ThemeVariant
		}{
			{theme.ThemeMocha, theme.ThemeVariantDark},
			{theme.ThemeLatte, theme.ThemeVariantLight},
			{theme.ThemeFrappe, theme.ThemeVariantDark},
			{theme.ThemeMacchiato, theme.ThemeVariantDark},
		}

		for _, tc := range testCases {
			originalState := &themeInfra.ThemeState{
				ThemeName:   tc.themeName,
				ThemeVariant: tc.themeVariant,
				SetAt:       time.Now(),
			}

			err := store.Save(ctx, originalState)
			require.NoError(t, err)

			loadedState, err := store.Load(ctx)
			require.NoError(t, err)

			assert.Equal(t, originalState.ThemeName, loadedState.ThemeName)
			assert.Equal(t, originalState.ThemeVariant, loadedState.ThemeVariant)
			assert.WithinDuration(t, originalState.SetAt, loadedState.SetAt, time.Second)
		}
	})
}

func TestFileThemeStateStore_ConcurrentAccess(t *testing.T) {
	t.Run("handles concurrent saves safely", func(t *testing.T) {
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "theme-state.json")

		store := themeInfra.NewFileThemeStateStore(stateFile)
		ctx := context.Background()

		// This is a basic test - real file locking would need OS-level locks
		// For now, we just verify no panics occur
		done := make(chan bool, 2)

		go func() {
			state := &themeInfra.ThemeState{
				ThemeName:   theme.ThemeLatte,
				ThemeVariant: theme.ThemeVariantLight,
				SetAt:       time.Now(),
			}
			_ = store.Save(ctx, state)
			done <- true
		}()

		go func() {
			state := &themeInfra.ThemeState{
				ThemeName:   theme.ThemeMocha,
				ThemeVariant: theme.ThemeVariantDark,
				SetAt:       time.Now(),
			}
			_ = store.Save(ctx, state)
			done <- true
		}()

		<-done
		<-done

		// Verify file is valid JSON
		loadedState, err := store.Load(ctx)
		require.NoError(t, err)
		assert.NotNil(t, loadedState)
	})
}
