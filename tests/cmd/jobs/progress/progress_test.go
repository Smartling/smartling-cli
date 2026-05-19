package progress

import (
	"os/exec"
	"strings"
	"testing"
)

func TestProgress(t *testing.T) {
	subCommands := []string{"jobs", "progress", "CLI uploads"}
	expectedOutputs := []string{"Total word count", "Percent complete"}
	unexpectedOutputs := []string{"DEBUG", "ERROR"}

	out, err := exec.Command(
		"./../../bin/smartling-cli", subCommands...).
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

func TestProgressJsonFormat(t *testing.T) {
	subCommands := []string{"jobs", "progress", "CLI uploads", "--output", "json"}
	expectedOutputs := []string{"totalWordCount", "percentComplete", "workflowProgressReportList"}
	unexpectedOutputs := []string{"DEBUG", "ERROR"}

	out, err := exec.Command(
		"./../../bin/smartling-cli", subCommands...).
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
