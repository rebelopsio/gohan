package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rebelopsio/gohan/internal/application/installation/dto"
	"github.com/rebelopsio/gohan/internal/config"
	"github.com/rebelopsio/gohan/internal/container"
	installTUI "github.com/rebelopsio/gohan/internal/tui/installation"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	components     []string
	gpuVendor      string
	availableSpace uint64
	requiredSpace  uint64
	useAPI         bool
	dryRun         bool
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Hyprland and components",
	Long: `Install Hyprland and selected components. Can run locally or connect
to a remote API server.

Examples:
  # Install with default components
  gohan install

  # Install specific components
  gohan install --components hyprland,waybar,rofi

  # Dry-run mode (no actual installation)
  gohan install --dry-run

  # Use remote API
  gohan install --use-api --api-url http://server:8080

  # Specify GPU vendor
  gohan install --gpu amd`,
	RunE: runInstall,
}

func init() {
	installCmd.Flags().StringSliceVar(&components, "components", []string{"hyprland"}, "Components to install (comma-separated)")
	installCmd.Flags().StringVar(&gpuVendor, "gpu", "", "GPU vendor (amd, nvidia, intel)")
	installCmd.Flags().Uint64Var(&availableSpace, "available-space", 107374182400, "Available disk space in bytes (default: 100GB)")
	installCmd.Flags().Uint64Var(&requiredSpace, "required-space", 10737418240, "Required disk space in bytes (default: 10GB)")
	installCmd.Flags().BoolVar(&useAPI, "use-api", false, "Use remote API instead of local execution")
	installCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Dry-run mode (no actual installation)")
}

func runInstall(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Build installation request
	request := buildInstallationRequest()

	logVerbose("Installation request: %+v", request)

	if useAPI {
		return runInstallViaAPI(ctx, request)
	}

	return runInstallLocal(ctx, request)
}

func buildInstallationRequest() dto.InstallationRequest {
	// Convert component names to requests
	var componentRequests []dto.ComponentRequest
	for _, comp := range components {
		componentRequests = append(componentRequests, dto.ComponentRequest{
			Name:    comp,
			Version: "latest", // TODO: Support version specification
		})
	}

	request := dto.InstallationRequest{
		Components:     componentRequests,
		AvailableSpace: availableSpace,
		RequiredSpace:  requiredSpace,
	}

	// Add GPU if specified
	if gpuVendor != "" {
		request.GPU = &dto.GPURequest{
			Vendor:         gpuVendor,
			RequiresDriver: true,
			DriverName:     gpuVendor + "_driver",
		}
	}

	return request
}

func runInstallLocal(ctx context.Context, request dto.InstallationRequest) error {
	fmt.Println("Starting local installation...")

	// Set dry-run mode in config if flag is set
	if dryRun {
		// Load config to modify it
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		cfg.Installation.DryRun = true
		// Save temporarily (will be reset on next load)
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Println("Running in dry-run mode (no actual installation)")
	}

	// Initialize dependency container (will use dry-run setting from config)
	c, err := container.New()
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}
	defer c.Close()

	// Start installation using pre-wired use cases
	response, err := c.StartInstallationUseCase.Execute(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to start installation: %w", err)
	}

	// Get package name and version for display
	packageName := "hyprland"
	packageVersion := "latest"
	if len(request.Components) > 0 {
		packageName = request.Components[0].Name
		packageVersion = request.Components[0].Version
	}

	// Create progress channel
	progressChan := make(chan installTUI.ProgressUpdate, 100)

	// Launch installation in a goroutine with progress updates
	go func() {
		defer close(progressChan)

		// Initial progress
		progressChan <- installTUI.ProgressUpdate{
			Phase:           "Starting Installation",
			PercentComplete: 0,
			Message:         fmt.Sprintf("Session created: %s", response.SessionID),
			ComponentsTotal: response.ComponentCount,
		}

		time.Sleep(500 * time.Millisecond)

		// Simulate progress updates (in real implementation, these would come from the use case)
		phases := []struct {
			name    string
			percent int
			message string
		}{
			{"Checking Requirements", 10, "Verifying system requirements"},
			{"Resolving Dependencies", 25, "Analyzing package dependencies"},
			{"Downloading Packages", 40, "Downloading required packages"},
			{"Installing Dependencies", 60, "Installing dependency packages"},
			{"Installing Main Package", 80, "Installing " + packageName},
			{"Configuring", 90, "Applying configuration"},
			{"Finalizing", 95, "Cleaning up temporary files"},
		}

		for i, phase := range phases {
			time.Sleep(300 * time.Millisecond)
			progressChan <- installTUI.ProgressUpdate{
				Phase:               phase.name,
				PercentComplete:     phase.percent,
				Message:             phase.message,
				ComponentsInstalled: i,
				ComponentsTotal:     response.ComponentCount,
			}
		}

		// Execute the actual installation
		progress, err := c.ExecuteInstallationUseCase.Execute(ctx, response.SessionID)

		// Final update
		if err != nil {
			progressChan <- installTUI.ProgressUpdate{
				Phase:           "Failed",
				PercentComplete: 100,
				Message:         "Installation failed",
				IsComplete:      true,
				IsError:         true,
				ErrorMessage:    err.Error(),
			}
		} else if progress.Status == "completed" {
			progressChan <- installTUI.ProgressUpdate{
				Phase:               "Completed",
				PercentComplete:     100,
				Message:             "Installation completed successfully!",
				ComponentsInstalled: progress.ComponentsInstalled,
				ComponentsTotal:     progress.ComponentsTotal,
				IsComplete:          true,
			}
		} else {
			progressChan <- installTUI.ProgressUpdate{
				Phase:        "Failed",
				Message:      progress.Message,
				IsComplete:   true,
				IsError:      true,
				ErrorMessage: progress.Message,
			}
		}
	}()

	// Run the TUI
	viewer := installTUI.NewProgressViewer(packageName, packageVersion, progressChan)
	p := tea.NewProgram(viewer, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run progress viewer: %w", err)
	}

	fmt.Println("\nView installation history with: gohan history browse")
	return nil
}

func runInstallViaAPI(ctx context.Context, request dto.InstallationRequest) error {
	fmt.Printf("Connecting to API server at %s...\n", apiURL)

	// Start installation via API
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(apiURL+"/api/installation/start", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned error: %s - %s", resp.Status, string(bodyBytes))
	}

	var startResponse dto.InstallationResponse
	if err := json.NewDecoder(resp.Body).Decode(&startResponse); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("Session created: %s\n", startResponse.SessionID)

	// Execute installation
	fmt.Println("Executing installation...")
	resp, err = http.Post(apiURL+"/api/installation/"+startResponse.SessionID+"/execute", "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to execute installation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned error: %s - %s", resp.Status, string(bodyBytes))
	}

	var progressResponse dto.InstallationProgressResponse
	if err := json.NewDecoder(resp.Body).Decode(&progressResponse); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("\nInstallation %s\n", progressResponse.Status)
	fmt.Printf("Progress: %d%%\n", progressResponse.PercentComplete)

	if progressResponse.Status == "completed" {
		fmt.Println("\n✓ Installation completed successfully!")
	} else {
		fmt.Printf("\n✗ Installation failed: %s\n", progressResponse.Message)
	}

	return nil
}
