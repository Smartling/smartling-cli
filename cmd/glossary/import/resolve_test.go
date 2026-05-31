package glimport

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	glossarycmd "github.com/Smartling/smartling-cli/cmd/glossary"
	srv "github.com/Smartling/smartling-cli/services/glossary"

	"github.com/Smartling/api-sdk-go/helpers/uid"
	"github.com/spf13/cobra"
)

func TestMain(m *testing.M) {
	rootcmd.ConfigureLogger()
	m.Run()
}

func makeCmd(t *testing.T, configPath, accountFlag string) *cobra.Command {
	t.Helper()
	root := rootcmd.NewRootCmd()
	importCmd := NewImportCmd(nil)
	root.AddCommand(importCmd)
	_ = root.PersistentFlags().Set("config", configPath)
	t.Cleanup(func() { _ = root.PersistentFlags().Set("config", "") })
	if accountFlag != "" {
		_ = root.PersistentFlags().Set("account", accountFlag)
	}
	return importCmd
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
	t.Setenv("SMARTLING_USER_ID", "test-user")
	t.Setenv("SMARTLING_SECRET", "test-secret")

	noConfigPath := filepath.Join(t.TempDir(), "no-smartling.yml")

	const (
		testAccount  = uid.AccountUID("test-account-uid")
		testGlossary = "my-glossary"
	)

	// baseWant returns an ImportParams with the common fields filled for the
	// given inFile; callers override only the fields they are testing.
	baseWant := func(inFile, mediaType string) srv.ImportParams {
		return srv.ImportParams{
			AccountUID:        testAccount,
			GlossaryUIDOrName: testGlossary,
			ImportFile: srv.ImportFile{
				Path:      inFile,
				Name:      filepath.Base(inFile),
				MediaType: mediaType,
			},
		}
	}

	tests := []struct {
		name              string
		setup             func(t *testing.T) *cobra.Command
		fileConfig        glossarycmd.FileConfig
		glossaryUIDOrName string
		inFile            string
		want              srv.ImportParams
		wantErr           bool
	}{
		// ── Account resolution ────────────────────────────────────────────────
		{
			name:              "account from --account flag",
			setup:             func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			glossaryUIDOrName: testGlossary,
			inFile:            "terms.csv",
			want:              baseWant("terms.csv", "text/csv"),
		},
		{
			name: "account from SMARTLING_CLI_ACCOUNT env var",
			setup: func(t *testing.T) *cobra.Command {
				t.Setenv("SMARTLING_CLI_ACCOUNT", string(testAccount))
				return makeCmd(t, noConfigPath, "")
			},
			glossaryUIDOrName: testGlossary,
			inFile:            "terms.csv",
			want:              baseWant("terms.csv", "text/csv"),
		},
		{
			name: "account from config file",
			setup: func(t *testing.T) *cobra.Command {
				cfgPath := writeConfig(t, "account_id: "+string(testAccount)+"\n")
				return makeCmd(t, cfgPath, "")
			},
			glossaryUIDOrName: testGlossary,
			inFile:            "terms.csv",
			want:              baseWant("terms.csv", "text/csv"),
		},
		// ── glossaryUIDOrName pass-through ────────────────────────────────────
		{
			name:              "glossaryUIDOrName is passed through",
			setup:             func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			glossaryUIDOrName: "00000000-0000-0000-0000-000000000001",
			inFile:            "terms.csv",
			want: func() srv.ImportParams {
				p := baseWant("terms.csv", "text/csv")
				p.GlossaryUIDOrName = "00000000-0000-0000-0000-000000000001"
				return p
			}(),
		},
		// ── inFile → ImportFile mapping ───────────────────────────────────────
		{
			name:              "inFile sets Path and Name from argument",
			setup:             func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			glossaryUIDOrName: testGlossary,
			inFile:            "/exports/2024/terms.csv",
			want:              baseWant("/exports/2024/terms.csv", "text/csv"),
		},
		// ── mediaType derived from extension ──────────────────────────────────
		{
			name:              "mediaType derived from .csv extension",
			setup:             func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			glossaryUIDOrName: testGlossary,
			inFile:            "terms.csv",
			want:              baseWant("terms.csv", "text/csv"),
		},
		{
			name:              "mediaType derived from .xlsx extension",
			setup:             func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			glossaryUIDOrName: testGlossary,
			inFile:            "terms.xlsx",
			want:              baseWant("terms.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"),
		},
		{
			name:              "mediaType derived from .tbx extension",
			setup:             func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			glossaryUIDOrName: testGlossary,
			inFile:            "terms.tbx",
			want:              baseWant("terms.tbx", "text/xml"),
		},
		{
			name:              "mediaType derived from .xml extension",
			setup:             func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			glossaryUIDOrName: testGlossary,
			inFile:            "terms.xml",
			want:              baseWant("terms.xml", "text/xml"),
		},
		{
			name:              "unknown extension leaves mediaType empty",
			setup:             func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			glossaryUIDOrName: testGlossary,
			inFile:            "terms.dat",
			want:              baseWant("terms.dat", ""),
		},
		// ── mediaType flag and fileConfig ─────────────────────────────────────
		{
			name: "--media-type flag overrides derived mediaType",
			setup: func(t *testing.T) *cobra.Command {
				cmd := makeCmd(t, noConfigPath, string(testAccount))
				_ = cmd.Flags().Set(mediaTypeFlag, "text/csv")
				return cmd
			},
			glossaryUIDOrName: testGlossary,
			inFile:            "terms.dat",
			want:              baseWant("terms.dat", "text/csv"),
		},
		{
			name:  "mediaType from fileConfig when no flag set",
			setup: func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			fileConfig: glossarycmd.FileConfig{Glossary: struct {
				Export glossarycmd.ExportConfig `yaml:"export,omitzero"`
				Create glossarycmd.CreateConfig `yaml:"create,omitzero"`
				Import glossarycmd.ImportConfig `yaml:"import,omitzero"`
			}{Import: glossarycmd.ImportConfig{MediaType: "text/xml"}}},
			glossaryUIDOrName: testGlossary,
			inFile:            "terms.dat",
			want:              baseWant("terms.dat", "text/xml"),
		},
		// ── archiveMode ───────────────────────────────────────────────────────
		{
			name: "--archive-mode flag sets archiveMode",
			setup: func(t *testing.T) *cobra.Command {
				cmd := makeCmd(t, noConfigPath, string(testAccount))
				_ = cmd.Flags().Set(archiveModeFlag, "true")
				return cmd
			},
			glossaryUIDOrName: testGlossary,
			inFile:            "terms.csv",
			want: func() srv.ImportParams {
				p := baseWant("terms.csv", "text/csv")
				p.ArchiveMode = true
				return p
			}(),
		},
		{
			name:  "archiveMode from fileConfig",
			setup: func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			fileConfig: glossarycmd.FileConfig{Glossary: struct {
				Export glossarycmd.ExportConfig `yaml:"export,omitzero"`
				Create glossarycmd.CreateConfig `yaml:"create,omitzero"`
				Import glossarycmd.ImportConfig `yaml:"import,omitzero"`
			}{Import: glossarycmd.ImportConfig{ArchiveMode: true}}},
			glossaryUIDOrName: testGlossary,
			inFile:            "terms.csv",
			want: func() srv.ImportParams {
				p := baseWant("terms.csv", "text/csv")
				p.ArchiveMode = true
				return p
			}(),
		},
		// ── error cases ───────────────────────────────────────────────────────
		{
			name:    "missing account — error",
			setup:   func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, "") },
			inFile:  "terms.csv",
			wantErr: true,
		},
		{
			name: "Config load failure — unreadable config path",
			setup: func(t *testing.T) *cobra.Command {
				return makeCmd(t, t.TempDir(), "")
			},
			inFile:  "terms.csv",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setup(t)
			got, err := resolveParams(cmd, tt.fileConfig, tt.glossaryUIDOrName, tt.inFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}
