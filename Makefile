BINARY_NAME=server
CMD_DIR=cmd/server
BIN_DIR=bin

all: build

build:
	@echo "Building the application..."
	@go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go

run: build
	@echo "Running the application..."
	@$(BIN_DIR)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test ./...
