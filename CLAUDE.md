# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the Smartling CLI tool - a command-line interface for managing translation files through the Smartling platform. The CLI provides commands for file operations (push, pull, list, status, rename, delete, import), project management, and machine translation services.

## Architecture

The codebase follows a layered architecture pattern:

- **cmd/**: Contains all CLI command definitions using Cobra framework
  - Each command has its own subdirectory with command definition and tests
  - Service initializers are used to inject dependencies into commands
- **services/**: Business logic layer containing service implementations
  - Each service corresponds to a command group (files, projects, init, mt)
  - Services interact with the Smartling API SDK
- **output/**: Rendering and formatting layer for CLI output
- **main.go**: Entry point that wires together all commands and their dependencies

The CLI uses dependency injection through service initializers, making it testable and modular.

## Common Development Commands

### Build
```bash
make all           # Clean, get dependencies, and build for all platforms
make build         # Build for darwin, windows, linux
go build          # Build for current platform
```

### Testing
```bash
make test_unit           # Run unit tests
make test_integration    # Run integration tests (requires binary in tests/cmd/bin/)
go test ./cmd/...        # Run specific unit tests
```

### Code Quality
```bash
make lint         # Run revive linter
make tidy         # Clean up go.mod
make mockery      # Generate mocks using mockery
```

### Package Building
```bash
make deb VERSION=1.0.0    # Build Debian package
make rpm VERSION=1.0.0    # Build RPM package
```

## Configuration

The CLI uses YAML configuration files (smartling.yml) that can be placed in the current directory or parent directories (git-like behavior). Key configuration includes:

- Authentication: user_id and secret (required)
- Project settings: project_id, account_id
- File-specific settings for push/pull operations
- Network settings: proxy, insecure mode

## Key Dependencies

- **Cobra**: CLI framework for command structure
- **Smartling API SDK Go**: Official SDK for Smartling API interaction
- **Bubble Tea**: For interactive TUI components
- **YAML**: Configuration file parsing
- **Mockery**: Mock generation for testing

## Testing Approach

- Unit tests are located alongside source files (*_test.go)
- Integration tests are in tests/cmd/ directory
- Mocks are generated using Mockery and stored in mocks/ subdirectories
- Service layer is fully mocked for command testing

## Development Workflow

1. Commands are defined in cmd/ with Cobra
2. Business logic is implemented in services/
3. Service interfaces are mocked for testing
4. Output formatting is handled in output/ package
5. Configuration is managed through config helpers

The service initializer pattern allows for clean dependency injection and makes the codebase highly testable.