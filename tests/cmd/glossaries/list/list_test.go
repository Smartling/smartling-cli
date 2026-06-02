package list

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGlossaryList(t *testing.T) {
	relativeDir := "../../bin/"
	absDir, err := filepath.Abs(relativeDir)
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	subCommands := []string{"glossaries", "list"}
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
			expectedOutputs:   []string{"List glossaries", "name", "output"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:            "name flag with extra positional arg is rejected",
			args:            append(subCommands, "unexpected-arg"),
			expectedOutputs: []string{"unknown command"},
			wantErr:         true,
		},
		{
			name:              "list all glossaries",
			args:              subCommands,
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "list with name filter",
			args:              append(subCommands, "--name", "test-integration-glossary-create"),
			expectedOutputs:   []string{"test-integration-glossary-create"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "table output format shows column headers",
			args:              []string{"glossaries", "--output", "table", "list"},
			expectedOutputs:   []string{"GLOSSARY UID", "NAME", "DESCRIPTION", "LOCALES"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "json output format contains field names",
			args:              []string{"glossaries", "--output", "json", "list", "--name", "test-integration-glossary-create"},
			expectedOutputs:   []string{"GlossaryUID", "Name", "test-integration-glossary-create"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:            "invalid output format is rejected",
			args:            []string{"glossaries", "--output", "invalid", "list"},
			expectedOutputs: []string{"invalid output"},
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
