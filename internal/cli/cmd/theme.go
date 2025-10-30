package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	themeApp "github.com/rebelopsio/gohan/internal/application/theme"
	"github.com/rebelopsio/gohan/internal/container"
	"github.com/rebelopsio/gohan/internal/domain/theme"
	themeInfra "github.com/rebelopsio/gohan/internal/infrastructure/theme"
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

// themeSetCmd applies a theme
var themeSetCmd = &cobra.Command{
	Use:   "set <theme-name>",
	Short: "Apply a theme",
	Long: `Apply a theme to your desktop environment.

This will update configuration files for Hyprland, Waybar, Kitty, and other
components to use the selected theme.

Examples:
  # Apply the latte theme
  gohan theme set latte

  # Apply the mocha theme
  gohan theme set mocha`,
	Args: cobra.ExactArgs(1),
	RunE: runThemeSet,
}

// themeRollbackCmd rolls back to previous theme
var themeRollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to previous theme",
	Long: `Rollback to the previously active theme.

This command allows you to undo theme changes and restore the previous theme.
You can use this command multiple times to step back through your theme history.

Examples:
  # Rollback to previous theme
  gohan theme rollback

  # Rollback multiple times
  gohan theme rollback
  gohan theme rollback`,
	Args: cobra.NoArgs,
	RunE: runThemeRollback,
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
	themeCmd.AddCommand(themeSetCmd)
	themeCmd.AddCommand(themeRollbackCmd)

	// Flags for list command
	themeListCmd.Flags().StringVar(&variantFilter, "variant", "", "Filter by variant (dark or light)")

	// Flags for show command
	themeShowCmd.Flags().BoolVarP(&verboseOutput, "verbose", "v", false, "Show detailed color information")
}

func runThemeList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Initialize theme registry with saved state
	registry, err := initializeThemeRegistry(ctx)
	if err != nil {
		return err
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

	fmt.Println("\nAvailable Themes:")

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

	// Initialize theme registry with saved state
	registry, err := initializeThemeRegistry(ctx)
	if err != nil {
		return err
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

	// Initialize theme registry with saved state
	registry, err := initializeThemeRegistry(ctx)
	if err != nil {
		return err
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

func runThemeSet(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	themeName := args[0]

	// Initialize dependency container
	c, err := container.New()
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}
	defer c.Close()

	// Initialize theme registry with saved state
	registry, err := initializeThemeRegistry(ctx)
	if err != nil {
		return err
	}

	// Load saved theme state if exists
	if err := loadSavedThemeState(ctx, registry, c.ThemeStateStore); err != nil {
		// Don't fail - just log warning and continue with default
		fmt.Printf("Warning: failed to load saved theme state: %v\n", err)
	}

	// Create use case with real theme applier, state store, and history store from container
	applyUseCase := themeApp.NewApplyThemeUseCase(registry, c.ThemeApplier, c.ThemeStateStore, c.ThemeHistoryStore)

	// Execute
	fmt.Printf("Applying theme '%s'...\n", themeName)
	result, err := applyUseCase.Execute(ctx, themeName)
	if err != nil {
		return fmt.Errorf("failed to apply theme: %w", err)
	}

	// Display result
	if result.Success {
		fmt.Printf("\n✓ %s\n\n", result.Message)
		fmt.Println("Theme applied successfully!")
		fmt.Println("Updated configuration files:")
		fmt.Println("  - Hyprland configuration")
		fmt.Println("  - Waybar configuration")
		fmt.Println("  - Kitty terminal colors")
		fmt.Println("  - Rofi/Fuzzel theme")
		fmt.Println("\nBackups have been created. Use 'gohan theme rollback' to restore previous theme.")
		fmt.Println()
	}

	return nil
}

// loadSavedThemeState loads the saved theme state and sets it as active
func loadSavedThemeState(ctx context.Context, registry theme.ThemeRegistry, stateStore themeInfra.ThemeStateStore) error {
	// Check if state exists
	exists, err := stateStore.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check theme state: %w", err)
	}

	if !exists {
		// No saved state - use default
		return nil
	}

	// Load state
	state, err := stateStore.Load(ctx)
	if err != nil {
		return fmt.Errorf("failed to load theme state: %w", err)
	}

	// Find the theme
	th, err := registry.FindByName(ctx, state.ThemeName)
	if err != nil {
		// Theme doesn't exist - use default
		return fmt.Errorf("saved theme '%s' not found: %w", state.ThemeName, err)
	}

	// Set as active
	if err := registry.SetActive(ctx, th.Name()); err != nil {
		return fmt.Errorf("failed to set active theme: %w", err)
	}

	return nil
}

func runThemeRollback(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Initialize dependency container
	c, err := container.New()
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}
	defer c.Close()

	// Initialize theme registry with saved state
	registry, err := initializeThemeRegistry(ctx)
	if err != nil {
		return err
	}

	// Load saved theme state if exists
	if err := loadSavedThemeState(ctx, registry, c.ThemeStateStore); err != nil {
		// Don't fail - just log warning and continue with default
		fmt.Printf("Warning: failed to load saved theme state: %v\n", err)
	}

	// Create rollback use case
	rollbackUseCase := themeApp.NewRollbackThemeUseCase(
		registry,
		c.ThemeHistoryStore,
		c.ThemeApplier,
		c.ThemeStateStore,
	)

	// Execute
	fmt.Println("Rolling back to previous theme...")
	result, err := rollbackUseCase.Execute(ctx)
	if err != nil {
		// Check if it's a "no history" error
		if err.Error() == "cannot rollback: no theme history available" {
			fmt.Println("\n✗ No theme history available")
			fmt.Println("You need to change themes at least once before you can rollback.")
			return nil
		}
		return fmt.Errorf("failed to rollback theme: %w", err)
	}

	// Display result
	if result.Success {
		fmt.Printf("\n✓ %s\n\n", result.Message)
		fmt.Println("Theme rolled back successfully!")
		fmt.Printf("Restored theme: %s\n", result.RestoredTheme)
		fmt.Printf("Previous theme: %s\n", result.PreviousTheme)
		fmt.Println("\nConfiguration files have been updated.")
		fmt.Println("Use 'gohan theme rollback' again to continue rolling back through history.")
		fmt.Println()
	}

	return nil
}

// initializeThemeRegistry initializes the theme registry and loads saved state
func initializeThemeRegistry(ctx context.Context) (theme.ThemeRegistry, error) {
	registry := theme.NewThemeRegistry()
	if err := theme.InitializeStandardThemes(registry); err != nil {
		return nil, fmt.Errorf("failed to initialize themes: %w", err)
	}

	// Load saved state if exists
	stateFilePath, _ := themeInfra.GetDefaultStateFilePath()
	stateStore := themeInfra.NewFileThemeStateStore(stateFilePath)

	if err := loadSavedThemeState(ctx, registry, stateStore); err != nil {
		// Don't fail - just use default theme
		// Silently ignore errors for commands that just read state
	}

	return registry, nil
}
