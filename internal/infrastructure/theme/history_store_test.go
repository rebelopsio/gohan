package theme_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/theme"
	themeInfra "github.com/rebelopsio/gohan/internal/infrastructure/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// Theme History Store - TDD Unit Tests
// ========================================

func TestFileThemeHistoryStore_Add(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "theme-history.json")
	
	store := themeInfra.NewFileThemeHistoryStore(historyFile)
	ctx := context.Background()
	
	t.Run("adds theme to history", func(t *testing.T) {
		err := store.Add(ctx, theme.ThemeMocha)
		require.NoError(t, err)
		
		history, err := store.GetHistory(ctx)
		require.NoError(t, err)
		assert.Len(t, history, 1)
		assert.Equal(t, theme.ThemeMocha, history[0])
	})
	
	t.Run("maintains order (newest first)", func(t *testing.T) {
		store := themeInfra.NewFileThemeHistoryStore(filepath.Join(tmpDir, "order-test.json"))
		
		themes := []theme.ThemeName{theme.ThemeMocha, theme.ThemeLatte, theme.ThemeFrappe}
		for _, th := range themes {
			err := store.Add(ctx, th)
			require.NoError(t, err)
		}
		
		history, err := store.GetHistory(ctx)
		require.NoError(t, err)
		
		// Should be in reverse order (newest first)
		assert.Equal(t, theme.ThemeFrappe, history[0])
		assert.Equal(t, theme.ThemeLatte, history[1])
		assert.Equal(t, theme.ThemeMocha, history[2])
	})
	
	t.Run("limits history to 10 entries", func(t *testing.T) {
		store := themeInfra.NewFileThemeHistoryStore(filepath.Join(tmpDir, "limit-test.json"))
		
		// Add 15 themes
		for i := 0; i < 15; i++ {
			err := store.Add(ctx, theme.ThemeMocha)
			require.NoError(t, err)
		}
		
		history, err := store.GetHistory(ctx)
		require.NoError(t, err)
		assert.Len(t, history, 10, "History should be limited to 10 entries")
	})
}

func TestFileThemeHistoryStore_GetPrevious(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "theme-history.json")
	
	store := themeInfra.NewFileThemeHistoryStore(historyFile)
	ctx := context.Background()
	
	t.Run("returns previous theme", func(t *testing.T) {
		err := store.Add(ctx, theme.ThemeMocha)
		require.NoError(t, err)
		
		err = store.Add(ctx, theme.ThemeLatte)
		require.NoError(t, err)
		
		previous, err := store.GetPrevious(ctx)
		require.NoError(t, err)
		assert.Equal(t, theme.ThemeMocha, previous)
	})
	
	t.Run("returns error when no history", func(t *testing.T) {
		store := themeInfra.NewFileThemeHistoryStore(filepath.Join(tmpDir, "empty.json"))
		
		_, err := store.GetPrevious(ctx)
		assert.Error(t, err)
		assert.ErrorIs(t, err, themeInfra.ErrNoThemeHistory)
	})
	
	t.Run("returns error with only one theme", func(t *testing.T) {
		store := themeInfra.NewFileThemeHistoryStore(filepath.Join(tmpDir, "single.json"))
		
		err := store.Add(ctx, theme.ThemeMocha)
		require.NoError(t, err)
		
		_, err = store.GetPrevious(ctx)
		assert.Error(t, err)
		assert.ErrorIs(t, err, themeInfra.ErrNoThemeHistory)
	})
}

func TestFileThemeHistoryStore_GetHistory(t *testing.T) {
	tmpDir := t.TempDir()
	
	t.Run("returns empty slice when no history", func(t *testing.T) {
		store := themeInfra.NewFileThemeHistoryStore(filepath.Join(tmpDir, "empty.json"))
		ctx := context.Background()
		
		history, err := store.GetHistory(ctx)
		require.NoError(t, err)
		assert.Empty(t, history)
	})
	
	t.Run("returns all themes in order", func(t *testing.T) {
		store := themeInfra.NewFileThemeHistoryStore(filepath.Join(tmpDir, "multi.json"))
		ctx := context.Background()
		
		themes := []theme.ThemeName{
			theme.ThemeMocha,
			theme.ThemeLatte,
			theme.ThemeFrappe,
		}
		
		for _, th := range themes {
			err := store.Add(ctx, th)
			require.NoError(t, err)
		}
		
		history, err := store.GetHistory(ctx)
		require.NoError(t, err)
		assert.Len(t, history, 3)
		
		// Newest first
		assert.Equal(t, theme.ThemeFrappe, history[0])
		assert.Equal(t, theme.ThemeLatte, history[1])
		assert.Equal(t, theme.ThemeMocha, history[2])
	})
}

func TestFileThemeHistoryStore_Clear(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "theme-history.json")
	
	store := themeInfra.NewFileThemeHistoryStore(historyFile)
	ctx := context.Background()
	
	// Add some history
	err := store.Add(ctx, theme.ThemeMocha)
	require.NoError(t, err)
	err = store.Add(ctx, theme.ThemeLatte)
	require.NoError(t, err)
	
	// Verify history exists
	history, err := store.GetHistory(ctx)
	require.NoError(t, err)
	assert.Len(t, history, 2)
	
	// Clear
	err = store.Clear(ctx)
	require.NoError(t, err)
	
	// Verify empty
	history, err = store.GetHistory(ctx)
	require.NoError(t, err)
	assert.Empty(t, history)
}

func TestFileThemeHistoryStore_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "theme-history.json")
	ctx := context.Background()
	
	t.Run("persists history across instances", func(t *testing.T) {
		// First instance - add history
		store1 := themeInfra.NewFileThemeHistoryStore(historyFile)
		err := store1.Add(ctx, theme.ThemeMocha)
		require.NoError(t, err)
		err = store1.Add(ctx, theme.ThemeLatte)
		require.NoError(t, err)
		
		// Second instance - load history
		store2 := themeInfra.NewFileThemeHistoryStore(historyFile)
		history, err := store2.GetHistory(ctx)
		require.NoError(t, err)
		
		assert.Len(t, history, 2)
		assert.Equal(t, theme.ThemeLatte, history[0])
		assert.Equal(t, theme.ThemeMocha, history[1])
	})
}

func TestFileThemeHistoryStore_Remove(t *testing.T) {
	tmpDir := t.TempDir()
	historyFile := filepath.Join(tmpDir, "theme-history.json")
	
	store := themeInfra.NewFileThemeHistoryStore(historyFile)
	ctx := context.Background()
	
	t.Run("removes most recent theme", func(t *testing.T) {
		// Add themes
		err := store.Add(ctx, theme.ThemeMocha)
		require.NoError(t, err)
		err = store.Add(ctx, theme.ThemeLatte)
		require.NoError(t, err)
		err = store.Add(ctx, theme.ThemeFrappe)
		require.NoError(t, err)
		
		// Remove last one
		err = store.RemoveLast(ctx)
		require.NoError(t, err)
		
		// Verify
		history, err := store.GetHistory(ctx)
		require.NoError(t, err)
		assert.Len(t, history, 2)
		assert.Equal(t, theme.ThemeLatte, history[0])
	})
}
