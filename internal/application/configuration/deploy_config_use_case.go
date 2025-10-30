package configuration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebelopsio/gohan/internal/infrastructure/installation/configservice"
	"github.com/rebelopsio/gohan/internal/infrastructure/installation/templates"
)

// DeployConfigRequest contains parameters for configuration deployment
type DeployConfigRequest struct {
	Components      []string // Which components to deploy (hyprland, waybar, kitty, etc.)
	SkipBackup      bool     // Skip backup of existing configurations
	DryRun          bool     // Preview without actually deploying
	Force           bool     // Overwrite without prompting
	CustomVars      map[string]string // Additional template variables
	ShowProgress    bool     // Show progress during deployment
}

// DeployConfigResponse contains deployment results
type DeployConfigResponse struct {
	DeployedFiles   []DeployedFileInfo
	BackupID        string
	BackupPath      string
	TotalFiles      int
	SuccessfulFiles int
	FailedFiles     int
	SkippedFiles    int
	DurationMs      int64
	DryRun          bool
}

// DeployedFileInfo contains information about a deployed file
type DeployedFileInfo struct {
	Component    string
	TargetPath   string
	Status       string // "deployed", "skipped", "failed", "dry-run"
	BackedUp     bool
	Error        string
}

// ProgressCallback is called for each file deployment
type ProgressCallback func(component string, filePath string, progress float64)

// ConfigDeployUseCase coordinates configuration deployment
type ConfigDeployUseCase struct {
	deployer       *configservice.ConfigDeployer
	templateEngine *templates.TemplateEngine
	homeDir        string
}

// NewConfigDeployUseCase creates a new use case instance
func NewConfigDeployUseCase(
	deployer *configservice.ConfigDeployer,
	templateEngine *templates.TemplateEngine,
) *ConfigDeployUseCase {
	homeDir, _ := os.UserHomeDir()
	return &ConfigDeployUseCase{
		deployer:       deployer,
		templateEngine: templateEngine,
		homeDir:        homeDir,
	}
}

// Execute runs configuration deployment
func (uc *ConfigDeployUseCase) Execute(ctx context.Context, req DeployConfigRequest) (*DeployConfigResponse, error) {
	// Build configuration file list
	configs := uc.buildConfigList(req.Components)

	if len(configs) == 0 {
		return nil, fmt.Errorf("no configurations to deploy for components: %v", req.Components)
	}

	// Prepare template variables
	vars := uc.prepareTemplateVars(req.CustomVars)

	response := &DeployConfigResponse{
		TotalFiles:      len(configs),
		DeployedFiles:   make([]DeployedFileInfo, 0),
		DryRun:          req.DryRun,
	}

	// Dry run - just show what would be deployed
	if req.DryRun {
		for _, config := range configs {
			response.DeployedFiles = append(response.DeployedFiles, DeployedFileInfo{
				Component:  extractComponent(config.TargetPath),
				TargetPath: config.TargetPath,
				Status:     "dry-run",
			})
		}
		return response, nil
	}

	// Deploy configurations
	for _, config := range configs {
		result, err := uc.deployer.DeployWithBackup(ctx, config, vars)

		fileInfo := DeployedFileInfo{
			Component:  extractComponent(config.TargetPath),
			TargetPath: result.FilePath,
		}

		if err != nil {
			fileInfo.Status = "failed"
			fileInfo.Error = err.Error()
			response.FailedFiles++
		} else if result.Success {
			fileInfo.Status = "deployed"
			fileInfo.BackedUp = result.BackupID != ""
			response.SuccessfulFiles++

			// Store backup info (use last one)
			if result.BackupID != "" {
				response.BackupID = result.BackupID
				response.BackupPath = result.BackupPath
			}
		}

		response.DeployedFiles = append(response.DeployedFiles, fileInfo)
	}

	return response, nil
}

// ExecuteWithProgress runs deployment with progress callbacks
func (uc *ConfigDeployUseCase) ExecuteWithProgress(
	ctx context.Context,
	req DeployConfigRequest,
	progressFn ProgressCallback,
) (*DeployConfigResponse, error) {
	// Build configuration file list
	configs := uc.buildConfigList(req.Components)

	if len(configs) == 0 {
		return nil, fmt.Errorf("no configurations to deploy for components: %v", req.Components)
	}

	// Prepare template variables
	vars := uc.prepareTemplateVars(req.CustomVars)

	response := &DeployConfigResponse{
		TotalFiles:    len(configs),
		DeployedFiles: make([]DeployedFileInfo, 0),
		DryRun:        req.DryRun,
	}

	// Dry run - just show what would be deployed
	if req.DryRun {
		for _, config := range configs {
			response.DeployedFiles = append(response.DeployedFiles, DeployedFileInfo{
				Component:  extractComponent(config.TargetPath),
				TargetPath: config.TargetPath,
				Status:     "dry-run",
			})
		}
		return response, nil
	}

	// Deploy with progress
	progressChan := make(chan configservice.DeploymentProgress)
	done := make(chan error)

	go func() {
		done <- uc.deployer.DeployConfigurations(ctx, configs, vars, progressChan)
	}()

	// Process progress updates
	for progress := range progressChan {
		if progressFn != nil {
			component := extractComponent(progress.FilePath)
			progressFn(component, progress.FilePath, progress.PercentComplete)
		}

		// Track results
		if progress.Status == "completed" {
			response.DeployedFiles = append(response.DeployedFiles, DeployedFileInfo{
				Component:  extractComponent(progress.FilePath),
				TargetPath: progress.FilePath,
				Status:     "deployed",
			})
			response.SuccessfulFiles++
		} else if progress.Status == "failed" {
			response.DeployedFiles = append(response.DeployedFiles, DeployedFileInfo{
				Component:  extractComponent(progress.FilePath),
				TargetPath: progress.FilePath,
				Status:     "failed",
				Error:      progress.Error.Error(),
			})
			response.FailedFiles++
		}
	}

	// Wait for completion
	err := <-done
	return response, err
}

func (uc *ConfigDeployUseCase) buildConfigList(components []string) []configservice.ConfigurationFile {
	configs := []configservice.ConfigurationFile{}

	configDir := filepath.Join(uc.homeDir, ".config")

	// If no components specified, deploy all
	if len(components) == 0 {
		components = []string{"hyprland", "waybar", "kitty", "fuzzel"}
	}

	for _, component := range components {
		switch component {
		case "hyprland":
			configs = append(configs, configservice.ConfigurationFile{
				SourceTemplate: "templates/hyprland/hyprland.conf",
				TargetPath:     filepath.Join(configDir, "hypr/hyprland.conf"),
				Permissions:    0644,
				BackupBefore:   true,
			})
		case "waybar":
			configs = append(configs, configservice.ConfigurationFile{
				SourceTemplate: "templates/waybar/config.jsonc",
				TargetPath:     filepath.Join(configDir, "waybar/config.jsonc"),
				Permissions:    0644,
				BackupBefore:   true,
			})
			configs = append(configs, configservice.ConfigurationFile{
				SourceTemplate: "templates/waybar/style.css",
				TargetPath:     filepath.Join(configDir, "waybar/style.css"),
				Permissions:    0644,
				BackupBefore:   true,
			})
		case "kitty":
			configs = append(configs, configservice.ConfigurationFile{
				SourceTemplate: "templates/kitty/kitty.conf",
				TargetPath:     filepath.Join(configDir, "kitty/kitty.conf"),
				Permissions:    0644,
				BackupBefore:   true,
			})
		case "fuzzel":
			configs = append(configs, configservice.ConfigurationFile{
				SourceTemplate: "templates/fuzzel/fuzzel.ini",
				TargetPath:     filepath.Join(configDir, "fuzzel/fuzzel.ini"),
				Permissions:    0644,
				BackupBefore:   true,
			})
		}
	}

	return configs
}

func (uc *ConfigDeployUseCase) prepareTemplateVars(customVars map[string]string) templates.TemplateVars {
	vars := templates.TemplateVars{
		"username":  os.Getenv("USER"),
		"home_dir":  uc.homeDir,
		"config_dir": filepath.Join(uc.homeDir, ".config"),
	}

	// Merge custom variables
	for k, v := range customVars {
		vars[k] = v
	}

	return vars
}

func extractComponent(path string) string {
	// Extract component name from path like ~/.config/hypr/hyprland.conf -> hyprland
	parts := filepath.SplitList(path)
	for i, part := range parts {
		if part == ".config" && i+1 < len(parts) {
			return parts[i+1]
		}
	}

	// Fallback: use parent directory name
	return filepath.Base(filepath.Dir(path))
}
