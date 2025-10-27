package history

import "github.com/charmbracelet/lipgloss"

var (
	// Color palette
	primaryColor   = lipgloss.Color("#7C3AED") // Purple
	successColor   = lipgloss.Color("#10B981") // Green
	errorColor     = lipgloss.Color("#EF4444") // Red
	mutedColor     = lipgloss.Color("#6B7280") // Gray
	highlightColor = lipgloss.Color("#F59E0B") // Amber

	// Title styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// List styles
	listHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Padding(0, 1)

	listItemStyle = lipgloss.NewStyle().
			Padding(0, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Foreground(highlightColor).
				Bold(true)

	// Status styles
	successStatusStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	failedStatusStyle = lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true)

	// Detail view styles
	detailLabelStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Width(20)

	detailValueStyle = lipgloss.NewStyle().
				Bold(true)

	detailSectionStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(1, 2).
				MarginBottom(1)

	// Help text styles
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			MarginTop(1)

	// Error styles
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor)

	// Footer styles
	footerStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1).
			Padding(1, 0)
)
