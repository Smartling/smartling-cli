package status

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilesStatus(t *testing.T) {
	relativeDir := "../../bin/"
	absDir, err := filepath.Abs(relativeDir)
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	subCommands := []string{"files", "status"}
	tests := []struct {
		name              string
		args              []string
		expectedOutputs   []string
		unexpectedOutputs []string
		wantErr           bool
	}{
		{
			name:              "Files status",
			args:              subCommands,
			expectedOutputs:   []string{"source", "website_menu.txt", "missing"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCmd := exec.Command("./smartling-cli", tt.args...)
			testCmd.Dir = absDir
			out, err := testCmd.CombinedOutput()
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
