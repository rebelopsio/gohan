package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information (set by build flags)
	version = "dev"
	commit  = "none"
	date    = "unknown"

	// Global flags
	apiURL  string
	verbose bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "gohan",
	Short: "Gohan - Omakase Hyprland for Debian",
	Long: `Gohan is an opinionated Hyprland installer for Debian that provides
a curated, batteries-included desktop environment experience.

Complete documentation is available at https://github.com/rebelopsio/gohan`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "http://localhost:8080", "API server URL")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(historyCmd)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gohan %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", date)
	},
}

// SetVersion sets version information (useful for testing and build)
func SetVersion(v, c, d string) {
	version = v
	commit = c
	date = d
}

// logVerbose prints verbose output if enabled
func logVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, "[verbose] "+format+"\n", args...)
	}
}
