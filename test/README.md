# Test Directory

This directory contains all test files for the K4A project.

## Structure

```
test/
├── integration/     # Integration tests
├── unit/           # Unit tests
└── fixtures/       # Test fixtures and mock data
```

## Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./test/...

# Run specific test package
go test ./test/unit/...
go test ./test/integration/...
```

## Writing Tests

Follow the table-driven test pattern:

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
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

## Test Coverage

Tests should maintain at least 70% code coverage for new features.

