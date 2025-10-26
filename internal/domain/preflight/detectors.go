package preflight

import "context"

// DebianDetector detects Debian version and distribution
type DebianDetector interface {
	// DetectVersion identifies the Debian version
	DetectVersion(ctx context.Context) (DebianVersion, error)

	// IsDebianBased checks if system is Debian-based
	IsDebianBased(ctx context.Context) bool
}

// GPUDetector detects installed GPUs
type GPUDetector interface {
	// DetectGPUs returns all detected GPUs
	DetectGPUs(ctx context.Context) ([]GPUType, error)

	// PrimaryGPU returns the primary GPU
	PrimaryGPU(ctx context.Context) (GPUType, error)
}

// DiskSpaceDetector detects available disk space
type DiskSpaceDetector interface {
	// DetectAvailableSpace checks disk space at path
	DetectAvailableSpace(ctx context.Context, path string) (DiskSpace, error)
}

// ConnectivityChecker checks internet connectivity
type ConnectivityChecker interface {
	// CheckInternetConnectivity tests internet access
	CheckInternetConnectivity(ctx context.Context) (InternetConnectivity, error)

	// CheckDebianRepositories tests Debian repo access
	CheckDebianRepositories(ctx context.Context) (bool, error)
}

// SourceRepositoryChecker checks source repository configuration
type SourceRepositoryChecker interface {
	// CheckSourceRepositories verifies deb-src configuration
	CheckSourceRepositories(ctx context.Context) (SourceRepositoryStatus, error)
}
