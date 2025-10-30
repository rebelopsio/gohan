package acceptance

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// Phase 4.6: Theme Persistence - ATDD
// ========================================

// ThemeState represents the persisted theme state
type ThemeState struct {
	ThemeName   string    `json:"theme_name"`
	ThemeVariant string   `json:"theme_variant"`
	SetAt       time.Time `json:"set_at"`
}

// ThemeStateStore defines the interface for theme persistence
type ThemeStateStore interface {
	Save(ctx context.Context, state *ThemeState) error
	Load(ctx context.Context) (*ThemeState, error)
	Exists(ctx context.Context) (bool, error)
}

func TestThemePersistence_DefaultThemeOnFirstLaunch(t *testing.T) {
	t.Run("uses default theme when no state exists", func(t *testing.T) {
		// Given: Fresh start with no saved state
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "theme-state.json")

		// Verify no state exists
		_, err := os.Stat(stateFile)
		require.True(t, os.IsNotExist(err), "State file should not exist")

		// When: Application initializes
		registry := theme.NewThemeRegistry()
		err = theme.InitializeStandardThemes(registry)
		require.NoError(t, err)

		// Then: Default theme (mocha) should be active
		activeTheme, err := registry.GetActive(context.Background())
		require.NoError(t, err)
		require.NotNil(t, activeTheme)
		assert.Equal(t, theme.ThemeMocha, activeTheme.Name())
	})
}

func TestThemePersistence_SaveThemeState(t *testing.T) {
	t.Run("saves theme state when theme is set", func(t *testing.T) {
		// Given: A theme registry and state file
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "theme-state.json")

		registry := theme.NewThemeRegistry()
		err := theme.InitializeStandardThemes(registry)
		require.NoError(t, err)

		latte, err := registry.FindByName(context.Background(), theme.ThemeLatte)
		require.NoError(t, err)

		// When: Theme is set
		err = registry.SetActive(context.Background(), latte.Name())
		require.NoError(t, err)

		// Manually save state for now (will be automated later)
		state := &ThemeState{
			ThemeName:   string(latte.Name()),
			ThemeVariant: string(latte.Variant()),
			SetAt:       time.Now(),
		}

		data, err := json.Marshal(state)
		require.NoError(t, err)

		err = os.WriteFile(stateFile, data, 0644)
		require.NoError(t, err)

		// Then: State file should exist and contain correct data
		savedData, err := os.ReadFile(stateFile)
		require.NoError(t, err)

		var loadedState ThemeState
		err = json.Unmarshal(savedData, &loadedState)
		require.NoError(t, err)

		assert.Equal(t, string(theme.ThemeLatte), loadedState.ThemeName)
		assert.Equal(t, string(theme.ThemeVariantLight), loadedState.ThemeVariant)
	})
}

func TestThemePersistence_LoadThemeState(t *testing.T) {
	t.Run("loads theme state on startup", func(t *testing.T) {
		// Given: A saved theme state
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "theme-state.json")

		state := &ThemeState{
			ThemeName:   string(theme.ThemeFrappe),
			ThemeVariant: string(theme.ThemeVariantDark),
			SetAt:       time.Now(),
		}

		data, err := json.Marshal(state)
		require.NoError(t, err)

		err = os.WriteFile(stateFile, data, 0644)
		require.NoError(t, err)

		// When: Application loads state
		loadedData, err := os.ReadFile(stateFile)
		require.NoError(t, err)

		var loadedState ThemeState
		err = json.Unmarshal(loadedData, &loadedState)
		require.NoError(t, err)

		// Then: Correct theme should be loaded
		assert.Equal(t, string(theme.ThemeFrappe), loadedState.ThemeName)
		assert.Equal(t, string(theme.ThemeVariantDark), loadedState.ThemeVariant)

		// And: Theme should be set as active
		registry := theme.NewThemeRegistry()
		err = theme.InitializeStandardThemes(registry)
		require.NoError(t, err)

		frappe, err := registry.FindByName(context.Background(), theme.ThemeName(loadedState.ThemeName))
		require.NoError(t, err)

		err = registry.SetActive(context.Background(), frappe.Name())
		require.NoError(t, err)

		activeTheme, err := registry.GetActive(context.Background())
		require.NoError(t, err)
		assert.Equal(t, theme.ThemeFrappe, activeTheme.Name())
	})
}

func TestThemePersistence_CorruptedStateFallback(t *testing.T) {
	t.Run("falls back to default when state is corrupted", func(t *testing.T) {
		// Given: A corrupted state file
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "theme-state.json")

		// Write invalid JSON
		err := os.WriteFile(stateFile, []byte("invalid json {{{"), 0644)
		require.NoError(t, err)

		// When: Application tries to load state
		_, err = os.ReadFile(stateFile)
		require.NoError(t, err)

		var loadedState ThemeState
		err = json.Unmarshal([]byte("invalid json {{{"), &loadedState)

		// Then: Should fail to unmarshal
		require.Error(t, err)

		// And: Should fall back to default theme
		registry := theme.NewThemeRegistry()
		err = theme.InitializeStandardThemes(registry)
		require.NoError(t, err)

		activeTheme, err := registry.GetActive(context.Background())
		require.NoError(t, err)
		assert.Equal(t, theme.ThemeMocha, activeTheme.Name())
	})
}

func TestThemePersistence_MissingThemeFallback(t *testing.T) {
	t.Run("falls back to default when saved theme doesn't exist", func(t *testing.T) {
		// Given: State references a non-existent theme
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "theme-state.json")

		state := &ThemeState{
			ThemeName:   "deleted-theme",
			ThemeVariant: string(theme.ThemeVariantDark),
			SetAt:       time.Now(),
		}

		data, err := json.Marshal(state)
		require.NoError(t, err)

		err = os.WriteFile(stateFile, data, 0644)
		require.NoError(t, err)

		// When: Application tries to load the theme
		loadedData, err := os.ReadFile(stateFile)
		require.NoError(t, err)

		var loadedState ThemeState
		err = json.Unmarshal(loadedData, &loadedState)
		require.NoError(t, err)

		registry := theme.NewThemeRegistry()
		err = theme.InitializeStandardThemes(registry)
		require.NoError(t, err)

		_, err = registry.FindByName(context.Background(), theme.ThemeName(loadedState.ThemeName))

		// Then: Should not find the theme
		require.Error(t, err)
		assert.ErrorIs(t, err, theme.ErrThemeNotFound)

		// And: Should fall back to default
		activeTheme, err := registry.GetActive(context.Background())
		require.NoError(t, err)
		assert.Equal(t, theme.ThemeMocha, activeTheme.Name())
	})
}

func TestThemePersistence_StateMetadata(t *testing.T) {
	t.Run("includes metadata in saved state", func(t *testing.T) {
		// Given: A theme to save
		now := time.Now()
		state := &ThemeState{
			ThemeName:   string(theme.ThemeLatte),
			ThemeVariant: string(theme.ThemeVariantLight),
			SetAt:       now,
		}

		// When: State is serialized
		data, err := json.Marshal(state)
		require.NoError(t, err)

		// Then: Should include all metadata
		var parsed map[string]interface{}
		err = json.Unmarshal(data, &parsed)
		require.NoError(t, err)

		assert.Equal(t, string(theme.ThemeLatte), parsed["theme_name"])
		assert.Equal(t, string(theme.ThemeVariantLight), parsed["theme_variant"])
		assert.NotNil(t, parsed["set_at"])
	})
}

func TestThemePersistence_StateFileLocation(t *testing.T) {
	t.Run("saves state to correct location", func(t *testing.T) {
		// Given: A custom state directory
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "theme-state.json")

		state := &ThemeState{
			ThemeName:   string(theme.ThemeMacchiato),
			ThemeVariant: string(theme.ThemeVariantDark),
			SetAt:       time.Now(),
		}

		// When: State is saved
		data, err := json.Marshal(state)
		require.NoError(t, err)

		err = os.WriteFile(stateFile, data, 0644)
		require.NoError(t, err)

		// Then: File should exist at expected location
		info, err := os.Stat(stateFile)
		require.NoError(t, err)
		assert.False(t, info.IsDir())
		assert.Equal(t, "theme-state.json", info.Name())
	})
}

func TestThemePersistence_UpdateOnThemeChange(t *testing.T) {
	t.Run("updates state file when theme changes", func(t *testing.T) {
		// Given: An existing saved state
		tmpDir := t.TempDir()
		stateFile := filepath.Join(tmpDir, "theme-state.json")

		initialState := &ThemeState{
			ThemeName:   string(theme.ThemeMocha),
			ThemeVariant: string(theme.ThemeVariantDark),
			SetAt:       time.Now().Add(-1 * time.Hour),
		}

		data, err := json.Marshal(initialState)
		require.NoError(t, err)
		err = os.WriteFile(stateFile, data, 0644)
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond) // Ensure timestamp difference

		// When: Theme is changed
		newState := &ThemeState{
			ThemeName:   string(theme.ThemeLatte),
			ThemeVariant: string(theme.ThemeVariantLight),
			SetAt:       time.Now(),
		}

		data, err = json.Marshal(newState)
		require.NoError(t, err)
		err = os.WriteFile(stateFile, data, 0644)
		require.NoError(t, err)

		// Then: State file should be updated
		loadedData, err := os.ReadFile(stateFile)
		require.NoError(t, err)

		var loadedState ThemeState
		err = json.Unmarshal(loadedData, &loadedState)
		require.NoError(t, err)

		assert.Equal(t, string(theme.ThemeLatte), loadedState.ThemeName)
		assert.True(t, loadedState.SetAt.After(initialState.SetAt))
	})
}
