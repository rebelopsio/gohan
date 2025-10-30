package repository

import (
	"fmt"
)

// DebianVersion represents a Debian-based distribution version
type DebianVersion struct {
	codename string
	version  string
}

// NewDebianVersion creates a new DebianVersion instance
func NewDebianVersion(codename, version string) (*DebianVersion, error) {
	return &DebianVersion{
		codename: codename,
		version:  version,
	}, nil
}

// Codename returns the distribution codename (e.g., "sid", "trixie", "bookworm")
func (dv *DebianVersion) Codename() string {
	return dv.codename
}

// Version returns the version string (e.g., "unstable", "testing", "12")
func (dv *DebianVersion) Version() string {
	return dv.version
}

// IsSupported returns true if this version is supported for Hyprland installation
// Only Debian Sid and Trixie are fully supported
func (dv *DebianVersion) IsSupported() bool {
	return dv.IsSid() || dv.IsTrixie()
}

// IsSid returns true if this is Debian Sid (unstable)
func (dv *DebianVersion) IsSid() bool {
	return dv.codename == "sid"
}

// IsTrixie returns true if this is Debian Trixie (testing)
func (dv *DebianVersion) IsTrixie() bool {
	return dv.codename == "trixie"
}

// IsBookworm returns true if this is Debian Bookworm (stable)
func (dv *DebianVersion) IsBookworm() bool {
	return dv.codename == "bookworm"
}

// IsUbuntu returns true if this is an Ubuntu-based distribution
func (dv *DebianVersion) IsUbuntu() bool {
	// Ubuntu codenames include: jammy, focal, noble, etc.
	// We can detect Ubuntu by checking if it's not a known Debian version
	// and if the version looks like Ubuntu versioning (e.g., "22.04", "20.04")
	if dv.IsSid() || dv.IsTrixie() || dv.IsBookworm() {
		return false
	}

	// Common Ubuntu codenames
	ubuntuCodenames := []string{"jammy", "focal", "noble", "mantic", "lunar", "kinetic", "impish", "hirsute", "groovy", "bionic", "xenial"}
	for _, name := range ubuntuCodenames {
		if dv.codename == name {
			return true
		}
	}

	return false
}

// String returns a human-readable representation
func (dv *DebianVersion) String() string {
	return fmt.Sprintf("Debian %s (%s)", dv.codename, dv.version)
}

// SupportMessage returns a message about support status
// Returns empty string if fully supported (Sid)
// Returns warning for partial support (Trixie)
// Returns error message for unsupported versions
func (dv *DebianVersion) SupportMessage() string {
	if dv.IsSid() {
		return ""
	}

	if dv.IsTrixie() {
		return "⚠️  Warning: Trixie (testing) may have outdated Hyprland packages.\n" +
			"For the latest Hyprland features, consider using Sid (unstable)."
	}

	if dv.IsBookworm() {
		return "❌ Error: Debian Bookworm (stable) is not supported.\n" +
			"Hyprland requires newer packages available in Sid or Trixie.\n" +
			"Please upgrade to Debian Sid (unstable) or Trixie (testing)."
	}

	if dv.IsUbuntu() {
		return "❌ Error: Ubuntu is not supported.\n" +
			"This installer is designed for Debian Sid or Trixie.\n" +
			"For Ubuntu, please use the official Hyprland installation methods."
	}

	return "❌ Error: This distribution is not supported.\n" +
		"Only Debian Sid (unstable) and Trixie (testing) are supported."
}
