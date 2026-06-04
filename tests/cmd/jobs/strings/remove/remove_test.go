package remove

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestJobsStringsRemove(t *testing.T) {
	absDir, err := filepath.Abs("../../../bin/")
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	subCommands := []string{"jobs", "strings", "remove"}
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
			expectedOutputs:   []string{"Detach strings (by hashcode) from an existing translation job", "hashcode", "locale"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:            "missing job argument is rejected",
			args:            append(subCommands, "--hashcode", "h1"),
			expectedOutputs: []string{"accepts 1 arg(s)"},
			wantErr:         true,
		},
		{
			name:              "remove strings from a job",
			args:              append(subCommands, "CLI uploads", "--hashcode", "ca51a04da69cf64dce022bb4f146c962"),
			expectedOutputs:   []string{"removed"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
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
