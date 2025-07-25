# Variables
BINARY_NAME=server
SRC=.
BUILD_DIR=build
LOG_FILE=server.log
PORT=8000

# Default target to build and run
.PHONY: all
all: build run

# Build the Go binary
.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC)

# Run the server as a background process
.PHONY: run
run:
	@echo "Starting the server on port $(PORT)..."
	nohup $(BUILD_DIR)/$(BINARY_NAME) > $(LOG_FILE) 2>&1 & \
	tail -f $(LOG_FILE)

# Stop the server process
.PHONY: stop
stop:
	@echo "Stopping the server..."
	-pkill -f $(BUILD_DIR)/$(BINARY_NAME)

# Clean up build files
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(LOG_FILE)

# Restart the server
.PHONY: restart
restart: stop all

