package push

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	cmdmocks "github.com/Smartling/smartling-cli/cmd/files/mocks"
	"github.com/Smartling/smartling-cli/services/files"
	srvmocks "github.com/Smartling/smartling-cli/services/files/mocks"

	"github.com/stretchr/testify/mock"
)

func TestNewPushCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	filesSrv := srvmocks.NewMockService(t)
	params := files.PushParams{
		URI:        "https://example.com:8080/path/to/resource?search=a",
		File:       "01.txt",
		Branch:     "testing",
		Locales:    []string{"en-US", "fr-FR"},
		Authorize:  true,
		FileType:   "text",
		Directory:  ".",
		Directives: []string{"01", "05"},
	}
	filesSrv.On("RunPush", params).Run(func(args mock.Arguments) {
		fmt.Fprintln(buf, fmt.Sprintf("RunPush was called with %d args", len(args)))
		fmt.Fprintln(buf, fmt.Sprintf("params: %v", args[0]))
	}).Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitFilesSrv").Return(filesSrv, nil)

	cmd := NewPushCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		params.File,
		params.URI,
		"--authorize",
		"--locale", params.Locales[0],
		"--locale", params.Locales[1],
		"--branch", params.Branch,
		"--type", params.FileType,
		"--directive", params.Directives[0],
		"--directive", params.Directives[1],
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := fmt.Sprintf(`RunPush was called with 1 args
params: %v
`, params)

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
