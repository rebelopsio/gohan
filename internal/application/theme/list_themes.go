package theme

import (
	"context"

	"github.com/rebelopsio/gohan/internal/domain/theme"
)

// ListThemesUseCase lists all available themes
type ListThemesUseCase struct {
	registry theme.ThemeRegistry
}

// NewListThemesUseCase creates a new list themes use case
func NewListThemesUseCase(registry theme.ThemeRegistry) *ListThemesUseCase {
	return &ListThemesUseCase{
		registry: registry,
	}
}

// Execute lists all available themes
func (uc *ListThemesUseCase) Execute(ctx context.Context) ([]ThemeInfo, error) {
	// Get all themes from registry
	themes := uc.registry.ListAll()

	// Get active theme to mark it
	activeTheme, err := uc.registry.GetActive(ctx)
	var activeThemeName theme.ThemeName
	if err == nil && activeTheme != nil {
		activeThemeName = activeTheme.Name()
	}

	// Convert to DTOs
	result := make([]ThemeInfo, len(themes))
	for i, th := range themes {
		result[i] = themeToDTO(th, th.Name() == activeThemeName)
	}

	return result, nil
}

// themeToDTO converts a domain theme to a DTO
func themeToDTO(th *theme.Theme, isActive bool) ThemeInfo {
	return ThemeInfo{
		Name:        string(th.Name()),
		DisplayName: th.DisplayName(),
		Author:      th.Author(),
		Description: th.Description(),
		Variant:     string(th.Variant()),
		IsActive:    isActive,
		PreviewURL:  th.PreviewURL(),
		ColorScheme: colorSchemeToMap(th.ColorScheme()),
	}
}

// colorSchemeToMap converts a color scheme to a map
func colorSchemeToMap(cs theme.ColorScheme) map[string]string {
	return map[string]string{
		// Base colors
		"base":    cs.Base().String(),
		"surface": cs.Surface().String(),
		"overlay": cs.Overlay().String(),
		"text":    cs.Text().String(),
		"subtext": cs.Subtext().String(),

		// Accent colors
		"rosewater": cs.Rosewater().String(),
		"flamingo":  cs.Flamingo().String(),
		"pink":      cs.Pink().String(),
		"mauve":     cs.Mauve().String(),
		"red":       cs.Red().String(),
		"maroon":    cs.Maroon().String(),
		"peach":     cs.Peach().String(),
		"yellow":    cs.Yellow().String(),
		"green":     cs.Green().String(),
		"teal":      cs.Teal().String(),
		"sky":       cs.Sky().String(),
		"sapphire":  cs.Sapphire().String(),
		"blue":      cs.Blue().String(),
		"lavender":  cs.Lavender().String(),
	}
}
