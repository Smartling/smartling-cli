package remove

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestJobsFilesRemove(t *testing.T) {
	absDir, err := filepath.Abs("../../../bin/")
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	subCommands := []string{"jobs", "files", "remove"}
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
			args:              append(subCommands, "--help"),
			expectedOutputs:   []string{"Detach files from an existing translation job", "file"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:            "missing --file is rejected",
			args:            append(subCommands, jobName),
			expectedOutputs: []string{`required flag(s) "file" not set`},
			wantErr:         true,
		},
		{
			name:              "remove files from a job",
			args:              append(subCommands, jobName, "--file", "test.json"),
			expectedOutputs:   []string{"Files removed"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
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
