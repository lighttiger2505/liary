# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**liary** is a CLI tool for fast diary creation and management written in Go. It allows users to quickly create, edit, and search markdown-based diary entries organized by date.

## Architecture

### Core Components

- **main.go**: Entry point using urfave/cli framework with commands for edit, append, list, find, grep, config, and move operations
- **cmd/**: Command implementations for each CLI subcommand
- **internal/**: Core business logic including:
  - **config.go**: YAML-based configuration management with workspace support
  - **diary.go**: Date-based file path generation and diary listing functionality
  - **date.go**: Date parsing and manipulation utilities
  - **file.go**: File system operations
  - **ui/**: User interface utilities
  - **cmd.go**: Helper functions for terminal interaction and external command execution

### File Organization Structure

Diary files are organized hierarchically by date:
```
diary_dir/
├── YYYY/
│   └── MM/
│       └── DD.md (or DD-suffix.md)
```

### Configuration System

- Config file location: XDG config directory + `config.yml`
- Supports multiple workspaces for different diary collections
- Configurable editor, editor options, and grep command
- Auto-creates default configuration with vim as default editor

## Development Commands

### Building and Testing
```bash
# Build the project
make build

# Install locally
make install

# Run tests
make test

# Run linting
make lint

# Generate test coverage report
make coverage
```

### Cross-platform Building
```bash
# Build for multiple platforms
make cross-build

# Create distribution packages
make dist
```

### Dependencies
- Uses Go 1.21+
- Key dependencies: urfave/cli, go-fuzzyfinder, golang.org/x/crypto, gopkg.in/yaml.v2

### Testing Strategy
- Unit tests for core functionality (date parsing, diary operations)
- Test files: `internal/*_test.go`
- CI runs tests via GitHub Actions with Go 1.21.x

### Code Style
- Uses golangci-lint for code quality
- CI enforces linting via golangci/golangci-lint-action@v3