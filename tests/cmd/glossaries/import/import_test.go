package glimport

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// exportCSV runs `glossary export` and returns the base filename of the
// exported CSV. This guarantees the column headers match exactly what the
// Smartling import API expects, without hard-coding locale display names.
func exportCSV(t *testing.T, absDir, glossaryName, outName string) string {
	t.Helper()
	outPath := filepath.Join(absDir, outName)
	cmd := exec.Command("./smartling-cli", "glossaries", "export", glossaryName, outName, "--file-type", "csv")
	cmd.Dir = absDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("export for import test: %v\n%s", err, out)
	}
	t.Cleanup(func() { os.Remove(outPath) })
	return outName
}

func TestGlossaryImport(t *testing.T) {
	relativeDir := "../../bin/"
	absDir, err := filepath.Abs(relativeDir)
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	subCommands := []string{"glossaries", "import"}
	// glossary created by the create integration tests
	testGlossary := "test-integration-glossary-create"
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
			expectedOutputs:   []string{"Upload a CSV, XLSX, or TBX file", "glossaryUID", "inFile"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:            "missing both arguments",
			args:            subCommands,
			expectedOutputs: []string{"accepts 2 arg"},
			wantErr:         true,
		},
		{
			name:            "missing inFile argument",
			args:            append(subCommands, testGlossary),
			expectedOutputs: []string{"accepts 2 arg"},
			wantErr:         true,
		},
		{
			name: "unknown extension without --media-type fails validation",
			// Validate() catches empty MediaType before any file read or API call.
			args:            append(subCommands, testGlossary, "terms.dat"),
			expectedOutputs: []string{"ImportFile.MediaType", "cannot be empty"},
			wantErr:         true,
		},
		{
			name: "import csv file",
			args: func() []string {
				f := exportCSV(t, absDir, testGlossary, "test-import.csv")
				return append(subCommands, testGlossary, f)
			}(),
			expectedOutputs:   []string{"Glossary UID:", "Import UID:", "Import status:", "Source file:"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name: "table output format shows column headers",
			args: func() []string {
				f := exportCSV(t, absDir, testGlossary, "test-import-table.csv")
				return []string{"glossaries", "--output", "table", "import", testGlossary, f}
			}(),
			expectedOutputs:   []string{"GLOSSARY UID", "IMPORT UID", "STATUS", "SOURCE FILE"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name: "json output format contains field names",
			args: func() []string {
				f := exportCSV(t, absDir, testGlossary, "test-import-json.csv")
				return []string{"glossaries", "--output", "json", "import", testGlossary, f}
			}(),
			expectedOutputs:   []string{"glossaryUid", "importUid", "importStatus", "sourceFile"},
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
