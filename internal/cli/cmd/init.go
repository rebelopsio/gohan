package cmd

import (
	"fmt"
	"os"

	"github.com/rebelopsio/gohan/internal/config"
	"github.com/rebelopsio/gohan/internal/container"
	"github.com/spf13/cobra"
)

var (
	forceInit bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gohan configuration and databases",
	Long: `Initialize the gohan environment by creating:
  - Configuration directory (~/.gohan)
  - Default configuration file (~/.gohan/config.yaml)
  - SQLite databases (history.db, installations.db)

This command is safe to run multiple times. By default, it will not
overwrite existing configuration. Use --force to recreate everything.

Example:
  gohan init
  gohan init --force`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&forceInit, "force", "f", false, "Force initialization, overwriting existing config")
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println("Initializing gohan environment...")
	fmt.Println()

	// Get config path
	configPath := config.GetConfigPath()
	dataDir := config.GetDataDir()

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil && !forceInit {
		fmt.Printf("Configuration already exists at %s\n", configPath)
		fmt.Println("Use --force to overwrite existing configuration")
		return fmt.Errorf("configuration already exists (use --force to overwrite)")
	}

	// Create data directory
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	fmt.Printf("✓ Created data directory: %s\n", dataDir)

	// Create default config
	cfg := config.DefaultConfig()
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	fmt.Printf("✓ Created configuration file: %s\n", configPath)

	// Ensure all directories exist
	if err := cfg.EnsureDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	fmt.Printf("✓ Created database directory: %s\n", cfg.Database.HistoryDB)
	fmt.Printf("✓ Created snapshots directory: %s\n", cfg.Installation.SnapshotDir)

	// Initialize container (this will create and initialize databases)
	fmt.Println()
	fmt.Println("Initializing databases...")
	c, err := container.New()
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}
	defer c.Close()

	fmt.Printf("✓ Initialized history database: %s\n", cfg.Database.HistoryDB)
	fmt.Printf("✓ Initialized installation database: %s\n", cfg.Database.InstallationDB)

	// Display summary
	fmt.Println()
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║  Gohan initialization complete!                           ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Printf("  Config file:     %s\n", configPath)
	fmt.Printf("  Data directory:  %s\n", dataDir)
	fmt.Println()
	fmt.Println("Databases:")
	fmt.Printf("  History:         %s\n", cfg.Database.HistoryDB)
	fmt.Printf("  Installations:   %s\n", cfg.Database.InstallationDB)
	fmt.Println()
	fmt.Println("Settings:")
	fmt.Printf("  API Host:        %s\n", cfg.API.Host)
	fmt.Printf("  API Port:        %d\n", cfg.API.Port)
	fmt.Printf("  Log Level:       %s\n", cfg.Logging.Level)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  • Run 'gohan check' to verify your system")
	fmt.Println("  • Run 'gohan install' to install Hyprland")
	fmt.Println("  • Run 'gohan history browse' to view installation history")
	fmt.Println()

	return nil
}
