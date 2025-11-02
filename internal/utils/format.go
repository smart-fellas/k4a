package utils

import (
	"fmt"
	"strings"
	"time"
)

// FormatBytes formats bytes to human readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDuration formats milliseconds to human readable duration
func FormatDuration(ms int64) string {
	duration := time.Duration(ms) * time.Millisecond

	if duration < time.Second {
		return fmt.Sprintf("%dms", ms)
	} else if duration < time.Minute {
		return fmt.Sprintf("%.1fs", duration.Seconds())
	} else if duration < time.Hour {
		return fmt.Sprintf("%.1fm", duration.Minutes())
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%.1fh", duration.Hours())
	} else {
		days := duration.Hours() / 24
		return fmt.Sprintf("%.1fd", days)
	}
}

// TruncateString truncates string to specified length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	if maxLen <= 3 {
		return s[:maxLen]
	}

	return s[:maxLen-3] + "..."
}

// PadRight pads string with spaces to the right
func PadRight(s string, length int) string {
	if len(s) >= length {
		return s
	}

	return s + strings.Repeat(" ", length-len(s))
}
