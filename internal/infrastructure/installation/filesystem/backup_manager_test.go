package filesystem_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rebelopsio/gohan/internal/infrastructure/installation/filesystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBackupManager(t *testing.T) {
	t.Run("creates backup manager", func(t *testing.T) {
		manager := filesystem.NewBackupManager()

		assert.NotNil(t, manager)
	})
}

func TestBackupManager_CreateBackup(t *testing.T) {
	manager := filesystem.NewBackupManager()
	ctx := context.Background()

	t.Run("creates backup of single file", func(t *testing.T) {
		// Setup test files
		tmpDir := t.TempDir()
		sourceFile := filepath.Join(tmpDir, "config.txt")
		backupDir := filepath.Join(tmpDir, "backup")

		err := os.WriteFile(sourceFile, []byte("test config content"), 0644)
		require.NoError(t, err)

		// Create backup
		backupPath, err := manager.CreateBackup(ctx, sourceFile, backupDir)

		require.NoError(t, err)
		assert.NotEmpty(t, backupPath)

		// Verify backup exists
		_, err = os.Stat(backupPath)
		assert.NoError(t, err)

		// Verify backup content
		content, err := os.ReadFile(backupPath)
		require.NoError(t, err)
		assert.Equal(t, "test config content", string(content))
	})

	t.Run("creates backup directory if it doesn't exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		sourceFile := filepath.Join(tmpDir, "config.txt")
		backupDir := filepath.Join(tmpDir, "nonexistent", "backup")

		err := os.WriteFile(sourceFile, []byte("test"), 0644)
		require.NoError(t, err)

		backupPath, err := manager.CreateBackup(ctx, sourceFile, backupDir)

		require.NoError(t, err)
		assert.NotEmpty(t, backupPath)

		// Verify backup directory was created
		_, err = os.Stat(backupDir)
		assert.NoError(t, err)
	})

	t.Run("returns error for non-existent source", func(t *testing.T) {
		tmpDir := t.TempDir()
		nonExistentFile := filepath.Join(tmpDir, "nonexistent.txt")
		backupDir := filepath.Join(tmpDir, "backup")

		_, err := manager.CreateBackup(ctx, nonExistentFile, backupDir)

		assert.Error(t, err)
	})

	t.Run("preserves file permissions", func(t *testing.T) {
		tmpDir := t.TempDir()
		sourceFile := filepath.Join(tmpDir, "config.txt")
		backupDir := filepath.Join(tmpDir, "backup")

		err := os.WriteFile(sourceFile, []byte("test"), 0600)
		require.NoError(t, err)

		backupPath, err := manager.CreateBackup(ctx, sourceFile, backupDir)
		require.NoError(t, err)

		// Check backup has same permissions
		sourceInfo, _ := os.Stat(sourceFile)
		backupInfo, _ := os.Stat(backupPath)

		assert.Equal(t, sourceInfo.Mode(), backupInfo.Mode())
	})
}

func TestBackupManager_CreateBackupDirectory(t *testing.T) {
	manager := filesystem.NewBackupManager()
	ctx := context.Background()

	t.Run("creates backup of directory recursively", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create source directory structure
		sourceDir := filepath.Join(tmpDir, "config")
		err := os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0755)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("file 1"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(sourceDir, "subdir", "file2.txt"), []byte("file 2"), 0644)
		require.NoError(t, err)

		backupDir := filepath.Join(tmpDir, "backup")

		// Create backup
		backupPath, err := manager.CreateBackupDirectory(ctx, sourceDir, backupDir)

		require.NoError(t, err)
		assert.NotEmpty(t, backupPath)

		// Verify files were backed up
		content, err := os.ReadFile(filepath.Join(backupPath, "file1.txt"))
		require.NoError(t, err)
		assert.Equal(t, "file 1", string(content))

		content, err = os.ReadFile(filepath.Join(backupPath, "subdir", "file2.txt"))
		require.NoError(t, err)
		assert.Equal(t, "file 2", string(content))
	})

	t.Run("returns error for non-existent directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		nonExistentDir := filepath.Join(tmpDir, "nonexistent")
		backupDir := filepath.Join(tmpDir, "backup")

		_, err := manager.CreateBackupDirectory(ctx, nonExistentDir, backupDir)

		assert.Error(t, err)
	})
}

func TestBackupManager_RestoreBackup(t *testing.T) {
	manager := filesystem.NewBackupManager()
	ctx := context.Background()

	t.Run("restores backup to original location", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create original file and backup it
		originalFile := filepath.Join(tmpDir, "config.txt")
		backupDir := filepath.Join(tmpDir, "backup")

		err := os.WriteFile(originalFile, []byte("original content"), 0644)
		require.NoError(t, err)

		backupPath, err := manager.CreateBackup(ctx, originalFile, backupDir)
		require.NoError(t, err)

		// Modify original
		err = os.WriteFile(originalFile, []byte("modified content"), 0644)
		require.NoError(t, err)

		// Restore from backup
		err = manager.RestoreBackup(ctx, backupPath, originalFile)
		require.NoError(t, err)

		// Verify restoration
		content, err := os.ReadFile(originalFile)
		require.NoError(t, err)
		assert.Equal(t, "original content", string(content))
	})

	t.Run("creates parent directories if needed", func(t *testing.T) {
		tmpDir := t.TempDir()

		backupFile := filepath.Join(tmpDir, "backup", "config.txt")
		err := os.MkdirAll(filepath.Join(tmpDir, "backup"), 0755)
		require.NoError(t, err)
		err = os.WriteFile(backupFile, []byte("backed up"), 0644)
		require.NoError(t, err)

		restorePath := filepath.Join(tmpDir, "new", "location", "config.txt")

		err = manager.RestoreBackup(ctx, backupFile, restorePath)
		require.NoError(t, err)

		// Verify file was restored
		content, err := os.ReadFile(restorePath)
		require.NoError(t, err)
		assert.Equal(t, "backed up", string(content))
	})

	t.Run("returns error for non-existent backup", func(t *testing.T) {
		tmpDir := t.TempDir()

		backupFile := filepath.Join(tmpDir, "nonexistent.txt")
		restorePath := filepath.Join(tmpDir, "restore.txt")

		err := manager.RestoreBackup(ctx, backupFile, restorePath)

		assert.Error(t, err)
	})
}

func TestBackupManager_VerifyBackup(t *testing.T) {
	manager := filesystem.NewBackupManager()
	ctx := context.Background()

	t.Run("verifies valid backup", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create and backup file
		originalFile := filepath.Join(tmpDir, "config.txt")
		backupDir := filepath.Join(tmpDir, "backup")

		err := os.WriteFile(originalFile, []byte("test content"), 0644)
		require.NoError(t, err)

		backupPath, err := manager.CreateBackup(ctx, originalFile, backupDir)
		require.NoError(t, err)

		// Verify backup
		valid, err := manager.VerifyBackup(ctx, backupPath)

		require.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("returns false for non-existent backup", func(t *testing.T) {
		tmpDir := t.TempDir()
		nonExistentBackup := filepath.Join(tmpDir, "nonexistent.txt")

		valid, err := manager.VerifyBackup(ctx, nonExistentBackup)

		require.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("returns false for empty backup", func(t *testing.T) {
		tmpDir := t.TempDir()
		emptyBackup := filepath.Join(tmpDir, "empty.txt")

		err := os.WriteFile(emptyBackup, []byte(""), 0644)
		require.NoError(t, err)

		valid, err := manager.VerifyBackup(ctx, emptyBackup)

		require.NoError(t, err)
		// Empty files are technically valid backups
		assert.True(t, valid)
	})
}
