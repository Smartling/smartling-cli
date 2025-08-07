package push

import (
	"bytes"
	"fmt"
	"log"
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
		Directives: map[string]string{"key1": "01", "key5": "05"},
	}
	filesSrv.On("RunPush", mock.Anything, params).Run(func(args mock.Arguments) {
		if _, err := fmt.Fprintf(buf, "RunPush was called with %d args\n", len(args)); err != nil {
			log.Panic(err)
		}
		if _, err := fmt.Fprintf(buf, "params: %v\n", args[1]); err != nil {
			log.Panic(err)
		}
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
		"--directive", "key1=01",
		"--directive", "key5=05",
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := fmt.Sprintf(`RunPush was called with 2 args
params: %v
`, params)

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
