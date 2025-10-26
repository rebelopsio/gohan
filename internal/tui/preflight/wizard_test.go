package preflight

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	domain "github.com/rebelopsio/gohan/internal/domain/preflight"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWizard(t *testing.T) {
	wizard := NewWizard()

	assert.NotNil(t, wizard)
	assert.Equal(t, stateWelcome, wizard.state)
	assert.NotNil(t, wizard.runner)
	assert.NotNil(t, wizard.spinner)
	assert.NotNil(t, wizard.progress)
	assert.NotNil(t, wizard.ctx)
	assert.NotNil(t, wizard.cancel)
}

func TestWizard_Init(t *testing.T) {
	wizard := NewWizard()

	cmd := wizard.Init()
	assert.NotNil(t, cmd, "Init should return spinner tick command")
}

func TestWizard_InitialState(t *testing.T) {
	wizard := NewWizard()

	assert.Equal(t, stateWelcome, wizard.state, "Initial state should be welcome")
	assert.Empty(t, wizard.progress, "Progress should be empty initially")
	assert.Nil(t, wizard.err, "Error should be nil initially")
}

func TestWizard_Update_WindowSize(t *testing.T) {
	wizard := NewWizard()

	msg := tea.WindowSizeMsg{
		Width:  80,
		Height: 24,
	}

	model, _ := wizard.Update(msg)
	w := model.(*Wizard)

	assert.Equal(t, 80, w.width)
	assert.Equal(t, 24, w.height)
}

func TestWizard_Update_KeyMsg_Quit(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"quit with q", "q"},
		{"quit with ctrl+c", "ctrl+c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wizard := NewWizard()
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}

			model, cmd := wizard.Update(msg)

			assert.NotNil(t, model)
			// Should return quit command (though we can't easily test cmd == tea.Quit)
			assert.NotNil(t, cmd)
		})
	}
}

func TestWizard_View_Welcome(t *testing.T) {
	wizard := NewWizard()

	view := wizard.View()

	assert.NotEmpty(t, view, "Welcome view should not be empty")
	assert.Contains(t, view, "Pre-flight Validation", "Should contain title")
	assert.Contains(t, view, "Debian version", "Should list Debian check")
	assert.Contains(t, view, "GPU support", "Should list GPU check")
	assert.Contains(t, view, "Disk space", "Should list disk check")
	assert.Contains(t, view, "Internet connectivity", "Should list internet check")
	assert.Contains(t, view, "Source repositories", "Should list source repos check")
	assert.Contains(t, view, "Enter", "Should show Enter key hint")
}

func TestWizard_StateTransitions(t *testing.T) {
	tests := []struct {
		name          string
		initialState  wizardState
		message       tea.Msg
		expectedState wizardState
	}{
		{
			name:          "welcome to running on validation started",
			initialState:  stateWelcome,
			message:       validationStartedMsg{},
			expectedState: stateRunning,
		},
		{
			name:          "running to results on validation complete",
			initialState:  stateRunning,
			message:       validationCompleteMsg{},
			expectedState: stateResults,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wizard := NewWizard()
			wizard.state = tt.initialState

			model, _ := wizard.Update(tt.message)
			w := model.(*Wizard)

			assert.Equal(t, tt.expectedState, w.state)
		})
	}
}

func TestWizard_Update_ProgressUpdate(t *testing.T) {
	wizard := NewWizard()
	wizard.state = stateRunning

	testReq := domain.RequirementName("test_requirement")
	msg := progressUpdateMsg{
		RequirementName: testReq,
		Status:          domain.StatusPass,
		Message:         "Test message",
	}

	model, _ := wizard.Update(msg)
	w := model.(*Wizard)

	assert.Contains(t, w.progress, testReq, "Progress should be recorded")
	assert.Equal(t, "Test message", w.progress[testReq].Message)
}

func TestWizard_View_Running(t *testing.T) {
	wizard := NewWizard()
	wizard.state = stateRunning

	view := wizard.View()

	assert.NotEmpty(t, view, "Running view should not be empty")
	assert.Contains(t, view, "Running", "Should indicate validation is running")
}

func TestWizard_View_Results(t *testing.T) {
	wizard := NewWizard()
	wizard.state = stateResults

	// Need to run validation first to have results
	// For now, just test that it doesn't panic
	view := wizard.View()
	assert.NotEmpty(t, view, "Results view should not be empty")
}

func TestWizard_View_Guidance(t *testing.T) {
	wizard := NewWizard()
	wizard.state = stateGuidance

	view := wizard.View()

	assert.NotEmpty(t, view, "Guidance view should not be empty")
	assert.Contains(t, view, "Blocked", "Should indicate installation is blocked")
}

func TestWizard_ContextCancellation(t *testing.T) {
	wizard := NewWizard()

	// Wizard should have a cancellation function
	assert.NotNil(t, wizard.cancel, "Should have cancel function")

	// Cancel should not panic
	require.NotPanics(t, func() {
		wizard.cancel()
	})
}

func TestWizard_ErrorHandling(t *testing.T) {
	wizard := NewWizard()

	msg := errMsg(assert.AnError)

	model, _ := wizard.Update(msg)
	w := model.(*Wizard)

	assert.Equal(t, assert.AnError, w.err, "Error should be stored")
}

func TestWizard_AllStates_RenderWithoutPanic(t *testing.T) {
	states := []wizardState{
		stateWelcome,
		stateRunning,
		stateResults,
		stateGuidance,
		stateComplete,
	}

	for _, state := range states {
		t.Run(string(rune(state)), func(t *testing.T) {
			wizard := NewWizard()
			wizard.state = state

			require.NotPanics(t, func() {
				_ = wizard.View()
			}, "View should not panic for state %d", state)
		})
	}
}

func TestWizard_MessageTypes(t *testing.T) {
	t.Run("validationStartedMsg", func(t *testing.T) {
		msg := validationStartedMsg{}
		assert.NotNil(t, msg)
	})

	t.Run("validationCompleteMsg", func(t *testing.T) {
		msg := validationCompleteMsg{}
		assert.NotNil(t, msg)
	})

	t.Run("progressUpdateMsg", func(t *testing.T) {
		msg := progressUpdateMsg{
			RequirementName: domain.RequirementName("test"),
			Status:          domain.StatusPass,
			Message:         "Test",
		}
		assert.Equal(t, domain.RequirementName("test"), msg.RequirementName)
		assert.Equal(t, domain.StatusPass, msg.Status)
		assert.Equal(t, "Test", msg.Message)
	})

	t.Run("errMsg", func(t *testing.T) {
		msg := errMsg(assert.AnError)
		assert.Equal(t, assert.AnError, error(msg))
	})
}

func TestStatusIcon_AllStatuses(t *testing.T) {
	tests := []struct {
		status           string
		shouldNotBeEmpty bool
	}{
		{"pass", true},
		{"success", true},
		{"fail", true},
		{"error", true},
		{"warning", true},
		{"running", true},
		{"unknown", true},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			icon := StatusIcon(tt.status)
			if tt.shouldNotBeEmpty {
				assert.NotEmpty(t, icon, "Icon should not be empty for status: %s", tt.status)
			}
		})
	}
}

func TestWizard_Integration_WelcomeToResults(t *testing.T) {
	// This is a higher-level integration test
	wizard := NewWizard()

	// Start in welcome state
	assert.Equal(t, stateWelcome, wizard.state)

	// User presses Enter
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, _ := wizard.Update(keyMsg)
	wizard = model.(*Wizard)

	// Should transition through validation
	// Note: We can't easily test the full flow here without mocking,
	// but we can verify the wizard handles the messages correctly
}
