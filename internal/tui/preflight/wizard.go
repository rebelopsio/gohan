package preflight

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rebelopsio/gohan/internal/domain/preflight"
)

// wizardState represents the current state of the wizard
type wizardState int

const (
	stateWelcome wizardState = iota
	stateRunning
	stateResults
	stateGuidance
	stateComplete
)

// Wizard is the main Bubble Tea model for preflight validation
type Wizard struct {
	state         wizardState
	runner        *ValidationRunner
	spinner       spinner.Model
	width         int
	height        int
	progress      map[preflight.RequirementName]ProgressUpdate
	currentCheck  preflight.RequirementName
	err           error
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewWizard creates a new preflight validation wizard
func NewWizard() *Wizard {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	ctx, cancel := context.WithCancel(context.Background())

	return &Wizard{
		state:    stateWelcome,
		runner:   NewValidationRunner(),
		spinner:  s,
		progress: make(map[preflight.RequirementName]ProgressUpdate),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Init initializes the wizard
func (w *Wizard) Init() tea.Cmd {
	return w.spinner.Tick
}

// Update handles messages
func (w *Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			w.cancel()
			return w, tea.Quit
		case "enter":
			return w.handleEnter()
		}

	case tea.WindowSizeMsg:
		w.width = msg.Width
		w.height = msg.Height
		return w, nil

	case spinner.TickMsg:
		if w.state == stateRunning {
			var cmd tea.Cmd
			w.spinner, cmd = w.spinner.Update(msg)
			return w, cmd
		}

	case validationStartedMsg:
		w.state = stateRunning
		return w, tea.Batch(w.spinner.Tick, w.waitForProgress())

	case progressUpdateMsg:
		update := ProgressUpdate(msg)
		w.progress[update.RequirementName] = update
		w.currentCheck = update.RequirementName
		return w, w.waitForProgress()

	case validationCompleteMsg:
		w.state = stateResults
		return w, nil

	case errMsg:
		w.err = error(msg)
		return w, nil
	}

	return w, nil
}

// View renders the wizard
func (w *Wizard) View() string {
	switch w.state {
	case stateWelcome:
		return w.viewWelcome()
	case stateRunning:
		return w.viewRunning()
	case stateResults:
		return w.viewResults()
	case stateGuidance:
		return w.viewGuidance()
	default:
		return ""
	}
}

func (w *Wizard) handleEnter() (tea.Model, tea.Cmd) {
	switch w.state {
	case stateWelcome:
		return w, w.startValidation
	case stateResults:
		session := w.runner.Session()
		if session.CanProceed() {
			w.state = stateComplete
			return w, tea.Quit
		}
		// Show guidance for failures
		w.state = stateGuidance
		return w, nil
	case stateGuidance:
		return w, tea.Quit
	default:
		return w, nil
	}
}

func (w *Wizard) viewWelcome() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Gohan Pre-flight Validation"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Checking system requirements before installation"))
	b.WriteString("\n\n")

	checks := []string{
		"✓ Debian version (Sid or Trixie required)",
		"✓ GPU support (AMD/NVIDIA recommended)",
		"✓ Disk space (minimum 10 GB)",
		"✓ Internet connectivity",
		"✓ Source repositories (recommended)",
	}

	b.WriteString(labelStyle.Render("The following checks will be performed:"))
	b.WriteString("\n\n")
	for _, check := range checks {
		b.WriteString(progressItemStyle.Render(check))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press Enter to begin • q to quit"))

	return boxStyle.Render(b.String())
}

func (w *Wizard) viewRunning() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Running Pre-flight Validation"))
	b.WriteString("\n\n")

	requirements := []preflight.RequirementName{
		preflight.RequirementDebianVersion,
		preflight.RequirementGPUSupport,
		preflight.RequirementDiskSpace,
		preflight.RequirementInternet,
		preflight.RequirementSourceRepos,
	}

	for _, req := range requirements {
		update, exists := w.progress[req]

		if !exists {
			// Not started yet
			icon := labelStyle.Render("·")
			label := labelStyle.Render(string(req))
			b.WriteString(fmt.Sprintf("%s %s\n", icon, label))
			continue
		}

		var icon string
		var style lipgloss.Style
		var label string

		switch update.Status {
		case preflight.StatusPass:
			icon = StatusIcon("pass")
			style = progressItemDoneStyle
			label = fmt.Sprintf("%s: %s", req, update.Message)
		case preflight.StatusFail:
			icon = StatusIcon("fail")
			style = progressItemFailedStyle
			label = fmt.Sprintf("%s: %s", req, update.Message)
		case preflight.StatusWarning:
			icon = StatusIcon("warning")
			style = progressItemDoneStyle
			label = fmt.Sprintf("%s: %s", req, update.Message)
		default:
			// Running
			if req == w.currentCheck {
				icon = w.spinner.View()
				style = progressItemCurrentStyle
				label = fmt.Sprintf("%s: %s", req, update.Message)
			} else {
				icon = labelStyle.Render("·")
				style = progressItemStyle
				label = string(req)
			}
		}

		b.WriteString(fmt.Sprintf("%s %s\n", icon, style.Render(label)))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Please wait... • ctrl+c to cancel"))

	return boxStyle.Render(b.String())
}

func (w *Wizard) viewResults() string {
	var b strings.Builder

	session := w.runner.Session()
	outcome := session.OverallResult()

	// Title based on outcome
	switch outcome {
	case preflight.OutcomeSuccess:
		b.WriteString(successStyle.Render("✓ All Checks Passed!"))
		b.WriteString("\n")
		b.WriteString(subtitleStyle.Render("System meets all requirements for Gohan installation"))
	case preflight.OutcomeWarnings:
		b.WriteString(warningStyle.Render("⚠ Passed with Warnings"))
		b.WriteString("\n")
		b.WriteString(subtitleStyle.Render("Some optional requirements not met, but installation can proceed"))
	case preflight.OutcomeBlocked:
		b.WriteString(errorStyle.Render("✗ Validation Failed"))
		b.WriteString("\n")
		b.WriteString(subtitleStyle.Render("Critical requirements not met - cannot proceed with installation"))
	default:
		b.WriteString(titleStyle.Render("Validation Complete"))
	}

	b.WriteString("\n\n")

	// Show all results
	results := session.Results()
	for _, result := range results {
		line := result.FormatMessage()

		switch result.Status() {
		case preflight.StatusPass:
			b.WriteString(progressItemDoneStyle.Render(line))
		case preflight.StatusFail:
			b.WriteString(progressItemFailedStyle.Render(line))
		case preflight.StatusWarning:
			b.WriteString(warningStyle.Render(line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(labelStyle.Render(fmt.Sprintf("Duration: %v", session.Duration())))
	b.WriteString("\n\n")

	// Show appropriate action
	if session.CanProceed() {
		b.WriteString(helpStyle.Render("Press Enter to continue with installation • q to quit"))
	} else {
		b.WriteString(helpStyle.Render("Press Enter to see guidance • q to quit"))
	}

	return boxStyle.Render(b.String())
}

func (w *Wizard) viewGuidance() string {
	var b strings.Builder

	session := w.runner.Session()

	b.WriteString(errorStyle.Render("Installation Blocked"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Please resolve the following issues before installing:"))
	b.WriteString("\n\n")

	blockers := session.BlockingResults()
	for i, result := range blockers {
		if i > 0 {
			b.WriteString("\n\n")
		}

		b.WriteString(errorStyle.Render(fmt.Sprintf("Issue %d: %s", i+1, result.RequirementName())))
		b.WriteString("\n")

		guidance := result.Guidance()
		b.WriteString(labelStyle.Render(guidance.Message()))
		b.WriteString("\n\n")

		steps := guidance.ActionableSteps()
		if len(steps) > 0 {
			b.WriteString(labelStyle.Render("Steps to resolve:"))
			b.WriteString("\n")
			for j, step := range steps {
				b.WriteString(fmt.Sprintf("  %d. %s\n", j+1, step))
			}
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press q to quit and resolve issues"))

	return guidanceBoxStyle.Render(b.String())
}

func (w *Wizard) startValidation() tea.Msg {
	go func() {
		// Run validation in background
		if err := w.runner.Run(w.ctx); err != nil {
			// Error already handled in runner
		}
	}()

	return validationStartedMsg{}
}

func (w *Wizard) waitForProgress() tea.Cmd {
	return func() tea.Msg {
		select {
		case update, ok := <-w.runner.Progress():
			if !ok {
				// Channel closed, validation complete
				return validationCompleteMsg{}
			}
			return progressUpdateMsg(update)
		case <-w.ctx.Done():
			return errMsg(w.ctx.Err())
		}
	}
}

// Messages
type validationStartedMsg struct{}
type validationCompleteMsg struct{}
type progressUpdateMsg ProgressUpdate
type errMsg error
