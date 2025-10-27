package templates

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// TemplateEngine handles template variable substitution
type TemplateEngine struct {
	// Can add caching or other optimizations later
}

// TemplateVars contains variables for template substitution
type TemplateVars struct {
	Username   string // Current user's username
	Home       string // User's home directory
	ConfigDir  string // User's config directory (~/.config)
	Hostname   string // System hostname
	Display    string // Primary display name (e.g., eDP-1)
	Resolution string // Primary display resolution (e.g., 1920x1080)
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{}
}

// ProcessTemplate processes a template string and substitutes variables
func (e *TemplateEngine) ProcessTemplate(content string, vars TemplateVars) (string, error) {
	result := content

	// Map of variable names to their values
	replacements := map[string]string{
		"{{username}}":   vars.Username,
		"{{home}}":       vars.Home,
		"{{config_dir}}": vars.ConfigDir,
		"{{hostname}}":   vars.Hostname,
		"{{display}}":    vars.Display,
		"{{resolution}}": vars.Resolution,
	}

	// Perform replacements
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result, nil
}

// ProcessFile reads a template file, processes it, and writes the result
func (e *TemplateEngine) ProcessFile(srcPath, dstPath string, vars TemplateVars) error {
	// Read source template
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", srcPath, err)
	}

	// Process template
	processed, err := e.ProcessTemplate(string(content), vars)
	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	// Ensure destination directory exists
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", dstDir, err)
	}

	// Write processed content
	if err := os.WriteFile(dstPath, []byte(processed), 0644); err != nil {
		return fmt.Errorf("failed to write output file %s: %w", dstPath, err)
	}

	return nil
}

// CollectSystemVars collects template variables from the current system
func CollectSystemVars() (TemplateVars, error) {
	vars := TemplateVars{}

	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		return vars, fmt.Errorf("failed to get current user: %w", err)
	}

	vars.Username = currentUser.Username
	vars.Home = currentUser.HomeDir
	vars.ConfigDir = filepath.Join(currentUser.HomeDir, ".config")

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		// Non-fatal - just leave empty
		hostname = ""
	}
	vars.Hostname = hostname

	// Display and Resolution are optional and would typically be detected
	// from the running Wayland/X11 session. For now, leave empty.
	// They can be set manually or detected by a display detection service.
	vars.Display = detectPrimaryDisplay()
	vars.Resolution = detectPrimaryResolution()

	return vars, nil
}

// detectPrimaryDisplay attempts to detect the primary display
// This is a simplified implementation - production would use wlr-randr or similar
func detectPrimaryDisplay() string {
	// Check if we're in a Wayland session
	if waylandDisplay := os.Getenv("WAYLAND_DISPLAY"); waylandDisplay != "" {
		// Try to detect using wlr-randr (if available)
		// For now, return empty - will be enhanced later
		return ""
	}

	// Check if we're in an X11 session
	if xDisplay := os.Getenv("DISPLAY"); xDisplay != "" {
		// Try to detect using xrandr (if available)
		// For now, return empty - will be enhanced later
		return ""
	}

	return ""
}

// detectPrimaryResolution attempts to detect the primary display resolution
func detectPrimaryResolution() string {
	// Similar to detectPrimaryDisplay, this would use wlr-randr or xrandr
	// For now, return empty - will be enhanced later
	return ""
}
