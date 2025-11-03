package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Primary color for main UI elements.
	Primary = lipgloss.Color("229")
	// Secondary color for secondary UI elements.
	Secondary = lipgloss.Color("86")
	// Success color for success states.
	Success = lipgloss.Color("42")
	// Warning color for warning states.
	Warning = lipgloss.Color("214")
	// Error color for error states.
	Error = lipgloss.Color("196")
	// Muted color for muted text.
	Muted = lipgloss.Color("241")

	// TableHeader style for table headers.
	TableHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(Muted)

	// TableSelected style for selected table rows.
	TableSelected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57"))

	// StatusRunning style for running status.
	StatusRunning = lipgloss.NewStyle().
			Foreground(Success)

	// StatusPaused style for paused status.
	StatusPaused = lipgloss.NewStyle().
			Foreground(Warning)

	// StatusFailed style for failed status.
	StatusFailed = lipgloss.NewStyle().
			Foreground(Error)

	// Title style for titles.
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary)

	// Subtitle style for subtitles.
	Subtitle = lipgloss.NewStyle().
			Foreground(Secondary)

	// MutedText style for muted text.
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
