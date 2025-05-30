# Makefile for GoTrain project

APP_NAME := critical_speed
CMD_DIR := ./cmd/critical_speed
BIN_DIR := ./bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)

# Default target
all: build

# Build the binary
build:
	@echo "🔧 Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_PATH) $(CMD_DIR)

# Run the application
run: build
	@echo "🚀 Running $(APP_NAME)..."
	@$(BIN_PATH)

# Clean build artifacts
clean:
	@echo "🧹 Cleaning up..."
	@rm -rf $(BIN_DIR)

# Tidy modules
tidy:
	@echo "🧽 Tidying up go.mod and go.sum..."
	go mod tidy

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...

.PHONY: all build run clean tidy test
