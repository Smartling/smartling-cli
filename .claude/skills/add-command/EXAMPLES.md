# Examples

## Two-repo workspace (read before adding SDK calls)

Local dev spans two sibling checkouts, wired by `go.work` — do **not**
`go mod edit -replace`:

- `smartling-cli/` — the CLI (this repo)
- `../api-sdk-go/` — the Smartling API SDK

`go.work` is gitignored; only `go.work.sample` (`use .` + `use ../api-sdk-go`)
is tracked. If `go.work` is missing, copy the sample and make sure its `go`
line is >= the `go` line in `go.mod`, or builds fail.

`go.work` makes the CLI consume the **local** SDK checkout, so SDK types you
see in the module cache may be stale. When inspecting or extending SDK types
(e.g. `api/glossary`, `api/batches`, `api/job`), read/edit them in the
`api-sdk-go` checkout. The SDK repo mirrors the CLI's branch names.

If `RunX` needs an API call the SDK doesn't expose yet:

1. Add the method to the relevant `api/<area>` **interface** in `api-sdk-go`
   (not just the impl, so it stays mockable), implemented on the `http<Area>`
   struct.
2. Make the request through `client *smclient.Client`
   (`GetJSON`/`PostJSON`/`DeleteJSON`/`Get`/`Post`/`UploadFile`), and return errors from
   `helpers/sm_error` (package `smerror`) where one fits — e.g.
   `smerror.ErrEmptyParam("projectID")`. Full guidance + template:
   [REFERENCE.md](REFERENCE.md#extending-the-sdk).
3. `go build ./...` in the CLI picks it up via `go.work` immediately — but only
   locally. Release builds use `GOWORK=off`, so once the SDK change is merged,
   bump `go.mod` (`go get github.com/Smartling/api-sdk-go@<sha>`, `make tidy`)
   and confirm with `GOWORK=off go build ./...`. See
   [REFERENCE.md](REFERENCE.md#bump-sdk).
4. Regenerate any SDK mocks the service tests use (e.g. `services/glossary/sdkmocks`).

---

## How modern commands are tested

Modern groups (jobs, glossaries) cover a command with **two** unit tests — and
**no** mocked `SrvInitializer` (that mock isn't even generated for them):

1. `cmd/<group>/<sub>/resolve_test.go` — `resolveParams` directly (no mocks).
2. `services/<group>/run_<sub>_test.go` — the **real service** against the **SDK
   mock** (`sdkmocks`), asserting on the returned typed output.

The `Service`+`SrvInitializer` command test (below) is only used by groups that
generate those mocks: the legacy groups (files/projects/init) and the hybrid mt
group. For a new modern command, write the two tests above and skip the
initializer mock. Fake IDs only.

<a name="resolveparams-test"></a>
## resolveParams test — `cmd/<group>/<sub>/resolve_test.go`

Always add this after creating `resolveParams` — it's a fast, mock-free unit
test of the flag→env→config→error precedence. Because `resolveParams` calls
`rootcmd.Config()`, the real pattern builds a root→child command tree and
controls config via a temp `smartling.yml` + the `--config`/`--account`
persistent flags (mirrors `cmd/glossaries/list/resolve_test.go`). Fake IDs only.

```go
package gllist

import (
	"os"
	"path/filepath"
	"testing"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	"github.com/Smartling/smartling-cli/services/glossary"

	"github.com/Smartling/api-sdk-go/helpers/uid"
	"github.com/spf13/cobra"
)

func TestMain(m *testing.M) {
	rootcmd.ConfigureLogger()
	m.Run()
}

// makeCmd builds a root→child tree and binds --config (and optionally --account)
// on the root so rootcmd.Config() is controllable from the test.
func makeCmd(t *testing.T, configPath, accountFlag string) *cobra.Command {
	t.Helper()
	root := rootcmd.NewRootCmd()
	child := &cobra.Command{Use: "test"}
	root.AddCommand(child)
	_ = root.PersistentFlags().Set("config", configPath)
	t.Cleanup(func() { _ = root.PersistentFlags().Set("config", "") })
	if accountFlag != "" {
		_ = root.PersistentFlags().Set("account", accountFlag)
	}
	return child
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "smartling.yml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func Test_resolveParams(t *testing.T) {
	// Config() always requires these.
	t.Setenv("SMARTLING_USER_ID", "test-user")
	t.Setenv("SMARTLING_SECRET", "test-secret")
	noConfig := filepath.Join(t.TempDir(), "no-smartling.yml") // missing file is tolerated

	tests := []struct {
		testName string
		setup    func(t *testing.T) *cobra.Command
		name     string
		want     glossary.ListParams
		wantErr  bool
	}{
		{
			testName: "account from --account flag",
			setup:    func(t *testing.T) *cobra.Command { return makeCmd(t, noConfig, "flag-account-uid") },
			want:     glossary.ListParams{AccountUID: uid.AccountUID("flag-account-uid")},
		},
		{
			testName: "account from config file",
			setup: func(t *testing.T) *cobra.Command {
				return makeCmd(t, writeConfig(t, "account_id: config-account-uid\n"), "")
			},
			want: glossary.ListParams{AccountUID: uid.AccountUID("config-account-uid")},
		},
		{
			testName: "--account flag wins over config",
			setup: func(t *testing.T) *cobra.Command {
				return makeCmd(t, writeConfig(t, "account_id: config-account-uid\n"), "flag-account-uid")
			},
			want: glossary.ListParams{AccountUID: uid.AccountUID("flag-account-uid")},
		},
		{
			testName: "name passed through",
			setup:    func(t *testing.T) *cobra.Command { return makeCmd(t, noConfig, "flag-account-uid") },
			name:     "Marketing",
			want:     glossary.ListParams{AccountUID: uid.AccountUID("flag-account-uid"), Name: "Marketing"},
		},
		{
			testName: "missing account — error",
			setup:    func(t *testing.T) *cobra.Command { return makeCmd(t, noConfig, "") },
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got, err := resolveParams(tt.setup(t), tt.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("resolveParams() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("resolveParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

Notes:
- Test the env tier with `t.Setenv`. The name is `SMARTLING_CLI_` +
  `strings.ToUpper(flagName)` with **hyphens kept**: `--account` →
  `SMARTLING_CLI_ACCOUNT`, `--file-type` → `SMARTLING_CLI_FILE-TYPE` (not
  `..._FILE_TYPE`).
- `got != tt.want` only compiles because `ListParams` has no slices. Most
  `XParams` contain slices — compare them with `reflect.DeepEqual` instead.
- Cover flag-vs-env-vs-config precedence, pass-throughs, and error branches
  (missing required value, unreadable config). Don't re-test the helpers.

---

## Service test — `services/<group>/run_<sub>_test.go`

The main unit test for a modern command's logic. Drive the real service against
the **SDK mock** (`services/<group>/sdkmocks`, `NewMock<Area>(t)`), set
expectations with `m.EXPECT().<Method>(...)`, call `RunX`, and assert the typed
output. Table-driven with `setup`/`check` funcs. Mirrors
`services/glossary/run_list_test.go`. Fake IDs only.

```go
package glossary

import (
	"errors"
	"testing"

	sdkmocks "github.com/Smartling/smartling-cli/services/glossary/sdkmocks"

	glossaryapi "github.com/Smartling/api-sdk-go/api/glossary"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

func Test_service_RunList(t *testing.T) {
	ctx := t.Context()
	const testAccountUID = uid.AccountUID("test-account-uid")

	tests := []struct {
		name    string
		setup   func(*sdkmocks.MockGlossary)
		params  ListParams
		wantErr bool
		check   func(*testing.T, ListOutput)
	}{
		{
			name:    "validation error — empty account UID",
			setup:   func(m *sdkmocks.MockGlossary) {},
			params:  ListParams{},
			wantErr: true,
		},
		{
			name: "API error",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().GetByName(ctx, testAccountUID, "").Return(nil, errors.New("API error"))
			},
			params:  ListParams{AccountUID: testAccountUID},
			wantErr: true,
		},
		{
			name: "name filter is forwarded to API",
			setup: func(m *sdkmocks.MockGlossary) {
				m.EXPECT().GetByName(ctx, testAccountUID, "Marketing").
					Return([]glossaryapi.GetGlossaryResponse{{GlossaryUID: "uid-3", Name: "Marketing"}}, nil)
			},
			params: ListParams{AccountUID: testAccountUID, Name: "Marketing"},
			check: func(t *testing.T, got ListOutput) {
				if len(got.Glossaries) != 1 || got.Glossaries[0].GlossaryUID != "uid-3" {
					t.Fatalf("unexpected result: %+v", got.Glossaries)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := sdkmocks.NewMockGlossary(t)
			tt.setup(m)
			got, err := service{glossaryApi: m}.RunList(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
```

Cover: `Validate()` failure (no SDK call), SDK error propagation, the filter
forwarded into the SDK request, and the shape of the typed output. Construct the
service directly as `service{glossaryApi: m}` (or `NewService(m)`); `t.Context()`
gives the test context.

---

## (Legacy) Command test — `cmd/<group>/<sub>/cmd_<sub>_test.go`

Only for groups that generate `Service` + `SrvInitializer` mocks
(legacy files/projects/init, hybrid mt). Two mocks: the service mock returns canned output, the
initializer mock returns the service mock. Modern commands skip this.

```go
package list // legacy files-style leaf package

import (
	"bytes"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/files/mocks"
	srvmocks "github.com/Smartling/smartling-cli/services/files/mocks"

	"github.com/stretchr/testify/mock"
)

func TestNewListCmd(t *testing.T) {
	buf := new(bytes.Buffer)

	filesSrv := srvmocks.NewMockService(t)
	filesSrv.On("RunList", mock.Anything, "any", true, "test-uri").Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitFilesSrv", mock.Anything).Return(filesSrv, nil)

	cmd := NewListCmd(initializer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"test-uri", "--short", "--format", "any"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}
}
```

---

## Integration test — `tests/cmd/<group>/<sub>/<sub>_test.go`

Black-box test of the **prebuilt binary** at `tests/cmd/bin/smartling-cli`. Hits
the real API — the developer builds the binary, drops it plus a `smartling.yml`
into `tests/cmd/bin/`, and runs the suite manually. **Don't run it yourself.**
Author the file and register its dir in the Makefile `test_integration` target.

```go
package export

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// createGlossary makes this package's own fixture; there is no delete
// command, so a per-run name keeps reruns from colliding.
func createGlossary(t *testing.T, absDir string) string {
	t.Helper()
	name := "test-integration-glossary-export-" + strconv.FormatInt(time.Now().Unix(), 10)
	cmd := exec.Command("./smartling-cli", "glossaries", "create", name, "--locale", "en-US")
	cmd.Dir = absDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("create fixture glossary: %v\n%s", err, out)
	}
	return name
}

func TestGlossaryExport(t *testing.T) {
	absDir, err := filepath.Abs("../../bin/") // "../../../bin/" for tests/cmd/<group>/<sub>/<verb>/
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	subCommands := []string{"glossaries", "export"}
	testGlossary := createGlossary(t, absDir)

	tests := []struct {
		name              string
		args              []string
		expectedOutputs   []string
		unexpectedOutputs []string
		wantErr           bool
		cleanup           func()
	}{
		{
			name:              "help flag shows command description",
			args:              append(subCommands, "--help"),
			expectedOutputs:   []string{"Export a glossary's entries to a file", "file-type"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:            "missing file type is rejected",
			args:            append(subCommands, testGlossary),
			expectedOutputs: []string{"file type is required"},
			wantErr:         true,
		},
		{
			name:              "table output format shows column headers",
			args:              []string{"glossaries", "--output", "table", "export", testGlossary, "--file-type", "csv"},
			expectedOutputs:   []string{"GLOSSARY UID", "OUTPUT FILE", "FILE TYPE"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name: "export to explicit output file",
			args: append(subCommands, testGlossary, "out.csv", "--file-type", "csv"),
			expectedOutputs: []string{"Output file:", "out.csv", "Bytes written:"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			cleanup: func() { _ = os.Remove(filepath.Join(absDir, "out.csv")) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCmd := exec.Command("./smartling-cli", tt.args...)
			testCmd.Dir = absDir
			out, err := testCmd.CombinedOutput()
			if tt.cleanup != nil {
				defer tt.cleanup()
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("error: %v, output: %s", err, string(out))
			}
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got none; output: %s", string(out))
			}
			for _, want := range tt.expectedOutputs {
				if !strings.Contains(string(out), want) {
					t.Errorf("output: %s\nwithout expected: %s", string(out), want)
				}
			}
			for _, notWant := range tt.unexpectedOutputs {
				if strings.Contains(string(out), notWant) {
					t.Errorf("output: %s\nwith unexpected: %s", string(out), notWant)
				}
			}
		})
	}
}
```

---

## Full checklist for one new subcommand

```
[ ] (if needed) add SDK method to api/<area> INTERFACE in api-sdk-go;
                 use client *smclient.Client + smerror errors
[ ] services/<group>/service.go     — add RunX to Service interface
[ ] services/<group>/run_x.go       — XParams, XOutput (Renderable), RunX impl
[ ] cmd/<group>/x/cmd_x.go          — NewXCmd(initializer), flags, RunE
[ ] cmd/<group>/x/resolve.go        — resolveParams
[ ] cmd/<group>/x/run.go            — run (init srv, call RunX, render)
[ ] main.go                         — gCmd.AddCommand(x.NewXCmd(gInit))
[ ] .mockery.yml                    — add SDK (+Service/SrvInitializer if used);
                                      then make mockery
[ ] cmd/<group>/x/resolve_test.go   — resolveParams test: root tree + temp config (fake IDs)
[ ] services/<group>/run_x_test.go  — service test vs SDK mock (fake IDs)
[ ] (legacy/mt groups only) cmd/<group>/x/cmd_x_test.go — command test w/ Service+Init mocks
[ ] tests/cmd/<group>/x/x_test.go   — integration test (author only; don't run)
[ ] Makefile test_integration       — add: go test ./tests/cmd/<group>/x/...
[ ] make mockery && make test_unit && make lint
    # NOT make test_integration, NOT make docs — the developer runs those
```