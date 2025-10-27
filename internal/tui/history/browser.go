package history

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rebelopsio/gohan/internal/application/history/services"
	"github.com/rebelopsio/gohan/internal/domain/history"
)

// browserState represents the current state of the browser
type browserState int

const (
	stateLoading browserState = iota
	stateList
	stateDetail
	stateError
)

// Browser is the main Bubble Tea model for browsing installation history
type Browser struct {
	state         browserState
	queryService  *services.HistoryQueryService
	records       []history.InstallationRecord
	selectedIndex int
	width         int
	height        int
	err           error
	ctx           context.Context
	cancel        context.CancelFunc
}

// Message types
type recordsLoadedMsg struct {
	records []history.InstallationRecord
}

type errorMsg struct {
	err error
}

// NewBrowser creates a new history browser
func NewBrowser(queryService *services.HistoryQueryService) *Browser {
	ctx, cancel := context.WithCancel(context.Background())

	return &Browser{
		state:         stateLoading,
		queryService:  queryService,
		records:       []history.InstallationRecord{},
		selectedIndex: 0,
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Init initializes the browser and loads records
func (b *Browser) Init() tea.Cmd {
	return b.loadRecords()
}

// Update handles messages
func (b *Browser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return b.handleKey(msg)

	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		return b, nil

	case recordsLoadedMsg:
		b.records = msg.records
		b.state = stateList
		return b, nil

	case errorMsg:
		b.err = msg.err
		b.state = stateError
		return b, nil
	}

	return b, nil
}

// handleKey handles keyboard input
func (b *Browser) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		b.cancel()
		return b, tea.Quit

	case "up", "k":
		if b.state == stateList && b.selectedIndex > 0 {
			b.selectedIndex--
		}

	case "down", "j":
		if b.state == stateList && b.selectedIndex < len(b.records)-1 {
			b.selectedIndex++
		}

	case "enter":
		if b.state == stateList && len(b.records) > 0 {
			b.state = stateDetail
		}

	case "esc":
		if b.state == stateDetail {
			b.state = stateList
		}

	case "r":
		if b.state == stateList {
			b.state = stateLoading
			return b, b.loadRecords()
		}
	}

	return b, nil
}

// View renders the current view
func (b *Browser) View() string {
	switch b.state {
	case stateLoading:
		return b.renderLoading()
	case stateList:
		return b.renderList()
	case stateDetail:
		return b.renderDetail()
	case stateError:
		return b.renderError()
	default:
		return "Unknown state"
	}
}

// loadRecords loads installation records
func (b *Browser) loadRecords() tea.Cmd {
	return func() tea.Msg {
		records, err := b.queryService.ListRecent(b.ctx, 50)
		if err != nil {
			return errorMsg{err: err}
		}
		return recordsLoadedMsg{records: records}
	}
}

// renderLoading renders the loading state
func (b *Browser) renderLoading() string {
	return titleStyle.Render("Installation History") + "\n\n" +
		"Loading records...\n"
}

// renderList renders the list view
func (b *Browser) renderList() string {
	var s strings.Builder

	// Title
	s.WriteString(titleStyle.Render("Installation History"))
	s.WriteString("\n\n")

	if len(b.records) == 0 {
		s.WriteString("📦 No installation records found.\n\n")
		s.WriteString("Your installation history is empty. Records will appear here after you:\n")
		s.WriteString("  • Complete an installation with 'gohan install'\n")
		s.WriteString("  • Run any package installation through gohan\n\n")
		s.WriteString(helpStyle.Render("Press 'r' to refresh • Press 'q' to quit"))
		return s.String()
	}

	// Header
	header := fmt.Sprintf("%-10s %-20s %-10s %-8s %-20s",
		"STATUS", "PACKAGE", "VERSION", "DURATION", "INSTALLED")
	s.WriteString(listHeaderStyle.Render(header))
	s.WriteString("\n\n")

	// Calculate visible range
	visibleHeight := b.height - 10 // Account for header, footer, padding
	if visibleHeight < 5 {
		visibleHeight = 5
	}

	startIdx := 0
	endIdx := len(b.records)

	// If we have more records than can fit, show a window around selected
	if len(b.records) > visibleHeight {
		startIdx = b.selectedIndex - visibleHeight/2
		if startIdx < 0 {
			startIdx = 0
		}
		endIdx = startIdx + visibleHeight
		if endIdx > len(b.records) {
			endIdx = len(b.records)
			startIdx = endIdx - visibleHeight
			if startIdx < 0 {
				startIdx = 0
			}
		}
	}

	// Records
	for i := startIdx; i < endIdx; i++ {
		record := b.records[i]
		line := b.formatRecordLine(record)

		if i == b.selectedIndex {
			s.WriteString(selectedItemStyle.Render("→ " + line))
		} else {
			s.WriteString(listItemStyle.Render("  " + line))
		}
		s.WriteString("\n")
	}

	// Footer with navigation help
	s.WriteString("\n")
	s.WriteString(footerStyle.Render(
		fmt.Sprintf("Showing %d-%d of %d records", startIdx+1, endIdx, len(b.records))))
	s.WriteString("\n")
	s.WriteString(helpStyle.Render(
		"↑/k: up • ↓/j: down • enter: details • r: refresh • q: quit"))

	return s.String()
}

// renderDetail renders the detail view
func (b *Browser) renderDetail() string {
	if b.selectedIndex >= len(b.records) {
		return "Invalid selection"
	}

	record := b.records[b.selectedIndex]

	var s strings.Builder

	// Title
	s.WriteString(titleStyle.Render("Installation Details"))
	s.WriteString("\n\n")

	// Basic info section
	basicInfo := b.renderBasicInfo(record)
	s.WriteString(detailSectionStyle.Render(basicInfo))

	// System context section
	systemInfo := b.renderSystemContext(record)
	s.WriteString(detailSectionStyle.Render(systemInfo))

	// Installed packages section
	if record.PackageCount() > 0 {
		packagesInfo := b.renderInstalledPackages(record)
		s.WriteString(detailSectionStyle.Render(packagesInfo))
	}

	// Failure details if present
	if record.HasFailureDetails() {
		failureInfo := b.renderFailureDetails(record)
		s.WriteString(detailSectionStyle.Render(failureInfo))
	}

	// Footer
	s.WriteString("\n")
	s.WriteString(helpStyle.Render("esc: back to list • q: quit"))

	return s.String()
}

// renderBasicInfo renders basic record information
func (b *Browser) renderBasicInfo(record history.InstallationRecord) string {
	var s strings.Builder

	s.WriteString(lipgloss.NewStyle().Bold(true).Render("Basic Information"))
	s.WriteString("\n\n")

	s.WriteString(detailLabelStyle.Render("Record ID:"))
	s.WriteString(detailValueStyle.Render(truncateID(record.ID().String(), 12)))
	s.WriteString("\n")

	s.WriteString(detailLabelStyle.Render("Package:"))
	s.WriteString(detailValueStyle.Render(record.PackageName()))
	s.WriteString("\n")

	s.WriteString(detailLabelStyle.Render("Version:"))
	s.WriteString(detailValueStyle.Render(record.TargetVersion()))
	s.WriteString("\n")

	s.WriteString(detailLabelStyle.Render("Status:"))
	s.WriteString(b.formatStatus(record.Outcome()))
	s.WriteString("\n")

	s.WriteString(detailLabelStyle.Render("Installed At:"))
	s.WriteString(detailValueStyle.Render(record.InstalledAt().Format(time.RFC3339)))
	s.WriteString("\n")

	s.WriteString(detailLabelStyle.Render("Duration:"))
	s.WriteString(detailValueStyle.Render(formatDuration(record.Duration())))

	return s.String()
}

// renderSystemContext renders system context information
func (b *Browser) renderSystemContext(record history.InstallationRecord) string {
	var s strings.Builder
	sysCtx := record.SystemContext()

	s.WriteString(lipgloss.NewStyle().Bold(true).Render("System Context"))
	s.WriteString("\n\n")

	s.WriteString(detailLabelStyle.Render("OS:"))
	s.WriteString(detailValueStyle.Render(sysCtx.OSVersion()))
	s.WriteString("\n")

	if sysCtx.KernelVersion() != "" {
		s.WriteString(detailLabelStyle.Render("Kernel:"))
		s.WriteString(detailValueStyle.Render(sysCtx.KernelVersion()))
		s.WriteString("\n")
	}

	s.WriteString(detailLabelStyle.Render("Gohan:"))
	s.WriteString(detailValueStyle.Render(sysCtx.GohanVersion()))
	s.WriteString("\n")

	if sysCtx.Hostname() != "" {
		s.WriteString(detailLabelStyle.Render("Hostname:"))
		s.WriteString(detailValueStyle.Render(sysCtx.Hostname()))
	}

	return s.String()
}

// renderInstalledPackages renders installed packages information
func (b *Browser) renderInstalledPackages(record history.InstallationRecord) string {
	var s strings.Builder
	metadata := record.Metadata()
	packages := metadata.InstalledPackages()

	s.WriteString(lipgloss.NewStyle().Bold(true).Render(
		fmt.Sprintf("Installed Packages (%d)", len(packages))))
	s.WriteString("\n\n")

	for _, pkg := range packages {
		s.WriteString(fmt.Sprintf("• %s %s (%s)\n",
			pkg.Name(),
			pkg.Version(),
			formatSize(pkg.SizeBytes())))
	}

	return strings.TrimRight(s.String(), "\n")
}

// renderFailureDetails renders failure details
func (b *Browser) renderFailureDetails(record history.InstallationRecord) string {
	var s strings.Builder
	fd := record.FailureDetails()

	s.WriteString(failedStatusStyle.Render("Failure Details"))
	s.WriteString("\n\n")

	s.WriteString(detailLabelStyle.Render("Reason:"))
	s.WriteString(detailValueStyle.Render(fd.Reason()))
	s.WriteString("\n")

	s.WriteString(detailLabelStyle.Render("Failed At:"))
	s.WriteString(detailValueStyle.Render(fd.FailedAt().Format(time.RFC3339)))
	s.WriteString("\n")

	s.WriteString(detailLabelStyle.Render("Phase:"))
	s.WriteString(detailValueStyle.Render(fd.Phase()))

	if fd.ErrorCode() != "" {
		s.WriteString("\n")
		s.WriteString(detailLabelStyle.Render("Error Code:"))
		s.WriteString(detailValueStyle.Render(fd.ErrorCode()))
	}

	return s.String()
}

// renderError renders the error state
func (b *Browser) renderError() string {
	return titleStyle.Render("Installation History") + "\n\n" +
		errorStyle.Render(fmt.Sprintf("Error: %v", b.err)) + "\n\n" +
		helpStyle.Render("Press 'q' to quit")
}

// Helper functions

func (b *Browser) formatRecordLine(record history.InstallationRecord) string {
	status := b.formatStatusShort(record.Outcome())
	pkg := truncateString(record.PackageName(), 20)
	version := truncateString(record.TargetVersion(), 10)
	duration := formatDuration(record.Duration())
	installed := record.InstalledAt().Format("2006-01-02 15:04")

	return fmt.Sprintf("%-10s %-20s %-10s %-8s %-20s",
		status, pkg, version, duration, installed)
}

func (b *Browser) formatStatus(outcome history.InstallationOutcome) string {
	switch {
	case outcome.IsSuccessful():
		return successStatusStyle.Render("✓ Success")
	case outcome.IsFailed():
		return failedStatusStyle.Render("✗ Failed")
	case outcome.IsRolledBack():
		return failedStatusStyle.Render("↻ Rolled Back")
	default:
		return outcome.String()
	}
}

func (b *Browser) formatStatusShort(outcome history.InstallationOutcome) string {
	switch {
	case outcome.IsSuccessful():
		return "✓ Success"
	case outcome.IsFailed():
		return "✗ Failed"
	case outcome.IsRolledBack():
		return "↻ Rollback"
	default:
		return outcome.String()
	}
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func truncateID(id string, length int) string {
	if len(id) <= length {
		return id
	}
	return id[:length]
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return "< 1s"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}

func formatSize(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
