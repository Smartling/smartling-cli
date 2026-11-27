# User Documentation

See the [Wiki](https://github.com/Smartling/smartling-cli/wiki) page for this repository. Per-command reference pages generated from the CLI help are in [`docs/`](docs/smartling-cli.md).

# Development
For developers interested in modifying the tool.

## Building

Local development build (writes `./smartling-cli` for the current platform):

```
go build
```

Full release build (cross-compiles all platforms and produces deb/rpm packages via [GoReleaser](https://goreleaser.com/) nfpms — requires `goreleaser` v2 installed locally):

```
make build
```

The version comes from `CliVersion` in `cmd/helpers/build/build.go`.

Outputs land in `bin/`:

- Per-platform binaries under `bin/smartling-cli_<os>_<arch>_<variant>/smartling-cli`
- Flat release binaries: `bin/smartling-cli_latest_<os>_<arch>` (`.exe` on Windows)
- Linux packages: `bin/smartling_latest_linux_amd64.deb`, `.rpm`, plus arm64 variants
- `bin/checksums.txt`

Package metadata (maintainer, description, license) and cross-compile targets are configured in `.goreleaser.yml`.

To build against a local checkout of [api-sdk-go](https://github.com/Smartling/api-sdk-go), copy `go.work.sample` to `go.work`.

## Testing

```
make test_unit
```

Integration tests (`make test_integration`) run against the live Smartling API and need a built binary and a `smartling.yml` with credentials in `tests/cmd/bin/`.

## Docs

After adding or changing commands or flags, regenerate `docs/`:

```
make docs
```
