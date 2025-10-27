package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rebelopsio/gohan/internal/application/installation/dto"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [session-id]",
	Short: "Get installation status",
	Long: `Get the status of an installation session. This command requires
an API server to be running.

Examples:
  # Get status of a specific session
  gohan status abc123-def456

  # Get status with custom API URL
  gohan status abc123-def456 --api-url http://server:8080`,
	Args: cobra.ExactArgs(1),
	RunE: runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	sessionID := args[0]

	fmt.Printf("Fetching status for session: %s\n", sessionID)
	logVerbose("API URL: %s", apiURL)

	// TODO: This endpoint needs to be implemented in the next phase
	// For now, we'll show what the API call would look like
	resp, err := http.Get(fmt.Sprintf("%s/api/installation/%s/status", apiURL, sessionID))
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned error: %s - %s", resp.Status, string(bodyBytes))
	}

	var statusResponse dto.InstallationProgressResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResponse); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Display status
	fmt.Println("\nInstallation Status:")
	fmt.Printf("  Session ID:    %s\n", statusResponse.SessionID)
	fmt.Printf("  Status:        %s\n", statusResponse.Status)
	fmt.Printf("  Phase:         %s\n", statusResponse.CurrentPhase)
	fmt.Printf("  Progress:      %d%%\n", statusResponse.PercentComplete)
	fmt.Printf("  Components:    %d/%d installed\n", statusResponse.ComponentsInstalled, statusResponse.ComponentsTotal)
	if statusResponse.EstimatedRemaining != "0s" {
		fmt.Printf("  Est. Time:     %s\n", statusResponse.EstimatedRemaining)
	}
	if statusResponse.Message != "" {
		fmt.Printf("  Message:       %s\n", statusResponse.Message)
	}

	return nil
}
