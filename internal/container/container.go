package container

import (
	"fmt"
	"os"
	"path/filepath"

	historyServices "github.com/rebelopsio/gohan/internal/application/history/services"
	"github.com/rebelopsio/gohan/internal/application/installation/usecases"
	"github.com/rebelopsio/gohan/internal/config"
	"github.com/rebelopsio/gohan/internal/domain/installation"
	historyRepo "github.com/rebelopsio/gohan/internal/infrastructure/history/repository"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/backup"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/configservice"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/packagemanager"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/repository"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/services"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/templates"
	preflightTUI "github.com/rebelopsio/gohan/internal/tui/preflight"
)

// Container holds all dependencies for the application
type Container struct {
	Config *config.Config

	// Repositories
	HistoryRepo      *historyRepo.SQLiteRepository
	InstallationRepo installation.InstallationSessionRepository

	// Services
	HistoryQueryService     *historyServices.HistoryQueryService
	HistoryRecordingService *historyServices.HistoryRecordingService
	ProgressEstimator       *services.ProgressEstimator
	ConfigMerger            *services.ConfigurationMerger
	PackageManager          *packagemanager.APTManager
	ConfigDeployer          *configservice.ConfigDeployer

	// Use Cases
	StartInstallationUseCase   *usecases.StartInstallationUseCase
	ExecuteInstallationUseCase *usecases.ExecuteInstallationUseCase
	GetStatusUseCase           *usecases.GetInstallationStatusUseCase
	ListInstallationsUseCase   *usecases.ListInstallationsUseCase
	CancelInstallationUseCase  *usecases.CancelInstallationUseCase
}

// New creates a new dependency container
func New() (*Container, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	c := &Container{
		Config: cfg,
	}

	// Initialize repositories
	if err := c.initRepositories(); err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	// Initialize services
	c.initServices()

	// Initialize use cases
	c.initUseCases()

	return c, nil
}

// initRepositories initializes all repositories
func (c *Container) initRepositories() error {
	// History repository
	historyRepo, err := historyRepo.NewSQLiteRepository(c.Config.Database.HistoryDB)
	if err != nil {
		return fmt.Errorf("failed to create history repository: %w", err)
	}
	c.HistoryRepo = historyRepo

	// Installation repository
	// Note: Using memory repository for now as SQLite reconstruction not fully implemented
	c.InstallationRepo = repository.NewMemorySessionRepository()

	// TODO: Switch to SQLite when reconstruction is complete
	// installationRepo, err := repository.NewSQLiteSessionRepository(c.Config.Database.InstallationDB)
	// if err != nil {
	// 	return fmt.Errorf("failed to create installation repository: %w", err)
	// }
	// c.InstallationRepo = installationRepo

	return nil
}

// initServices initializes all application services
func (c *Container) initServices() {
	// History services
	c.HistoryQueryService = historyServices.NewHistoryQueryService(c.HistoryRepo)
	c.HistoryRecordingService = historyServices.NewHistoryRecordingService(c.HistoryRepo)

	// Installation services
	c.ProgressEstimator = services.NewProgressEstimator()
	c.ConfigMerger = services.NewConfigurationMerger()

	// Choose package manager based on dry-run setting
	if c.Config.Installation.DryRun {
		c.PackageManager = packagemanager.NewAPTManagerDryRun()
	} else {
		c.PackageManager = packagemanager.NewAPTManager()
	}

	// Configuration deployment services
	homeDir, _ := os.UserHomeDir()
	backupDir := filepath.Join(homeDir, ".config", "gohan", "backups")
	templateEngine := templates.NewTemplateEngine()
	backupService := backup.NewBackupService(backupDir)
	c.ConfigDeployer = configservice.NewConfigDeployer(templateEngine, backupService)
}

// initUseCases initializes all use cases
func (c *Container) initUseCases() {
	c.StartInstallationUseCase = usecases.NewStartInstallationUseCase(c.InstallationRepo)

	c.ExecuteInstallationUseCase = usecases.NewExecuteInstallationUseCase(
		c.InstallationRepo,
		c.PackageManager, // ConflictResolver
		c.ProgressEstimator,
		c.ConfigMerger,
		c.PackageManager, // PackageManager
		c.HistoryRecordingService, // HistoryRecorder
		preflightTUI.NewValidationRunner(), // PreflightValidator
		c.ConfigDeployer,
	)

	c.GetStatusUseCase = usecases.NewGetInstallationStatusUseCase(c.InstallationRepo)
	c.ListInstallationsUseCase = usecases.NewListInstallationsUseCase(c.InstallationRepo)
	c.CancelInstallationUseCase = usecases.NewCancelInstallationUseCase(c.InstallationRepo)
}

// Close closes all resources
func (c *Container) Close() error {
	var errs []error

	if c.HistoryRepo != nil {
		if err := c.HistoryRepo.Close(); err != nil {
			errs = append(errs, fmt.Errorf("history repo: %w", err))
		}
	}

	// Close installation repo if it implements io.Closer
	if closer, ok := c.InstallationRepo.(interface{ Close() error }); ok && closer != nil {
		if err := closer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("installation repo: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %v", errs)
	}

	return nil
}
