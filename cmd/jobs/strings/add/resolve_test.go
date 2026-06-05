package add

import (
	"os"
	"path/filepath"
	"testing"

	rootcmd "github.com/Smartling/smartling-cli/cmd"
)

func TestMain(m *testing.M) {
	rootcmd.ConfigureLogger()
	os.Exit(m.Run())
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

	root := rootcmd.NewRootCmd()
	cfgPath := writeConfig(t, "project_id: config-project-id\n")
	if err := root.PersistentFlags().Set("config", cfgPath); err != nil {
		t.Fatalf("set config flag: %v", err)
	}
	t.Cleanup(func() { _ = root.PersistentFlags().Set("config", "") })

	params, err := resolveParams("aabbccdd1122", []string{"h1", "h2"}, []string{"fr-FR"}, true)
	if err != nil {
		t.Fatalf("resolveParams() error = %v", err)
	}
	if params.ProjectID != "config-project-id" {
		t.Errorf("ProjectID = %q, want config-project-id", params.ProjectID)
	}
	if params.JobUIDOrName != "aabbccdd1122" {
		t.Errorf("JobUIDOrName = %q, want aabbccdd1122", params.JobUIDOrName)
	}
	if len(params.Hashcodes) != 2 || params.Hashcodes[0] != "h1" {
		t.Errorf("Hashcodes = %v, want [h1 h2]", params.Hashcodes)
	}
	if len(params.TargetLocaleIDs) != 1 || params.TargetLocaleIDs[0] != "fr-FR" {
		t.Errorf("TargetLocaleIDs = %v, want [fr-FR]", params.TargetLocaleIDs)
	}
	if !params.MoveEnabled {
		t.Error("MoveEnabled = false, want true")
	}
}
