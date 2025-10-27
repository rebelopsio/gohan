package history_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/history"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSystemContext(t *testing.T) {
	tests := []struct {
		name          string
		osVersion     string
		kernelVersion string
		gohanVersion  string
		hostname      string
		wantErr       error
	}{
		{
			name:          "valid context with all fields",
			osVersion:     "Debian GNU/Linux 13 (trixie)",
			kernelVersion: "6.1.0-13-amd64",
			gohanVersion:  "1.0.0",
			hostname:      "myserver",
			wantErr:       nil,
		},
		{
			name:          "valid context with minimal fields",
			osVersion:     "Debian GNU/Linux 13",
			kernelVersion: "",
			gohanVersion:  "",
			hostname      : "",
			wantErr:       nil,
		},
		{
			name:          "with whitespace trimmed",
			osVersion:     "  Debian GNU/Linux 13  ",
			kernelVersion: "  6.1.0-13-amd64  ",
			gohanVersion:  "  1.0.0  ",
			hostname:      "  myserver  ",
			wantErr:       nil,
		},
		{
			name:          "empty OS version",
			osVersion:     "",
			kernelVersion: "6.1.0-13-amd64",
			gohanVersion:  "1.0.0",
			hostname:      "myserver",
			wantErr:       history.ErrInvalidSystemContext,
		},
		{
			name:          "whitespace-only OS version",
			osVersion:     "   ",
			kernelVersion: "6.1.0-13-amd64",
			gohanVersion:  "1.0.0",
			hostname:      "myserver",
			wantErr:       history.ErrInvalidSystemContext,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := history.NewSystemContext(
				tt.osVersion,
				tt.kernelVersion,
				tt.gohanVersion,
				tt.hostname,
			)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, ctx.OSVersion())
			}
		})
	}
}

func TestSystemContext_Accessors(t *testing.T) {
	ctx, err := history.NewSystemContext(
		"Debian GNU/Linux 13 (trixie)",
		"6.1.0-13-amd64",
		"1.0.0",
		"myserver",
	)
	require.NoError(t, err)

	assert.Equal(t, "Debian GNU/Linux 13 (trixie)", ctx.OSVersion())
	assert.Equal(t, "6.1.0-13-amd64", ctx.KernelVersion())
	assert.Equal(t, "1.0.0", ctx.GohanVersion())
	assert.Equal(t, "myserver", ctx.Hostname())
}

func TestSystemContext_TrimsWhitespace(t *testing.T) {
	ctx, err := history.NewSystemContext(
		"  Debian GNU/Linux 13  ",
		"  6.1.0-13-amd64  ",
		"  1.0.0  ",
		"  myserver  ",
	)
	require.NoError(t, err)

	assert.Equal(t, "Debian GNU/Linux 13", ctx.OSVersion())
	assert.Equal(t, "6.1.0-13-amd64", ctx.KernelVersion())
	assert.Equal(t, "1.0.0", ctx.GohanVersion())
	assert.Equal(t, "myserver", ctx.Hostname())
}

func TestSystemContext_AllowsEmptyOptionalFields(t *testing.T) {
	ctx, err := history.NewSystemContext(
		"Debian GNU/Linux 13",
		"",
		"",
		"",
	)
	require.NoError(t, err)

	assert.Equal(t, "Debian GNU/Linux 13", ctx.OSVersion())
	assert.Empty(t, ctx.KernelVersion())
	assert.Empty(t, ctx.GohanVersion())
	assert.Empty(t, ctx.Hostname())
}
