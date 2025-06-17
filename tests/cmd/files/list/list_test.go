package list

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilesList(t *testing.T) {
	relativeDir := "../../bin/"
	absDir, err := filepath.Abs(relativeDir)
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	subCommands := []string{"files", "list"}
	tests := []struct {
		name              string
		args              []string
		expectedOutputs   []string
		unexpectedOutputs []string
		wantErr           bool
	}{
		{
			name:              "List of all files in project",
			args:              subCommands,
			expectedOutputs:   []string{"plainText", "website_menu.txt", "texts"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "List files by mask",
			args:              append(subCommands, "*.txt", "--short"),
			expectedOutputs:   []string{"website_menu.txt"},
			unexpectedOutputs: []string{"plainText", "DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "Custom output format",
			args:              append(subCommands, "*.txt", "--format", "{{.FileType}}\\t||\\t{{.FileURI}};\\n"),
			expectedOutputs:   []string{"plainText", "website_menu.txt", "||"},
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
