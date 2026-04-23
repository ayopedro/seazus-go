# Project variables
BINARY_NAME=seazus-api
BUILD_DIR=bin
MAIN_PATH=./cmd/api

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# Linker flags for production
# -s: Omit the symbol table and debug information
# -w: Omit the DWARF symbol table
LDFLAGS=-ldflags="-s -w"

.PHONY: all
all: tidy format test build

.PHONY: tidy
tidy:
	$(GOMOD) tidy
	$(GOMOD) verify

.PHONY: format
format:
	$(GOCMD) fmt ./...

.PHONY: test
test:
	$(GOTEST) -v -race -cover ./...

.PHONY: dev
dev:
	air -c .air.toml

## Build targets
.PHONY: build
build: clean
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

.PHONY: build-linux
build-linux: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build for current OS"
	@echo "  build-linux   Cross-compile for Linux (Production/Docker)"
	@echo "  test          Run tests with race detection"
	@echo "  tidy          Clean up go.mod and verify dependencies"