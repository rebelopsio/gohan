package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rebelopsio/gohan/internal/tui/preflight"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run preflight checks",
	Long: `Run preflight checks to verify system compatibility and requirements
for Hyprland installation. This will check:
  - Debian version
  - Hardware capabilities (GPU, RAM, disk space)
  - Required dependencies
  - Network connectivity

The wizard will guide you through the process interactively.`,
	RunE: runCheck,
}

func init() {
	checkCmd.Flags().Bool("interactive", true, "Run interactive wizard (default)")
	checkCmd.Flags().Bool("json", false, "Output results as JSON")
}

func runCheck(cmd *cobra.Command, args []string) error {
	interactive, _ := cmd.Flags().GetBool("interactive")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	if jsonOutput {
		// TODO: Implement non-interactive JSON output
		return fmt.Errorf("JSON output not yet implemented")
	}

	if !interactive {
		// TODO: Implement non-interactive mode
		return fmt.Errorf("non-interactive mode not yet implemented")
	}

	// Run interactive wizard
	logVerbose("Starting preflight check wizard")

	wizard := preflight.NewWizard()
	p := tea.NewProgram(wizard, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running preflight wizard: %w", err)
	}

	return nil
}
