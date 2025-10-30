package theme

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/theme"
)

// ThemeState represents the persisted state of the active theme
type ThemeState struct {
	ThemeName   theme.ThemeName   `json:"theme_name"`
	ThemeVariant theme.ThemeVariant `json:"theme_variant"`
	SetAt       time.Time          `json:"set_at"`
}

// ThemeStateStore defines the interface for persisting theme state
type ThemeStateStore interface {
	Save(ctx context.Context, state *ThemeState) error
	Load(ctx context.Context) (*ThemeState, error)
	Exists(ctx context.Context) (bool, error)
}

// FileThemeStateStore implements ThemeStateStore using file-based persistence
type FileThemeStateStore struct {
	filePath string
}

// NewFileThemeStateStore creates a new file-based theme state store
func NewFileThemeStateStore(filePath string) *FileThemeStateStore {
	return &FileThemeStateStore{
		filePath: filePath,
	}
}

// Save persists the theme state to disk
func (s *FileThemeStateStore) Save(ctx context.Context, state *ThemeState) error {
	// Ensure directory exists
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	// Marshal state to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal theme state: %w", err)
	}

	// Write to file atomically by writing to temp file first
	tmpFile := s.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write theme state: %w", err)
	}

	// Rename temp file to actual file (atomic on most filesystems)
	if err := os.Rename(tmpFile, s.filePath); err != nil {
		os.Remove(tmpFile) // Clean up temp file
		return fmt.Errorf("failed to save theme state: %w", err)
	}

	return nil
}

// Load reads the theme state from disk
func (s *FileThemeStateStore) Load(ctx context.Context) (*ThemeState, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	var state ThemeState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal theme state: %w", err)
	}

	return &state, nil
}

// Exists checks if the theme state file exists
func (s *FileThemeStateStore) Exists(ctx context.Context) (bool, error) {
	_, err := os.Stat(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetDefaultStateFilePath returns the default location for the theme state file
func GetDefaultStateFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}

	return filepath.Join(configDir, "gohan", "theme-state.json"), nil
}
