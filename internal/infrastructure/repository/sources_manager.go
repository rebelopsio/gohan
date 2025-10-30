package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	domainRepo "github.com/rebelopsio/gohan/internal/domain/repository"
)

// FileSourcesManager manages sources.list files on the filesystem
type FileSourcesManager struct{}

// NewFileSourcesManager creates a new file-based sources.list manager
func NewFileSourcesManager() *FileSourcesManager {
	return &FileSourcesManager{}
}

// ReadConfig reads and parses a sources.list file into a RepositoryConfig
func (m *FileSourcesManager) ReadConfig(path string) (*domainRepo.RepositoryConfig, error) {
	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// Parse entries
	entries, err := ParseSourcesFile(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse sources.list: %w", err)
	}

	// Create config
	config, err := domainRepo.NewRepositoryConfig(entries)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository config: %w", err)
	}

	return config, nil
}

// WriteConfig writes a RepositoryConfig to a sources.list file
func (m *FileSourcesManager) WriteConfig(path string, config *domainRepo.RepositoryConfig) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Convert config to string
	content := config.String() + "\n"

	// Write file with appropriate permissions (root-owned files)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// Backup creates a timestamped backup of a sources.list file
func (m *FileSourcesManager) Backup(path string, backupDir string) (string, error) {
	// Read current content
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	// Ensure backup directory exists
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory %s: %w", backupDir, err)
	}

	// Create timestamped backup filename
	timestamp := time.Now().Format("20060102-150405")
	backupName := fmt.Sprintf("%s.backup-%s", filepath.Base(path), timestamp)
	backupPath := filepath.Join(backupDir, backupName)

	// Write backup
	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup file %s: %w", backupPath, err)
	}

	return backupPath, nil
}

// Exists checks if a file exists
func (m *FileSourcesManager) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// SystemVersionDetector detects the Debian version from the actual system
type SystemVersionDetector struct{}

// NewSystemVersionDetector creates a new system version detector
func NewSystemVersionDetector() *SystemVersionDetector {
	return &SystemVersionDetector{}
}

// DetectVersion detects the current Debian version from /etc/os-release
func (d *SystemVersionDetector) DetectVersion() (*domainRepo.DebianVersion, error) {
	return DetectDebianVersion()
}
