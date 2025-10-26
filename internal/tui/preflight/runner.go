package preflight

import (
	"context"
	"fmt"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
	"github.com/rebelopsio/gohan/internal/infrastructure/preflight/detectors"
)

// ValidationRunner orchestrates the validation flow
type ValidationRunner struct {
	debianDetector       *detectors.DebianVersionDetector
	gpuDetector          *detectors.SystemGPUDetector
	diskSpaceDetector    *detectors.SystemDiskSpaceDetector
	connectivityChecker  *detectors.SystemConnectivityChecker
	sourceRepoChecker    *detectors.SystemSourceRepositoryChecker
	session              *preflight.ValidationSession
	progressChan         chan ProgressUpdate
}

// ProgressUpdate represents a validation progress event
type ProgressUpdate struct {
	RequirementName preflight.RequirementName
	Status          preflight.ValidationStatus
	Message         string
	Result          *preflight.ValidationResult
}

// NewValidationRunner creates a new validation runner
func NewValidationRunner() *ValidationRunner {
	return &ValidationRunner{
		debianDetector:      detectors.NewDebianVersionDetector(),
		gpuDetector:         detectors.NewSystemGPUDetector(),
		diskSpaceDetector:   detectors.NewSystemDiskSpaceDetector(),
		connectivityChecker: detectors.NewSystemConnectivityChecker(),
		sourceRepoChecker:   detectors.NewSystemSourceRepositoryChecker(),
		session:             preflight.NewValidationSession(),
		progressChan:        make(chan ProgressUpdate, 10),
	}
}

// Run executes all validation checks
func (r *ValidationRunner) Run(ctx context.Context) error {
	defer close(r.progressChan)

	// Run each validation in sequence
	validations := []func(context.Context) error{
		r.validateDebianVersion,
		r.validateGPU,
		r.validateDiskSpace,
		r.validateConnectivity,
		r.validateSourceRepositories,
	}

	for _, validate := range validations {
		if err := validate(ctx); err != nil {
			// Continue even on error - we want to complete all validations
			// Individual validation errors are captured in results
			continue
		}
	}

	// Mark session as complete
	r.session.Complete()

	return nil
}

// Session returns the validation session
func (r *ValidationRunner) Session() *preflight.ValidationSession {
	return r.session
}

// Progress returns the progress update channel
func (r *ValidationRunner) Progress() <-chan ProgressUpdate {
	return r.progressChan
}

func (r *ValidationRunner) validateDebianVersion(ctx context.Context) error {
	r.sendProgress(preflight.RequirementDebianVersion, "running", "Detecting Debian version...")

	version, err := r.debianDetector.DetectVersion(ctx)
	if err != nil {
		result := preflight.NewValidationResult(
			preflight.RequirementDebianVersion,
			preflight.StatusFail,
			preflight.SeverityCritical,
			nil,
			"Debian Sid or Trixie",
			preflight.NewUserGuidance(
				"Unable to detect Debian version",
				"Failed to read or parse /etc/os-release",
				[]string{
					"Ensure /etc/os-release exists and contains VERSION_CODENAME",
					"Verify you are running Debian Sid or Trixie",
					"Check that the system is properly configured",
				},
				"",
			),
		)
		r.session.AddResult(result)
		r.sendProgressWithResult(preflight.RequirementDebianVersion, preflight.StatusFail, "Failed to detect Debian version", &result)
		return err
	}

	if !version.IsSupported() {
		result := preflight.NewValidationResult(
			preflight.RequirementDebianVersion,
			preflight.StatusFail,
			preflight.SeverityCritical,
			version.Codename(),
			"Debian Sid or Trixie",
			preflight.NewUserGuidance(
				fmt.Sprintf("Debian %s is not supported. Gohan requires Debian Sid or Trixie.", version.Codename()),
				"This Debian version does not meet Hyprland's requirements",
				[]string{
					"Upgrade to Debian Sid (unstable) or Trixie (testing)",
					"Visit https://wiki.debian.org/DebianUnstable for upgrade instructions",
					"Backup your system before upgrading",
				},
				"https://wiki.debian.org/DebianUnstable",
			),
		)
		r.session.AddResult(result)
		r.sendProgressWithResult(preflight.RequirementDebianVersion, preflight.StatusFail, fmt.Sprintf("Unsupported version: %s", version), &result)
		return nil
	}

	result := preflight.NewValidationResult(
		preflight.RequirementDebianVersion,
		preflight.StatusPass,
		preflight.SeverityLow,
		version.Codename(),
		"Debian Sid or Trixie",
		preflight.UserGuidance{},
	)
	r.session.AddResult(result)
	r.sendProgressWithResult(preflight.RequirementDebianVersion, preflight.StatusPass, fmt.Sprintf("Detected: %s", version), &result)
	return nil
}

func (r *ValidationRunner) validateGPU(ctx context.Context) error {
	r.sendProgress(preflight.RequirementGPUSupport, "running", "Detecting GPU...")

	gpus, err := r.gpuDetector.DetectGPUs(ctx)
	if err != nil {
		result := preflight.NewValidationResult(
			preflight.RequirementGPUSupport,
			preflight.StatusWarning,
			preflight.SeverityMedium,
			nil,
			"AMD or NVIDIA GPU",
			preflight.NewUserGuidance(
				"No GPU detected. Hyprland may run with reduced performance.",
				"lspci did not detect any VGA or 3D controllers",
				[]string{
					"Hyprland can run on integrated graphics but performance will be limited",
					"Consider installing a dedicated AMD or NVIDIA GPU for best experience",
					"Check if GPU is properly seated in PCIe slot",
				},
				"",
			),
		)
		r.session.AddResult(result)
		r.sendProgressWithResult(preflight.RequirementGPUSupport, preflight.StatusWarning, "No GPU detected", &result)
		return nil
	}

	primaryGPU := gpus[0]
	result := preflight.NewValidationResult(
		preflight.RequirementGPUSupport,
		preflight.StatusPass,
		preflight.SeverityLow,
		primaryGPU.Vendor(),
		"AMD or NVIDIA GPU",
		preflight.UserGuidance{},
	)
	r.session.AddResult(result)
	r.sendProgressWithResult(preflight.RequirementGPUSupport, preflight.StatusPass, fmt.Sprintf("Detected: %s", primaryGPU), &result)
	return nil
}

func (r *ValidationRunner) validateDiskSpace(ctx context.Context) error {
	r.sendProgress(preflight.RequirementDiskSpace, "running", "Checking disk space...")

	diskSpace, err := r.diskSpaceDetector.DetectAvailableSpace(ctx, "/")
	if err != nil {
		result := preflight.NewValidationResult(
			preflight.RequirementDiskSpace,
			preflight.StatusFail,
			preflight.SeverityHigh,
			nil,
			"10 GB available",
			preflight.NewUserGuidance(
				"Unable to check disk space",
				"Failed to query filesystem statistics",
				[]string{
					"Verify filesystem is mounted correctly",
					"Check disk health with 'smartctl -a /dev/sda'",
					"Ensure at least 10 GB of free space on root partition",
				},
				"",
			),
		)
		r.session.AddResult(result)
		r.sendProgressWithResult(preflight.RequirementDiskSpace, preflight.StatusFail, "Failed to check disk space", &result)
		return err
	}

	const minSpaceGB = 10.0
	if !diskSpace.MeetsMinimum(minSpaceGB) {
		result := preflight.NewValidationResult(
			preflight.RequirementDiskSpace,
			preflight.StatusFail,
			preflight.SeverityHigh,
			fmt.Sprintf("%.2f GB", diskSpace.AvailableGB()),
			fmt.Sprintf("%.2f GB", minSpaceGB),
			preflight.NewUserGuidance(
				fmt.Sprintf("Insufficient disk space. Found %.2f GB, need %.2f GB", diskSpace.AvailableGB(), minSpaceGB),
				"Hyprland and dependencies require significant disk space",
				[]string{
					"Free up disk space by removing unnecessary files",
					"Use 'apt clean' to remove cached packages",
					"Use 'du -sh /*' to find large directories",
					"Consider resizing partitions or adding storage",
				},
				"",
			),
		)
		r.session.AddResult(result)
		r.sendProgressWithResult(preflight.RequirementDiskSpace, preflight.StatusFail, fmt.Sprintf("Only %.2f GB available", diskSpace.AvailableGB()), &result)
		return nil
	}

	result := preflight.NewValidationResult(
		preflight.RequirementDiskSpace,
		preflight.StatusPass,
		preflight.SeverityLow,
		fmt.Sprintf("%.2f GB", diskSpace.AvailableGB()),
		fmt.Sprintf("%.2f GB", minSpaceGB),
		preflight.UserGuidance{},
	)
	r.session.AddResult(result)
	r.sendProgressWithResult(preflight.RequirementDiskSpace, preflight.StatusPass, fmt.Sprintf("%.2f GB available", diskSpace.AvailableGB()), &result)
	return nil
}

func (r *ValidationRunner) validateConnectivity(ctx context.Context) error {
	r.sendProgress(preflight.RequirementInternet, "running", "Testing internet connectivity...")

	connectivity, err := r.connectivityChecker.CheckInternetConnectivity(ctx)
	if err != nil {
		result := preflight.NewValidationResult(
			preflight.RequirementInternet,
			preflight.StatusFail,
			preflight.SeverityHigh,
			nil,
			"Internet access",
			preflight.NewUserGuidance(
				"Unable to test internet connectivity",
				"Network connectivity test failed",
				[]string{
					"Check network configuration with 'ip addr' and 'ip route'",
					"Verify DNS resolution with 'ping -c 3 debian.org'",
					"Check firewall settings",
					"Ensure network cable is connected or WiFi is enabled",
				},
				"",
			),
		)
		r.session.AddResult(result)
		r.sendProgressWithResult(preflight.RequirementInternet, preflight.StatusFail, "Failed to check connectivity", &result)
		return err
	}

	if !connectivity.IsConnected() {
		result := preflight.NewValidationResult(
			preflight.RequirementInternet,
			preflight.StatusFail,
			preflight.SeverityHigh,
			"No connection",
			"Internet access",
			preflight.NewUserGuidance(
				"No internet connection detected. Cannot reach Debian repositories.",
				"All connectivity tests failed to reach Debian servers",
				[]string{
					"Check network cable or WiFi connection",
					"Verify network configuration: 'ip addr' and 'ip route'",
					"Test DNS: 'ping -c 3 debian.org'",
					"Check if proxy settings are required",
					"Temporarily disable firewall: 'systemctl stop ufw'",
				},
				"",
			),
		)
		r.session.AddResult(result)
		r.sendProgressWithResult(preflight.RequirementInternet, preflight.StatusFail, "No internet connection", &result)
		return nil
	}

	result := preflight.NewValidationResult(
		preflight.RequirementInternet,
		preflight.StatusPass,
		preflight.SeverityLow,
		"Connected",
		"Internet access",
		preflight.UserGuidance{},
	)
	r.session.AddResult(result)
	r.sendProgressWithResult(preflight.RequirementInternet, preflight.StatusPass, fmt.Sprintf("Connected (avg latency: %v)", connectivity.AverageLatency()), &result)
	return nil
}

func (r *ValidationRunner) validateSourceRepositories(ctx context.Context) error {
	r.sendProgress(preflight.RequirementSourceRepos, "running", "Checking source repositories...")

	status, err := r.sourceRepoChecker.CheckSourceRepositories(ctx)
	if err != nil {
		result := preflight.NewValidationResult(
			preflight.RequirementSourceRepos,
			preflight.StatusWarning,
			preflight.SeverityLow,
			nil,
			"deb-src configured",
			preflight.NewUserGuidance(
				"Unable to check source repositories",
				"Could not read /etc/apt/sources.list or sources.list.d/",
				[]string{
					"This is optional but recommended for building packages from source",
					"Manually verify /etc/apt/sources.list contains deb-src lines",
				},
				"",
			),
		)
		r.session.AddResult(result)
		r.sendProgressWithResult(preflight.RequirementSourceRepos, preflight.StatusWarning, "Could not check source repos", &result)
		return nil
	}

	if !status.HasDebSrc() {
		result := preflight.NewValidationResult(
			preflight.RequirementSourceRepos,
			preflight.StatusWarning,
			preflight.SeverityLow,
			"Not configured",
			"deb-src configured",
			preflight.NewUserGuidance(
				"Source repositories (deb-src) are not configured. Recommended for building packages.",
				"No deb-src lines found in apt configuration",
				[]string{
					"Edit /etc/apt/sources.list",
					"Uncomment lines starting with 'deb-src' or add them if missing",
					"Run 'apt update' after making changes",
					"This is optional but helpful for building custom packages",
				},
				"",
			),
		)
		r.session.AddResult(result)
		r.sendProgressWithResult(preflight.RequirementSourceRepos, preflight.StatusWarning, "deb-src not configured", &result)
		return nil
	}

	result := preflight.NewValidationResult(
		preflight.RequirementSourceRepos,
		preflight.StatusPass,
		preflight.SeverityLow,
		"Configured",
		"deb-src configured",
		preflight.UserGuidance{},
	)
	r.session.AddResult(result)
	r.sendProgressWithResult(preflight.RequirementSourceRepos, preflight.StatusPass, "deb-src configured", &result)
	return nil
}

func (r *ValidationRunner) sendProgress(req preflight.RequirementName, status, message string) {
	// Convert string status to ValidationStatus
	var validationStatus preflight.ValidationStatus
	switch status {
	case "running":
		validationStatus = "" // No status yet
	case "pass":
		validationStatus = preflight.StatusPass
	case "fail":
		validationStatus = preflight.StatusFail
	case "warning":
		validationStatus = preflight.StatusWarning
	}

	r.progressChan <- ProgressUpdate{
		RequirementName: req,
		Status:          validationStatus,
		Message:         message,
	}
}

func (r *ValidationRunner) sendProgressWithResult(req preflight.RequirementName, status preflight.ValidationStatus, message string, result *preflight.ValidationResult) {
	r.progressChan <- ProgressUpdate{
		RequirementName: req,
		Status:          status,
		Message:         message,
		Result:          result,
	}
}
