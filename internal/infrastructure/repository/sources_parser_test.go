package repository_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSourcesFile(t *testing.T) {
	t.Run("parses basic sources.list", func(t *testing.T) {
		content := `deb http://deb.debian.org/debian sid main
deb-src http://deb.debian.org/debian sid main`

		entries, err := repository.ParseSourcesFile(content)
		require.NoError(t, err)
		assert.Len(t, entries, 2)

		assert.Equal(t, "deb", entries[0].Type)
		assert.Equal(t, "http://deb.debian.org/debian", entries[0].URI)
		assert.Equal(t, "sid", entries[0].Suite)
		assert.Equal(t, []string{"main"}, entries[0].Components)

		assert.Equal(t, "deb-src", entries[1].Type)
	})

	t.Run("parses multiple components", func(t *testing.T) {
		content := `deb http://deb.debian.org/debian sid main contrib non-free non-free-firmware`

		entries, err := repository.ParseSourcesFile(content)
		require.NoError(t, err)
		assert.Len(t, entries, 1)
		assert.Equal(t, []string{"main", "contrib", "non-free", "non-free-firmware"}, entries[0].Components)
	})

	t.Run("skips comments", func(t *testing.T) {
		content := `# This is a comment
deb http://deb.debian.org/debian sid main
# Another comment`

		entries, err := repository.ParseSourcesFile(content)
		require.NoError(t, err)
		assert.Len(t, entries, 1)
	})

	t.Run("skips empty lines", func(t *testing.T) {
		content := `
deb http://deb.debian.org/debian sid main

deb-src http://deb.debian.org/debian sid main
`

		entries, err := repository.ParseSourcesFile(content)
		require.NoError(t, err)
		assert.Len(t, entries, 2)
	})

	t.Run("handles commented out entries", func(t *testing.T) {
		content := `deb http://deb.debian.org/debian sid main
#deb http://security.debian.org/ sid-security main`

		entries, err := repository.ParseSourcesFile(content)
		require.NoError(t, err)
		assert.Len(t, entries, 1)
	})

	t.Run("handles security and updates repos", func(t *testing.T) {
		content := `deb http://deb.debian.org/debian sid main
deb http://security.debian.org/ sid-security main
deb http://deb.debian.org/debian sid-updates main`

		entries, err := repository.ParseSourcesFile(content)
		require.NoError(t, err)
		assert.Len(t, entries, 3)
		assert.Equal(t, "sid-security", entries[1].Suite)
		assert.Equal(t, "sid-updates", entries[2].Suite)
	})

	t.Run("returns empty slice for empty content", func(t *testing.T) {
		entries, err := repository.ParseSourcesFile("")
		require.NoError(t, err)
		assert.Len(t, entries, 0)
	})
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    *repository.ParsedEntry
		wantErr bool
	}{
		{
			name: "valid line",
			line: "deb http://deb.debian.org/debian sid main",
			want: &repository.ParsedEntry{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main"},
			},
			wantErr: false,
		},
		{
			name: "multiple components",
			line: "deb http://deb.debian.org/debian sid main contrib non-free",
			want: &repository.ParsedEntry{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main", "contrib", "non-free"},
			},
			wantErr: false,
		},
		{
			name:    "comment line",
			line:    "# This is a comment",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty line",
			line:    "",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "whitespace only",
			line:    "   ",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "commented out entry",
			line:    "#deb http://deb.debian.org/debian sid main",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "invalid line - too few fields",
			line:    "deb http://deb.debian.org/debian",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid line - no components",
			line:    "deb http://deb.debian.org/debian sid",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repository.ParseLine(tt.line)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.want == nil {
					assert.Nil(t, got)
				} else {
					require.NotNil(t, got)
					assert.Equal(t, tt.want.Type, got.Type)
					assert.Equal(t, tt.want.URI, got.URI)
					assert.Equal(t, tt.want.Suite, got.Suite)
					assert.Equal(t, tt.want.Components, got.Components)
				}
			}
		})
	}
}
