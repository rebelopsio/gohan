package repository_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepositoryConfig(t *testing.T) {
	t.Run("creates valid repository config", func(t *testing.T) {
		entries := []repository.SourceEntry{
			{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main", "contrib", "non-free", "non-free-firmware"},
			},
		}

		config, err := repository.NewRepositoryConfig(entries)
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Len(t, config.Entries(), 1)
	})

	t.Run("rejects empty entries", func(t *testing.T) {
		_, err := repository.NewRepositoryConfig([]repository.SourceEntry{})
		assert.Error(t, err)
	})
}

func TestRepositoryConfig_HasComponent(t *testing.T) {
	tests := []struct {
		name       string
		entries    []repository.SourceEntry
		component  string
		wantResult bool
	}{
		{
			name: "has main component",
			entries: []repository.SourceEntry{
				{
					Type:       "deb",
					URI:        "http://deb.debian.org/debian",
					Suite:      "sid",
					Components: []string{"main"},
				},
			},
			component:  "main",
			wantResult: true,
		},
		{
			name: "has non-free component",
			entries: []repository.SourceEntry{
				{
					Type:       "deb",
					URI:        "http://deb.debian.org/debian",
					Suite:      "sid",
					Components: []string{"main", "non-free"},
				},
			},
			component:  "non-free",
			wantResult: true,
		},
		{
			name: "missing non-free component",
			entries: []repository.SourceEntry{
				{
					Type:       "deb",
					URI:        "http://deb.debian.org/debian",
					Suite:      "sid",
					Components: []string{"main"},
				},
			},
			component:  "non-free",
			wantResult: false,
		},
		{
			name: "has non-free-firmware",
			entries: []repository.SourceEntry{
				{
					Type:       "deb",
					URI:        "http://deb.debian.org/debian",
					Suite:      "sid",
					Components: []string{"main", "non-free-firmware"},
				},
			},
			component:  "non-free-firmware",
			wantResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := repository.NewRepositoryConfig(tt.entries)
			require.NoError(t, err)

			result := config.HasComponent(tt.component)
			assert.Equal(t, tt.wantResult, result)
		})
	}
}

func TestRepositoryConfig_HasDebSrc(t *testing.T) {
	t.Run("has deb-src entries", func(t *testing.T) {
		entries := []repository.SourceEntry{
			{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main"},
			},
			{
				Type:       "deb-src",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main"},
			},
		}

		config, err := repository.NewRepositoryConfig(entries)
		require.NoError(t, err)

		assert.True(t, config.HasDebSrc())
	})

	t.Run("missing deb-src entries", func(t *testing.T) {
		entries := []repository.SourceEntry{
			{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main"},
			},
		}

		config, err := repository.NewRepositoryConfig(entries)
		require.NoError(t, err)

		assert.False(t, config.HasDebSrc())
	})
}

func TestRepositoryConfig_AddComponent(t *testing.T) {
	t.Run("adds component to all entries", func(t *testing.T) {
		entries := []repository.SourceEntry{
			{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main"},
			},
		}

		config, err := repository.NewRepositoryConfig(entries)
		require.NoError(t, err)

		err = config.AddComponent("non-free")
		require.NoError(t, err)

		assert.True(t, config.HasComponent("non-free"))
	})

	t.Run("does not add duplicate component", func(t *testing.T) {
		entries := []repository.SourceEntry{
			{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main", "non-free"},
			},
		}

		config, err := repository.NewRepositoryConfig(entries)
		require.NoError(t, err)

		err = config.AddComponent("non-free")
		require.NoError(t, err)

		// Verify it's still there and not duplicated
		assert.True(t, config.HasComponent("non-free"))
	})
}

func TestRepositoryConfig_EnableDebSrc(t *testing.T) {
	t.Run("adds deb-src entries", func(t *testing.T) {
		entries := []repository.SourceEntry{
			{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main"},
			},
		}

		config, err := repository.NewRepositoryConfig(entries)
		require.NoError(t, err)

		assert.False(t, config.HasDebSrc())

		config.EnableDebSrc()

		assert.True(t, config.HasDebSrc())
		assert.Len(t, config.Entries(), 2) // Original + deb-src
	})

	t.Run("idempotent - does not duplicate deb-src", func(t *testing.T) {
		entries := []repository.SourceEntry{
			{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main"},
			},
		}

		config, err := repository.NewRepositoryConfig(entries)
		require.NoError(t, err)

		config.EnableDebSrc()
		config.EnableDebSrc() // Call twice

		assert.True(t, config.HasDebSrc())
		assert.Len(t, config.Entries(), 2) // Should still be 2, not 3
	})
}

func TestRepositoryConfig_String(t *testing.T) {
	t.Run("formats sources.list content", func(t *testing.T) {
		entries := []repository.SourceEntry{
			{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main", "contrib"},
			},
			{
				Type:       "deb-src",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main", "contrib"},
			},
		}

		config, err := repository.NewRepositoryConfig(entries)
		require.NoError(t, err)

		output := config.String()
		assert.Contains(t, output, "deb http://deb.debian.org/debian sid main contrib")
		assert.Contains(t, output, "deb-src http://deb.debian.org/debian sid main contrib")
	})
}

func TestSourceEntry_String(t *testing.T) {
	tests := []struct {
		name  string
		entry repository.SourceEntry
		want  string
	}{
		{
			name: "basic entry",
			entry: repository.SourceEntry{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main"},
			},
			want: "deb http://deb.debian.org/debian sid main",
		},
		{
			name: "multiple components",
			entry: repository.SourceEntry{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main", "contrib", "non-free"},
			},
			want: "deb http://deb.debian.org/debian sid main contrib non-free",
		},
		{
			name: "deb-src entry",
			entry: repository.SourceEntry{
				Type:       "deb-src",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main"},
			},
			want: "deb-src http://deb.debian.org/debian sid main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.entry.String()
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestSourceEntry_Validate(t *testing.T) {
	tests := []struct {
		name    string
		entry   repository.SourceEntry
		wantErr bool
	}{
		{
			name: "valid deb entry",
			entry: repository.SourceEntry{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main"},
			},
			wantErr: false,
		},
		{
			name: "valid deb-src entry",
			entry: repository.SourceEntry{
				Type:       "deb-src",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main"},
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			entry: repository.SourceEntry{
				Type:       "rpm",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{"main"},
			},
			wantErr: true,
		},
		{
			name: "missing URI",
			entry: repository.SourceEntry{
				Type:       "deb",
				URI:        "",
				Suite:      "sid",
				Components: []string{"main"},
			},
			wantErr: true,
		},
		{
			name: "missing suite",
			entry: repository.SourceEntry{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "",
				Components: []string{"main"},
			},
			wantErr: true,
		},
		{
			name: "missing components",
			entry: repository.SourceEntry{
				Type:       "deb",
				URI:        "http://deb.debian.org/debian",
				Suite:      "sid",
				Components: []string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
