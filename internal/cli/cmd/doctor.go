package cmd

import (
	"context"
	"fmt"
	"strings"

	verificationApp "github.com/rebelopsio/gohan/internal/application/verification"
	verificationInfra "github.com/rebelopsio/gohan/internal/infrastructure/verification/checkers"
	"github.com/spf13/cobra"
)

// doctorCmd represents the system health check command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run system health checks",
	Long: `Verify your Hyprland installation and configuration.

The doctor command checks critical components to ensure your
system is correctly configured and functioning properly. It verifies:
- Hyprland binary installation
- Configuration files
- Theme application
- And more...

Examples:
  # Run full health check
  gohan doctor

  # Run with progress output
  gohan doctor --progress

  # Quick check (critical only)
  gohan doctor --quick`,
	RunE: runDoctor,
}

// Flags
var (
	quickCheck bool
)

func init() {
	rootCmd.AddCommand(doctorCmd)

	// Flags
	doctorCmd.Flags().BoolVar(&quickCheck, "quick", false, "Run only critical checks")
	doctorCmd.Flags().BoolVar(&showProgress, "progress", false, "Show progress during checks")
}

func runDoctor(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Create checkers
	checkers := verificationApp.Checkers{
		HyprlandChecker: verificationInfra.NewHyprlandChecker(),
		ThemeChecker:    verificationInfra.NewThemeChecker(),
		ConfigChecker:   verificationInfra.NewConfigChecker(),
	}

	// Create use case
	useCase := verificationApp.NewDoctorUseCase(checkers)

	// Execute with or without progress
	var resp *verificationApp.DoctorResponse
	var err error

	if showProgress {
		fmt.Println("🔍 Running system health checks...")
		fmt.Println()

		resp, err = useCase.ExecuteWithProgress(
			ctx,
			verificationApp.DoctorRequest{
				ShowProgress: true,
				QuickCheck:   quickCheck,
			},
			func(checkerName string, result verificationApp.CheckResultDTO) {
				// Display progress
				status := getStatusIcon(result.Status)
				fmt.Printf("%s %s\n", status, checkerName)
			},
		)
	} else {
		resp, err = useCase.Execute(ctx, verificationApp.DoctorRequest{
			QuickCheck: quickCheck,
		})
	}

	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	// Display results
	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Printf("  SYSTEM HEALTH CHECK RESULTS\n")
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
	if resp.CriticalIssues > 0 {
		fmt.Printf("Critical Issues: %d 🔴\n", resp.CriticalIssues)
	}
	fmt.Printf("Duration:        %dms\n", resp.DurationMs)
	fmt.Println()

	// Detailed results (show failures/warnings if progress was already shown)
	if showProgress {
		for _, result := range resp.Results {
			if result.Status != "pass" {
				displayDoctorResult(result)
			}
		}
	} else {
		for _, result := range resp.Results {
			displayDoctorResult(result)
		}
	}

	// Overall status
	fmt.Println(strings.Repeat("─", 60))
	statusIcon := getStatusIcon(resp.OverallStatus)
	fmt.Printf("%s  Overall Status: %s\n", statusIcon, strings.ToUpper(resp.OverallStatus))

	// Recommendations
	if len(resp.Recommendations) > 0 {
		fmt.Println("\n💡 Recommendations:")
		seen := make(map[string]bool)
		for _, rec := range resp.Recommendations {
			if !seen[rec] {
				fmt.Printf("  • %s\n", rec)
				seen[rec] = true
			}
		}
	}

	fmt.Println(strings.Repeat("─", 60) + "\n")

	// Return error if critical issues found
	if resp.CriticalIssues > 0 {
		return fmt.Errorf("found %d critical issue(s)", resp.CriticalIssues)
	}

	return nil
}

func displayDoctorResult(result verificationApp.CheckResultDTO) {
	statusIcon := getStatusIcon(result.Status)

	fmt.Printf("\n%s %s\n", statusIcon, result.Component)
	fmt.Printf("   %s\n", result.Message)

	if len(result.Details) > 0 {
		for _, detail := range result.Details {
			fmt.Printf("   • %s\n", detail)
		}
	}

	if result.Status != "pass" && len(result.Suggestions) > 0 {
		fmt.Println("   Suggestions:")
		for _, suggestion := range result.Suggestions {
			fmt.Printf("   → %s\n", suggestion)
		}
	}
}

func getStatusIcon(status string) string {
	switch status {
	case "pass":
		return "✓"
	case "warning":
		return "⚠"
	case "fail":
		return "✗"
	default:
		return "?"
	}
}
