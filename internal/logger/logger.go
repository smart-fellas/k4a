package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logFile      *os.File
	logMutex     sync.Mutex
	debugEnabled bool
)

// Init initializes the debug logger
func Init(debug bool) error {
	debugEnabled = debug

	if !debug {
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	logDir := filepath.Join(homeDir, ".local", "k4a")
	if mkdirErr := os.MkdirAll(logDir, 0o755); mkdirErr != nil {
		return fmt.Errorf("failed to create log directory: %w", mkdirErr)
	}

	logPath := filepath.Join(logDir, "debug.log")

	// Open log file in append mode, create if doesn't exist
	logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	Debugf("=== K4A Debug Session Started ===")
	Debugf("Log file: %s", logPath)

	return nil
}

// Close closes the log file
func Close() {
	if logFile != nil {
		Debugf("=== K4A Debug Session Ended ===")
		_ = logFile.Close() // Errors during close are non-fatal
	}
}

// Debugf logs a debug message with printf-style formatting
func Debugf(format string, args ...any) {
	if !debugEnabled || logFile == nil {
		return
	}

	logMutex.Lock()
	defer logMutex.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, message)

	if _, err := logFile.WriteString(logLine); err != nil {
		// Can't do much if logging fails, but at least don't crash
		return
	}
	// Sync errors are also non-fatal for logging
	//nolint:errcheck // Sync errors are explicitly non-fatal for logging
	logFile.Sync()
}

// Debug is a convenience wrapper for Debugf with no formatting
func Debug(msg string) {
	Debugf("%s", msg)
}

// IsEnabled returns whether debug logging is enabled
func IsEnabled() bool {
	return debugEnabled
}
