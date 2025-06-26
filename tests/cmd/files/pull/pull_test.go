package pull

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilesPull(t *testing.T) {
	filename := "website_menu.txt"
	mdDilename := "readme.md"

	relativeDir := "../../bin/"
	absDir, err := filepath.Abs(relativeDir)
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}

	before, after := preparation(t, filepath.Join(absDir, filename), filepath.Join(absDir, mdDilename))
	before()
	defer after()

	subCommands := []string{"files", "pull"}
	tests := []struct {
		name              string
		args              []string
		expectedOutputs   []string
		unexpectedOutputs []string
		wantErr           bool
	}{
		{
			name:              "Download unavailable files",
			args:              append(subCommands, "**test.xemel", "--source"),
			expectedOutputs:   []string{"ERROR", "failed to run pull", "no files found on the remote server"},
			unexpectedOutputs: []string{"DEBUG"},
			wantErr:           false,
		},
		{
			name:              "Download translated files files",
			args:              append(subCommands, "*.txt", "-l", "uk-UA"),
			expectedOutputs:   []string{"downloaded", "txt"},
			unexpectedOutputs: []string{"ERROR", "DEBUG"},
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

func preparation(t *testing.T, filename, mdFilename string) (func(), func()) {
	before := func() {
		f, err := os.Create(filename)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
		}()
		if _, err := f.WriteString("Home\n"); err != nil {
			t.Fatal(err)
		}
		if _, err := f.WriteString("About us\n"); err != nil {
			t.Fatal(err)
		}
		if _, err := f.WriteString("News\n"); err != nil {
			t.Fatal(err)
		}

		mdFile, err := os.Create(mdFilename)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := mdFile.Close(); err != nil {
				t.Fatal(err)
			}
		}()
		if _, err := mdFile.WriteString("Readme\n"); err != nil {
			t.Fatal(err)
		}
	}
	after := func() {
		if err := os.Remove(filename); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(mdFilename); err != nil {
			t.Fatal(err)
		}
	}
	return before, after
}
