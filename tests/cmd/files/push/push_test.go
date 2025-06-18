package push

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilesPush(t *testing.T) {
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

	subCommands := []string{"files", "push"}
	tests := []struct {
		name              string
		args              []string
		expectedOutputs   []string
		unexpectedOutputs []string
		wantErr           bool
	}{
		{
			name:              "Simplest one-file upload",
			args:              append(subCommands, filename),
			expectedOutputs:   []string{filename, "(plaintext)", "strings", "words"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "One-file upload with URI",
			args:              append(subCommands, filename, "/texts/"+filename),
			expectedOutputs:   []string{filename, "(plaintext)", "strings", "words"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "Override file type",
			args:              append(subCommands, filename, "--type", "plaintext"),
			expectedOutputs:   []string{filename, "(plaintext)", "strings", "words"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "Upload files by mask",
			args:              append(subCommands, "../../cmd/bin/**.txt"),
			expectedOutputs:   []string{filename, "(plaintext)", "strings", "words"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "Branching versioning",
			args:              append(subCommands, "../../cmd/bin/**.txt", "-b", "testing"),
			expectedOutputs:   []string{filename, "(plaintext)", "strings", "words"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
		{
			name:              "Branching @auto",
			args:              append(subCommands, "../../cmd/bin/**.txt", "-b", "@auto"),
			expectedOutputs:   []string{filename, "(plaintext)", "strings", "words"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
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

func TestFilesPushWithDirective(t *testing.T) {
	filename := "custom.json"

	relativeDir := "../../bin/"
	absDir, err := filepath.Abs(relativeDir)
	if err != nil {
		t.Fatalf("Failed to get abs path: %v", err)
	}

	before, after := prepareJson(t, filepath.Join(absDir, filename))
	before()
	defer after()

	subCommands := []string{"files", "push"}
	directiveValue := `translate_paths={"path":"*/string","key":"{*}/string","instruction":"*/instruction"}`
	tests := []struct {
		name              string
		command           string
		args              []string
		expectedOutputs   []string
		unexpectedOutputs []string
		wantErr           bool
	}{
		{
			name: "Setting Smartling directives to customize parsing (string capture) behavior",
			args: append(subCommands, "custom.json",
				"--directive", directiveValue,
			),
			command:           `./smartling-cli  files push custom.json --directive=` + "`translate_paths=\"{\\\"path\\\":\\\"*/string\\\",\\\"key\\\":\\\"{*}/string\\\",\\\"instruction\\\":\\\"*/instruction\\\"}\"`",
			expectedOutputs:   []string{filename, "(json)", "strings", "words"},
			unexpectedOutputs: []string{"DEBUG", "ERROR"},
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCmd := exec.Command("bash", "-c", tt.command)
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

func prepareJson(t *testing.T, filename string) (func(), func()) {
	content := `{"name": "Products","description": "More info"}`
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
		if _, err := f.WriteString(content); err != nil {
			t.Fatal(err)
		}
	}
	after := func() {
		if err := os.Remove(filename); err != nil {
			t.Fatal(err)
		}
	}
	return before, after
}
