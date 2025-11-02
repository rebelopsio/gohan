//go:build integration
// +build integration

package integration

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rebelopsio/gohan/internal/application/configuration"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/backup"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/configservice"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTemplateFiles copies actual template files into the test working directory
func setupTemplateFiles(t *testing.T) {
	t.Helper()

	// Get project root (assuming tests/integration is 2 levels deep)
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	require.NoError(t, err)

	sourceTemplateDir := filepath.Join(projectRoot, "templates")

	// Template files to copy
	templateFiles := []struct {
		src string
		dst string
	}{
		{
			src: filepath.Join(sourceTemplateDir, "hyprland", "hyprland.conf.tmpl"),
			dst: filepath.Join("templates", "hyprland", "hyprland.conf.tmpl"),
		},
		{
			src: filepath.Join(sourceTemplateDir, "waybar", "config.jsonc"),
			dst: filepath.Join("templates", "waybar", "config.jsonc"),
		},
		{
			src: filepath.Join(sourceTemplateDir, "waybar", "style.css.tmpl"),
			dst: filepath.Join("templates", "waybar", "style.css.tmpl"),
		},
		{
			src: filepath.Join(sourceTemplateDir, "kitty", "kitty.conf.tmpl"),
			dst: filepath.Join("templates", "kitty", "kitty.conf.tmpl"),
		},
		{
			src: filepath.Join(sourceTemplateDir, "fuzzel", "fuzzel.ini.tmpl"),
			dst: filepath.Join("templates", "fuzzel", "fuzzel.ini.tmpl"),
		},
	}

	// Copy each template file
	for _, tf := range templateFiles {
		// Ensure destination directory exists
		err := os.MkdirAll(filepath.Dir(tf.dst), 0755)
		require.NoError(t, err, "Failed to create directory for %s", tf.dst)

		// Copy file
		src, err := os.Open(tf.src)
		require.NoError(t, err, "Failed to open source file %s", tf.src)
		defer src.Close()

		dst, err := os.Create(tf.dst)
		require.NoError(t, err, "Failed to create destination file %s", tf.dst)
		defer dst.Close()

		_, err = io.Copy(dst, src)
		require.NoError(t, err, "Failed to copy %s to %s", tf.src, tf.dst)
	}
}

// TestConfigurationDeployment_FreshSystem corresponds to:
// Feature: Configuration Deployment
// Scenario: Deploy configurations to fresh system
func TestConfigurationDeployment_FreshSystem(t *testing.T) {
	setupTemplateFiles(t)
	tmpDir := t.TempDir()
	ctx := context.Background()

	// Setup services
	templateEngine := templates.NewTemplateEngine()
	backupRoot := filepath.Join(tmpDir, "backups")
	backupService := backup.NewBackupService(backupRoot)
	deployer := configservice.NewConfigDeployer(templateEngine, backupService)
	useCase := configuration.NewConfigDeployUseCase(deployer, templateEngine)

	// Given I have no existing Hyprland configurations
	configDir := filepath.Join(tmpDir, ".config")
	_, err := os.Stat(configDir)
	require.True(t, os.IsNotExist(err), "Config directory should not exist yet")

	// When the configuration deployment starts (deploy all components)
	request := configuration.DeployConfigRequest{
		Components: []string{"hyprland", "waybar", "kitty", "fuzzel"},
		CustomVars: map[string]string{
			"username": "testuser",
			"home":     tmpDir,
			"home_dir": tmpDir,
		},
		DryRun: false,
	}

	resp, err := useCase.Execute(ctx, request)

	// Then all configurations should be deployed successfully
	require.NoError(t, err)
	assert.Equal(t, 5, resp.TotalFiles) // hyprland(1) + waybar(2) + kitty(1) + fuzzel(1)
	assert.Equal(t, 5, resp.SuccessfulFiles)
	assert.Equal(t, 0, resp.FailedFiles)

	// And Hyprland configurations should be created in ~/.config/hypr/
	hyprlandConf := filepath.Join(tmpDir, ".config", "hypr", "hyprland.conf")
	assert.FileExists(t, hyprlandConf)

	// And Waybar configurations should be created in ~/.config/waybar/
	waybarConfig := filepath.Join(tmpDir, ".config", "waybar", "config.jsonc")
	waybarStyle := filepath.Join(tmpDir, ".config", "waybar", "style.css")
	assert.FileExists(t, waybarConfig)
	assert.FileExists(t, waybarStyle)

	// And Kitty configurations should be created in ~/.config/kitty/
	kittyConf := filepath.Join(tmpDir, ".config", "kitty", "kitty.conf")
	assert.FileExists(t, kittyConf)

	// And Fuzzel configurations should be created in ~/.config/fuzzel/
	fuzzelIni := filepath.Join(tmpDir, ".config", "fuzzel", "fuzzel.ini")
	assert.FileExists(t, fuzzelIni)

	// And all configuration files should have correct permissions
	for _, file := range []string{hyprlandConf, waybarConfig, waybarStyle, kittyConf, fuzzelIni} {
		info, err := os.Stat(file)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0644), info.Mode().Perm())
	}
}

// TestConfigurationDeployment_BackupExisting corresponds to:
// Scenario: Backup existing configurations before overwriting
func TestConfigurationDeployment_BackupExisting(t *testing.T) {
	setupTemplateFiles(t)
	tmpDir := t.TempDir()
	ctx := context.Background()

	// Setup services
	templateEngine := templates.NewTemplateEngine()
	backupRoot := filepath.Join(tmpDir, "backups")
	backupService := backup.NewBackupService(backupRoot)
	deployer := configservice.NewConfigDeployer(templateEngine, backupService)
	useCase := configuration.NewConfigDeployUseCase(deployer, templateEngine)

	// Given I have existing Hyprland configurations in ~/.config/hypr/
	configDir := filepath.Join(tmpDir, ".config", "hypr")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	existingConfig := filepath.Join(configDir, "hyprland.conf")
	originalContent := "# My existing config\nmonitor = eDP-1,1920x1080,auto,1"
	err = os.WriteFile(existingConfig, []byte(originalContent), 0644)
	require.NoError(t, err)

	// When the configuration deployment starts
	request := configuration.DeployConfigRequest{
		Components: []string{"hyprland"},
		CustomVars: map[string]string{
			"username": "testuser",
			"home":     tmpDir,
		},
		DryRun: false,
	}

	resp, err := useCase.Execute(ctx, request)

	// Then the system should create a timestamped backup
	require.NoError(t, err)
	assert.NotEmpty(t, resp.BackupID, "Backup ID should be set")
	assert.NotEmpty(t, resp.BackupPath, "Backup path should be set")

	// And the backup should contain all existing configurations
	backups, err := backupService.ListBackups(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, backups, "Should have created a backup")

	// And original files should be preserved in backup
	backupInfo := backups[0] // Most recent backup
	assert.NotEmpty(t, backupInfo.Files)
	assert.Contains(t, backupInfo.Files[0].OriginalPath, "hyprland.conf")

	// And then the system should deploy new configurations
	assert.Equal(t, 1, resp.SuccessfulFiles)

	// Verify new config was deployed
	newContent, err := os.ReadFile(existingConfig)
	require.NoError(t, err)
	assert.NotEqual(t, originalContent, string(newContent), "Config should be updated")
	assert.Contains(t, string(newContent), "Gohan", "Should contain Gohan header")
}

// TestConfigurationDeployment_PersonalizeForSystem corresponds to:
// Scenario: Personalize configurations for my system
func TestConfigurationDeployment_PersonalizeForSystem(t *testing.T) {
	setupTemplateFiles(t)
	tmpDir := t.TempDir()
	ctx := context.Background()

	// Setup services
	templateEngine := templates.NewTemplateEngine()
	backupRoot := filepath.Join(tmpDir, "backups")
	backupService := backup.NewBackupService(backupRoot)
	deployer := configservice.NewConfigDeployer(templateEngine, backupService)
	useCase := configuration.NewConfigDeployUseCase(deployer, templateEngine)

	// Given I am logged in as "alice"
	// And my home directory is set to tmpDir for testing
	request := configuration.DeployConfigRequest{
		Components: []string{"hyprland", "kitty"},
		CustomVars: map[string]string{
			"username": "alice",
			"home":     tmpDir,
			"home_dir": tmpDir,
		},
		DryRun: false,
	}

	// When configurations are deployed
	resp, err := useCase.Execute(ctx, request)

	// Then configurations should reference my username and paths
	require.NoError(t, err)
	assert.Equal(t, 2, resp.SuccessfulFiles)

	// Verify hyprland config has personalized values
	hyprlandConf := filepath.Join(tmpDir, ".config", "hypr", "hyprland.conf")
	content, err := os.ReadFile(hyprlandConf)
	require.NoError(t, err)

	configStr := string(content)
	assert.Contains(t, configStr, "alice", "Should contain username")
	assert.NotContains(t, configStr, "{{username}}", "Should not contain template placeholders")
	assert.NotContains(t, configStr, "{{home}}", "Should not contain template placeholders")

	// Verify kitty config has theme colors substituted
	kittyConf := filepath.Join(tmpDir, ".config", "kitty", "kitty.conf")
	kittyContent, err := os.ReadFile(kittyConf)
	require.NoError(t, err)

	kittyStr := string(kittyContent)
	assert.Contains(t, kittyStr, "alice", "Should contain username")
	assert.Contains(t, kittyStr, "Catppuccin Mocha", "Should contain theme name")
	assert.NotContains(t, kittyStr, "{{theme_", "Should not contain theme placeholders")

	// And configurations should be ready to use without manual editing
	assert.NotContains(t, configStr, "{{", "No template placeholders should remain")
}

// TestConfigurationDeployment_SelectiveDeployment corresponds to:
// Scenario: Selective configuration deployment
func TestConfigurationDeployment_SelectiveDeployment(t *testing.T) {
	setupTemplateFiles(t)
	tmpDir := t.TempDir()
	ctx := context.Background()

	// Setup services
	templateEngine := templates.NewTemplateEngine()
	backupRoot := filepath.Join(tmpDir, "backups")
	backupService := backup.NewBackupService(backupRoot)
	deployer := configservice.NewConfigDeployer(templateEngine, backupService)
	useCase := configuration.NewConfigDeployUseCase(deployer, templateEngine)

	// Given I only want to deploy Hyprland configurations
	request := configuration.DeployConfigRequest{
		Components: []string{"hyprland"}, // Only hyprland
		CustomVars: map[string]string{
			"username": "testuser",
			"home":     tmpDir,
		},
		DryRun: false,
	}

	// When I specify selective deployment
	resp, err := useCase.Execute(ctx, request)

	// Then only Hyprland configs should be deployed
	require.NoError(t, err)
	assert.Equal(t, 1, resp.TotalFiles)
	assert.Equal(t, 1, resp.SuccessfulFiles)

	hyprlandConf := filepath.Join(tmpDir, ".config", "hypr", "hyprland.conf")
	assert.FileExists(t, hyprlandConf)

	// And Waybar configs should not be deployed
	waybarConfig := filepath.Join(tmpDir, ".config", "waybar", "config.jsonc")
	assert.NoFileExists(t, waybarConfig)

	// And Kitty configs should not be deployed
	kittyConf := filepath.Join(tmpDir, ".config", "kitty", "kitty.conf")
	assert.NoFileExists(t, kittyConf)
}

// TestConfigurationDeployment_TemplateSubstitution validates that template
// variables are properly replaced in configuration files
func TestConfigurationDeployment_TemplateSubstitution(t *testing.T) {
	setupTemplateFiles(t)
	tmpDir := t.TempDir()
	ctx := context.Background()

	// Setup services
	templateEngine := templates.NewTemplateEngine()
	backupRoot := filepath.Join(tmpDir, "backups")
	backupService := backup.NewBackupService(backupRoot)
	deployer := configservice.NewConfigDeployer(templateEngine, backupService)
	useCase := configuration.NewConfigDeployUseCase(deployer, templateEngine)

	// Deploy with custom template variables
	request := configuration.DeployConfigRequest{
		Components: []string{"hyprland", "waybar", "kitty"},
		CustomVars: map[string]string{
			"username": "bob",
			"home":     tmpDir,
			"home_dir": tmpDir,
		},
		DryRun: false,
	}

	resp, err := useCase.Execute(ctx, request)
	require.NoError(t, err)
	assert.Equal(t, 4, resp.SuccessfulFiles) // hyprland(1) + waybar(2) + kitty(1)

	// Verify hyprland template variables were replaced
	hyprlandConf := filepath.Join(tmpDir, ".config", "hypr", "hyprland.conf")
	content, err := os.ReadFile(hyprlandConf)
	require.NoError(t, err)

	configStr := string(content)

	// User variables should be substituted
	assert.Contains(t, configStr, "User: bob")
	assert.Contains(t, configStr, tmpDir, "Should contain tmpDir as home path")

	// Theme variables should be substituted (Catppuccin Mocha by default)
	assert.Contains(t, configStr, "Catppuccin Mocha")
	assert.Contains(t, configStr, "mocha")

	// No template placeholders should remain
	assert.NotContains(t, configStr, "{{username}}")
	assert.NotContains(t, configStr, "{{home}}")
	assert.NotContains(t, configStr, "{{theme_")

	// Verify waybar style has theme colors substituted
	waybarStyle := filepath.Join(tmpDir, ".config", "waybar", "style.css")
	styleContent, err := os.ReadFile(waybarStyle)
	require.NoError(t, err)

	styleStr := string(styleContent)

	// Should have hex colors (without # prefix)
	assert.Contains(t, styleStr, "1e1e2e") // theme_base
	assert.Contains(t, styleStr, "cdd6f4") // theme_text
	assert.Contains(t, styleStr, "cba6f7") // theme_mauve

	// No template placeholders for colors
	assert.NotContains(t, styleStr, "{{theme_base}}")
	assert.NotContains(t, styleStr, "{{theme_text}}")

	// Verify kitty config has colors
	kittyConf := filepath.Join(tmpDir, ".config", "kitty", "kitty.conf")
	kittyContent, err := os.ReadFile(kittyConf)
	require.NoError(t, err)

	kittyStr := string(kittyContent)

	// Should have theme colors for cursor, etc.
	lines := strings.Split(kittyStr, "\n")
	foundCursor := false
	for _, line := range lines {
		if strings.HasPrefix(line, "cursor ") {
			// Should be hex color, not template
			assert.NotContains(t, line, "{{")
			assert.NotContains(t, line, "}}")
			foundCursor = true
		}
	}
	assert.True(t, foundCursor, "Should have cursor color line")
}

// TestConfigurationDeployment_FilePermissions validates that deployed
// configuration files have appropriate permissions
func TestConfigurationDeployment_FilePermissions(t *testing.T) {
	setupTemplateFiles(t)
	tmpDir := t.TempDir()
	ctx := context.Background()

	// Setup services
	templateEngine := templates.NewTemplateEngine()
	backupRoot := filepath.Join(tmpDir, "backups")
	backupService := backup.NewBackupService(backupRoot)
	deployer := configservice.NewConfigDeployer(templateEngine, backupService)
	useCase := configuration.NewConfigDeployUseCase(deployer, templateEngine)

	// Deploy configurations
	request := configuration.DeployConfigRequest{
		Components: []string{"hyprland", "waybar", "kitty"},
		CustomVars: map[string]string{
			"username": "testuser",
			"home":     tmpDir,
		},
		DryRun: false,
	}

	_, err := useCase.Execute(ctx, request)
	require.NoError(t, err)

	// Verify file permissions
	configFiles := []struct {
		path         string
		expectedPerm os.FileMode
	}{
		{filepath.Join(tmpDir, ".config", "hypr", "hyprland.conf"), 0644},
		{filepath.Join(tmpDir, ".config", "waybar", "config.jsonc"), 0644},
		{filepath.Join(tmpDir, ".config", "waybar", "style.css"), 0644},
		{filepath.Join(tmpDir, ".config", "kitty", "kitty.conf"), 0644},
	}

	for _, file := range configFiles {
		info, err := os.Stat(file.path)
		require.NoError(t, err, "Config file %s should exist", file.path)

		actualPerm := info.Mode().Perm()
		assert.Equal(t, file.expectedPerm, actualPerm,
			"File %s should have correct permissions", file.path)
	}
}

// TestConfigurationDeployment_SetupDirectoryStructure corresponds to:
// Scenario: Set up required directory structure
func TestConfigurationDeployment_SetupDirectoryStructure(t *testing.T) {
	setupTemplateFiles(t)
	tmpDir := t.TempDir()
	ctx := context.Background()

	// Setup services
	templateEngine := templates.NewTemplateEngine()
	backupRoot := filepath.Join(tmpDir, "backups")
	backupService := backup.NewBackupService(backupRoot)
	deployer := configservice.NewConfigDeployer(templateEngine, backupService)
	useCase := configuration.NewConfigDeployUseCase(deployer, templateEngine)

	// Given I have no Hyprland configuration directories
	configDir := filepath.Join(tmpDir, ".config")
	_, err := os.Stat(configDir)
	require.True(t, os.IsNotExist(err), "Config directory should not exist yet")

	// When configuration deployment starts
	request := configuration.DeployConfigRequest{
		Components: []string{"hyprland"},
		CustomVars: map[string]string{
			"username": "testuser",
			"home":     tmpDir,
		},
		DryRun: false,
	}

	resp, err := useCase.Execute(ctx, request)

	// Then necessary directories should be created
	require.NoError(t, err)
	assert.DirExists(t, filepath.Join(tmpDir, ".config"))
	assert.DirExists(t, filepath.Join(tmpDir, ".config", "hypr"))

	// And directories should be accessible by me
	info, err := os.Stat(filepath.Join(tmpDir, ".config", "hypr"))
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// And configuration files should be deployed successfully
	assert.Equal(t, 1, resp.SuccessfulFiles)
	assert.FileExists(t, filepath.Join(tmpDir, ".config", "hypr", "hyprland.conf"))
}

// Remaining skipped tests for future implementation

func TestConfigurationDeployment_ReportProgress(t *testing.T) {
	t.Skip("Covered by unit tests - see deploy_config_use_case_test.go ExecuteWithProgress")
}

func TestConfigurationDeployment_CannotWrite(t *testing.T) {
	t.Skip("TODO: Implement permission error handling test")
}

func TestConfigurationDeployment_VerifyDeployment(t *testing.T) {
	t.Skip("Covered by other tests - file existence, permissions, and template substitution")
}

func TestConfigurationDeployment_PreserveCustomizations(t *testing.T) {
	t.Skip("TODO: Implement custom file preservation test when feature is added")
}

func TestConfigurationDeployment_RollbackOnFailure(t *testing.T) {
	t.Skip("TODO: Implement rollback mechanism test")
}
