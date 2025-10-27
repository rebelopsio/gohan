package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/rebelopsio/gohan/internal/application/history/services"
	"github.com/rebelopsio/gohan/internal/domain/history"
	historyRepo "github.com/rebelopsio/gohan/internal/infrastructure/history/repository"
	historyTUI "github.com/rebelopsio/gohan/internal/tui/history"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	// Flags for history list command
	listLimit  int
	listStatus string
	listFrom   string
	listTo     string

	// Flags for history export command
	exportOutput string
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Manage installation history",
	Long: `View, query, and manage the permanent record of all installation activities.

The history command provides access to a complete audit trail of installations,
including successful installs, failures, timestamps, and system context.`,
}

// historyListCmd represents the history list command
var historyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installation history",
	Long: `Display a list of installation records with filtering options.

Examples:
  # List all installations
  gohan history list

  # List only successful installations
  gohan history list --status success

  # List only failed installations
  gohan history list --status failed

  # List last 10 installations
  gohan history list --limit 10

  # List installations in a date range
  gohan history list --from 2025-10-01 --to 2025-10-31`,
	RunE: runHistoryList,
}

// historyShowCmd represents the history show command
var historyShowCmd = &cobra.Command{
	Use:   "show <record-id>",
	Short: "Show detailed installation record",
	Long: `Display complete details for a specific installation record.

The show command displays:
  - Package name and version
  - Installation timestamp and duration
  - Success/failure status
  - System context (OS, hostname)
  - All installed packages
  - Failure details (if applicable)

Example:
  gohan history show abc123-def456-ghi789`,
	Args: cobra.ExactArgs(1),
	RunE: runHistoryShow,
}

// historyBrowseCmd represents the history browse command
var historyBrowseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse installation history interactively",
	Long: `Launch an interactive terminal UI to browse installation history.

The browse command provides a rich interface with:
  - List view with keyboard navigation (↑/↓)
  - Detailed view for selected records (Enter)
  - Real-time filtering and search
  - Beautiful formatting and colors

Navigation:
  ↑/k: Move up
  ↓/j: Move down
  Enter: View details
  Esc: Back to list
  r: Refresh
  q: Quit

Example:
  gohan history browse`,
	RunE: runHistoryBrowse,
}

func init() {
	// Add subcommands to history
	historyCmd.AddCommand(historyListCmd)
	historyCmd.AddCommand(historyShowCmd)
	historyCmd.AddCommand(historyBrowseCmd)

	// Flags for list command
	historyListCmd.Flags().IntVarP(&listLimit, "limit", "n", 20, "Limit number of results")
	historyListCmd.Flags().StringVar(&listStatus, "status", "", "Filter by status (success/failed)")
	historyListCmd.Flags().StringVar(&listFrom, "from", "", "Start date (YYYY-MM-DD)")
	historyListCmd.Flags().StringVar(&listTo, "to", "", "End date (YYYY-MM-DD)")
}

func runHistoryList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Initialize repository and service
	dbPath := getHistoryDBPath()
	repo, err := historyRepo.NewSQLiteRepository(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open history database: %w", err)
	}
	defer repo.Close()

	service := services.NewHistoryQueryService(repo)

	// Build filter
	filter := history.NewRecordFilter()

	// Apply status filter
	if listStatus != "" {
		outcome, err := history.NewInstallationOutcome(listStatus)
		if err != nil {
			return fmt.Errorf("invalid status: %w", err)
		}
		filter = filter.WithOutcome(outcome)
	}

	// Apply date range filter
	if listFrom != "" || listTo != "" {
		from, to, err := parseDateRange(listFrom, listTo)
		if err != nil {
			return fmt.Errorf("invalid date range: %w", err)
		}

		period, err := history.NewInstallationPeriod(from, to)
		if err != nil {
			return fmt.Errorf("invalid period: %w", err)
		}
		filter = filter.WithPeriod(period)
	}

	// Query records
	var records []history.InstallationRecord
	if listLimit > 0 {
		records, err = service.ListRecent(ctx, listLimit)
	} else {
		records, err = service.ListRecords(ctx, filter)
	}

	if err != nil {
		return fmt.Errorf("failed to query history: %w", err)
	}

	// Display results
	if len(records) == 0 {
		fmt.Println("No installation records found.")
		return nil
	}

	displayRecordsList(records)
	return nil
}

func runHistoryShow(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	recordIDStr := args[0]

	// Initialize repository and service
	dbPath := getHistoryDBPath()
	repo, err := historyRepo.NewSQLiteRepository(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open history database: %w", err)
	}
	defer repo.Close()

	service := services.NewHistoryQueryService(repo)

	// Parse record ID
	recordID, err := history.ParseRecordID(recordIDStr)
	if err != nil {
		return fmt.Errorf("invalid record ID: %w", err)
	}

	// Retrieve record
	record, err := service.GetRecordByID(ctx, recordID)
	if err != nil {
		if err == history.ErrRecordNotFound {
			return fmt.Errorf("record not found: %s", recordIDStr)
		}
		return fmt.Errorf("failed to retrieve record: %w", err)
	}

	// Display detailed record
	displayRecordDetails(record)
	return nil
}

func runHistoryBrowse(cmd *cobra.Command, args []string) error {
	// Initialize repository and service
	dbPath := getHistoryDBPath()
	repo, err := historyRepo.NewSQLiteRepository(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open history database: %w", err)
	}
	defer repo.Close()

	service := services.NewHistoryQueryService(repo)

	// Create and run browser TUI
	browser := historyTUI.NewBrowser(service)
	p := tea.NewProgram(browser, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run history browser: %w", err)
	}

	return nil
}

// displayRecordsList displays records in a table format
func displayRecordsList(records []history.InstallationRecord) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Header
	fmt.Fprintln(w, "ID\tPACKAGE\tVERSION\tSTATUS\tINSTALLED\tDURATION")
	fmt.Fprintln(w, strings.Repeat("-", 80))

	// Records
	for _, record := range records {
		id := truncateID(record.ID().String(), 8)
		pkg := record.PackageName()
		version := record.TargetVersion()
		status := formatStatus(record.Outcome())
		installed := record.InstalledAt().Format("2006-01-02 15:04")
		duration := formatDuration(record.Duration())

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			id, pkg, version, status, installed, duration)
	}

	fmt.Fprintf(w, "\nTotal: %d record(s)\n", len(records))
}

// displayRecordDetails displays complete details for a single record
func displayRecordDetails(record history.InstallationRecord) {
	fmt.Println("=== Installation Record ===")
	fmt.Println()

	// Basic info
	fmt.Printf("Record ID:      %s\n", record.ID().String())
	fmt.Printf("Session ID:     %s\n", record.SessionID())
	fmt.Printf("Package:        %s\n", record.PackageName())
	fmt.Printf("Version:        %s\n", record.TargetVersion())
	fmt.Printf("Status:         %s\n", formatStatus(record.Outcome()))
	fmt.Println()

	// Timing
	fmt.Printf("Installed At:   %s\n", record.InstalledAt().Format(time.RFC3339))
	fmt.Printf("Recorded At:    %s\n", record.RecordedAt().Format(time.RFC3339))
	fmt.Printf("Duration:       %s\n", formatDuration(record.Duration()))
	fmt.Println()

	// System context
	sysCtx := record.SystemContext()
	fmt.Println("System Context:")
	fmt.Printf("  OS:           %s\n", sysCtx.OSVersion())
	if sysCtx.KernelVersion() != "" {
		fmt.Printf("  Kernel:       %s\n", sysCtx.KernelVersion())
	}
	fmt.Printf("  Gohan:        %s\n", sysCtx.GohanVersion())
	if sysCtx.Hostname() != "" {
		fmt.Printf("  Hostname:     %s\n", sysCtx.Hostname())
	}
	fmt.Println()

	// Installed packages
	metadata := record.Metadata()
	packages := metadata.InstalledPackages()
	if len(packages) > 0 {
		fmt.Printf("Installed Packages (%d):\n", len(packages))
		for _, pkg := range packages {
			fmt.Printf("  - %s %s (%s)\n",
				pkg.Name(),
				pkg.Version(),
				formatSize(pkg.SizeBytes()))
		}
		fmt.Println()
	}

	// Failure details
	if record.HasFailureDetails() {
		fd := record.FailureDetails()
		fmt.Println("Failure Details:")
		fmt.Printf("  Reason:       %s\n", fd.Reason())
		fmt.Printf("  Failed At:    %s\n", fd.FailedAt().Format(time.RFC3339))
		fmt.Printf("  Phase:        %s\n", fd.Phase())
		if fd.ErrorCode() != "" {
			fmt.Printf("  Error Code:   %s\n", fd.ErrorCode())
		}
	}
}

// Helper functions

func getHistoryDBPath() string {
	// TODO: Make configurable via flag or config file
	homeDir, _ := os.UserHomeDir()
	gohanDir := filepath.Join(homeDir, ".gohan")

	// Ensure the directory exists
	if err := os.MkdirAll(gohanDir, 0755); err != nil {
		// Log but don't fail - let the database open fail with clearer error
		fmt.Fprintf(os.Stderr, "Warning: failed to create gohan directory: %v\n", err)
	}

	return filepath.Join(gohanDir, "history.db")
}

func truncateID(id string, length int) string {
	if len(id) <= length {
		return id
	}
	return id[:length]
}

func formatStatus(outcome history.InstallationOutcome) string {
	switch {
	case outcome.IsSuccessful():
		return "✓ Success"
	case outcome.IsFailed():
		return "✗ Failed"
	case outcome.IsRolledBack():
		return "↻ Rolled Back"
	default:
		return outcome.String()
	}
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

func parseDateRange(from, to string) (time.Time, time.Time, error) {
	layout := "2006-01-02"

	var fromTime, toTime time.Time
	var err error

	if from != "" {
		fromTime, err = time.Parse(layout, from)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid from date: %w", err)
		}
	} else {
		fromTime = time.Time{} // Zero time (far past)
	}

	if to != "" {
		toTime, err = time.Parse(layout, to)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid to date: %w", err)
		}
		// Set to end of day
		toTime = toTime.Add(24*time.Hour - time.Second)
	} else {
		toTime = time.Now() // Current time
	}

	return fromTime, toTime, nil
}
