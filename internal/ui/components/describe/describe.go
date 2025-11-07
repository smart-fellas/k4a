package describe

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents a dedicated describe view for resource details
type Model struct {
	viewport     viewport.Model
	content      string
	resourceName string
	resourceType string
	width        int
	height       int
	ready        bool
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	viewportStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("62"))
)

// New creates a new describe view
func New() Model {
	return Model{
		viewport: viewport.New(80, 20),
	}
}

// SetContent sets the YAML content to display
func (m *Model) SetContent(content string) {
	m.content = content
	if m.ready {
		m.viewport.SetContent(content)
	}
}

// SetResource sets the resource name and type for display
func (m *Model) SetResource(name, resourceType string) {
	m.resourceName = name
	m.resourceType = resourceType
}

// SetSize updates the viewport dimensions
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	headerHeight := 3 // Title + help text + border
	footerHeight := 1

	m.viewport.Width = width - 2 // Account for borders
	m.viewport.Height = height - headerHeight - footerHeight

	if m.content != "" {
		m.viewport.SetContent(m.content)
	}

	m.ready = true
}

// Update handles viewport updates
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the describe view
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Title bar
	title := titleStyle.Width(m.width).Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.resourceType+" ",
			lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Render(m.resourceName),
		),
	)

	// Help text
	help := headerStyle.Width(m.width).Render(
		"↑/↓ navigate • g/G top/bottom • ESC back",
	)

	// Viewport content
	viewportContent := viewportStyle.
		Width(m.width - 2).
		Height(m.viewport.Height).
		Render(m.viewport.View())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		help,
		viewportContent,
	)
}

// KeyBindings for describe view
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Top      key.Binding
	Bottom   key.Binding
	Back     key.Binding
}

// DefaultKeyMap returns the default key bindings for describe view
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b", "ctrl+u"),
			key.WithHelp("pgup/b", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "f", "ctrl+d"),
			key.WithHelp("pgdn/f", "page down"),
		),
		Top: key.NewBinding(
			key.WithKeys("g", "home"),
			key.WithHelp("g", "top"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("G", "end"),
			key.WithHelp("G", "bottom"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "q"),
			key.WithHelp("esc", "back"),
		),
	}
}
