package installation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestComponentName_IsCore(t *testing.T) {
	tests := []struct {
		name      string
		component ComponentName
		want      bool
	}{
		{
			name:      "Hyprland is core",
			component: ComponentHyprland,
			want:      true,
		},
		{
			name:      "Hyprpaper is not core",
			component: ComponentHyprpaper,
			want:      false,
		},
		{
			name:      "Hyprlock is not core",
			component: ComponentHyprlock,
			want:      false,
		},
		{
			name:      "Waybar is not core",
			component: ComponentWaybar,
			want:      false,
		},
		{
			name:      "Rofi is not core",
			component: ComponentRofi,
			want:      false,
		},
		{
			name:      "Kitty is not core",
			component: ComponentKitty,
			want:      false,
		},
		{
			name:      "Default config is not core",
			component: ComponentDefaultConfig,
			want:      false,
		},
		{
			name:      "AMD driver is not core",
			component: ComponentAMDDriver,
			want:      false,
		},
		{
			name:      "NVIDIA driver is not core",
			component: ComponentNVIDIADriver,
			want:      false,
		},
		{
			name:      "Intel driver is not core",
			component: ComponentIntelDriver,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.component.IsCore()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestComponentName_IsDriver(t *testing.T) {
	tests := []struct {
		name      string
		component ComponentName
		want      bool
	}{
		{
			name:      "Hyprland is not a driver",
			component: ComponentHyprland,
			want:      false,
		},
		{
			name:      "Hyprpaper is not a driver",
			component: ComponentHyprpaper,
			want:      false,
		},
		{
			name:      "AMD driver is a driver",
			component: ComponentAMDDriver,
			want:      true,
		},
		{
			name:      "NVIDIA driver is a driver",
			component: ComponentNVIDIADriver,
			want:      true,
		},
		{
			name:      "Intel driver is a driver",
			component: ComponentIntelDriver,
			want:      true,
		},
		{
			name:      "Waybar is not a driver",
			component: ComponentWaybar,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.component.IsDriver()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestComponentName_String(t *testing.T) {
	tests := []struct {
		name      string
		component ComponentName
		want      string
	}{
		{
			name:      "Hyprland string",
			component: ComponentHyprland,
			want:      "hyprland",
		},
		{
			name:      "AMD driver string",
			component: ComponentAMDDriver,
			want:      "amd_driver",
		},
		{
			name:      "Default config string",
			component: ComponentDefaultConfig,
			want:      "default_config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.component.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInstallationStatus_IsTerminal(t *testing.T) {
	tests := []struct {
		name   string
		status InstallationStatus
		want   bool
	}{
		{
			name:   "Pending is not terminal",
			status: StatusPending,
			want:   false,
		},
		{
			name:   "Preparation is not terminal",
			status: StatusPreparation,
			want:   false,
		},
		{
			name:   "Downloading is not terminal",
			status: StatusDownloading,
			want:   false,
		},
		{
			name:   "Installing is not terminal",
			status: StatusInstalling,
			want:   false,
		},
		{
			name:   "Configuring is not terminal",
			status: StatusConfiguring,
			want:   false,
		},
		{
			name:   "Verifying is not terminal",
			status: StatusVerifying,
			want:   false,
		},
		{
			name:   "Completed is terminal",
			status: StatusCompleted,
			want:   true,
		},
		{
			name:   "Failed is terminal",
			status: StatusFailed,
			want:   true,
		},
		{
			name:   "Rolling back is not terminal",
			status: StatusRollingBack,
			want:   false,
		},
		{
			name:   "Rolled back is terminal",
			status: StatusRolledBack,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsTerminal()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInstallationStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name      string
		from      InstallationStatus
		to        InstallationStatus
		want      bool
		rationale string
	}{
		// Valid transitions
		{
			name:      "Pending to Preparation",
			from:      StatusPending,
			to:        StatusPreparation,
			want:      true,
			rationale: "Initial transition",
		},
		{
			name:      "Preparation to Downloading",
			from:      StatusPreparation,
			to:        StatusDownloading,
			want:      true,
			rationale: "Ready to download",
		},
		{
			name:      "Preparation to Installing",
			from:      StatusPreparation,
			to:        StatusInstalling,
			want:      true,
			rationale: "Can skip download if cached",
		},
		{
			name:      "Downloading to Installing",
			from:      StatusDownloading,
			to:        StatusInstalling,
			want:      true,
			rationale: "Download complete",
		},
		{
			name:      "Installing to Configuring",
			from:      StatusInstalling,
			to:        StatusConfiguring,
			want:      true,
			rationale: "Installation complete",
		},
		{
			name:      "Configuring to Verifying",
			from:      StatusConfiguring,
			to:        StatusVerifying,
			want:      true,
			rationale: "Configuration complete",
		},
		{
			name:      "Verifying to Completed",
			from:      StatusVerifying,
			to:        StatusCompleted,
			want:      true,
			rationale: "Verification successful",
		},
		{
			name:      "Rolling back to Rolled back",
			from:      StatusRollingBack,
			to:        StatusRolledBack,
			want:      true,
			rationale: "Rollback complete",
		},

		// Can always transition to Failed or RollingBack from non-terminal states
		{
			name:      "Pending to Failed",
			from:      StatusPending,
			to:        StatusFailed,
			want:      true,
			rationale: "Can fail from any state",
		},
		{
			name:      "Downloading to Failed",
			from:      StatusDownloading,
			to:        StatusFailed,
			want:      true,
			rationale: "Can fail from any state",
		},
		{
			name:      "Installing to RollingBack",
			from:      StatusInstalling,
			to:        StatusRollingBack,
			want:      true,
			rationale: "Can rollback from any state",
		},

		// Invalid transitions - skipping states
		{
			name:      "Pending to Downloading",
			from:      StatusPending,
			to:        StatusDownloading,
			want:      false,
			rationale: "Cannot skip Preparation",
		},
		{
			name:      "Pending to Installing",
			from:      StatusPending,
			to:        StatusInstalling,
			want:      false,
			rationale: "Cannot skip to Installing",
		},
		{
			name:      "Downloading to Configuring",
			from:      StatusDownloading,
			to:        StatusConfiguring,
			want:      false,
			rationale: "Cannot skip Installing",
		},
		{
			name:      "Installing to Verifying",
			from:      StatusInstalling,
			to:        StatusVerifying,
			want:      false,
			rationale: "Cannot skip Configuring",
		},

		// Invalid transitions - backwards
		{
			name:      "Installing to Downloading",
			from:      StatusInstalling,
			to:        StatusDownloading,
			want:      false,
			rationale: "Cannot go backwards",
		},
		{
			name:      "Verifying to Installing",
			from:      StatusVerifying,
			to:        StatusInstalling,
			want:      false,
			rationale: "Cannot go backwards",
		},

		// Invalid transitions - from terminal states
		{
			name:      "Completed to anything",
			from:      StatusCompleted,
			to:        StatusPending,
			want:      false,
			rationale: "Cannot transition from terminal state",
		},
		{
			name:      "Failed to Installing",
			from:      StatusFailed,
			to:        StatusInstalling,
			want:      false,
			rationale: "Cannot transition from terminal state",
		},
		{
			name:      "Rolled back to Pending",
			from:      StatusRolledBack,
			to:        StatusPending,
			want:      false,
			rationale: "Cannot transition from terminal state",
		},
		{
			name:      "Completed to Failed",
			from:      StatusCompleted,
			to:        StatusFailed,
			want:      false,
			rationale: "Cannot transition from terminal state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.from.CanTransitionTo(tt.to)
			assert.Equal(t, tt.want, got, "Rationale: %s", tt.rationale)
		})
	}
}

func TestInstallationStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status InstallationStatus
		want   string
	}{
		{
			name:   "Pending status",
			status: StatusPending,
			want:   "pending",
		},
		{
			name:   "Completed status",
			status: StatusCompleted,
			want:   "completed",
		},
		{
			name:   "Rolling back status",
			status: StatusRollingBack,
			want:   "rolling_back",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInstallationPhase_String(t *testing.T) {
	tests := []struct {
		name  string
		phase InstallationPhase
		want  string
	}{
		{
			name:  "Snapshot phase",
			phase: PhaseSnapshot,
			want:  "snapshot",
		},
		{
			name:  "Download phase",
			phase: PhaseDownload,
			want:  "download",
		},
		{
			name:  "Verification phase",
			phase: PhaseVerification,
			want:  "verification",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.phase.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResolutionAction_String(t *testing.T) {
	tests := []struct {
		name   string
		action ResolutionAction
		want   string
	}{
		{
			name:   "Remove action",
			action: ActionRemove,
			want:   "remove",
		},
		{
			name:   "Replace action",
			action: ActionReplace,
			want:   "replace",
		},
		{
			name:   "Skip action",
			action: ActionSkip,
			want:   "skip",
		},
		{
			name:   "Abort action",
			action: ActionAbort,
			want:   "abort",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.action.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

// testEvent is a test implementation of DomainEvent
type testEvent struct {
	occurredAt time.Time
	eventType  string
}

func (te *testEvent) OccurredAt() time.Time {
	return te.occurredAt
}

func (te *testEvent) EventType() string {
	return te.eventType
}

func TestDomainEvent_Interface(t *testing.T) {
	// Test that we can implement the DomainEvent interface
	te := &testEvent{
		occurredAt: time.Now(),
		eventType:  "test_event",
	}

	// This should compile if the interface is defined correctly
	var _ DomainEvent = te

	t.Run("implements DomainEvent interface", func(t *testing.T) {
		assert.NotNil(t, te)
		assert.Equal(t, "test_event", te.EventType())
		assert.False(t, te.OccurredAt().IsZero())
	})
}
