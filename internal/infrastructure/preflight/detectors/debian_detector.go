package detectors

import (
	"context"
	"os"
	"strings"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
)

// DebianVersionDetector implements preflight.DebianDetector
type DebianVersionDetector struct {
	osReleasePath string
}

// NewDebianVersionDetector creates a new Debian detector
func NewDebianVersionDetector() *DebianVersionDetector {
	return &DebianVersionDetector{
		osReleasePath: "/etc/os-release",
	}
}

// DetectVersion identifies the Debian version
func (d *DebianVersionDetector) DetectVersion(ctx context.Context) (preflight.DebianVersion, error) {
	content, err := os.ReadFile(d.osReleasePath)
	if err != nil {
		return preflight.DebianVersion{}, err
	}

	codename := ""
	versionID := ""

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "VERSION_CODENAME=") {
			codename = strings.Trim(strings.TrimPrefix(line, "VERSION_CODENAME="), "\"")
		}

		if strings.HasPrefix(line, "VERSION_ID=") {
			versionID = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
		}
	}

	if codename == "" {
		return preflight.DebianVersion{}, preflight.ErrInvalidDebianVersion
	}

	return preflight.NewDebianVersion(codename, versionID)
}

// IsDebianBased checks if system is Debian-based
func (d *DebianVersionDetector) IsDebianBased(ctx context.Context) bool {
	content, err := os.ReadFile(d.osReleasePath)
	if err != nil {
		return false
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "ID=") {
			id := strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
			return id == "debian"
		}

		if strings.HasPrefix(line, "ID_LIKE=") {
			idLike := strings.Trim(strings.TrimPrefix(line, "ID_LIKE="), "\"")
			return strings.Contains(idLike, "debian")
		}
	}

	return false
}
