package postinstall

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// SystemdServiceManager implements ServiceManager using systemd
type SystemdServiceManager struct{}

// NewSystemdServiceManager creates a new systemd service manager
func NewSystemdServiceManager() *SystemdServiceManager {
	return &SystemdServiceManager{}
}

// Enable enables a systemd service
func (s *SystemdServiceManager) Enable(ctx context.Context, service string) error {
	cmd := exec.CommandContext(ctx, "sudo", "systemctl", "enable", service)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable service %s: %w, output: %s", service, err, string(output))
	}
	return nil
}

// Disable disables a systemd service
func (s *SystemdServiceManager) Disable(ctx context.Context, service string) error {
	cmd := exec.CommandContext(ctx, "sudo", "systemctl", "disable", service)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable service %s: %w, output: %s", service, err, string(output))
	}
	return nil
}

// Start starts a systemd service
func (s *SystemdServiceManager) Start(ctx context.Context, service string) error {
	cmd := exec.CommandContext(ctx, "sudo", "systemctl", "start", service)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start service %s: %w, output: %s", service, err, string(output))
	}
	return nil
}

// Stop stops a systemd service
func (s *SystemdServiceManager) Stop(ctx context.Context, service string) error {
	cmd := exec.CommandContext(ctx, "sudo", "systemctl", "stop", service)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop service %s: %w, output: %s", service, err, string(output))
	}
	return nil
}

// IsEnabled checks if a service is enabled
func (s *SystemdServiceManager) IsEnabled(ctx context.Context, service string) (bool, error) {
	cmd := exec.CommandContext(ctx, "systemctl", "is-enabled", service)
	output, err := cmd.Output()
	if err != nil {
		// If exit code is not 0, service is not enabled
		return false, nil
	}
	return strings.TrimSpace(string(output)) == "enabled", nil
}

// IsActive checks if a service is active (running)
func (s *SystemdServiceManager) IsActive(ctx context.Context, service string) (bool, error) {
	cmd := exec.CommandContext(ctx, "systemctl", "is-active", service)
	output, err := cmd.Output()
	if err != nil {
		// If exit code is not 0, service is not active
		return false, nil
	}
	return strings.TrimSpace(string(output)) == "active", nil
}
