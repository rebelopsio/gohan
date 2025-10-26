package preflight

import (
	"fmt"
	"strings"
)

// UserGuidance provides actionable information for validation failures
type UserGuidance struct {
	message          string
	reason           string
	actionableSteps  []string
	documentationURL string
}

// NewUserGuidance creates new user guidance
func NewUserGuidance(
	message string,
	reason string,
	actionableSteps []string,
	documentationURL string,
) UserGuidance {
	// Ensure non-nil slice
	if actionableSteps == nil {
		actionableSteps = []string{}
	}

	return UserGuidance{
		message:          message,
		reason:           reason,
		actionableSteps:  actionableSteps,
		documentationURL: documentationURL,
	}
}

// Message returns the main guidance message
func (g UserGuidance) Message() string {
	return g.message
}

// Reason returns why the validation failed
func (g UserGuidance) Reason() string {
	return g.reason
}

// ActionableSteps returns steps to resolve the issue
func (g UserGuidance) ActionableSteps() []string {
	if g.actionableSteps == nil {
		return nil
	}
	result := make([]string, len(g.actionableSteps))
	copy(result, g.actionableSteps)
	return result
}

// DocumentationURL returns link to relevant documentation
func (g UserGuidance) DocumentationURL() string {
	return g.documentationURL
}

// HasSteps returns true if actionable steps are provided
func (g UserGuidance) HasSteps() bool {
	return len(g.actionableSteps) > 0
}

// Format returns a formatted guidance message
func (g UserGuidance) Format() string {
	var builder strings.Builder

	builder.WriteString(g.message)
	builder.WriteString("\n")

	if g.reason != "" {
		builder.WriteString("\nReason: ")
		builder.WriteString(g.reason)
		builder.WriteString("\n")
	}

	if g.HasSteps() {
		builder.WriteString("\nHow to fix:\n")
		for i, step := range g.actionableSteps {
			builder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, step))
		}
	}

	if g.documentationURL != "" {
		builder.WriteString("\nLearn more: ")
		builder.WriteString(g.documentationURL)
		builder.WriteString("\n")
	}

	return builder.String()
}
