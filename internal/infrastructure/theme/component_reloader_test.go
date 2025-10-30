package theme_test

import (
	"context"
	"testing"

	themeInfra "github.com/rebelopsio/gohan/internal/infrastructure/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// Component Reloader - TDD Unit Tests
// ========================================

func TestComponentReloader_ReloadHyprland(t *testing.T) {
	t.Run("executes hyprctl reload command", func(t *testing.T) {
		// Create reloader with mock executor
		executor := &MockCommandExecutor{
			commands: []string{},
		}
		reloader := themeInfra.NewComponentReloader(executor)

		ctx := context.Background()
		err := reloader.ReloadHyprland(ctx)

		require.NoError(t, err)
		assert.Contains(t, executor.commands, "hyprctl reload")
	})

	t.Run("returns error if hyprctl fails", func(t *testing.T) {
		executor := &MockCommandExecutor{
			shouldError: true,
			errorMsg:    "hyprctl not found",
		}
		reloader := themeInfra.NewComponentReloader(executor)

		ctx := context.Background()
		err := reloader.ReloadHyprland(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to reload Hyprland")
	})
}

func TestComponentReloader_ReloadWaybar(t *testing.T) {
	t.Run("kills and restarts waybar process", func(t *testing.T) {
		executor := &MockCommandExecutor{
			commands: []string{},
		}
		reloader := themeInfra.NewComponentReloader(executor)

		ctx := context.Background()
		err := reloader.ReloadWaybar(ctx)

		require.NoError(t, err)
		// Should kill waybar and then start it again
		assert.Contains(t, executor.commands, "killall waybar")
		assert.Contains(t, executor.commands, "waybar &")
	})

	t.Run("continues if waybar is not running", func(t *testing.T) {
		executor := &MockCommandExecutor{
			failOnCommand: "killall waybar",
			errorMsg:      "no process found",
		}
		reloader := themeInfra.NewComponentReloader(executor)

		ctx := context.Background()
		err := reloader.ReloadWaybar(ctx)

		// Should not fail if waybar is not running
		require.NoError(t, err)
	})
}

func TestComponentReloader_ReloadAll(t *testing.T) {
	t.Run("reloads all components in order", func(t *testing.T) {
		executor := &MockCommandExecutor{
			commands: []string{},
		}
		reloader := themeInfra.NewComponentReloader(executor)

		ctx := context.Background()
		err := reloader.ReloadAll(ctx)

		require.NoError(t, err)
		// Should reload Hyprland first, then Waybar
		assert.Contains(t, executor.commands, "hyprctl reload")
		assert.Contains(t, executor.commands, "killall waybar")
		assert.Contains(t, executor.commands, "waybar &")

		// Verify order: hyprctl reload should come before killall waybar
		hyprctlIndex := -1
		waybarKillIndex := -1
		for i, cmd := range executor.commands {
			if cmd == "hyprctl reload" {
				hyprctlIndex = i
			}
			if cmd == "killall waybar" {
				waybarKillIndex = i
			}
		}
		assert.True(t, hyprctlIndex < waybarKillIndex, "Hyprland should reload before Waybar")
	})

	t.Run("continues on individual component failure", func(t *testing.T) {
		executor := &MockCommandExecutor{
			failOnCommand: "killall waybar",
			errorMsg:      "waybar not running",
			commands:      []string{},
		}
		reloader := themeInfra.NewComponentReloader(executor)

		ctx := context.Background()
		err := reloader.ReloadAll(ctx)

		// Should not fail overall
		require.NoError(t, err)
		// Hyprland should still have been reloaded
		assert.Contains(t, executor.commands, "hyprctl reload")
	})
}

func TestComponentReloader_DryRun(t *testing.T) {
	t.Run("does not execute commands in dry-run mode", func(t *testing.T) {
		executor := &MockCommandExecutor{
			commands: []string{},
		}
		reloader := themeInfra.NewComponentReloaderWithOptions(executor, true) // dry-run enabled

		ctx := context.Background()
		err := reloader.ReloadAll(ctx)

		require.NoError(t, err)
		// No commands should be executed
		assert.Empty(t, executor.commands)
	})
}

// MockCommandExecutor is a mock implementation of command execution
type MockCommandExecutor struct {
	commands      []string
	shouldError   bool
	errorMsg      string
	failOnCommand string
}

func (m *MockCommandExecutor) Execute(ctx context.Context, command string, args ...string) error {
	fullCommand := command
	if len(args) > 0 {
		fullCommand = command + " " + args[0]
	}

	m.commands = append(m.commands, fullCommand)

	if m.shouldError {
		return &CommandError{message: m.errorMsg}
	}

	if m.failOnCommand != "" && fullCommand == m.failOnCommand {
		return &CommandError{message: m.errorMsg}
	}

	return nil
}

type CommandError struct {
	message string
}

func (e *CommandError) Error() string {
	return e.message
}
