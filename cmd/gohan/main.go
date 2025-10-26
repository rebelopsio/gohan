package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rebelopsio/gohan/internal/tui/preflight"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Printf("gohan %s (commit: %s, built: %s)\n", version, commit, date)
			return
		case "init":
			runInit()
			return
		}
	}

	fmt.Println("Gohan - Omakase Hyprland for Debian")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  gohan init              Installation wizard")
	fmt.Println("  gohan version           Show version information")
	fmt.Println()
	fmt.Println("Run 'gohan init' to begin installation")
}

func runInit() {
	wizard := preflight.NewWizard()
	p := tea.NewProgram(wizard, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running installation wizard: %v\n", err)
		os.Exit(1)
	}
}
