# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the Smartling CLI tool - a command-line interface for managing translation files through the Smartling platform. The CLI provides commands for file operations (push, pull, list, status, rename, delete, import), project management, machine translation services, and job progress tracking.

## Architecture

The codebase follows a layered architecture pattern:

- **cmd/**: Contains all CLI command definitions using Cobra framework
  - Each command has its own subdirectory with command definition and tests
  - Service initializers are used to inject dependencies into commands
- **services/**: Business logic layer containing service implementations
  - Each service corresponds to a command group (files, projects, init, mt, jobs)
  - Services interact with the Smartling API SDK
  - **services/helpers/**: Shared utilities for config management, error handling, progress rendering, format compilation, and thread pooling
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
make test_unit                                    # Run all unit tests
make test_integration                             # Run all integration tests (requires binary in tests/cmd/bin/)
go test ./cmd/...                                 # Run all unit tests in cmd/
go test ./cmd/files/push/                         # Run tests for a specific command
go test ./tests/cmd/files/push/...                # Run specific integration test
go test -v -run TestSpecificFunction ./cmd/...    # Run specific test function
```

### Code Quality
```bash
make lint         # Run golangci-lint and revive linter
make tidy         # Clean up go.mod
make mockery      # Generate mocks using mockery (config: .mockery.yml)
make docs         # Generate command documentation
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

### Service Initializer Pattern

Each command group (files, projects, init, mt, jobs) follows the same dependency injection pattern:

1. **Command Group** (e.g., `cmd/files/cmd_files.go`): Defines the `SrvInitializer` interface and factory function
2. **Service Initializer** (e.g., `cmd/files/cmd_files.go`): Implements the initializer that wires up SDK clients and configuration
3. **Service** (e.g., `services/files/service.go`): Defines the Service interface with business logic methods
4. **Command Implementation** (e.g., `cmd/files/push/cmd_push.go`): Uses the initializer to get service instance and executes operations

This pattern enables:
- Easy mocking of services in command tests
- Centralized client and configuration setup
- Clear separation between CLI interface and business logic