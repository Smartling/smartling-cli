package status

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

func TestNewStatusCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	filesSrv := srvmocks.NewMockService(t)
	params := files.StatusParams{
		URI:       "https://example.com:8080/path/to/resource?search=a",
		Directory: "text",
		Format:    "txt",
	}
	filesSrv.On("RunStatus", params).Run(func(args mock.Arguments) {
		if _, err := fmt.Fprintf(buf, "RunStatus was called with %d args\n", len(args)); err != nil {
			t.Fatal(err)
		}
		if _, err := fmt.Fprintf(buf, "params: %v\n", args[0]); err != nil {
			t.Fatal(err)
		}
	}).Return(nil)

	initializer := cmdmocks.NewMockSrvInitializer(t)
	initializer.On("InitFilesSrv").Return(filesSrv, nil)

	cmd := NewStatusCmd(initializer)

	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		params.URI,
		"--format", params.Format,
		"--directory", params.Directory,
	})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() returned an error: %v", err)
	}

	output := buf.String()
	expected := fmt.Sprintf(`RunStatus was called with 1 args
params: %v
`, params)

	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, got %q", expected, output)
	}
}
