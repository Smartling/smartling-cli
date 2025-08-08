package importcmd

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

func TestNewImportCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	filesSrv := srvmocks.NewMockService(t)
	params := files.ImportParams{
		URI:             "https://example.com:8080/path/to/resource?search=a",
		File:            "a.txt",
		Locale:          "en-US",
		FileType:        "txt",
		PostTranslation: false,
		Overwrite:       true,
	}
	filesSrv.On("RunImport", params).Run(func(args mock.Arguments) {
		fmt.Fprintf(buf, "RunImport was called with %d args\n", len(args))
		fmt.Fprintf(buf, "params: %v\n", args[0])
	}).Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitFilesSrv").Return(filesSrv, nil)

	cmd := NewImportCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		params.URI,
		params.File,
		params.Locale,
		"--type", params.FileType,
		"--overwrite", fmt.Sprintf("%v", params.Overwrite),
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := fmt.Sprintf(`RunImport was called with 1 args
params: %v
`, params)

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
