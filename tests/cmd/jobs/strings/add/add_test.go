package add

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestJobsStringsAdd(t *testing.T) {
	absDir, err := filepath.Abs("../../../bin/")
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	subCommands := []string{"jobs", "strings", "add"}
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
			expectedOutputs:   []string{"Assign strings (by hashcode) to an existing translation job", "hashcode", "target-locale"},
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
			name:              "add strings to a job",
			args:              append(subCommands, "CLI uploads", "--hashcode", "ca51a04da69cf64dce022bb4f146c962"),
			expectedOutputs:   []string{"added"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "json output format contains field names",
			args:              []string{"jobs", "--output", "json", "strings", "add", "CLI uploads", "--hashcode", "ca51a04da69cf64dce022bb4f146c962"},
			expectedOutputs:   []string{"action", "translationJobUid", "hashcodes"},
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
