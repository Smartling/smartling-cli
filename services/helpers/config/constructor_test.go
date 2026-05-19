package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Smartling/smartling-cli/services/helpers/rlog"
)

func TestMain(m *testing.M) {
	rlog.Init()
	os.Exit(m.Run())
}

func writeTempConfig(t *testing.T, contents string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "smartling.yml")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return dir
}

func TestBuildConfigFromFlags_SourcesAllDefault(t *testing.T) {
	t.Setenv("SMARTLING_USER_ID", "")
	t.Setenv("SMARTLING_SECRET", "")
	t.Setenv("SMARTLING_PROJECT_ID", "")

	dir := writeTempConfig(t, "")

	cfg, err := BuildConfigFromFlags(Params{
		Directory: dir,
		IsInit:    true,
	})
	if err != nil {
		t.Fatalf("BuildConfigFromFlags: %v", err)
	}

	if cfg.Sources.UserID != SourceDefault {
		t.Errorf("UserID source = %q, want %q", cfg.Sources.UserID, SourceDefault)
	}
	if cfg.Sources.AccountID != SourceDefault {
		t.Errorf("AccountID source = %q, want %q", cfg.Sources.AccountID, SourceDefault)
	}
	if cfg.Sources.ProjectID != SourceDefault {
		t.Errorf("ProjectID source = %q, want %q", cfg.Sources.ProjectID, SourceDefault)
	}
}

func TestBuildConfigFromFlags_SourcesFromFile(t *testing.T) {
	t.Setenv("SMARTLING_USER_ID", "")
	t.Setenv("SMARTLING_SECRET", "")
	t.Setenv("SMARTLING_PROJECT_ID", "")

	dir := writeTempConfig(t, `
user_id: "file-user"
secret: "file-secret"
account_id: "file-account"
project_id: "file-project"
`)

	cfg, err := BuildConfigFromFlags(Params{Directory: dir})
	if err != nil {
		t.Fatalf("BuildConfigFromFlags: %v", err)
	}

	if cfg.Sources.UserID != SourceConfig {
		t.Errorf("UserID source = %q, want %q", cfg.Sources.UserID, SourceConfig)
	}
	if cfg.Sources.AccountID != SourceConfig {
		t.Errorf("AccountID source = %q, want %q", cfg.Sources.AccountID, SourceConfig)
	}
	if cfg.Sources.ProjectID != SourceConfig {
		t.Errorf("ProjectID source = %q, want %q", cfg.Sources.ProjectID, SourceConfig)
	}
}

func TestBuildConfigFromFlags_SourcesFromEnv(t *testing.T) {
	t.Setenv("SMARTLING_USER_ID", "env-user")
	t.Setenv("SMARTLING_SECRET", "env-secret")
	t.Setenv("SMARTLING_PROJECT_ID", "env-project")

	dir := writeTempConfig(t, "")

	cfg, err := BuildConfigFromFlags(Params{Directory: dir})
	if err != nil {
		t.Fatalf("BuildConfigFromFlags: %v", err)
	}

	if cfg.Sources.UserID != SourceEnv {
		t.Errorf("UserID source = %q, want %q", cfg.Sources.UserID, SourceEnv)
	}
	if cfg.Sources.ProjectID != SourceEnv {
		t.Errorf("ProjectID source = %q, want %q", cfg.Sources.ProjectID, SourceEnv)
	}
	if cfg.Sources.AccountID != SourceDefault {
		t.Errorf("AccountID source = %q, want %q", cfg.Sources.AccountID, SourceDefault)
	}
}

func TestBuildConfigFromFlags_SourcesFromFlags(t *testing.T) {
	t.Setenv("SMARTLING_USER_ID", "env-user")
	t.Setenv("SMARTLING_SECRET", "env-secret")
	t.Setenv("SMARTLING_PROJECT_ID", "env-project")

	dir := writeTempConfig(t, `
user_id: "file-user"
secret: "file-secret"
account_id: "file-account"
project_id: "file-project"
`)

	cfg, err := BuildConfigFromFlags(Params{
		Directory: dir,
		User:      "flag-user",
		Secret:    "flag-secret",
		Account:   "flag-account",
		Project:   "flag-project",
	})
	if err != nil {
		t.Fatalf("BuildConfigFromFlags: %v", err)
	}

	if cfg.Sources.UserID != SourceFlag {
		t.Errorf("UserID source = %q, want %q", cfg.Sources.UserID, SourceFlag)
	}
	if cfg.Sources.AccountID != SourceFlag {
		t.Errorf("AccountID source = %q, want %q", cfg.Sources.AccountID, SourceFlag)
	}
	if cfg.Sources.ProjectID != SourceFlag {
		t.Errorf("ProjectID source = %q, want %q", cfg.Sources.ProjectID, SourceFlag)
	}

	if cfg.UserID != "flag-user" {
		t.Errorf("UserID = %q, want %q", cfg.UserID, "flag-user")
	}
}
