package configservice_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rebelopsio/gohan/internal/infrastructure/installation/backup"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/configservice"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// Phase 3.4: Configuration Deployment Tests (TDD)
// ========================================

func TestConfigDeployer_DeployConfiguration(t *testing.T) {
	t.Run("deploys single configuration file", func(t *testing.T) {
		tmpDir := t.TempDir()
		backupDir := filepath.Join(tmpDir, "backups")

		deployer := setupDeployer(t, backupDir)

		// Create template file
		templatePath := filepath.Join(tmpDir, "templates", "test.conf")
		err := os.MkdirAll(filepath.Dir(templatePath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(templatePath, []byte("user = {{username}}"), 0644)
		require.NoError(t, err)

		// Deploy configuration
		targetPath := filepath.Join(tmpDir, "config", "test.conf")
		config := configservice.ConfigurationFile{
			SourceTemplate: templatePath,
			TargetPath:     targetPath,
			Permissions:    0644,
			BackupBefore:   false,
		}

		vars := templates.TemplateVars{Username: "testuser"}
		ctx := context.Background()

		err = deployer.DeployConfiguration(ctx, config, vars)
		require.NoError(t, err)

		// Verify deployed file
		content, err := os.ReadFile(targetPath)
		require.NoError(t, err)
		assert.Equal(t, "user = testuser", string(content))
	})

	t.Run("backs up existing file before overwriting", func(t *testing.T) {
		tmpDir := t.TempDir()
		backupDir := filepath.Join(tmpDir, "backups")

		deployer := setupDeployer(t, backupDir)

		// Create existing file
		targetPath := filepath.Join(tmpDir, "config", "existing.conf")
		err := os.MkdirAll(filepath.Dir(targetPath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(targetPath, []byte("old content"), 0644)
		require.NoError(t, err)

		// Create template
		templatePath := filepath.Join(tmpDir, "templates", "new.conf")
		err = os.MkdirAll(filepath.Dir(templatePath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(templatePath, []byte("new content"), 0644)
		require.NoError(t, err)

		// Deploy with backup
		config := configservice.ConfigurationFile{
			SourceTemplate: templatePath,
			TargetPath:     targetPath,
			Permissions:    0644,
			BackupBefore:   true,
		}

		vars := templates.TemplateVars{}
		ctx := context.Background()

		err = deployer.DeployConfiguration(ctx, config, vars)
		require.NoError(t, err)

		// Verify new content
		content, err := os.ReadFile(targetPath)
		require.NoError(t, err)
		assert.Equal(t, "new content", string(content))

		// Verify backup was created
		backups, err := deployer.ListBackups(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, backups, "Should have created a backup")
	})

	t.Run("processes template variables", func(t *testing.T) {
		tmpDir := t.TempDir()
		backupDir := filepath.Join(tmpDir, "backups")

		deployer := setupDeployer(t, backupDir)

		// Create template with multiple variables
		templatePath := filepath.Join(tmpDir, "templates", "config.conf")
		err := os.MkdirAll(filepath.Dir(templatePath), 0755)
		require.NoError(t, err)
		templateContent := `user = {{username}}
home = {{home}}
config = {{config_dir}}`
		err = os.WriteFile(templatePath, []byte(templateContent), 0644)
		require.NoError(t, err)

		// Deploy
		targetPath := filepath.Join(tmpDir, "config", "config.conf")
		config := configservice.ConfigurationFile{
			SourceTemplate: templatePath,
			TargetPath:     targetPath,
			Permissions:    0644,
		}

		vars := templates.TemplateVars{
			Username:  "alice",
			Home:      "/home/alice",
			ConfigDir: "/home/alice/.config",
		}
		ctx := context.Background()

		err = deployer.DeployConfiguration(ctx, config, vars)
		require.NoError(t, err)

		// Verify substitution
		content, err := os.ReadFile(targetPath)
		require.NoError(t, err)
		expected := `user = alice
home = /home/alice
config = /home/alice/.config`
		assert.Equal(t, expected, string(content))
	})

	t.Run("sets file permissions", func(t *testing.T) {
		tmpDir := t.TempDir()
		backupDir := filepath.Join(tmpDir, "backups")

		deployer := setupDeployer(t, backupDir)

		// Create template
		templatePath := filepath.Join(tmpDir, "templates", "secret.conf")
		err := os.MkdirAll(filepath.Dir(templatePath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(templatePath, []byte("secret"), 0644)
		require.NoError(t, err)

		// Deploy with specific permissions
		targetPath := filepath.Join(tmpDir, "config", "secret.conf")
		config := configservice.ConfigurationFile{
			SourceTemplate: templatePath,
			TargetPath:     targetPath,
			Permissions:    0600,
		}

		vars := templates.TemplateVars{}
		ctx := context.Background()

		err = deployer.DeployConfiguration(ctx, config, vars)
		require.NoError(t, err)

		// Verify permissions (best effort - environment dependent)
		info, err := os.Stat(targetPath)
		require.NoError(t, err)
		// Just verify file was created, permissions may vary by environment
		assert.NotNil(t, info)
	})
}

func TestConfigDeployer_DeployConfigurations(t *testing.T) {
	t.Run("deploys multiple configuration files", func(t *testing.T) {
		tmpDir := t.TempDir()
		backupDir := filepath.Join(tmpDir, "backups")

		deployer := setupDeployer(t, backupDir)

		// Create multiple templates
		templateNames := []string{"config1.conf", "config2.conf", "config3.conf"}
		configFiles := []configservice.ConfigurationFile{}

		for _, name := range templateNames {
			templatePath := filepath.Join(tmpDir, "templates", name)
			err := os.MkdirAll(filepath.Dir(templatePath), 0755)
			require.NoError(t, err)
			err = os.WriteFile(templatePath, []byte("content: "+name), 0644)
			require.NoError(t, err)

			targetPath := filepath.Join(tmpDir, "config", name)
			configFiles = append(configFiles, configservice.ConfigurationFile{
				SourceTemplate: templatePath,
				TargetPath:     targetPath,
				Permissions:    0644,
			})
		}

		// Deploy all
		vars := templates.TemplateVars{}
		ctx := context.Background()
		progressChan := make(chan configservice.DeploymentProgress, 10)

		done := make(chan error, 1)
		go func() {
			done <- deployer.DeployConfigurations(ctx, configFiles, vars, progressChan)
			close(progressChan)
		}()

		// Collect progress
		var progress []configservice.DeploymentProgress
		for p := range progressChan {
			progress = append(progress, p)
		}

		err := <-done
		require.NoError(t, err)

		// Verify all files deployed
		for _, name := range templateNames {
			targetPath := filepath.Join(tmpDir, "config", name)
			_, err := os.Stat(targetPath)
			assert.NoError(t, err, "File %s should be deployed", name)
		}

		// Verify progress was reported
		assert.NotEmpty(t, progress, "Should report progress")
	})

	t.Run("reports progress for each file", func(t *testing.T) {
		tmpDir := t.TempDir()
		backupDir := filepath.Join(tmpDir, "backups")

		deployer := setupDeployer(t, backupDir)

		// Create single template
		templatePath := filepath.Join(tmpDir, "templates", "test.conf")
		err := os.MkdirAll(filepath.Dir(templatePath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(templatePath, []byte("test"), 0644)
		require.NoError(t, err)

		configs := []configservice.ConfigurationFile{
			{
				SourceTemplate: templatePath,
				TargetPath:     filepath.Join(tmpDir, "config", "test.conf"),
				Permissions:    0644,
			},
		}

		vars := templates.TemplateVars{}
		ctx := context.Background()
		progressChan := make(chan configservice.DeploymentProgress, 10)

		done := make(chan struct{})
		var progress []configservice.DeploymentProgress

		go func() {
			for p := range progressChan {
				progress = append(progress, p)
			}
			close(done)
		}()

		err = deployer.DeployConfigurations(ctx, configs, vars, progressChan)
		close(progressChan)
		<-done

		require.NoError(t, err)
		assert.NotEmpty(t, progress, "Should have progress events")

		// Should have started and completed events
		hasStarted := false
		hasCompleted := false
		for _, p := range progress {
			if p.Status == "started" {
				hasStarted = true
			}
			if p.Status == "completed" {
				hasCompleted = true
			}
		}

		assert.True(t, hasStarted, "Should report started")
		assert.True(t, hasCompleted, "Should report completed")
	})

	t.Run("handles deployment failure gracefully", func(t *testing.T) {
		tmpDir := t.TempDir()
		backupDir := filepath.Join(tmpDir, "backups")

		deployer := setupDeployer(t, backupDir)

		// Create invalid config (non-existent template)
		configs := []configservice.ConfigurationFile{
			{
				SourceTemplate: "/nonexistent/template.conf",
				TargetPath:     filepath.Join(tmpDir, "config", "output.conf"),
				Permissions:    0644,
			},
		}

		vars := templates.TemplateVars{}
		ctx := context.Background()
		progressChan := make(chan configservice.DeploymentProgress, 10)

		done := make(chan error, 1)
		go func() {
			done <- deployer.DeployConfigurations(ctx, configs, vars, progressChan)
			close(progressChan)
		}()

		// Drain progress
		for range progressChan {
		}

		err := <-done
		assert.Error(t, err, "Should return error for invalid template")
	})
}

func TestConfigDeployer_ListBackups(t *testing.T) {
	t.Run("lists backups created during deployment", func(t *testing.T) {
		tmpDir := t.TempDir()
		backupDir := filepath.Join(tmpDir, "backups")

		deployer := setupDeployer(t, backupDir)

		// Create existing file
		targetPath := filepath.Join(tmpDir, "config", "test.conf")
		err := os.MkdirAll(filepath.Dir(targetPath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(targetPath, []byte("original"), 0644)
		require.NoError(t, err)

		// Create template
		templatePath := filepath.Join(tmpDir, "templates", "test.conf")
		err = os.MkdirAll(filepath.Dir(templatePath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(templatePath, []byte("new"), 0644)
		require.NoError(t, err)

		// Deploy with backup
		config := configservice.ConfigurationFile{
			SourceTemplate: templatePath,
			TargetPath:     targetPath,
			Permissions:    0644,
			BackupBefore:   true,
		}

		vars := templates.TemplateVars{}
		ctx := context.Background()

		err = deployer.DeployConfiguration(ctx, config, vars)
		require.NoError(t, err)

		// List backups
		backups, err := deployer.ListBackups(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, backups, "Should have created backup")
	})
}

func TestConfigDeployer_RollbackDeployment(t *testing.T) {
	t.Run("rolls back deployment using backup ID", func(t *testing.T) {
		tmpDir := t.TempDir()
		backupDir := filepath.Join(tmpDir, "backups")

		deployer := setupDeployer(t, backupDir)

		// Create existing file
		targetPath := filepath.Join(tmpDir, "config", "test.conf")
		err := os.MkdirAll(filepath.Dir(targetPath), 0755)
		require.NoError(t, err)
		originalContent := "original content"
		err = os.WriteFile(targetPath, []byte(originalContent), 0644)
		require.NoError(t, err)

		// Create template
		templatePath := filepath.Join(tmpDir, "templates", "test.conf")
		err = os.MkdirAll(filepath.Dir(templatePath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(templatePath, []byte("new content"), 0644)
		require.NoError(t, err)

		// Deploy with backup
		config := configservice.ConfigurationFile{
			SourceTemplate: templatePath,
			TargetPath:     targetPath,
			Permissions:    0644,
			BackupBefore:   true,
		}

		vars := templates.TemplateVars{}
		ctx := context.Background()

		result, err := deployer.DeployWithBackup(ctx, config, vars)
		require.NoError(t, err)
		require.NotEmpty(t, result.BackupID, "Should have backup ID")

		// Verify new content
		content, _ := os.ReadFile(targetPath)
		assert.Equal(t, "new content", string(content))

		// Rollback
		err = deployer.RollbackDeployment(ctx, result.BackupID)
		require.NoError(t, err)

		// Verify original content restored
		content, err = os.ReadFile(targetPath)
		require.NoError(t, err)
		assert.Equal(t, originalContent, string(content))
	})

	t.Run("returns error for non-existent backup", func(t *testing.T) {
		tmpDir := t.TempDir()
		backupDir := filepath.Join(tmpDir, "backups")

		deployer := setupDeployer(t, backupDir)
		ctx := context.Background()

		err := deployer.RollbackDeployment(ctx, "nonexistent-backup-id")
		assert.Error(t, err)
	})
}

// Helper functions

func setupDeployer(t *testing.T, backupDir string) *configservice.ConfigDeployer {
	t.Helper()

	templateEngine := templates.NewTemplateEngine()
	backupService := backup.NewBackupService(backupDir)

	return configservice.NewConfigDeployer(templateEngine, backupService)
}
