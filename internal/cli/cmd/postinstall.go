package cmd

import (
	"context"
	"fmt"
	"strings"

	postinstallApp "github.com/rebelopsio/gohan/internal/application/postinstall"
	"github.com/rebelopsio/gohan/internal/domain/postinstall"
	postinstallInfra "github.com/rebelopsio/gohan/internal/infrastructure/postinstall"
	"github.com/spf13/cobra"
)

// postInstallCmd represents the post-installation command
var postInstallCmd = &cobra.Command{
	Use:   "post-install",
	Short: "Run post-installation configuration",
	Long: `Configure essential components after Hyprland installation.

The post-install command sets up:
- Display manager (SDDM, GDM, or TTY launch)
- Shell configuration with theme
- Audio system (PipeWire)
- Network manager
- Wallpaper cache generation

Examples:
  # Interactive setup (prompts for each component)
  gohan post-install

  # Setup with specific display manager
  gohan post-install --display-manager sddm

  # Complete setup with all components
  gohan post-install --display-manager sddm --shell zsh --audio --network

  # TTY launch (no display manager)
  gohan post-install --display-manager tty

  # Show progress during setup
  gohan post-install --progress`,
	RunE: runPostInstall,
}

// Flags
var (
	displayManagerFlag string
	shellFlag          string
	shellThemeFlag     string
	setupAudioFlag     bool
	setupNetworkFlag   bool
	wallpaperDirFlag   string
)

func init() {
	rootCmd.AddCommand(postInstallCmd)

	// Flags
	postInstallCmd.Flags().StringVar(&displayManagerFlag, "display-manager", "", "Display manager to install (sddm, gdm, tty, none)")
	postInstallCmd.Flags().StringVar(&shellFlag, "shell", "", "Shell to configure (zsh, bash, fish)")
	postInstallCmd.Flags().StringVar(&shellThemeFlag, "shell-theme", "default", "Theme for shell configuration")
	postInstallCmd.Flags().BoolVar(&setupAudioFlag, "audio", false, "Setup audio system (PipeWire)")
	postInstallCmd.Flags().BoolVar(&setupNetworkFlag, "network", false, "Setup network manager")
	postInstallCmd.Flags().StringVar(&wallpaperDirFlag, "wallpaper-dir", "/usr/share/wallpapers", "Directory containing wallpapers")
	postInstallCmd.Flags().BoolVar(&showProgress, "progress", false, "Show progress during setup")
}

func runPostInstall(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Parse display manager
	var dm postinstall.DisplayManager
	if displayManagerFlag != "" {
		dm = postinstall.DisplayManager(displayManagerFlag)
		if !dm.IsValid() {
			return fmt.Errorf("invalid display manager: %s (valid options: sddm, gdm, tty, none)", displayManagerFlag)
		}
	}

	// Parse shell
	var shell postinstall.Shell
	if shellFlag != "" {
		shell = postinstall.Shell(shellFlag)
		if !shell.IsValid() {
			return fmt.Errorf("invalid shell: %s (valid options: zsh, bash, fish)", shellFlag)
		}
	}

	// Create infrastructure components
	packageMgr := postinstallInfra.NewAPTPackageManagerAdapter()
	serviceMgr := postinstallInfra.NewSystemdServiceManager()

	// Create installers
	installers := postinstallApp.Installers{}

	if dm != "" && dm != postinstall.DisplayManagerNone {
		installers.DisplayManagerInstaller = postinstallInfra.NewDisplayManagerInstaller(dm, packageMgr, serviceMgr)
	}

	if shell != "" {
		installers.ShellInstaller = postinstallInfra.NewShellInstaller(shell, shellThemeFlag, packageMgr)
	}

	if setupAudioFlag {
		installers.AudioInstaller = postinstallInfra.NewAudioInstaller(packageMgr, serviceMgr)
	}

	if setupNetworkFlag {
		installers.NetworkInstaller = postinstallInfra.NewNetworkInstaller(packageMgr, serviceMgr)
	}

	if wallpaperDirFlag != "" {
		installers.WallpaperGenerator = postinstallInfra.NewWallpaperCacheGenerator(wallpaperDirFlag)
	}

	// Create use case
	useCase := postinstallApp.NewPostInstallUseCase(installers)

	// Build request
	request := postinstallApp.PostInstallRequest{
		DisplayManager: dm,
		Shell:          shell,
		ShellTheme:     shellThemeFlag,
		SetupAudio:     setupAudioFlag,
		SetupNetwork:   setupNetworkFlag,
		WallpaperDir:   wallpaperDirFlag,
		ShowProgress:   showProgress,
	}

	// Execute with or without progress
	var resp *postinstallApp.PostInstallResponse
	var err error

	if showProgress {
		fmt.Println("🔧 Running post-installation setup...")
		fmt.Println()

		resp, err = useCase.ExecuteWithProgress(
			ctx,
			request,
			func(installerName string, result postinstallApp.ComponentResultDTO) {
				// Display progress
				status := getSetupStatusIcon(result.Status)
				fmt.Printf("%s %s\n", status, installerName)
			},
		)
	} else {
		resp, err = useCase.Execute(ctx, request)
	}

	if err != nil {
		return fmt.Errorf("post-installation failed: %w", err)
	}

	// Display results
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Printf("  POST-INSTALLATION RESULTS\n")
	fmt.Println(strings.Repeat("═", 60) + "\n")

	// Component status
	fmt.Printf("Display Manager:  %s\n", formatBoolStatus(resp.DisplayManagerConfigured))
	fmt.Printf("Shell:            %s\n", formatBoolStatus(resp.ShellConfigured))
	fmt.Printf("Audio:            %s\n", formatBoolStatus(resp.AudioConfigured))
	fmt.Printf("Network:          %s\n", formatBoolStatus(resp.NetworkManagerConfigured))
	fmt.Printf("Wallpaper Cache:  %s\n", formatBoolStatus(resp.WallpaperCacheGenerated))
	fmt.Printf("Duration:         %dms\n", resp.DurationMs)
	fmt.Println()

	// Detailed results (show failures/warnings if progress was already shown)
	if showProgress {
		for _, result := range resp.ComponentResults {
			if result.Status != "completed" {
				displayPostInstallResult(result)
			}
		}
	} else {
		for _, result := range resp.ComponentResults {
			displayPostInstallResult(result)
		}
	}

	// Overall status
	fmt.Println(strings.Repeat("─", 60))
	statusIcon := getSetupStatusIcon(resp.OverallStatus)
	fmt.Printf("%s  Overall Status: %s\n", statusIcon, strings.ToUpper(resp.OverallStatus))

	// Recommendations
	if len(resp.Recommendations) > 0 {
		fmt.Println("\n💡 Recommendations:")
		for _, rec := range resp.Recommendations {
			fmt.Printf("  • %s\n", rec)
		}
	}

	// Failures
	if len(resp.ServicesFailed) > 0 {
		fmt.Println("\n⚠ Failed Components:")
		for _, fail := range resp.ServicesFailed {
			fmt.Printf("  • %s\n", fail)
		}
	}

	fmt.Println(strings.Repeat("─", 60) + "\n")

	// Return error if setup failed
	if resp.OverallStatus == "failed" {
		return fmt.Errorf("post-installation completed with failures")
	}

	return nil
}

func displayPostInstallResult(result postinstallApp.ComponentResultDTO) {
	statusIcon := getSetupStatusIcon(result.Status)

	fmt.Printf("\n%s %s\n", statusIcon, result.Component)
	fmt.Printf("   %s\n", result.Message)

	if len(result.Details) > 0 {
		for _, detail := range result.Details {
			fmt.Printf("   • %s\n", detail)
		}
	}

	if result.Error != "" {
		fmt.Printf("   Error: %s\n", result.Error)
	}
}

func getSetupStatusIcon(status string) string {
	switch status {
	case "completed":
		return "✓"
	case "in_progress":
		return "⏳"
	case "failed":
		return "✗"
	case "skipped":
		return "⊘"
	case "pending":
		return "○"
	default:
		return "?"
	}
}

func formatBoolStatus(enabled bool) string {
	if enabled {
		return "✓ Configured"
	}
	return "○ Not configured"
}
