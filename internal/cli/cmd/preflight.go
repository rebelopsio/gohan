package cmd

import (
	"context"
	"fmt"
	"strings"

	preflightApp "github.com/rebelopsio/gohan/internal/application/preflight"
	preflightInfra "github.com/rebelopsio/gohan/internal/infrastructure/preflight/detectors"
	"github.com/spf13/cobra"
)

// preflightCmd represents the preflight checks command
var preflightCmd = &cobra.Command{
	Use:   "preflight",
	Short: "Run preflight checks before installation",
	Long: `Run system validation checks to ensure your environment is ready
for Hyprland installation.

The preflight command checks:
- Debian version compatibility (Sid or Trixie required)
- GPU detection and driver requirements
- Available disk space (minimum 10GB)
- Internet connectivity
- Source repository configuration`,
}

// preflightCheckCmd runs all preflight checks
var preflightCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Run all preflight checks",
	Long: `Execute all system validation checks to verify installation readiness.

This will check your system configuration and report any issues that
would prevent successful installation of Hyprland. Blocking issues must
be resolved before installation, while warnings are recommended fixes.

Examples:
  # Run preflight checks
  gohan preflight check

  # Run with progress output
  gohan preflight check --progress`,
	RunE: runPreflightCheck,
}

// Flags
var (
	showProgress bool
)

func init() {
	rootCmd.AddCommand(preflightCmd)

	// Add subcommands
	preflightCmd.AddCommand(preflightCheckCmd)

	// Flags
	preflightCheckCmd.Flags().BoolVar(&showProgress, "progress", false, "Show progress as checks run")
}

func runPreflightCheck(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create detectors
	detectors := preflightApp.Detectors{
		DebianDetector:          preflightInfra.NewDebianVersionDetector(),
		GPUDetector:             preflightInfra.NewSystemGPUDetector(),
		DiskSpaceDetector:       preflightInfra.NewSystemDiskSpaceDetector(),
		ConnectivityChecker:     preflightInfra.NewSystemConnectivityChecker(),
		SourceRepositoryChecker: preflightInfra.NewSystemSourceRepositoryChecker(),
	}

	// Create use case
	useCase := preflightApp.NewRunPreflightUseCase(detectors)

	// Execute with or without progress
	var resp *preflightApp.RunPreflightResponse
	var err error

	if showProgress {
		fmt.Println("🔍 Running preflight checks...\n")

		resp, err = useCase.ExecuteWithProgress(
			ctx,
			preflightApp.RunPreflightRequest{ShowProgress: true},
			func(validatorName string, result preflightApp.CheckResult) {
				// Display progress
				status := "✓"
				if !result.Passed {
					if result.Blocking {
						status = "✗"
					} else {
						status = "⚠"
					}
				}
				fmt.Printf("%s %s\n", status, validatorName)
			},
		)
	} else {
		resp, err = useCase.Execute(ctx, preflightApp.RunPreflightRequest{})
	}

	if err != nil {
		return fmt.Errorf("preflight checks failed: %w", err)
	}

	// Display results
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Printf("  PREFLIGHT CHECK RESULTS\n")
	fmt.Println(strings.Repeat("═", 60) + "\n")

	// Summary stats
	fmt.Printf("Total Checks:    %d\n", resp.TotalChecks)
	fmt.Printf("Passed:          %d ✓\n", resp.PassedChecks)
	if resp.WarningChecks > 0 {
		fmt.Printf("Warnings:        %d ⚠\n", resp.WarningChecks)
	}
	if resp.FailedChecks > 0 {
		fmt.Printf("Failed:          %d ✗\n", resp.FailedChecks)
	}
	fmt.Println()

	// Detailed results
	if !showProgress {
		for _, result := range resp.Results {
			displayCheckResult(result)
		}
	} else {
		// Only show failures and warnings if progress was already shown
		for _, result := range resp.Results {
			if !result.Passed {
				displayCheckResult(result)
			}
		}
	}

	// Overall status
	fmt.Println(strings.Repeat("─", 60))
	if resp.Passed {
		if resp.HasWarnings {
			fmt.Printf("⚠  %s\n", resp.OverallMessage)
			fmt.Println("\nInstallation can proceed, but some warnings should be addressed.")
		} else {
			fmt.Printf("✓  %s\n", resp.OverallMessage)
		}
	} else {
		fmt.Printf("✗  %s\n", resp.OverallMessage)
		fmt.Println("\nPlease resolve the blocking issues above before attempting installation.")
		return fmt.Errorf("preflight checks failed with %d blocking issue(s)", resp.FailedChecks)
	}
	fmt.Println(strings.Repeat("─", 60) + "\n")

	return nil
}

func displayCheckResult(result preflightApp.CheckResult) {
	// Status icon
	status := "✓"
	statusColor := "green"
	if !result.Passed {
		if result.Blocking {
			status = "✗"
			statusColor = "red"
		} else {
			status = "⚠"
			statusColor = "yellow"
		}
	}

	// Result name and status
	fmt.Printf("\n%s %s\n", status, result.Name)

	// Message
	if result.Message != "" {
		fmt.Printf("   %s\n", result.Message)
	}

	// Guidance (only for failures/warnings)
	if !result.Passed && result.Guidance != "" {
		fmt.Printf("\n   💡 %s\n", result.Guidance)
	}

	_ = statusColor // For future color output support
}
