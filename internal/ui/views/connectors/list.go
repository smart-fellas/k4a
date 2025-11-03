package connectors

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smart-fellas/k4a/internal/kafkactl"
	"github.com/smart-fellas/k4a/internal/ui/components/dialog"
	"github.com/smart-fellas/k4a/internal/ui/keys"
	"github.com/smart-fellas/k4a/internal/ui/styles"
)

type Model struct {
	client     *kafkactl.Client
	table      table.Model
	connectors []map[string]any
	keys       keys.KeyMap
	width      int
	height     int
	loading    bool
	err        error

	// Detail view
	showDetail   bool
	detailDialog dialog.Model
}

func New(client *kafkactl.Client) Model {
	columns := []table.Column{
		{Title: "Name", Width: 40},
		{Title: "Class", Width: 40},
		{Title: "Type", Width: 10},
		{Title: "State", Width: 10},
		{Title: "Tasks", Width: 10},
		{Title: "Connect Cluster", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return Model{
		client:       client,
		table:        t,
		keys:         keys.DefaultKeyMap(),
		detailDialog: dialog.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadConnectors
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle detail view
	if m.showDetail {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, m.keys.Back) || key.Matches(msg, m.keys.Quit) {
				m.showDetail = false
				return m, nil
			}
		}

		newDialog, cmd := m.detailDialog.Update(msg)
		m.detailDialog = newDialog
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Describe):
			if len(m.connectors) > 0 {
				m.showDetail = true
				return m, m.loadConnectorDetail
			}

		case key.Matches(msg, m.keys.Refresh):
			return m, m.loadConnectors

		case msg.String() == "p":
			// Pause connector
			return m, m.pauseConnector

		case msg.String() == "r":
			// Resume connector
			return m, m.resumeConnector

		case msg.String() == "R":
			// Restart connector
			return m, m.restartConnector
		}

	case connectorsLoadedMsg:
		m.connectors = msg.connectors
		m.loading = false
		m.updateTable()

	case connectorDetailMsg:
		m.detailDialog.SetContent(msg.yaml)
		m.showDetail = true

	case connectorActionMsg:
		// Refresh after action
		return m, m.loadConnectors
	}

	newTable, cmd := m.table.Update(msg)
	m.table = newTable
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.showDetail {
		return m.detailDialog.View()
	}

	if m.loading {
		return "Loading connectors..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	return m.table.View()
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.table.SetHeight(height - 2)
}

func (m *Model) updateTable() {
	rows := []table.Row{}

	for _, connector := range m.connectors {
		metadata, ok := connector["metadata"].(map[string]any)
		if !ok {
			continue
		}

		spec, ok := connector["spec"].(map[string]any)
		if !ok {
			continue
		}

		name, ok := metadata["name"].(string)
		if !ok {
			continue
		}

		connectorClass := "-"
		connectorType := "source"
		state := "RUNNING"
		tasks := "1"
		cluster := "-"

		if config, configOk := spec["config"].(map[string]any); configOk {
			if class, classOk := config["connector.class"].(string); classOk {
				connectorClass = class
				if strings.Contains(strings.ToLower(class), "sink") {
					connectorType = "sink"
				}
			}
			if t, tasksOk := config["tasks.max"]; tasksOk {
				tasks = fmt.Sprintf("%v", t)
			}
		}

		if c, clusterOk := spec["connectCluster"].(string); clusterOk {
			cluster = c
		}

		displayName := styles.StatusDot(state) + " " + name

		rows = append(rows, table.Row{
			displayName,
			connectorClass,
			connectorType,
			state,
			tasks,
			cluster,
		})
	}

	m.table.SetRows(rows)
}

type connectorsLoadedMsg struct {
	connectors []map[string]any
	err        error
}

type connectorDetailMsg struct {
	yaml string
}

type connectorActionMsg struct {
	action string
	result string
}

func (m *Model) loadConnectors() tea.Msg {
	connectors, err := m.client.GetConnectors()
	return connectorsLoadedMsg{connectors: connectors, err: err}
}

func (m *Model) loadConnectorDetail() tea.Msg {
	selectedRow := m.table.SelectedRow()
	if len(selectedRow) == 0 {
		return nil
	}

	connectorName := selectedRow[0]
	yaml, err := m.client.GetResourceYAML("connector", connectorName)
	if err != nil {
		return connectorDetailMsg{yaml: fmt.Sprintf("Error loading connector details: %v", err)}
	}

	return connectorDetailMsg{yaml: yaml}
}

func (m *Model) pauseConnector() tea.Msg {
	selectedRow := m.table.SelectedRow()
	if len(selectedRow) == 0 {
		return nil
	}

	connectorName := selectedRow[0]
	_, err := m.client.ExecuteCommand("connector", "pause", connectorName)
	if err != nil {
		return connectorActionMsg{action: "pause", result: fmt.Sprintf("Error: %v", err)}
	}

	return connectorActionMsg{action: "pause", result: "success"}
}

func (m *Model) resumeConnector() tea.Msg {
	selectedRow := m.table.SelectedRow()
	if len(selectedRow) == 0 {
		return nil
	}

	connectorName := selectedRow[0]
	_, err := m.client.ExecuteCommand("connector", "resume", connectorName)
	if err != nil {
		return connectorActionMsg{action: "resume", result: fmt.Sprintf("Error: %v", err)}
	}

	return connectorActionMsg{action: "resume", result: "success"}
}

func (m *Model) restartConnector() tea.Msg {
	selectedRow := m.table.SelectedRow()
	if len(selectedRow) == 0 {
		return nil
	}

	connectorName := selectedRow[0]
	_, err := m.client.ExecuteCommand("connector", "restart", connectorName)
	if err != nil {
		return connectorActionMsg{action: "restart", result: fmt.Sprintf("Error: %v", err)}
	}

	return connectorActionMsg{action: "restart", result: "success"}
}
