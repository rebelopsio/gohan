package preflight

import (
	"strings"
)

// DebianVersion represents a Debian release
type DebianVersion struct {
	codename      string
	versionNumber string
}

var (
	// Supported versions
	DebianSid    = DebianVersion{codename: "sid", versionNumber: "unstable"}
	DebianTrixie = DebianVersion{codename: "trixie", versionNumber: "13"}
)

// NewDebianVersion creates a new Debian version value object
func NewDebianVersion(codename string, versionNumber string) (DebianVersion, error) {
	codename = strings.ToLower(strings.TrimSpace(codename))

	if codename == "" {
		return DebianVersion{}, ErrInvalidDebianVersion
	}

	return DebianVersion{
		codename:      codename,
		versionNumber: versionNumber,
	}, nil
}

// Codename returns the Debian codename
func (v DebianVersion) Codename() string {
	return v.codename
}

// VersionNumber returns the version number
func (v DebianVersion) VersionNumber() string {
	return v.versionNumber
}

// IsSupported returns true if this version is supported by Gohan
func (v DebianVersion) IsSupported() bool {
	return v.IsSid() || v.IsTrixie()
}

// IsSid returns true if this is Debian Sid
func (v DebianVersion) IsSid() bool {
	return v.codename == "sid"
}

// IsTrixie returns true if this is Debian Trixie
func (v DebianVersion) IsTrixie() bool {
	return v.codename == "trixie"
}

// IsBookworm returns true if this is Debian Bookworm
func (v DebianVersion) IsBookworm() bool {
	return v.codename == "bookworm"
}

// String returns the string representation
func (v DebianVersion) String() string {
	if v.versionNumber != "" {
		return v.codename + " (" + v.versionNumber + ")"
	}
	return v.codename
}

// Equals compares two Debian versions
func (v DebianVersion) Equals(other DebianVersion) bool {
	return v.codename == other.codename
}
