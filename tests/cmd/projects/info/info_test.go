package cmd

import (
	"os/exec"
	"strings"
	"testing"
)

func TestProjectInfo_verbose(t *testing.T) {
	subCommands := []string{"projects", "info"}
	tests := []struct {
		name              string
		args              []string
		expectedOutputs   []string
		unexpectedOutputs []string
		wantErr           bool
	}{
		{
			name: "debug output with verbose flag",
			args: append(subCommands,
				[]string{
					"-vv",
				}...),
			expectedOutputs:   []string{"DEBUG", "ID", "ACCOUNT", "NAME", "LOCALE", "STATUS"},
			unexpectedOutputs: []string{"ERROR"},
			wantErr:           false,
		},
		{
			name: "debug output without verbose flag",
			args: append(subCommands,
				[]string{
					"-v",
				}...),
			expectedOutputs:   []string{"ID", "ACCOUNT", "NAME", "LOCALE", "STATUS"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := exec.Command(
				"./../../bin/smartling-cli", tt.args...).
				CombinedOutput()
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
