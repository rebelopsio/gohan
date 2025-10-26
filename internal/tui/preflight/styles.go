package preflight

import "github.com/charmbracelet/lipgloss"

var (
	// Theme colors
	colorPrimary   = lipgloss.Color("#89b4fa") // Blue
	colorSuccess   = lipgloss.Color("#a6e3a1") // Green
	colorWarning   = lipgloss.Color("#f9e2af") // Yellow
	colorError     = lipgloss.Color("#f38ba8") // Red
	colorSubtle    = lipgloss.Color("#6c7086") // Gray
	colorHighlight = lipgloss.Color("#cba6f7") // Purple

	// Title styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Margin(1, 0)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorSubtle).
			Margin(0, 0, 1, 0)

	// Status styles
	successStyle = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(colorWarning).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true)

	// Info styles
	labelStyle = lipgloss.NewStyle().
			Foreground(colorSubtle)

	valueStyle = lipgloss.NewStyle().
			Bold(true)

	// Box styles
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(1, 2).
			Margin(1, 0)

	guidanceBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorWarning).
				Padding(1, 2).
				Margin(1, 0)

	// Progress styles
	spinnerStyle = lipgloss.NewStyle().
			Foreground(colorPrimary)

	progressItemStyle = lipgloss.NewStyle().
				Margin(0, 0, 0, 2)

	progressItemDoneStyle = lipgloss.NewStyle().
				Foreground(colorSuccess).
				Margin(0, 0, 0, 2)

	progressItemFailedStyle = lipgloss.NewStyle().
				Foreground(colorError).
				Margin(0, 0, 0, 2)

	progressItemCurrentStyle = lipgloss.NewStyle().
					Foreground(colorHighlight).
					Margin(0, 0, 0, 2)

	// Help styles
	helpStyle = lipgloss.NewStyle().
			Foreground(colorSubtle).
			Margin(1, 0, 0, 0)
)

// StatusIcon returns the appropriate icon for a status
func StatusIcon(status string) string {
	switch status {
	case "pass", "success":
		return successStyle.Render("✓")
	case "fail", "error":
		return errorStyle.Render("✗")
	case "warning":
		return warningStyle.Render("⚠")
	case "running":
		return spinnerStyle.Render("⋯")
	default:
		return labelStyle.Render("·")
	}
}
