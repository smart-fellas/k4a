.PHONY: build run clean test test-unit test-integration test-coverage test-short install build-all release fmt lint deps dev

BINARY_NAME=k4a
MAIN_PATH=cmd/k4a/main.go
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)"

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

install:
	go install $(MAIN_PATH)

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-arm64.exe $(MAIN_PATH)
	@echo "Build complete!"

clean:
	go clean
	rm -f bin/$(BINARY_NAME)

test:
	go test -v ./test/...

test-unit:
	go test -v ./test/unit/...

test-integration:
	go test -v ./test/integration/...

test-coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./test/...
	go tool cover -html=coverage.txt -o coverage.html

test-short:
	go test -short -v ./test/...

deps:
	go mod download
	go mod tidy

# Development helpers
dev:
	air -c .air.toml

fmt:
	go fmt ./...
	gofumpt -w .

lint:
	golangci-lint run

.DEFAULT_GOAL := build