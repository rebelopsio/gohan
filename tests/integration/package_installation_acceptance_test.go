//go:build integration
// +build integration

package integration

import (
	"context"
	"path/filepath"
	"testing"
	"time"

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

// TestPackageInstallation_MinimalProfile corresponds to:
// Feature: Package Installation
// Scenario: Install minimal package profile
func TestPackageInstallation_MinimalProfile(t *testing.T) {
	// Given I am running Debian Sid
	// And I have network connectivity
	// And I have sufficient disk space
	tmpDir := t.TempDir()
	ctx := context.Background()

	cfg := createTestConfig(t, tmpDir)
	services := setupInstallationServices(t, cfg)

	// Given I select the "minimal" installation profile
	startReq := dto.InstallationRequest{
		Profile:         "minimal",
		AvailableSpace:  50000000000, // 50GB
		RequiredSpace:   5000000000,  // 5GB
		BackupDirectory: filepath.Join(tmpDir, "backups"),
	}

	// When I start the installation
	startResp, err := services.startUseCase.Execute(ctx, startReq)
	require.NoError(t, err, "Installation should start successfully")
	require.NotEmpty(t, startResp.SessionID)

	sessionID := startResp.SessionID

	// Execute installation asynchronously
	done := make(chan error, 1)
	go func() {
		_, err := services.executeUseCase.Execute(ctx, sessionID, nil)
		done <- err
	}()

	// Wait for installation to complete
	session := waitForInstallationCompletion(t, ctx, services.installRepo, sessionID, 30*time.Second)

	// Then the system should install Hyprland core packages
	// And the system should install essential Wayland tools
	// And the system should install a terminal emulator
	// And the system should install fonts
	installedPackages := session.InstalledPackages()
	assert.NotEmpty(t, installedPackages, "Should have installed packages")

	// Essential packages for minimal profile
	essentialPackages := []string{
		"hyprland",
		"waybar",
		"fuzzel",
		"kitty",
		"fonts-noto-color-emoji",
	}

	for _, pkg := range essentialPackages {
		assert.Contains(t, installedPackages, pkg,
			"Minimal profile should include %s", pkg)
	}

	// And all installed packages should be functional
	assert.Equal(t, installation.StatusCompleted, session.Status(),
		"Installation should complete successfully")
}

// TestPackageInstallation_RecommendedProfile corresponds to:
// Scenario: Install recommended package profile
func TestPackageInstallation_RecommendedProfile(t *testing.T) {
	t.Skip("TODO: Implement once batch installation with progress is available")

	// Given I select the "recommended" installation profile
	// When I start the installation
	// Then the system should install all minimal packages
	// And the system should install clipboard history tools
	// And the system should install media control tools
	// And the system should install network management tools
	// And all installed packages should be functional
}

// TestPackageInstallation_NetworkError corresponds to:
// Scenario: Network error during installation
func TestPackageInstallation_NetworkError(t *testing.T) {
	t.Skip("TODO: Implement once error handling and recovery are available")

	// Given I select the "minimal" installation profile
	// And network connectivity becomes unavailable during installation
	// When the installation encounters the network issue
	// Then the system should recover gracefully
	// And I should be notified of the issue
	// And the system should remain in a consistent state
}

// TestPackageInstallation_StayInformed corresponds to:
// Scenario: Stay informed during installation
func TestPackageInstallation_StayInformed(t *testing.T) {
	t.Skip("TODO: Implement once progress reporting is available")

	// Given I select the "recommended" installation profile
	// When the installation is running
	// Then I should know what is currently happening
	// And I should understand how much work remains
	// And I should know when installation is complete
}

// TestPackageInstallation_SkipAlreadyInstalled corresponds to:
// Scenario: Skip already installed packages
func TestPackageInstallation_SkipAlreadyInstalled(t *testing.T) {
	t.Skip("TODO: Implement once package detection is available")

	// Given the "hyprland" package is already installed
	// And I select the "minimal" installation profile
	// When I start the installation
	// Then the system should detect existing packages
	// And the system should skip already installed packages
	// And the system should only install missing packages
}

// TestPackageInstallation_VerifyIntegrity corresponds to:
// Scenario: Verify package integrity after installation
func TestPackageInstallation_VerifyIntegrity(t *testing.T) {
	t.Skip("TODO: Implement once post-installation verification is available")

	// Given I select the "minimal" installation profile
	// When the installation completes
	// Then all installed packages should be verified
	// And all packages should be in "installed" state
	// And no packages should be in "broken" state
}

// TestPackageInstallation_InsufficientSpace corresponds to:
// Scenario: Handle insufficient disk space
func TestPackageInstallation_InsufficientSpace(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := context.Background()

	cfg := createTestConfig(t, tmpDir)
	services := setupInstallationServices(t, cfg)

	// Given I select the "full" installation profile
	// And I have insufficient disk space for all packages
	startReq := dto.InstallationRequest{
		Profile:         "full",
		AvailableSpace:  1000000,    // Only 1MB available
		RequiredSpace:   5000000000, // 5GB required
		BackupDirectory: filepath.Join(tmpDir, "backups"),
	}

	// When I start the installation
	_, err := services.startUseCase.Execute(ctx, startReq)

	// Then the system should check disk space before installation
	// And the system should report insufficient disk space error
	// And the system should not start installing packages
	// And I should see how much space is needed
	require.Error(t, err, "Should fail with insufficient disk space")
	assert.Contains(t, err.Error(), "insufficient",
		"Error should indicate insufficient space")
}

// TestPackageInstallation_ConflictingPackages corresponds to:
// Scenario: Resolve conflicting packages
func TestPackageInstallation_ConflictingPackages(t *testing.T) {
	t.Skip("TODO: Implement once conflict detection is available")

	// Given I have a package that conflicts with Hyprland
	// And I select the "minimal" installation profile
	// When I start the installation
	// Then I should be informed of the conflict
	// And I should be able to resolve it before continuing
	// And the system should proceed only when conflict is resolved
}

// TestPackageInstallation_LatestPackages corresponds to:
// Scenario: Install latest available packages
func TestPackageInstallation_LatestPackages(t *testing.T) {
	t.Skip("TODO: Implement once package cache update is integrated")

	// Given package information is outdated
	// And I select the "minimal" installation profile
	// When I start the installation
	// Then the latest package versions should be installed
	// And I should have current software
}

// Helper function to create test configuration
func createTestConfig(t *testing.T, tmpDir string) *config.Config {
	t.Helper()

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

	err := cfg.EnsureDirectories()
	require.NoError(t, err)

	return cfg
}

// Helper to set up installation services
type installationServices struct {
	startUseCase   *usecases.StartInstallationUseCase
	executeUseCase *usecases.ExecuteInstallationUseCase
	installRepo    installation.SessionRepository
}

func setupInstallationServices(t *testing.T, cfg *config.Config) *installationServices {
	t.Helper()

	// Initialize repositories
	historyRepository, err := historyRepo.NewSQLiteRepository(cfg.Database.HistoryDB)
	require.NoError(t, err)
	t.Cleanup(func() { historyRepository.Close() })

	installationRepository := repository.NewMemorySessionRepository()

	// Initialize services
	// TODO: Replace with actual services once implemented
	progressEstimator := installServices.NewProgressEstimator()
	configMerger := installServices.NewConfigurationMerger()
	packageManager := packagemanager.NewAPTManagerDryRun()

	// Initialize use cases
	startUseCase := usecases.NewStartInstallationUseCase(installationRepository)
	executeUseCase := usecases.NewExecuteInstallationUseCase(
		installationRepository,
		packageManager,
		progressEstimator,
		configMerger,
		packageManager,
		nil, // History recording not needed for these tests
	)

	return &installationServices{
		startUseCase:   startUseCase,
		executeUseCase: executeUseCase,
		installRepo:    installationRepository,
	}
}

// Helper to wait for installation completion
func waitForInstallationCompletion(
	t *testing.T,
	ctx context.Context,
	repo installation.SessionRepository,
	sessionID string,
	timeout time.Duration,
) *installation.InstallationSession {
	t.Helper()

	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutChan:
			t.Fatal("Installation did not complete within timeout")
		case <-ticker.C:
			session, err := repo.FindByID(ctx, sessionID)
			require.NoError(t, err, "Failed to get session")

			if session.Status().IsTerminal() {
				return session
			}
		}
	}
}
