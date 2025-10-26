package main

import (
	"fmt"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("gohan %s (commit: %s, built: %s)\n", version, commit, date)
		return
	}

	fmt.Println("Gohan - Omakase Hyprland for Debian")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  gohan init              Installation wizard")
	fmt.Println("  gohan version           Show version information")
	fmt.Println()
	fmt.Println("Run 'gohan init' to begin installation")
}
