package help

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	sections []Section
}

type Section struct {
	Title    string
	Commands []Command
}

type Command struct {
	Key  string
	Desc string
}

var (
	helpStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			MarginBottom(1).
			Align(lipgloss.Center)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			MarginTop(1).
			MarginBottom(1)

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Width(15)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))
)

func New() Model {
	return Model{
		sections: defaultSections(),
	}
}

func defaultSections() []Section {
	return []Section{
		{
			Title: "Navigation",
			Commands: []Command{
				{"↑/k", "Move up"},
				{"↓/j", "Move down"},
				{"←/h", "Move left"},
				{"→/l", "Move right"},
				{"g", "Go to top"},
				{"G", "Go to bottom"},
				{"PgUp", "Page up"},
				{"PgDn", "Page down"},
			},
		},
		{
			Title: "Actions",
			Commands: []Command{
				{"enter", "Select/Drill down"},
				{"esc", "Go back"},
				{"d", "Describe resource"},
				{"e", "Edit resource"},
				{"ctrl+d", "Delete resource"},
				{"r", "Refresh view"},
				{"/", "Filter resources"},
				{"ctrl+r", "Force refresh"},
			},
		},
		{
			Title: "View Commands",
			Commands: []Command{
				{":topics", "Switch to topics view"},
				{":schemas", "Switch to schemas view"},
				{":connectors", "Switch to connectors view"},
				{":consumers", "Switch to consumers view"},
				{":acls", "Switch to ACLs view"},
				{":ctx", "Switch context"},
				{":ns", "Switch namespace"},
			},
		},
		{
			Title: "Connector Actions",
			Commands: []Command{
				{"p", "Pause connector"},
				{"r", "Resume connector"},
				{"R", "Restart connector"},
			},
		},
		{
			Title: "General",
			Commands: []Command{
				{":", "Enter command mode"},
				{"?", "Show this help"},
				{"q", "Quit application"},
				{"ctrl+c", "Force quit"},
			},
		},
	}
}

func (m Model) View() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("K4A Help"))
	content.WriteString("\n\n")
	content.WriteString("Press ? or ESC to close\n\n")

	for _, section := range m.sections {
		content.WriteString(sectionStyle.Render(section.Title))
		content.WriteString("\n")

		for _, cmd := range section.Commands {
			content.WriteString(fmt.Sprintf("  %s %s\n",
				keyStyle.Render(cmd.Key),
				descStyle.Render(cmd.Desc),
			))
		}
		content.WriteString("\n")
	}

	return helpStyle.Render(content.String())
}
