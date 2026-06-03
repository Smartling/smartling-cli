package list

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestJobsList(t *testing.T) {
	relativeDir := "../../bin/"
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
			name:              "help flag shows command description",
			args:              []string{"jobs", "list", "--help"},
			expectedOutputs:   []string{"List", "name", "output"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:              "list jobs in project",
			args:              []string{"jobs", "list"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:              "list with name filter",
			args:              []string{"jobs", "list", "--name", jobName},
			expectedOutputs:   []string{jobName},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:              "table output format shows column headers",
			args:              []string{"jobs", "--output", "table", "list", "--name", jobName},
			expectedOutputs:   []string{"TRANSLATION JOB UID", "NAME", "STATUS"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:              "json output format contains field names",
			args:              []string{"jobs", "--output", "json", "list", "--name", jobName},
			expectedOutputs:   []string{"translationJobUid", "jobName", jobName},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:            "invalid output format is rejected",
			args:            []string{"jobs", "--output", "invalid", "list"},
			expectedOutputs: []string{"invalid output"},
			wantErr:         true,
		},
		{
			name:            "extra positional arg is rejected",
			args:            []string{"jobs", "list", "unexpected-arg"},
			expectedOutputs: []string{"unknown command"},
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
