# Project variables
BINARY_NAME=seazus-api
BUILD_DIR=bin
MAIN_PATH=./cmd/api

# Migrations
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://ayotunde:password@localhost:5432/seazus-go?sslmode=disable
MIGRATIONS_DIR=./migrations

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# Linker flags for production
LDFLAGS=-ldflags="-s -w"

# ----------------------------
# Default target
# ----------------------------
.PHONY: all
all: tidy format test build

# ----------------------------
# Dependency management
# ----------------------------
.PHONY: tidy
tidy:
	$(GOMOD) tidy
	$(GOMOD) verify

# ----------------------------
# Code quality
# ----------------------------
.PHONY: format
format:
	$(GOCMD) fmt ./...

.PHONY: test
test: tidy
	$(GOTEST) -v -race -cover ./...

# ----------------------------
# Development
# ----------------------------
.PHONY: dev
dev: db-up
	air -c .air.toml

# ----------------------------
# Build
# ----------------------------
.PHONY: build
build: clean
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

.PHONY: build-linux
build-linux: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)

# ----------------------------
# Cleanup
# ----------------------------
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# ----------------------------
# Migrations (Goose)
# ----------------------------
.PHONY: migrate-create db-up db-down db-status db-reset

migrate-create:
	goose -dir=$(MIGRATIONS_DIR) create $(name) sql

db-up:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING="$(GOOSE_DBSTRING)" \
	goose -dir=$(MIGRATIONS_DIR) up

db-down:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING="$(GOOSE_DBSTRING)" \
	goose -dir=$(MIGRATIONS_DIR) down

db-status:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING="$(GOOSE_DBSTRING)" \
	goose -dir=$(MIGRATIONS_DIR) status

db-reset:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING="$(GOOSE_DBSTRING)" \
	goose -dir=$(MIGRATIONS_DIR) reset

# ----------------------------
# Help
# ----------------------------
.PHONY: help
help:
	@echo "Usage:"
	@echo ""
	@echo "  make build                 Build binary"
	@echo "  make build-linux           Cross-compile for Linux"
	@echo "  make test                  Run tests"
	@echo "  make dev                   Run dev server"
	@echo ""
	@echo "  make migrate-create name=xxx   Create migration"
	@echo "  make migrate-up                Run migrations"
	@echo "  make migrate-down              Rollback migration"
	@echo "  make migrate-status            Show migration status"
	@echo "  make migrate-reset             Reset all migrations"