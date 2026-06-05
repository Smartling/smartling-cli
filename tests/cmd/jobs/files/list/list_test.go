package list

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestJobsFiles(t *testing.T) {
	relativeDir := "../../../bin/"
	absDir, err := filepath.Abs(relativeDir)
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}

	jobName := "CLI uploads"

	tests := []struct {
		name              string
		args              []string
		expectedOutputs   []string
		unexpectedOutputs []string
		wantErr           bool
	}{
		{
			name:              "group help lists subcommands",
			args:              []string{"jobs", "files", "--help"},
			expectedOutputs:   []string{"Add, remove, or list", "add", "list", "remove"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:              "list files of a job by name includes the added file",
			args:              []string{"jobs", "files", "list", jobName},
			expectedOutputs:   []string{"test.json"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:              "table output format shows column headers",
			args:              []string{"jobs", "--output", "table", "files", "list", jobName},
			expectedOutputs:   []string{"FILE URI", "LOCALES"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:              "json output format contains field names",
			args:              []string{"jobs", "--output", "json", "files", "list", jobName},
			expectedOutputs:   []string{"fileUri", "localeIds"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:            "invalid output format is rejected",
			args:            []string{"jobs", "--output", "invalid", "files", "list", jobName},
			expectedOutputs: []string{"invalid output"},
			wantErr:         true,
		},
		{
			name:            "missing positional arg is rejected",
			args:            []string{"jobs", "files", "list"},
			expectedOutputs: []string{"accepts 1 arg(s), received 0"},
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
