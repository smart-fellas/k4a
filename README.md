# K4A - Kafka for All

A terminal UI for Kafka management using [kafkactl](https://github.com/michelin/kafkactl), inspired by [k9s](https://k9scli.io/).

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)
[![CI](https://github.com/smart-fellas/k4a/actions/workflows/ci.yml/badge.svg)](https://github.com/smart-fellas/k4a/actions/workflows/ci.yml)
[![Release](https://github.com/smart-fellas/k4a/actions/workflows/release.yml/badge.svg)](https://github.com/smart-fellas/k4a/actions/workflows/release.yml)
[![CodeQL](https://github.com/smart-fellas/k4a/actions/workflows/codeql.yml/badge.svg)](https://github.com/smart-fellas/k4a/actions/workflows/codeql.yml)

## Features

- 🎯 **k9s-like Interface**: Familiar navigation and commands for Kubernetes users
- 📊 **Resource Management**: Topics, Schemas, Connectors, Consumer Groups, ACLs
- 🔍 **Quick Navigation**: Use `:` commands to switch between views
- 📝 **YAML View**: Press `d` to describe resources in YAML format
- 🔄 **Real-time Updates**: Refresh views with `r`
- 🎮 **Connector Control**: Pause, Resume, and Restart connectors and their tasks
- 🔐 **Multi-context Support**: Switch between different Kafka environments

## Installation

### Prerequisites

- [kafkactl](https://github.com/michelin/kafkactl) installed and configured
- Go 1.21+ (for building from source)

### Download Pre-built Binary

Download the latest release for your platform from the [releases page](https://github.com/smart-fellas/k4a/releases):

```bash
# Linux (amd64)
curl -LO https://github.com/smart-fellas/k4a/releases/latest/download/k4a-linux-amd64
chmod +x k4a-linux-amd64
sudo mv k4a-linux-amd64 /usr/local/bin/k4a

# macOS (arm64/Apple Silicon)
curl -LO https://github.com/smart-fellas/k4a/releases/latest/download/k4a-darwin-arm64
chmod +x k4a-darwin-arm64
sudo mv k4a-darwin-arm64 /usr/local/bin/k4a

# macOS (amd64/Intel)
curl -LO https://github.com/smart-fellas/k4a/releases/latest/download/k4a-darwin-amd64
chmod +x k4a-darwin-amd64
sudo mv k4a-darwin-amd64 /usr/local/bin/k4a
```

### Docker

Run k4a using Docker:

```bash
# Pull the image
docker pull ghcr.io/smart-fellas/k4a:latest

# Run with your kafkactl config
docker run -it --rm \
  -v ~/.kafkactl:/home/k4a/.kafkactl:ro \
  ghcr.io/smart-fellas/k4a:latest
```

### From Source
```bash
git clone https://github.com/smart-fellas/k4a.git
cd k4a
make build
./bin/k4a
```

### Install from Source
```bash
git clone https://github.com/smart-fellas/k4a.git
cd k4a
make install
k4a
```

## Configuration

K4A uses the same configuration as kafkactl. Ensure you have a properly configured `~/.kafkactl/config.yml`:
```yaml
kafkactl:
  contexts:
    - name: dev
      context:
        api: https://ns4kafka-dev-api.domain.com
        user-token: my_gitlab_token
        namespace: my_namespace
    - name: prod
      context:
        api: https://ns4kafka-prod-api.domain.com
        user-token: my_gitlab_token
        namespace: my_namespace
```

## Usage

### Basic Navigation

- `↑/↓` or `k/j` - Navigate up/down
- `Enter` - Select/drill down
- `ESC` - Go back
- `q` - Quit

### View Commands

- `:topics` - Switch to topics view
- `:schemas` - Switch to schemas view
- `:connectors` - Switch to connectors view
- `:consumers` - Switch to consumer groups view
- `:acls` - Switch to ACLs view

### Resource Actions

- `d` - Describe resource (show YAML)
- `e` - Edit resource
- `Ctrl+d` - Delete resource
- `r` - Refresh view
- `/` - Filter resources

### Connector Actions

- `p` - Pause connector
- `r` - Resume connector
- `R` - Restart connector

### Help

Press `?` to show the help dialog with all available commands.

## Development
```bash
# Run development mode with hot reload
make dev

# Run tests
make test

# Format code
make fmt

# Run linter
make lint
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [k9s](https://k9scli.io/) for the inspiration
- [kafkactl](https://github.com/michelin/kafkactl) for the Kafka management capabilities
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling the terminal UI
## Project Structure

```
k4a/
├── cmd/
│   └── k4a/
│       └── main.go              # Entry point
├── internal/
│   ├── app/
│   │   └── app.go               # Main application model
│   ├── config/
│   │   ├── config.go            # Configuration management
│   │   └── context.go           # Kafka context management
│   ├── ui/
│   │   ├── components/
│   │   │   ├── table/
│   │   │   │   └── table.go    # Reusable table component
│   │   │   ├── header/
│   │   │   │   └── header.go   # Header component
│   │   │   ├── footer/
│   │   │   │   └── footer.go   # Footer with keybindings
│   │   │   ├── dialog/
│   │   │   │   └── dialog.go   # Confirmation/YAML viewer dialogs
│   │   │   ├── command/
│   │   │   │   └── command.go  # Command input (: commands)
│   │   │   └── help/
│   │   │       └── help.go     # Help dialog
│   │   ├── views/
│   │   │   ├── topics/
│   │   │   │   ├── list.go     # Topics list view
│   │   │   │   └── detail.go   # Topic YAML detail view
│   │   │   ├── schemas/
│   │   │   │   ├── list.go     # Schemas list view
│   │   │   │   └── detail.go   # Schema detail view
│   │   │   ├── connectors/
│   │   │   │   ├── list.go     # Connectors list view
│   │   │   │   └── detail.go   # Connector detail view
│   │   │   ├── consumers/
│   │   │   │   └── list.go     # Consumer groups for topic
│   │   │   └── acls/
│   │   │       └── list.go     # ACLs view
│   │   ├── styles/
│   │   │   └── styles.go       # Lipgloss styles
│   │   └── keys/
│   │       └── keys.go         # Keybinding definitions
│   ├── kafkactl/
│   │   ├── client.go           # Kafkactl CLI wrapper
│   │   ├── executor.go         # Command executor
│   │   └── parser.go           # YAML parser
│   └── utils/
│       ├── format.go           # Formatting utilities
│       └── helpers.go          # General helpers
├── pkg/
│   └── models/
│       ├── resource.go         # Base resource interface
│       ├── topic.go            # Topic data model
│       ├── schema.go           # Schema data model
│       ├── connector.go        # Connector data model
│       └── consumer.go         # Consumer group data model
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── .gitignore
```

## Installation

```bash
go install github.com/yourusername/k4a/cmd/k4a@latest
```

## Usage

```bash
# Launch with default context
k4a

# Launch with specific context
k4a --context production

# Launch with custom config
k4a --config ~/.kafkactl/config.yml
```

## Key Bindings

- `?` - Show help
- `/` - Search
- `q` - Quit
- `Tab` - Switch views
- `Enter` - Select/Drill down
- `Esc` - Go back
- `r` - Refresh
- `d` - Describe
- `c` - Create new resource