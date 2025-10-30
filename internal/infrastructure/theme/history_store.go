package theme

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rebelopsio/gohan/internal/domain/theme"
)

// Errors
var (
	ErrNoThemeHistory = errors.New("no theme history available")
)

// ThemeHistoryEntry represents a single theme in history
type ThemeHistoryEntry struct {
	ThemeName theme.ThemeName `json:"theme_name"`
}

// ThemeHistoryData represents the persisted history data
type ThemeHistoryData struct {
	History []ThemeHistoryEntry `json:"history"`
}

// ThemeHistoryStore defines the interface for theme history persistence
type ThemeHistoryStore interface {
	Add(ctx context.Context, themeName theme.ThemeName) error
	GetPrevious(ctx context.Context) (theme.ThemeName, error)
	GetHistory(ctx context.Context) ([]theme.ThemeName, error)
	RemoveLast(ctx context.Context) error
	Clear(ctx context.Context) error
}

// FileThemeHistoryStore implements ThemeHistoryStore using file-based persistence
type FileThemeHistoryStore struct {
	filePath string
	maxEntries int
}

// NewFileThemeHistoryStore creates a new file-based theme history store
func NewFileThemeHistoryStore(filePath string) *FileThemeHistoryStore {
	return &FileThemeHistoryStore{
		filePath:   filePath,
		maxEntries: 10, // Limit to 10 entries
	}
}

// Add records a theme change in history
func (s *FileThemeHistoryStore) Add(ctx context.Context, themeName theme.ThemeName) error {
	// Load existing history
	data, err := s.load()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load history: %w", err)
	}
	
	if data == nil {
		data = &ThemeHistoryData{History: []ThemeHistoryEntry{}}
	}
	
	// Add new entry at the beginning (newest first)
	entry := ThemeHistoryEntry{ThemeName: themeName}
	data.History = append([]ThemeHistoryEntry{entry}, data.History...)
	
	// Limit to maxEntries
	if len(data.History) > s.maxEntries {
		data.History = data.History[:s.maxEntries]
	}
	
	// Save
	return s.save(data)
}

// GetPrevious returns the previous theme (second in history)
func (s *FileThemeHistoryStore) GetPrevious(ctx context.Context) (theme.ThemeName, error) {
	data, err := s.load()
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrNoThemeHistory
		}
		return "", err
	}
	
	// Need at least 2 entries (current and previous)
	if len(data.History) < 2 {
		return "", ErrNoThemeHistory
	}
	
	// Return the second entry (index 1)
	return data.History[1].ThemeName, nil
}

// GetHistory returns all themes in history (newest first)
func (s *FileThemeHistoryStore) GetHistory(ctx context.Context) ([]theme.ThemeName, error) {
	data, err := s.load()
	if err != nil {
		if os.IsNotExist(err) {
			return []theme.ThemeName{}, nil
		}
		return nil, err
	}
	
	// Convert to slice of theme names
	result := make([]theme.ThemeName, len(data.History))
	for i, entry := range data.History {
		result[i] = entry.ThemeName
	}
	
	return result, nil
}

// RemoveLast removes the most recent theme from history
func (s *FileThemeHistoryStore) RemoveLast(ctx context.Context) error {
	data, err := s.load()
	if err != nil {
		return err
	}
	
	if len(data.History) == 0 {
		return nil
	}
	
	// Remove first entry
	data.History = data.History[1:]
	
	return s.save(data)
}

// Clear removes all history
func (s *FileThemeHistoryStore) Clear(ctx context.Context) error {
	data := &ThemeHistoryData{History: []ThemeHistoryEntry{}}
	return s.save(data)
}

// load reads history from disk
func (s *FileThemeHistoryStore) load() (*ThemeHistoryData, error) {
	fileData, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}
	
	var data ThemeHistoryData
	if err := json.Unmarshal(fileData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal history: %w", err)
	}
	
	return &data, nil
}

// save writes history to disk
func (s *FileThemeHistoryStore) save(data *ThemeHistoryData) error {
	// Ensure directory exists
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create history directory: %w", err)
	}
	
	// Marshal to JSON
	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}
	
	// Write atomically
	tmpFile := s.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, fileData, 0644); err != nil {
		return fmt.Errorf("failed to write history: %w", err)
	}
	
	if err := os.Rename(tmpFile, s.filePath); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to save history: %w", err)
	}
	
	return nil
}

// GetDefaultHistoryFilePath returns the default location for the theme history file
func GetDefaultHistoryFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}
	
	return filepath.Join(configDir, "gohan", "theme-history.json"), nil
}
