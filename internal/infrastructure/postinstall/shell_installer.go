package postinstall

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rebelopsio/gohan/internal/domain/postinstall"
)

// ShellInstaller handles shell installation and configuration
type ShellInstaller struct {
	shell      postinstall.Shell
	theme      string
	packageMgr PackageManager
	homeDir    string
}

// NewShellInstaller creates a new shell installer
func NewShellInstaller(
	shell postinstall.Shell,
	theme string,
	packageMgr PackageManager,
) *ShellInstaller {
	homeDir, _ := os.UserHomeDir()
	return &ShellInstaller{
		shell:      shell,
		theme:      theme,
		packageMgr: packageMgr,
		homeDir:    homeDir,
	}
}

// Name returns the installer name
func (i *ShellInstaller) Name() string {
	return fmt.Sprintf("Shell (%s)", i.shell)
}

// Component returns the component type
func (i *ShellInstaller) Component() postinstall.ComponentType {
	return postinstall.ComponentShell
}

// Install performs the installation
func (i *ShellInstaller) Install(ctx context.Context) (postinstall.ComponentResult, error) {
	result := postinstall.NewComponentResult(
		postinstall.ComponentShell,
		postinstall.StatusInProgress,
		"Installing shell",
	)

	switch i.shell {
	case postinstall.ShellZsh:
		return i.installZsh(ctx, result)
	case postinstall.ShellBash:
		return i.installBash(ctx, result)
	case postinstall.ShellFish:
		return i.installFish(ctx, result)
	default:
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentShell,
			"Invalid shell",
			fmt.Errorf("unknown shell: %s", i.shell),
		), nil
	}
}

func (i *ShellInstaller) installZsh(ctx context.Context, result postinstall.ComponentResult) (postinstall.ComponentResult, error) {
	details := []string{}

	// Check if already installed
	installed, err := i.packageMgr.IsInstalled(ctx, "zsh")
	if err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentShell,
			"Failed to check zsh installation status",
			err,
		), err
	}

	// Install if not present
	if !installed {
		if err := i.packageMgr.Install(ctx, "zsh"); err != nil {
			return postinstall.NewComponentResultWithError(
				postinstall.ComponentShell,
				"Failed to install zsh",
				err,
			), err
		}
		details = append(details, "zsh package installed")
	} else {
		details = append(details, "zsh already installed")
	}

	// Set as default shell
	zshPath, err := exec.LookPath("zsh")
	if err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentShell,
			"Failed to locate zsh binary",
			err,
		), err
	}

	// Note: Changing default shell requires chsh which may need user interaction
	// We'll just note it for now
	details = append(details, fmt.Sprintf("zsh available at: %s", zshPath))
	details = append(details, "Run 'chsh -s $(which zsh)' to set as default")

	// Apply theme configuration if specified
	if i.theme != "" {
		if err := i.applyZshTheme(); err != nil {
			return postinstall.NewComponentResultWithError(
				postinstall.ComponentShell,
				"Failed to apply zsh theme",
				err,
			), err
		}
		details = append(details, fmt.Sprintf("Theme '%s' configured", i.theme))
	}

	return result.
		WithDetails(details...).
		Complete(postinstall.StatusCompleted), nil
}

func (i *ShellInstaller) installBash(ctx context.Context, result postinstall.ComponentResult) (postinstall.ComponentResult, error) {
	// Bash is usually already installed
	details := []string{"bash is the system default"}

	// Apply theme configuration if specified
	if i.theme != "" {
		if err := i.applyBashTheme(); err != nil {
			return postinstall.NewComponentResultWithError(
				postinstall.ComponentShell,
				"Failed to apply bash theme",
				err,
			), err
		}
		details = append(details, fmt.Sprintf("Theme '%s' configured", i.theme))
	}

	return result.
		WithDetails(details...).
		Complete(postinstall.StatusCompleted), nil
}

func (i *ShellInstaller) installFish(ctx context.Context, result postinstall.ComponentResult) (postinstall.ComponentResult, error) {
	details := []string{}

	// Check if already installed
	installed, err := i.packageMgr.IsInstalled(ctx, "fish")
	if err != nil {
		return postinstall.NewComponentResultWithError(
			postinstall.ComponentShell,
			"Failed to check fish installation status",
			err,
		), err
	}

	// Install if not present
	if !installed {
		if err := i.packageMgr.Install(ctx, "fish"); err != nil {
			return postinstall.NewComponentResultWithError(
				postinstall.ComponentShell,
				"Failed to install fish",
				err,
			), err
		}
		details = append(details, "fish package installed")
	} else {
		details = append(details, "fish already installed")
	}

	// Apply theme configuration if specified
	if i.theme != "" {
		if err := i.applyFishTheme(); err != nil {
			return postinstall.NewComponentResultWithError(
				postinstall.ComponentShell,
				"Failed to apply fish theme",
				err,
			), err
		}
		details = append(details, fmt.Sprintf("Theme '%s' configured", i.theme))
	}

	return result.
		WithDetails(details...).
		Complete(postinstall.StatusCompleted), nil
}

func (i *ShellInstaller) applyZshTheme() error {
	zshrc := filepath.Join(i.homeDir, ".zshrc")

	// Create basic .zshrc with theme
	content := fmt.Sprintf("# Gohan zsh configuration\n# Theme: %s\n\n", i.theme)

	// Append to existing or create new
	f, err := os.OpenFile(zshrc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func (i *ShellInstaller) applyBashTheme() error {
	bashrc := filepath.Join(i.homeDir, ".bashrc")

	content := fmt.Sprintf("\n# Gohan bash configuration\n# Theme: %s\n\n", i.theme)

	// Append to existing
	f, err := os.OpenFile(bashrc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func (i *ShellInstaller) applyFishTheme() error {
	fishConfig := filepath.Join(i.homeDir, ".config/fish/config.fish")

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fishConfig), 0755); err != nil {
		return err
	}

	content := fmt.Sprintf("# Gohan fish configuration\n# Theme: %s\n\n", i.theme)

	// Append to existing or create new
	f, err := os.OpenFile(fishConfig, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

// Verify checks if the shell is properly configured
func (i *ShellInstaller) Verify(ctx context.Context) (bool, error) {
	switch i.shell {
	case postinstall.ShellZsh:
		return i.packageMgr.IsInstalled(ctx, "zsh")
	case postinstall.ShellBash:
		return true, nil // Bash is always available
	case postinstall.ShellFish:
		return i.packageMgr.IsInstalled(ctx, "fish")
	default:
		return false, fmt.Errorf("unknown shell: %s", i.shell)
	}
}

// Rollback reverts the installation
func (i *ShellInstaller) Rollback(ctx context.Context) error {
	// We don't uninstall shells or remove config (too destructive)
	// Just note that theme config remains
	return nil
}
