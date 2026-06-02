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

// makeCmd builds a minimal root→child command tree and binds --config and
// optionally --account on the root so that rootcmd.Config() can be controlled
// from tests. Cobra's StringVarP binding propagates the flag value to the
// package-level global in the cmd package; t.Cleanup resets it afterward.
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
	// SMARTLING_USER_ID and SMARTLING_SECRET are always required by Config().
	t.Setenv("SMARTLING_USER_ID", "test-user")
	t.Setenv("SMARTLING_SECRET", "test-secret")

	// noConfigPath: file doesn't exist — LoadConfigFromFile tolerates this and
	// returns an empty Config without error.
	noConfigPath := filepath.Join(t.TempDir(), "no-smartling.yml")

	tests := []struct {
		testName string
		setup    func(t *testing.T) *cobra.Command
		name     string
		want     glossary.ListParams
		wantErr  bool
	}{
		{
			testName: "account from --account flag",
			setup:    func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, "flag-account-uid") },
			want:     glossary.ListParams{AccountUID: uid.AccountUID("flag-account-uid")},
		},
		{
			testName: "account from SMARTLING_CLI_ACCOUNT env var",
			setup: func(t *testing.T) *cobra.Command {
				t.Setenv("SMARTLING_CLI_ACCOUNT", "env-account-uid")
				return makeCmd(t, noConfigPath, "")
			},
			want: glossary.ListParams{AccountUID: uid.AccountUID("env-account-uid")},
		},
		{
			testName: "account from config file",
			setup: func(t *testing.T) *cobra.Command {
				cfgPath := writeConfig(t, "account_id: config-account-uid\n")
				return makeCmd(t, cfgPath, "")
			},
			want: glossary.ListParams{AccountUID: uid.AccountUID("config-account-uid")},
		},
		{
			testName: "--account flag wins over config file",
			setup: func(t *testing.T) *cobra.Command {
				cfgPath := writeConfig(t, "account_id: config-account-uid\n")
				return makeCmd(t, cfgPath, "flag-account-uid")
			},
			want: glossary.ListParams{AccountUID: uid.AccountUID("flag-account-uid")},
		},
		{
			testName: "name filter is passed through",
			setup:    func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, "flag-account-uid") },
			name:     "Marketing",
			want:     glossary.ListParams{AccountUID: uid.AccountUID("flag-account-uid"), Name: "Marketing"},
		},
		{
			testName: "missing account — error",
			setup:    func(t *testing.T) *cobra.Command { return makeCmd(t, noConfigPath, "") },
			wantErr:  true,
		},
		{
			testName: "Config load failure — unreadable config path",
			setup: func(t *testing.T) *cobra.Command {
				// t.TempDir() is a directory, not a file; os.ReadFile on a
				// directory returns an error that is not os.IsNotExist, so
				// LoadConfigFromFile propagates it through Config().
				return makeCmd(t, t.TempDir(), "")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cmd := tt.setup(t)
			got, err := resolveParams(cmd, tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("resolveParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}
