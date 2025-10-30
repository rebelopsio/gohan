package preflight

import (
	"context"
	"fmt"

	"github.com/rebelopsio/gohan/internal/domain/preflight"
)

// RunPreflightRequest contains parameters for running preflight checks
type RunPreflightRequest struct {
	// ShowProgress enables progress callbacks
	ShowProgress bool
}

// RunPreflightResponse contains the result of preflight checks
type RunPreflightResponse struct {
	SessionID      string
	Passed         bool
	HasBlockers    bool
	HasWarnings    bool
	TotalChecks    int
	PassedChecks   int
	WarningChecks  int
	FailedChecks   int
	Results        []CheckResult
	OverallMessage string
}

// CheckResult represents a single check result for display
type CheckResult struct {
	Name           string
	Passed         bool
	Blocking       bool
	Message        string
	Guidance       string
	RequirementMet bool
}

// ProgressCallback is called for each validation step
type ProgressCallback func(validatorName string, result CheckResult)

// Detectors aggregates all system detectors
type Detectors struct {
	DebianDetector          preflight.DebianDetector
	GPUDetector             preflight.GPUDetector
	DiskSpaceDetector       preflight.DiskSpaceDetector
	ConnectivityChecker     preflight.ConnectivityChecker
	SourceRepositoryChecker preflight.SourceRepositoryChecker
}

// RunPreflightUseCase coordinates all preflight validations
type RunPreflightUseCase struct {
	detectors Detectors
}

// NewRunPreflightUseCase creates a new use case instance
func NewRunPreflightUseCase(detectors Detectors) *RunPreflightUseCase {
	return &RunPreflightUseCase{
		detectors: detectors,
	}
}

// Execute runs all preflight checks
func (uc *RunPreflightUseCase) Execute(ctx context.Context, req RunPreflightRequest) (*RunPreflightResponse, error) {
	// Create validators
	validators, err := uc.createValidators(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create validators: %w", err)
	}

	// Create orchestrator
	orchestrator := preflight.NewValidationOrchestrator(validators)

	// Execute validations
	var session *preflight.ValidationSession
	if req.ShowProgress {
		session = orchestrator.ExecuteValidationsWithProgress(ctx, func(name string, result preflight.ValidationResult) {
			// Progress callback handled by CLI layer
		})
	} else {
		session = orchestrator.ExecuteValidations(ctx)
	}

	// Convert to response
	return uc.buildResponse(session), nil
}

// ExecuteWithProgress runs checks with progress callbacks
func (uc *RunPreflightUseCase) ExecuteWithProgress(
	ctx context.Context,
	req RunPreflightRequest,
	progressFn ProgressCallback,
) (*RunPreflightResponse, error) {
	// Create validators
	validators, err := uc.createValidators(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create validators: %w", err)
	}

	// Create orchestrator
	orchestrator := preflight.NewValidationOrchestrator(validators)

	// Execute with progress
	session := orchestrator.ExecuteValidationsWithProgress(ctx, func(name string, result preflight.ValidationResult) {
		if progressFn != nil {
			progressFn(name, uc.convertResult(result))
		}
	})

	return uc.buildResponse(session), nil
}

func (uc *RunPreflightUseCase) createValidators(ctx context.Context) ([]preflight.Validator, error) {
	validators := make([]preflight.Validator, 0)

	// Debian Version Validator
	debianVersion, err := uc.detectors.DebianDetector.DetectVersion(ctx)
	if err == nil {
		validators = append(validators, NewDebianVersionValidator(debianVersion))
	}

	// GPU Validator
	gpu, err := uc.detectors.GPUDetector.PrimaryGPU(ctx)
	if err == nil {
		validators = append(validators, NewGPUValidator(gpu))
	}

	// Disk Space Validator
	diskSpace, err := uc.detectors.DiskSpaceDetector.DetectAvailableSpace(ctx, "/")
	if err == nil {
		validators = append(validators, NewDiskSpaceValidator(diskSpace))
	}

	// Connectivity Validator
	connectivity, err := uc.detectors.ConnectivityChecker.CheckInternetConnectivity(ctx)
	if err == nil {
		validators = append(validators, NewConnectivityValidator(connectivity))
	}

	// Source Repository Validator
	sourceRepos, err := uc.detectors.SourceRepositoryChecker.CheckSourceRepositories(ctx)
	if err == nil {
		validators = append(validators, NewSourceRepositoryValidator(sourceRepos))
	}

	if len(validators) == 0 {
		return nil, fmt.Errorf("no validators could be created")
	}

	return validators, nil
}

func (uc *RunPreflightUseCase) buildResponse(session *preflight.ValidationSession) *RunPreflightResponse {
	results := session.Results()

	response := &RunPreflightResponse{
		SessionID:   session.ID(),
		HasBlockers: session.HasBlockers(),
		TotalChecks: len(results),
		Results:     make([]CheckResult, 0, len(results)),
	}

	passedCount := 0
	warningCount := 0
	failedCount := 0

	for _, result := range results {
		checkResult := uc.convertResult(result)
		response.Results = append(response.Results, checkResult)

		if checkResult.Passed {
			passedCount++
		} else if !checkResult.Blocking {
			warningCount++
		} else {
			failedCount++
		}
	}

	response.PassedChecks = passedCount
	response.WarningChecks = warningCount
	response.FailedChecks = failedCount
	response.Passed = !response.HasBlockers
	response.HasWarnings = warningCount > 0

	// Overall message
	if response.Passed {
		if response.HasWarnings {
			response.OverallMessage = "Preflight checks passed with warnings. Installation can proceed."
		} else {
			response.OverallMessage = "All preflight checks passed! Ready to install."
		}
	} else {
		response.OverallMessage = fmt.Sprintf(
			"Preflight checks failed. %d critical issue(s) must be resolved before installation.",
			failedCount,
		)
	}

	return response
}

func (uc *RunPreflightUseCase) convertResult(result preflight.ValidationResult) CheckResult {
	return CheckResult{
		Name:           string(result.RequirementName()),
		Passed:         result.IsPassing(),
		Blocking:       result.IsBlocking(),
		Message:        result.FormatMessage(),
		Guidance:       result.Guidance().Message(),
		RequirementMet: result.IsPassing(),
	}
}

// Validator implementations that wrap domain logic

type debianVersionValidator struct {
	version preflight.DebianVersion
}

func NewDebianVersionValidator(version preflight.DebianVersion) preflight.Validator {
	return &debianVersionValidator{version: version}
}

func (v *debianVersionValidator) Name() string {
	return "Debian Version"
}

func (v *debianVersionValidator) RequirementName() preflight.RequirementName {
	return preflight.RequirementDebianVersion
}

func (v *debianVersionValidator) Validate(ctx context.Context) preflight.ValidationResult {
	// Check if version is supported
	if v.version.IsSupported() {
		return preflight.NewValidationResult(
			preflight.RequirementDebianVersion,
			preflight.StatusPass,
			preflight.SeverityLow,
			v.version.String(),
			"Debian Sid or Trixie",
			preflight.NewUserGuidance("", "", nil, ""),
		)
	}

	// Unsupported version
	var message, reason string
	var steps []string

	if v.version.IsBookworm() {
		message = "Debian Bookworm (stable) is not supported"
		reason = "Hyprland requires cutting-edge packages only available in Sid or Trixie"
		steps = []string{
			"Upgrade to Debian Sid: sudo sed -i 's/bookworm/sid/g' /etc/apt/sources.list && sudo apt update && sudo apt full-upgrade",
			"Or upgrade to Debian Trixie: sudo sed -i 's/bookworm/trixie/g' /etc/apt/sources.list && sudo apt update && sudo apt full-upgrade",
		}
	} else {
		message = fmt.Sprintf("Debian version '%s' is not supported", v.version.String())
		reason = "Only Debian Sid (unstable) and Trixie (testing) are supported"
		steps = []string{
			"Install on Debian Sid or Trixie instead",
			"See documentation: https://gohan.sh/docs/installation",
		}
	}

	guidance := preflight.NewUserGuidance(message, reason, steps, "https://gohan.sh/docs/installation")

	return preflight.NewValidationResult(
		preflight.RequirementDebianVersion,
		preflight.StatusFail,
		preflight.SeverityCritical,
		v.version.String(),
		"Debian Sid or Trixie",
		guidance,
	)
}

type gpuValidator struct {
	gpu preflight.GPUType
}

func NewGPUValidator(gpu preflight.GPUType) preflight.Validator {
	return &gpuValidator{gpu: gpu}
}

func (v *gpuValidator) Name() string {
	return "GPU Detection"
}

func (v *gpuValidator) RequirementName() preflight.RequirementName {
	return preflight.RequirementGPUSupport
}

func (v *gpuValidator) Validate(ctx context.Context) preflight.ValidationResult {
	// NVIDIA GPUs require proprietary drivers
	if v.gpu.IsNVIDIA() {
		guidance := preflight.NewUserGuidance(
			"NVIDIA GPU detected - proprietary drivers required",
			"NVIDIA GPUs need non-free repository and nvidia-driver package",
			[]string{
				"Enable non-free repositories: gohan repo enable-nonfree",
				"Update package lists: sudo apt update",
				"Install NVIDIA drivers: sudo apt install nvidia-driver",
			},
			"https://gohan.sh/docs/nvidia-setup",
		)
		return preflight.NewValidationResult(
			preflight.RequirementGPUSupport,
			preflight.StatusWarning,
			preflight.SeverityMedium,
			v.gpu.String(),
			"GPU with open-source drivers",
			guidance,
		)
	}

	// AMD/Intel GPUs work out of the box
	return preflight.NewValidationResult(
		preflight.RequirementGPUSupport,
		preflight.StatusPass,
		preflight.SeverityLow,
		v.gpu.String(),
		"GPU with open-source drivers",
		preflight.NewUserGuidance("", "", nil, ""),
	)
}

type diskSpaceValidator struct {
	space preflight.DiskSpace
}

func NewDiskSpaceValidator(space preflight.DiskSpace) preflight.Validator {
	return &diskSpaceValidator{space: space}
}

func (v *diskSpaceValidator) Name() string {
	return "Disk Space"
}

func (v *diskSpaceValidator) RequirementName() preflight.RequirementName {
	return preflight.RequirementDiskSpace
}

func (v *diskSpaceValidator) Validate(ctx context.Context) preflight.ValidationResult {
	const requiredGB uint64 = 10

	if v.space.MeetsMinimum(requiredGB) {
		return preflight.NewValidationResult(
			preflight.RequirementDiskSpace,
			preflight.StatusPass,
			preflight.SeverityLow,
			v.space.String(),
			fmt.Sprintf("%d GB minimum", requiredGB),
			preflight.NewUserGuidance("", "", nil, ""),
		)
	}

	// Insufficient disk space
	guidance := preflight.NewUserGuidance(
		fmt.Sprintf("Insufficient disk space: %.2f GB available, %d GB required", v.space.AvailableGB(), requiredGB),
		"Hyprland and its dependencies require at least 10GB of free disk space",
		[]string{
			"Free up disk space by removing unused packages: sudo apt autoremove",
			"Clean package cache: sudo apt clean",
			"Remove old files or move data to external storage",
		},
		"https://gohan.sh/docs/troubleshooting#disk-space",
	)

	return preflight.NewValidationResult(
		preflight.RequirementDiskSpace,
		preflight.StatusFail,
		preflight.SeverityHigh,
		v.space.String(),
		fmt.Sprintf("%d GB minimum", requiredGB),
		guidance,
	)
}

type connectivityValidator struct {
	connectivity preflight.InternetConnectivity
}

func NewConnectivityValidator(connectivity preflight.InternetConnectivity) preflight.Validator {
	return &connectivityValidator{connectivity: connectivity}
}

func (v *connectivityValidator) Name() string {
	return "Internet Connectivity"
}

func (v *connectivityValidator) RequirementName() preflight.RequirementName {
	return preflight.RequirementInternet
}

func (v *connectivityValidator) Validate(ctx context.Context) preflight.ValidationResult {
	if v.connectivity.IsConnected() {
		return preflight.NewValidationResult(
			preflight.RequirementInternet,
			preflight.StatusPass,
			preflight.SeverityLow,
			v.connectivity.String(),
			"Internet connection required",
			preflight.NewUserGuidance("", "", nil, ""),
		)
	}

	// No internet connection
	guidance := preflight.NewUserGuidance(
		"No internet connection detected",
		"Internet connection is required to download Hyprland and dependencies",
		[]string{
			"Check your network connection",
			"Verify network cable is connected or WiFi is enabled",
			"Check firewall settings",
			"Test connectivity: ping debian.org",
		},
		"https://gohan.sh/docs/troubleshooting#connectivity",
	)

	return preflight.NewValidationResult(
		preflight.RequirementInternet,
		preflight.StatusFail,
		preflight.SeverityCritical,
		v.connectivity.String(),
		"Internet connection required",
		guidance,
	)
}

type sourceRepositoryValidator struct {
	status preflight.SourceRepositoryStatus
}

func NewSourceRepositoryValidator(status preflight.SourceRepositoryStatus) preflight.Validator {
	return &sourceRepositoryValidator{status: status}
}

func (v *sourceRepositoryValidator) Name() string {
	return "Source Repositories"
}

func (v *sourceRepositoryValidator) RequirementName() preflight.RequirementName {
	return preflight.RequirementSourceRepos
}

func (v *sourceRepositoryValidator) Validate(ctx context.Context) preflight.ValidationResult {
	if v.status.IsEnabled() {
		return preflight.NewValidationResult(
			preflight.RequirementSourceRepos,
			preflight.StatusPass,
			preflight.SeverityLow,
			v.status.String(),
			"deb-src repositories enabled",
			preflight.NewUserGuidance("", "", nil, ""),
		)
	}

	// Source repositories not enabled
	guidance := preflight.NewUserGuidance(
		"Source repositories (deb-src) are not enabled",
		"Source repositories are needed for building some Hyprland ecosystem packages",
		[]string{
			"Enable source repositories: gohan repo enable-debsrc",
			"Or manually edit /etc/apt/sources.list to add deb-src lines",
			"Update package lists: sudo apt update",
		},
		"https://gohan.sh/docs/repository-setup#source-repos",
	)

	return preflight.NewValidationResult(
		preflight.RequirementSourceRepos,
		preflight.StatusWarning,
		preflight.SeverityMedium,
		v.status.String(),
		"deb-src repositories enabled",
		guidance,
	)
}
