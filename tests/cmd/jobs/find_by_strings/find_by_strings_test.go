package findbystrings

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestJobsFindByStrings(t *testing.T) {
	absDir, err := filepath.Abs("../../bin/")
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	subCommands := []string{"jobs", "find-by-strings"}
	tests := []struct {
		name              string
		args              []string
		expectedOutputs   []string
		unexpectedOutputs []string
		wantErr           bool
	}{
		{
			name:              "help flag shows command description",
			args:              append(subCommands, "--help"),
			expectedOutputs:   []string{"Find the translation jobs", "hashcode", "locale"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:            "missing required hashcode is rejected",
			args:            subCommands,
			expectedOutputs: []string{"required flag(s)", "hashcode"},
			wantErr:         true,
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
