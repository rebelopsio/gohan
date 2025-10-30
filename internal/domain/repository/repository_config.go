package repository

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrEmptyEntries is returned when repository config has no entries
	ErrEmptyEntries = errors.New("repository config must have at least one entry")
	// ErrInvalidEntryType is returned when entry type is not "deb" or "deb-src"
	ErrInvalidEntryType = errors.New("entry type must be 'deb' or 'deb-src'")
	// ErrMissingURI is returned when entry URI is empty
	ErrMissingURI = errors.New("entry URI cannot be empty")
	// ErrMissingSuite is returned when entry suite is empty
	ErrMissingSuite = errors.New("entry suite cannot be empty")
	// ErrMissingComponents is returned when entry has no components
	ErrMissingComponents = errors.New("entry must have at least one component")
)

// SourceEntry represents a single line in sources.list
type SourceEntry struct {
	Type       string   // "deb" or "deb-src"
	URI        string   // Repository URL
	Suite      string   // Distribution suite (e.g., "sid", "trixie")
	Components []string // Components (e.g., "main", "contrib", "non-free")
}

// String formats the source entry as a sources.list line
func (se SourceEntry) String() string {
	return fmt.Sprintf("%s %s %s %s", se.Type, se.URI, se.Suite, strings.Join(se.Components, " "))
}

// Validate checks if the source entry is valid
func (se SourceEntry) Validate() error {
	if se.Type != "deb" && se.Type != "deb-src" {
		return ErrInvalidEntryType
	}

	if se.URI == "" {
		return ErrMissingURI
	}

	if se.Suite == "" {
		return ErrMissingSuite
	}

	if len(se.Components) == 0 {
		return ErrMissingComponents
	}

	return nil
}

// HasComponent checks if this entry includes a specific component
func (se SourceEntry) HasComponent(component string) bool {
	for _, c := range se.Components {
		if c == component {
			return true
		}
	}
	return false
}

// AddComponent adds a component to this entry if not already present
func (se *SourceEntry) AddComponent(component string) {
	if !se.HasComponent(component) {
		se.Components = append(se.Components, component)
	}
}

// RepositoryConfig represents the complete apt sources configuration
type RepositoryConfig struct {
	entries []SourceEntry
}

// NewRepositoryConfig creates a new repository configuration
func NewRepositoryConfig(entries []SourceEntry) (*RepositoryConfig, error) {
	if len(entries) == 0 {
		return nil, ErrEmptyEntries
	}

	// Validate all entries
	for _, entry := range entries {
		if err := entry.Validate(); err != nil {
			return nil, fmt.Errorf("invalid entry: %w", err)
		}
	}

	return &RepositoryConfig{
		entries: entries,
	}, nil
}

// Entries returns a copy of the entries
func (rc *RepositoryConfig) Entries() []SourceEntry {
	entriesCopy := make([]SourceEntry, len(rc.entries))
	copy(entriesCopy, rc.entries)
	return entriesCopy
}

// HasComponent checks if any entry includes a specific component
func (rc *RepositoryConfig) HasComponent(component string) bool {
	for _, entry := range rc.entries {
		if entry.HasComponent(component) {
			return true
		}
	}
	return false
}

// HasDebSrc checks if there are any deb-src entries
func (rc *RepositoryConfig) HasDebSrc() bool {
	for _, entry := range rc.entries {
		if entry.Type == "deb-src" {
			return true
		}
	}
	return false
}

// AddComponent adds a component to all deb entries
func (rc *RepositoryConfig) AddComponent(component string) error {
	for i := range rc.entries {
		if rc.entries[i].Type == "deb" {
			rc.entries[i].AddComponent(component)
		}
	}
	return nil
}

// EnableDebSrc adds deb-src entries for all deb entries that don't have them
func (rc *RepositoryConfig) EnableDebSrc() {
	// Collect deb entries that need deb-src counterparts
	var newEntries []SourceEntry

	for _, entry := range rc.entries {
		if entry.Type == "deb" {
			// Check if corresponding deb-src already exists
			hasDebSrc := false
			for _, e := range rc.entries {
				if e.Type == "deb-src" && e.URI == entry.URI && e.Suite == entry.Suite {
					hasDebSrc = true
					break
				}
			}

			// Add deb-src if it doesn't exist
			if !hasDebSrc {
				debSrcEntry := SourceEntry{
					Type:       "deb-src",
					URI:        entry.URI,
					Suite:      entry.Suite,
					Components: make([]string, len(entry.Components)),
				}
				copy(debSrcEntry.Components, entry.Components)
				newEntries = append(newEntries, debSrcEntry)
			}
		}
	}

	// Add new deb-src entries
	rc.entries = append(rc.entries, newEntries...)
}

// String formats the entire configuration as sources.list content
func (rc *RepositoryConfig) String() string {
	var lines []string
	for _, entry := range rc.entries {
		lines = append(lines, entry.String())
	}
	return strings.Join(lines, "\n")
}
