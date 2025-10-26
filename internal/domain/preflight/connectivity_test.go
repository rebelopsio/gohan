package preflight_test

import (
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
	"github.com/stretchr/testify/assert"
)

func TestNewInternetConnectivity(t *testing.T) {
	tests := []struct {
		name            string
		isConnected     bool
		testedEndpoints []preflight.ConnectivityTest
		wantAvgLatency  time.Duration
	}{
		{
			name:        "connected with single endpoint",
			isConnected: true,
			testedEndpoints: []preflight.ConnectivityTest{
				{
					Endpoint: "debian.org",
					Success:  true,
					Latency:  50 * time.Millisecond,
					ErrorMsg: "",
				},
			},
			wantAvgLatency: 50 * time.Millisecond,
		},
		{
			name:        "connected with multiple endpoints",
			isConnected: true,
			testedEndpoints: []preflight.ConnectivityTest{
				{
					Endpoint: "debian.org",
					Success:  true,
					Latency:  50 * time.Millisecond,
				},
				{
					Endpoint: "deb.debian.org",
					Success:  true,
					Latency:  100 * time.Millisecond,
				},
			},
			wantAvgLatency: 75 * time.Millisecond,
		},
		{
			name:        "not connected",
			isConnected: false,
			testedEndpoints: []preflight.ConnectivityTest{
				{
					Endpoint: "debian.org",
					Success:  false,
					Latency:  0,
					ErrorMsg: "timeout",
				},
			},
			wantAvgLatency: 0,
		},
		{
			name:        "mixed success and failure",
			isConnected: true,
			testedEndpoints: []preflight.ConnectivityTest{
				{
					Endpoint: "debian.org",
					Success:  true,
					Latency:  50 * time.Millisecond,
				},
				{
					Endpoint: "example.org",
					Success:  false,
					Latency:  0,
					ErrorMsg: "timeout",
				},
			},
			wantAvgLatency: 50 * time.Millisecond, // Only successful tests count
		},
		{
			name:            "no endpoints tested",
			isConnected:     false,
			testedEndpoints: []preflight.ConnectivityTest{},
			wantAvgLatency:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connectivity := preflight.NewInternetConnectivity(tt.isConnected, tt.testedEndpoints)

			assert.Equal(t, tt.isConnected, connectivity.IsConnected())
			assert.Equal(t, len(tt.testedEndpoints), len(connectivity.TestedEndpoints()))
			assert.Equal(t, tt.wantAvgLatency, connectivity.AverageLatency())
		})
	}
}

func TestInternetConnectivity_CanReachDebianRepos(t *testing.T) {
	tests := []struct {
		name            string
		isConnected     bool
		testedEndpoints []preflight.ConnectivityTest
		wantCanReach    bool
	}{
		{
			name:        "can reach debian.org",
			isConnected: true,
			testedEndpoints: []preflight.ConnectivityTest{
				{
					Endpoint: "debian.org",
					Success:  true,
					Latency:  50 * time.Millisecond,
				},
			},
			wantCanReach: true,
		},
		{
			name:        "can reach deb.debian.org",
			isConnected: true,
			testedEndpoints: []preflight.ConnectivityTest{
				{
					Endpoint: "deb.debian.org",
					Success:  true,
					Latency:  50 * time.Millisecond,
				},
			},
			wantCanReach: true,
		},
		{
			name:        "cannot reach debian repos",
			isConnected: true,
			testedEndpoints: []preflight.ConnectivityTest{
				{
					Endpoint: "debian.org",
					Success:  false,
					ErrorMsg: "timeout",
				},
			},
			wantCanReach: false,
		},
		{
			name:        "connected but no debian endpoints tested",
			isConnected: true,
			testedEndpoints: []preflight.ConnectivityTest{
				{
					Endpoint: "google.com",
					Success:  true,
					Latency:  30 * time.Millisecond,
				},
			},
			wantCanReach: true, // Falls back to isConnected
		},
		{
			name:            "not connected with no endpoints",
			isConnected:     false,
			testedEndpoints: []preflight.ConnectivityTest{},
			wantCanReach:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connectivity := preflight.NewInternetConnectivity(tt.isConnected, tt.testedEndpoints)
			assert.Equal(t, tt.wantCanReach, connectivity.CanReachDebianRepos())
		})
	}
}

func TestInternetConnectivity_String(t *testing.T) {
	tests := []struct {
		name            string
		isConnected     bool
		testedEndpoints []preflight.ConnectivityTest
		wantContains    string
	}{
		{
			name:        "connected shows latency",
			isConnected: true,
			testedEndpoints: []preflight.ConnectivityTest{
				{
					Endpoint: "debian.org",
					Success:  true,
					Latency:  50 * time.Millisecond,
				},
			},
			wantContains: "Connected",
		},
		{
			name:            "not connected",
			isConnected:     false,
			testedEndpoints: []preflight.ConnectivityTest{},
			wantContains:    "Not connected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connectivity := preflight.NewInternetConnectivity(tt.isConnected, tt.testedEndpoints)
			str := connectivity.String()
			assert.Contains(t, str, tt.wantContains)
		})
	}
}

func TestInternetConnectivity_AverageLatencyCalculation(t *testing.T) {
	// Test with multiple successful connections
	connectivity := preflight.NewInternetConnectivity(true, []preflight.ConnectivityTest{
		{Endpoint: "debian.org", Success: true, Latency: 30 * time.Millisecond},
		{Endpoint: "deb.debian.org", Success: true, Latency: 90 * time.Millisecond},
		{Endpoint: "security.debian.org", Success: true, Latency: 60 * time.Millisecond},
	})

	// Average should be (30 + 90 + 60) / 3 = 60ms
	assert.Equal(t, 60*time.Millisecond, connectivity.AverageLatency())
}

func TestInternetConnectivity_OnlySuccessfulTestsCountForAverage(t *testing.T) {
	connectivity := preflight.NewInternetConnectivity(true, []preflight.ConnectivityTest{
		{Endpoint: "debian.org", Success: true, Latency: 50 * time.Millisecond},
		{Endpoint: "example.com", Success: false, Latency: 0}, // Failed test
		{Endpoint: "deb.debian.org", Success: true, Latency: 100 * time.Millisecond},
	})

	// Average should only include successful tests: (50 + 100) / 2 = 75ms
	assert.Equal(t, 75*time.Millisecond, connectivity.AverageLatency())
}
