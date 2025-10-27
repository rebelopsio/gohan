//go:build acceptance
// +build acceptance

package acceptance

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestThemeRollback covers scenarios from theme-rollback.feature
func TestThemeRollback(t *testing.T) {
	ctx := context.Background()

	t.Run("Restore previous theme", func(t *testing.T) {
		// Given the theme system is initialized
		// And the "mocha" theme is active
		themeService := setupThemeRollbackService(t)
		err := themeService.SetActiveTheme(ctx, "mocha")
		require.NoError(t, err)

		// And I recently changed from "mocha" to "latte"
		_, err = themeService.SwitchTheme(ctx, "latte")
		require.NoError(t, err)

		// When I undo the theme change
		result, err := themeService.UndoThemeChange(ctx)
		require.NoError(t, err)

		// Then my desktop should display the "mocha" theme again
		activeTheme, err := themeService.GetActiveTheme(ctx)
		require.NoError(t, err)
		assert.Equal(t, "mocha", activeTheme.Name)

		// And my previous settings should be restored
		assert.True(t, result.Success)
		assert.NotEmpty(t, result.RestoredFiles)

		// And I should see confirmation of the restoration
		assert.Contains(t, result.Message, "restored")
	})

	t.Run("Restore to specific earlier theme", func(t *testing.T) {
		// Given the theme system is initialized
		themeService := setupThemeRollbackService(t)

		// And I have changed themes multiple times:
		// | from      | to        | when                |
		// | mocha     | latte     | earlier today       |
		// | latte     | frappe    | a few hours ago     |
		// | frappe    | macchiato | recently            |
		err := themeService.SetActiveTheme(ctx, "mocha")
		require.NoError(t, err)

		_, err = themeService.SwitchTheme(ctx, "latte")
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		switchToFrappeResult, err := themeService.SwitchTheme(ctx, "frappe")
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		_, err = themeService.SwitchTheme(ctx, "macchiato")
		require.NoError(t, err)

		// When I restore my appearance from a few hours ago
		result, err := themeService.RestoreToPoint(ctx, switchToFrappeResult.BackupID)
		require.NoError(t, err)

		// Then my desktop should display the "latte" theme
		activeTheme, err := themeService.GetActiveTheme(ctx)
		require.NoError(t, err)
		assert.Equal(t, "latte", activeTheme.Name)

		// And my desktop should look as it did then
		assert.True(t, result.Success)
	})

	t.Run("Restoration shows progress", func(t *testing.T) {
		// Given the theme system is initialized
		// And I recently changed from "mocha" to "gohan"
		themeService := setupThemeRollbackService(t)
		err := themeService.SetActiveTheme(ctx, "mocha")
		require.NoError(t, err)
		_, err = themeService.SwitchTheme(ctx, "gohan")
		require.NoError(t, err)

		// When I undo the theme change
		progressChan := make(chan RestoreProgress, 10)
		done := make(chan error, 1)

		go func() {
			done <- themeService.UndoThemeChangeWithProgress(ctx, progressChan)
			close(progressChan)
		}()

		// Then I should see the restoration in progress
		var progressUpdates []RestoreProgress
		for progress := range progressChan {
			progressUpdates = append(progressUpdates, progress)
		}

		// And I should be notified when complete
		err = <-done
		require.NoError(t, err)
		assert.NotEmpty(t, progressUpdates, "should receive progress updates")

		lastProgress := progressUpdates[len(progressUpdates)-1]
		assert.Equal(t, "completed", lastProgress.Status)
	})

	t.Run("View my theme change history", func(t *testing.T) {
		// Given the theme system is initialized
		themeService := setupThemeRollbackService(t)

		// And I changed themes twice today:
		// | from  | to     | when            |
		// | mocha | latte  | 2 hours ago     |
		// | latte | frappe | 30 minutes ago  |
		err := themeService.SetActiveTheme(ctx, "mocha")
		require.NoError(t, err)

		_, err = themeService.SwitchTheme(ctx, "latte")
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		_, err = themeService.SwitchTheme(ctx, "frappe")
		require.NoError(t, err)

		// When I view my theme history
		history, err := themeService.GetThemeHistory(ctx)
		require.NoError(t, err)

		// Then I should see 2 previous themes I can restore
		assert.GreaterOrEqual(t, len(history), 2)

		// And they should be sorted newest first
		assert.True(t, history[0].Timestamp.After(history[1].Timestamp) ||
			history[0].Timestamp.Equal(history[1].Timestamp))

		// And each should show the theme transition
		for _, entry := range history {
			assert.NotEmpty(t, entry.FromTheme)
			assert.NotEmpty(t, entry.ToTheme)
			assert.NotEmpty(t, entry.BackupID)
		}
	})

	t.Run("Attempt restore with no history", func(t *testing.T) {
		// Given no theme changes have been made
		themeService := setupThemeRollbackService(t)

		// When I attempt to undo changes
		_, err := themeService.UndoThemeChange(ctx)

		// Then I should receive an error
		require.Error(t, err)

		// And the error should indicate no previous themes are available
		assert.Contains(t, err.Error(), "no previous")

		// And the current theme should remain active
		activeTheme, err := themeService.GetActiveTheme(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, activeTheme.Name)
	})

	t.Run("Restore to invalid point in history", func(t *testing.T) {
		// Given the theme system is initialized
		themeService := setupThemeRollbackService(t)

		// When I attempt to restore to a non-existent point in history
		_, err := themeService.RestoreToPoint(ctx, "invalid-backup-id")

		// Then I should receive an error
		require.Error(t, err)

		// And the error should indicate that point was not found
		assert.Contains(t, err.Error(), "not found")

		// And the current theme should remain unchanged
		activeTheme, err := themeService.GetActiveTheme(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, activeTheme.Name)
	})

	t.Run("Failed restoration preserves current state", func(t *testing.T) {
		// Given the theme system is initialized
		// And I switched from "mocha" to "latte"
		themeService := setupThemeRollbackService(t)
		err := themeService.SetActiveTheme(ctx, "mocha")
		require.NoError(t, err)
		switchResult, err := themeService.SwitchTheme(ctx, "latte")
		require.NoError(t, err)

		// And the saved settings are corrupted
		themeService.CorruptBackup(switchResult.BackupID)

		// When I attempt to undo the change
		result, err := themeService.UndoThemeChange(ctx)

		// Then I should receive an error
		require.Error(t, err)

		// And the error should describe the problem
		assert.NotEmpty(t, result.ErrorMessage)

		// And the "latte" theme should remain active
		activeTheme, err := themeService.GetActiveTheme(ctx)
		require.NoError(t, err)
		assert.Equal(t, "latte", activeTheme.Name)

		// And my current settings should be unchanged
		assert.False(t, result.Success)
	})

	t.Run("Undo multiple theme changes sequentially", func(t *testing.T) {
		// Given the theme system is initialized
		themeService := setupThemeRollbackService(t)

		// And I switched from "mocha" to "latte"
		err := themeService.SetActiveTheme(ctx, "mocha")
		require.NoError(t, err)
		_, err = themeService.SwitchTheme(ctx, "latte")
		require.NoError(t, err)

		// And I switched from "latte" to "frappe"
		_, err = themeService.SwitchTheme(ctx, "frappe")
		require.NoError(t, err)

		// And I switched from "frappe" to "macchiato"
		_, err = themeService.SwitchTheme(ctx, "macchiato")
		require.NoError(t, err)

		// When I undo the most recent change
		_, err = themeService.UndoThemeChange(ctx)
		require.NoError(t, err)

		// Then the active theme should be "frappe"
		activeTheme, err := themeService.GetActiveTheme(ctx)
		require.NoError(t, err)
		assert.Equal(t, "frappe", activeTheme.Name)

		// When I undo again
		_, err = themeService.UndoThemeChange(ctx)
		require.NoError(t, err)

		// Then the active theme should be "latte"
		activeTheme, err = themeService.GetActiveTheme(ctx)
		require.NoError(t, err)
		assert.Equal(t, "latte", activeTheme.Name)

		// When I undo again
		_, err = themeService.UndoThemeChange(ctx)
		require.NoError(t, err)

		// Then the active theme should be "mocha"
		activeTheme, err = themeService.GetActiveTheme(ctx)
		require.NoError(t, err)
		assert.Equal(t, "mocha", activeTheme.Name)
	})
}

// setupThemeRollbackService initializes a theme service for rollback tests
func setupThemeRollbackService(t *testing.T) ThemeRollbackService {
	t.Helper()
	// TODO: Implement once we have the domain models
	t.Skip("Theme rollback service not yet implemented")
	return nil
}

// ThemeRollbackService extends ThemeSwitchingService with rollback operations
type ThemeRollbackService interface {
	ThemeSwitchingService
	UndoThemeChange(ctx context.Context) (RestoreResult, error)
	UndoThemeChangeWithProgress(ctx context.Context, progressChan chan<- RestoreProgress) error
	RestoreToPoint(ctx context.Context, backupID string) (RestoreResult, error)
	GetThemeHistory(ctx context.Context) ([]ThemeHistoryEntry, error)
	CorruptBackup(backupID string) // For testing failure scenarios
}

// RestoreProgress represents progress during theme restoration
type RestoreProgress struct {
	Component       string
	Status          string // "started", "restoring", "completed", "failed"
	PercentComplete float64
	Error           error
}

// ThemeHistoryEntry represents a theme change in history
type ThemeHistoryEntry struct {
	Timestamp time.Time
	FromTheme string
	ToTheme   string
	BackupID  string
}
