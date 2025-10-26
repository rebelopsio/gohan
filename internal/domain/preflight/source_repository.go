package preflight

import "strings"

// SourceRepositoryStatus represents the status of deb-src repositories
type SourceRepositoryStatus struct {
	isEnabled         bool
	configuredSources []string
}

// NewSourceRepositoryStatus creates a new source repository status
func NewSourceRepositoryStatus(
	isEnabled bool,
	configuredSources []string,
) SourceRepositoryStatus {
	// Ensure non-nil slice
	if configuredSources == nil {
		configuredSources = []string{}
	}

	return SourceRepositoryStatus{
		isEnabled:         isEnabled,
		configuredSources: configuredSources,
	}
}

// IsEnabled returns true if deb-src is enabled
func (s SourceRepositoryStatus) IsEnabled() bool {
	return s.isEnabled
}

// ConfiguredSources returns all configured source repos
func (s SourceRepositoryStatus) ConfiguredSources() []string {
	if s.configuredSources == nil {
		return nil
	}
	result := make([]string, len(s.configuredSources))
	copy(result, s.configuredSources)
	return result
}

// HasDebSrc checks if deb-src lines exist
func (s SourceRepositoryStatus) HasDebSrc() bool {
	for _, source := range s.configuredSources {
		if strings.HasPrefix(strings.TrimSpace(source), "deb-src") {
			return true
		}
	}
	return false
}

// String returns human-readable representation
func (s SourceRepositoryStatus) String() string {
	if s.isEnabled {
		return "Source repositories enabled"
	}
	return "Source repositories not enabled"
}
