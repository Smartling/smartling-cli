package view

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestJobsView(t *testing.T) {
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
			args:              []string{"jobs", "view", "--help"},
			expectedOutputs:   []string{"translation job", "output"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:              "view job by name",
			args:              []string{"jobs", "view", jobName},
			expectedOutputs:   []string{"Name:", "Status:", jobName},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:              "table output format shows fields",
			args:              []string{"jobs", "--output", "table", "view", jobName},
			expectedOutputs:   []string{"FIELD", "VALUE", "NAME", "STATUS"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:              "json output format contains field names",
			args:              []string{"jobs", "--output", "json", "view", jobName},
			expectedOutputs:   []string{"TranslationJobUID", "JobName", jobName},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:            "invalid output format is rejected",
			args:            []string{"jobs", "--output", "invalid", "view", jobName},
			expectedOutputs: []string{"invalid output"},
			wantErr:         true,
		},
		{
			name:            "missing positional arg is rejected",
			args:            []string{"jobs", "view"},
			expectedOutputs: []string{"wrong argument quantity"},
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
