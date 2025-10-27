package history

import (
	"strings"
)

// SystemContext is a value object capturing system information at installation time
type SystemContext struct {
	osVersion     string
	kernelVersion string
	gohanVersion  string
	hostname      string
}

// NewSystemContext creates system context value object
func NewSystemContext(
	osVersion string,
	kernelVersion string,
	gohanVersion string,
	hostname string,
) (SystemContext, error) {
	// Trim all fields
	osVersion = strings.TrimSpace(osVersion)
	kernelVersion = strings.TrimSpace(kernelVersion)
	gohanVersion = strings.TrimSpace(gohanVersion)
	hostname = strings.TrimSpace(hostname)

	// OS version is required
	if osVersion == "" {
		return SystemContext{}, ErrInvalidSystemContext
	}

	return SystemContext{
		osVersion:     osVersion,
		kernelVersion: kernelVersion,
		gohanVersion:  gohanVersion,
		hostname:      hostname,
	}, nil
}

// OSVersion returns the OS version
func (s SystemContext) OSVersion() string {
	return s.osVersion
}

// KernelVersion returns the kernel version
func (s SystemContext) KernelVersion() string {
	return s.kernelVersion
}

// GohanVersion returns the Gohan version
func (s SystemContext) GohanVersion() string {
	return s.gohanVersion
}

// Hostname returns the hostname
func (s SystemContext) Hostname() string {
	return s.hostname
}
