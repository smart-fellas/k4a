package topics

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smart-fellas/k4a/internal/cache"
	"github.com/smart-fellas/k4a/internal/kafkactl"
	"github.com/smart-fellas/k4a/internal/logger"
	"github.com/smart-fellas/k4a/internal/ui/components/describe"
	"github.com/smart-fellas/k4a/internal/ui/keys"
	"gopkg.in/yaml.v3"
)

type Model struct {
	client  *kafkactl.Client
	table   table.Model
	topics  []map[string]any
	keys    keys.KeyMap
	width   int
	height  int
	loading bool
	err     error

	// Describe view
	showDescribe bool
	describeView describe.Model

	// Consumer groups view
	showConsumers  bool
	consumersTable table.Model

	// Auto-refresh
	refreshInterval time.Duration
	lastRefresh     time.Time
}

func New(client *kafkactl.Client) Model {
	columns := []table.Column{
		{Title: "Name", Width: 40},
		{Title: "Partitions", Width: 12},
		{Title: "Replication", Width: 12},
		{Title: "Retention", Width: 15},
		{Title: "Description", Width: 30},
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
		m.loadTopics(false),
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

	// Handle consumers view
	if m.showConsumers {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, m.keys.Back) || key.Matches(msg, m.keys.Quit) {
				m.showConsumers = false
				return m, nil
			}
		}

		newTable, cmd := m.consumersTable.Update(msg)
		m.consumersTable = newTable
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			// Show consumer groups for selected topic
			if len(m.topics) > 0 {
				m.showConsumers = true
				return m, m.loadConsumerGroups(false)
			}

		case key.Matches(msg, m.keys.Describe):
			// Show YAML detail
			logger.Debugf("topics: Describe key pressed, topics count=%d", len(m.topics))
			if len(m.topics) > 0 {
				logger.Debugf("topics: Calling loadTopicDetail")
				return m, m.loadTopicDetail()
			}
			logger.Debugf("topics: No topics available to describe")

		case key.Matches(msg, m.keys.Refresh):
			// Check if it's Shift+R (force refresh)
			forceRefresh := msg.String() == "R"
			return m, m.loadTopics(forceRefresh)
		}

	case tickRefreshMsg:
		// Auto-refresh timer tick
		return m, tea.Batch(
			m.loadTopics(false),
			m.tickRefresh(),
		)

	case topicsLoadedMsg:
		m.topics = msg.topics
		m.loading = false
		m.lastRefresh = time.Now()
		m.updateTable()

	case topicDetailMsg:
		logger.Debugf("topics: Received topicDetailMsg for '%s' with %d bytes of YAML", msg.name, len(msg.yaml))
		m.describeView.SetContent(msg.yaml)
		m.describeView.SetResource(msg.name, "Topic")
		m.describeView.SetSize(m.width, m.height)
		m.showDescribe = true
		logger.Debugf("topics: showDescribe set to true, width=%d, height=%d", m.width, m.height)

	case consumerGroupsMsg:
		m.updateConsumersTable(msg.groups)
		m.showConsumers = true
	}

	newTable, cmd := m.table.Update(msg)
	m.table = newTable
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.showDescribe {
		logger.Debugf("topics.View: Rendering describe view")
		return m.describeView.View()
	}

	if m.showConsumers {
		logger.Debugf("topics.View: Rendering consumers table")
		return m.consumersTable.View()
	}

	if m.loading {
		return "Loading topics..."
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

	for _, topic := range m.topics {
		metadata, ok := topic["metadata"].(map[string]any)
		if !ok {
			continue
		}

		spec, ok := topic["spec"].(map[string]any)
		if !ok {
			continue
		}

		name, ok := metadata["name"].(string)
		if !ok {
			continue
		}

		partitions := fmt.Sprintf("%v", spec["partitions"])
		replication := fmt.Sprintf("%v", spec["replicationFactor"])

		retention := "-"
		if configs, configsOk := spec["configs"].(map[string]any); configsOk {
			if ret, retOk := configs["retention.ms"]; retOk {
				retention = fmt.Sprintf("%v", ret)
			}
		}

		description := "-"
		if desc, descOk := spec["description"].(string); descOk {
			description = desc
		}

		rows = append(rows, table.Row{
			name,
			partitions,
			replication,
			retention,
			description,
		})
	}

	m.table.SetRows(rows)
}

func (m *Model) updateConsumersTable(groups []map[string]any) {
	columns := []table.Column{
		{Title: "Group ID", Width: 30},
		{Title: "State", Width: 15},
		{Title: "Members", Width: 10},
		{Title: "Lag", Width: 15},
	}

	m.consumersTable = table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(m.height-2),
	)

	// Add consumer group rows
	rows := []table.Row{}
	for _, group := range groups {
		// Parse group data and add rows
		// This is a placeholder - actual implementation depends on kafkactl output structure
		groupId := "-"
		if id, ok := group["metadata"].(map[string]any)["name"].(string); ok {
			groupId = id
		}
		rows = append(rows, table.Row{groupId, "Stable", "3", "0"})
	}

	m.consumersTable.SetRows(rows)
}

// Command messages.
type topicsLoadedMsg struct {
	topics []map[string]any
	err    error
}

type topicDetailMsg struct {
	name string
	yaml string
}

type consumerGroupsMsg struct {
	groups []map[string]any
}

type tickRefreshMsg time.Time

func (m *Model) tickRefresh() tea.Cmd {
	return tea.Tick(m.refreshInterval, func(t time.Time) tea.Msg {
		return tickRefreshMsg(t)
	})
}

func (m *Model) loadTopics(forceRefresh bool) tea.Cmd {
	return func() tea.Msg {
		topics, err := m.client.GetTopics(forceRefresh)
		return topicsLoadedMsg{topics: topics, err: err}
	}
}

func (m *Model) loadTopicDetail() tea.Cmd {
	return func() tea.Msg {
		logger.Debugf("topics.loadTopicDetail: Starting, topics count=%d", len(m.topics))

		if len(m.topics) == 0 {
			logger.Debugf("topics.loadTopicDetail: No topics available")
			return nil
		}

		selectedRow := m.table.SelectedRow()
		if len(selectedRow) == 0 {
			logger.Debugf("topics.loadTopicDetail: No row selected")
			return nil
		}

		topicName := selectedRow[0]
		logger.Debugf("topics.loadTopicDetail: Selected topic name: '%s'", topicName)

		// Find the topic in the cached list
		var topicData map[string]any
		for _, topic := range m.topics {
			if metadata, ok := topic["metadata"].(map[string]any); ok {
				if name, nameOk := metadata["name"].(string); nameOk && name == topicName {
					topicData = topic
					logger.Debugf("topics.loadTopicDetail: Found topic '%s' in cache", topicName)
					break
				}
			}
		}

		if topicData == nil {
			logger.Debugf("topics.loadTopicDetail: Topic '%s' NOT found in cache", topicName)
			return topicDetailMsg{name: topicName, yaml: fmt.Sprintf("Topic '%s' not found in cache", topicName)}
		}

		// Convert the topic data back to YAML
		yamlBytes, err := yaml.Marshal(topicData)
		if err != nil {
			logger.Debugf("topics.loadTopicDetail: Error marshaling YAML: %v", err)
			return topicDetailMsg{name: topicName, yaml: fmt.Sprintf("Error serializing topic details: %v", err)}
		}

		logger.Debugf("topics.loadTopicDetail: Successfully marshaled %d bytes of YAML", len(yamlBytes))
		return topicDetailMsg{name: topicName, yaml: string(yamlBytes)}
	}
}

func (m *Model) loadConsumerGroups(forceRefresh bool) tea.Cmd {
	return func() tea.Msg {
		selectedRow := m.table.SelectedRow()
		if len(selectedRow) == 0 {
			return nil
		}

		topicName := selectedRow[0]
		groups, err := m.client.GetConsumerGroups(topicName, forceRefresh)
		if err != nil {
			return nil
		}

		return consumerGroupsMsg{groups: groups}
	}
}
