//go:build integration
// +build integration

package preflight_test

import (
	"context"
	"testing"
	"time"

	"github.com/rebelopsio/gohan/internal/infrastructure/preflight/detectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebianVersionDetector_Integration(t *testing.T) {
	detector := detectors.NewDebianVersionDetector()
	ctx := context.Background()

	t.Run("DetectVersion", func(t *testing.T) {
		version, err := detector.DetectVersion(ctx)
		require.NoError(t, err)

		assert.NotEmpty(t, version.Codename(), "Should detect Debian codename")
		t.Logf("Detected Debian version: %s", version)
	})

	t.Run("IsDebianBased", func(t *testing.T) {
		isDebian := detector.IsDebianBased(ctx)
		assert.True(t, isDebian, "Should detect Debian-based system")
	})
}

func TestSystemGPUDetector_Integration(t *testing.T) {
	detector := detectors.NewSystemGPUDetector()
	ctx := context.Background()

	t.Run("DetectGPUs", func(t *testing.T) {
		gpus, err := detector.DetectGPUs(ctx)

		if err != nil {
			t.Skipf("No GPU detected or lspci not available: %v", err)
			return
		}

		assert.NotEmpty(t, gpus, "Should detect at least one GPU")

		for i, gpu := range gpus {
			t.Logf("GPU %d: %s", i, gpu)
			assert.NotEmpty(t, gpu.Vendor(), "GPU should have vendor")
		}
	})

	t.Run("PrimaryGPU", func(t *testing.T) {
		gpu, err := detector.PrimaryGPU(ctx)

		if err != nil {
			t.Skipf("No primary GPU detected: %v", err)
			return
		}

		t.Logf("Primary GPU: %s", gpu)
		assert.NotEmpty(t, gpu.Vendor(), "Primary GPU should have vendor")
	})
}

func TestSystemDiskSpaceDetector_Integration(t *testing.T) {
	detector := detectors.NewSystemDiskSpaceDetector()
	ctx := context.Background()

	t.Run("DetectAvailableSpace for root", func(t *testing.T) {
		diskSpace, err := detector.DetectAvailableSpace(ctx, "/")
		require.NoError(t, err)

		assert.Greater(t, diskSpace.Total(), uint64(0), "Total space should be greater than 0")
		assert.LessOrEqual(t, diskSpace.Available(), diskSpace.Total(), "Available should be <= Total")

		t.Logf("Disk space at /: %s", diskSpace)
		t.Logf("Available: %.2f GB, Total: %.2f GB, Usage: %.1f%%",
			diskSpace.AvailableGB(), diskSpace.TotalGB(), diskSpace.UsagePercent())
	})

	t.Run("DetectAvailableSpace for /tmp", func(t *testing.T) {
		diskSpace, err := detector.DetectAvailableSpace(ctx, "/tmp")
		require.NoError(t, err)

		assert.Greater(t, diskSpace.Total(), uint64(0), "Total space should be greater than 0")
		t.Logf("Disk space at /tmp: %s", diskSpace)
	})

	t.Run("DetectAvailableSpace with empty path defaults to root", func(t *testing.T) {
		diskSpace, err := detector.DetectAvailableSpace(ctx, "")
		require.NoError(t, err)

		assert.Equal(t, "/", diskSpace.Path(), "Empty path should default to /")
		assert.Greater(t, diskSpace.Total(), uint64(0), "Total space should be greater than 0")
	})
}

func TestSystemConnectivityChecker_Integration(t *testing.T) {
	detector := detectors.NewSystemConnectivityChecker()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("CheckInternetConnectivity", func(t *testing.T) {
		connectivity, err := detector.CheckInternetConnectivity(ctx)
		require.NoError(t, err)

		t.Logf("Internet connectivity: %s", connectivity)
		t.Logf("Is connected: %v", connectivity.IsConnected())
		t.Logf("Average latency: %v", connectivity.AverageLatency())

		endpoints := connectivity.TestedEndpoints()
		assert.NotEmpty(t, endpoints, "Should test at least one endpoint")

		for _, test := range endpoints {
			t.Logf("Endpoint %s: Success=%v, Latency=%v, Error=%s",
				test.Endpoint, test.Success, test.Latency, test.ErrorMsg)
		}

		// Note: This may fail in environments without internet
		if !connectivity.IsConnected() {
			t.Log("WARNING: No internet connectivity detected. This may be expected in isolated test environments.")
		}
	})

	t.Run("CheckDebianRepositories", func(t *testing.T) {
		canReach, err := detector.CheckDebianRepositories(ctx)
		require.NoError(t, err)

		t.Logf("Can reach Debian repositories: %v", canReach)

		// Note: This may fail in environments without internet
		if !canReach {
			t.Log("WARNING: Cannot reach Debian repositories. This may be expected in isolated test environments.")
		}
	})
}

func TestSystemSourceRepositoryChecker_Integration(t *testing.T) {
	checker := detectors.NewSystemSourceRepositoryChecker()
	ctx := context.Background()

	t.Run("CheckSourceRepositories", func(t *testing.T) {
		status, err := checker.CheckSourceRepositories(ctx)
		require.NoError(t, err)

		t.Logf("Source repositories status: %s", status)
		t.Logf("Is enabled: %v", status.IsEnabled())
		t.Logf("Has deb-src: %v", status.HasDebSrc())

		sources := status.ConfiguredSources()
		t.Logf("Found %d source lines", len(sources))

		if len(sources) > 0 {
			t.Log("Sample sources:")
			for i, source := range sources {
				if i >= 5 {
					t.Logf("... and %d more", len(sources)-5)
					break
				}
				t.Logf("  %s", source)
			}
		}

		// Check if deb-src is configured
		if !status.HasDebSrc() {
			t.Log("NOTE: No deb-src repositories found. This is common on minimal installations.")
		}
	})
}

func TestAllDetectors_RealWorldScenario_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	t.Log("=== Running Complete Pre-flight Check Simulation ===")

	// 1. Debian Version
	debianDetector := detectors.NewDebianVersionDetector()
	version, err := debianDetector.DetectVersion(ctx)
	if err != nil {
		t.Errorf("Failed to detect Debian version: %v", err)
	} else {
		t.Logf("✓ Debian Version: %s (Supported: %v)", version, version.IsSupported())
	}

	// 2. GPU Detection
	gpuDetector := detectors.NewSystemGPUDetector()
	gpus, err := gpuDetector.DetectGPUs(ctx)
	if err != nil {
		t.Logf("⚠ GPU Detection: %v (may be expected in VM/container)", err)
	} else {
		t.Logf("✓ Detected %d GPU(s):", len(gpus))
		for i, gpu := range gpus {
			t.Logf("  GPU %d: %s (Special config needed: %v)",
				i, gpu, gpu.RequiresSpecialConfiguration())
		}
	}

	// 3. Disk Space
	diskDetector := detectors.NewSystemDiskSpaceDetector()
	diskSpace, err := diskDetector.DetectAvailableSpace(ctx, "/")
	if err != nil {
		t.Errorf("Failed to detect disk space: %v", err)
	} else {
		meetsRequirement := diskSpace.MeetsMinimum(10)
		t.Logf("✓ Disk Space: %.2f GB available / %.2f GB total (Meets 10GB requirement: %v)",
			diskSpace.AvailableGB(), diskSpace.TotalGB(), meetsRequirement)
	}

	// 4. Internet Connectivity
	connChecker := detectors.NewSystemConnectivityChecker()
	connectivity, err := connChecker.CheckInternetConnectivity(ctx)
	if err != nil {
		t.Errorf("Failed to check connectivity: %v", err)
	} else {
		t.Logf("✓ Internet Connectivity: %v (Avg latency: %v, Can reach Debian repos: %v)",
			connectivity.IsConnected(), connectivity.AverageLatency(), connectivity.CanReachDebianRepos())
	}

	// 5. Source Repositories
	srcChecker := detectors.NewSystemSourceRepositoryChecker()
	srcStatus, err := srcChecker.CheckSourceRepositories(ctx)
	if err != nil {
		t.Errorf("Failed to check source repositories: %v", err)
	} else {
		t.Logf("✓ Source Repositories: deb-src enabled: %v (%d total sources)",
			srcStatus.HasDebSrc(), len(srcStatus.ConfiguredSources()))
	}

	t.Log("=== Pre-flight Check Simulation Complete ===")
}
