package verification

// ComponentName represents a verifiable system component
type ComponentName string

const (
	ComponentHyprland       ComponentName = "hyprland"
	ComponentPortal         ComponentName = "portal"
	ComponentTheme          ComponentName = "theme"
	ComponentDisplayManager ComponentName = "display_manager"
	ComponentServices       ComponentName = "services"
	ComponentConfiguration  ComponentName = "configuration"
	ComponentDependencies   ComponentName = "dependencies"
	ComponentGPU            ComponentName = "gpu"
	ComponentAudio          ComponentName = "audio"
	ComponentNetwork        ComponentName = "network"
	ComponentShell          ComponentName = "shell"
	ComponentWallpaper      ComponentName = "wallpaper"
	ComponentPermissions    ComponentName = "permissions"
)

// CheckStatus represents the outcome of a verification check
type CheckStatus string

const (
	StatusPass    CheckStatus = "pass"
	StatusWarning CheckStatus = "warning"
	StatusFail    CheckStatus = "fail"
)

// CheckSeverity indicates how critical a failure is
type CheckSeverity string

const (
	SeverityCritical CheckSeverity = "critical" // Blocks functionality
	SeverityHigh     CheckSeverity = "high"     // Significant issue
	SeverityMedium   CheckSeverity = "medium"   // Minor issue
	SeverityLow      CheckSeverity = "low"      // Informational
)

// String returns the string representation of ComponentName
func (c ComponentName) String() string {
	return string(c)
}

// String returns the string representation of CheckStatus
func (s CheckStatus) String() string {
	return string(s)
}

// String returns the string representation of CheckSeverity
func (s CheckSeverity) String() string {
	return string(s)
}
