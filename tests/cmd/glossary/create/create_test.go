package create

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGlossaryCreate(t *testing.T) {
	relativeDir := "../../bin/"
	absDir, err := filepath.Abs(relativeDir)
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}
	subCommands := []string{"glossary", "create"}
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
			expectedOutputs:   []string{"Create a new glossary", "glossaryName", "locale"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:            "missing glossary name argument",
			args:            subCommands,
			expectedOutputs: []string{"accepts 1 arg"},
			wantErr:         true,
		},
		{
			name:            "missing locale flag returns validation error",
			args:            append(subCommands, "test-integration-glossary-create"),
			expectedOutputs: []string{"LocaleIDs"},
			wantErr:         true,
		},
		{
			name:              "create glossary with locale",
			args:              append(subCommands, "test-integration-glossary-create", "--locale", "en-US"),
			expectedOutputs:   []string{"Glossary UID:", "Glossary name:", "test-integration-glossary-create"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "create glossary with description",
			args:              append(subCommands, "test-integration-glossary-create-desc", "--locale", "en-US", "--description", "Integration test glossary"),
			expectedOutputs:   []string{"Glossary UID:", "Glossary name:", "test-integration-glossary-create-desc"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "create glossary with multiple locales",
			args:              append(subCommands, "test-integration-glossary-create-locales", "--locale", "en-US", "--locale", "es-ES"),
			expectedOutputs:   []string{"Glossary UID:", "Glossary name:", "test-integration-glossary-create-locales"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "table output format shows column headers",
			args:              []string{"glossary", "--output", "table", "create", "test-integration-glossary-create-table", "--locale", "en-US"},
			expectedOutputs:   []string{"GLOSSARY UID", "GLOSSARY NAME", "test-integration-glossary-create-table"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "json output format",
			args:              []string{"glossary", "--output", "json", "create", "test-integration-glossary-create-json", "--locale", "en-US"},
			expectedOutputs:   []string{"glossaryUid", "glossaryName", "test-integration-glossary-create-json"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "create glossary with fallback locale",
			args:              append(subCommands, "test-integration-glossary-create-fallback", "--locale", "en-US", "--locale", "es-ES", "--fallback-locale", "es:es-ES"),
			expectedOutputs:   []string{"Glossary UID:", "Glossary name:", "test-integration-glossary-create-fallback"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:            "invalid fallback locale format",
			args:            append(subCommands, "test-integration-glossary-create-bad-fallback", "--locale", "en-US", "--fallback-locale", "no-colon-here"),
			expectedOutputs: []string{"fallback locale must use format"},
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
