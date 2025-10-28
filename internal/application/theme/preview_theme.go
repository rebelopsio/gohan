package theme

import (
	"context"
	"fmt"

	"github.com/rebelopsio/gohan/internal/domain/theme"
)

// PreviewThemeUseCase generates a preview of a theme without applying it
type PreviewThemeUseCase struct {
	registry theme.ThemeRegistry
}

// NewPreviewThemeUseCase creates a new preview theme use case
func NewPreviewThemeUseCase(registry theme.ThemeRegistry) *PreviewThemeUseCase {
	return &PreviewThemeUseCase{
		registry: registry,
	}
}

// Execute generates a preview of the specified theme
func (uc *PreviewThemeUseCase) Execute(ctx context.Context, themeName string) (*ThemePreview, error) {
	// Find the theme
	th, err := uc.registry.FindByName(ctx, theme.ThemeName(themeName))
	if err != nil {
		return nil, err
	}

	// Generate preview
	preview := &ThemePreview{
		Name:        string(th.Name()),
		DisplayName: th.DisplayName(),
		Author:      th.Author(),
		Description: th.Description(),
		Variant:     string(th.Variant()),
		ColorScheme: colorSchemeToMap(th.ColorScheme()),
		PreviewText: generatePreviewText(th),
	}

	return preview, nil
}

// generatePreviewText creates a visual preview text for the theme
func generatePreviewText(th *theme.Theme) string {
	cs := th.ColorScheme()

	var preview string
	preview += fmt.Sprintf("%s (%s theme)\n\n", th.DisplayName(), th.Variant())
	preview += "Color Palette:\n"
	preview += fmt.Sprintf("  Background: %s\n", cs.Base())
	preview += fmt.Sprintf("  Text:       %s\n", cs.Text())
	preview += fmt.Sprintf("  Accent:     %s\n", cs.Mauve())
	preview += fmt.Sprintf("  Success:    %s\n", cs.Green())
	preview += fmt.Sprintf("  Error:      %s\n", cs.Red())

	return preview
}
