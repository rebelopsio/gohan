package theme

import (
	"context"
	"fmt"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/theme"
	themeInfra "github.com/rebelopsio/gohan/internal/infrastructure/theme"
)

// RollbackThemeResult contains the result of rolling back a theme
type RollbackThemeResult struct {
	Success       bool
	RestoredTheme theme.ThemeName
	PreviousTheme theme.ThemeName
	Message       string
}

// RollbackThemeUseCase rolls back to the previous theme
type RollbackThemeUseCase struct {
	registry     theme.ThemeRegistry
	historyStore themeInfra.ThemeHistoryStore
	applier      ThemeApplier
	stateStore   themeInfra.ThemeStateStore
}

// NewRollbackThemeUseCase creates a new rollback theme use case
func NewRollbackThemeUseCase(
	registry theme.ThemeRegistry,
	historyStore themeInfra.ThemeHistoryStore,
	applier ThemeApplier,
	stateStore themeInfra.ThemeStateStore,
) *RollbackThemeUseCase {
	return &RollbackThemeUseCase{
		registry:     registry,
		historyStore: historyStore,
		applier:      applier,
		stateStore:   stateStore,
	}
}

// Execute rolls back to the previous theme in history
func (uc *RollbackThemeUseCase) Execute(ctx context.Context) (*RollbackThemeResult, error) {
	// Get current active theme for result reporting
	currentTheme, err := uc.registry.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active theme: %w", err)
	}
	currentThemeName := currentTheme.Name()

	// Get previous theme from history
	previousThemeName, err := uc.historyStore.GetPrevious(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot rollback: %w", err)
	}

	// Find the previous theme in registry
	previousTheme, err := uc.registry.FindByName(ctx, previousThemeName)
	if err != nil {
		return nil, fmt.Errorf("previous theme not found: %w", err)
	}

	// Apply theme if applier is available
	if uc.applier != nil {
		if err := uc.applier.ApplyTheme(ctx, previousTheme); err != nil {
			return nil, fmt.Errorf("failed to apply theme: %w", err)
		}
	}

	// Set as active in registry
	if err := uc.registry.SetActive(ctx, previousTheme.Name()); err != nil {
		return nil, fmt.Errorf("failed to set active theme: %w", err)
	}

	// Save state to disk if state store is available
	if uc.stateStore != nil {
		state := &themeInfra.ThemeState{
			ThemeName:    previousTheme.Name(),
			ThemeVariant: previousTheme.Variant(),
			SetAt:        time.Now(),
		}

		if err := uc.stateStore.Save(ctx, state); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: failed to save theme state: %v\n", err)
		}
	}

	// Remove the most recent entry from history
	// This allows sequential rollbacks to work correctly
	if err := uc.historyStore.RemoveLast(ctx); err != nil {
		// Log error but don't fail the operation
		// The theme has already been applied successfully
		fmt.Printf("Warning: failed to update history: %v\n", err)
	}

	return &RollbackThemeResult{
		Success:       true,
		RestoredTheme: previousTheme.Name(),
		PreviousTheme: currentThemeName,
		Message:       fmt.Sprintf("Successfully rolled back from '%s' to '%s'", currentThemeName, previousTheme.Name()),
	}, nil
}
