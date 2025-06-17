package init

import (
	"os/exec"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	subCommands := []string{"init"}
	expectedOutputs := []string{"Smartling API V2.0", "User Identifier", "Token Secret",
		"Project ID", "Enter a value", "Connection is successful"}
	unexpectedOutputs := []string{"DEBUG", "ERROR"}

	out, err := exec.Command(
		"./../bin/smartling-cli", subCommands...).
		CombinedOutput()
	if err != nil {
		t.Fatalf("error: %v, output: %s", err, string(out))
	}
	if len(expectedOutputs) > 0 {
		for _, expectedOutput := range expectedOutputs {
			if !strings.Contains(string(out), expectedOutput) {
				t.Errorf("output: %s\nwithout expected: %s", string(out), expectedOutput)
			}
		}
	}
	if len(unexpectedOutputs) > 0 {
		for _, unexpectedOutput := range unexpectedOutputs {
			if strings.Contains(string(out), unexpectedOutput) {
				t.Errorf("output: %s\nwith unexpected: %s", string(out), unexpectedOutput)
			}
		}
	}
}
