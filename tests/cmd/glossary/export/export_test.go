package export

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGlossaryExport(t *testing.T) {
	relativeDir := "../../bin/"
	absDir, err := filepath.Abs(relativeDir)
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	subCommands := []string{"glossary", "export"}
	// glossary created by the create integration tests
	testGlossary := "test-integration-glossary-create"
	tests := []struct {
		name              string
		args              []string
		expectedOutputs   []string
		unexpectedOutputs []string
		wantErr           bool
		cleanup           func()
	}{
		{
			name:              "help flag shows command description",
			args:              append(subCommands, "--help"),
			expectedOutputs:   []string{"Export glossary entries", "file-type", "locale"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:            "missing --file-type required flag",
			args:            append(subCommands, testGlossary),
			expectedOutputs: []string{"file-type"},
			wantErr:         true,
		},
		{
			name:            "tbx file type without --tbx-version is rejected",
			args:            append(subCommands, testGlossary, "--file-type", "tbx"),
			expectedOutputs: []string{"tbx-version"},
			wantErr:         true,
		},
		{
			name:            "unsupported file type is rejected",
			args:            append(subCommands, testGlossary, "--file-type", "docx"),
			expectedOutputs: []string{"unsupported file type"},
			wantErr:         true,
		},
		{
			name:              "export to csv",
			args:              append(subCommands, testGlossary, "--file-type", "csv"),
			expectedOutputs:   []string{"Glossary UID:", "Output file:", "File type:", "Bytes written:"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name: "export to explicit output file",
			args: append(subCommands, testGlossary, "test-export-output.csv", "--file-type", "csv"),
			expectedOutputs: []string{
				"Glossary UID:", "Output file:", "test-export-output.csv", "File type:", "Bytes written:",
			},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
			cleanup: func() {
				_ = os.Remove(filepath.Join(absDir, "test-export-output.csv"))
			},
		},
		{
			name:              "table output format shows column headers",
			args:              []string{"glossary", "--output", "table", "export", testGlossary, "--file-type", "csv"},
			expectedOutputs:   []string{"GLOSSARY UID", "OUTPUT FILE", "FILE TYPE", "BYTES WRITTEN"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "json output format contains field names",
			args:              []string{"glossary", "--output", "json", "export", testGlossary, "--file-type", "csv"},
			expectedOutputs:   []string{"glossaryUid", "outFile", "fileType", "bytesWritten"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "export xlsx format",
			args:              append(subCommands, testGlossary, "--file-type", "xlsx"),
			expectedOutputs:   []string{"Glossary UID:", "File type:", "xlsx"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "export with locale filter",
			args:              append(subCommands, testGlossary, "--file-type", "csv", "--locale", "en-US"),
			expectedOutputs:   []string{"Glossary UID:", "File type:", "csv"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCmd := exec.Command("./smartling-cli", tt.args...)
			testCmd.Dir = absDir
			out, err := testCmd.CombinedOutput()
			if tt.cleanup != nil {
				defer tt.cleanup()
			}
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
