package cmd

import (
	"os/exec"
	"strings"
	"testing"
)

func TestProjectInfo_verbose(t *testing.T) {
	subCommands := []string{"projects", "info"}
	tests := []struct {
		name             string
		args             []string
		expectedOutput   string
		unexpectedOutput string
		wantErr          bool
	}{
		{
			name: "debug output with verbose flag",
			args: append(subCommands,
				[]string{
					"-vv",
				}...),
			expectedOutput: "DEBUG",
			wantErr:        false,
		},
		{
			name: "debug output without verbose flag",
			args: append(subCommands,
				[]string{
					"-v",
				}...),
			unexpectedOutput: "DEBUG",
			wantErr:          false,
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
			if tt.expectedOutput != "" {
				if !strings.Contains(string(out), tt.expectedOutput) {
					t.Errorf("output: %s\nwithout expected: %s", string(out), tt.expectedOutput)
				}
			}
			if tt.unexpectedOutput != "" {
				if strings.Contains(string(out), tt.unexpectedOutput) {
					t.Errorf("output: %s\nwith unexpected: %s", string(out), tt.expectedOutput)
				}
			}
		})
	}
}
