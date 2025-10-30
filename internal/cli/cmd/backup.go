package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	backupApp "github.com/rebelopsio/gohan/internal/application/backup"
	backupInfra "github.com/rebelopsio/gohan/internal/infrastructure/installation/backup"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage configuration backups",
	Long: `Create, restore, and manage backups of your configuration files.

The backup command helps you safely store and restore your Hyprland and related
configuration files, ensuring you can recover from mistakes or system changes.`,
}

// backupCreateCmd creates a new backup
var backupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new backup",
	Long: `Create a backup of your current configuration files.

This creates a timestamped backup of your Hyprland, Waybar, Kitty, and other
configuration files.

Examples:
  # Create a backup with description
  gohan backup create --description "Pre-update backup"

  # Create backup of specific paths
  gohan backup create --paths ~/.config/hypr,~/.config/waybar`,
	RunE: runBackupCreate,
}

// backupListCmd lists all available backups
var backupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all backups",
	Long: `Display a list of all available backups.

Shows backup ID, description, creation time, file count, and total size.

Examples:
  # List all backups
  gohan backup list`,
	RunE: runBackupList,
}

// backupRestoreCmd restores a backup
var backupRestoreCmd = &cobra.Command{
	Use:   "restore <backup-id>",
	Short: "Restore a backup",
	Long: `Restore configuration files from a backup.

This replaces your current configuration files with those from the backup.

Examples:
  # Restore a specific backup
  gohan backup restore 2025-01-29_120000

  # Restore only Hyprland configs
  gohan backup restore 2025-01-29_120000 --selective ~/.config/hypr`,
	Args: cobra.ExactArgs(1),
	RunE: runBackupRestore,
}

// backupCleanupCmd removes old backups
var backupCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up old backups",
	Long: `Remove old backups based on retention policy.

This helps manage disk space by removing backups older than the specified
retention period, while always keeping a minimum number of backups.

Examples:
  # Remove backups older than 30 days, keep at least 5
  gohan backup cleanup --retention-days 30 --keep-minimum 5

  # Dry run to see what would be removed
  gohan backup cleanup --retention-days 30 --dry-run`,
	RunE: runBackupCleanup,
}

// Flags
var (
	backupDescription string
	backupPaths       []string
	backupSelective   []string
	retentionDays     int
	keepMinimum       int
	backupDryRun      bool
)

func init() {
	rootCmd.AddCommand(backupCmd)

	// Add subcommands
	backupCmd.AddCommand(backupCreateCmd)
	backupCmd.AddCommand(backupListCmd)
	backupCmd.AddCommand(backupRestoreCmd)
	backupCmd.AddCommand(backupCleanupCmd)

	// Create flags
	backupCreateCmd.Flags().StringVar(&backupDescription, "description", "", "Backup description")
	backupCreateCmd.Flags().StringSliceVar(&backupPaths, "paths", nil, "Specific paths to backup")

	// Restore flags
	backupRestoreCmd.Flags().StringSliceVar(&backupSelective, "selective", nil, "Restore only specific paths")

	// Cleanup flags
	backupCleanupCmd.Flags().IntVar(&retentionDays, "retention-days", 30, "Keep backups newer than this many days")
	backupCleanupCmd.Flags().IntVar(&keepMinimum, "keep-minimum", 5, "Always keep at least this many backups")
	backupCleanupCmd.Flags().BoolVar(&backupDryRun, "dry-run", false, "Show what would be removed without actually removing")
}

func runBackupCreate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get backup root
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	backupRoot := filepath.Join(homeDir, ".config", "gohan", "backups")

	// Create repository
	repo := backupInfra.NewRepositoryAdapter(backupRoot)

	// Create use case
	useCase := backupApp.NewCreateBackupUseCase(repo)

	// Determine paths to back up
	paths := backupPaths
	if len(paths) == 0 {
		// Default paths
		configDir := filepath.Join(homeDir, ".config")
		paths = []string{
			filepath.Join(configDir, "hypr"),
			filepath.Join(configDir, "waybar"),
			filepath.Join(configDir, "kitty"),
			filepath.Join(configDir, "rofi"),
			filepath.Join(configDir, "mako"),
		}
	}

	// Execute use case
	req := backupApp.CreateBackupRequest{
		Description: backupDescription,
		FilePaths:   paths,
		BackupRoot:  backupRoot,
	}

	resp, err := useCase.Execute(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Display success
	fmt.Printf("✓ Backup created successfully\n\n")
	fmt.Printf("Backup ID:    %s\n", resp.BackupID)
	fmt.Printf("Files:        %d\n", resp.FileCount)
	fmt.Printf("Total Size:   %s\n", formatBytes(resp.TotalSize))
	fmt.Printf("Location:     %s\n", resp.BackupPath)

	return nil
}

func runBackupList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get backup root
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	backupRoot := filepath.Join(homeDir, ".config", "gohan", "backups")

	// Create repository
	repo := backupInfra.NewRepositoryAdapter(backupRoot)

	// Create use case
	useCase := backupApp.NewListBackupsUseCase(repo)

	// Execute use case
	resp, err := useCase.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	if resp.Total == 0 {
		fmt.Println("No backups found.")
		return nil
	}

	// Display backups in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "ID\tDESCRIPTION\tAGE\tFILES\tSIZE\tSTATUS\n")
	fmt.Fprintf(w, "──\t───────────\t───\t─────\t────\t──────\n")

	for _, b := range resp.Backups {
		desc := b.Description
		if desc == "" {
			desc = "-"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n",
			b.ID, desc, b.Age, b.FileCount, b.TotalSize, b.Status)
	}

	w.Flush()
	fmt.Printf("\nTotal: %d backups\n", resp.Total)

	return nil
}

func runBackupRestore(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	backupID := args[0]

	// Get backup root
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	backupRoot := filepath.Join(homeDir, ".config", "gohan", "backups")

	// Create repository
	repo := backupInfra.NewRepositoryAdapter(backupRoot)

	// Create use case
	useCase := backupApp.NewRestoreBackupUseCase(repo)

	// Execute use case
	req := backupApp.RestoreBackupRequest{
		BackupID:  backupID,
		Selective: backupSelective,
	}

	resp, err := useCase.Execute(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	// Display success
	fmt.Printf("✓ Backup restored successfully\n\n")
	fmt.Printf("Backup ID:       %s\n", resp.BackupID)
	fmt.Printf("Files Restored:  %d\n", resp.FilesRestored)

	if len(resp.RestoredPaths) > 0 && len(resp.RestoredPaths) <= 10 {
		fmt.Println("\nRestored files:")
		for _, path := range resp.RestoredPaths {
			fmt.Printf("  - %s\n", path)
		}
	}

	return nil
}

func runBackupCleanup(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get backup root
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	backupRoot := filepath.Join(homeDir, ".config", "gohan", "backups")

	// Create repository
	repo := backupInfra.NewRepositoryAdapter(backupRoot)

	// Create use case
	useCase := backupApp.NewCleanupBackupsUseCase(repo)

	// Execute use case
	req := backupApp.CleanupBackupsRequest{
		RetentionDays: retentionDays,
		KeepMinimum:   keepMinimum,
		DryRun:        backupDryRun,
	}

	resp, err := useCase.Execute(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to cleanup backups: %w", err)
	}

	// Display results
	if backupDryRun {
		fmt.Printf("Dry run - would remove %d backups\n\n", resp.RemovedCount)
	} else {
		fmt.Printf("✓ Cleanup completed\n\n")
	}

	fmt.Printf("Removed:        %d backups\n", resp.RemovedCount)
	fmt.Printf("Freed:          %s\n", formatBytes(resp.FreedBytes))
	fmt.Printf("Remaining:      %d backups\n", resp.RemainingCount)

	if len(resp.RemovedIDs) > 0 && len(resp.RemovedIDs) <= 10 {
		fmt.Println("\nRemoved backups:")
		for _, id := range resp.RemovedIDs {
			fmt.Printf("  - %s\n", id)
		}
	}

	return nil
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
