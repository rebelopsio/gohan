package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	// Database paths
	Database DatabaseConfig `yaml:"database"`

	// API settings
	API APIConfig `yaml:"api"`

	// Installation settings
	Installation InstallationConfig `yaml:"installation"`

	// Logging settings
	Logging LoggingConfig `yaml:"logging"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	// History database path
	HistoryDB string `yaml:"history_db"`

	// Installation session database path
	InstallationDB string `yaml:"installation_db"`
}

// APIConfig holds API server configuration
type APIConfig struct {
	// Server host
	Host string `yaml:"host"`

	// Server port
	Port int `yaml:"port"`

	// Enable CORS
	EnableCORS bool `yaml:"enable_cors"`
}

// InstallationConfig holds installation-specific settings
type InstallationConfig struct {
	// Default snapshot directory
	SnapshotDir string `yaml:"snapshot_dir"`

	// Enable automatic backups
	AutoBackup bool `yaml:"auto_backup"`

	// History retention days (0 = unlimited)
	HistoryRetentionDays int `yaml:"history_retention_days"`

	// Dry-run mode (no actual package installation)
	DryRun bool `yaml:"dry_run"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	// Log level (debug, info, warn, error)
	Level string `yaml:"level"`

	// Log file path (empty = stdout)
	File string `yaml:"file"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	gohanDir := filepath.Join(homeDir, ".gohan")

	return &Config{
		Database: DatabaseConfig{
			HistoryDB:      filepath.Join(gohanDir, "history.db"),
			InstallationDB: filepath.Join(gohanDir, "installations.db"),
		},
		API: APIConfig{
			Host:       "localhost",
			Port:       8080,
			EnableCORS: true,
		},
		Installation: InstallationConfig{
			SnapshotDir:          filepath.Join(gohanDir, "snapshots"),
			AutoBackup:           true,
			HistoryRetentionDays: 90,
		},
		Logging: LoggingConfig{
			Level: "info",
			File:  "",
		},
	}
}

// Load loads configuration from file, falling back to defaults
func Load() (*Config, error) {
	// Start with defaults
	cfg := DefaultConfig()

	// Try to load from config file
	configPath := GetConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		// Config file exists, load it
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Ensure directories exist
	if err := cfg.EnsureDirectories(); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	return cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath := GetConfigPath()

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// EnsureDirectories ensures all required directories exist
func (c *Config) EnsureDirectories() error {
	dirs := []string{
		filepath.Dir(c.Database.HistoryDB),
		filepath.Dir(c.Database.InstallationDB),
		c.Installation.SnapshotDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".gohan", "config.yaml")
}

// GetDataDir returns the gohan data directory
func GetDataDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".gohan")
}
