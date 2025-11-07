package schemas

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smart-fellas/k4a/internal/cache"
	"github.com/smart-fellas/k4a/internal/kafkactl"
	"github.com/smart-fellas/k4a/internal/ui/components/dialog"
	"github.com/smart-fellas/k4a/internal/ui/keys"
	"gopkg.in/yaml.v3"
)

type Model struct {
	client  *kafkactl.Client
	table   table.Model
	schemas []map[string]any
	keys    keys.KeyMap
	width   int
	height  int
	loading bool
	err     error

	// Detail view
	showDetail   bool
	detailDialog dialog.Model

	// Auto-refresh
	refreshInterval time.Duration
	lastRefresh     time.Time
}

func New(client *kafkactl.Client) Model {
	columns := []table.Column{
		{Title: "Subject", Width: 50},
		{Title: "Version", Width: 10},
		{Title: "ID", Width: 10},
		{Title: "Type", Width: 15},
		{Title: "Compatibility", Width: 20},
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
		detailDialog:    dialog.New(),
		refreshInterval: cache.DefaultRefreshInterval,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadSchemas(false),
		m.tickRefresh(),
	)
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
			if len(m.schemas) > 0 {
				m.showDetail = true
				return m, m.loadSchemaDetail
			}

		case key.Matches(msg, m.keys.Refresh):
			// Check if it's Shift+R (force refresh)
			forceRefresh := msg.String() == "R"
			return m, m.loadSchemas(forceRefresh)
		}

	case tickRefreshMsg:
		// Auto-refresh timer tick
		return m, tea.Batch(
			m.loadSchemas(false),
			m.tickRefresh(),
		)

	case schemasLoadedMsg:
		m.schemas = msg.schemas
		m.loading = false
		m.lastRefresh = time.Now()
		m.updateTable()

	case schemaDetailMsg:
		m.detailDialog.SetContent(msg.yaml)
		m.showDetail = true
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
		return "Loading schemas..."
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

	for _, schema := range m.schemas {
		metadata, ok := schema["metadata"].(map[string]any)
		if !ok {
			continue
		}

		subject, ok := metadata["name"].(string)
		if !ok {
			continue
		}

		version := "latest"
		id := "-"
		schemaType := "AVRO"
		compatibility := "BACKWARD"

		if spec, specOk := schema["spec"].(map[string]any); specOk {
			if v, versionOk := spec["version"]; versionOk {
				version = fmt.Sprintf("%v", v)
			}
			if i, idOk := spec["id"]; idOk {
				id = fmt.Sprintf("%v", i)
			}
			if t, typeOk := spec["type"]; typeOk {
				schemaType = fmt.Sprintf("%v", t)
			}
			if c, compatOk := spec["compatibility"]; compatOk {
				compatibility = fmt.Sprintf("%v", c)
			}
		}

		rows = append(rows, table.Row{
			subject,
			version,
			id,
			schemaType,
			compatibility,
		})
	}

	m.table.SetRows(rows)
}

type schemasLoadedMsg struct {
	schemas []map[string]any
	err     error
}

type schemaDetailMsg struct {
	yaml string
}

type tickRefreshMsg time.Time

func (m *Model) tickRefresh() tea.Cmd {
	return tea.Tick(m.refreshInterval, func(t time.Time) tea.Msg {
		return tickRefreshMsg(t)
	})
}

func (m *Model) loadSchemas(forceRefresh bool) tea.Cmd {
	return func() tea.Msg {
		schemas, err := m.client.GetSchemas(forceRefresh)
		return schemasLoadedMsg{schemas: schemas, err: err}
	}
}

func (m *Model) loadSchemaDetail() tea.Msg {
	if len(m.schemas) == 0 {
		return nil
	}

	selectedRow := m.table.SelectedRow()
	if len(selectedRow) == 0 {
		return nil
	}

	schemaName := selectedRow[0]

	// Find the schema in the cached list
	var schemaData map[string]any
	for _, schema := range m.schemas {
		if metadata, ok := schema["metadata"].(map[string]any); ok {
			if name, ok := metadata["name"].(string); ok && name == schemaName {
				schemaData = schema
				break
			}
		}
	}

	if schemaData == nil {
		return schemaDetailMsg{yaml: fmt.Sprintf("Schema '%s' not found in cache", schemaName)}
	}

	// Convert the schema data back to YAML
	yamlBytes, err := yaml.Marshal(schemaData)
	if err != nil {
		return schemaDetailMsg{yaml: fmt.Sprintf("Error serializing schema details: %v", err)}
	}

	return schemaDetailMsg{yaml: string(yamlBytes)}
}
