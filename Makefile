.PHONY: help build run test clean deps

# Default target
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make deps     - Install dependencies"
	@echo "  make lint     - Run linter"
	@echo "  make bench    - Run benchmarks"

# Build the application
build:
	@echo "Building application..."
	go build -o karyawan-app main.go
	@echo "Build complete!"

# Run the application
run:
	@echo "Starting application..."
	go run main.go

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
	golangci-lint run

# Generate dummy data
dummy:
	@echo "Generating dummy data..."
	go run -c "package main; import 'karyawan-app/dummy_data'; dummy_data.RunDummyDataGenerator()" .

# Development mode with hot reload (requires air)
dev:
	@echo "Starting development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi
