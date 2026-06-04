package locales

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestJobsLocalesGroup(t *testing.T) {
	absDir, err := filepath.Abs("../../bin/")
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	tests := []struct {
		name              string
		args              []string
		expectedOutputs   []string
		unexpectedOutputs []string
		wantErr           bool
	}{
		{
			name:              "bare group shows help with subcommands",
			args:              []string{"jobs", "locales"},
			expectedOutputs:   []string{"Add or remove target locales", "Available Commands:", "add", "remove"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			// Regression: an invalid subcommand used to be swallowed as positional
			// args and silently succeed; it must now surface the usage/command list.
			name:              "invalid subcommand surfaces usage instead of silently succeeding",
			args:              []string{"jobs", "locales", "ad", "fr-FR"},
			expectedOutputs:   []string{"Available Commands:", "add", "remove"},
			unexpectedOutputs: []string{"added", "removed", "DEBUG", "ERROR"},
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCmd := exec.Command("./smartling-cli", tt.args...)
			testCmd.Dir = absDir
			out, err := testCmd.CombinedOutput()
			if !tt.wantErr && err != nil {
				t.Fatalf("error: %v, output: %s", err, string(out))
			}
			for _, expectedOutput := range tt.expectedOutputs {
				if !strings.Contains(string(out), expectedOutput) {
					t.Errorf("output: %s\nwithout expected: %s", string(out), expectedOutput)
				}
			}
			for _, unexpectedOutput := range tt.unexpectedOutputs {
				if strings.Contains(string(out), unexpectedOutput) {
					t.Errorf("output: %s\nwith unexpected: %s", string(out), unexpectedOutput)
				}
			}
		})
	}
}
