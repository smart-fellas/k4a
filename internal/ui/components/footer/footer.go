package footer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	keybindings []Keybinding
	message     string
	width       int
}

type Keybinding struct {
	Key  string
	Desc string
}

var (
	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Bold(true)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220"))
)

func New() Model {
	return Model{
		keybindings: DefaultKeybindings(),
	}
}

func DefaultKeybindings() []Keybinding {
	return []Keybinding{
		{"↑↓", "navigate"},
		{"enter", "select"},
		{"d", "describe"},
		{"r", "refresh"},
		{"/", "filter"},
		{":", "command"},
		{"?", "help"},
		{"q", "quit"},
	}
}

func (m *Model) SetMessage(msg string) {
	m.message = msg
}

func (m *Model) ClearMessage() {
	m.message = ""
}

func (m *Model) SetKeybindings(kb []Keybinding) {
	m.keybindings = kb
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m Model) View() string {
	var parts []string

	for _, kb := range m.keybindings {
		part := fmt.Sprintf("%s %s",
			keyStyle.Render(kb.Key),
			descStyle.Render(kb.Desc),
		)
		parts = append(parts, part)
	}

	keysLine := strings.Join(parts, "  ")

	if m.message != "" {
		keysLine += "  |  " + messageStyle.Render(m.message)
	}

	// Ensure footer spans full width
	if m.width > 0 {
		return footerStyle.Width(m.width).Render(keysLine)
	}

	return footerStyle.Render(keysLine)
}
