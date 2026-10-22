# Reference: smartling-cli command templates

All paths relative to the smartling-cli repo root.

## Architecture recap

```
cmd/<group>/
├── cmd_<group>.go        # group command: New<Group>Cmd()  (name may differ,
│                         #   e.g. cmd/glossaries/cmd_glossary.go)
├── srv_initializer.go    # SrvInitializer interface + NewSrvInitializer (modern)
│                         #   (legacy groups define this inside cmd_<group>.go)
├── mocks/<group>.go      # legacy/mt only: generated SrvInitializer mock
└── <sub>/
    ├── cmd_<sub>.go      # New<Sub>Cmd(initializer) *cobra.Command
    ├── resolve.go        # resolveParams(...)  (modern + mt)
    ├── resolve_test.go   # modern: resolveParams test
    ├── run.go            # run(...)            (modern + mt)
    └── cmd_<sub>_test.go # legacy/mt only: command test w/ mocked Service+SrvInitializer

services/<group>/         # package name may be singular (services/glossary)
├── service.go            # Service interface + service struct + NewService
├── run_<sub>.go          # RunX impl + XParams + XOutput
├── run_<sub>_test.go     # modern: real service vs SDK mock
├── sdkmocks/<area>.go    # generated mock of the SDK api/<area> interface
└── mocks/<group>.go      # generated mock of Service (only if a command test uses it)

output/
├── params.go             # output.Params{ Mode, Format, Template string }
│                         #   modern commands set Format only; Mode is used by mt
│                         #   ("dynamic" render) via resolve.OutputParams
└── static/base.go        # Renderable, OutputFormat[T], GetOutputFormat[T]
```

Root helpers (modern commands use these):
- `rootcmd.Config()` → `(config.Config, error)`
- `rootcmd.Client(ctx)` → `(sdk.HttpAPIClient, error)`
- `rootcmd.ConfigFile()` → `string`
- `rootcmd.RunRootPersistentPreRun(cmd)` — call from a group `PersistentPreRunE`
- `cmd/helpers/resolve` — e.g. `resolve.FallbackAccount(flag, cnf.AccountID)`
- `clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"` — `clierror.UIError`
- `rlog` — `rlog.Debugf/Errorf` logging

---

## MODERN PATTERN

### 1. Group command — `cmd/<group>/cmd_<group>.go`

```go
package glossaries

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Smartling/smartling-cli/cmd"

	"github.com/spf13/cobra"
)

const outputFormatFlag = "output"

var (
	outputFormat   string
	allowedOutputs = []string{"table", "json", "simple"}
	joinedAllowedOutputs = strings.Join(allowedOutputs, ", ")
)

// NewGlossariesCmd builds the `glossaries` command.
func NewGlossariesCmd() *cobra.Command {
	glossariesCmd := &cobra.Command{
		Use:   "glossaries",
		Short: "Manage Smartling glossaries",
		Long:  `...`,
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			if err := cmd.RunRootPersistentPreRun(c); err != nil {
				return err
			}
			if !slices.Contains(allowedOutputs, outputFormat) {
				return fmt.Errorf("invalid output: %s (allowed: %s)", outputFormat, joinedAllowedOutputs)
			}
			return nil
		},
	}
	glossariesCmd.PersistentFlags().StringVar(&outputFormat, outputFormatFlag, "simple", "Output format: "+joinedAllowedOutputs)
	return glossariesCmd
}
```

No `RunE` on the group: a non-runnable group shows usage for a missing or
unknown subcommand instead of silently succeeding.

### 2. Service initializer — `cmd/<group>/srv_initializer.go`

```go
package glossaries

import (
	"context"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	srv "github.com/Smartling/smartling-cli/services/glossary"

	glossaryapi "github.com/Smartling/api-sdk-go/api/glossary"
)

// SrvInitializer defines glossary service initializer
type SrvInitializer interface {
	InitGlossarySrv(ctx context.Context) (srv.Service, error)
}

// NewSrvInitializer returns new SrvInitializer implementation
func NewSrvInitializer() SrvInitializer { return srvInitializer{} }

type srvInitializer struct{}

func (i srvInitializer) InitGlossarySrv(ctx context.Context) (srv.Service, error) {
	client, err := rootcmd.Client(ctx)
	if err != nil {
		return nil, err
	}
	glossaryApi := glossaryapi.NewGlossary(client.Client)
	return srv.NewService(glossaryApi), nil
}
```

### 3. Leaf command — `cmd/<group>/<sub>/cmd_<sub>.go`

```go
package gllist

import (
	"fmt"

	glossariescmd "github.com/Smartling/smartling-cli/cmd/glossaries"
	"github.com/Smartling/smartling-cli/output"

	"github.com/spf13/cobra"
)

const nameFlag = "name"

func NewListCmd(initializer glossariescmd.SrvInitializer) *cobra.Command {
	var name string
	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List glossaries in the current account",
		Long:    `...`,
		Example: `  smartling-cli glossaries list --name "CLI"`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			params, err := resolveParams(cmd, name)
			if err != nil {
				return err
			}
			format, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}
			return run(ctx, initializer, params, output.Params{Format: format})
		},
	}
	listCmd.Flags().StringVar(&name, nameFlag, "", "Filter glossaries by name.")
	return listCmd
}
```

Return `resolveParams` errors **unwrapped**: `main.go` renders errors via
`output/mt.RenderError`, which uses a type assertion `err.(clierror.UIError)`
(not `errors.As`), so wrapping with `fmt.Errorf("...: %w", err)` loses the
structured UIError rendering. (The real `cmd/glossaries/list/cmd_list.go`
wraps it — don't copy that.)

Flag handling: bind each flag inline with `StringVar(&v, ...)` (etc.). Default
to resolving inputs inside `resolveParams` with a helper (as
`cmd/jobs/list/resolve.go` does for `--name`), so they get the env tier.
Passing the raw `StringVar` value in, as `glossaries list` does here, is an
accepted alternative for positional args and inputs nobody would set via env
— it skips the env tier. The account UID always goes through
`FallbackAccount`.

### 4. Resolve — `cmd/<group>/<sub>/resolve.go`

`resolveParams` builds the typed `srv.XParams` the service method consumes.
**One rule for inputs:** anything persistable (it has a `FileConfig` field) or
that a user may reasonably set via env var goes through the
`cmd/helpers/resolve` fallback helpers, so it gets the flag → env → config
chain. Positional args and transient per-invocation inputs may be passed raw
from their `StringVar` value. Don't call `cmd.Flags().GetString(...)` for
service inputs. (The `output` format flag is read directly off the group's
persistent flag in `cmd_x.go`.)

```go
package gllist

import (
	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	srv "github.com/Smartling/smartling-cli/services/glossary"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"

	"github.com/spf13/cobra"
)

func resolveParams(cmd *cobra.Command, name string) (srv.ListParams, error) {
	rlog.Debugf("resolving list params")
	cnf, err := rootcmd.Config()
	if err != nil {
		return srv.ListParams{}, clierror.UIError{
			Operation: "config", Err: err, Description: "failed to read config",
		}
	}
	accountUID, err := resolve.FallbackAccount(
		cmd.Root().PersistentFlags().Lookup("account"), cnf.AccountID)
	if err != nil {
		return srv.ListParams{}, err
	}
	return srv.ListParams{AccountUID: accountUID, Name: name}, nil
}
```

Here `name` is already the flag value (passed from `RunE`), so it's used
directly. Values that aren't simple flag pass-throughs — e.g. the account UID —
go through the fallback helpers below. `rootcmd.Config()` reads the discoverable
`smartling.yml`; the `resolve_test.go` example shows how to control it from a
test.

<a name="file-config"></a>
#### Per-group file config — `cmd/<group>/config.go`

A group that has persistable options defines a typed `FileConfig` mirroring the
YAML under `smartling.yml`, plus `BindFileConfig(cmd) (FileConfig, error)`.
Structure: a top-level `FileConfig` with one nested struct per subcommand,
nested again for grouped flags (e.g. filters).

**IMPORTANT — not every CLI flag belongs in the config file.** Only persistable
options that a user reasonably wants as a default go in `FileConfig` (e.g.
export `file_type`, `locale_ids`, filter settings). Transient/per-invocation
inputs — positional args, the thing being acted on, one-off filters like
`list --name` — stay flag-only and get **no** config field (pass `nil` / `""`
as the `Config` to the fallback helper, or skip the helper entirely).

```go
package glossaries

import (
	"os"

	"github.com/Smartling/smartling-cli/cmd/helpers/resolve"
	"github.com/Smartling/smartling-cli/services/helpers/config"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// FileConfig is nested under the group name in smartling.yml.
type FileConfig struct {
	Glossaries struct {
		Export ExportConfig `yaml:"export,omitzero"`
		Create CreateConfig `yaml:"create,omitzero"`
		Import ImportConfig `yaml:"import,omitzero"`
	} `yaml:"glossaries,omitzero"`
}

// ExportConfig mirrors the *persistable* flags of `glossaries export`.
type ExportConfig struct {
	FileType  string   `yaml:"file_type,omitzero"`
	LocaleIDs []string `yaml:"locale_ids,omitzero"`
	// ...one field per persistable flag, omitzero
}

// BindFileConfig loads the group section from the discovered smartling.yml.
// A missing file is not an error — it returns a zero FileConfig.
func BindFileConfig(cmd *cobra.Command) (FileConfig, error) {
	dir := resolve.ConfigDirectory(cmd)
	filename := resolve.ConfigFile(cmd)
	path, err := config.GetPath(dir, filename, false)
	if err != nil {
		return FileConfig{}, err
	}
	var cfg FileConfig
	data, err := os.ReadFile(path)
	if err != nil && os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
```

**Wiring:** when a subcommand reads file config, its `RunE` binds first, then
passes the config into `resolveParams`, which feeds each field as the `Config`
pointer to a fallback helper:

```go
// cmd_x.go RunE
fileConfig, err := glossariescmd.BindFileConfig(cmd)
if err != nil {
	return err
}
params, err := resolveParams(cmd, fileConfig, glossaryUIDOrName, outFile)

// resolve.go
func resolveParams(cmd *cobra.Command, fileConfig glossariescmd.FileConfig,
	glossaryUIDOrName, outFile string) (srv.ExportParams, error) {
	cfg := fileConfig.Glossaries.Export        // per-subcommand section
	fileType := resolve.FallbackString(
		cmd.Flags().Lookup(fileTypeFlag),
		resolve.StringParam{FlagName: fileTypeFlag, Config: &cfg.FileType})
	// glossaryUIDOrName / outFile are positional/transient — no config field
	// ...
}
```

`&cfg.FileType` is non-nil even when the YAML omits it, so the flag default is
skipped. That's harmless here because `--file-type` defaults to `""`. For flags
with a non-zero default, see the pointer trap under
[Resolving params](#resolving-params).

A command with no persistable options (e.g. `glossaries list`) skips
`config.go` / `BindFileConfig` entirely and resolves only from flags + the
account fallback.

<a name="resolving-params"></a>
#### Resolving params — `cmd/helpers/resolve` helpers

All helpers live in `cmd/helpers/resolve/fallbacks.go`. The env var name is
`SMARTLING_CLI_` + `strings.ToUpper(flagName)` (`env.VarNameFromCLIFlagName`).
**Hyphens are kept**: `--file-type` ↔ `SMARTLING_CLI_FILE-TYPE`.

| Helper | Signature | Precedence |
|--------|-----------|------------|
| `FallbackString` | `(flag *pflag.Flag, StringParam{FlagName, Config *string}) string` | flag → env → config → flag default |
| `FallbackInt` | `(flag *pflag.Flag, IntParam{FlagName, Config *int}) (int, error)` | flag → env → config → flag default |
| `FallbackBool` | `(flag *pflag.Flag, BoolParam{FlagName, Config *bool}) bool` | see traps below |
| `FallbackAccount` | `(flag *pflag.Flag, accountIDConfig string) (uid.AccountUID, error)` | flag → env → config (if non-empty) → flag default; validates |
| `FallbackStringArray` | `(cmd *cobra.Command, flagName string, configVal []string) []string` | flag → env (comma-split) → config (if non-empty) → `nil` — **no flag default** |
| `FallbackDate` | `(cmd *cobra.Command, flagName string, configVal string) (time.Time, error)` | flag → env → config → zero time — **no flag default**; RFC3339 |

Look up a flag with `cmd.Flags().Lookup("name")`, or for a root persistent flag
`cmd.Root().PersistentFlags().Lookup("name")`.

**Traps:**
- **A non-nil `Config` pointer beats the flag default**, even when it points to
  `""`/`0`/`false`. A pointer to a `FileConfig` field is always non-nil, so the
  flag default is never used. That's fine when the default is the zero value
  (e.g. `--file-type`, default `""`); for a flag with a non-zero default, pass
  `nil` when config is empty, as `FallbackAccount` does:
  ```go
  var sortByConfig *string
  if cfg.SortBy != "" {
  	sortByConfig = &cfg.SortBy
  }
  ```
- **`FallbackBool` returns `true`** when the flag is `Changed` (so `--x=false`
  gives `true`), when the env var is set to *any* value, and when the flag exists
  unchanged with no env/config. Always pass an explicit default via `Config`
  (from `cmd/jobs/list/resolve.go`):
  ```go
  withPriorityDefault := false
  withPriority := resolve.FallbackBool(cmd.Flags().Lookup(withPriorityFlag),
  	resolve.BoolParam{FlagName: withPriorityFlag, Config: &withPriorityDefault})
  ```
- **`FallbackStringArray` reads the flag with `GetStringArray`**: declare array
  flags with `StringArrayVar`. A `StringSliceVar` flag silently resolves to `nil`
  (the error is discarded).
- **`IntParam.Config` is `*int`**: convert config fields of other int types
  (e.g. `config.Config.Threads` is `uint32`) before passing them.

```go
// from cmd/glossaries/export/resolve.go (cfg := fileConfig.Glossaries.Export)
fileType := resolve.FallbackString(
	cmd.Flags().Lookup(fileTypeFlag),
	resolve.StringParam{FlagName: fileTypeFlag, Config: &cfg.FileType},
)

localeIDs := resolve.FallbackStringArray(cmd, localeFlag, cfg.LocaleIDs)

createdDate, err := resolve.FallbackDate(cmd, createdDateFlag, cfg.Filter.Created.Date)
if err != nil {
	return srv.ExportParams{}, err
}

// from cmd/jobs/list/resolve.go (no config field)
limit, err := resolve.FallbackInt(cmd.Flags().Lookup(limitFlag), resolve.IntParam{FlagName: limitFlag})
if err != nil {
	return srv.ListParams{}, fmt.Errorf("invalid value for --%s: %w", limitFlag, err)
}
```

### 5. Run — `cmd/<group>/<sub>/run.go`

`run` is the orchestration layer and owns exactly three responsibilities, in
order: **(1) init the service** via the initializer, **(2) execute** the service
method with the resolved params, **(3) render** the typed output via
`static.GetOutputFormat[T](format).FormatAndRender(...)`. Keep business logic
out of here — it stays in the service. Init failures return a
`clierror.UIError`; service errors propagate unchanged.

```go
package gllist

import (
	"context"

	glossariescmd "github.com/Smartling/smartling-cli/cmd/glossaries"
	"github.com/Smartling/smartling-cli/output"
	"github.com/Smartling/smartling-cli/output/static"
	srv "github.com/Smartling/smartling-cli/services/glossary"
	clierror "github.com/Smartling/smartling-cli/services/helpers/cli_error"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
)

func run(ctx context.Context,
	initializer glossariescmd.SrvInitializer,
	params srv.ListParams,
	outputParams output.Params,
) error {
	rlog.Debugf("running glossary list with params: %v", params)
	glossarySrv, err := initializer.InitGlossarySrv(ctx)
	if err != nil {
		return clierror.UIError{
			Operation: "init", Err: err, Description: "unable to initialize Glossary service",
		}
	}
	listOutput, err := glossarySrv.RunList(ctx, params)
	if err != nil {
		return err
	}
	static.GetOutputFormat[srv.ListOutput](outputParams.Format).FormatAndRender(listOutput)
	return nil
}
```

### 6. Service interface — `services/<group>/service.go`

```go
package glossary

import (
	"context"

	api "github.com/Smartling/api-sdk-go/api/glossary"
)

type Service interface {
	RunList(ctx context.Context, params ListParams) (ListOutput, error)
	// ...other RunX methods
}

func NewService(glossaryApi api.Glossary) Service {
	return service{glossaryApi: glossaryApi}
}

type service struct {
	glossaryApi api.Glossary
}
```

### 7. Service impl + Params + Output — `services/<group>/run_<sub>.go`

The `*Output` type must implement `static.Renderable`:
`JSONBytes() []byte`, `SimpleLines() []string`, `TableData() ([]string, [][]string)`.

```go
package glossary

import (
	"context"
	"encoding/json"
	"fmt"

	api "github.com/Smartling/api-sdk-go/api/glossary"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

type ListParams struct {
	AccountUID uid.AccountUID
	Name       string
}

func (p ListParams) Validate() error { return p.AccountUID.Validate() }

type Item struct {
	GlossaryUID, Name, Description string
	LocaleIDs                     []string
}

type ListOutput struct {
	Glossaries []Item
	JSON       []byte
}

func (o ListOutput) JSONBytes() []byte { return o.JSON }

func (o ListOutput) SimpleLines() []string {
	if len(o.Glossaries) == 0 {
		return []string{"No glossaries found."}
	}
	lines := make([]string, 0, len(o.Glossaries))
	for _, g := range o.Glossaries {
		lines = append(lines, fmt.Sprintf("%s  %s  %s", g.GlossaryUID, g.Name, g.Description))
	}
	return lines
}

func (o ListOutput) TableData() ([]string, [][]string) {
	headers := []string{"GLOSSARY UID", "NAME", "DESCRIPTION", "LOCALES"}
	rows := make([][]string, 0, len(o.Glossaries))
	for _, g := range o.Glossaries {
		rows = append(rows, []string{g.GlossaryUID, g.Name, g.Description, fmt.Sprint(g.LocaleIDs)})
	}
	return headers, rows
}

func (s service) RunList(ctx context.Context, params ListParams) (ListOutput, error) {
	if err := params.Validate(); err != nil {
		return ListOutput{}, fmt.Errorf("invalid list params: %w", err)
	}
	resp, err := s.glossaryApi.GetByName(ctx, params.AccountUID, params.Name)
	if err != nil {
		return ListOutput{}, fmt.Errorf("failed to list glossaries: %w", err)
	}
	return toListOutput(resp)
}

func toListOutput(resp []api.GetGlossaryResponse) (ListOutput, error) {
	items := make([]Item, 0, len(resp))
	for _, g := range resp {
		items = append(items, Item{
			GlossaryUID: g.GlossaryUID,
			Name:        g.Name,
			Description: g.Description,
			LocaleIDs:   g.LocaleIDs,
		})
	}
	b, err := json.Marshal(items)
	if err != nil {
		return ListOutput{}, fmt.Errorf("failed to marshal list output: %w", err)
	}
	return ListOutput{Glossaries: items, JSON: b}, nil
}
```

Adapted from `services/glossary/run_list.go` (the real `toListOutput` logs the
marshal error with `rlog` instead of returning it).

### 8. Register — `main.go`

```go
gCmd := glossaries.NewGlossariesCmd()
rootCmd.AddCommand(gCmd)
gInit := glossaries.NewSrvInitializer()
gCmd.AddCommand(gllist.NewListCmd(gInit))
```

### 9. Mocks — `.mockery.yml`

All mocks are generated by mockery (`make mockery`) from `.mockery.yml`. Global
settings there: `template: testify`, generated package `{{.SrcPackageName}}mocks`,
file `{{.SrcPackageName}}.go`, `unroll-variadic: true`. Each entry maps a source
package + interface to an output `dir`. Three kinds of mock exist:

| Mock of | Output dir | Generated pkg | Constructor | Used by |
|---------|-----------|---------------|-------------|---------|
| SDK `api/<area>` interface | `services/<group>/sdkmocks` | `<area>mocks` | `NewMock<Area>(t)` | **service** tests (`run_x_test.go`) |
| service `Service` | `services/<group>/mocks` | `<group>mocks` | `NewMockService(t)` | command tests (legacy groups) |
| `SrvInitializer` | `cmd/<group>/mocks` | `<group>mocks` | `NewMockSrvInitializer(t)` | command tests (legacy groups) |

**What modern groups actually generate:** for jobs/glossaries the file has the
**SDK mock** (e.g. `api/job` → `services/jobs/sdkmocks`) and **no
`SrvInitializer` mock**. (A jobs `Service` mock is generated too, but nothing
imports `services/jobs/mocks` — don't treat it as a precedent.) Modern commands are covered by
`resolve_test.go` + a service test against the SDK mock, so they don't need an
initializer mock. Add only the entries you actually test against.

For a new modern command, the only entry you normally add is the SDK interface
you call (skip it if `.mockery.yml` already lists it):

```yaml
packages:
  github.com/Smartling/api-sdk-go/api/<area>:
    interfaces:
      <Area>:
        config:
          dir: services/<group>/sdkmocks
```

`Service` (`services/<group>/mocks`) and `SrvInitializer` (`cmd/<group>/mocks`)
entries exist only for groups with command-level tests — files, projects, init,
mt (plus the unused jobs `Service` mock). Don't add them for a modern group; copy an existing entry
only if you are extending one of those groups.

Then `make mockery`.

---

<a name="extending-the-sdk"></a>
## EXTENDING THE SDK (api-sdk-go)

Only needed when a command requires an API call the SDK doesn't expose yet.
Work in the **api-sdk-go** checkout (`../api-sdk-go`, as listed in `go.work`);
`go.work` makes the CLI pick it up immediately.

Each API area is one package under `api/<area>` (e.g. `api/job`, `api/glossary`)
with the shape: a public **interface**, a private `http<Area>` struct holding
`client *smclient.Client`, a `New<Area>(client *smclient.Client)` constructor,
and one method per endpoint. Template **adapted from** `api/job/job.go` (the
real `GetJob` has no `smerror` guard; `api/job/locale/locale.go` is a real file
that uses `smerror.ErrEmptyParam`):

```go
package job

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

// Job defines the job behaviour. ADD NEW METHODS HERE.
type Job interface {
	GetJob(ctx context.Context, projectID, jobUID string) (GetJobResponse, error)
	// ...add your method to this interface so it stays mockable
}

func NewJob(client *smclient.Client) Job { return newHttpJob(client) }

type httpJob struct{ client *smclient.Client }

func newHttpJob(client *smclient.Client) httpJob { return httpJob{client: client} }

func (h httpJob) GetJob(ctx context.Context, projectID, jobUID string) (GetJobResponse, error) {
	if projectID == "" {
		return GetJobResponse{}, smerror.ErrEmptyParam("projectID")
	}
	reqURL := path.Join("/jobs-api/v3/projects/", url.PathEscape(projectID),
		"jobs", url.PathEscape(jobUID))

	var response getJobResponse
	_, code, err := h.client.GetJSON(ctx, reqURL, nil, &response.Response.Data)
	if err != nil && code == http.StatusNotFound {
		return GetJobResponse{}, ErrNotFound
	}
	if err != nil {
		return GetJobResponse{}, fmt.Errorf("failed to get job: %w", err)
	}
	return toGetJobResponse(response)
}
```

Rules:
- **If the method doesn't exist, add it to the area's interface** (not just the
  impl) — the service layer depends on the interface and mocks are generated
  from it.
- **Make requests via `client *smclient.Client`** — use its helpers rather than
  building `http.Request`s by hand:
  - `GetJSON(ctx, url, query, &out)` / `DeleteJSON(...)` (removals, e.g.
    `api/job/locale/locale.go`)
  - `PostJSON(ctx, url, payload []byte, &out)` — the body is **already
    marshalled** bytes: `json.Marshal` the request first
  - `Get(...)` / `Post(...)` for raw responses
  - `UploadFile(...)` is **Files-API only** (hard-wired to
    `/files-api/v2/projects/%s/file` with `smfile.FileUploadRequest`), not a
    generic multipart helper — other multipart endpoints build the body
    themselves (see `api/glossary/glossary_import.go`)
- **Return errors from `helpers/sm_error` (package `smerror`) where one fits**
  instead of ad-hoc `errors.New`:
  - `smerror.ErrEmptyParam(name)` — missing required param (`IsErrEmptyParam` to test)
  - `smerror.NotFoundError`, `smerror.NotAuthorizedError`
  - `smerror.APIError`, `smerror.JSONError`, `smerror.ValidationError`
  Wrap lower-level errors with `fmt.Errorf("...: %w", err)`; a package-level
  sentinel like `ErrNotFound` is fine for domain-specific not-found mapping.

After changing an interface, regenerate the SDK mocks the service tests consume
(e.g. `services/<group>/sdkmocks`) via `make mockery`, then `go build ./...` in
the CLI.

<a name="bump-sdk"></a>
**Bump `go.mod` before the CLI change ships.** `go.work` is local only:
`.goreleaser.yml` builds with `GOWORK=off`, and `go.mod` pins an
`api-sdk-go` pseudo-version. Until it's bumped, the CLI builds for you but
breaks `make build` and CI. Once the SDK change is merged and pushed:

```bash
go get github.com/Smartling/api-sdk-go@<sdk-commit-sha>
make tidy
GOWORK=off go build ./...   # proves the build no longer needs the local SDK
```

Tell the developer this step is pending if the SDK change isn't merged yet.

---

<a name="sub-resource-commands"></a>
## SUB-RESOURCE COMMANDS (the `jobs locales|strings|files` pattern)

Adding `add`/`remove`/`list` for a job sub-resource — a **nested subgroup**
(`jobs <sub> <verb>`), as used by `jobs locales|strings|files`.

1. **Nested group command** `cmd/jobs/<sub>/cmd_job_<sub>.go` (`New…Cmd()`): real
   Short/Long/Example, and **no `RunE` and no `PersistentPreRunE`**. Omitting
   `RunE` makes the group non-runnable, so an unknown subcommand (`jobs locales
   ad`) shows usage instead of **silently succeeding**; omitting
   `PersistentPreRunE` lets it inherit the parent `jobs` group's `--output` flag
   and auth/output-validation pre-run. Defining either on the subgroup overrides
   the parent's and breaks inheritance.
2. **Own initializer** `cmd/jobs/<sub>/srv_initializer.go` (`InitJob<Sub>Srv`)
   builds the sub-resource SDK API **and** `jobapi.NewJob(client.Client)` (needed
   for name→UID resolution). **Own service** in `services/jobs/<sub>/` (e.g.
   `services/jobs/locales`), with its SDK mock generated to
   `services/jobs/<sub>/sdkmocks` (see the `api/job/locale` entry in
   `.mockery.yml`).
3. **SDK sub-package** `api/job/<sub>` (e.g. `api/job/file`): `interface` +
   `http<Sub>` impl + `New<Sub>(client *smclient.Client)`. Batch mutations whose
   endpoint reports counts (`file`, `string`) return
   `Result{SuccessCount, FailCount}`; single-item ones (`locale` `Add`/`Remove`)
   return only `error`. `smerror.ErrEmptyParam` guards; segments
   `url.PathEscape`d. `httptest` tests: success-envelope, 404→wrap, empty-param
   table. (Each sub-package redeclares `jobBasePath`, and `Result` where needed —
   accepted duplication for package independence.)
4. **Mutation semantics (add/remove) — "no error" ≠ "it worked".** The API
   returns HTTP 200 + `{successCount, failCount}` **even when nothing matched**
   (nonexistent hashcode/fileUri is silently ignored — counted in neither).
   So pass a real result (never `result=nil`), surface the counts in the Output,
   add a hint when `0/0`, and return a **non-zero `clierror`** from `run.go` if
   any item errored.
5. **Name→UID**: `jobresolver.GetJobUID(ctx, jobAPI, projectID, jobUIDOrName)`;
   map `jobapi.ErrNotFound` → "no job found for %q" in `run.go`.
6. **Required flags**: `cmd.MarkFlagRequired("<flag>")` right after declaring
   (fail fast at parse), guarding its setup error the way jobs does —
   `rlog.Errorf` + `os.Exit(1)` (not `panic`; mt instead uses
   `output.RenderAndExitIfErr(clierror.UIError{...})`, with `output` imported as
   `output "github.com/Smartling/smartling-cli/output/mt"`). Keep service
   `Validate()` as defense-in-depth.
7. **Glob pattern flags** (`--file`): expand against the project's files with a
   single `ListAllFiles` call + `gobwas/glob` — **not** a per-pattern
   `globfiles.Remote` loop (it errors on no-match and lists per call). One API
   call per resolved item; aggregate per-item results, report unmatched
   patterns. **Beware unbounded fan-out** from broad patterns (`**`).
8. **Prefer logging in `run.go`, not the service layer.** `rlog.Debugf`
   nil-panics unless the logger is initialized, so a service package that logs
   needs a `TestMain` calling `rlog.Init()` (as `services/jobs/main_test.go`
   does for `RunList`/`RunView`/`RunFindByStrings`).
9. **Minimize coupling**: if a consumer needs only one method of an SDK
   interface, inject a **function type** (e.g. `ListJobFilesFn`), not the whole
   interface.

### Restructuring a leaf into a group
When a leaf (`jobs files <job>`) becomes a group + subcommand (`jobs files list`):
- `mv cmd/jobs/files/cmd_files.go cmd/jobs/files/list/cmd_list.go` (+ its
  `run.go`) — plain `mv`, not `git mv`, which stages the change; change
  `package`, rename the constructor, fix `Use`.
- Add the group command (no `RunE`) + initializer.
- Migrate the leaf's service method into the group's service and rename it
  (`RunFiles`→`RunList`); remove it from the old service; `make mockery`.
- Update `main.go` and the relocated integration test (the binary path depth
  changes, e.g. `../../bin` → `../../../bin`).

---

## HYBRID PATTERN (mt)

`mt` sits between the two: leaf commands use `RunE` with `resolve.go` +
`run.go`, but `run.go` renders through `mtoutput.InitRender` (a streaming
render loop) rather than `static.GetOutputFormat`, and the group keeps a
`SrvInitializer` mock plus `cmd_<sub>_test.go`/`run_test.go`. Output params
come from `resolve.OutputParams(cmd, fileConfig.MT.FileFormat)`, which sets
`Format`, `Template` (`--format`) and `Mode` (`--output-mode`; `"dynamic"`
selects the live renderer). Follow `cmd/mt/detect/` when extending mt; use the
modern pattern everywhere else.

---

## LEGACY PATTERN (match only when editing files/projects/init groups)

- `SrvInitializer` is defined inside `cmd/<group>/cmd_<group>.go`.
- Initializer also pulls `cmd.Config()` + `cnf.GetFileConfig(cmd.ConfigFile())`
  and injects multiple SDK clients into `NewService(...)`.
- Leaf command uses `Run:` (not `RunE`), reads `args`, calls
  `initializer.InitFilesSrv(ctx)`, then `s.RunX(...)`, and on error
  `rlog.Errorf(...)` + `os.Exit(1)`.
- Service method signature is `RunX(ctx, ...) error` and does its own
  rendering (e.g. via `services/helpers/table` + `os.Stdout`).
- Flags declared with `cmd.Flags().BoolVarP / StringVar` on package-level vars.

See `cmd/files/list/cmd_list.go` and `services/files/run_list.go` for the
canonical legacy example.

---

<a name="integration-tests"></a>
## INTEGRATION TESTS

Integration tests live under `tests/cmd/<group>/<sub>/<sub>_test.go` (package
named after the subcommand, except Go keywords or collisions — e.g.
`tests/cmd/glossaries/import` uses `package glimport`). They are **black-box**: each test execs the
**prebuilt binary** at `tests/cmd/bin/smartling-cli` and asserts on its combined
stdout/stderr — no Go imports from the CLI, no mocks. They hit the **real
Smartling API**, so the developer runs them manually after placing the binary
**and a `smartling.yml` config (real creds)** in `tests/cmd/bin/`. **Never run
them yourself** (`make test_integration` is off-limits) — just author the file.

Shape (based on `tests/cmd/files/pull/pull_test.go`, which checks both
`wantErr` directions; many other integration tests, e.g. the glossaries ones,
only check one — don't copy that gap):

- Resolve `absDir` relative to the test's depth: `"../../bin/"` for
  `tests/cmd/<group>/<sub>/`, `"../../../bin/"` for a nested subgroup
  `tests/cmd/<group>/<sub>/<verb>/` (e.g. `tests/cmd/jobs/locales/add`).
  `exec.Command("./smartling-cli", args...)` with `testCmd.Dir = absDir`;
  capture `CombinedOutput()`.
- Check both directions of `wantErr`: fail on an error when `!wantErr`, **and**
  fail when `wantErr` but the command exited 0.
- Table-driven cases with: `args`, `expectedOutputs` (substrings that must
  appear), `unexpectedOutputs` (must NOT appear — always include `"DEBUG"`,
  `"ERROR"` for success cases), `wantErr`, optional `cleanup func()`.
- Group persistent flags go before the subcommand in `args`, e.g.
  `{"glossaries", "--output", "table", "export", glossaryUID, "--file-type", "csv"}`.
- **Self-contained data**: each test package creates (and, where a delete
  command exists, removes) the data it needs. Don't depend on another package's
  test having run first — packages run in parallel, in no fixed order. If there
  is no delete command, use unique per-run names (e.g. a timestamp suffix) so
  reruns don't collide.

Cases worth covering: `--help` text; each validation/error branch (`wantErr:
true` + the expected message); the happy path; and each output format
(`--output table` → column headers, `--output json` → JSON field names).

**Register the new dir** in the Makefile `test_integration` target so it runs
with the suite — **unless an existing `.../...` line already covers it** (e.g.
`./tests/cmd/jobs/locales/...`, `./tests/cmd/projects/...`), which would run it
twice:

```make
test_integration:
	...
	go test ./tests/cmd/<group>/<sub>/...
```

---

## Commands

Prefer the **Makefile** targets — they fix the exact package sets. Don't
hand-roll `go test`/`mockery`. They do **not** pin tool versions: `make mockery`
and `make lint` run whatever is on `PATH`. The pinned versions live in the
`_mockery-install` (mockery v3.3.4) and `_linter` targets, which nothing
depends on. Before `make mockery`, check `mockery version`; if it isn't v3.3.4,
run `make _mockery-install` (or tell the developer). `make mockery` regenerates
**every** mock, so revert unrelated mock diffs (e.g. `interface{}` → `any`
churn) and keep only the mocks your change needs.

```bash
go build ./...        # quick compile check (only raw go command used directly)
make _mockery-install # installs the pinned mockery v3.3.4
make mockery          # mockery --config .mockery.yml (after interface changes)
make test_unit        # go test ./cmd/... ./services/... ./output/...
make test_integration # DO NOT RUN — slow; the developer runs it manually
make lint             # golangci-lint run ./... && revive -config .revive.toml ./...
make docs             # DO NOT RUN — the developer regenerates ./docs at the end
make tidy             # go mod tidy
```