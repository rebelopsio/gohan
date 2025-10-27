package installation

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProgressUpdate represents a progress update from the installation
type ProgressUpdate struct {
	Phase              string
	PercentComplete    int
	Message            string
	ComponentsTotal    int
	ComponentsInstalled int
	CurrentComponent   string
	IsComplete         bool
	IsError            bool
	ErrorMessage       string
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp time.Time
	Level     string // "success", "error", "info"
	Message   string
}

// ProgressViewer is the Bubble Tea model for installation progress
type ProgressViewer struct {
	packageName    string
	packageVersion string
	startTime      time.Time

	// Current state
	currentUpdate  ProgressUpdate
	logs           []LogEntry
	maxLogs        int

	// UI state
	spinner        spinner.Model
	width          int
	height         int

	// Channels
	progressChan   <-chan ProgressUpdate
	done           bool
}

// NewProgressViewer creates a new progress viewer
func NewProgressViewer(packageName, packageVersion string, progressChan <-chan ProgressUpdate) *ProgressViewer {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return &ProgressViewer{
		packageName:    packageName,
		packageVersion: packageVersion,
		startTime:      time.Now(),
		maxLogs:        10,
		spinner:        s,
		progressChan:   progressChan,
		logs:           []LogEntry{},
	}
}

type progressMsg ProgressUpdate
type tickMsg time.Time

func waitForProgress(progressChan <-chan ProgressUpdate) tea.Cmd {
	return func() tea.Msg {
		update := <-progressChan
		return progressMsg(update)
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init initializes the progress viewer
func (m *ProgressViewer) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		waitForProgress(m.progressChan),
		tick(),
	)
}

// Update handles messages
func (m *ProgressViewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Allow quit with q or ctrl+c, but only after completion
		if m.done {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
		return m, nil

	case progressMsg:
		// Update progress
		m.currentUpdate = ProgressUpdate(msg)

		// Add log entry
		if msg.Message != "" {
			level := "info"
			if msg.IsError {
				level = "error"
			} else if msg.IsComplete {
				level = "success"
			}

			m.addLog(LogEntry{
				Timestamp: time.Now(),
				Level:     level,
				Message:   msg.Message,
			})
		}

		// Check if complete
		if msg.IsComplete || msg.IsError {
			m.done = true
			return m, nil
		}

		// Wait for next update
		return m, waitForProgress(m.progressChan)

	case tickMsg:
		// Update spinner and time display
		if !m.done {
			return m, tick()
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the progress viewer
func (m *ProgressViewer) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader())
	b.WriteString("\n\n")

	// Progress section
	b.WriteString(m.renderProgress())
	b.WriteString("\n\n")

	// Logs section
	b.WriteString(m.renderLogs())
	b.WriteString("\n")

	// Footer
	if m.done {
		b.WriteString("\n")
		b.WriteString(m.renderFooter())
	}

	return b.String()
}

func (m *ProgressViewer) renderHeader() string {
	title := fmt.Sprintf("Installing %s %s",
		packageNameStyle.Render(m.packageName),
		packageVersionStyle.Render(m.packageVersion))

	return headerBoxStyle.Render(title)
}

func (m *ProgressViewer) renderProgress() string {
	var b strings.Builder

	// Current phase
	phase := m.currentUpdate.Phase
	if phase == "" {
		phase = "Initializing"
	}

	if !m.done {
		b.WriteString(fmt.Sprintf("%s %s\n\n",
			m.spinner.View(),
			phaseStyle.Render(phase)))
	} else {
		if m.currentUpdate.IsError {
			b.WriteString(fmt.Sprintf("✗ %s\n\n", errorStyle.Render("Failed")))
		} else {
			b.WriteString(fmt.Sprintf("✓ %s\n\n", successStyle.Render("Completed")))
		}
	}

	// Progress bar
	percent := m.currentUpdate.PercentComplete
	bar := renderProgressBar(percent, progressBarWidth)
	b.WriteString(fmt.Sprintf("%s %3d%%\n\n", bar, percent))

	// Components progress
	if m.currentUpdate.ComponentsTotal > 0 {
		b.WriteString(fmt.Sprintf("Components: %d/%d\n\n",
			m.currentUpdate.ComponentsInstalled,
			m.currentUpdate.ComponentsTotal))
	}

	// Time information
	elapsed := time.Since(m.startTime)
	b.WriteString(timeStyle.Render(fmt.Sprintf("Time elapsed: %s", formatDuration(elapsed))))

	// Estimate remaining time if not done and we have progress
	if !m.done && percent > 0 && percent < 100 {
		estimatedTotal := time.Duration(float64(elapsed) / float64(percent) * 100)
		remaining := estimatedTotal - elapsed
		if remaining > 0 {
			b.WriteString(timeStyle.Render(fmt.Sprintf(" | Est. remaining: %s", formatDuration(remaining))))
		}
	}

	return boxStyle.Render(b.String())
}

func (m *ProgressViewer) renderLogs() string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().
		Foreground(dimColor).
		Bold(true).
		Render("Activity Log"))
	b.WriteString("\n\n")

	if len(m.logs) == 0 {
		b.WriteString(logDimStyle.Render("  Waiting for activity..."))
	} else {
		// Show last N logs
		start := 0
		if len(m.logs) > m.maxLogs {
			start = len(m.logs) - m.maxLogs
		}

		for _, log := range m.logs[start:] {
			prefix := "  "
			var style lipgloss.Style

			switch log.Level {
			case "success":
				prefix = "  ✓ "
				style = logSuccessStyle
			case "error":
				prefix = "  ✗ "
				style = logErrorStyle
			case "info":
				prefix = "  → "
				style = logInfoStyle
			}

			b.WriteString(prefix)
			b.WriteString(style.Render(log.Message))
			b.WriteString("\n")
		}
	}

	return boxStyle.Render(b.String())
}

func (m *ProgressViewer) renderFooter() string {
	if m.currentUpdate.IsError {
		msg := "Installation failed"
		if m.currentUpdate.ErrorMessage != "" {
			msg = m.currentUpdate.ErrorMessage
		}
		return helpStyle.Render(fmt.Sprintf("\n%s\n\nPress 'q' to exit", errorStyle.Render(msg)))
	}

	return helpStyle.Render("\nInstallation completed successfully!\n\nPress 'q' to exit")
}

func (m *ProgressViewer) addLog(entry LogEntry) {
	m.logs = append(m.logs, entry)
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return "< 1s"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes())%60)
}
