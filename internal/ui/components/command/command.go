package command

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textInput textinput.Model
	submitted bool
}

var (
	commandStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)
)

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter command (topics, schemas, connectors, quit)..."
	ti.CharLimit = 100
	ti.Width = 50
	ti.Prompt = ":"
	ti.Focus() // Important: Focus by default

	return Model{
		textInput: ti,
	}
}

func (m Model) Focus() tea.Cmd {
	m.textInput.Focus()
	return textinput.Blink
}

func (m *Model) Blur() {
	m.textInput.Blur()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.submitted = true
			return m, nil
		case tea.KeyEsc:
			m.Reset()
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return commandStyle.Render(m.textInput.View())
}

func (m Model) Value() string {
	return m.textInput.Value()
}

func (m Model) Submitted() bool {
	return m.submitted
}

func (m *Model) Reset() {
	m.textInput.SetValue("")
	m.submitted = false
	m.textInput.Blur()
}

func (m *Model) SetValue(value string) {
	m.textInput.SetValue(value)
}
