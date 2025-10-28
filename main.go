package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Topic represents the relevant fields from the YAML
// Only the fields we want to display are included
// Add more fields as needed

type Topic struct {
	Metadata struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		ReplicationFactor interface{} `yaml:"replicationFactor"`
		Partitions        interface{} `yaml:"partitions"`
	} `yaml:"spec"`
	Status struct {
		Phase   string `yaml:"phase"`
		Message string `yaml:"message"`
	} `yaml:"status"`
}

// parseTopics parses multiple YAML documents from r into a slice of Topic
func parseTopics(r io.Reader) ([]Topic, error) {
	var topics []Topic
	dec := yaml.NewDecoder(r)
	for {
		var t Topic
		err := dec.Decode(&t)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// Only add if Name is not empty (skip empty docs)
		if t.Metadata.Name != "" {
			topics = append(topics, t)
		}
	}
	return topics, nil
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table         table.Model
	topics        []Topic
	viewMode      string // "table" or "detail"
	selectedTopic int    // index in topics
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.viewMode {
		case "table":
			switch msg.String() {
			case "esc":
				if m.table.Focused() {
					m.table.Blur()
				} else {
					m.table.Focus()
				}
			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter", "right":
				m.selectedTopic = m.table.Cursor()
				m.viewMode = "detail"
				return m, nil
			}
		case "detail":
			switch msg.String() {
			case "esc", "q", "left":
				m.viewMode = "table"
				return m, nil
			}
		}
	}
	if m.viewMode == "table" {
		m.table, cmd = m.table.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	switch m.viewMode {
	case "detail":
		if m.selectedTopic >= 0 && m.selectedTopic < len(m.topics) {
			topic := m.topics[m.selectedTopic]
			data, err := yaml.Marshal(topic)
			if err != nil {
				return "Error displaying topic details"
			}
			return baseStyle.Render(string(data)) + "\n\nPress Esc or q to go back."
		}
		return "No topic selected. Press Esc or q to go back."
	default:
		return baseStyle.Render(m.table.View()) + "\n"
	}
}

func main() {
	// Run 'cat output.yaml' and capture its output
	cmd := exec.Command("cat", "output.yaml")
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run 'cat output.yaml': %v\n", err)
		os.Exit(1)
	}

	topics, err := parseTopics(strings.NewReader(string(output)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse YAML: %v\n", err)
		os.Exit(1)
	}

	columns := []table.Column{
		{Title: "Name", Width: 40},
		{Title: "Namespace", Width: 15},
		{Title: "Phase", Width: 10},
		{Title: "Replicas", Width: 8},
		{Title: "Partitions", Width: 10},
	}

	var rows []table.Row
	for _, t := range topics {
		replicas := "-"
		if t.Spec.ReplicationFactor != nil && strings.TrimSpace(fmt.Sprintf("%v", t.Spec.ReplicationFactor)) != "" {
			replicas = fmt.Sprintf("%v", t.Spec.ReplicationFactor)
		}
		partitions := "-"
		if t.Spec.Partitions != nil && strings.TrimSpace(fmt.Sprintf("%v", t.Spec.Partitions)) != "" {
			partitions = fmt.Sprintf("%v", t.Spec.Partitions)
		}
		phase := t.Status.Phase
		if strings.TrimSpace(phase) == "" {
			phase = "-"
		}
		rows = append(rows, table.Row{
			t.Metadata.Name,
			t.Metadata.Namespace,
			phase,
			replicas,
			partitions,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
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

	m := model{
		table:         t,
		topics:        topics,
		viewMode:      "table",
		selectedTopic: 0,
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
