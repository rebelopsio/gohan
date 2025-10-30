package postinstall

// DisplayManager represents the type of display manager
type DisplayManager string

const (
	DisplayManagerSDDM DisplayManager = "sddm"
	DisplayManagerGDM  DisplayManager = "gdm"
	DisplayManagerTTY  DisplayManager = "tty"
	DisplayManagerNone DisplayManager = "none"
)

// String returns the string representation
func (d DisplayManager) String() string {
	return string(d)
}

// IsValid checks if the display manager type is valid
func (d DisplayManager) IsValid() bool {
	switch d {
	case DisplayManagerSDDM, DisplayManagerGDM, DisplayManagerTTY, DisplayManagerNone:
		return true
	default:
		return false
	}
}

// RequiresPackage returns true if the display manager requires package installation
func (d DisplayManager) RequiresPackage() bool {
	return d == DisplayManagerSDDM || d == DisplayManagerGDM
}

// Shell represents the type of shell
type Shell string

const (
	ShellZsh  Shell = "zsh"
	ShellBash Shell = "bash"
	ShellFish Shell = "fish"
)

// String returns the string representation
func (s Shell) String() string {
	return string(s)
}

// IsValid checks if the shell type is valid
func (s Shell) IsValid() bool {
	switch s {
	case ShellZsh, ShellBash, ShellFish:
		return true
	default:
		return false
	}
}

// ComponentType represents different post-installation components
type ComponentType string

const (
	ComponentDisplayManager ComponentType = "display_manager"
	ComponentShell          ComponentType = "shell"
	ComponentAudio          ComponentType = "audio"
	ComponentNetwork        ComponentType = "network"
	ComponentServices       ComponentType = "services"
	ComponentWallpaper      ComponentType = "wallpaper"
)

// String returns the string representation
func (c ComponentType) String() string {
	return string(c)
}

// SetupStatus represents the status of a component setup
type SetupStatus string

const (
	StatusPending    SetupStatus = "pending"
	StatusInProgress SetupStatus = "in_progress"
	StatusCompleted  SetupStatus = "completed"
	StatusFailed     SetupStatus = "failed"
	StatusSkipped    SetupStatus = "skipped"
)

// String returns the string representation
func (s SetupStatus) String() string {
	return string(s)
}

// IsSuccess returns true if the status indicates success
func (s SetupStatus) IsSuccess() bool {
	return s == StatusCompleted
}

// IsFailure returns true if the status indicates failure
func (s SetupStatus) IsFailure() bool {
	return s == StatusFailed
}
