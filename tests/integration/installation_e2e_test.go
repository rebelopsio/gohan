//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/application/history/services"
	"github.com/rebelopsio/gohan/internal/application/installation/dto"
	"github.com/rebelopsio/gohan/internal/application/installation/usecases"
	"github.com/rebelopsio/gohan/internal/config"
	"github.com/rebelopsio/gohan/internal/domain/installation"
	historyRepo "github.com/rebelopsio/gohan/internal/infrastructure/history/repository"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/packagemanager"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/repository"
	installServices "github.com/rebelopsio/gohan/internal/infrastructure/installation/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInstallationE2E tests the complete installation flow from start to finish,
// including history recording and persistence.
func TestInstallationE2E(t *testing.T) {
	// Create temporary directory for test databases
	tmpDir := t.TempDir()

	// Create test configuration
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			HistoryDB:      filepath.Join(tmpDir, "history.db"),
			InstallationDB: filepath.Join(tmpDir, "installations.db"),
		},
		Installation: config.InstallationConfig{
			SnapshotDir:          filepath.Join(tmpDir, "snapshots"),
			AutoBackup:           true,
			HistoryRetentionDays: 90,
		},
	}

	// Ensure directories exist
	err := cfg.EnsureDirectories()
	require.NoError(t, err)

	// Initialize repositories
	historyRepository, err := historyRepo.NewSQLiteRepository(cfg.Database.HistoryDB)
	require.NoError(t, err)
	defer historyRepository.Close()

	// Use memory repository for installation sessions (SQLite reconstruction not fully implemented)
	installationRepository := repository.NewMemorySessionRepository()

	// Initialize services
	historyQueryService := services.NewHistoryQueryService(historyRepository)
	historyRecordingService := services.NewHistoryRecordingService(historyRepository)
	progressEstimator := installServices.NewProgressEstimator()
	configMerger := installServices.NewConfigurationMerger()
	packageManager := packagemanager.NewAPTManagerDryRun() // Dry-run mode for testing

	// Initialize use cases
	startUseCase := usecases.NewStartInstallationUseCase(installationRepository)
	executeUseCase := usecases.NewExecuteInstallationUseCase(
		installationRepository,
		packageManager,
		progressEstimator,
		configMerger,
		packageManager,
		historyRecordingService, // History recorder is wired!
	)

	ctx := context.Background()

	t.Run("complete installation flow with history recording", func(t *testing.T) {
		// Step 1: Start installation
		startReq := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{
					Name:        "hyprland",
					Version:     "latest",
					PackageName: "hyprland",
					SizeBytes:   10000000,
				},
			},
			AvailableSpace:      50000000,
			RequiredSpace:       10000000,
			MergeExistingConfig: false,
			BackupDirectory:     filepath.Join(tmpDir, "backups"),
		}

		startResp, err := startUseCase.Execute(ctx, startReq)
		require.NoError(t, err)
		require.NotEmpty(t, startResp.SessionID)

		sessionID := startResp.SessionID
		t.Logf("Started installation session: %s", sessionID)

		// Step 2: Execute installation
		// Execute in a goroutine and wait for completion
		go func() {
			_, err := executeUseCase.Execute(ctx, sessionID)
			if err != nil {
				t.Logf("Installation execution error: %v", err)
			}
		}()

		// Step 3: Wait for installation to complete
		// Poll for completion with timeout
		timeout := time.After(10 * time.Second)
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		var finalSession *installation.InstallationSession
		for {
			select {
			case <-timeout:
				t.Fatal("Installation did not complete within timeout")
			case <-ticker.C:
				session, err := installationRepository.FindByID(ctx, sessionID)
				if err != nil {
					t.Fatalf("Failed to get session: %v", err)
				}

				status := session.Status()
				if status.IsTerminal() {
					finalSession = session
					goto completed
				}
			}
		}

	completed:
		require.NotNil(t, finalSession)
		t.Logf("Installation completed with status: %s", finalSession.Status())

		// Step 4: Verify installation session was persisted
		retrievedSession, err := installationRepository.FindByID(ctx, sessionID)
		require.NoError(t, err)
		assert.Equal(t, sessionID, retrievedSession.ID())

		// Step 5: Verify history was recorded
		// Wait a moment for history to be written
		time.Sleep(100 * time.Millisecond)

		records, err := historyQueryService.ListRecent(ctx, 10)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(records), 1, "Expected at least one history record")

		record := records[0]
		assert.NotEmpty(t, record.PackageName())
		assert.NotEmpty(t, record.ID())
		assert.NotZero(t, record.InstalledAt())
		assert.NotZero(t, record.RecordedAt())

		// Verify system context was recorded
		sysCtx := record.SystemContext()
		assert.NotEmpty(t, sysCtx.OSVersion())
		assert.NotEmpty(t, sysCtx.GohanVersion())

		// Verify metadata
		metadata := record.Metadata()
		assert.NotNil(t, metadata)

		t.Logf("History record created: %s", record.ID())
		t.Logf("  Package: %s", record.PackageName())
		t.Logf("  Outcome: %s", record.Outcome())
		t.Logf("  Installed: %s", record.InstalledAt().Format(time.RFC3339))
		t.Logf("  Duration: %s", record.Duration())
	})

	t.Run("verify history persistence across repository reopens", func(t *testing.T) {
		// Close history repository
		historyRepository.Close()

		// Reopen history repository
		historyRepo2, err := historyRepo.NewSQLiteRepository(cfg.Database.HistoryDB)
		require.NoError(t, err)
		defer historyRepo2.Close()

		// Create new query service
		queryService := services.NewHistoryQueryService(historyRepo2)

		// Verify history is still accessible
		records, err := queryService.ListRecent(ctx, 10)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(records), 1, "History records should persist across reopens")

		t.Logf("Successfully retrieved %d history records after reopen", len(records))
	})
}

// TestInstallationE2E_MultipleInstallations tests recording multiple installations
func TestInstallationE2E_MultipleInstallations(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			HistoryDB:      filepath.Join(tmpDir, "history.db"),
			InstallationDB: filepath.Join(tmpDir, "installations.db"),
		},
		Installation: config.InstallationConfig{
			SnapshotDir: filepath.Join(tmpDir, "snapshots"),
		},
	}

	err := cfg.EnsureDirectories()
	require.NoError(t, err)

	historyRepository, err := historyRepo.NewSQLiteRepository(cfg.Database.HistoryDB)
	require.NoError(t, err)
	defer historyRepository.Close()

	// Use memory repository for installation sessions
	installationRepository := repository.NewMemorySessionRepository()

	historyQueryService := services.NewHistoryQueryService(historyRepository)
	historyRecordingService := services.NewHistoryRecordingService(historyRepository)
	progressEstimator := installServices.NewProgressEstimator()
	configMerger := installServices.NewConfigurationMerger()
	packageManager := packagemanager.NewAPTManagerDryRun()

	startUseCase := usecases.NewStartInstallationUseCase(installationRepository)
	executeUseCase := usecases.NewExecuteInstallationUseCase(
		installationRepository,
		packageManager,
		progressEstimator,
		configMerger,
		packageManager,
		historyRecordingService,
	)

	ctx := context.Background()

	// Install multiple packages
	packages := []string{"hyprland", "waybar", "kitty"}

	for _, pkg := range packages {
		startReq := dto.InstallationRequest{
			Components: []dto.ComponentRequest{
				{
					Name:        pkg,
					Version:     "latest",
					PackageName: pkg,
					SizeBytes:   10000000,
				},
			},
			AvailableSpace:      50000000,
			RequiredSpace:       10000000,
			MergeExistingConfig: false,
			BackupDirectory:     filepath.Join(tmpDir, "backups"),
		}

		startResp, err := startUseCase.Execute(ctx, startReq)
		require.NoError(t, err)

		sessionID := startResp.SessionID

		// Execute and wait
		go executeUseCase.Execute(ctx, sessionID)

		// Simple wait for completion
		time.Sleep(500 * time.Millisecond)
	}

	// Verify all installations were recorded
	records, err := historyQueryService.ListRecent(ctx, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(records), len(packages),
		"Should have recorded all installations")

	t.Logf("Successfully recorded %d installations", len(records))
}

func TestMain(m *testing.M) {
	// Setup
	code := m.Run()
	// Teardown
	os.Exit(code)
}
