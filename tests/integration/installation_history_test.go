//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	historyServices "github.com/rebelopsio/gohan/internal/application/history/services"
	"github.com/rebelopsio/gohan/internal/application/installation/usecases"
	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/rebelopsio/gohan/internal/domain/installation"
	installationRepo "github.com/rebelopsio/gohan/internal/infrastructure/installation/repository"
	"github.com/rebelopsio/gohan/internal/infrastructure/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestInstallationHistoryIntegration tests the complete flow from installation to history recording
func TestInstallationHistoryIntegration(t *testing.T) {
	t.Run("successful installation is recorded to history", func(t *testing.T) {
		// Setup repositories and services
		sessionRepo := installationRepo.NewMemorySessionRepository()
		historyRepo := memory.NewHistoryRepository()
		historyRecordingService := historyServices.NewHistoryRecordingService(historyRepo)

		// Setup mocks for dependencies
		mockConflictResolver := &MockConflictResolver{}
		mockProgressEstimator := &MockProgressEstimator{}
		mockConfigMerger := &MockConfigurationMerger{}
		mockPkgManager := &MockPackageManager{}

		// Create use case with real history recording service
		useCase := usecases.NewExecuteInstallationUseCase(
			sessionRepo,
			mockConflictResolver,
			mockProgressEstimator,
			mockConfigMerger,
			mockPkgManager,
			historyRecordingService,
		)

		// Create a test installation session
		session := createTestSession(t)
		ctx := context.Background()

		// Save session to repository
		require.NoError(t, sessionRepo.Save(ctx, session))

		// Setup mock expectations
		mockConflictResolver.On("DetectConflicts", mock.Anything, mock.Anything).
			Return([]installation.PackageConflict{}, nil)
		mockProgressEstimator.On("CalculatePhaseProgress", mock.Anything, mock.Anything, mock.Anything).
			Return(50)
		mockProgressEstimator.On("EstimateRemainingTime", mock.Anything, mock.Anything, mock.Anything).
			Return(time.Duration(0))
		mockPkgManager.On("InstallPackage", mock.Anything, "hyprland", "0.45.0").
			Return(nil)

		// Execute installation
		response, err := useCase.Execute(ctx, session.ID())

		// Verify installation completed successfully
		require.NoError(t, err)
		assert.Equal(t, "completed", response.Status)
		assert.Equal(t, 1, response.ComponentsInstalled)

		// Verify history record was created
		records, err := historyRepo.FindAll(ctx, history.NewRecordFilter())
		require.NoError(t, err)
		assert.Len(t, records, 1, "Expected exactly one history record")

		record := records[0]
		assert.Equal(t, session.ID(), record.SessionID())
		assert.True(t, record.WasSuccessful())
		assert.Equal(t, "hyprland", record.PackageName())
		assert.False(t, record.HasFailureDetails())
		assert.Equal(t, 1, record.PackageCount())
	})

	t.Run("failed installation is recorded to history", func(t *testing.T) {
		// Setup repositories and services
		sessionRepo := installationRepo.NewMemorySessionRepository()
		historyRepo := memory.NewHistoryRepository()
		historyRecordingService := historyServices.NewHistoryRecordingService(historyRepo)

		// Setup mocks for dependencies
		mockConflictResolver := &MockConflictResolver{}
		mockProgressEstimator := &MockProgressEstimator{}
		mockConfigMerger := &MockConfigurationMerger{}
		mockPkgManager := &MockPackageManager{}

		// Create use case with real history recording service
		useCase := usecases.NewExecuteInstallationUseCase(
			sessionRepo,
			mockConflictResolver,
			mockProgressEstimator,
			mockConfigMerger,
			mockPkgManager,
			historyRecordingService,
		)

		// Create a test installation session
		session := createTestSession(t)
		ctx := context.Background()

		// Save session to repository
		require.NoError(t, sessionRepo.Save(ctx, session))

		// Setup mock expectations - package installation will fail
		mockConflictResolver.On("DetectConflicts", mock.Anything, mock.Anything).
			Return([]installation.PackageConflict{}, nil)
		mockProgressEstimator.On("CalculatePhaseProgress", mock.Anything, mock.Anything, mock.Anything).
			Return(50)
		mockProgressEstimator.On("EstimateRemainingTime", mock.Anything, mock.Anything, mock.Anything).
			Return(time.Duration(0))
		installErr := assert.AnError
		mockPkgManager.On("InstallPackage", mock.Anything, "hyprland", "0.45.0").
			Return(installErr)

		// Execute installation (will fail)
		response, err := useCase.Execute(ctx, session.ID())

		// Verify installation failed
		require.NoError(t, err) // Use case returns response, not error
		assert.Equal(t, "failed", response.Status)

		// Verify history record was created for failed installation
		records, err := historyRepo.FindAll(ctx, history.NewRecordFilter())
		require.NoError(t, err)
		assert.Len(t, records, 1, "Expected exactly one history record for failed installation")

		record := records[0]
		assert.Equal(t, session.ID(), record.SessionID())
		assert.True(t, record.WasFailed())
		assert.True(t, record.HasFailureDetails())
	})
}

// Helper function to create a test session
func createTestSession(t *testing.T) *installation.InstallationSession {
	// Create component
	pkgInfo, err := installation.NewPackageInfo(
		"hyprland",
		"0.45.0",
		15728640, // 15 MB
		[]string{},
	)
	require.NoError(t, err)

	component, err := installation.NewComponentSelection(
		installation.ComponentHyprland,
		"0.45.0",
		&pkgInfo,
	)
	require.NoError(t, err)

	// Create disk space
	diskSpace, err := installation.NewDiskSpace(21474836480, 0) // 20 GB available
	require.NoError(t, err)

	// Create configuration
	config, err := installation.NewInstallationConfiguration(
		[]installation.ComponentSelection{component},
		nil, // No GPU support specified
		diskSpace,
		false,
	)
	require.NoError(t, err)

	// Create session
	session, err := installation.NewInstallationSession(config)
	require.NoError(t, err)

	return session
}

// Mock implementations (reusing from usecases tests)
type MockConflictResolver struct {
	mock.Mock
}

func (m *MockConflictResolver) DetectConflicts(ctx context.Context, components []installation.ComponentSelection) ([]installation.PackageConflict, error) {
	args := m.Called(ctx, components)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]installation.PackageConflict), args.Error(1)
}

func (m *MockConflictResolver) ResolveConflict(ctx context.Context, conflict installation.PackageConflict, strategy installation.ResolutionAction) error {
	args := m.Called(ctx, conflict, strategy)
	return args.Error(0)
}

type MockProgressEstimator struct {
	mock.Mock
}

func (m *MockProgressEstimator) EstimateRemainingTime(currentPhase installation.InstallationStatus, percentComplete int, elapsedTime time.Duration) time.Duration {
	args := m.Called(currentPhase, percentComplete, elapsedTime)
	return args.Get(0).(time.Duration)
}

func (m *MockProgressEstimator) CalculatePhaseProgress(phase installation.InstallationStatus, totalItems, completedItems int) int {
	args := m.Called(phase, totalItems, completedItems)
	return args.Int(0)
}

type MockConfigurationMerger struct {
	mock.Mock
}

func (m *MockConfigurationMerger) MergeConfigurations(ctx context.Context, existing, new installation.InstallationConfiguration) (installation.InstallationConfiguration, error) {
	args := m.Called(ctx, existing, new)
	if args.Get(0) == nil {
		return installation.InstallationConfiguration{}, args.Error(1)
	}
	return args.Get(0).(installation.InstallationConfiguration), args.Error(1)
}

func (m *MockConfigurationMerger) ShouldBackupExisting(ctx context.Context, path string) (bool, error) {
	args := m.Called(ctx, path)
	return args.Bool(0), args.Error(1)
}

type MockPackageManager struct {
	mock.Mock
}

func (m *MockPackageManager) InstallPackage(ctx context.Context, packageName, version string) error {
	args := m.Called(ctx, packageName, version)
	return args.Error(0)
}

func (m *MockPackageManager) IsPackageInstalled(ctx context.Context, packageName string) (bool, error) {
	args := m.Called(ctx, packageName)
	return args.Bool(0), args.Error(1)
}
