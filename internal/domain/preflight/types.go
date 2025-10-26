package preflight

import "time"

// ValidationStatus represents the outcome of a validation check
type ValidationStatus string

const (
	StatusPass    ValidationStatus = "pass"
	StatusFail    ValidationStatus = "fail"
	StatusWarning ValidationStatus = "warning"
)

// Severity indicates the impact level of a validation failure
type Severity string

const (
	SeverityCritical Severity = "critical" // Blocks installation
	SeverityHigh     Severity = "high"     // Blocks installation
	SeverityMedium   Severity = "medium"   // Warning only
	SeverityLow      Severity = "low"      // Informational
)

// ValidationOutcome represents the overall result
type ValidationOutcome string

const (
	OutcomeSuccess        ValidationOutcome = "success"         // All validations passed
	OutcomeBlocked        ValidationOutcome = "blocked"         // Critical failures exist
	OutcomeWarnings       ValidationOutcome = "warnings"        // Warnings exist, can proceed
	OutcomePartialSuccess ValidationOutcome = "partial_success" // Some passed, some warned
)

// RequirementName identifies specific requirements
type RequirementName string

const (
	RequirementDebianVersion RequirementName = "debian_version"
	RequirementGPUSupport    RequirementName = "gpu_support"
	RequirementDiskSpace     RequirementName = "disk_space"
	RequirementInternet      RequirementName = "internet_connectivity"
	RequirementSourceRepos   RequirementName = "source_repositories"
	RequirementDistribution  RequirementName = "distribution"
)

// GPUVendor represents GPU manufacturers
type GPUVendor string

const (
	GPUVendorAMD     GPUVendor = "amd"
	GPUVendorNVIDIA  GPUVendor = "nvidia"
	GPUVendorIntel   GPUVendor = "intel"
	GPUVendorUnknown GPUVendor = "unknown"
)

// DomainEvent is the base interface for all domain events
type DomainEvent interface {
	OccurredAt() time.Time
	EventType() string
}
