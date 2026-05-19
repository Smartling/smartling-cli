package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestProjectInfo_verbose(t *testing.T) {
	subCommands := []string{"projects", "info"}
	tests := []struct {
		name              string
		args              []string
		expectedOutputs   []string
		unexpectedOutputs []string
		wantErr           bool
	}{
		{
			name:              "debug output with debug verbose flag",
			args:              append(subCommands, "-vv"),
			expectedOutputs:   []string{"DEBUG", "ID", "ACCOUNT", "NAME", "LOCALE", "STATUS"},
			unexpectedOutputs: []string{"ERROR"},
			wantErr:           false,
		},
		{
			name:              "debug output without debug verbose flag",
			args:              append(subCommands, "-v"),
			expectedOutputs:   []string{"ID", "ACCOUNT", "NAME", "LOCALE", "STATUS"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := exec.Command(
				"./../../bin/smartling-cli", tt.args...).
				CombinedOutput()
			if err != nil {
				t.Fatalf("error: %v, output: %s", err, string(out))
			}
			if len(tt.expectedOutputs) > 0 {
				for _, expectedOutput := range tt.expectedOutputs {
					if !strings.Contains(string(out), expectedOutput) {
						t.Errorf("output: %s\nwithout expected: %s", string(out), expectedOutput)
					}
				}
			}
			if len(tt.unexpectedOutputs) > 0 {
				for _, unexpectedOutput := range tt.unexpectedOutputs {
					if strings.Contains(string(out), unexpectedOutput) {
						t.Errorf("output: %s\nwith unexpected: %s", string(out), unexpectedOutput)
					}
				}
			}
		})
	}
}

func TestProjectInfo_ShowConfig(t *testing.T) {
	// `go test` runs without a TTY, so --show-config prints the banner to
	// stdout and proceeds without prompting.
	cmd := exec.Command("./../../bin/smartling-cli", "projects", "info", "--show-config")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("error: %v\nstderr:\n%s\nstdout:\n%s", err, stderr.String(), stdout.String())
	}

	out := stdout.String()
	for _, want := range []string{
		"> Smartling CLI configuration:",
		">   Config file:",
		">   User ID:",
		">   Account ID:",
		">   Project ID:",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout missing banner line %q\nfull stdout:\n%s", want, out)
		}
	}

	if strings.Contains(stderr.String(), "Continue?") {
		t.Errorf("non-TTY run must not prompt:\n%s", stderr.String())
	}

	for _, want := range []string{"ID", "ACCOUNT", "NAME", "LOCALE", "STATUS", "USER", "CONFIG", "SOURCES"} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout missing expected info row label %q\nfull stdout:\n%s", want, out)
		}
	}
}
