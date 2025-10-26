package preflight

import (
	"time"
)

// ConnectivityTest represents a single connectivity check
type ConnectivityTest struct {
	Endpoint string
	Success  bool
	Latency  time.Duration
	ErrorMsg string
}

// InternetConnectivity represents network connectivity status
type InternetConnectivity struct {
	isConnected     bool
	testedEndpoints []ConnectivityTest
	avgLatency      time.Duration
}

// NewInternetConnectivity creates a new connectivity value object
func NewInternetConnectivity(
	isConnected bool,
	testedEndpoints []ConnectivityTest,
) InternetConnectivity {
	var totalLatency time.Duration
	successCount := 0

	for _, test := range testedEndpoints {
		if test.Success {
			totalLatency += test.Latency
			successCount++
		}
	}

	avgLatency := time.Duration(0)
	if successCount > 0 {
		avgLatency = totalLatency / time.Duration(successCount)
	}

	return InternetConnectivity{
		isConnected:     isConnected,
		testedEndpoints: testedEndpoints,
		avgLatency:      avgLatency,
	}
}

// IsConnected returns true if internet is available
func (c InternetConnectivity) IsConnected() bool {
	return c.isConnected
}

// TestedEndpoints returns all tested endpoints
func (c InternetConnectivity) TestedEndpoints() []ConnectivityTest {
	return c.testedEndpoints
}

// AverageLatency returns average response time
func (c InternetConnectivity) AverageLatency() time.Duration {
	return c.avgLatency
}

// CanReachDebianRepos checks if Debian repos are accessible
func (c InternetConnectivity) CanReachDebianRepos() bool {
	for _, test := range c.testedEndpoints {
		if test.Endpoint == "deb.debian.org" || test.Endpoint == "debian.org" {
			return test.Success
		}
	}
	return c.isConnected
}

// String returns human-readable representation
func (c InternetConnectivity) String() string {
	if c.isConnected {
		return "Connected (avg latency: " + c.avgLatency.String() + ")"
	}
	return "Not connected"
}
