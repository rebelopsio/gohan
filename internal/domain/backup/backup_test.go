package backup_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/backup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBackup(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
	}{
		{
			name:        "creates backup with valid description",
			description: "Pre-installation backup",
			wantErr:     false,
		},
		{
			name:        "creates backup with empty description",
			description: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := backup.NewBackup(tt.description)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, b.ID())
			assert.Equal(t, tt.description, b.Description())
			assert.NotZero(t, b.CreatedAt())
			assert.Empty(t, b.Files())
			assert.Equal(t, backup.StatusPending, b.Status())
		})
	}
}

func TestBackup_AddFile(t *testing.T) {
	t.Run("adds file to backup", func(t *testing.T) {
		b, err := backup.NewBackup("test backup")
		require.NoError(t, err)

		file := backup.BackupFile{
			OriginalPath: "/home/user/.config/hypr/hyprland.conf",
			BackupPath:   "/backups/2025-01-01_120000/hyprland.conf",
			Permissions:  0644,
			SizeBytes:    1024,
		}

		err = b.AddFile(file)
		require.NoError(t, err)

		files := b.Files()
		assert.Len(t, files, 1)
		assert.Equal(t, file.OriginalPath, files[0].OriginalPath)
	})

	t.Run("prevents adding duplicate files", func(t *testing.T) {
		b, err := backup.NewBackup("test backup")
		require.NoError(t, err)

		file := backup.BackupFile{
			OriginalPath: "/home/user/.config/hypr/hyprland.conf",
			BackupPath:   "/backups/2025-01-01_120000/hyprland.conf",
			Permissions:  0644,
			SizeBytes:    1024,
		}

		err = b.AddFile(file)
		require.NoError(t, err)

		// Try to add same file again
		err = b.AddFile(file)
		assert.Error(t, err)
		assert.ErrorIs(t, err, backup.ErrDuplicateFile)

		// Should still only have one file
		assert.Len(t, b.Files(), 1)
	})
}

func TestBackup_MarkComplete(t *testing.T) {
	t.Run("marks backup as complete", func(t *testing.T) {
		b, err := backup.NewBackup("test backup")
		require.NoError(t, err)

		assert.Equal(t, backup.StatusPending, b.Status())

		b.MarkComplete()

		assert.Equal(t, backup.StatusComplete, b.Status())
	})
}

func TestBackup_MarkFailed(t *testing.T) {
	t.Run("marks backup as failed with error", func(t *testing.T) {
		b, err := backup.NewBackup("test backup")
		require.NoError(t, err)

		testErr := "disk full"
		b.MarkFailed(testErr)

		assert.Equal(t, backup.StatusFailed, b.Status())
		assert.Equal(t, testErr, b.Error())
	})
}

func TestBackup_TotalSize(t *testing.T) {
	t.Run("calculates total size of all files", func(t *testing.T) {
		b, err := backup.NewBackup("test backup")
		require.NoError(t, err)

		file1 := backup.BackupFile{
			OriginalPath: "/config/file1",
			BackupPath:   "/backup/file1",
			Permissions:  0644,
			SizeBytes:    1024,
		}

		file2 := backup.BackupFile{
			OriginalPath: "/config/file2",
			BackupPath:   "/backup/file2",
			Permissions:  0644,
			SizeBytes:    2048,
		}

		require.NoError(t, b.AddFile(file1))
		require.NoError(t, b.AddFile(file2))

		assert.Equal(t, int64(3072), b.TotalSize())
	})

	t.Run("returns zero for empty backup", func(t *testing.T) {
		b, err := backup.NewBackup("empty backup")
		require.NoError(t, err)

		assert.Equal(t, int64(0), b.TotalSize())
	})
}

func TestBackup_Contains(t *testing.T) {
	t.Run("checks if backup contains specific file", func(t *testing.T) {
		b, err := backup.NewBackup("test backup")
		require.NoError(t, err)

		file := backup.BackupFile{
			OriginalPath: "/config/hyprland.conf",
			BackupPath:   "/backup/hyprland.conf",
			Permissions:  0644,
			SizeBytes:    1024,
		}

		require.NoError(t, b.AddFile(file))

		assert.True(t, b.Contains("/config/hyprland.conf"))
		assert.False(t, b.Contains("/config/waybar.conf"))
	})
}

func TestBackup_Age(t *testing.T) {
	t.Run("calculates age of backup", func(t *testing.T) {
		b, err := backup.NewBackup("test backup")
		require.NoError(t, err)

		// Sleep briefly to ensure non-zero age
		time.Sleep(10 * time.Millisecond)

		age := b.Age()
		assert.Greater(t, age, time.Duration(0))
	})
}

func TestBackupFile_Validate(t *testing.T) {
	tests := []struct {
		name    string
		file    backup.BackupFile
		wantErr bool
	}{
		{
			name: "valid file",
			file: backup.BackupFile{
				OriginalPath: "/config/file",
				BackupPath:   "/backup/file",
				Permissions:  0644,
				SizeBytes:    100,
			},
			wantErr: false,
		},
		{
			name: "missing original path",
			file: backup.BackupFile{
				OriginalPath: "",
				BackupPath:   "/backup/file",
				Permissions:  0644,
				SizeBytes:    100,
			},
			wantErr: true,
		},
		{
			name: "missing backup path",
			file: backup.BackupFile{
				OriginalPath: "/config/file",
				BackupPath:   "",
				Permissions:  0644,
				SizeBytes:    100,
			},
			wantErr: true,
		},
		{
			name: "negative size",
			file: backup.BackupFile{
				OriginalPath: "/config/file",
				BackupPath:   "/backup/file",
				Permissions:  0644,
				SizeBytes:    -100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.file.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
