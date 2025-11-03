# GitHub Copilot Instructions for K4A Project

## Project Overview

**K4A (Kafka for All)** is a Terminal User Interface (TUI) application for managing Kafka resources, inspired by k9s. It provides an intuitive, keyboard-driven interface for interacting with Kafka through the kafkactl API.

### Tech Stack
- **Language**: Go 1.25+
- **UI Framework**: Bubble Tea (Charm.sh) - TUI framework
- **UI Components**: Bubbles (Charm.sh) - Reusable TUI components
- **Styling**: Lipgloss (Charm.sh) - Terminal styling
- **External Dependency**: kafkactl CLI tool (Michelin)

### Repository
- **GitHub**: `github.com/smart-fellas/k4a`
- **License**: Apache 2.0
- **Type**: Open Source

## Project Structure

```
k4a/
├── cmd/k4a/              # Application entry point
│   └── main.go           # Main function, CLI flags, version info
├── internal/             # Private application code
│   ├── app/              # Core application logic
│   │   └── app.go        # Bubble Tea model and update logic
│   ├── config/           # Configuration management
│   │   ├── config.go     # Config loading (reads ~/.kafkactl/config.yml)
│   │   └── context/      # Context switching logic
│   ├── kafkactl/         # Kafkactl CLI wrapper
│   │   └── client.go     # Executes kafkactl commands
│   └── ui/               # User interface components
│       ├── components/   # Reusable UI components
│       │   ├── command/  # Command input handler
│       │   ├── dialog/   # Dialog/modal windows
│       │   ├── footer/   # Bottom status bar
│       │   ├── header/   # Top navigation bar
│       │   └── help/     # Help screen
│       ├── keys/         # Keyboard bindings
│       ├── styles/       # Lipgloss styles
│       └── views/        # Main views (screens)
│           ├── connectors/ # Kafka Connect connectors view
│           ├── schemas/    # Schema Registry view
│           └── topics/     # Kafka topics view
├── pkg/                  # Public packages (reusable)
│   └── models/           # Shared data models
└── bin/                  # Compiled binaries (not in git)
```

## Architecture Patterns

### 1. Bubble Tea Pattern (Elm Architecture)
The application follows the Model-View-Update pattern:
- **Model**: Application state (in `internal/app/app.go`)
- **View**: Rendering logic (returns strings)
- **Update**: Handles messages and updates state

```go
type Model struct {
    // State fields
}

func (m Model) Init() tea.Cmd { /* Initialize */ }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { /* Handle events */ }
func (m Model) View() string { /* Render UI */ }
```

### 2. Component-Based UI
Each UI component (header, footer, list, dialog) is self-contained with its own:
- State
- Update logic
- View rendering
- Key bindings

### 3. Command Pattern
User commands (`:topics`, `:schemas`, etc.) are parsed and executed to switch views.

### 4. External Process Execution
The app shells out to `kafkactl` CLI for all Kafka operations. No direct Kafka library usage.

## Code Style & Conventions

### General Go Style
- Follow standard Go conventions (effective Go)
- Use `gofmt` for formatting (enforced in CI)
- Use `golangci-lint` for linting (enforced in CI)
- Keep functions small and focused
- Use early returns to reduce nesting
- Prefer named return values for complex functions

### Naming Conventions
- **Packages**: lowercase, single word (e.g., `config`, `app`, `ui`)
- **Files**: lowercase with underscores (e.g., `main.go`, `list.go`)
- **Types**: PascalCase (e.g., `Model`, `KeyMap`, `Config`)
- **Functions/Methods**: camelCase for private, PascalCase for exported
- **Variables**: camelCase (e.g., `currentView`, `selectedIndex`)
- **Constants**: PascalCase or UPPER_SNAKE_CASE for exported

### Project-Specific Conventions

#### 1. Error Handling
```go
// Always check errors immediately
if err != nil {
    // Handle or wrap errors with context
    return fmt.Errorf("failed to load config: %w", err)
}

// For UI errors, often return as tea.Cmd
return func() tea.Msg {
    return errorMsg{err: err}
}
```

#### 2. Bubble Tea Messages
Define custom message types for events:
```go
type tickMsg time.Time
type resourceLoadedMsg struct {
    resources []Resource
}
type errorMsg struct {
    err error
}
```

#### 3. Key Bindings
Use `key.Binding` from bubbles/key:
```go
type KeyMap struct {
    Up    key.Binding
    Down  key.Binding
    Enter key.Binding
}

var DefaultKeyMap = KeyMap{
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "up"),
    ),
}
```

#### 4. Styling with Lipgloss
Define styles as package variables:
```go
var (
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("170"))
    
    selectedStyle = lipgloss.NewStyle().
        Background(lipgloss.Color("62"))
)
```

#### 5. Component Structure
Each component should have:
```go
type Component struct {
    // State fields
    width  int
    height int
}

func New() Component { /* Constructor */ }
func (c Component) Update(msg tea.Msg) (Component, tea.Cmd) { /* Update */ }
func (c Component) View() string { /* Render */ }
func (c *Component) SetSize(width, height int) { /* Resize */ }
```

#### 6. Config File Location
- Always use `~/.kafkactl/config.yml` for configuration
- Support multiple contexts (dev, prod, etc.)
- Never hardcode credentials

### Comments & Documentation
```go
// Package comment at the top of each package
// Explains the purpose of the package

// Exported functions/types must have doc comments
// Format: starts with the name being documented

// New creates a new instance of Component.
// It initializes default values and returns a ready-to-use component.
func New() Component {
    // Implementation
}

// Complex internal logic should have inline comments
func complexFunction() {
    // Step 1: Validate input
    // ...
    
    // Step 2: Process data
    // ...
}
```

### Testing
```go
// Test files: *_test.go
// Test functions: func TestFunctionName(t *testing.T)
// Table-driven tests preferred

func TestConfig_Load(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    Config
        wantErr bool
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

## Build & Release

### Version Information
The build process injects version info:
```go
var (
    version = "dev"      // Set via -ldflags at build time
    commit  = "none"     // Git commit SHA
    date    = "unknown"  // Build timestamp
)
```

### Build Commands
```bash
make build       # Build for current platform
make build-all   # Build for all platforms
make test        # Run tests
make fmt         # Format code
make lint        # Run linter
```

### Release Process
1. Use conventional commits: `feat:`, `fix:`, `docs:`, etc.
2. Push tag with `v` prefix: `git tag -a v1.0.0 -m "Release v1.0.0"`
3. CI automatically builds multi-platform binaries and Docker images
4. Changelog generated from commit messages

## CI/CD

### Workflows
- **CI**: Lint, format, test, build on every push/PR
- **Release**: Build binaries and Docker on tag push
- **CodeQL**: Security scanning (weekly + on PR)
- **Auto-Format**: Automatically formats code on PRs

### Supported Platforms
- Linux: amd64, arm64
- macOS: amd64 (Intel), arm64 (Apple Silicon)
- Windows: amd64, arm64

### Docker
- Images published to GHCR: `ghcr.io/smart-fellas/k4a`
- Multi-arch support: `linux/amd64`, `linux/arm64`
- Based on Alpine for minimal size

## Common Patterns in Codebase

### 1. View Switching
```go
switch msg.String() {
case ":topics":
    m.currentView = TopicsView
case ":schemas":
    m.currentView = SchemasView
}
```

### 2. Resource Listing
```go
// Execute kafkactl command
output, err := exec.Command("kafkactl", "get", "topics").Output()
if err != nil {
    return err
}

// Parse output (typically YAML/JSON)
// Update model state
```

### 3. Keyboard Navigation
```go
case key.Matches(msg, m.keys.Up):
    if m.selectedIndex > 0 {
        m.selectedIndex--
    }
case key.Matches(msg, m.keys.Down):
    if m.selectedIndex < len(m.items)-1 {
        m.selectedIndex++
    }
```

### 4. Responsive Sizing
```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    m.header.SetSize(msg.Width, headerHeight)
    m.content.SetSize(msg.Width, msg.Height-headerHeight-footerHeight)
    m.footer.SetSize(msg.Width, footerHeight)
```

## Dependencies Management

### Direct Dependencies
```go
github.com/charmbracelet/bubbles   // TUI components
github.com/charmbracelet/bubbletea // TUI framework
github.com/charmbracelet/lipgloss  // Terminal styling
gopkg.in/yaml.v3                   // YAML parsing
```

### External Tools
- **kafkactl**: Must be installed and configured on user's system
- Required for all Kafka operations

## Development Workflow

1. Create feature branch: `feat/feature-name`
2. Make changes following conventions
3. Run `make fmt` to format code
4. Run `make lint` to check for issues
5. Run `make test` to verify tests pass
6. Commit with conventional commit message
7. Push and create PR
8. CI automatically checks and formats code
9. After merge, tag for release if needed

## Security Considerations

- Never log or expose user credentials
- Configuration files may contain sensitive tokens
- Docker runs as non-root user
- CodeQL scans for security vulnerabilities
- Dependabot updates dependencies weekly

## Future Development Guidelines

### When Adding New Views
1. Create new package under `internal/ui/views/`
2. Implement Bubble Tea Model pattern
3. Register view in main app router
4. Add command in command handler (e.g., `:newview`)
5. Add keyboard shortcut if needed
6. Update help screen
7. Add documentation

### When Adding New Features
1. Check if kafkactl supports the operation
2. Implement in `internal/kafkactl/client.go`
3. Add UI component if needed
4. Update key bindings
5. Write tests
6. Update documentation
7. Use conventional commit message

### Performance Considerations
- Avoid blocking operations in Update()
- Use tea.Cmd for async operations
- Cache kafkactl results when appropriate
- Debounce rapid key presses
- Optimize View() rendering (it's called frequently)

## Common Gotchas

1. **Bubble Tea Update() must not block**: Use tea.Cmd for long operations
2. **Lipgloss styles are immutable**: Chain methods return new styles
3. **Terminal size changes**: Always handle tea.WindowSizeMsg
4. **Key binding conflicts**: Check existing bindings before adding new ones
5. **kafkactl output format**: May vary between versions, parse carefully

## Useful Commands for Copilot

When asked to:
- **"Add a new view"**: Follow the views pattern in `internal/ui/views/`
- **"Add a command"**: Update `internal/ui/components/command/`
- **"Style something"**: Use Lipgloss in `internal/ui/styles/`
- **"Add keyboard shortcut"**: Update `internal/ui/keys/`
- **"Fix formatting"**: Run `make fmt`
- **"Add a test"**: Create `*_test.go` with table-driven tests

## Resources

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Lipgloss Styling](https://github.com/charmbracelet/lipgloss)
- [kafkactl Documentation](https://github.com/michelin/kafkactl)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

---

**Last Updated**: 2025-11-03
**Project Status**: Active Development / POC
**Version**: v0.1.0+

