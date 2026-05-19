package projects

import (
	"bytes"
	"strings"
	"testing"

	projectconfig "github.com/Smartling/smartling-cli/services/projects/config"
)

func testExtendedConfig() projectconfig.Extended {
	return projectconfig.Extended{
		ProjectID:  "p-789",
		AccountUID: "a-456",
		Name:       "Acme Localization",
		Locale:     "en-US: English (United States)",
		Status:     "active",
		UserID:     "u-123",
		ConfigFile: "/tmp/smartling.yml",
		Sources:    "project=flag  account=config  user=env",
	}
}

func TestRenderPlain_AllFieldsSet(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderPlain(&buf, testExtendedConfig()); err != nil {
		t.Fatalf("RenderPlain: %v", err)
	}
	out := buf.String()

	for _, want := range []string{
		"Smartling CLI configuration:",
		"Config file: /tmp/smartling.yml",
		"User:     u-123",
		"Account:  a-456",
		"Project:  p-789",
		"Project Name:  Acme Localization",
		"Locale:  en-US: English (United States)",
		"Status:  active",
		"Sources:  project=flag  account=config  user=env",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("banner missing %q\nfull output:\n%s", want, out)
		}
	}
}

func TestRenderPlain_EmptyFields(t *testing.T) {
	cfg := projectconfig.Extended{ConfigFile: "/tmp/smartling.yml"}

	var buf bytes.Buffer
	if err := RenderPlain(&buf, cfg); err != nil {
		t.Fatalf("RenderPlain: %v", err)
	}
	out := buf.String()

	for _, want := range []string{
		"User:     ",
		"Account:  ",
		"Project:  ",
		"Project Name:  ",
		"Locale:  ",
		"Status:  ",
		"Sources:  ",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("banner missing %q for empty field:\n%s", want, out)
		}
	}
}

func TestRenderPlain_LinePrefixes(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderPlain(&buf, testExtendedConfig()); err != nil {
		t.Fatalf("RenderPlain: %v", err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 9 {
		t.Fatalf("expected 9 lines, got %d:\n%s", len(lines), buf.String())
	}
	for i, line := range lines {
		if !strings.HasPrefix(line, "> ") {
			t.Errorf("line %d does not start with '> ': %q", i, line)
		}
	}
}
