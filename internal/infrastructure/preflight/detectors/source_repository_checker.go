package detectors

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
)

// SystemSourceRepositoryChecker implements preflight.SourceRepositoryChecker
type SystemSourceRepositoryChecker struct {
	sourcesListPath string
	sourcesListDir  string
}

// NewSystemSourceRepositoryChecker creates a new source repository checker
func NewSystemSourceRepositoryChecker() *SystemSourceRepositoryChecker {
	return &SystemSourceRepositoryChecker{
		sourcesListPath: "/etc/apt/sources.list",
		sourcesListDir:  "/etc/apt/sources.list.d",
	}
}

// CheckSourceRepositories verifies deb-src configuration
func (c *SystemSourceRepositoryChecker) CheckSourceRepositories(ctx context.Context) (preflight.SourceRepositoryStatus, error) {
	sources := make([]string, 0)

	// Read main sources.list
	mainSources, err := c.readSourcesFile(c.sourcesListPath)
	if err == nil {
		sources = append(sources, mainSources...)
	}

	// Read sources.list.d/*.list files
	dirSources, err := c.readSourcesDir(c.sourcesListDir)
	if err == nil {
		sources = append(sources, dirSources...)
	}

	// Check if any deb-src lines exist
	hasDebSrc := false
	for _, source := range sources {
		trimmed := strings.TrimSpace(source)
		if strings.HasPrefix(trimmed, "deb-src") {
			hasDebSrc = true
			break
		}
	}

	return preflight.NewSourceRepositoryStatus(hasDebSrc, sources), nil
}

func (c *SystemSourceRepositoryChecker) readSourcesFile(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	sources := make([]string, 0)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Only include deb and deb-src lines
		if strings.HasPrefix(trimmed, "deb") {
			sources = append(sources, trimmed)
		}
	}

	return sources, nil
}

func (c *SystemSourceRepositoryChecker) readSourcesDir(dirPath string) ([]string, error) {
	sources := make([]string, 0)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only process .list files
		if !strings.HasSuffix(entry.Name(), ".list") {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		fileSources, err := c.readSourcesFile(filePath)
		if err == nil {
			sources = append(sources, fileSources...)
		}
	}

	return sources, nil
}
