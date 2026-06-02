package glexport

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
	glossariescmd "github.com/Smartling/smartling-cli/cmd/glossaries"
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
	exportCmd := NewExportCmd(nil)
	root.AddCommand(exportCmd)
	_ = root.PersistentFlags().Set("config", configPath)
	t.Cleanup(func() { _ = root.PersistentFlags().Set("config", "") })
	if accountFlag != "" {
		_ = root.PersistentFlags().Set("account", accountFlag)
	}
	return exportCmd
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "smartling.yml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

// newFileConfig returns a FileConfig with only the Export section populated.
func newFileConfig(export glossariescmd.ExportConfig) glossariescmd.FileConfig {
	var fc glossariescmd.FileConfig
	fc.Glossaries.Export = export
	return fc
}

func Test_resolveParams(t *testing.T) {
	t.Setenv("SMARTLING_USER_ID", "test-user")
	t.Setenv("SMARTLING_SECRET", "test-secret")

	noConfigPath := filepath.Join(t.TempDir(), "no-smartling.yml")

	const (
		testAccount  = uid.AccountUID("test-account-uid")
		testGlossary = "my-glossary"
		testOutFile  = "out.csv"
	)

	// baseWant is a zeroed ExportParams with the invariant fields set.
	// Individual test cases override only the fields under test.
	baseWant := func() srv.ExportParams {
		return srv.ExportParams{
			AccountUID:        testAccount,
			GlossaryUIDOrName: testGlossary,
			OutFile:           testOutFile,
			FileType:          "csv",
		}
	}

	defaultSetup := func(t *testing.T) *cobra.Command {
		cmd := makeCmd(t, noConfigPath, string(testAccount))
		_ = cmd.Flags().Set(fileTypeFlag, "csv")
		return cmd
	}

	tests := []struct {
		name              string
		setup             func(t *testing.T) *cobra.Command
		fileConfig        glossariescmd.FileConfig
		glossaryUIDOrName string
		outFile           string
		want              srv.ExportParams
		wantErr           bool
	}{
		// ── Account resolution ────────────────────────────────────────────────
		{
			name:              "account from --account flag",
			setup:             defaultSetup,
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want:              baseWant(),
		},
		{
			name: "account from SMARTLING_CLI_ACCOUNT env var",
			setup: func(t *testing.T) *cobra.Command {
				t.Setenv("SMARTLING_CLI_ACCOUNT", string(testAccount))
				cmd := makeCmd(t, noConfigPath, "")
				_ = cmd.Flags().Set(fileTypeFlag, "csv")
				return cmd
			},
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want:              baseWant(),
		},
		{
			name: "account from config file",
			setup: func(t *testing.T) *cobra.Command {
				cfgPath := writeConfig(t, "account_id: "+string(testAccount)+"\n")
				cmd := makeCmd(t, cfgPath, "")
				_ = cmd.Flags().Set(fileTypeFlag, "csv")
				return cmd
			},
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want:              baseWant(),
		},
		// ── Pass-through fields ───────────────────────────────────────────────
		{
			name:              "glossaryUIDOrName and outFile are passed through",
			setup:             defaultSetup,
			glossaryUIDOrName: "00000000-0000-0000-0000-000000000001",
			outFile:           "/tmp/export.xlsx",
			want: func() srv.ExportParams {
				p := baseWant()
				p.GlossaryUIDOrName = "00000000-0000-0000-0000-000000000001"
				p.OutFile = "/tmp/export.xlsx"
				return p
			}(),
		},
		// ── fileType ──────────────────────────────────────────────────────────
		{
			name: "fileType is lowercased",
			setup: func(t *testing.T) *cobra.Command {
				cmd := makeCmd(t, noConfigPath, string(testAccount))
				_ = cmd.Flags().Set(fileTypeFlag, "CSV")
				return cmd
			},
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want:              baseWant(), // "csv" after ToLower
		},
		{
			name: "fileType from fileConfig when no flag",
			setup: func(t *testing.T) *cobra.Command {
				return makeCmd(t, noConfigPath, string(testAccount))
			},
			fileConfig:        newFileConfig(glossariescmd.ExportConfig{FileType: "xlsx"}),
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want: func() srv.ExportParams {
				p := baseWant()
				p.FileType = "xlsx"
				return p
			}(),
		},
		// ── tbxVersion ────────────────────────────────────────────────────────
		{
			name: "tbxVersion is lowercased",
			setup: func(t *testing.T) *cobra.Command {
				cmd := makeCmd(t, noConfigPath, string(testAccount))
				_ = cmd.Flags().Set(fileTypeFlag, "tbx")
				_ = cmd.Flags().Set(tbxVersionFlag, "V3")
				return cmd
			},
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want: func() srv.ExportParams {
				p := baseWant()
				p.FileType = "tbx"
				p.TbxVersion = "v3"
				return p
			}(),
		},
		{
			name: "tbxVersion from fileConfig when no flag",
			setup: func(t *testing.T) *cobra.Command {
				cmd := makeCmd(t, noConfigPath, string(testAccount))
				_ = cmd.Flags().Set(fileTypeFlag, "tbx")
				return cmd
			},
			fileConfig:        newFileConfig(glossariescmd.ExportConfig{TbxVersion: "v2"}),
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want: func() srv.ExportParams {
				p := baseWant()
				p.FileType = "tbx"
				p.TbxVersion = "v2"
				return p
			}(),
		},
		// ── localeIDs ─────────────────────────────────────────────────────────
		{
			name: "localeIDs from --locale flags",
			setup: func(t *testing.T) *cobra.Command {
				cmd := defaultSetup(t)
				_ = cmd.Flags().Set(localeFlag, "en-US")
				_ = cmd.Flags().Set(localeFlag, "es-ES")
				return cmd
			},
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want: func() srv.ExportParams {
				p := baseWant()
				p.LocaleIDs = []string{"en-US", "es-ES"}
				return p
			}(),
		},
		{
			name:              "localeIDs from fileConfig",
			setup:             func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			fileConfig:        newFileConfig(glossariescmd.ExportConfig{FileType: "csv", LocaleIDs: []string{"fr-FR"}}),
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want: func() srv.ExportParams {
				p := baseWant()
				p.LocaleIDs = []string{"fr-FR"}
				return p
			}(),
		},
		// ── skipEntries ───────────────────────────────────────────────────────
		{
			name: "skipEntries from --skip-entries flag",
			setup: func(t *testing.T) *cobra.Command {
				cmd := defaultSetup(t)
				_ = cmd.Flags().Set(skipEntriesFlag, "true")
				return cmd
			},
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want: func() srv.ExportParams {
				p := baseWant()
				p.SkipEntries = true
				return p
			}(),
		},
		{
			name:              "skipEntries from fileConfig",
			setup:             func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			fileConfig:        newFileConfig(glossariescmd.ExportConfig{FileType: "csv", SkipEntries: true}),
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want: func() srv.ExportParams {
				p := baseWant()
				p.SkipEntries = true
				return p
			}(),
		},
		// ── Filter fields from flags ──────────────────────────────────────────
		{
			name: "filter fields from flags",
			setup: func(t *testing.T) *cobra.Command {
				cmd := defaultSetup(t)
				_ = cmd.Flags().Set(filterQueryFlag, "checkout")
				_ = cmd.Flags().Set(filterLocaleFlag, "en-US")
				_ = cmd.Flags().Set(filterEntryStateFlag, "ACTIVE")
				_ = cmd.Flags().Set(filterReturnFallbackTranslationsFlag, "true")
				_ = cmd.Flags().Set(filterCreatedLevelFlag, "ACCOUNT")
				_ = cmd.Flags().Set(filterCreatedByUserIDFlag, "user-1")
				return cmd
			},
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want: func() srv.ExportParams {
				p := baseWant()
				p.Filter.Query = "checkout"
				p.Filter.LocaleID = []string{"en-US"}
				p.Filter.EntryState = "ACTIVE"
				p.Filter.ReturnFallbackTranslations = true
				p.Filter.Created.Level = "ACCOUNT"
				p.Filter.CreatedBy.UserIDs = []string{"user-1"}
				return p
			}(),
		},
		// ── Filter fields from fileConfig ─────────────────────────────────────
		{
			name:  "filter fields from fileConfig",
			setup: func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, string(testAccount)) },
			fileConfig: newFileConfig(glossariescmd.ExportConfig{
				FileType: "csv",
				Filter: glossariescmd.ExportFilterConfig{
					Query:      "checkout",
					EntryState: "ACTIVE",
					EntryUIDs:  []string{"uid-1", "uid-2"},
					LastModified: glossariescmd.LastModifiedConfig{
						Level: "GLOSSARY",
						Type:  "AFTER",
					},
				},
			}),
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want: func() srv.ExportParams {
				p := baseWant()
				p.Filter.Query = "checkout"
				p.Filter.EntryState = "ACTIVE"
				p.Filter.EntryUIDs = []string{"uid-1", "uid-2"}
				p.Filter.LastModified.Level = "GLOSSARY"
				p.Filter.LastModified.Type = "AFTER"
				return p
			}(),
		},
		// ── Date filters ──────────────────────────────────────────────────────
		{
			name: "filter-created-date parsed from RFC3339 flag",
			setup: func(t *testing.T) *cobra.Command {
				cmd := defaultSetup(t)
				_ = cmd.Flags().Set(createdDateFlag, "2026-01-02T15:04:05Z")
				return cmd
			},
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want: func() srv.ExportParams {
				p := baseWant()
				p.Filter.Created.Date = time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)
				return p
			}(),
		},
		{
			name: "filter-last-modified-date parsed from RFC3339 flag",
			setup: func(t *testing.T) *cobra.Command {
				cmd := defaultSetup(t)
				_ = cmd.Flags().Set(lastModifiedDateFlag, "2026-06-15T00:00:00Z")
				return cmd
			},
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			want: func() srv.ExportParams {
				p := baseWant()
				p.Filter.LastModified.Date = time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
				return p
			}(),
		},
		// ── Error cases ───────────────────────────────────────────────────────
		{
			name: "invalid RFC3339 date — error",
			setup: func(t *testing.T) *cobra.Command {
				cmd := defaultSetup(t)
				_ = cmd.Flags().Set(createdDateFlag, "not-a-date")
				return cmd
			},
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			wantErr:           true,
		},
		{
			name:              "missing account — error",
			setup:             func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, "") },
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			wantErr:           true,
		},
		{
			name: "Config load failure — unreadable config path",
			setup: func(t *testing.T) *cobra.Command {
				return makeCmd(t, t.TempDir(), "")
			},
			glossaryUIDOrName: testGlossary,
			outFile:           testOutFile,
			wantErr:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setup(t)
			got, err := resolveParams(cmd, tt.fileConfig, tt.glossaryUIDOrName, tt.outFile)
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
