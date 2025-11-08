package connectors

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smart-fellas/k4a/internal/cache"
	"github.com/smart-fellas/k4a/internal/kafkactl"
	"github.com/smart-fellas/k4a/internal/ui/components/describe"
	"github.com/smart-fellas/k4a/internal/ui/keys"
	"github.com/smart-fellas/k4a/internal/ui/styles"
	"gopkg.in/yaml.v3"
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

	// Describe view
	showDescribe bool
	describeView describe.Model

	// Auto-refresh
	refreshInterval time.Duration
	lastRefresh     time.Time
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
		client:          client,
		table:           t,
		keys:            keys.DefaultKeyMap(),
		describeView:    describe.New(),
		refreshInterval: cache.DefaultRefreshInterval,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadConnectors(false),
		m.tickRefresh(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle describe view
	if m.showDescribe {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, m.keys.Back) || key.Matches(msg, m.keys.Quit) {
				m.showDescribe = false
				return m, nil
			}
		}

		newDescribe, cmd := m.describeView.Update(msg)
		m.describeView = newDescribe
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Describe):
			if len(m.connectors) > 0 {
				return m, m.loadConnectorDetail()
			}

		case key.Matches(msg, m.keys.Refresh):
			// Check if it's Shift+R (force refresh)
			forceRefresh := msg.String() == "R"
			return m, m.loadConnectors(forceRefresh)

		case msg.String() == "p":
			// Pause connector
			return m, m.pauseConnector

		case msg.String() == "s":
			// Resume connector (changed from 'r' to 's' for start/resume)
			return m, m.resumeConnector

		case msg.String() == "t":
			// Restart connector (changed from 'R' to 't' for restart)
			return m, m.restartConnector
		}

	case tickRefreshMsg:
		// Auto-refresh timer tick
		return m, tea.Batch(
			m.loadConnectors(false),
			m.tickRefresh(),
		)

	case connectorsLoadedMsg:
		m.connectors = msg.connectors
		m.loading = false
		m.lastRefresh = time.Now()
		m.updateTable()

	case connectorDetailMsg:
		m.describeView.SetContent(msg.yaml)
		m.describeView.SetResource(msg.name, "Connector")
		m.describeView.SetSize(m.width, m.height)
		m.showDescribe = true

	case connectorActionMsg:
		// Refresh after action
		return m, m.loadConnectors(false)
	}

	newTable, cmd := m.table.Update(msg)
	m.table = newTable
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.showDescribe {
		return m.describeView.View()
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
	m.describeView.SetSize(width, height)
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
	name string
	yaml string
}

type connectorActionMsg struct {
	action string
	result string
}

type tickRefreshMsg time.Time

func (m *Model) tickRefresh() tea.Cmd {
	return tea.Tick(m.refreshInterval, func(t time.Time) tea.Msg {
		return tickRefreshMsg(t)
	})
}

func (m *Model) loadConnectors(forceRefresh bool) tea.Cmd {
	return func() tea.Msg {
		connectors, err := m.client.GetConnectors(forceRefresh)
		return connectorsLoadedMsg{connectors: connectors, err: err}
	}
}

func (m *Model) loadConnectorDetail() tea.Cmd {
	return func() tea.Msg {
		if len(m.connectors) == 0 {
			return nil
		}

		selectedRow := m.table.SelectedRow()
		if len(selectedRow) == 0 {
			return nil
		}

		// Remove the status dot from the connector name
		connectorName := strings.TrimPrefix(selectedRow[0], "● ")
		connectorName = strings.TrimPrefix(connectorName, "○ ")
		connectorName = strings.TrimSpace(connectorName)

		// Find the connector in the cached list
		var connectorData map[string]any
		for _, connector := range m.connectors {
			if metadata, ok := connector["metadata"].(map[string]any); ok {
				if name, nameOk := metadata["name"].(string); nameOk && name == connectorName {
					connectorData = connector
					break
				}
			}
		}

		if connectorData == nil {
			return connectorDetailMsg{name: connectorName, yaml: fmt.Sprintf("Connector '%s' not found in cache", connectorName)}
		}

		// Convert the connector data back to YAML
		yamlBytes, err := yaml.Marshal(connectorData)
		if err != nil {
			return connectorDetailMsg{name: connectorName, yaml: fmt.Sprintf("Error serializing connector details: %v", err)}
		}

		return connectorDetailMsg{name: connectorName, yaml: string(yamlBytes)}
	}
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
