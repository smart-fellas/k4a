# KafkaCtl TUI

A k9s-inspired Terminal User Interface for Michelin's kafkactl, built with Bubble Tea.

## Features

- Interactive cluster browsing
- Topic management and monitoring
- Consumer group inspection
- Real-time partition and offset monitoring
- Configuration management
- Keyboard-driven navigation

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