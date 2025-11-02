package theme

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rebelopsio/gohan/internal/domain/theme"
)

// KeyMap defines keyboard shortcuts
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
	Help   key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "select theme"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q/esc", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

// Model represents the theme picker state
type Model struct {
	themes         []theme.Theme
	currentTheme   string
	cursor         int
	selected       *theme.Theme
	err            error
	quitting       bool
	showHelp       bool
	keys           KeyMap
	width          int
	height         int
}

// New creates a new theme picker model
func New(themes []theme.Theme, currentTheme string) Model {
	return Model{
		themes:       themes,
		currentTheme: currentTheme,
		cursor:       findCurrentThemeIndex(themes, currentTheme),
		keys:         DefaultKeyMap(),
		width:        80,
		height:       24,
	}
}

func findCurrentThemeIndex(themes []theme.Theme, currentTheme string) int {
	for i, t := range themes {
		if string(t.Name()) == currentTheme {
			return i
		}
	}
	return 0
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.themes)-1 {
				m.cursor++
			}

		case key.Matches(msg, m.keys.Select):
			m.selected = &m.themes[m.cursor]
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		if m.selected != nil {
			return successStyle.Render(fmt.Sprintf("✓ Selected theme: %s", m.selected.DisplayName))
		}
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("🎨 Theme Selector"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Choose a theme for your Hyprland environment"))
	b.WriteString("\n\n")

	// Theme list
	for i, t := range m.themes {
		cursor := " "
		if m.cursor == i {
			cursor = "▶"
		}

		// Active indicator
		activeIndicator := ""
		if string(t.Name()) == m.currentTheme {
			activeIndicator = activeIndicatorStyle.Render(" [ACTIVE]")
		}

		// Theme name
		themeName := fmt.Sprintf("%s %s", getThemeIcon(string(t.Name())), t.DisplayName())

		// Style based on selection
		var line string
		if m.cursor == i {
			line = selectedStyle.Render(cursor + " " + themeName + activeIndicator)
		} else {
			line = unselectedStyle.Render(cursor + " " + themeName + activeIndicator)
		}

		b.WriteString(line)
		b.WriteString("\n")

		// Description (for selected theme)
		if m.cursor == i {
			desc := t.Description()
			if desc == "" {
				desc = getThemeDescription(string(t.Name()))
			}
			b.WriteString(descriptionStyle.Render(desc))
			b.WriteString("\n")
		}
	}

	// Color preview for selected theme
	if m.cursor < len(m.themes) {
		b.WriteString("\n")
		b.WriteString(m.renderColorPreview(m.themes[m.cursor]))
		b.WriteString("\n")
	}

	// Help text
	if m.showHelp {
		b.WriteString("\n")
		b.WriteString(m.renderHelp())
	} else {
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Press ? for help"))
	}

	return b.String()
}

// renderColorPreview shows color swatches for the theme
func (m Model) renderColorPreview(t theme.Theme) string {
	var preview strings.Builder
	preview.WriteString(lipgloss.NewStyle().Bold(true).Render("Color Preview:"))
	preview.WriteString("\n\n")

	// Get color scheme
	cs := t.ColorScheme()

	// Key colors to display
	colors := []struct {
		name  string
		value string
	}{
		{"base", string(cs.Base())},
		{"surface", string(cs.Surface())},
		{"text", string(cs.Text())},
		{"mauve", string(cs.Mauve())},
		{"pink", string(cs.Pink())},
		{"blue", string(cs.Blue())},
		{"green", string(cs.Green())},
		{"yellow", string(cs.Yellow())},
		{"red", string(cs.Red())},
	}

	// Render color swatches in rows
	cols := 3
	for i := 0; i < len(colors); i += cols {
		for j := 0; j < cols && i+j < len(colors); j++ {
			color := colors[i+j]
			colorValue := color.value

			// Add # if not present
			if !strings.HasPrefix(colorValue, "#") {
				colorValue = "#" + colorValue
			}

			// Color name
			name := colorNameStyle.Render(color.name + ":")
			// Color swatch
			swatch := renderColorSwatch(colorValue)
			// Color hex value
			value := colorValueStyle.Render(colorValue)

			preview.WriteString(fmt.Sprintf("  %s %s %s", name, swatch, value))
		}
		preview.WriteString("\n")
	}

	return previewBoxStyle.Render(preview.String())
}

// renderHelp shows keyboard shortcuts
func (m Model) renderHelp() string {
	var help strings.Builder

	help.WriteString(lipgloss.NewStyle().Bold(true).Render("Keyboard Shortcuts:"))
	help.WriteString("\n\n")

	shortcuts := []struct {
		keys string
		desc string
	}{
		{"↑/k", "Move up"},
		{"↓/j", "Move down"},
		{"Enter/Space", "Select theme"},
		{"?", "Toggle help"},
		{"q/Esc", "Quit"},
	}

	for _, s := range shortcuts {
		key := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#89b4fa")).
			Bold(true).
			Width(15).
			Render(s.keys)

		desc := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#cdd6f4")).
			Render(s.desc)

		help.WriteString(fmt.Sprintf("  %s %s\n", key, desc))
	}

	return helpStyle.Render(help.String())
}

// getThemeIcon returns an emoji icon for the theme
func getThemeIcon(themeName string) string {
	icons := map[string]string{
		"mocha":     "🌙",
		"latte":     "☀️",
		"frappe":    "🌆",
		"macchiato": "🌸",
		"gohan":     "🚀",
	}

	if icon, ok := icons[themeName]; ok {
		return icon
	}
	return "🎨"
}

// getThemeDescription returns a description for the theme
func getThemeDescription(themeName string) string {
	descriptions := map[string]string{
		"mocha":     "Warm dark theme with muted colors - perfect for evening work",
		"latte":     "Soft light theme for daytime use - easy on the eyes",
		"frappe":    "Cool dark theme with blue tones - refreshing and calm",
		"macchiato": "Purple-tinted dark theme - elegant and sophisticated",
		"gohan":     "Custom branded theme - unique Gohan experience",
	}

	if desc, ok := descriptions[themeName]; ok {
		return desc
	}
	return "A beautiful theme for your Hyprland environment"
}

// SelectedTheme returns the theme selected by the user
func (m Model) SelectedTheme() *theme.Theme {
	return m.selected
}
