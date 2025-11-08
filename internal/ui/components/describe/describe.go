package describe

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smart-fellas/k4a/internal/logger"
)

// Model represents a dedicated describe view for resource details
type Model struct {
	viewport     viewport.Model
	content      string
	resourceName string
	resourceType string
	width        int
	height       int
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
	logger.Debugf("describe.SetContent called with %d bytes of content", len(content))
	logger.Debugf("Content preview (first 200 chars): %.200s", content)
	m.content = content
	m.viewport.SetContent(content)
	logger.Debugf("Viewport content set, viewport height=%d, width=%d", m.viewport.Height, m.viewport.Width)
}

// SetResource sets the resource name and type for display
func (m *Model) SetResource(name, resourceType string) {
	logger.Debugf("describe.SetResource called: name=%s, type=%s", name, resourceType)
	m.resourceName = name
	m.resourceType = resourceType
}

// SetSize updates the viewport dimensions
func (m *Model) SetSize(width, height int) {
	logger.Debugf("describe.SetSize called: width=%d, height=%d", width, height)
	m.width = width
	m.height = height

	headerHeight := 3 // Title + help text + border
	footerHeight := 1

	m.viewport.Width = width - 2 // Account for borders
	m.viewport.Height = height - headerHeight - footerHeight

	logger.Debugf("Viewport dimensions set: width=%d, height=%d", m.viewport.Width, m.viewport.Height)

	// Always update content when resizing
	if m.content != "" {
		logger.Debugf("Re-setting viewport content (%d bytes)", len(m.content))
		m.viewport.SetContent(m.content)
	} else {
		logger.Debugf("WARNING: No content to set in viewport")
	}
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
	logger.Debugf("describe.View called: width=%d, height=%d, content_len=%d, resource=%s %s",
		m.width, m.height, len(m.content), m.resourceType, m.resourceName)

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
	viewportView := m.viewport.View()
	logger.Debugf("Viewport.View() returned %d bytes", len(viewportView))

	viewportContent := viewportStyle.
		Width(m.width - 2).
		Height(m.viewport.Height).
		Render(viewportView)

	result := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		help,
		viewportContent,
	)

	logger.Debugf("describe.View returning %d bytes", len(result))
	return result
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
