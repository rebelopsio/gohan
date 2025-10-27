package main

import (
	"os"

	"github.com/rebelopsio/gohan/internal/cli/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version information
	cmd.SetVersion(version, commit, date)

	// Execute CLI
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
