package glcreate

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
	createCmd := NewCreateCmd(nil)
	root.AddCommand(createCmd)
	_ = root.PersistentFlags().Set("config", configPath)
	t.Cleanup(func() { _ = root.PersistentFlags().Set("config", "") })
	if accountFlag != "" {
		_ = root.PersistentFlags().Set("account", accountFlag)
	}
	return createCmd
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "smartling.yml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func newFileConfig(create glossarycmd.CreateConfig) glossarycmd.FileConfig {
	var fc glossarycmd.FileConfig
	fc.Glossary.Create = create
	return fc
}

func Test_resolveParams(t *testing.T) {
	t.Setenv("SMARTLING_USER_ID", "test-user")
	t.Setenv("SMARTLING_SECRET", "test-secret")

	noConfigPath := filepath.Join(t.TempDir(), "no-smartling.yml")

	const (
		testAccount      = uid.AccountUID("test-account-uid")
		testGlossaryName = "My Glossary"
	)

	// parseFallbackLocales always returns a non-nil slice, so the zero-value
	// expected params must use an empty slice literal, not nil.
	baseWant := func() srv.CreateParams {
		return srv.CreateParams{
			AccountUID:      testAccount,
			GlossaryName:    testGlossaryName,
			FallbackLocales: []srv.FallbackLocale{},
		}
	}

	defaultSetup := func(t *testing.T) *cobra.Command {
		return makeCmd(t, noConfigPath, string(testAccount))
	}

	tests := []struct {
		name         string
		setup        func(t *testing.T) *cobra.Command
		fileConfig   glossarycmd.FileConfig
		glossaryName string
		want         srv.CreateParams
		wantErr      bool
	}{
		// ── Account resolution ────────────────────────────────────────────────
		{
			name:         "account from --account flag",
			setup:        defaultSetup,
			glossaryName: testGlossaryName,
			want:         baseWant(),
		},
		{
			name: "account from SMARTLING_CLI_ACCOUNT env var",
			setup: func(t *testing.T) *cobra.Command {
				t.Setenv("SMARTLING_CLI_ACCOUNT", string(testAccount))
				return makeCmd(t, noConfigPath, "")
			},
			glossaryName: testGlossaryName,
			want:         baseWant(),
		},
		{
			name: "account from config file",
			setup: func(t *testing.T) *cobra.Command {
				cfgPath := writeConfig(t, "account_id: "+string(testAccount)+"\n")
				return makeCmd(t, cfgPath, "")
			},
			glossaryName: testGlossaryName,
			want:         baseWant(),
		},
		// ── glossaryName pass-through ─────────────────────────────────────────
		{
			name:         "glossaryName is passed through",
			setup:        defaultSetup,
			glossaryName: "Custom Glossary Name",
			want: func() srv.CreateParams {
				p := baseWant()
				p.GlossaryName = "Custom Glossary Name"
				return p
			}(),
		},
		// ── description ───────────────────────────────────────────────────────
		{
			name: "description from --description flag",
			setup: func(t *testing.T) *cobra.Command {
				cmd := defaultSetup(t)
				_ = cmd.Flags().Set(descriptionFlag, "A helpful glossary")
				return cmd
			},
			glossaryName: testGlossaryName,
			want: func() srv.CreateParams {
				p := baseWant()
				p.Description = "A helpful glossary"
				return p
			}(),
		},
		// ── verificationMode ──────────────────────────────────────────────────
		{
			name: "verificationMode from --verification-mode flag",
			setup: func(t *testing.T) *cobra.Command {
				cmd := defaultSetup(t)
				_ = cmd.Flags().Set(verificationModeFlag, "true")
				return cmd
			},
			glossaryName: testGlossaryName,
			want: func() srv.CreateParams {
				p := baseWant()
				p.VerificationMode = true
				return p
			}(),
		},
		{
			name:         "verificationMode from fileConfig",
			setup:        defaultSetup,
			fileConfig:   newFileConfig(glossarycmd.CreateConfig{VerificationMode: true}),
			glossaryName: testGlossaryName,
			want: func() srv.CreateParams {
				p := baseWant()
				p.VerificationMode = true
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
			glossaryName: testGlossaryName,
			want: func() srv.CreateParams {
				p := baseWant()
				p.LocaleIDs = []string{"en-US", "es-ES"}
				return p
			}(),
		},
		{
			name:         "localeIDs from fileConfig",
			setup:        defaultSetup,
			fileConfig:   newFileConfig(glossarycmd.CreateConfig{LocaleIDs: []string{"fr-FR", "de-DE"}}),
			glossaryName: testGlossaryName,
			want: func() srv.CreateParams {
				p := baseWant()
				p.LocaleIDs = []string{"fr-FR", "de-DE"}
				return p
			}(),
		},
		// ── fallbackLocales ───────────────────────────────────────────────────
		{
			name: "fallbackLocale from --fallback-locale flag — single locale",
			setup: func(t *testing.T) *cobra.Command {
				cmd := defaultSetup(t)
				_ = cmd.Flags().Set(fallbackLocaleFlag, "es:es-ES")
				return cmd
			},
			glossaryName: testGlossaryName,
			want: func() srv.CreateParams {
				p := baseWant()
				p.FallbackLocales = []srv.FallbackLocale{
					{FallbackLocaleID: "es", LocaleIDs: []string{"es-ES"}},
				}
				return p
			}(),
		},
		{
			name: "fallbackLocale from --fallback-locale flag — multiple locales",
			setup: func(t *testing.T) *cobra.Command {
				cmd := defaultSetup(t)
				_ = cmd.Flags().Set(fallbackLocaleFlag, "es:es-MX,es-AR")
				return cmd
			},
			glossaryName: testGlossaryName,
			want: func() srv.CreateParams {
				p := baseWant()
				p.FallbackLocales = []srv.FallbackLocale{
					{FallbackLocaleID: "es", LocaleIDs: []string{"es-MX", "es-AR"}},
				}
				return p
			}(),
		},
		{
			name: "multiple --fallback-locale flags",
			setup: func(t *testing.T) *cobra.Command {
				cmd := defaultSetup(t)
				_ = cmd.Flags().Set(fallbackLocaleFlag, "es:es-ES")
				_ = cmd.Flags().Set(fallbackLocaleFlag, "pt:pt-BR,pt-PT")
				return cmd
			},
			glossaryName: testGlossaryName,
			want: func() srv.CreateParams {
				p := baseWant()
				p.FallbackLocales = []srv.FallbackLocale{
					{FallbackLocaleID: "es", LocaleIDs: []string{"es-ES"}},
					{FallbackLocaleID: "pt", LocaleIDs: []string{"pt-BR", "pt-PT"}},
				}
				return p
			}(),
		},
		{
			name:  "fallbackLocales from fileConfig",
			setup: defaultSetup,
			fileConfig: newFileConfig(glossarycmd.CreateConfig{
				FallbackLocales: []string{"es:es-ES", "pt:pt-BR"},
			}),
			glossaryName: testGlossaryName,
			want: func() srv.CreateParams {
				p := baseWant()
				p.FallbackLocales = []srv.FallbackLocale{
					{FallbackLocaleID: "es", LocaleIDs: []string{"es-ES"}},
					{FallbackLocaleID: "pt", LocaleIDs: []string{"pt-BR"}},
				}
				return p
			}(),
		},
		// ── Error cases ───────────────────────────────────────────────────────
		{
			name: "invalid fallback locale format — error",
			setup: func(t *testing.T) *cobra.Command {
				cmd := defaultSetup(t)
				_ = cmd.Flags().Set(fallbackLocaleFlag, "no-colon-here")
				return cmd
			},
			glossaryName: testGlossaryName,
			wantErr:      true,
		},
		{
			name:         "missing account — error",
			setup:        func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, "") },
			glossaryName: testGlossaryName,
			wantErr:      true,
		},
		{
			name:         "Config load failure — unreadable config path",
			setup:        func(t *testing.T) *cobra.Command { return makeCmd(t, t.TempDir(), "") },
			glossaryName: testGlossaryName,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setup(t)
			got, err := resolveParams(cmd, tt.fileConfig, tt.glossaryName)
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

func Test_parseFallbackLocales(t *testing.T) {
	tests := []struct {
		name    string
		raws    []string
		want    []srv.FallbackLocale
		wantErr bool
	}{
		{
			name: "nil input returns empty slice",
			raws: nil,
			want: []srv.FallbackLocale{},
		},
		{
			name: "single entry — one locale",
			raws: []string{"es:es-ES"},
			want: []srv.FallbackLocale{
				{FallbackLocaleID: "es", LocaleIDs: []string{"es-ES"}},
			},
		},
		{
			name: "single entry — multiple locales",
			raws: []string{"es:es-MX,es-AR,es-ES"},
			want: []srv.FallbackLocale{
				{FallbackLocaleID: "es", LocaleIDs: []string{"es-MX", "es-AR", "es-ES"}},
			},
		},
		{
			name: "multiple entries",
			raws: []string{"es:es-ES", "pt:pt-BR,pt-PT"},
			want: []srv.FallbackLocale{
				{FallbackLocaleID: "es", LocaleIDs: []string{"es-ES"}},
				{FallbackLocaleID: "pt", LocaleIDs: []string{"pt-BR", "pt-PT"}},
			},
		},
		{
			name:    "no colon — error",
			raws:    []string{"no-colon"},
			wantErr: true,
		},
		{
			name:    "empty fallback locale ID — error",
			raws:    []string{":es-ES"},
			wantErr: true,
		},
		{
			name:    "empty locale IDs after colon — error",
			raws:    []string{"es:"},
			wantErr: true,
		},
		{
			name:    "empty string — error",
			raws:    []string{""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFallbackLocales(tt.raws)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFallbackLocales() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFallbackLocales() got = %v, want %v", got, tt.want)
			}
		})
	}
}
