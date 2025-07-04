package detect

import (
	"bytes"
	"strings"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/mt/mocks"
	"github.com/Smartling/smartling-cli/services/helpers/rlog"
	srv "github.com/Smartling/smartling-cli/services/mt"
)

func TestNewDetectCmdNoArgs(t *testing.T) {
	rlog.Init()
	buf := new(bytes.Buffer)
	params := srv.DetectParams{
		InputDirectory: "/inputs/",
		FileType:       "text",
		FileOrPattern:  "*.txt",
	}

	initializer := cmdmocks.NewMockSrvInitializer(t)

	cmd := NewDetectCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		"--input-directory", params.InputDirectory,
		"--type", params.FileType,
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() expected with error")
	}

	output := buf.String()
	expected := "wrong argument quantity"

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
