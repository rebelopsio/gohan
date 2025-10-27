package dto

// InstallationRequest represents a request to start an installation
type InstallationRequest struct {
	// Components to install with their versions
	Components []ComponentRequest

	// GPU configuration
	GPU *GPURequest

	// Available disk space in bytes
	AvailableSpace uint64

	// Required disk space in bytes
	RequiredSpace uint64

	// Whether to merge with existing configuration
	MergeExistingConfig bool

	// Backup directory for configuration files
	BackupDirectory string
}

// ComponentRequest represents a component to install
type ComponentRequest struct {
	Name    string
	Version string

	// Optional package information
	PackageName string
	SizeBytes   uint64
}

// GPURequest represents GPU configuration
type GPURequest struct {
	Vendor         string
	RequiresDriver bool
	DriverName     string
}

// InstallationResponse represents the result of starting an installation
type InstallationResponse struct {
	SessionID   string
	Status      string
	Message     string
	StartedAt   string
	ComponentCount int
}

// InstallationProgressResponse represents installation progress
type InstallationProgressResponse struct {
	SessionID         string
	Status            string
	CurrentPhase      string
	PercentComplete   int
	Message           string
	EstimatedRemaining string
	ComponentsInstalled int
	ComponentsTotal     int
}

// InstallationCompleteResponse represents completed installation
type InstallationCompleteResponse struct {
	SessionID           string
	Status              string
	Duration            string
	ComponentsInstalled []InstalledComponentDTO
	BackupPath          string
}

// InstalledComponentDTO represents an installed component
type InstalledComponentDTO struct {
	Name       string
	Version    string
	InstalledAt string
	Verified   bool
}

// InstallationErrorResponse represents an installation error
type InstallationErrorResponse struct {
	SessionID    string
	Status       string
	Phase        string
	ErrorMessage string
	Recoverable  bool
	BackupPath   string
}

// ListInstallationsResponse represents a list of installation sessions
type ListInstallationsResponse struct {
	Sessions   []InstallationSessionSummary
	TotalCount int
}

// InstallationSessionSummary represents a summary of an installation session
type InstallationSessionSummary struct {
	SessionID           string
	Status              string
	CurrentPhase        string
	PercentComplete     int
	ComponentsInstalled int
	ComponentsTotal     int
	StartedAt           string
	CompletedAt         string
}
