package cmd

import (
	"context"
	"fmt"

	repoApp "github.com/rebelopsio/gohan/internal/application/repository"
	repoInfra "github.com/rebelopsio/gohan/internal/infrastructure/repository"
	"github.com/spf13/cobra"
)

// repoCmd represents the repository management command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage Debian repositories",
	Long: `Manage Debian repository configuration for Hyprland installation.

The repo command helps you configure Debian repositories correctly for
installing Hyprland and its dependencies. It can detect your Debian version,
check repository status, enable non-free components, and manage sources.list.`,
}

// detectVersionCmd detects the Debian version
var detectVersionCmd = &cobra.Command{
	Use:   "detect-version",
	Short: "Detect Debian version",
	Long: `Detect the current Debian version and check compatibility.

This command detects whether you're running Debian Sid (unstable),
Trixie (testing), Bookworm (stable), or Ubuntu. It also checks if
the version is supported for Hyprland installation.

Examples:
  # Detect current Debian version
  gohan repo detect-version`,
	RunE: runDetectVersion,
}

// repoCheckCmd checks repository configuration
var repoCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check repository configuration",
	Long: `Check the current sources.list configuration.

This command analyzes your /etc/apt/sources.list file to check which
repository components are enabled (main, contrib, non-free, etc.) and
whether deb-src entries are available.

Examples:
  # Check current repository configuration
  gohan repo check`,
	RunE: runRepoCheck,
}

// enableNonFreeCmd enables non-free repositories
var enableNonFreeCmd = &cobra.Command{
	Use:   "enable-nonfree",
	Short: "Enable non-free repositories",
	Long: `Enable non-free and non-free-firmware repository components.

This is required for NVIDIA GPU users who need proprietary drivers.
The command will backup your sources.list before making changes.

Examples:
  # Enable non-free repositories
  gohan repo enable-nonfree

  # Enable without creating backup
  gohan repo enable-nonfree --no-backup`,
	RunE: runEnableNonFree,
}

// enableDebSrcCmd enables deb-src repositories
var enableDebSrcCmd = &cobra.Command{
	Use:   "enable-debsrc",
	Short: "Enable source repositories (deb-src)",
	Long: `Enable deb-src entries for all repository sources.

Source repositories are needed when building packages from source.
This is required for installing some Hyprland ecosystem components
that aren't available as pre-built packages.

Examples:
  # Enable deb-src repositories
  gohan repo enable-debsrc

  # Enable without creating backup
  gohan repo enable-debsrc --no-backup`,
	RunE: runEnableDebSrc,
}

// backupSourcesCmd backs up sources.list
var backupSourcesCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup sources.list",
	Long: `Create a timestamped backup of /etc/apt/sources.list.

This creates a backup that can be restored if repository changes
cause problems. Backups are stored in /etc/apt/ by default.

Examples:
  # Backup sources.list
  gohan repo backup

  # Backup to custom directory
  gohan repo backup --dir /path/to/backup`,
	RunE: runBackupSources,
}

// Flags
var (
	noBackup  bool
	backupDir string
)

func init() {
	rootCmd.AddCommand(repoCmd)

	// Add subcommands
	repoCmd.AddCommand(detectVersionCmd)
	repoCmd.AddCommand(repoCheckCmd)
	repoCmd.AddCommand(enableNonFreeCmd)
	repoCmd.AddCommand(enableDebSrcCmd)
	repoCmd.AddCommand(backupSourcesCmd)

	// Flags for enable commands
	enableNonFreeCmd.Flags().BoolVar(&noBackup, "no-backup", false, "Skip creating backup before changes")
	enableDebSrcCmd.Flags().BoolVar(&noBackup, "no-backup", false, "Skip creating backup before changes")

	// Flags for backup command
	backupSourcesCmd.Flags().StringVar(&backupDir, "dir", "/etc/apt", "Directory to store backup")
}

func runDetectVersion(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create version detector
	detector := repoInfra.NewSystemVersionDetector()

	// Create use case
	useCase := repoApp.NewDetectVersionUseCase(detector)

	// Execute
	resp, err := useCase.Execute(ctx, repoApp.DetectVersionRequest{})
	if err != nil {
		return fmt.Errorf("failed to detect version: %w", err)
	}

	// Display results
	fmt.Printf("🐧 %s\n\n", resp.DisplayString)

	fmt.Printf("Codename:    %s\n", resp.Codename)
	fmt.Printf("Version:     %s\n", resp.Version)

	if resp.IsSupported {
		fmt.Printf("Supported:   ✓ Yes\n")
	} else {
		fmt.Printf("Supported:   ✗ No\n")
	}

	// Version-specific info
	if resp.IsSid {
		fmt.Printf("Type:        Sid (unstable) - bleeding edge\n")
	} else if resp.IsTrixie {
		fmt.Printf("Type:        Trixie (testing) - 10 days behind Sid\n")
	} else if resp.IsBookworm {
		fmt.Printf("Type:        Bookworm (stable) - NOT recommended\n")
	} else if resp.IsUbuntu {
		fmt.Printf("Type:        Ubuntu - NOT supported\n")
	}

	// Support message
	if resp.SupportMessage != "" {
		fmt.Printf("\n%s\n", resp.SupportMessage)
	}

	return nil
}

func runRepoCheck(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create sources manager
	manager := repoInfra.NewFileSourcesManager()

	// Create use case
	useCase := repoApp.NewCheckRepositoryUseCase(manager)

	// Execute
	resp, err := useCase.Execute(ctx, repoApp.CheckRepositoryRequest{
		SourcesListPath: "/etc/apt/sources.list",
	})
	if err != nil {
		return fmt.Errorf("failed to check repositories: %w", err)
	}

	// Display results
	fmt.Printf("📦 Repository Configuration\n\n")

	fmt.Printf("Components:\n")
	fmt.Printf("  main:              %s\n", formatBool(resp.HasMain))
	fmt.Printf("  contrib:           %s\n", formatBool(resp.HasContrib))
	fmt.Printf("  non-free:          %s\n", formatBool(resp.HasNonFree))
	fmt.Printf("  non-free-firmware: %s\n", formatBool(resp.HasNonFreeFirmware))

	fmt.Printf("\nEntries:\n")
	fmt.Printf("  deb entries:       %d\n", resp.DebEntries)
	fmt.Printf("  deb-src entries:   %d\n", resp.DebSrcEntries)
	fmt.Printf("  total entries:     %d\n", resp.TotalEntries)

	fmt.Printf("\nSource repositories: %s\n", formatBool(resp.HasDebSrc))

	// Recommendations
	fmt.Printf("\n💡 Recommendations:\n")
	if !resp.HasNonFree {
		fmt.Printf("  • Enable non-free for NVIDIA drivers: gohan repo enable-nonfree\n")
	}
	if !resp.HasDebSrc {
		fmt.Printf("  • Enable source repos for building packages: gohan repo enable-debsrc\n")
	}
	if resp.HasNonFree && resp.HasDebSrc {
		fmt.Printf("  ✓ Repository configuration looks good!\n")
	}

	return nil
}

func runEnableNonFree(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create sources manager
	manager := repoInfra.NewFileSourcesManager()

	// Create use case
	useCase := repoApp.NewEnableNonFreeUseCase(manager)

	// Execute
	resp, err := useCase.Execute(ctx, repoApp.EnableNonFreeRequest{
		SourcesListPath: "/etc/apt/sources.list",
		BackupFirst:     !noBackup,
	})
	if err != nil {
		return fmt.Errorf("failed to enable non-free: %w", err)
	}

	// Display results
	if !resp.Modified {
		fmt.Printf("ℹ️  Non-free repositories already enabled\n")
		return nil
	}

	fmt.Printf("✓ Non-free repositories enabled\n\n")

	if resp.BackupPath != "" {
		fmt.Printf("Backup created: %s\n", resp.BackupPath)
	}

	fmt.Printf("Components added:\n")
	for _, comp := range resp.ComponentsAdded {
		fmt.Printf("  • %s\n", comp)
	}

	fmt.Printf("\n💡 Next steps:\n")
	fmt.Printf("  1. Update package lists: sudo apt update\n")
	fmt.Printf("  2. Install NVIDIA drivers: sudo apt install nvidia-driver\n")

	return nil
}

func runEnableDebSrc(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create sources manager
	manager := repoInfra.NewFileSourcesManager()

	// Create use case
	useCase := repoApp.NewEnableDebSrcUseCase(manager)

	// Execute
	resp, err := useCase.Execute(ctx, repoApp.EnableDebSrcRequest{
		SourcesListPath: "/etc/apt/sources.list",
		BackupFirst:     !noBackup,
	})
	if err != nil {
		return fmt.Errorf("failed to enable deb-src: %w", err)
	}

	// Display results
	if !resp.Modified {
		fmt.Printf("ℹ️  Source repositories already enabled\n")
		return nil
	}

	fmt.Printf("✓ Source repositories enabled\n\n")

	if resp.BackupPath != "" {
		fmt.Printf("Backup created: %s\n", resp.BackupPath)
	}

	fmt.Printf("Entries added:  %d\n", resp.EntriesAdded)

	fmt.Printf("\n💡 Next steps:\n")
	fmt.Printf("  1. Update package lists: sudo apt update\n")
	fmt.Printf("  2. Now you can build packages from source with apt-get source\n")

	return nil
}

func runBackupSources(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create sources manager
	manager := repoInfra.NewFileSourcesManager()

	// Create use case
	useCase := repoApp.NewBackupSourcesListUseCase(manager)

	// Execute
	resp, err := useCase.Execute(ctx, repoApp.BackupSourcesListRequest{
		SourcesListPath: "/etc/apt/sources.list",
		BackupDir:       backupDir,
	})
	if err != nil {
		return fmt.Errorf("failed to backup sources.list: %w", err)
	}

	// Display results
	fmt.Printf("✓ sources.list backed up successfully\n\n")
	fmt.Printf("Backup path:  %s\n", resp.BackupPath)
	fmt.Printf("Timestamp:    %s\n", resp.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("Size:         %d bytes\n", resp.Size)

	return nil
}

func formatBool(b bool) string {
	if b {
		return "✓"
	}
	return "✗"
}
