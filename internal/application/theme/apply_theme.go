package theme

import (
	"context"
	"fmt"

	"github.com/rebelopsio/gohan/internal/domain/theme"
)

// ThemeApplier is the interface for applying themes to the system
type ThemeApplier interface {
	ApplyTheme(ctx context.Context, th *theme.Theme) error
}

// ApplyThemeResult contains the result of applying a theme
type ApplyThemeResult struct {
	Success            bool
	ThemeName          string
	Message            string
	AffectedComponents []string
	BackupID           string
}

// ApplyThemeUseCase applies a theme to the system
type ApplyThemeUseCase struct {
	registry theme.ThemeRegistry
	applier  ThemeApplier
}

// NewApplyThemeUseCase creates a new apply theme use case
func NewApplyThemeUseCase(registry theme.ThemeRegistry, applier ThemeApplier) *ApplyThemeUseCase {
	return &ApplyThemeUseCase{
		registry: registry,
		applier:  applier,
	}
}

// Execute applies the specified theme
func (uc *ApplyThemeUseCase) Execute(ctx context.Context, themeName string) (*ApplyThemeResult, error) {
	// Find the theme
	th, err := uc.registry.FindByName(ctx, theme.ThemeName(themeName))
	if err != nil {
		return nil, fmt.Errorf("theme not found: %w", err)
	}

	// Check if already active
	activeTheme, err := uc.registry.GetActive(ctx)
	if err == nil && activeTheme.Name() == th.Name() {
		return &ApplyThemeResult{
			Success:   true,
			ThemeName: themeName,
			Message:   fmt.Sprintf("Theme '%s' is already active", themeName),
		}, nil
	}

	// Apply theme if applier is available
	if uc.applier != nil {
		if err := uc.applier.ApplyTheme(ctx, th); err != nil {
			return nil, fmt.Errorf("failed to apply theme: %w", err)
		}
	}

	// Set as active
	if err := uc.registry.SetActive(ctx, th.Name()); err != nil {
		return nil, fmt.Errorf("failed to set active theme: %w", err)
	}

	return &ApplyThemeResult{
		Success:   true,
		ThemeName: themeName,
		Message:   fmt.Sprintf("Successfully applied theme '%s'", th.DisplayName()),
	}, nil
}
