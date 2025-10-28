package theme

import (
	"context"

	"github.com/rebelopsio/gohan/internal/domain/theme"
)

// GetActiveThemeUseCase gets the currently active theme
type GetActiveThemeUseCase struct {
	registry theme.ThemeRegistry
}

// NewGetActiveThemeUseCase creates a new get active theme use case
func NewGetActiveThemeUseCase(registry theme.ThemeRegistry) *GetActiveThemeUseCase {
	return &GetActiveThemeUseCase{
		registry: registry,
	}
}

// Execute gets the currently active theme
func (uc *GetActiveThemeUseCase) Execute(ctx context.Context) (*ThemeInfo, error) {
	activeTheme, err := uc.registry.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	themeInfo := themeToDTO(activeTheme, true)
	return &themeInfo, nil
}
