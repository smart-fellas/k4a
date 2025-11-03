package dialog

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport viewport.Model
	content  string
	width    int
	height   int
	title    string
}

var (
	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			MarginBottom(1)
)

func New() Model {
	vp := viewport.New(80, 20)
	return Model{
		viewport: vp,
		title:    "Resource Details (ESC to close)",
	}
}

func (m *Model) SetContent(content string) {
	m.content = content
	m.viewport.SetContent(content)
}

func (m *Model) SetTitle(title string) {
	m.title = title
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = m.width - 6
		m.viewport.Height = m.height - 6
		m.viewport.SetContent(m.content)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	titleBar := titleStyle.Render(m.title)
	content := lipgloss.JoinVertical(lipgloss.Left, titleBar, m.viewport.View())

	dialog := dialogBoxStyle.
		Width(m.width - 4).
		Height(m.height - 4).
		Render(content)

	// Center the dialog
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}
