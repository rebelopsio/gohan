package templates_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rebelopsio/gohan/internal/infrastructure/installation/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// Phase 3.2: Template Engine Tests (TDD)
// ========================================

func TestTemplateEngine_ProcessTemplate(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		vars     templates.TemplateVars
		expected string
		wantErr  bool
	}{
		{
			name:     "no variables",
			content:  "Hello World",
			vars:     templates.TemplateVars{},
			expected: "Hello World",
			wantErr:  false,
		},
		{
			name:    "single variable substitution",
			content: "Hello {{username}}",
			vars: templates.TemplateVars{
				"username": "alice",
			},
			expected: "Hello alice",
			wantErr:  false,
		},
		{
			name:    "multiple variable substitutions",
			content: "User: {{username}}, Home: {{home}}",
			vars: templates.TemplateVars{
				"username": "alice",
				"home":     "/home/alice",
			},
			expected: "User: alice, Home: /home/alice",
			wantErr:  false,
		},
		{
			name:    "all supported variables",
			content: "{{username}} {{home}} {{config_dir}} {{hostname}} {{display}} {{resolution}}",
			vars: templates.TemplateVars{
				"username":   "bob",
				"home":       "/home/bob",
				"config_dir": "/home/bob/.config",
				"hostname":   "debian-box",
				"display":    "eDP-1",
				"resolution": "1920x1080",
			},
			expected: "bob /home/bob /home/bob/.config debian-box eDP-1 1920x1080",
			wantErr:  false,
		},
		{
			name:    "variable in configuration file",
			content: "exec-once = waybar\nmonitor = {{display}},{{resolution}},auto,1\n",
			vars: templates.TemplateVars{
				"display":    "HDMI-A-1",
				"resolution": "2560x1440",
			},
			expected: "exec-once = waybar\nmonitor = HDMI-A-1,2560x1440,auto,1\n",
			wantErr:  false,
		},
		{
			name:    "repeated variable substitution",
			content: "{{username}} likes {{username}}",
			vars: templates.TemplateVars{
				"username": "alice",
			},
			expected: "alice likes alice",
			wantErr:  false,
		},
		{
			name:     "empty template",
			content:  "",
			vars:     templates.TemplateVars{},
			expected: "",
			wantErr:  false,
		},
		{
			name:    "undefined variable - left as-is for visibility",
			content: "Hello {{username}}",
			vars:    templates.TemplateVars{}, // Variable not in map
			expected: "Hello {{username}}",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := templates.NewTemplateEngine()

			result, err := engine.ProcessTemplate(tt.content, tt.vars)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateEngine_ProcessFile(t *testing.T) {
	t.Run("processes template file and writes output", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create source template file
		srcPath := filepath.Join(tmpDir, "hyprland.conf.tmpl")
		templateContent := `# Hyprland config for {{username}}
$mainMod = SUPER

exec-once = waybar
monitor = {{display}},{{resolution}},auto,1

# User home: {{home}}
# Config dir: {{config_dir}}
`
		err := os.WriteFile(srcPath, []byte(templateContent), 0644)
		require.NoError(t, err)

		// Process template
		dstPath := filepath.Join(tmpDir, "hyprland.conf")
		engine := templates.NewTemplateEngine()
		vars := templates.TemplateVars{
			"username":   "testuser",
			"home":       "/home/testuser",
			"config_dir": "/home/testuser/.config",
			"display":    "eDP-1",
			"resolution": "1920x1080",
		}

		err = engine.ProcessFile(srcPath, dstPath, vars)
		require.NoError(t, err)

		// Verify output
		content, err := os.ReadFile(dstPath)
		require.NoError(t, err)

		expected := `# Hyprland config for testuser
$mainMod = SUPER

exec-once = waybar
monitor = eDP-1,1920x1080,auto,1

# User home: /home/testuser
# Config dir: /home/testuser/.config
`
		assert.Equal(t, expected, string(content))
	})

	t.Run("handles non-existent source file", func(t *testing.T) {
		tmpDir := t.TempDir()

		engine := templates.NewTemplateEngine()
		vars := templates.TemplateVars{"username": "test"}

		srcPath := filepath.Join(tmpDir, "nonexistent.tmpl")
		dstPath := filepath.Join(tmpDir, "output.conf")

		err := engine.ProcessFile(srcPath, dstPath, vars)
		assert.Error(t, err)
	})

	t.Run("creates destination directory if needed", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create source
		srcPath := filepath.Join(tmpDir, "source.tmpl")
		err := os.WriteFile(srcPath, []byte("Hello {{username}}"), 0644)
		require.NoError(t, err)

		// Destination in nested directory that doesn't exist
		dstPath := filepath.Join(tmpDir, "nested", "deep", "output.conf")

		engine := templates.NewTemplateEngine()
		vars := templates.TemplateVars{"username": "alice"}

		err = engine.ProcessFile(srcPath, dstPath, vars)
		require.NoError(t, err)

		// Verify file was created
		content, err := os.ReadFile(dstPath)
		require.NoError(t, err)
		assert.Equal(t, "Hello alice", string(content))
	})
}

func TestTemplateVars_CollectFromSystem(t *testing.T) {
	t.Run("collects system variables", func(t *testing.T) {
		vars, err := templates.CollectSystemVars()

		require.NoError(t, err)
		assert.NotEmpty(t, vars["username"], "Should detect current username")
		assert.NotEmpty(t, vars["home"], "Should detect home directory")
		assert.NotEmpty(t, vars["config_dir"], "Should set config directory")
		assert.NotEmpty(t, vars["hostname"], "Should detect hostname")

		// Display and Resolution might be empty in non-graphical environment
		// Just verify the function doesn't error
	})

	t.Run("config dir is home/.config", func(t *testing.T) {
		vars, err := templates.CollectSystemVars()

		require.NoError(t, err)
		expected := filepath.Join(vars["home"], ".config")
		assert.Equal(t, expected, vars["config_dir"])
	})
}

func TestTemplateEngine_RealWorldUsage(t *testing.T) {
	t.Run("processes complete Hyprland config", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Realistic Hyprland config template
		templateContent := `# Generated Hyprland Configuration
# User: {{username}}@{{hostname}}

$mainMod = SUPER
$terminal = kitty
$fileManager = thunar
$menu = fuzzel

# Monitors
monitor = {{display}},{{resolution}},auto,1
monitor = ,preferred,auto,1

# Environment
env = XDG_CONFIG_HOME,{{config_dir}}
env = HOME,{{home}}

# Execute at launch
exec-once = waybar
exec-once = mako

# Keybindings
bind = $mainMod, Return, exec, $terminal
bind = $mainMod SHIFT, Q, killactive
bind = $mainMod, Space, exec, $menu

# Window rules
windowrulev2 = workspace 1, class:^(kitty)$
`

		srcPath := filepath.Join(tmpDir, "hyprland.conf.tmpl")
		err := os.WriteFile(srcPath, []byte(templateContent), 0644)
		require.NoError(t, err)

		// Process
		engine := templates.NewTemplateEngine()
		vars := templates.TemplateVars{
			"username":   "developer",
			"home":       "/home/developer",
			"config_dir": "/home/developer/.config",
			"hostname":   "workstation",
			"display":    "DP-1",
			"resolution": "3840x2160",
		}

		dstPath := filepath.Join(tmpDir, "hyprland.conf")
		err = engine.ProcessFile(srcPath, dstPath, vars)
		require.NoError(t, err)

		// Verify all variables were substituted
		content, err := os.ReadFile(dstPath)
		require.NoError(t, err)

		assert.Contains(t, string(content), "developer@workstation")
		assert.Contains(t, string(content), "DP-1,3840x2160")
		assert.Contains(t, string(content), "/home/developer/.config")
		assert.NotContains(t, string(content), "{{", "Should not contain any unprocessed templates")
	})
}

func TestTemplateEngine_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		vars     templates.TemplateVars
		expected string
	}{
		{
			name:     "braces without variable name",
			content:  "Hello {{}} World",
			vars:     templates.TemplateVars{},
			expected: "Hello {{}} World",
		},
		{
			name:     "single brace",
			content:  "Hello { World",
			vars:     templates.TemplateVars{},
			expected: "Hello { World",
		},
		{
			name:    "variable at start",
			content: "{{username}} is here",
			vars:    templates.TemplateVars{"username": "alice"},
			expected: "alice is here",
		},
		{
			name:    "variable at end",
			content: "Hello {{username}}",
			vars:    templates.TemplateVars{"username": "bob"},
			expected: "Hello bob",
		},
		{
			name:     "unknown variable - left as-is for visibility",
			content:  "Hello {{unknown_var}}",
			vars:     templates.TemplateVars{},
			expected: "Hello {{unknown_var}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := templates.NewTemplateEngine()
			result, err := engine.ProcessTemplate(tt.content, tt.vars)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
