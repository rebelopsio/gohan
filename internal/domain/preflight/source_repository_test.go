package preflight_test

import (
	"testing"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
	"github.com/stretchr/testify/assert"
)

func TestNewSourceRepositoryStatus(t *testing.T) {
	tests := []struct {
		name              string
		isEnabled         bool
		configuredSources []string
		wantIsEnabled     bool
		wantSourceCount   int
	}{
		{
			name:      "enabled with sources",
			isEnabled: true,
			configuredSources: []string{
				"deb http://deb.debian.org/debian sid main",
				"deb-src http://deb.debian.org/debian sid main",
			},
			wantIsEnabled:   true,
			wantSourceCount: 2,
		},
		{
			name:              "not enabled with no sources",
			isEnabled:         false,
			configuredSources: []string{},
			wantIsEnabled:     false,
			wantSourceCount:   0,
		},
		{
			name:      "enabled with only binary repos",
			isEnabled: false,
			configuredSources: []string{
				"deb http://deb.debian.org/debian sid main",
				"deb http://security.debian.org/debian-security sid-security main",
			},
			wantIsEnabled:   false,
			wantSourceCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := preflight.NewSourceRepositoryStatus(tt.isEnabled, tt.configuredSources)

			assert.Equal(t, tt.wantIsEnabled, status.IsEnabled())
			assert.Equal(t, tt.wantSourceCount, len(status.ConfiguredSources()))
		})
	}
}

func TestSourceRepositoryStatus_HasDebSrc(t *testing.T) {
	tests := []struct {
		name              string
		configuredSources []string
		wantHasDebSrc     bool
	}{
		{
			name: "has deb-src line",
			configuredSources: []string{
				"deb http://deb.debian.org/debian sid main",
				"deb-src http://deb.debian.org/debian sid main",
			},
			wantHasDebSrc: true,
		},
		{
			name: "has deb-src with leading whitespace",
			configuredSources: []string{
				"deb http://deb.debian.org/debian sid main",
				"  deb-src http://deb.debian.org/debian sid main",
			},
			wantHasDebSrc: true,
		},
		{
			name: "has deb-src with tab",
			configuredSources: []string{
				"deb http://deb.debian.org/debian sid main",
				"\tdeb-src http://deb.debian.org/debian sid main",
			},
			wantHasDebSrc: true,
		},
		{
			name: "no deb-src line",
			configuredSources: []string{
				"deb http://deb.debian.org/debian sid main",
				"deb http://security.debian.org/debian-security sid-security main",
			},
			wantHasDebSrc: false,
		},
		{
			name:              "empty sources",
			configuredSources: []string{},
			wantHasDebSrc:     false,
		},
		{
			name: "commented deb-src doesn't count",
			configuredSources: []string{
				"deb http://deb.debian.org/debian sid main",
				"# deb-src http://deb.debian.org/debian sid main",
			},
			wantHasDebSrc: false,
		},
		{
			name: "multiple deb-src lines",
			configuredSources: []string{
				"deb-src http://deb.debian.org/debian sid main",
				"deb-src http://deb.debian.org/debian sid contrib",
				"deb-src http://deb.debian.org/debian sid non-free",
			},
			wantHasDebSrc: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := preflight.NewSourceRepositoryStatus(true, tt.configuredSources)
			assert.Equal(t, tt.wantHasDebSrc, status.HasDebSrc())
		})
	}
}

func TestSourceRepositoryStatus_String(t *testing.T) {
	tests := []struct {
		name         string
		isEnabled    bool
		wantContains string
	}{
		{
			name:         "enabled",
			isEnabled:    true,
			wantContains: "enabled",
		},
		{
			name:         "not enabled",
			isEnabled:    false,
			wantContains: "not enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := preflight.NewSourceRepositoryStatus(tt.isEnabled, []string{})
			str := status.String()
			assert.Contains(t, str, tt.wantContains)
		})
	}
}

func TestSourceRepositoryStatus_ConfiguredSources(t *testing.T) {
	sources := []string{
		"deb http://deb.debian.org/debian sid main",
		"deb-src http://deb.debian.org/debian sid main",
		"deb http://security.debian.org/debian-security sid-security main",
	}

	status := preflight.NewSourceRepositoryStatus(true, sources)

	retrieved := status.ConfiguredSources()
	assert.Equal(t, len(sources), len(retrieved))
	for i, source := range sources {
		assert.Equal(t, source, retrieved[i])
	}
}

func TestSourceRepositoryStatus_NilSources(t *testing.T) {
	status := preflight.NewSourceRepositoryStatus(true, nil)

	assert.True(t, status.IsEnabled())
	assert.NotNil(t, status.ConfiguredSources(), "Should return non-nil slice")
	assert.Len(t, status.ConfiguredSources(), 0, "Should return empty slice for nil input")
	assert.False(t, status.HasDebSrc())
}

func TestSourceRepositoryStatus_EmptySources(t *testing.T) {
	status := preflight.NewSourceRepositoryStatus(false, []string{})

	assert.False(t, status.IsEnabled())
	assert.Len(t, status.ConfiguredSources(), 0)
	assert.False(t, status.HasDebSrc())
}
