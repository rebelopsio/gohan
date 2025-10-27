package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/application/installation/usecases"
	"github.com/rebelopsio/gohan/internal/domain/installation"
	"github.com/rebelopsio/gohan/internal/domain/preflight"
	preflightTUI "github.com/rebelopsio/gohan/internal/tui/preflight"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockInstallationSessionRepository is a mock implementation of the session repository
type MockInstallationSessionRepository struct {
	mock.Mock
}

func (m *MockInstallationSessionRepository) Save(ctx context.Context, session *installation.InstallationSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockInstallationSessionRepository) FindByID(ctx context.Context, id string) (*installation.InstallationSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*installation.InstallationSession), args.Error(1)
}

func (m *MockInstallationSessionRepository) List(ctx context.Context) ([]*installation.InstallationSession, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*installation.InstallationSession), args.Error(1)
}

// MockConflictResolver is a mock implementation of ConflictResolver
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

// MockProgressEstimator is a mock implementation of ProgressEstimator
type MockProgressEstimator struct {
	mock.Mock
}

func (m *MockProgressEstimator) EstimateRemainingTime(
	currentPhase installation.InstallationStatus,
	percentComplete int,
	elapsedTime time.Duration,
) time.Duration {
	args := m.Called(currentPhase, percentComplete, elapsedTime)
	return args.Get(0).(time.Duration)
}

func (m *MockProgressEstimator) CalculatePhaseProgress(
	phase installation.InstallationStatus,
	totalItems, completedItems int,
) int {
	args := m.Called(phase, totalItems, completedItems)
	return args.Int(0)
}

// MockConfigurationMerger is a mock implementation of ConfigurationMerger
type MockConfigurationMerger struct {
	mock.Mock
}

func (m *MockConfigurationMerger) MergeConfigurations(
	ctx context.Context,
	existing, new installation.InstallationConfiguration,
) (installation.InstallationConfiguration, error) {
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

// MockPackageManager is a mock implementation of package manager
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

// MockPreflightValidator is a mock implementation of preflight validator
type MockPreflightValidator struct {
	mock.Mock
	progressChan chan preflightTUI.ProgressUpdate
	session      *preflight.ValidationSession
}

func NewMockPreflightValidator() *MockPreflightValidator {
	session := preflight.NewValidationSession()
	session.Complete() // Mark as complete with no failures
	return &MockPreflightValidator{
		progressChan: make(chan preflightTUI.ProgressUpdate),
		session:      session,
	}
}

func (m *MockPreflightValidator) Run(ctx context.Context) error {
	args := m.Called(ctx)
	close(m.progressChan) // Close channel to signal completion
	return args.Error(0)
}

func (m *MockPreflightValidator) Session() *preflight.ValidationSession {
	return m.session
}

func (m *MockPreflightValidator) Progress() <-chan preflightTUI.ProgressUpdate {
	return m.progressChan
}

func TestExecuteInstallationUseCase_Execute(t *testing.T) {
	t.Run("successfully executes installation with no conflicts", func(t *testing.T) {
		// Create a valid installation session
		components, err := createTestComponents()
		require.NoError(t, err)

		diskSpace, err := installation.NewDiskSpace(
			100*uint64(installation.GB),
			10*uint64(installation.GB),
		)
		require.NoError(t, err)

		config, err := installation.NewInstallationConfiguration(
			components,
			nil,
			diskSpace,
			false,
		)
		require.NoError(t, err)

		session, err := installation.NewInstallationSession(config)
		require.NoError(t, err)

		// Setup mocks
		mockRepo := new(MockInstallationSessionRepository)
		mockConflictResolver := new(MockConflictResolver)
		mockProgressEstimator := new(MockProgressEstimator)
		mockConfigMerger := new(MockConfigurationMerger)
		mockPkgManager := new(MockPackageManager)
		mockPreflight := NewMockPreflightValidator()

		// Mock expectations
		mockRepo.On("FindByID", mock.Anything, session.ID()).
			Return(session, nil)
		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*installation.InstallationSession")).
			Return(nil)

		mockConflictResolver.On("DetectConflicts", mock.Anything, mock.Anything).
			Return([]installation.PackageConflict{}, nil)

		mockProgressEstimator.On("CalculatePhaseProgress", mock.Anything, mock.Anything, mock.Anything).
			Return(50)
		mockProgressEstimator.On("EstimateRemainingTime", mock.Anything, mock.Anything, mock.Anything).
			Return(5 * time.Minute)

		mockPkgManager.On("InstallPackage", mock.Anything, "hyprland", "0.35.0").
			Return(nil)

		mockPreflight.On("Run", mock.Anything).Return(nil)

		// Execute use case
		useCase := usecases.NewExecuteInstallationUseCase(
			mockRepo,
			mockConflictResolver,
			mockProgressEstimator,
			mockConfigMerger,
			mockPkgManager,
			nil, // historyRecorder not needed for this test
			mockPreflight,
		)
		ctx := context.Background()

		response, err := useCase.Execute(ctx, session.ID(), nil)

		require.NoError(t, err)
		assert.Equal(t, session.ID(), response.SessionID)
		assert.NotEmpty(t, response.Status)
		mockRepo.AssertExpectations(t)
		mockConflictResolver.AssertExpectations(t)
	})

	t.Run("detects and handles conflicts", func(t *testing.T) {
		components, err := createTestComponents()
		require.NoError(t, err)

		diskSpace, err := installation.NewDiskSpace(
			100*uint64(installation.GB),
			10*uint64(installation.GB),
		)
		require.NoError(t, err)

		config, err := installation.NewInstallationConfiguration(
			components,
			nil,
			diskSpace,
			false,
		)
		require.NoError(t, err)

		session, err := installation.NewInstallationSession(config)
		require.NoError(t, err)

		// Create a conflict
		conflict, err := installation.NewPackageConflict(
			"hyprland",
			"hyprland-git",
			"conflicting package versions",
		)
		require.NoError(t, err)

		// Setup mocks
		mockRepo := new(MockInstallationSessionRepository)
		mockConflictResolver := new(MockConflictResolver)
		mockProgressEstimator := new(MockProgressEstimator)
		mockConfigMerger := new(MockConfigurationMerger)
		mockPkgManager := new(MockPackageManager)
		mockPreflight := NewMockPreflightValidator()

		mockRepo.On("FindByID", mock.Anything, session.ID()).
			Return(session, nil)
		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*installation.InstallationSession")).
			Return(nil)

		// Return a conflict
		mockConflictResolver.On("DetectConflicts", mock.Anything, mock.Anything).
			Return([]installation.PackageConflict{conflict}, nil)

		// Conflict should be resolved with Remove action
		mockConflictResolver.On("ResolveConflict", mock.Anything, conflict, installation.ActionRemove).
			Return(nil)

		mockProgressEstimator.On("CalculatePhaseProgress", mock.Anything, mock.Anything, mock.Anything).
			Return(50)
		mockProgressEstimator.On("EstimateRemainingTime", mock.Anything, mock.Anything, mock.Anything).
			Return(5 * time.Minute)

		mockPkgManager.On("InstallPackage", mock.Anything, "hyprland", "0.35.0").
			Return(nil)

		mockPreflight.On("Run", mock.Anything).Return(nil)

		// Execute use case
		useCase := usecases.NewExecuteInstallationUseCase(
			mockRepo,
			mockConflictResolver,
			mockProgressEstimator,
			mockConfigMerger,
			mockPkgManager,
			nil, // historyRecorder not needed for this test
			mockPreflight,
		)
		ctx := context.Background()

		response, err := useCase.Execute(ctx, session.ID(), nil)

		require.NoError(t, err)
		assert.Equal(t, session.ID(), response.SessionID)
		mockConflictResolver.AssertExpectations(t)
	})

	t.Run("returns error for non-existent session", func(t *testing.T) {
		mockRepo := new(MockInstallationSessionRepository)
		mockConflictResolver := new(MockConflictResolver)
		mockProgressEstimator := new(MockProgressEstimator)
		mockConfigMerger := new(MockConfigurationMerger)
		mockPkgManager := new(MockPackageManager)
		mockPreflight := NewMockPreflightValidator()

		mockRepo.On("FindByID", mock.Anything, "nonexistent-id").
			Return(nil, installation.ErrSessionNotFound)

		useCase := usecases.NewExecuteInstallationUseCase(
			mockRepo,
			mockConflictResolver,
			mockProgressEstimator,
			mockConfigMerger,
			mockPkgManager,
			nil, // historyRecorder not needed for this test
			mockPreflight,
		)
		ctx := context.Background()

		_, err := useCase.Execute(ctx, "nonexistent-id", nil)

		assert.Error(t, err)
		assert.ErrorIs(t, err, installation.ErrSessionNotFound)
	})

	t.Run("handles installation errors and marks session as failed", func(t *testing.T) {
		components, err := createTestComponents()
		require.NoError(t, err)

		diskSpace, err := installation.NewDiskSpace(
			100*uint64(installation.GB),
			10*uint64(installation.GB),
		)
		require.NoError(t, err)

		config, err := installation.NewInstallationConfiguration(
			components,
			nil,
			diskSpace,
			false,
		)
		require.NoError(t, err)

		session, err := installation.NewInstallationSession(config)
		require.NoError(t, err)

		mockRepo := new(MockInstallationSessionRepository)
		mockConflictResolver := new(MockConflictResolver)
		mockProgressEstimator := new(MockProgressEstimator)
		mockConfigMerger := new(MockConfigurationMerger)
		mockPkgManager := new(MockPackageManager)
		mockPreflight := NewMockPreflightValidator()

		mockRepo.On("FindByID", mock.Anything, session.ID()).
			Return(session, nil)
		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*installation.InstallationSession")).
			Return(nil)

		mockConflictResolver.On("DetectConflicts", mock.Anything, mock.Anything).
			Return([]installation.PackageConflict{}, nil)

		mockProgressEstimator.On("CalculatePhaseProgress", mock.Anything, mock.Anything, mock.Anything).
			Return(50)
		mockProgressEstimator.On("EstimateRemainingTime", mock.Anything, mock.Anything, mock.Anything).
			Return(5 * time.Minute)

		// Package installation fails
		installErr := assert.AnError
		mockPkgManager.On("InstallPackage", mock.Anything, "hyprland", "0.35.0").
			Return(installErr)

		mockPreflight.On("Run", mock.Anything).Return(nil)

		useCase := usecases.NewExecuteInstallationUseCase(
			mockRepo,
			mockConflictResolver,
			mockProgressEstimator,
			mockConfigMerger,
			mockPkgManager,
			nil, // historyRecorder not needed for this test
			mockPreflight,
		)
		ctx := context.Background()

		response, err := useCase.Execute(ctx, session.ID(), nil)

		require.NoError(t, err) // Use case doesn't error, but marks session as failed
		assert.Equal(t, "failed", response.Status)
		mockPkgManager.AssertExpectations(t)
	})

	t.Run("blocks installation when preflight checks fail", func(t *testing.T) {
		// Create a valid installation session
		components, err := createTestComponents()
		require.NoError(t, err)

		diskSpace, err := installation.NewDiskSpace(
			100*uint64(installation.GB),
			10*uint64(installation.GB),
		)
		require.NoError(t, err)

		config, err := installation.NewInstallationConfiguration(
			components,
			nil,
			diskSpace,
			false,
		)
		require.NoError(t, err)

		session, err := installation.NewInstallationSession(config)
		require.NoError(t, err)

		// Setup mocks
		mockRepo := new(MockInstallationSessionRepository)
		mockConflictResolver := new(MockConflictResolver)
		mockProgressEstimator := new(MockProgressEstimator)
		mockConfigMerger := new(MockConfigurationMerger)
		mockPkgManager := new(MockPackageManager)

		// Create a mock preflight that will fail with a blocker
		mockPreflight := &MockPreflightValidator{
			Mock:         mock.Mock{},
			progressChan: make(chan preflightTUI.ProgressUpdate),
			session:      createFailedPreflightSession(),
		}

		mockRepo.On("FindByID", mock.Anything, session.ID()).
			Return(session, nil)
		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*installation.InstallationSession")).
			Return(nil)

		mockPreflight.On("Run", mock.Anything).Return(nil)

		// Execute use case
		useCase := usecases.NewExecuteInstallationUseCase(
			mockRepo,
			mockConflictResolver,
			mockProgressEstimator,
			mockConfigMerger,
			mockPkgManager,
			nil, // historyRecorder not needed for this test
			mockPreflight,
		)
		ctx := context.Background()

		response, err := useCase.Execute(ctx, session.ID(), nil)

		// Should return an error because preflight failed
		require.Error(t, err)
		assert.Contains(t, err.Error(), "preflight checks failed")
		assert.NotNil(t, response)
		assert.Equal(t, "failed", response.Status)
		assert.Equal(t, "Preflight Checks", response.CurrentPhase)
	})
}

// createFailedPreflightSession creates a preflight session with a blocking failure
func createFailedPreflightSession() *preflight.ValidationSession {
	session := preflight.NewValidationSession()

	// Create guidance for fixing the issue
	guidance := preflight.NewUserGuidance(
		"Debian version check failed",
		"Debian noble is not supported. Gohan requires Debian Sid or Trixie.",
		[]string{"Upgrade to Debian Sid or Trixie"},
		"https://www.debian.org/releases/",
	)

	// Create a blocking validation result (critical severity)
	result := preflight.NewValidationResult(
		preflight.RequirementDebianVersion,
		preflight.StatusFail,
		preflight.SeverityCritical,
		"noble",
		"sid or trixie",
		guidance,
	)

	// Add the result to the session
	session.AddResult(result)
	session.Complete()

	return session
}

// Helper function to create test components
func createTestComponents() ([]installation.ComponentSelection, error) {
	pkg, err := installation.NewPackageInfo("hyprland", "0.35.0", 50*uint64(installation.MB), nil)
	if err != nil {
		return nil, err
	}

	component, err := installation.NewComponentSelection(
		installation.ComponentHyprland,
		"0.35.0",
		&pkg,
	)
	if err != nil {
		return nil, err
	}

	return []installation.ComponentSelection{component}, nil
}
