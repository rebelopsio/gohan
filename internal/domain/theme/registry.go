package theme

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	// ErrThemeNotFound indicates the requested theme does not exist
	ErrThemeNotFound = errors.New("theme not found")
	// ErrThemeAlreadyRegistered indicates a theme with that name already exists
	ErrThemeAlreadyRegistered = errors.New("theme already registered")
	// ErrNoThemesRegistered indicates no themes are available
	ErrNoThemesRegistered = errors.New("no themes registered")
	// ErrNilTheme indicates a nil theme was provided
	ErrNilTheme = errors.New("theme cannot be nil")
)

// ThemeRegistry manages available themes and tracks the active theme
type ThemeRegistry interface {
	// Register adds a new theme to the registry
	Register(theme *Theme) error

	// FindByName retrieves a theme by name
	FindByName(ctx context.Context, name ThemeName) (*Theme, error)

	// ListAll returns all registered themes
	ListAll() []*Theme

	// ListByVariant returns themes filtered by variant (dark/light)
	ListByVariant(variant ThemeVariant) []*Theme

	// GetActive returns the currently active theme
	GetActive(ctx context.Context) (*Theme, error)

	// SetActive sets the active theme
	SetActive(ctx context.Context, name ThemeName) error
}

// InMemoryThemeRegistry implements ThemeRegistry with in-memory storage
type InMemoryThemeRegistry struct {
	mu         sync.RWMutex
	themes     map[ThemeName]*Theme
	activeTheme ThemeName
}

// NewThemeRegistry creates a new in-memory theme registry
func NewThemeRegistry() *InMemoryThemeRegistry {
	return &InMemoryThemeRegistry{
		themes: make(map[ThemeName]*Theme),
	}
}

// Register adds a new theme to the registry
func (r *InMemoryThemeRegistry) Register(theme *Theme) error {
	if theme == nil {
		return ErrNilTheme
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if theme already exists
	if _, exists := r.themes[theme.Name()]; exists {
		return fmt.Errorf("%w: %s", ErrThemeAlreadyRegistered, theme.Name())
	}

	r.themes[theme.Name()] = theme

	// If this is the first theme or it's Mocha, set it as default active
	if len(r.themes) == 1 || theme.Name() == ThemeMocha {
		r.activeTheme = theme.Name()
	}

	return nil
}

// FindByName retrieves a theme by name
func (r *InMemoryThemeRegistry) FindByName(ctx context.Context, name ThemeName) (*Theme, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	theme, exists := r.themes[name]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrThemeNotFound, name)
	}

	return theme, nil
}

// ListAll returns all registered themes
func (r *InMemoryThemeRegistry) ListAll() []*Theme {
	r.mu.RLock()
	defer r.mu.RUnlock()

	themes := make([]*Theme, 0, len(r.themes))
	for _, theme := range r.themes {
		themes = append(themes, theme)
	}

	return themes
}

// ListByVariant returns themes filtered by variant
func (r *InMemoryThemeRegistry) ListByVariant(variant ThemeVariant) []*Theme {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*Theme
	for _, theme := range r.themes {
		if theme.Variant() == variant {
			filtered = append(filtered, theme)
		}
	}

	return filtered
}

// GetActive returns the currently active theme
func (r *InMemoryThemeRegistry) GetActive(ctx context.Context) (*Theme, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.themes) == 0 {
		return nil, ErrNoThemesRegistered
	}

	// If no active theme set, default to Mocha
	if r.activeTheme == "" {
		if mochaTheme, exists := r.themes[ThemeMocha]; exists {
			return mochaTheme, nil
		}
		// If Mocha doesn't exist, return error
		return nil, fmt.Errorf("%w: no active theme set and default (mocha) not available", ErrThemeNotFound)
	}

	theme, exists := r.themes[r.activeTheme]
	if !exists {
		return nil, fmt.Errorf("%w: active theme %s", ErrThemeNotFound, r.activeTheme)
	}

	return theme, nil
}

// SetActive sets the active theme
func (r *InMemoryThemeRegistry) SetActive(ctx context.Context, name ThemeName) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verify theme exists
	if _, exists := r.themes[name]; !exists {
		return fmt.Errorf("%w: %s", ErrThemeNotFound, name)
	}

	r.activeTheme = name
	return nil
}
