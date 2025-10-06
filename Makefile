# Makefile for GoTrain project

APP1_NAME := critical_speed
APP2_NAME := runner

CMD1_DIR := ./cmd/critical_speed
CMD2_DIR := ./cmd/runner

BIN_DIR := ./bin
BIN1_PATH := $(BIN_DIR)/$(APP1_NAME)
BIN2_PATH := $(BIN_DIR)/$(APP2_NAME)

# Default target: build everything
all: build

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

# Tidy modules
tidy:
	@echo "🧽 Tidying go.mod and go.sum..."
	@go mod tidy

# Build all apps
build: fmt tidy $(BIN1_PATH) $(BIN2_PATH)

# Build critical_speed binary
$(BIN1_PATH):
	@echo "🔧 Building $(APP1_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN1_PATH) $(CMD1_DIR)

# Build runner binary
$(BIN2_PATH):
	@echo "🔧 Building $(APP2_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN2_PATH) $(CMD2_DIR)

# Run critical_speed
run-critical: $(BIN1_PATH)
	@echo "🚀 Running $(APP1_NAME)..."
	@$(BIN1_PATH)

# Run runner
run-runner: $(BIN2_PATH)
	@echo "🚀 Running $(APP2_NAME)..."
	@$(BIN2_PATH)

# Clean build artifacts
clean:
	@echo "🧹 Cleaning up..."
	@rm -rf $(BIN_DIR)

# Run tests
test: fmt tidy
	@echo "🧪 Running tests..."
	go test ./...

.PHONY: all build clean fmt tidy test run-critical run-runner
