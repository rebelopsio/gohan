package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rebelopsio/gohan/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	require.NotNil(t, cfg)

	// API defaults
	assert.Equal(t, "localhost", cfg.API.Host)
	assert.Equal(t, 8080, cfg.API.Port)
	assert.True(t, cfg.API.EnableCORS)

	// Database defaults
	homeDir, _ := os.UserHomeDir()
	gohanDir := filepath.Join(homeDir, ".gohan")
	assert.Equal(t, filepath.Join(gohanDir, "history.db"), cfg.Database.HistoryDB)
	assert.Equal(t, filepath.Join(gohanDir, "installations.db"), cfg.Database.InstallationDB)

	// Installation defaults
	assert.Equal(t, filepath.Join(gohanDir, "snapshots"), cfg.Installation.SnapshotDir)
	assert.True(t, cfg.Installation.AutoBackup)
	assert.Equal(t, 90, cfg.Installation.HistoryRetentionDays)

	// Logging defaults
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "", cfg.Logging.File)
}

func TestLoad(t *testing.T) {
	t.Run("loads default configuration when no file exists", func(t *testing.T) {
		cfg, err := config.Load()

		require.NoError(t, err)
		require.NotNil(t, cfg)

		// Should return defaults
		assert.Equal(t, "localhost", cfg.API.Host)
		assert.Equal(t, 8080, cfg.API.Port)
		assert.Equal(t, "info", cfg.Logging.Level)
	})
}

func TestConfig_EnsureDirectories(t *testing.T) {
	// Create a temporary config with custom paths
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			HistoryDB:      filepath.Join(tmpDir, "db", "history.db"),
			InstallationDB: filepath.Join(tmpDir, "db", "installations.db"),
		},
		Installation: config.InstallationConfig{
			SnapshotDir: filepath.Join(tmpDir, "snapshots"),
		},
	}

	err := cfg.EnsureDirectories()
	require.NoError(t, err)

	// Check that directories were created
	assert.DirExists(t, filepath.Join(tmpDir, "db"))
	assert.DirExists(t, filepath.Join(tmpDir, "snapshots"))
}

func TestGetConfigPath(t *testing.T) {
	path := config.GetConfigPath()

	homeDir, _ := os.UserHomeDir()
	expected := filepath.Join(homeDir, ".gohan", "config.yaml")

	assert.Equal(t, expected, path)
}

func TestGetDataDir(t *testing.T) {
	dir := config.GetDataDir()

	homeDir, _ := os.UserHomeDir()
	expected := filepath.Join(homeDir, ".gohan")

	assert.Equal(t, expected, dir)
}
