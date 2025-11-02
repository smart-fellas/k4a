.PHONY: build run clean test install

BINARY_NAME=k4a
MAIN_PATH=cmd/k4a/main.go

build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

install:
	go install $(MAIN_PATH)

clean:
	go clean
	rm -f bin/$(BINARY_NAME)

test:
	go test ./...

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