# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the Smartling CLI tool - a command-line interface for working with the Smartling translation platform. Top-level command groups (all wired in `main.go`):

- `files`: push, pull, list, status, rename, delete, import
- `projects`: list, info, locales
- `jobs`: list, view, progress, find-by-strings, plus the `files` (add/list/remove), `locales` (add/remove) and `strings` (add/list/remove) sub-groups
- `glossaries`: create, list, import, export
- `mt`: detect, translate (machine translation)
- `init`: creates `smartling.yml` interactively
- `build`: prints version/build info; `docs` (hidden) regenerates `docs/`

## Architecture

The codebase follows a layered architecture pattern:

- **main.go**: Entry point that builds the root command and wires every command group to its service initializer
- **cmd/**: CLI command definitions using Cobra
  - `cmd/cmd_root.go`: root command and global persistent flags (`--config`, `--project`, `--account`, `--user`, `--secret`, `--operation-directory`, `--insecure`, `--proxy`, `--smartling-url`, `--verbose`, `--show-config`)
  - One subdirectory per command group, and one per subcommand (command definition + tests)
  - `cmd/helpers/build/`: version variables (`CliVersion` etc.), injected via ldflags in release builds
  - `cmd/helpers/resolve/`: resolves a value from flag → env var → config → flag default
- **services/**: Business logic, one package per command group (`files`, `projects`, `init`, `mt`, `jobs`, and `glossary` for glossaries). The jobs sub-groups have their own packages: `services/jobs/{files,locales,strings}`
  - Services call the Smartling API SDK (`github.com/Smartling/api-sdk-go`)
  - **services/helpers/**: shared utilities: config loading (`config`), SDK/HTTP client construction (`client`), user-facing error types (`cli_error`), file-name templates (`format`), progress rendering (`progress`), secret-redacting logger (`rlog`), table output (`table`), globbing, and stdin reading
- **output/**: Rendering for `mt` (static and Bubble Tea dynamic renderers) and `projects` output

The CLI uses dependency injection through service initializers, making it testable and modular.

## Common Development Commands

### Build

```bash
go build           # Build ./smartling-cli for the current platform (fast dev iteration; gitignored)
make build         # Cross-compile via GoReleaser (requires `goreleaser` v2 installed)
make all           # clean + go mod download + make build
```

`go build ./...` only checks that the packages compile. It does not write a binary.

`make build` reads the version from `CliVersion` in `cmd/helpers/build/build.go` (bump it there for a release) and runs GoReleaser (`.goreleaser.yml`) in snapshot mode. It cross-compiles for linux/darwin/windows × amd64/arm64 and, for linux, also produces `.deb` and `.rpm` packages via nfpms (package metadata lives in the `nfpms` block).

**Build outputs (under `bin/`):**
- Per-target build dirs: `bin/smartling-cli_<os>_<arch>_<variant>/smartling-cli` (e.g. `bin/smartling-cli_linux_amd64_v1/smartling-cli`, `bin/smartling-cli_windows_amd64_v1/smartling-cli.exe`)
- Flat release binaries: `bin/smartling-cli_latest_<os>_<arch>` (`.exe` on windows)
- Linux packages: `bin/smartling_latest_linux_amd64.deb`, `.rpm`, plus arm64 variants. They install the binary at `/usr/bin/smartling-cli` with a `/usr/bin/smartling` symlink
- `bin/checksums.txt`: SHA256 sums for verification

Release file names contain the literal string `latest`, not the version number.

**Build time:** Cross-compile typically 5-15 seconds; `go build` for one platform is under 5 seconds.

**Local SDK development:** copy `go.work.sample` to `go.work` (gitignored) to build against a sibling `../api-sdk-go` checkout. GoReleaser sets `GOWORK=off`, so release builds always use the SDK version pinned in `go.mod`.

### Testing

#### Unit Tests
Cover command flag resolution and service logic with mocked dependencies:

```bash
make test_unit                                    # go test ./cmd/... ./services/... ./output/...
go test ./cmd/files/push/                         # Run tests for a specific command
go test -v -run TestSpecificFunction ./cmd/...    # Run specific test function with verbose output
```

**Expected results:** All unit tests pass (about 38 packages with tests across `cmd/`, `services/` and `output/`, covering every command group).

**Test time:** About 5 seconds (less when cached).

#### Integration Tests
End-to-end tests in `tests/cmd/` execute a built binary against the live Smartling API.

**Prerequisites:** place a built binary at `tests/cmd/bin/smartling-cli` and a `smartling.yml` with real credentials at `tests/cmd/bin/smartling.yml` (the directory is gitignored).

```bash
make test_integration                             # Runs a hard-coded list of tests/cmd packages
go test ./tests/cmd/files/push/...                # Run a specific integration test
```

`make test_integration` does not include `tests/cmd/glossaries/...`, `tests/cmd/jobs/list`, or `tests/cmd/jobs/view`; run those with `go test` directly.

#### Post-release Smoke Tests
`tests/post-release-tests/` holds a bash smoke suite (`smartling_test_suite.sh`, optionally run in Docker via `docker_setup.sh`); see `test_readme.md` there.

### Code Quality
```bash
make lint         # golangci-lint + revive (install revive with `make _linter`; golangci-lint separately)
make tidy         # go mod tidy
make mockery      # Regenerate mocks from .mockery.yml (install mockery v3 with `make _mockery-install`)
make docs         # Regenerate docs/*.md from the Cobra command tree; rerun after changing commands or flags
```

### Release and CI

- `Jenkinsfile` runs `make build` inside the `goreleaser/goreleaser:v2.17.0` image. On `master`, after a manual "Release to PROD?" approval, it uploads the `bin/` artifacts to the public S3 bucket `smartling-connectors-releases/cli/`.
- CI does not run tests or lint. Run `make test_unit` and `make lint` locally before pushing.
- `/publish-release-notes [tag]` (`.claude/commands/`) drafts GitHub release notes for a tag and publishes them after approval.

## Configuration

The CLI reads a YAML config file, `smartling.yml` (annotated example: `smartling.yml.example`):

- Keys: `user_id` and `secret` (required), `account_id`, `project_id`, `threads` (`files pull` concurrency), `proxy`, and `files` (per-glob push `type`/`directives` and pull `format`)
- `--insecure` / `-k` (skip TLS verification) is a flag only, not a config key
- Lookup: `--config` (or `SMARTLING_CLI_CONFIG`). Otherwise `<operation-directory>/smartling.yml` is searched for starting from `filepath.Dir(os.Args[0])` and walking up the parent directories (`services/helpers/config/constructor.go`). When the binary is invoked through `PATH`, the search starts in the current directory; when it is invoked by path (e.g. `./bin/smartling-cli`), it starts in the binary's directory. The integration tests rely on this.
- Credential env vars `SMARTLING_USER_ID`, `SMARTLING_SECRET`, `SMARTLING_PROJECT_ID` apply only when the config file leaves that value empty; the `--user`/`--secret`/`--project` flags override both.
- Flags resolved through `cmd/helpers/resolve` can also be set with `SMARTLING_CLI_<FLAG>` (the flag name uppercased, hyphens kept as-is). Precedence: flag > env var > config > default.

## Key Dependencies

- **Cobra**: CLI framework (also generates `docs/`)
- **Smartling API SDK Go** (`github.com/Smartling/api-sdk-go`): Smartling API client
- **Bubble Tea / Bubbles / Lip Gloss**: interactive TUI and dynamic progress output
- **goccy/go-yaml**: config file parsing
- **testify + Mockery v3**: assertions and mock generation
- **golang.org/x/sync/errgroup**: concurrency in `files pull` and `mt`

## Testing Approach

- Unit tests are located alongside source files (`*_test.go`); integration tests are in `tests/cmd/`
- Mocks are generated by Mockery from `.mockery.yml`: mocks of this repo's interfaces go in `mocks/` (each service's `Service`, and `SrvInitializer` for files, projects, init, mt), and mocks of api-sdk-go interfaces go in `sdkmocks/`. Register new interfaces in `.mockery.yml` and run `make mockery`; don't hand-edit generated mocks
- Legacy groups (files, projects, init, mt): command tests mock the `SrvInitializer` and `Service`
- Newer groups (jobs, glossaries): command tests cover flag-to-params resolution (`resolve_test.go`), and service tests mock the SDK via `sdkmocks/`

## Development Workflow

1. Commands are defined in cmd/ with Cobra
2. Business logic is implemented in services/
3. Service and SDK interfaces are mocked for testing
4. Output formatting is handled in output/ package
5. Configuration is managed through config helpers

### Service Initializer Pattern

Each command group (files, projects, init, mt, jobs, glossaries) follows the same dependency injection pattern:

1. **Command Group** (e.g., `cmd/files/cmd_files.go`): Defines the group command. Legacy groups (files, projects, init) also define the `SrvInitializer` interface and factory function here; newer groups (mt, jobs, glossaries) put them in `cmd/<group>/srv_initializer.go`
2. **Service Initializer** (e.g., `cmd/files/cmd_files.go`, `cmd/glossaries/srv_initializer.go`): Implements the initializer that wires up SDK clients and configuration
3. **Service** (e.g., `services/files/service.go`): Defines the Service interface with business logic methods
4. **Command Implementation** (e.g., `cmd/files/push/cmd_push.go`): Uses the initializer to get service instance and executes operations

Nested sub-groups (`jobs files`, `jobs locales`, `jobs strings`) repeat the pattern one level down: each has its own `cmd/jobs/<sub>/srv_initializer.go` and its own service in `services/jobs/<sub>/`, wired separately in `main.go`.

To add a new command, use the `add-command` project skill (`.claude/skills/add-command/`).

This pattern enables:
- Easy mocking of services in command tests
- Centralized client and configuration setup
- Clear separation between CLI interface and business logic
