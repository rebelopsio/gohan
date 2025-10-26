package preflight_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDebianVersion(t *testing.T) {
	tests := []struct {
		name             string
		codename         string
		versionNumber    string
		wantErr          bool
		errType          error
		errMessageContains string
	}{
		{
			name:          "valid sid",
			codename:      "sid",
			versionNumber: "unstable",
			wantErr:       false,
		},
		{
			name:          "valid trixie",
			codename:      "trixie",
			versionNumber: "13",
			wantErr:       false,
		},
		{
			name:          "valid bookworm",
			codename:      "bookworm",
			versionNumber: "12",
			wantErr:       false,
		},
		{
			name:               "empty codename",
			codename:           "",
			versionNumber:      "13",
			wantErr:            true,
			errType:            preflight.ErrInvalidDebianVersion,
			errMessageContains: "codename",
		},
		{
			name:               "whitespace only codename",
			codename:           "   ",
			versionNumber:      "13",
			wantErr:            true,
			errType:            preflight.ErrInvalidDebianVersion,
			errMessageContains: "codename",
		},
		{
			name:          "codename with uppercase",
			codename:      "Sid",
			versionNumber: "unstable",
			wantErr:       false,
		},
		{
			name:          "codename with extra spaces",
			codename:      "  trixie  ",
			versionNumber: "13",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := preflight.NewDebianVersion(tt.codename, tt.versionNumber)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				if tt.errMessageContains != "" {
					assert.Contains(t, err.Error(), tt.errMessageContains,
						"Error message should contain helpful context")
				}
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, version.Codename())
			assert.Equal(t, tt.versionNumber, version.VersionNumber())
		})
	}
}

func TestDebianVersion_IsSupported(t *testing.T) {
	tests := []struct {
		name       string
		version    preflight.DebianVersion
		wantResult bool
	}{
		{
			name:       "sid is supported",
			version:    preflight.DebianSid,
			wantResult: true,
		},
		{
			name:       "trixie is supported",
			version:    preflight.DebianTrixie,
			wantResult: true,
		},
		{
			name: "bookworm is not supported",
			version: func() preflight.DebianVersion {
				v, _ := preflight.NewDebianVersion("bookworm", "12")
				return v
			}(),
			wantResult: false,
		},
		{
			name: "bullseye is not supported",
			version: func() preflight.DebianVersion {
				v, _ := preflight.NewDebianVersion("bullseye", "11")
				return v
			}(),
			wantResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantResult, tt.version.IsSupported())
		})
	}
}

func TestDebianVersion_IsSid(t *testing.T) {
	assert.True(t, preflight.DebianSid.IsSid())
	assert.False(t, preflight.DebianTrixie.IsSid())

	bookworm, _ := preflight.NewDebianVersion("bookworm", "12")
	assert.False(t, bookworm.IsSid())
}

func TestDebianVersion_IsTrixie(t *testing.T) {
	assert.True(t, preflight.DebianTrixie.IsTrixie())
	assert.False(t, preflight.DebianSid.IsTrixie())

	bookworm, _ := preflight.NewDebianVersion("bookworm", "12")
	assert.False(t, bookworm.IsTrixie())
}

func TestDebianVersion_IsBookworm(t *testing.T) {
	bookworm, _ := preflight.NewDebianVersion("bookworm", "12")
	assert.True(t, bookworm.IsBookworm())
	assert.False(t, preflight.DebianSid.IsBookworm())
	assert.False(t, preflight.DebianTrixie.IsBookworm())
}

func TestDebianVersion_String(t *testing.T) {
	tests := []struct {
		name       string
		version    preflight.DebianVersion
		wantString string
	}{
		{
			name:       "sid with version",
			version:    preflight.DebianSid,
			wantString: "sid (unstable)",
		},
		{
			name:       "trixie with version",
			version:    preflight.DebianTrixie,
			wantString: "trixie (13)",
		},
		{
			name: "codename only",
			version: func() preflight.DebianVersion {
				v, _ := preflight.NewDebianVersion("test", "")
				return v
			}(),
			wantString: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantString, tt.version.String())
		})
	}
}

func TestDebianVersion_Equals(t *testing.T) {
	sid1 := preflight.DebianSid
	sid2, _ := preflight.NewDebianVersion("sid", "unstable")
	trixie := preflight.DebianTrixie

	assert.True(t, sid1.Equals(sid2))
	assert.False(t, sid1.Equals(trixie))
}
