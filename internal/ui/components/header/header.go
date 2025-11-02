package header

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	context     string
	namespace   string
	api         string
	currentView string
	width       int
}

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229"))

	asciiStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	viewStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))
)

var asciiArt = []string{
	` ____  __.  _____         `,
	`|    |/ _| /  |  |_____   `,
	`|      <  /   |  |\__  \  `,
	`|    |  \/    ^   // __ \_ `,
	`|____|__ \____   |(____  / `,
	`        \/    |__|     \/  `,
}

func New(context, namespace, api string) Model {
	return Model{
		context:     context,
		namespace:   namespace,
		api:         api,
		currentView: "topics",
	}
}

func (m *Model) SetView(view string) {
	m.currentView = view
}

func (m *Model) SetContext(context string) {
	m.context = context
}

func (m *Model) SetNamespace(namespace string) {
	m.namespace = namespace
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m Model) View() string {
	// Build ASCII art section
	var asciiSection strings.Builder
	for _, line := range asciiArt {
		asciiSection.WriteString(asciiStyle.Render(line))
		asciiSection.WriteString("\n")
	}

	// Build info section
	infoLines := []string{
		fmt.Sprintf("%s %s", labelStyle.Render("Context:  "), infoStyle.Render(m.context)),
		fmt.Sprintf("%s %s", labelStyle.Render("Namespace:"), infoStyle.Render(m.namespace)),
		fmt.Sprintf("%s %s", labelStyle.Render("API:      "), infoStyle.Render(truncateAPI(m.api))),
		fmt.Sprintf("%s %s", labelStyle.Render("View:     "), viewStyle.Render(":"+m.currentView)),
		fmt.Sprintf("%s %s", labelStyle.Render("Time:     "), infoStyle.Render(time.Now().Format("15:04:05"))),
	}

	// Calculate spacing
	padding := 3

	// Combine ASCII art and info side by side
	var result strings.Builder
	for i, asciiLine := range strings.Split(asciiSection.String(), "\n") {
		if i < len(asciiArt) {
			result.WriteString(asciiLine)
			result.WriteString(strings.Repeat(" ", padding))
			if i < len(infoLines) {
				result.WriteString(infoLines[i])
			}
			result.WriteString("\n")
		}
	}

	// Add a separator line
	result.WriteString(strings.Repeat("─", m.width))

	return result.String()
}

func truncateAPI(api string) string {
	// Remove https:// prefix for display
	api = strings.TrimPrefix(api, "https://")
	api = strings.TrimPrefix(api, "http://")

	// Truncate if too long
	if len(api) > 40 {
		return api[:37] + "..."
	}
	return api
}
