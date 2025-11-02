package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	Primary   = lipgloss.Color("229")
	Secondary = lipgloss.Color("86")
	Success   = lipgloss.Color("42")
	Warning   = lipgloss.Color("214")
	Error     = lipgloss.Color("196")
	Muted     = lipgloss.Color("241")

	// Table Styles
	TableHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(Muted)

	TableSelected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57"))

	// Status Styles
	StatusRunning = lipgloss.NewStyle().
			Foreground(Success)

	StatusPaused = lipgloss.NewStyle().
			Foreground(Warning)

	StatusFailed = lipgloss.NewStyle().
			Foreground(Error)

	// General Styles
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary)

	Subtitle = lipgloss.NewStyle().
			Foreground(Secondary)

	MutedText = lipgloss.NewStyle().
			Foreground(Muted)
)

func StatusDot(status string) string {
	dot := "●"
	switch status {
	case "RUNNING", "ACTIVE", "SUCCESS":
		return StatusRunning.Render(dot)
	case "PAUSED", "PENDING", "WARNING":
		return StatusPaused.Render(dot)
	case "FAILED", "ERROR", "DELETED":
		return StatusFailed.Render(dot)
	default:
		return MutedText.Render(dot)
	}
}
