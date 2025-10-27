// +build e2e

package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/application/installation/dto"
	"github.com/rebelopsio/gohan/internal/application/installation/usecases"
	httpinfra "github.com/rebelopsio/gohan/internal/infrastructure/http"
	"github.com/rebelopsio/gohan/internal/infrastructure/http/handlers"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/packagemanager"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/repository"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_E2E_InstallationFlow(t *testing.T) {
	// Setup complete server
	server := setupTestServer(t)
	defer server.Shutdown(context.Background())

	baseURL := "http://localhost:18080" // Use test port

	t.Run("complete installation flow", func(t *testing.T) {
		// Step 1: Start installation
		startRequest := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{
					Name:    "hyprland",
					Version: "0.35.0",
				},
			},
			AvailableSpace: 100 * 1024 * 1024 * 1024, // 100 GB
			RequiredSpace:  10 * 1024 * 1024 * 1024,  // 10 GB
		}

		body, _ := json.Marshal(startRequest)
		resp, err := http.Post(baseURL+"/api/installation/start", "application/json", bytes.NewReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var startResponse dto.InstallationResponse
		err = json.NewDecoder(resp.Body).Decode(&startResponse)
		require.NoError(t, err)

		assert.NotEmpty(t, startResponse.SessionID)
		assert.Equal(t, "pending", startResponse.Status)
		assert.Equal(t, 1, startResponse.ComponentCount)

		sessionID := startResponse.SessionID

		// Step 2: Execute installation
		resp, err = http.Post(baseURL+"/api/installation/"+sessionID+"/execute", "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var progressResponse dto.InstallationProgressResponse
		err = json.NewDecoder(resp.Body).Decode(&progressResponse)
		require.NoError(t, err)

		assert.Equal(t, sessionID, progressResponse.SessionID)
		assert.Equal(t, "completed", progressResponse.Status)
		assert.Equal(t, 100, progressResponse.PercentComplete)
	})

	t.Run("health check endpoint", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var health map[string]string
		err = json.NewDecoder(resp.Body).Decode(&health)
		require.NoError(t, err)

		assert.Equal(t, "ok", health["status"])
	})
}

func setupTestServer(t *testing.T) *httpinfra.Server {
	t.Helper()

	// Initialize dependencies
	sessionRepo := repository.NewMemorySessionRepository()
	aptManager := packagemanager.NewAPTManagerDryRun() // Dry-run for testing
	progressEstimator := services.NewProgressEstimator()
	configMerger := services.NewConfigurationMerger()

	startUseCase := usecases.NewStartInstallationUseCase(sessionRepo)
	executeUseCase := usecases.NewExecuteInstallationUseCase(
		sessionRepo,
		aptManager,
		progressEstimator,
		configMerger,
		aptManager,
	)

	installationHandler := handlers.NewInstallationHandler(startUseCase, executeUseCase)

	// Create server config
	config := httpinfra.Config{
		Host:         "localhost",
		Port:         18080,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server := httpinfra.NewServer(config, installationHandler)

	// Start server in background
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	return server
}
