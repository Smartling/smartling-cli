package cmd

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

func TestConfirmContinue(t *testing.T) {
	tests := []struct {
		name, input string
		want        bool
	}{
		{"y lowercase", "y\n", true},
		{"Y uppercase", "Y\n", true},
		{"yes", "yes\n", true},
		{"YES", "YES\n", true},
		{"n", "n\n", false},
		{"N", "N\n", false},
		{"no", "no\n", false},
		{"empty line", "\n", false},
		{"eof", "", false},
		{"whitespace around y", "  y  \n", true},
		{"whitespace around yes", "  yes  \n", true},
		{"unrelated text", "foo\n", false},
		{"y-prefixed word yeti", "yeti\n", false},
		{"y-prefixed word yikes", "yikes\n", false},
		{"yy", "yy\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := confirmContinue(strings.NewReader(tt.input)); got != tt.want {
				t.Errorf("confirmContinue(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestShowConfigAndMaybePrompt_NonTTY(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := showConfigAndMaybePrompt(testExtendedConfig(), &stdout, &stderr, strings.NewReader(""), false)

	if !got {
		t.Errorf("showConfigAndMaybePrompt should return true in non-TTY case, got false")
	}
	if !strings.Contains(stdout.String(), "Smartling CLI configuration:") {
		t.Errorf("stdout should contain banner, got: %q", stdout.String())
	}
	if stderr.String() != "" {
		t.Errorf("stderr should be empty (no prompt) in non-TTY case, got: %q", stderr.String())
	}
}

func TestShowConfigAndMaybePrompt_TTYConfirm(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := showConfigAndMaybePrompt(testExtendedConfig(), &stdout, &stderr, strings.NewReader("y\n"), true)

	if !got {
		t.Errorf("showConfigAndMaybePrompt should return true on 'y', got false")
	}
	if !strings.Contains(stdout.String(), "Smartling CLI configuration:") {
		t.Errorf("stdout should contain banner, got: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Continue?") {
		t.Errorf("stderr should contain prompt, got: %q", stderr.String())
	}
}

func TestShowConfigAndMaybePrompt_TTYAbort(t *testing.T) {
	var stdout, stderr bytes.Buffer
	got := showConfigAndMaybePrompt(testExtendedConfig(), &stdout, &stderr, strings.NewReader("n\n"), true)

	if got {
		t.Errorf("showConfigAndMaybePrompt should return false on 'n', got true")
	}
	if !strings.Contains(stderr.String(), "aborted") {
		t.Errorf("stderr should contain 'aborted' on abort, got: %q", stderr.String())
	}
}
