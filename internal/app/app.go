package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/smart-fellas/k4a/internal/config"
	"github.com/smart-fellas/k4a/internal/kafkactl"
	"github.com/smart-fellas/k4a/internal/ui/components/command"
	"github.com/smart-fellas/k4a/internal/ui/components/footer"
	"github.com/smart-fellas/k4a/internal/ui/components/header"
	"github.com/smart-fellas/k4a/internal/ui/components/help"
	"github.com/smart-fellas/k4a/internal/ui/keys"
	"github.com/smart-fellas/k4a/internal/ui/views/connectors"
	"github.com/smart-fellas/k4a/internal/ui/views/schemas"
	"github.com/smart-fellas/k4a/internal/ui/views/topics"
)

type ViewType string

const (
	TopicsView     ViewType = "topics"
	SchemasView    ViewType = "schemas"
	ConnectorsView ViewType = "connectors"
	ConsumersView  ViewType = "consumers"
	ACLsView       ViewType = "acls"
)

type Model struct {
	config      *config.Config
	client      *kafkactl.Client
	currentView ViewType
	width       int
	height      int

	// Components
	header  header.Model
	footer  footer.Model
	command command.Model
	help    help.Model

	// Views
	topicsView     topics.Model
	schemasView    schemas.Model
	connectorsView connectors.Model

	// State
	commandMode bool
	helpVisible bool
	keys        keys.KeyMap
}

func New(cfg *config.Config) Model {
	client := kafkactl.NewClient(cfg)

	// Get current context details
	ctx, err := cfg.GetCurrentContext()
	contextName := cfg.CurrentContext
	namespace := ""
	api := ""

	if err == nil && ctx != nil {
		namespace = ctx.Context.Namespace
		api = ctx.Context.API
	}

	return Model{
		config:         cfg,
		client:         client,
		currentView:    TopicsView,
		header:         header.New(contextName, namespace, api),
		footer:         footer.New(),
		command:        command.New(),
		help:           help.New(),
		topicsView:     topics.New(client),
		schemasView:    schemas.New(client),
		connectorsView: connectors.New(client),
		keys:           keys.DefaultKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.topicsView.Init(),
		tea.EnterAltScreen,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(m.width)
		m.updateLayout()

	case tea.KeyMsg:
		// Handle command mode
		if m.commandMode {
			return m.handleCommandMode(msg)
		}

		// Handle help toggle
		if key.Matches(msg, m.keys.Help) {
			m.helpVisible = !m.helpVisible
			return m, nil
		}

		// Close help if visible
		if m.helpVisible {
			if msg.Type == tea.KeyEsc || key.Matches(msg, m.keys.Help) {
				m.helpVisible = false
				return m, nil
			}
		}

		// Handle quit
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

		// Handle colon command - check for ":" specifically
		if msg.String() == ":" {
			m.commandMode = true
			m.command = command.New() // Reset command
			cmd := m.command.Focus()
			return m, cmd
		}

		// Handle direct view switching commands (when typed quickly)
		msgStr := msg.String()
		if strings.HasPrefix(msgStr, ":") {
			switch strings.TrimPrefix(msgStr, ":") {
			case "topics", "topic":
				m.switchView(TopicsView)
				return m, nil
			case "schemas", "schema":
				m.switchView(SchemasView)
				return m, nil
			case "connectors", "connector":
				m.switchView(ConnectorsView)
				return m, nil
			}
		}
	}

	// Don't update views if help is visible
	if m.helpVisible {
		return m, nil
	}

	// Update current view
	switch m.currentView {
	case TopicsView:
		newView, cmd := m.topicsView.Update(msg)
		if tv, ok := newView.(topics.Model); ok {
			m.topicsView = tv
		}
		cmds = append(cmds, cmd)

	case SchemasView:
		newView, cmd := m.schemasView.Update(msg)
		if sv, ok := newView.(schemas.Model); ok {
			m.schemasView = sv
		}
		cmds = append(cmds, cmd)

	case ConnectorsView:
		newView, cmd := m.connectorsView.Update(msg)
		if cv, ok := newView.(connectors.Model); ok {
			m.connectorsView = cv
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.helpVisible {
		return m.help.View()
	}

	var content string

	// Show command input if in command mode
	if m.commandMode {
		content = m.command.View()
	} else {
		// Show current view
		switch m.currentView {
		case TopicsView:
			content = m.topicsView.View()
		case SchemasView:
			content = m.schemasView.View()
		case ConnectorsView:
			content = m.connectorsView.View()
		}
	}

	return m.header.View() + "\n" + content + "\n" + m.footer.View()
}

func (m *Model) switchView(view ViewType) {
	m.currentView = view
	m.header.SetView(string(view))

	// Update footer keybindings based on view
	switch view {
	case ConnectorsView:
		m.footer.SetKeybindings([]footer.Keybinding{
			{Key: "↑↓", Desc: "navigate"},
			{Key: "enter", Desc: "select"},
			{Key: "d", Desc: "describe"},
			{Key: "p", Desc: "pause"},
			{Key: "r", Desc: "resume"},
			{Key: "R", Desc: "restart"},
			{Key: ":", Desc: "command"},
			{Key: "?", Desc: "help"},
			{Key: "q", Desc: "quit"},
		})
	default:
		m.footer.SetKeybindings(footer.DefaultKeybindings())
	}
}

func (m *Model) updateLayout() {
	headerHeight := 6 // ASCII art is 5 lines + separator
	footerHeight := 2
	contentHeight := m.height - headerHeight - footerHeight

	if contentHeight < 1 {
		contentHeight = 1
	}

	m.topicsView.SetSize(m.width, contentHeight)
	m.schemasView.SetSize(m.width, contentHeight)
	m.connectorsView.SetSize(m.width, contentHeight)
}

func (m Model) handleCommandMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// ESC exits command mode
	if msg.Type == tea.KeyEsc {
		m.commandMode = false
		m.command.Reset()
		return m, nil
	}

	// Update command input
	newCmd, cmd := m.command.Update(msg)
	m.command = newCmd

	// Check if command was submitted
	if m.command.Submitted() {
		cmdText := strings.TrimSpace(m.command.Value())
		m.commandMode = false
		m.command.Reset()

		// Process command
		switch cmdText {
		case "topics", "topic":
			m.switchView(TopicsView)
		case "schemas", "schema":
			m.switchView(SchemasView)
		case "connectors", "connector":
			m.switchView(ConnectorsView)
		case "consumers", "consumer":
			m.switchView(ConsumersView)
		case "acls", "acl":
			m.switchView(ACLsView)
		case "q", "quit":
			return m, tea.Quit
		}
	}

	return m, cmd
}
