# User Documentation

See the [Wiki](https://github.com/Smartling/smartling-cli/wiki) page for this repository.

# Development
For developers interested in modifying the tool.

## Building

Local development build:

```
go build ./...
```

Full release build (cross-compiles all platforms and produces deb/rpm packages via [GoReleaser](https://goreleaser.com/) nfpms — requires `goreleaser` installed locally):

```
make build
```

Outputs land in `bin/`:

- Per-platform binaries under `bin/smartling-cli_<os>_<arch>_<variant>/smartling-cli`
- Linux packages: `bin/smartling_<version>_linux_amd64.deb`, `.rpm`, plus arm64 variants
- `bin/checksums.txt`

Package metadata (maintainer, description, license) and cross-compile targets are configured in `.goreleaser.yml`.
