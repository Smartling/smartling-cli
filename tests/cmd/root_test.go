package cmd

import (
	"os/exec"
	"strings"
	"testing"
)

func TestRoot(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "config flag",
			args:       []string{"--config", "test.yaml"},
			wantOutput: "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := exec.Command(
				"./bin/smartling-cli", tt.args...).
				CombinedOutput()
			if err != nil {
				t.Fatalf("error: %v, output: %s", err, string(out))
			}
			if !strings.Contains(string(out), tt.wantOutput) {
				t.Errorf("unexpected output: %s", string(out))
			}
		})
	}
}
