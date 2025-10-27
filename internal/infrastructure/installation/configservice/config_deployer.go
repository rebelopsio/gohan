package configservice

import (
	"context"
	"fmt"
	"os"

	"github.com/rebelopsio/gohan/internal/infrastructure/installation/backup"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/templates"
)

// ConfigDeployer handles deployment of configuration files
type ConfigDeployer struct {
	templateEngine *templates.TemplateEngine
	backupService  *backup.BackupService
}

// ConfigurationFile represents a configuration file to deploy
type ConfigurationFile struct {
	SourceTemplate string      // Path to template file
	TargetPath     string      // Where to deploy
	Permissions    os.FileMode // File permissions
	BackupBefore   bool        // Whether to backup before overwriting
}

// DeploymentProgress represents progress for a single configuration deployment
type DeploymentProgress struct {
	FilePath        string
	Status          string // "started", "processing", "completed", "failed"
	PercentComplete float64
	Error           error
}

// DeploymentResult contains the result of a deployment operation
type DeploymentResult struct {
	FilePath   string
	Success    bool
	BackupID   string // ID of backup if created
	BackupPath string // Path to backup if created
	Error      error
}

// NewConfigDeployer creates a new configuration deployer
func NewConfigDeployer(templateEngine *templates.TemplateEngine, backupService *backup.BackupService) *ConfigDeployer {
	return &ConfigDeployer{
		templateEngine: templateEngine,
		backupService:  backupService,
	}
}

// DeployConfiguration deploys a single configuration file
func (cd *ConfigDeployer) DeployConfiguration(ctx context.Context, config ConfigurationFile, vars templates.TemplateVars) error {
	// Backup if requested and file exists
	if config.BackupBefore {
		if _, err := os.Stat(config.TargetPath); err == nil {
			// File exists, back it up
			_, err := cd.backupService.CreateBackup(
				ctx,
				[]string{config.TargetPath},
				fmt.Sprintf("Backup before deploying %s", config.TargetPath),
			)
			if err != nil {
				return fmt.Errorf("failed to backup existing file: %w", err)
			}
		}
	}

	// Process template and deploy
	if err := cd.templateEngine.ProcessFile(config.SourceTemplate, config.TargetPath, vars); err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	// Set permissions
	if err := os.Chmod(config.TargetPath, config.Permissions); err != nil {
		// Non-fatal - permissions are best effort
		// Log but continue
	}

	return nil
}

// DeployConfigurations deploys multiple configuration files with progress reporting
func (cd *ConfigDeployer) DeployConfigurations(
	ctx context.Context,
	configs []ConfigurationFile,
	vars templates.TemplateVars,
	progressChan chan<- DeploymentProgress) error {

	if len(configs) == 0 {
		return nil
	}

	totalFiles := len(configs)

	for i, config := range configs {
		// Check context
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled: %w", err)
		}

		percentComplete := float64(i) / float64(totalFiles) * 100

		// Report started
		if progressChan != nil {
			progressChan <- DeploymentProgress{
				FilePath:        config.TargetPath,
				Status:          "started",
				PercentComplete: percentComplete,
			}
		}

		// Report processing
		if progressChan != nil {
			progressChan <- DeploymentProgress{
				FilePath:        config.TargetPath,
				Status:          "processing",
				PercentComplete: percentComplete + (50.0 / float64(totalFiles)),
			}
		}

		// Deploy the file
		err := cd.DeployConfiguration(ctx, config, vars)
		if err != nil {
			// Report failure
			if progressChan != nil {
				progressChan <- DeploymentProgress{
					FilePath:        config.TargetPath,
					Status:          "failed",
					PercentComplete: percentComplete,
					Error:           err,
				}
			}
			return fmt.Errorf("failed to deploy %s: %w", config.TargetPath, err)
		}

		// Report completed
		if progressChan != nil {
			progressChan <- DeploymentProgress{
				FilePath:        config.TargetPath,
				Status:          "completed",
				PercentComplete: float64(i+1) / float64(totalFiles) * 100,
			}
		}
	}

	return nil
}

// DeployWithBackup deploys a configuration and returns backup information
func (cd *ConfigDeployer) DeployWithBackup(
	ctx context.Context,
	config ConfigurationFile,
	vars templates.TemplateVars) (*DeploymentResult, error) {

	result := &DeploymentResult{
		FilePath: config.TargetPath,
	}

	// Check if target exists
	if _, err := os.Stat(config.TargetPath); err == nil {
		// Create backup
		metadata, err := cd.backupService.CreateBackup(
			ctx,
			[]string{config.TargetPath},
			fmt.Sprintf("Deployment backup for %s", config.TargetPath),
		)
		if err != nil {
			result.Error = fmt.Errorf("failed to create backup: %w", err)
			return result, result.Error
		}

		result.BackupID = metadata.ID
		result.BackupPath = metadata.Path
	}

	// Deploy
	if err := cd.DeployConfiguration(ctx, config, vars); err != nil {
		result.Error = err
		return result, err
	}

	result.Success = true
	return result, nil
}

// ListBackups lists all available backups
func (cd *ConfigDeployer) ListBackups(ctx context.Context) ([]*backup.BackupMetadata, error) {
	return cd.backupService.ListBackups(ctx)
}

// RollbackDeployment rolls back a deployment using the backup ID
func (cd *ConfigDeployer) RollbackDeployment(ctx context.Context, backupID string) error {
	return cd.backupService.RestoreBackup(ctx, backupID)
}

// GetBackupInfo retrieves information about a specific backup
func (cd *ConfigDeployer) GetBackupInfo(ctx context.Context, backupID string) (*backup.BackupMetadata, error) {
	return cd.backupService.GetBackupInfo(ctx, backupID)
}
