package installation

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#7C3AED") // Purple
	successColor   = lipgloss.Color("#10B981") // Green
	errorColor     = lipgloss.Color("#EF4444") // Red
	warningColor   = lipgloss.Color("#F59E0B") // Orange
	infoColor      = lipgloss.Color("#3B82F6") // Blue
	textColor      = lipgloss.Color("#E5E7EB") // Light gray
	dimColor       = lipgloss.Color("#9CA3AF") // Dim gray
	borderColor    = lipgloss.Color("#6B7280") // Border gray

	// Title styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	// Box styles
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			MarginTop(1)

	headerBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 2)

	// Progress bar styles
	progressBarWidth    = 40
	progressBarComplete = lipgloss.NewStyle().Foreground(successColor)
	progressBarEmpty    = lipgloss.NewStyle().Foreground(dimColor)

	// Phase styles
	phaseStyle = lipgloss.NewStyle().
			Foreground(infoColor).
			Bold(true)

	// Package info styles
	packageNameStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	packageVersionStyle = lipgloss.NewStyle().
				Foreground(dimColor)

	// Status styles
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	// Log entry styles
	logSuccessStyle = lipgloss.NewStyle().Foreground(successColor)
	logErrorStyle   = lipgloss.NewStyle().Foreground(errorColor)
	logInfoStyle    = lipgloss.NewStyle().Foreground(textColor)
	logDimStyle     = lipgloss.NewStyle().Foreground(dimColor)

	// Time styles
	timeStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Italic(true)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Italic(true).
			Align(lipgloss.Center)

	// Spinner style
	spinnerStyle = lipgloss.NewStyle().Foreground(primaryColor)
)

// renderProgressBar renders a progress bar with percentage
func renderProgressBar(percent int, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := (percent * width) / 100
	empty := width - filled

	bar := ""
	if filled > 0 {
		bar += progressBarComplete.Render(repeatString("█", filled))
	}
	if empty > 0 {
		bar += progressBarEmpty.Render(repeatString("░", empty))
	}

	return bar
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
