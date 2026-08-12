.PHONY: help build run test test-coverage clean deps lint bench dev dummy

# Default target
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make deps           - Install dependencies"
	@echo "  make lint           - Run linter"
	@echo "  make bench          - Run benchmarks"
	@echo "  make dummy          - Generate dummy data"
	@echo "  make dev            - Development mode with hot reload"

# Build the application
build:
	@echo "Building application..."
	go build -o karyawan-app cmd/server/main.go
	@echo "Build complete!"

# Run the application
run:
	@echo "Starting application..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f karyawan-app
	rm -f coverage.out
	rm -f coverage.html
	@echo "Clean complete!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod verify
	@echo "Dependencies installed!"

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# Generate dummy data
dummy:
	@echo "Generating dummy data..."
	@echo "Note: Dummy data generation needs to be implemented as a separate command"

# Development mode with hot reload (requires air)
dev:
	@echo "Starting development mode..."
	@if command -v air > /dev/null; then \
		air -cmd.run "go run cmd/server/main.go"; \
	else \
		echo "Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air -cmd.run "go run cmd/server/main.go"; \
	fi
