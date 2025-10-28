package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	themeApp "github.com/rebelopsio/gohan/internal/application/theme"
	"github.com/rebelopsio/gohan/internal/domain/theme"
	"github.com/spf13/cobra"
)

// themeCmd represents the theme command
var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage visual themes",
	Long: `Manage visual themes for your Hyprland environment.

The theme command allows you to list, preview, and apply themes that change
the appearance of your desktop environment including Hyprland, Waybar, Kitty,
and other components.`,
}

// themeListCmd lists all available themes
var themeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available themes",
	Long: `Display a list of all available themes.

Examples:
  # List all themes
  gohan theme list

  # List only dark themes
  gohan theme list --variant dark

  # List only light themes
  gohan theme list --variant light`,
	RunE: runThemeList,
}

// themeShowCmd shows the currently active theme
var themeShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show active theme",
	Long: `Display information about the currently active theme.

Examples:
  # Show active theme
  gohan theme show

  # Show active theme with detailed colors
  gohan theme show --verbose`,
	RunE: runThemeShow,
}

// themePreviewCmd previews a theme without applying it
var themePreviewCmd = &cobra.Command{
	Use:   "preview <theme-name>",
	Short: "Preview a theme",
	Long: `Preview a theme's colors and appearance without applying it.

This allows you to see what a theme looks like before switching to it.

Examples:
  # Preview the latte theme
  gohan theme preview latte

  # Preview the mocha theme
  gohan theme preview mocha`,
	Args: cobra.ExactArgs(1),
	RunE: runThemePreview,
}

var (
	variantFilter string
	verboseOutput bool
)

func init() {
	rootCmd.AddCommand(themeCmd)
	themeCmd.AddCommand(themeListCmd)
	themeCmd.AddCommand(themeShowCmd)
	themeCmd.AddCommand(themePreviewCmd)

	// Flags for list command
	themeListCmd.Flags().StringVar(&variantFilter, "variant", "", "Filter by variant (dark or light)")

	// Flags for show command
	themeShowCmd.Flags().BoolVarP(&verboseOutput, "verbose", "v", false, "Show detailed color information")
}

func runThemeList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Initialize theme registry
	registry := theme.NewThemeRegistry()
	if err := theme.InitializeStandardThemes(registry); err != nil {
		return fmt.Errorf("failed to initialize themes: %w", err)
	}

	// Create use case
	listUseCase := themeApp.NewListThemesUseCase(registry)

	// Execute
	themes, err := listUseCase.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to list themes: %w", err)
	}

	// Filter by variant if specified
	if variantFilter != "" {
		filtered := make([]themeApp.ThemeInfo, 0)
		for _, t := range themes {
			if t.Variant == variantFilter {
				filtered = append(filtered, t)
			}
		}
		themes = filtered
	}

	// Display results
	if len(themes) == 0 {
		fmt.Println("No themes found.")
		return nil
	}

	fmt.Println("\nAvailable Themes:\n")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, t := range themes {
		activeMarker := " "
		if t.IsActive {
			activeMarker = "●"
		}

		variantDisplay := fmt.Sprintf("(%s)", t.Variant)
		fmt.Fprintf(w, "  %s %s\t%s\t%s\t[%s]\n",
			activeMarker,
			t.Name,
			t.DisplayName,
			variantDisplay,
			t.Author)
	}
	w.Flush()

	fmt.Println("\nUse 'gohan theme preview <name>' to preview a theme")
	return nil
}

func runThemeShow(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Initialize theme registry
	registry := theme.NewThemeRegistry()
	if err := theme.InitializeStandardThemes(registry); err != nil {
		return fmt.Errorf("failed to initialize themes: %w", err)
	}

	// Create use case
	getActiveUseCase := themeApp.NewGetActiveThemeUseCase(registry)

	// Execute
	activeTheme, err := getActiveUseCase.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active theme: %w", err)
	}

	// Display results
	fmt.Printf("\nActive Theme: %s\n\n", activeTheme.DisplayName)
	fmt.Printf("  Name:    %s\n", activeTheme.Name)
	fmt.Printf("  Author:  %s\n", activeTheme.Author)
	fmt.Printf("  Variant: %s\n", activeTheme.Variant)

	if activeTheme.Description != "" {
		fmt.Printf("  Description: %s\n", activeTheme.Description)
	}

	if verboseOutput && len(activeTheme.ColorScheme) > 0 {
		fmt.Println("\nColors:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for name, color := range activeTheme.ColorScheme {
			fmt.Fprintf(w, "  %s:\t%s\n", name, color)
		}
		w.Flush()
	}

	fmt.Println()
	return nil
}

func runThemePreview(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	themeName := args[0]

	// Initialize theme registry
	registry := theme.NewThemeRegistry()
	if err := theme.InitializeStandardThemes(registry); err != nil {
		return fmt.Errorf("failed to initialize themes: %w", err)
	}

	// Create use case
	previewUseCase := themeApp.NewPreviewThemeUseCase(registry)

	// Execute
	preview, err := previewUseCase.Execute(ctx, themeName)
	if err != nil {
		return fmt.Errorf("failed to preview theme: %w", err)
	}

	// Display preview
	fmt.Println(preview.PreviewText)

	return nil
}
