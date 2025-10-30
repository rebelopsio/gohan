package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	configApp "github.com/rebelopsio/gohan/internal/application/configuration"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/backup"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/configservice"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/templates"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Hyprland configurations",
	Long: `Deploy, backup, and manage Hyprland configuration files.

The config command handles deployment of pre-configured Hyprland
configurations optimized for your system.`,
}

// configDeployCmd deploys configurations
var configDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy Hyprland configurations",
	Long: `Deploy configuration files for Hyprland and related components.

Automatically backs up existing configurations before deploying new ones.
Templates are personalized for your system (username, paths, etc.).

Examples:
  # Deploy all configurations
  gohan config deploy

  # Deploy specific components
  gohan config deploy --components hyprland,waybar

  # Preview without deploying
  gohan config deploy --dry-run

  # Deploy with progress
  gohan config deploy --progress`,
	RunE: runConfigDeploy,
}

// configListCmd lists available configurations
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available configuration components",
	Run:   runConfigList,
}

// Flags
var (
	configComponents  []string
	configDryRun      bool
	configForce       bool
	configSkipBackup  bool
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configDeployCmd)
	configCmd.AddCommand(configListCmd)

	// Deploy flags
	configDeployCmd.Flags().StringSliceVar(&configComponents, "components", []string{}, "Components to deploy (hyprland,waybar,kitty,fuzzel)")
	configDeployCmd.Flags().BoolVar(&configDryRun, "dry-run", false, "Preview deployment without making changes")
	configDeployCmd.Flags().BoolVar(&configForce, "force", false, "Force deployment without prompting")
	configDeployCmd.Flags().BoolVar(&configSkipBackup, "skip-backup", false, "Skip backup of existing configurations")
	configDeployCmd.Flags().BoolVar(&showProgress, "progress", false, "Show progress during deployment")
}

func runConfigDeploy(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create infrastructure components
	templateEngine := templates.NewTemplateEngine()

	// Use default backup location
	homeDir, _ := os.UserHomeDir()
	backupRoot := filepath.Join(homeDir, ".local/share/gohan/backups")
	backupService := backup.NewBackupService(backupRoot)

	deployer := configservice.NewConfigDeployer(templateEngine, backupService)

	// Create use case
	useCase := configApp.NewConfigDeployUseCase(deployer, templateEngine)

	// Build request
	request := configApp.DeployConfigRequest{
		Components:   configComponents,
		DryRun:       configDryRun,
		Force:        configForce,
		SkipBackup:   configSkipBackup,
		ShowProgress: showProgress,
		CustomVars:   make(map[string]string),
	}

	// Execute with or without progress
	var resp *configApp.DeployConfigResponse
	var err error

	if showProgress {
		fmt.Println("📦 Deploying configurations...")
		fmt.Println()

		resp, err = useCase.ExecuteWithProgress(
			ctx,
			request,
			func(component string, filePath string, progress float64) {
				fmt.Printf("  [%.0f%%] %s: %s\n", progress, component, filePath)
			},
		)
	} else {
		resp, err = useCase.Execute(ctx, request)
	}

	if err != nil {
		return fmt.Errorf("configuration deployment failed: %w", err)
	}

	// Display results
	fmt.Println("\n" + strings.Repeat("═", 60))
	if resp.DryRun {
		fmt.Printf("  CONFIGURATION DEPLOYMENT PREVIEW\n")
	} else {
		fmt.Printf("  CONFIGURATION DEPLOYMENT RESULTS\n")
	}
	fmt.Println(strings.Repeat("═", 60) + "\n")

	// Summary stats
	fmt.Printf("Total Files:      %d\n", resp.TotalFiles)
	if !resp.DryRun {
		fmt.Printf("Deployed:         %d ✓\n", resp.SuccessfulFiles)
		if resp.FailedFiles > 0 {
			fmt.Printf("Failed:           %d ✗\n", resp.FailedFiles)
		}
		if resp.SkippedFiles > 0 {
			fmt.Printf("Skipped:          %d ⊘\n", resp.SkippedFiles)
		}
	}

	// Backup info
	if resp.BackupID != "" {
		fmt.Printf("\n📦 Backup Created:\n")
		fmt.Printf("   ID:   %s\n", resp.BackupID)
		fmt.Printf("   Path: %s\n", resp.BackupPath)
	}

	fmt.Println()

	// Detailed file list
	if len(resp.DeployedFiles) > 0 {
		fmt.Println("Files:")
		for _, file := range resp.DeployedFiles {
			icon := getDeployStatusIcon(file.Status)
			fmt.Printf("  %s [%s] %s\n", icon, file.Component, file.TargetPath)
			if file.Error != "" {
				fmt.Printf("     Error: %s\n", file.Error)
			}
		}
	}

	fmt.Println(strings.Repeat("─", 60))

	// Status message
	if resp.DryRun {
		fmt.Println("ℹ️  This was a dry-run. Run without --dry-run to deploy.")
	} else if resp.FailedFiles > 0 {
		fmt.Println("⚠  Some files failed to deploy. Check errors above.")
	} else {
		fmt.Println("✓  Configuration deployment completed successfully!")
	}

	fmt.Println(strings.Repeat("─", 60) + "\n")

	// Return error if any files failed
	if !resp.DryRun && resp.FailedFiles > 0 {
		return fmt.Errorf("%d file(s) failed to deploy", resp.FailedFiles)
	}

	return nil
}

func runConfigList(cmd *cobra.Command, args []string) {
	fmt.Println("Available Configuration Components:\n")

	components := []struct {
		name        string
		description string
		files       []string
	}{
		{
			name:        "hyprland",
			description: "Core Hyprland window manager configuration",
			files:       []string{"~/.config/hypr/hyprland.conf"},
		},
		{
			name:        "waybar",
			description: "Status bar configuration",
			files:       []string{"~/.config/waybar/config.jsonc", "~/.config/waybar/style.css"},
		},
		{
			name:        "kitty",
			description: "Terminal emulator configuration",
			files:       []string{"~/.config/kitty/kitty.conf"},
		},
		{
			name:        "fuzzel",
			description: "Application launcher configuration",
			files:       []string{"~/.config/fuzzel/fuzzel.ini"},
		},
	}

	for _, comp := range components {
		fmt.Printf("• %s\n", comp.name)
		fmt.Printf("  %s\n", comp.description)
		fmt.Printf("  Files:\n")
		for _, file := range comp.files {
			fmt.Printf("    - %s\n", file)
		}
		fmt.Println()
	}

	fmt.Println("Deploy specific components:")
	fmt.Println("  gohan config deploy --components hyprland,waybar")
}

func getDeployStatusIcon(status string) string {
	switch status {
	case "deployed":
		return "✓"
	case "failed":
		return "✗"
	case "skipped":
		return "⊘"
	case "dry-run":
		return "ℹ"
	default:
		return "?"
	}
}
