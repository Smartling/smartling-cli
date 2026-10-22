---
name: add-command
description: Step-by-step recipe for adding a new command to the smartling-cli Go codebase, following its Cobra + service-initializer DI conventions. Use when adding/scaffolding a new CLI command or subcommand (e.g. "add a files X command", "new glossaries subcommand", "scaffold a smartling-cli command"), or wiring a command to a service and tests.
---

# Adding a smartling-cli command

Scaffolds a new Cobra command in the layered architecture: **cmd/** (CLI) →
**services/** (business logic) → **output/** (rendering), wired by a
per-group `SrvInitializer`.

**Two-repo workspace:** this repo + an `api-sdk-go` checkout next to it
(`../api-sdk-go`), wired by a local, gitignored `go.work` copied from the
tracked `go.work.sample` (never `go mod edit -replace`). Add missing API calls
in the **api-sdk-go** checkout; `go.work` picks them up. See
[EXAMPLES.md](EXAMPLES.md#two-repo-workspace-read-before-adding-sdk-calls).

## Decide the pattern first

**Use the modern pattern for all new commands**, unless the target group still
uses the legacy one (then match it). Templates for both in [REFERENCE.md](REFERENCE.md).

- **Modern** (glossaries, jobs; `cmd/glossaries/*` is the canonical precedent):
  `RunE` + split `cmd_*.go`/`resolve.go`/`run.go` (some simple jobs commands —
  `view`, `progress`, `files list` — skip `resolve.go`; new commands shouldn't);
  service returns a typed `*Output` (`static.Renderable`) rendered via
  `static.GetOutputFormat[T](format)`; errors as `clierror.UIError`.
- **Hybrid** (mt): `RunE` + `resolve.go`/`run.go` split, but renders through
  `mtoutput.InitRender` and keeps `SrvInitializer` mocks + `cmd_*_test.go`.
  Match it when extending mt; don't copy it into other groups.
- **Legacy** (files, projects, init): `Run` + `os.Exit(1)`; service self-renders.

## How to execute

Work **interactively in this single session — no sub-agents, no parallelizing.**
Do one step at a time and **pause for the developer to approve** (show the diff,
wait for an OK) before the next. The developer owns commits, `make docs`, and
integration tests.

## Recipe (modern pattern)

Adding subcommand `X` to group `G`, one approved step at a time:

1. **Group exists?** If `cmd/G/` is missing, create the group: `cmd/G/cmd_G.go`
   (`NewGCmd()`) + `cmd/G/srv_initializer.go` (`SrvInitializer` interface +
   `NewSrvInitializer`). Add a `services/G/service.go` with the `Service`
   interface + `NewService`. Existing names don't always follow the template
   (`cmd/glossaries/cmd_glossary.go`, `services/glossary`) — check the group
   before assuming a path. **Nested subgroup** (`jobs locales|strings|files`):
   group has **no `RunE`/`PersistentPreRunE`** — see
   [REFERENCE.md](REFERENCE.md#sub-resource-commands).
2. **SDK call** (only if the API isn't already exposed): add the method to the
   relevant `api/<area>` **interface** + its `http<Area>` impl in the
   **api-sdk-go** checkout. Make requests through `client *smclient.Client`
   (`GetJSON`/`PostJSON`/`DeleteJSON`/`Get`/`Post`/`UploadFile`) and return errors from
   `helpers/sm_error` (package `smerror`) where one fits. See
   [REFERENCE.md](REFERENCE.md#extending-the-sdk). Before the CLI change
   ships, `go.mod` must be bumped to the merged SDK commit (`go.work` is
   local; release builds use `GOWORK=off`) — see
   [REFERENCE.md](REFERENCE.md#bump-sdk).
3. **Service method**: add `RunX(ctx, params) (XOutput, error)` to the `Service`
   interface in `services/G/service.go`; implement in `services/G/run_x.go`.
   Define `XParams` (with `Validate()`) and `XOutput` (implementing
   `JSONBytes()`, `SimpleLines()`, `TableData()`) in that file.
4. **Command files** in `cmd/G/x/`:
   - `cmd_x.go` — `NewXCmd(initializer G.SrvInitializer) *cobra.Command`, `RunE`
     calling `resolveParams` then `run`. Declare flags inline (`StringVar`, …).
     Build `output.Params{Format: …}` from `--output`; modern commands don't set
     `Mode` (only mt uses it, via `resolve.OutputParams`).
   - `resolve.go` — `resolveParams(cmd, [fileConfig,] …) (srv.XParams, error)`:
     read config via `rootcmd.Config()` (some jobs commands read it in `RunE`
     and pass values in — also fine). **Rule:** resolve inputs with the
     `cmd/helpers/resolve` fallback helpers (flag → env → config) by default;
     passing a raw flag value from `RunE` is acceptable only for positional
     args and inputs nobody would set via env. Return `resolveParams` errors
     unwrapped so `clierror.UIError` keeps its rendering. Mind the helper traps (`FallbackBool`, non-nil `Config`
     pointers, `StringArrayVar`, hyphenated env names) — see
     [REFERENCE.md](REFERENCE.md#resolving-params).
   - **If** the command has persistable options, add `cmd/G/config.go` (typed
     `FileConfig` + `BindFileConfig`) and bind it in `RunE` before resolving.
     **Not every flag belongs in config** — transient inputs (positional args,
     one-off filters) stay flag-only. See [REFERENCE.md](REFERENCE.md#file-config).
   - `run.go` — `run(ctx, initializer, params, outputParams) error`: init
     service, call `RunX`, render via `static.GetOutputFormat`.
5. **Register** in `main.go`: `gCmd.AddCommand(x.NewXCmd(gSrvInitializer))`.
6. **Mocks**: add new package(s) to `.mockery.yml` if needed, then `make mockery`
   (regenerate SDK mocks too if you changed an `api/<area>` interface).
7. **Tests** (fake IDs only). Modern commands use two tests and **no**
   `SrvInitializer` mock:
   - `cmd/G/x/resolve_test.go` — **always** add. Build a root→child command tree,
     control `rootcmd.Config()` via a temp `smartling.yml` + `--config`/`--account`
     flags, call `resolveParams` directly; cover flag→env→config precedence and
     errors. See [EXAMPLES.md](EXAMPLES.md#resolveparams-test).
   - `services/G/run_x_test.go` — the main logic test: real service vs SDK mock,
     set expectations, call `RunX`, assert the typed output.
   - `tests/cmd/G/x/x_test.go` — black-box integration test that execs the
     prebuilt binary. **Author it + register its dir in the Makefile
     `test_integration` target (unless an existing `.../...` line already
     covers it), but do not run it.** Keep its data self-contained — don't rely
     on another test package's fixtures. See
     [REFERENCE.md](REFERENCE.md#integration-tests).
   - *Legacy/mt groups only*: a `cmd_x_test.go` mocking `Service` + `SrvInitializer`.
8. **Verify**: `make mockery`, `make test_unit` (runs `cmd/... services/...
   output/...`), then `make lint`. Use the **Makefile** targets, don't hand-roll
   `go test`/`mockery`. The targets don't pin tool versions — check
   `mockery version` is v3.3.4 (else `make _mockery-install`) and revert
   unrelated regenerated-mock diffs. If an SDK change is involved, also run
   `GOWORK=off go build ./...` once `go.mod` is bumped. **Do not run `make test_integration`** (slow, real API)
   **or `make docs`** — the developer runs those at the end.

## Conventions to respect

- Package names: leaf packages are short (`gllist`, not `list`) when the dir
  name would collide; match the existing group's style.
- Errors in modern commands wrap with `clierror.UIError{Operation, Err, Description}`.
- Config file mirrors only *persistable* flags; transient inputs stay flag-only.
- Required flags: `cmd.MarkFlagRequired(...)` (fail fast); handle its setup
  error the way the group already does — jobs: `rlog.Errorf` + `os.Exit(1)`,
  mt: `output.RenderAndExitIfErr(clierror.UIError{...})` with
  `output "github.com/Smartling/smartling-cli/output/mt"`. Keep service
  `Validate()` too.
- Mutating add/remove: API returns 200 + `{successCount,failCount}` even for
  no-ops — surface the counts (not `result=nil`), non-zero exit on item errors.
- Prefer logging in `run.go`. If a service logs with `rlog`, its package needs a
  `TestMain` calling `rlog.Init()` (see `services/jobs/main_test.go`), or tests
  nil-panic.
- No real credentials/IDs in tests (use `"test-account-uid"`).
- Minimal comments — only where the *why* is non-obvious.
- **Never commit.** Do not run `git add`/`commit`/`push` — the developer makes
  all commits. Leave changes in the working tree for them to review.

See [REFERENCE.md](REFERENCE.md) for copy-paste templates and
[EXAMPLES.md](EXAMPLES.md) for a complete worked command.