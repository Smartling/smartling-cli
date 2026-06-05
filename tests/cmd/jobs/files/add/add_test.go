package add

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestJobsFilesAdd(t *testing.T) {
	absDir, err := filepath.Abs("../../../bin/")
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	subCommands := []string{"jobs", "files", "add"}
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
			expectedOutputs:   []string{"Attach files to an existing translation job", "file", "target-locale"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
		},
		{
			name:            "missing --file is rejected",
			args:            append(subCommands, jobName),
			expectedOutputs: []string{`required flag(s) "file" not set`},
			wantErr:         true,
		},
		{
			// `jobs files add` attaches an existing project file by URI; it does
			// not upload. Push the local fixture (project-only) first so the exact
			// --file test.json pattern resolves to a single file.
			name:              "push fixture file to project",
			args:              []string{"files", "push", "test.json", "test.json", "--nojob"},
			expectedOutputs:   []string{"test.json"},
			unexpectedOutputs: []string{"ERROR"},
		},
		{
			name:              "add file to a job",
			args:              append(subCommands, jobName, "--file", "test.json"),
			expectedOutputs:   []string{"Files added", "test.json"},
			unexpectedOutputs: []string{"No files matched", "DEBUG", "ERROR"},
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
