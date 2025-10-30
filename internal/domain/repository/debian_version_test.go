package repository_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDebianVersion(t *testing.T) {
	tests := []struct {
		name       string
		codename   string
		version    string
		wantErr    bool
		wantSupport bool
	}{
		{
			name:        "Debian Sid",
			codename:    "sid",
			version:     "unstable",
			wantErr:     false,
			wantSupport: true,
		},
		{
			name:        "Debian Trixie",
			codename:    "trixie",
			version:     "testing",
			wantErr:     false,
			wantSupport: true,
		},
		{
			name:        "Debian Bookworm",
			codename:    "bookworm",
			version:     "12",
			wantErr:     false,
			wantSupport: false,
		},
		{
			name:        "Ubuntu",
			codename:    "jammy",
			version:     "22.04",
			wantErr:     false,
			wantSupport: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dv, err := repository.NewDebianVersion(tt.codename, tt.version)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.codename, dv.Codename())
			assert.Equal(t, tt.version, dv.Version())
			assert.Equal(t, tt.wantSupport, dv.IsSupported())
		})
	}
}

func TestDebianVersion_IsSid(t *testing.T) {
	t.Run("returns true for Sid", func(t *testing.T) {
		dv, err := repository.NewDebianVersion("sid", "unstable")
		require.NoError(t, err)

		assert.True(t, dv.IsSid())
		assert.False(t, dv.IsTrixie())
		assert.False(t, dv.IsBookworm())
	})
}

func TestDebianVersion_IsTrixie(t *testing.T) {
	t.Run("returns true for Trixie", func(t *testing.T) {
		dv, err := repository.NewDebianVersion("trixie", "testing")
		require.NoError(t, err)

		assert.True(t, dv.IsTrixie())
		assert.False(t, dv.IsSid())
		assert.False(t, dv.IsBookworm())
	})
}

func TestDebianVersion_IsBookworm(t *testing.T) {
	t.Run("returns true for Bookworm", func(t *testing.T) {
		dv, err := repository.NewDebianVersion("bookworm", "12")
		require.NoError(t, err)

		assert.True(t, dv.IsBookworm())
		assert.False(t, dv.IsSid())
		assert.False(t, dv.IsTrixie())
	})
}

func TestDebianVersion_IsUbuntu(t *testing.T) {
	t.Run("detects Ubuntu", func(t *testing.T) {
		dv, err := repository.NewDebianVersion("jammy", "22.04")
		require.NoError(t, err)

		assert.True(t, dv.IsUbuntu())
		assert.False(t, dv.IsSid())
	})
}

func TestDebianVersion_String(t *testing.T) {
	tests := []struct {
		name     string
		codename string
		version  string
		want     string
	}{
		{
			name:     "Sid",
			codename: "sid",
			version:  "unstable",
			want:     "Debian sid (unstable)",
		},
		{
			name:     "Trixie",
			codename: "trixie",
			version:  "testing",
			want:     "Debian trixie (testing)",
		},
		{
			name:     "Bookworm",
			codename: "bookworm",
			version:  "12",
			want:     "Debian bookworm (12)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dv, err := repository.NewDebianVersion(tt.codename, tt.version)
			require.NoError(t, err)

			assert.Equal(t, tt.want, dv.String())
		})
	}
}

func TestDebianVersion_SupportMessage(t *testing.T) {
	t.Run("Sid has no warning", func(t *testing.T) {
		dv, err := repository.NewDebianVersion("sid", "unstable")
		require.NoError(t, err)

		assert.Empty(t, dv.SupportMessage())
	})

	t.Run("Trixie has warning", func(t *testing.T) {
		dv, err := repository.NewDebianVersion("trixie", "testing")
		require.NoError(t, err)

		msg := dv.SupportMessage()
		assert.NotEmpty(t, msg)
		assert.Contains(t, msg, "Trixie")
		assert.Contains(t, msg, "Sid")
	})

	t.Run("Bookworm has error", func(t *testing.T) {
		dv, err := repository.NewDebianVersion("bookworm", "12")
		require.NoError(t, err)

		msg := dv.SupportMessage()
		assert.NotEmpty(t, msg)
		assert.Contains(t, msg, "not supported")
	})
}
