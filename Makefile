# CanIFly Makefile - Simplified Version
# Focuses on developer convenience commands

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
.PHONY: dev
dev: ## Start all development servers
	npm start

# Testing
.PHONY: test
test: ## Run all tests
	npm test

# Building
.PHONY: build
build: ## Build both backend and frontend
	npm run build

# Packaging
.PHONY: package
package: ## Package the application for distribution
	npm run package:app

# Cleaning
.PHONY: clean
clean: ## Clean all build artifacts
	npm run clean

# Quick commands
.PHONY: run
run: dev ## Alias for 'make dev'

# Release
.PHONY: release
release: ## Create a release tag
	npm run release

# Installation
.PHONY: install
install: ## Install all dependencies
	npm install
	cd renderer && npm install
	go mod download

# Version check
.PHONY: version
version: ## Show current version
	@echo "Version: $$(cat version)"
	@echo "Package.json: $$(grep -o '"version": "[^"]*"' package.json | cut -d'"' -f4)"
