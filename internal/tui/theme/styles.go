package theme

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Title styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#cba6f7")). // Catppuccin Mauve
			MarginTop(1).
			MarginBottom(1)

	// Subtitle style
	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#bac2de")). // Catppuccin Subtext
			MarginBottom(1)

	// Selected theme style
	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#cba6f7")). // Catppuccin Mauve
			Background(lipgloss.Color("#313244")). // Catppuccin Surface0
			Padding(0, 2).
			Margin(0, 1)

	// Unselected theme style
	unselectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#cdd6f4")). // Catppuccin Text
				Padding(0, 2).
				Margin(0, 1)

	// Active theme indicator
	activeIndicatorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#a6e3a1")). // Catppuccin Green
				Bold(true)

	// Color swatch style
	swatchStyle = lipgloss.NewStyle().
			Width(4).
			Align(lipgloss.Center)

	// Description style
	descriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6c7086")). // Catppuccin Overlay0
				Italic(true).
				MarginLeft(8)

	// Color preview box
	previewBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#45475a")). // Catppuccin Surface2
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// Color name style
	colorNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#89b4fa")). // Catppuccin Blue
			Width(12)

	// Color value style
	colorValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")). // Catppuccin Overlay0
			Italic(true)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")). // Catppuccin Overlay0
			MarginTop(1).
			Italic(true)

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8")). // Catppuccin Red
			Bold(true).
			Margin(1)

	// Success style
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6e3a1")). // Catppuccin Green
			Bold(true).
			Margin(1)
)

// renderColorSwatch creates a colored block for color preview
func renderColorSwatch(color string) string {
	return swatchStyle.
		Background(lipgloss.Color(color)).
		Render("    ")
}
