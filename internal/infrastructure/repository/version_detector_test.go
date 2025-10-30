package repository_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOSRelease(t *testing.T) {
	t.Run("parses Debian Sid", func(t *testing.T) {
		content := `PRETTY_NAME="Debian GNU/Linux trixie/sid"
NAME="Debian GNU/Linux"
VERSION_CODENAME=sid
ID=debian
HOME_URL="https://www.debian.org/"
SUPPORT_URL="https://www.debian.org/support"
BUG_REPORT_URL="https://bugs.debian.org/"`

		version, err := repository.ParseOSRelease(content)
		require.NoError(t, err)
		assert.Equal(t, "sid", version.Codename)
		assert.Equal(t, "", version.VersionID) // Sid doesn't have VERSION_ID in os-release
	})

	t.Run("parses Debian Trixie", func(t *testing.T) {
		content := `PRETTY_NAME="Debian GNU/Linux trixie/sid"
NAME="Debian GNU/Linux"
VERSION_CODENAME=trixie
ID=debian
HOME_URL="https://www.debian.org/"
SUPPORT_URL="https://www.debian.org/support"
BUG_REPORT_URL="https://bugs.debian.org/"`

		version, err := repository.ParseOSRelease(content)
		require.NoError(t, err)
		assert.Equal(t, "trixie", version.Codename)
		assert.Equal(t, "", version.VersionID) // Trixie doesn't have VERSION_ID in os-release
	})

	t.Run("parses Debian Bookworm", func(t *testing.T) {
		content := `PRETTY_NAME="Debian GNU/Linux 12 (bookworm)"
NAME="Debian GNU/Linux"
VERSION_ID="12"
VERSION="12 (bookworm)"
VERSION_CODENAME=bookworm
ID=debian
HOME_URL="https://www.debian.org/"
SUPPORT_URL="https://www.debian.org/support"
BUG_REPORT_URL="https://bugs.debian.org/"`

		version, err := repository.ParseOSRelease(content)
		require.NoError(t, err)
		assert.Equal(t, "bookworm", version.Codename)
		assert.Equal(t, "12", version.VersionID)
	})

	t.Run("parses Ubuntu", func(t *testing.T) {
		content := `PRETTY_NAME="Ubuntu 22.04.3 LTS"
NAME="Ubuntu"
VERSION_ID="22.04"
VERSION="22.04.3 LTS (Jammy Jellyfish)"
VERSION_CODENAME=jammy
ID=ubuntu
ID_LIKE=debian
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"`

		version, err := repository.ParseOSRelease(content)
		require.NoError(t, err)
		assert.Equal(t, "jammy", version.Codename)
		assert.Equal(t, "22.04", version.VersionID)
	})

	t.Run("handles quoted values", func(t *testing.T) {
		content := `NAME="Debian GNU/Linux"
VERSION_ID="12"
VERSION_CODENAME=bookworm`

		version, err := repository.ParseOSRelease(content)
		require.NoError(t, err)
		assert.Equal(t, "bookworm", version.Codename)
		assert.Equal(t, "12", version.VersionID)
	})

	t.Run("handles unquoted values", func(t *testing.T) {
		content := `NAME=Debian
VERSION_ID=12
VERSION_CODENAME=bookworm`

		version, err := repository.ParseOSRelease(content)
		require.NoError(t, err)
		assert.Equal(t, "bookworm", version.Codename)
		assert.Equal(t, "12", version.VersionID)
	})

	t.Run("returns error for missing codename", func(t *testing.T) {
		content := `NAME="Debian GNU/Linux"
VERSION_ID="12"`

		_, err := repository.ParseOSRelease(content)
		assert.Error(t, err)
	})
}

func TestDetectVersionFromOSRelease(t *testing.T) {
	t.Run("detects Debian Sid", func(t *testing.T) {
		osRelease := &repository.OSReleaseInfo{
			Name:        "Debian GNU/Linux",
			VersionID:   "",
			Codename:    "sid",
			PrettyName:  "Debian GNU/Linux trixie/sid",
		}

		dv, err := repository.DetectVersionFromOSRelease(osRelease)
		require.NoError(t, err)
		assert.Equal(t, "sid", dv.Codename())
		assert.Equal(t, "unstable", dv.Version())
		assert.True(t, dv.IsSid())
	})

	t.Run("detects Debian Trixie", func(t *testing.T) {
		osRelease := &repository.OSReleaseInfo{
			Name:        "Debian GNU/Linux",
			VersionID:   "",
			Codename:    "trixie",
			PrettyName:  "Debian GNU/Linux trixie/sid",
		}

		dv, err := repository.DetectVersionFromOSRelease(osRelease)
		require.NoError(t, err)
		assert.Equal(t, "trixie", dv.Codename())
		assert.Equal(t, "testing", dv.Version())
		assert.True(t, dv.IsTrixie())
	})

	t.Run("detects Debian Bookworm", func(t *testing.T) {
		osRelease := &repository.OSReleaseInfo{
			Name:        "Debian GNU/Linux",
			VersionID:   "12",
			Codename:    "bookworm",
			PrettyName:  "Debian GNU/Linux 12 (bookworm)",
		}

		dv, err := repository.DetectVersionFromOSRelease(osRelease)
		require.NoError(t, err)
		assert.Equal(t, "bookworm", dv.Codename())
		assert.Equal(t, "12", dv.Version())
		assert.True(t, dv.IsBookworm())
	})
}
