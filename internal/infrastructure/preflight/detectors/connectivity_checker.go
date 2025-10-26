package detectors

import (
	"context"
	"net/http"
	"time"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
)

// SystemConnectivityChecker implements preflight.ConnectivityChecker using HTTP
type SystemConnectivityChecker struct {
	client   *http.Client
	endpoints []string
}

// NewSystemConnectivityChecker creates a new connectivity checker
func NewSystemConnectivityChecker() *SystemConnectivityChecker {
	return &SystemConnectivityChecker{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		endpoints: []string{
			"https://deb.debian.org",
			"https://debian.org",
			"https://security.debian.org",
		},
	}
}

// CheckInternetConnectivity tests internet access
func (c *SystemConnectivityChecker) CheckInternetConnectivity(ctx context.Context) (preflight.InternetConnectivity, error) {
	tests := make([]preflight.ConnectivityTest, 0, len(c.endpoints))
	hasConnection := false

	for _, endpoint := range c.endpoints {
		test := c.testEndpoint(ctx, endpoint)
		tests = append(tests, test)

		if test.Success {
			hasConnection = true
		}
	}

	return preflight.NewInternetConnectivity(hasConnection, tests), nil
}

// CheckDebianRepositories tests Debian repo access
func (c *SystemConnectivityChecker) CheckDebianRepositories(ctx context.Context) (bool, error) {
	connectivity, err := c.CheckInternetConnectivity(ctx)
	if err != nil {
		return false, err
	}

	return connectivity.CanReachDebianRepos(), nil
}

func (c *SystemConnectivityChecker) testEndpoint(ctx context.Context, endpoint string) preflight.ConnectivityTest {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "HEAD", endpoint, nil)
	if err != nil {
		return preflight.ConnectivityTest{
			Endpoint: endpoint,
			Success:  false,
			Latency:  0,
			ErrorMsg: err.Error(),
		}
	}

	resp, err := c.client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return preflight.ConnectivityTest{
			Endpoint: endpoint,
			Success:  false,
			Latency:  latency,
			ErrorMsg: err.Error(),
		}
	}
	defer resp.Body.Close()

	// Consider 2xx and 3xx as success
	success := resp.StatusCode >= 200 && resp.StatusCode < 400

	return preflight.ConnectivityTest{
		Endpoint: endpoint,
		Success:  success,
		Latency:  latency,
		ErrorMsg: "",
	}
}
