package repository

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	domainRepo "github.com/rebelopsio/gohan/internal/domain/repository"
)

var (
	// ErrMissingCodename is returned when VERSION_CODENAME is not found
	ErrMissingCodename = errors.New("VERSION_CODENAME not found in os-release")
)

// OSReleaseInfo contains information from /etc/os-release
type OSReleaseInfo struct {
	Name       string
	VersionID  string
	Codename   string
	PrettyName string
}

// ParseOSRelease parses the content of /etc/os-release
func ParseOSRelease(content string) (*OSReleaseInfo, error) {
	info := &OSReleaseInfo{}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes
		value = strings.Trim(value, `"`)

		switch key {
		case "NAME":
			info.Name = value
		case "VERSION_ID":
			info.VersionID = value
		case "VERSION_CODENAME":
			info.Codename = value
		case "PRETTY_NAME":
			info.PrettyName = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read os-release content: %w", err)
	}

	// Codename is required
	if info.Codename == "" {
		return nil, ErrMissingCodename
	}

	return info, nil
}

// ReadOSRelease reads and parses /etc/os-release
func ReadOSRelease() (*OSReleaseInfo, error) {
	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return nil, fmt.Errorf("failed to read /etc/os-release: %w", err)
	}

	return ParseOSRelease(string(content))
}

// DetectVersionFromOSRelease converts OSReleaseInfo to DebianVersion
func DetectVersionFromOSRelease(info *OSReleaseInfo) (*domainRepo.DebianVersion, error) {
	codename := info.Codename
	version := info.VersionID

	// Map Debian codenames to version strings
	switch codename {
	case "sid":
		version = "unstable"
	case "trixie":
		if version == "" {
			version = "testing"
		}
	case "bookworm":
		if version == "" {
			version = "12"
		}
	}

	// For Ubuntu, use VERSION_ID as version
	// Ubuntu uses version numbers like "22.04", "20.04", etc.
	if version == "" {
		version = codename
	}

	return domainRepo.NewDebianVersion(codename, version)
}

// DetectDebianVersion detects the current Debian version from the system
func DetectDebianVersion() (*domainRepo.DebianVersion, error) {
	info, err := ReadOSRelease()
	if err != nil {
		return nil, err
	}

	return DetectVersionFromOSRelease(info)
}
