package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smart-fellas/k4a/internal/app"
	"github.com/smart-fellas/k4a/internal/config"
	"github.com/smart-fellas/k4a/internal/logger"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.BoolVar(versionFlag, "v", false, "Print version information (shorthand)")
	debugFlag := flag.Bool("debug", false, "Enable debug logging to ~/.local/k4a/debug.log")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("k4a version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built at: %s\n", date)
		os.Exit(0)
	}

	// Initialize logger
	if err := logger.Init(*debugFlag); err != nil {
		fmt.Printf("Warning: failed to initialize logger: %v\n", err)
	}
	defer logger.Close()

	logger.Debugf("Starting k4a version %s (commit: %s, built: %s)", version, commit, date)

	cfg, err := config.Load()
	if err != nil {
		logger.Debugf("Error loading config: %v", err)
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	logger.Debugf("Config loaded successfully, current context: %s", cfg.CurrentContext)

	logger.Debugf("Creating Bubble Tea program")
	p := tea.NewProgram(
		app.New(cfg),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	logger.Debugf("Starting program")
	if _, runErr := p.Run(); runErr != nil {
		logger.Debugf("Error running program: %v", runErr)
		fmt.Printf("Error running program: %v\n", runErr)
		os.Exit(1)
	}

	logger.Debugf("Program exited normally")
}
